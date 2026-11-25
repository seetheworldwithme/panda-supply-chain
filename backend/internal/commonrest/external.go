package commonrest

// UploadController 的本地文件/图片/分片上传链。工作流控制器源文件已全部注释，故不注册 workflow。
// Redis share token 用带过期内存表替代；preShareFIle 忠实保留读取字面键 "credit" 的 Java bug。
import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type shareValue struct {
	id      string
	expires time.Time
}

var shareTokens = struct {
	sync.Mutex
	values map[string]shareValue
}{values: map[string]shareValue{}}

func extUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = b[6]&15 | 64
	b[8] = b[8]&63 | 128
	s := hex.EncodeToString(b)
	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}
func storeShareToken(id string) string {
	key := extUUID()
	shareTokens.Lock()
	defer shareTokens.Unlock()
	shareTokens.values[key] = shareValue{id, time.Now().Add(24 * time.Hour)}
	return key
}
func shareTokenExists(key string) bool {
	shareTokens.Lock()
	defer shareTokens.Unlock()
	v, ok := shareTokens.values[key]
	if ok && time.Now().After(v.expires) {
		delete(shareTokens.values, key)
		return false
	}
	return ok
}
func externalRoot() string {
	if s := os.Getenv("PANDA_UPLOAD_DIR"); s != "" {
		return s
	}
	return "data/uploads"
}
func registerExternal(mux *http.ServeMux, s *Service) {
	reg := func(method, path string, h http.HandlerFunc) {
		if method == "POST" && (path == "/upload/getByRelId" || path == "/upload/getShareToken" || path == "/upload/preFIle" || path == "/upload/getError" || path == "/TenOcrController/getOcrResult") {
			next := h
			h = func(w http.ResponseWriter, r *http.Request) {
				b, err := io.ReadAll(io.LimitReader(r.Body, 32<<20))
				var m map[string]any
				if err != nil || json.Unmarshal(b, &m) != nil || m == nil {
					platform.Unknown(w)
					return
				}
				r.Body = io.NopCloser(bytes.NewReader(b))
				next(w, r)
			}
		}
		mux.HandleFunc(method+" /saas-common"+path, authGuard(h))
	}
	reg("POST", "/upload/uploadFiles", func(w http.ResponseWriter, r *http.Request) { s.extUpload(w, r, true) })
	reg("POST", "/upload/upload", func(w http.ResponseWriter, r *http.Request) { s.extUpload(w, r, false) })
	reg("POST", "/upload/preFIle", s.extPreview)
	reg("POST", "/upload/getShareToken", func(w http.ResponseWriter, r *http.Request) {
		m := bodyMap(r)
		if m["id"] == nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, storeShareToken(str(m, "id")))
	})
	reg("GET", "/upload/preShareFIle", func(w http.ResponseWriter, r *http.Request) { platform.Fail(w, 601, "链接已过期") })
	reg("POST", "/upload/readExcel", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil || r.MultipartForm == nil || len(r.MultipartForm.File["file"]) == 0 {
			platform.Unknown(w)
		}
	})
	reg("POST", "/upload/getByRelId", s.extByRel)
	reg("POST", "/upload/getImgByRelId", s.extImagesByRel)
	for _, path := range []string{"/upload/downById", "/upload/zjjgDownById"} {
		reg("GET", path, s.extDownload)
	}
	reg("GET", "/upload/imageShow", s.extImage)
	reg("GET", "/upload/previewStream", s.extStream)
	reg("GET", "/upload/zjjgPreviewStream", s.extStream)
	// @RequestMapping has no method constraint in Java for getVideo/play.
	mux.HandleFunc("/saas-common/upload/getVideo", authGuard(func(w http.ResponseWriter, r *http.Request) { s.extVideo(w, r, false) }))
	mux.HandleFunc("/saas-common/upload/play", authGuard(func(w http.ResponseWriter, r *http.Request) { s.extVideo(w, r, true) }))
	reg("POST", "/upload/getError", func(w http.ResponseWriter, r *http.Request) { platform.Unknown(w) })
	reg("POST", "/upload/bigUpload", s.extBigUpload)
	reg("GET", "/download", func(w http.ResponseWriter, r *http.Request) {
		p := os.Getenv("PANDA_APK_FILE")
		if p == "" {
			p = "/home/ubuntu/uploadfiles/xmmc.apk"
		}
		b, err := os.ReadFile(p)
		if err != nil {
			platform.Unknown(w)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=xmmc_1.17.apk")
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		w.Write(b)
	})
	for _, path := range []string{"/api/v1/download/getFile", "/api/v1/download/limit"} {
		reg("GET", path, func(w http.ResponseWriter, r *http.Request) {
			serveLocalDownload(w, r, "D:/xmmc.apk", strings.HasSuffix(r.URL.Path, "/limit"))
		})
	}
	reg("POST", "/TenOcrController/getOcrResult", s.extOCR)
}

// DownloadController:本地固定 D:/xmmc.apk，不存在时返回200纯文本而非JsonResult。
func serveLocalDownload(w http.ResponseWriter, r *http.Request, p string, limited bool) {
	b, err := os.ReadFile(p)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		io.WriteString(w, "本地测试文件["+p+"]不存在")
		return
	}
	w.Header().Set("Content-Type", "application/x-zip-compressed")
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(p))
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	if !limited {
		w.Write(b)
		return
	}
	speed := 500
	if v := r.URL.Query().Get("speed"); v != "" {
		n, e := strconv.Atoi(v)
		if e != nil {
			platform.Unknown(w)
			return
		}
		speed = n
	}
	start := time.Now()
	for i := 0; i < len(b); i += 1024 {
		end := i + 1024
		if end > len(b) {
			end = len(b)
		}
		w.Write(b[i:end])
		if speed > 0 {
			wait := time.Duration(end)*time.Second/time.Duration(speed*1024) - time.Since(start)
			if wait > 0 {
				select {
				case <-r.Context().Done():
					return
				case <-time.After(wait):
				}
			}
		}
	}
}
func extBase64(p string) string { b, _ := os.ReadFile(p); return base64.StdEncoding.EncodeToString(b) }
func extIsImage(p string) bool {
	switch strings.ToLower(filepath.Ext(p)) {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	}
	return false
}
func (s *Service) extFile(r *http.Request, name string) (map[string]any, error) {
	d := s.pick(r)
	q := "SELECT * FROM chain_common_file WHERE file_name=? OR association_id=?"
	args := []any{name, name}
	if _, err := strconv.ParseInt(name, 10, 32); err == nil {
		q = "SELECT * FROM chain_common_file WHERE id=?"
		args = []any{name}
	}
	rows, err := queryRows(d, q, args...)
	if err != nil {
		return nil, err
	}
	if len(rows) > 1 {
		return nil, fmt.Errorf("selectOne returned more than one")
	}
	if len(rows) == 1 {
		if !extRowUsable(rows[0], name) {
			return nil, nil
		}
		return rows[0], nil
	}
	imgs, err := queryRows(d, "SELECT * FROM chain_common_img WHERE file_name=? OR id=?", name, name)
	if err != nil {
		return nil, err
	}
	if len(imgs) > 1 {
		return nil, fmt.Errorf("selectOne returned more than one")
	}
	if len(imgs) == 1 {
		if !extRowUsable(imgs[0], name) {
			return nil, nil
		}
		return map[string]any{"file_path": imgs[0]["file_path"]}, nil
	}
	return nil, nil
}
func (s *Service) extDownload(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("fileName")
	if name != "" && (name != filepath.Base(name) || name == "..") {
		platform.FailMsg(w, "未找到对应文件") // 不接受路径分隔符:否则能读到 SFTP 目录之外的任意文件
		return
	}
	p := os.Getenv("PANDA_TL_SFTP_DIR") + name
	real := name
	if _, err := os.Stat(p); err != nil {
		f, e := s.extFile(r, name)
		if e != nil {
			platform.Unknown(w)
			return
		}
		if f == nil {
			platform.FailMsg(w, "未找到对应文件")
			return
		}
		if !extRowUsable(f, "") {
			platform.FailMsg(w, "未找到对应文件")
			return
		}
		real = str(f, "file_real_name")
		p = filepath.Join(str(f, "file_path"), str(f, "file_name"))
	}
	b, err := os.ReadFile(p)
	if err != nil {
		platform.Unknown(w)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename="+real)
	w.Write(javaDownloadBytes(b))
}
func (s *Service) extImage(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.URL.Query()["credit"]; ok && !shareTokenExists(r.URL.Query().Get("credit")) {
		return
	}
	f, err := s.extFile(r, r.URL.Query().Get("fileName"))
	if err != nil {
		platform.Unknown(w)
		return
	}
	if f == nil {
		platform.FailMsg(w, "未找到对应文件")
		return
	}
	name := str(f, "file_name")
	if name == "" {
		platform.Unknown(w)
		return
	}
	w.Header().Set("Content-Type", "image/"+strings.TrimPrefix(filepath.Ext(name), "."))
	b, _ := os.ReadFile(filepath.Join(str(f, "file_path"), name))
	w.Write(b)
}
func (s *Service) extStream(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("fileName")
	f, err := s.extFile(r, name)
	if err != nil || f == nil || name == "" || name != filepath.Base(name) {
		platform.Unknown(w)
		return
	}
	p := str(f, "file_path")
	real := str(f, "file_name")
	if strings.HasSuffix(r.URL.Path, "zjjgPreviewStream") {
		name = real
	}
	if strings.HasSuffix(strings.ToLower(real), ".doc") || strings.HasSuffix(strings.ToLower(real), ".docx") {
		rows, err := queryRows(s.pick(r), "SELECT * FROM chain_common_img WHERE file_id=?", f["id"])
		if err != nil || len(rows) == 0 {
			platform.Unknown(w)
			return
		}
		p = str(rows[0], "file_path")
		name = str(rows[0], "file_name")
	}
	b, e := os.ReadFile(filepath.Join(p, name))
	if e != nil {
		return
	}
	typ := mime.TypeByExtension(filepath.Ext(name))
	if typ == "" {
		typ = "application/octet-stream"
	}
	// 仅图片与 PDF 允许内联;HTML/SVG 等一律强制下载,否则同域下传一个 .html 就能盗走登录 token
	if !extIsImage(name) && strings.ToLower(filepath.Ext(name)) != ".pdf" {
		typ = "application/octet-stream"
		w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": name}))
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", typ)
	w.Write(b)
}
func (s *Service) extVideo(w http.ResponseWriter, r *http.Request, play bool) {
	name := r.URL.Query().Get("fileName")
	f, e := s.extFile(r, name)
	if e != nil {
		platform.Unknown(w)
		return
	}
	if f == nil {
		platform.FailMsg(w, "未找到对应文件")
		return
	}
	real := str(f, "file_name")
	if !play && real == "" {
		platform.Unknown(w)
		return
	}
	if !play {
		name = real
	}
	serveVideo(w, r, filepath.Join(str(f, "file_path"), name), play)
}
func serveVideo(w http.ResponseWriter, r *http.Request, p string, play bool) {
	f, e := os.Open(p)
	if e != nil {
		if !play {
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, "file:"+p+" not exists")
		}
		return
	}
	defer f.Close()
	info, _ := f.Stat()
	size := info.Size()
	start := int64(0)
	rg := r.Header.Get("Range")
	if rg != "" {
		v := strings.SplitN(strings.TrimPrefix(rg, "bytes="), "-", 2)
		start, e = strconv.ParseInt(v[0], 10, 64)
		if e != nil || start < 0 {
			return
		}
		f.Seek(start, io.SeekStart)
	}
	if play {
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, size-1, size))
		b := make([]byte, 1<<20)
		n, _ := f.Read(b)
		w.Header().Set("Content-Length", strconv.Itoa(n))
		w.WriteHeader(206)
		w.Write(b[:n])
		return
	}
	if rg != "" {
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Length", strconv.FormatInt(size-start, 10))
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, size-1, size))
		w.WriteHeader(206)
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
		w.Header().Set("Content-Disposition", "attachment; filename="+p)
	}
	io.Copy(w, f)
}

func (s *Service) extPreview(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if len(m) == 0 { // 空条件=全表:逐行读盘转 base64 会打满内存
		platform.Unknown(w)
		return
	}
	cond := "1=1"
	args := []any{}
	for _, k := range platform.SortedKeys(m) {
		if !externalColumn(k) {
			platform.Unknown(w)
			return
		}
		if m[k] == nil {
			cond += " AND " + k + " IS NULL"
		} else {
			cond += " AND " + k + "=?"
			args = append(args, m[k])
		}
	}
	rows, err := queryRows(s.pick(r), "SELECT * FROM chain_common_img WHERE "+cond, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	out := []map[string]any{}
	for _, row := range rows {
		if !extRowUsable(row, "") {
			continue
		}
		out = append(out, map[string]any{"fileName": row["file_name"], "filePath": "data:image/jpg;base64," + extBase64(filepath.Join(str(row, "file_path"), str(row, "file_name")))})
	}
	if len(out) == 0 {
		if m["file_id"] == nil {
			platform.Unknown(w)
			return
		}
		f, err := platform.GetByID(s.pick(r), "chain_common_file", m["file_id"])
		if err != nil || f == nil {
			platform.Unknown(w)
			return
		}
		if !extRowUsable(f, "") {
			platform.Unknown(w)
			return
		}
		name := str(f, "file_name")
		prefix := ""
		if extIsImage(name) {
			prefix = "data:image/jpg;base64,"
		} else if dot := strings.Index(name, "."); dot >= 0 && strings.ToLower(name[dot:]) == ".pdf" {
			prefix = "data: application/pdf;base64,"
		}
		if prefix != "" {
			out = append(out, map[string]any{"fileName": name, "filePath": prefix + extBase64(filepath.Join(str(f, "file_path"), name))})
		}
	}
	platform.Success(w, out)
}
func externalColumn(k string) bool {
	if k == "" {
		return false
	}
	for _, c := range k {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_') {
			return false
		}
	}
	return true
}
func (s *Service) extByRel(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond := "1=1"
	args := []any{}
	allowed := map[string]string{"id": "id", "companyId": "company_id", "fileRealName": "file_real_name", "fileName": "file_name", "filePath": "file_path", "fileType": "file_type", "ownerType": "owner_type", "ownerId": "owner_id", "relationId": "relation_id", "relationType": "relation_type", "fileDescription": "file_description", "fType": "f_type", "isCompressedFile": "is_compressed_file", "effective": "effective", "createTime": "create_time", "updateTime": "update_time", "remark1": "remark1", "remark2": "remark2", "remark3": "remark3", "remark4": "remark4", "remark5": "remark5", "sourceId": "source_id", "transStatus": "trans_status", "compressedId": "compressed_id", "associationId": "association_id"}
	for _, k := range platform.SortedKeys(m) {
		if col, ok := allowed[k]; ok {
			cond += " AND " + col + " LIKE ?"
			args = append(args, m[k])
		}
	}
	rows, err := queryRows(s.pick(r), "SELECT id id,file_real_name name,file_name fileName,file_type fileType FROM chain_common_file WHERE "+cond, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, rows)
}
func (s *Service) extUpload(w http.ResponseWriter, r *http.Request, logged bool) {
	if err := r.ParseMultipartForm(64 << 20); err != nil || r.MultipartForm == nil || len(r.MultipartForm.File["file"]) == 0 {
		platform.Unknown(w)
		return
	}
	defer r.MultipartForm.RemoveAll()
	fhs := r.MultipartForm.File["file"]
	if !logged && len([]rune(fhs[0].Filename)) > 200 {
		platform.Fail(w, 999, "文件名称过长")
		return
	}
	dir := filepath.Join(externalRoot(), time.Now().Format("2006-01-02"))
	cu := currentUser(r)
	if logged {
		dir = externalRoot()
		if cu.EnterpriseCode == "" {
			u, e := platform.GetByID(s.pick(r), "chain_auth_user", cu.UserID)
			if e != nil || u == nil {
				platform.Unknown(w)
				return
			}
			dir = filepath.Join(dir, str(u, "phone"))
		}
	}
	out := []map[string]any{}
	for _, fh := range fhs {
		suffix := filepath.Ext(fh.Filename)
		if suffix == "" {
			if !logged {
				platform.Unknown(w)
				return
			}
			break
		}
		if e := os.MkdirAll(dir, 0755); e != nil {
			if !logged {
				platform.Unknown(w)
				return
			}
			break
		}
		src, e := fh.Open()
		if e != nil {
			if !logged {
				platform.Unknown(w)
				return
			}
			break
		}
		name := extUUID() + suffix
		p := filepath.Join(dir, name)
		dst, e := os.Create(p)
		if e != nil {
			src.Close()
			if !logged {
				platform.Unknown(w)
				return
			}
			break
		}
		_, e = io.Copy(dst, src)
		dst.Close()
		src.Close()
		if e != nil {
			if !logged {
				platform.Unknown(w)
				return
			}
			break
		}
		vals := map[string]any{"file_real_name": fh.Filename, "file_name": name, "file_path": dir}
		if logged {
			vals["file_type"] = "5"
			vals["owner_type"] = ""
			vals["owner_id"] = cu.UserID
			vals["company_id"] = cu.CompanyID
		}
		id, e := platform.InsertNoFill(s.pick(r), "chain_common_file", vals)
		if e != nil {
			if !logged {
				platform.Unknown(w)
				return
			}
			break
		}
		item := map[string]any{"id": strconv.FormatInt(id, 10), "fileName": name, "name": fh.Filename}
		if logged && extIsImage(name) {
			cid := s.compressFile(r, id)
			if _, e := platform.UpdateByIDNoFill(s.pick(r), "chain_common_file", id, map[string]any{"compressed_id": cid}); e != nil {
				break
			}
			item["compressedId"] = strconv.FormatInt(cid, 10)
		}
		if logged && (strings.Contains(name, ".xls") || strings.Contains(name, ".xlsx")) {
			imgdir := filepath.Join(dir, strings.SplitN(name, ".", 2)[0]+"_img")
			imgs, e := platform.RenderExcelUpload(r.Context(), p, imgdir)
			if e != nil {
				break
			}
			for _, img := range imgs {
				if _, e = platform.InsertNoFill(s.pick(r), "chain_common_img", map[string]any{"file_id": id, "file_name": img, "file_path": imgdir}); e != nil {
					break
				}
			}
			if e != nil {
				break
			}
		}
		if logged && (strings.Contains(name, ".doc") || strings.Contains(name, ".docx")) {
			imgdir := filepath.Join(externalRoot(), "convertFiles")
			preview := platform.PreviewName()
			if e = platform.RenderDocumentPDF(r.Context(), p, filepath.Join(imgdir, preview)); e != nil {
				break
			}
			if _, e = platform.InsertNoFill(s.pick(r), "chain_common_img", map[string]any{"file_id": id, "file_name": preview, "file_path": imgdir}); e != nil {
				break
			}
		}
		if logged && extIsImage(name) {
			item["url"] = "data:image/jpg;base64," + extBase64(p)
		}
		out = append(out, item)
	}
	platform.Success(w, out)
}
func (s *Service) extBigUpload(w http.ResponseWriter, r *http.Request) {
	if e := r.ParseMultipartForm(64 << 20); e != nil || r.MultipartForm == nil {
		platform.Unknown(w)
		return
	}
	defer r.MultipartForm.RemoveAll()
	fhs := r.MultipartForm.File["file"]
	name := r.FormValue("fileName")
	chunk, e1 := strconv.Atoi(r.FormValue("chunkNumber"))
	total, e2 := strconv.Atoi(r.FormValue("totalChunks"))
	if len(fhs) == 0 || name == "" || e1 != nil || e2 != nil {
		platform.Unknown(w)
		return
	}
	if name != filepath.Base(name) || name == ".." { // 只允许纯文件名:带分隔符就能把分片写到上传目录之外
		platform.Unknown(w)
		return
	}
	dir := os.Getenv("PANDA_BIG_UPLOAD_DIR")
	if dir == "" {
		dir = filepath.Join(externalRoot(), "bigUpload")
	}
	if e := os.MkdirAll(dir, 0755); e != nil {
		platform.Unknown(w)
		return
	}
	part := filepath.Join(dir, name+".part"+strconv.Itoa(chunk))
	src, e := fhs[0].Open()
	if e != nil {
		platform.Unknown(w)
		return
	}
	defer src.Close()
	dst, e := os.Create(part)
	if e != nil {
		platform.Unknown(w)
		return
	}
	_, e = io.Copy(dst, src)
	dst.Close()
	if e != nil {
		platform.Unknown(w)
		return
	}
	files, _ := os.ReadDir(dir)
	count := 0
	for _, f := range files {
		if strings.HasPrefix(f.Name(), name) {
			count++
		}
	}
	out := map[string]string{}
	if count == total {
		for i := 1; i <= total; i++ {
			p := filepath.Join(dir, name+".part"+strconv.Itoa(i))
			b, e := os.ReadFile(p)
			if e != nil {
				platform.Unknown(w)
				return
			}
			f, e := os.OpenFile(filepath.Join(dir, name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if e != nil {
				platform.Unknown(w)
				return
			}
			_, e = f.Write(b)
			f.Close()
			if e != nil {
				platform.Unknown(w)
				return
			}
			os.Remove(p)
		}
		cu := currentUser(r)
		id, e := platform.InsertNoFill(s.pick(r), "chain_common_file", map[string]any{"company_id": cu.CompanyID, "owner_id": cu.UserID, "owner_type": "", "file_real_name": fhs[0].Filename, "file_name": name, "file_path": dir, "effective": true})
		if e != nil {
			platform.Unknown(w)
			return
		}
		out = map[string]string{"id": strconv.FormatInt(id, 10), "fileName": name, "name": fhs[0].Filename}
	}
	platform.Success(w, out)
}

// OCR is an explicit external adapter: never copy Java's embedded credentials.
// PANDA_OCR_ADAPTER_URL accepts {imageBase64,type} and returns the Java JsonResult.
// Missing input/file follows local Java 999; valid unsupported type returns empty success.
// Without an adapter configured, active OCR types fail 701 explicitly rather than fabricate recognition.
func (s *Service) extOCR(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if m["id"] == nil || m["type"] == nil {
		platform.Unknown(w)
		return
	}
	f, e := platform.GetByID(s.pick(r), "chain_common_file", m["id"])
	if e != nil || f == nil {
		platform.Unknown(w)
		return
	}
	typ, e := strconv.Atoi(str(m, "type"))
	if e != nil {
		platform.Unknown(w)
		return
	}
	if typ < 1 || typ > 3 {
		platform.Success(w, "")
		return
	}
	url := os.Getenv("PANDA_OCR_ADAPTER_URL")
	if url == "" {
		action := map[int]string{1: "IDCardOCR", 2: "BizLicenseOCR", 3: "GeneralBasicOCR"}[typ]
		response, err := platform.TencentJSON(r.Context(), "ocr", action, "2018-11-19", "ap-beijing", map[string]any{"ImageBase64": "data:image/jpg;base64," + extBase64(filepath.Join(str(f, "file_path"), str(f, "file_name")))})
		if err != nil {
			code := "ClientNetworkError"
			var provider *platform.TencentError
			if errors.As(err, &provider) {
				code = provider.Code
			}
			platform.FailData(w, 701, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(legacyOCRResult(response, typ))
		return
	}
	result, e := callOCRAdapter(r.Context(), url, "data:image/jpg;base64,"+extBase64(filepath.Join(str(f, "file_path"), str(f, "file_name"))), typ)
	if e != nil {
		platform.FailData(w, 701, e.Error(), "ExternalServiceUnavailable")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(result)
}

// FileUtil.download writes buff.length rather than the read length: preserve padding/stale tail.
func javaDownloadBytes(b []byte) []byte {
	out := make([]byte, 0, ((len(b)+1023)/1024)*1024)
	buf := make([]byte, 1024)
	for len(b) > 0 {
		n := copy(buf, b)
		out = append(out, buf...)
		b = b[n:]
	}
	return out
}

func callOCRAdapter(ctx context.Context, url, image string, typ int) (platform.JsonResult, error) {
	var result platform.JsonResult
	payload, e := json.Marshal(map[string]any{"imageBase64": image, "type": typ})
	if e != nil {
		return result, e
	}
	req, e := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if e != nil {
		return result, e
	}
	req.Header.Set("Content-Type", "application/json")
	resp, e := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if e != nil {
		return result, e
	}
	defer resp.Body.Close()
	e = json.NewDecoder(resp.Body).Decode(&result)
	return result, e
}

// extImagesByRel 保留 Java 的空对象、可选 yyzz 和有序 lbts 字段。
func (s *Service) extImagesByRel(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if len(m) == 0 {
		platform.Success(w, map[string]any{}) // 空条件=全表读盘
		return
	}
	cond := "1=1"
	args := []any{}
	for _, key := range platform.SortedKeys(m) {
		if !externalColumn(key) {
			platform.Unknown(w)
			return
		}
		if m[key] == nil {
			cond += " AND " + key + " IS NULL"
		} else {
			cond += " AND " + key + "=?"
			args = append(args, m[key])
		}
	}
	rows, err := queryRows(s.pick(r), "SELECT * FROM chain_common_file WHERE "+cond, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	result := map[string]any{}
	images := []map[string]any{}
	for _, row := range rows {
		if !extRowUsable(row, "") {
			continue
		}
		preview := "data:image/jpg;base64," + extBase64(filepath.Join(str(row, "file_path"), str(row, "file_name")))
		switch str(row, "file_type") {
		case "3":
			result["yyzz"] = preview
		case "4":
			images = append(images, map[string]any{"id": row["id"], "url": preview})
		}
		result["lbts"] = images
	}
	platform.Success(w, result)
}

// extRowUsable 文件行的磁盘路径只接受"非空目录 + 纯文件名":
// 目录为空时会相对进程工作目录读文件(.env、keys/ 等),文件名带分隔符则可穿越到任意路径。
func extRowUsable(row map[string]any, name string) bool {
	if name == "" {
		name = str(row, "file_name")
	}
	return str(row, "file_path") != "" && name != "" && name != ".." && name == filepath.Base(name)
}
