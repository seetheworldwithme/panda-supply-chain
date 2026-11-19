package platform

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// CompressUploadImage 对照 PhotoCompressUtil：<800KiB 原样返回；否则循环缩放 85% 质量 85。
// 大图字节与 Thumbnailator 编码器存在差异，尺寸/大小契约一致（见差异报告）。
func CompressUploadImage(b []byte) ([]byte, error) {
	if len(b) < 800*1024 {
		return b, nil
	}
	im, format, e := image.Decode(bytes.NewReader(b))
	if e != nil {
		return b, nil
	}
	for len(b) > 800*1024 {
		bounds := im.Bounds()
		ww := int(float64(bounds.Dx()) * .85)
		hh := int(float64(bounds.Dy()) * .85)
		if ww < 1 {
			ww = 1
		}
		if hh < 1 {
			hh = 1
		}
		dst := image.NewNRGBA(image.Rect(0, 0, ww, hh))
		for y := 0; y < hh; y++ {
			for x := 0; x < ww; x++ {
				dst.Set(x, y, im.At(bounds.Min.X+x*bounds.Dx()/ww, bounds.Min.Y+y*bounds.Dy()/hh))
			}
		}
		var out bytes.Buffer
		switch format {
		case "png":
			e = png.Encode(&out, dst)
		default:
			e = jpeg.Encode(&out, dst, &jpeg.Options{Quality: 85})
		}
		if e != nil {
			return b, e
		}
		b = out.Bytes()
		im = dst
		if ww == 1 && hh == 1 {
			break
		}
	}
	return b, nil
}

// CompressFileRow 对照 compressedFile：写压缩副本并插 source_id 关联行，失败返回 0（Java catch 同值）。
func CompressFileRow(d *sql.DB, id int64) int64 {
	f, e := GetByID(d, "chain_common_file", id)
	if e != nil || f == nil {
		return 0
	}
	text := func(k string) string {
		s, _ := f[k].(string)
		return s
	}
	p := filepath.Join(text("file_path"), text("file_name"))
	b, e := os.ReadFile(p)
	if e != nil {
		return 0
	}
	b, e = CompressUploadImage(b)
	if e != nil {
		return 0
	}
	name := "compressed" + text("file_name")
	if os.WriteFile(filepath.Join(text("file_path"), name), b, 0644) != nil {
		return 0
	}
	vals := map[string]any{"file_name": name, "file_real_name": name, "file_path": f["file_path"], "f_type": f["f_type"], "owner_type": f["owner_type"], "owner_id": f["owner_id"], "relation_id": f["relation_id"], "relation_type": f["relation_type"], "effective": true, "association_id": f["association_id"], "source_id": id, "is_compressed_file": true}
	cid, e := InsertNoFill(d, "chain_common_file", vals)
	if e != nil {
		return 0
	}
	return cid
}

// RenderExcelUpload 对照 ExcelUtils.toImageBySpire：LibreOffice SinglePageSheets 每表一页 PDF，
// pdftoppm 转 PNG。与商业 Spire 的渲染引擎/字体差异已记录，不静默假成功。
func RenderExcelUpload(ctx context.Context, src, outdir string) ([]string, error) {
	bin := os.Getenv("PANDA_SOFFICE_BIN")
	if bin == "" {
		bin = "soffice"
	}
	raster := os.Getenv("PANDA_PDFTOPPM_BIN")
	if raster == "" {
		raster = "pdftoppm"
	}
	if e := os.MkdirAll(outdir, 0755); e != nil {
		return nil, e
	}
	tmp, e := os.MkdirTemp("", "panda-sheet-render-")
	if e != nil {
		return nil, e
	}
	defer os.RemoveAll(tmp)
	src, e = filepath.Abs(src)
	if e != nil {
		return nil, e
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	profile := (&url.URL{Scheme: "file", Path: filepath.Join(tmp, "profile")}).String()
	cmd := exec.CommandContext(ctx, bin, "-env:UserInstallation="+profile, "--headless", "--convert-to", `pdf:calc_pdf_Export:{"SinglePageSheets":{"type":"boolean","value":"true"}}`, "--outdir", tmp, src)
	if out, e := cmd.CombinedOutput(); e != nil {
		return nil, fmt.Errorf("spreadsheet renderer: %w: %s", e, out)
	}
	pdf := filepath.Join(tmp, strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))+".pdf")
	if _, e = os.Stat(pdf); e != nil {
		return nil, e
	}
	if out, e := exec.CommandContext(ctx, raster, "-png", "-r", "96", pdf, filepath.Join(tmp, "sheet")).CombinedOutput(); e != nil {
		return nil, fmt.Errorf("spreadsheet rasterizer: %w: %s", e, out)
	}
	pages, e := filepath.Glob(filepath.Join(tmp, "sheet-*.png"))
	if e != nil || len(pages) == 0 {
		return nil, fmt.Errorf("spreadsheet produced no preview")
	}
	sort.Slice(pages, func(i, j int) bool {
		number := func(p string) int {
			p = strings.TrimSuffix(filepath.Base(p), ".png")
			n, _ := strconv.Atoi(strings.TrimPrefix(p, "sheet-"))
			return n
		}
		return number(pages[i]) < number(pages[j])
	})
	names := []string{}
	stamp := time.Now().UnixMilli()
	for i, p := range pages {
		b, e := os.ReadFile(p)
		if e != nil {
			return nil, e
		}
		name := strconv.FormatInt(stamp+int64(i), 10) + ".png"
		if e = os.WriteFile(filepath.Join(outdir, name), b, 0644); e != nil {
			return nil, e
		}
		names = append(names, name)
	}
	return names, nil
}

// RenderDocumentPDF 对应服务器实际 Doc2PdfUtil 的文档预览，复用既有 LibreOffice。
func RenderDocumentPDF(ctx context.Context, src, target string) error {
	bin := os.Getenv("PANDA_SOFFICE_BIN")
	if bin == "" {
		bin = "soffice"
	}
	tmp, err := os.MkdirTemp("", "panda-doc-render-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	profile := (&url.URL{Scheme: "file", Path: filepath.Join(tmp, "profile")}).String()
	cmd := exec.CommandContext(ctx, bin, "-env:UserInstallation="+profile, "--headless", "--convert-to", "pdf:writer_pdf_Export", "--outdir", tmp, src)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("文档转换失败: %w", err)
	}
	data, err := os.ReadFile(filepath.Join(tmp, strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))+".pdf"))
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	return os.WriteFile(target, data, 0644)
}
