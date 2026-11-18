package commonrest

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// PhotoCompressUtil: <800KiB bytes unchanged; otherwise repeatedly shrink 85% at quality85.
// The Go encoder differs in large-image bytes from Thumbnailator, while the size/dimensions contract holds.
func compressUploadImage(b []byte) ([]byte, error) {
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

func (s *Service) compressFile(r *http.Request, id int64) int64 {
	f, e := platform.GetByID(s.pick(r), "chain_common_file", id)
	if e != nil || f == nil {
		return 0
	}
	p := filepath.Join(str(f, "file_path"), str(f, "file_name"))
	b, e := os.ReadFile(p)
	if e != nil {
		return 0
	}
	b, e = compressUploadImage(b)
	if e != nil {
		return 0
	}
	name := "compressed" + str(f, "file_name")
	if os.WriteFile(filepath.Join(str(f, "file_path"), name), b, 0644) != nil {
		return 0
	}
	vals := map[string]any{"file_name": name, "file_real_name": name, "file_path": f["file_path"], "f_type": f["f_type"], "owner_type": f["owner_type"], "owner_id": f["owner_id"], "relation_id": f["relation_id"], "relation_type": f["relation_type"], "effective": true, "association_id": f["association_id"], "source_id": id, "is_compressed_file": true}
	cid, e := platform.InsertNoFill(s.pick(r), "chain_common_file", vals)
	if e != nil {
		return 0
	}
	return cid
}

// ExcelUtils renders each worksheet into one preview. LibreOffice SinglePageSheets gives one
// PDF page per sheet, then pdftoppm produces PNGs. No Java runtime or MQ required. Rendering
// engine/font differences from proprietary Spire are documented, rather than hidden.
func renderExcelUpload(ctx context.Context, src, outdir string) ([]string, error) {
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
