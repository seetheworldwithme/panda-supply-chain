package contract

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// 原 Java 合同服务通过 Feign 调用公共签章服务，复用同一上传、发起流程。
func (s *Service) contractSign(r *http.Request, m map[string]any, update bool) (string, error) {
	ids := strings.Split(str(m, "contractIds"), ",")
	file, err := platform.GetByID(s.pick(r), "chain_common_file", ids[0])
	if err != nil || file == nil {
		return "", fmt.Errorf("合同签章附件不存在")
	}
	data, err := os.ReadFile(filepath.Join(str(file, "file_path"), str(file, "file_name")))
	if err != nil {
		return "", err
	}
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	part, err := form.CreateFormFile("file", str(file, "file_name"))
	if err != nil {
		return "", err
	}
	if _, err = part.Write(data); err != nil {
		return "", err
	}
	if err = form.Close(); err != nil {
		return "", err
	}
	request := r.Clone(r.Context())
	request.Body = io.NopCloser(&body)
	request.Header.Set("Content-Type", form.FormDataContentType())
	request.MultipartForm, request.Form, request.PostForm = nil, nil, nil
	response := httptest.NewRecorder()
	s.esignUpload(response, request)
	var uploaded struct {
		Data map[string]any `json:"data"`
	}
	if err = json.Unmarshal(response.Body.Bytes(), &uploaded); err != nil || uploaded.Data == nil {
		return "", fmt.Errorf("签章上传返回空数据")
	}
	flow, ok := m["eqbFlow"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("签章流程为空")
	}
	flow["field"] = uploaded.Data["field"]
	if update {
		flow["contractId"] = m["id"]
	}
	encoded, err := json.Marshal(flow)
	if err != nil {
		return "", err
	}
	request = r.Clone(r.Context())
	request.Body = io.NopCloser(bytes.NewReader(encoded))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	s.esignFlow(response, request, true)
	var result platform.JsonResult
	if err = json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		return "", err
	}
	if !update && result.Data == nil {
		return "", fmt.Errorf("签章流程返回空数据")
	}
	return platform.JnullAny(result.Data), nil
}

func (s *Service) esignCallback(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil || !platform.VerifyEsignCallback(r, body) {
		platform.Unknown(w)
		return
	}
	var m map[string]any
	if json.Unmarshal(body, &m) != nil || m == nil {
		platform.Unknown(w)
		return
	}
	// 原 AOP 将 get* 设为 NOT_SUPPORTED，第二条失败不回滚已更新的流程状态。
	apply := func() error {
		if m["action"] != "SIGN_FLOW_FINISH" {
			return nil
		}
		status := int(num(m["flowStatus"]))
		if _, err := s.pick(r).Exec("UPDATE chain_eqb_flow SET flow_status=? WHERE eqb_flow_id=?", status, m["flowId"]); err != nil {
			return err
		}
		if status == 2 {
			_, err := s.pick(r).Exec("UPDATE chain_contract SET is_esign_complete=1 WHERE eqb_flow_id=?", m["flowId"])
			return err
		}
		return nil
	}
	err = apply()
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

func (s *Service) esignUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil || r.MultipartForm == nil || len(r.MultipartForm.File["file"]) == 0 {
		platform.Unknown(w)
		return
	}
	defer r.MultipartForm.RemoveAll()
	fh := r.MultipartForm.File["file"][0]
	file, err := fh.Open()
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, 32<<20)) // e签宝上传文件封顶 32MB
	if err != nil {
		platform.Unknown(w)
		return
	}
	u := currentUser(r)
	root := os.Getenv("PANDA_UPLOAD_DIR")
	if root == "" {
		root = "data/uploads"
	}
	dir := filepath.Join(root, u.EnterpriseCode, "sign")
	if err = os.MkdirAll(dir, 0755); err != nil {
		platform.Unknown(w)
		return
	}
	random := make([]byte, 16)
	if _, err = rand.Read(random); err != nil {
		platform.Unknown(w)
		return
	}
	name := fmt.Sprintf("%x", random) + filepath.Ext(fh.Filename)
	if err = os.WriteFile(filepath.Join(dir, name), data, 0644); err != nil {
		platform.Unknown(w)
		return
	}
	result := map[string]any{}
	err = platform.WithTx(s.pick(r), func(tx *sql.Tx) error {
		id, err := platform.TxInsertNoFill(tx, "chain_common_file", map[string]any{"company_id": u.CompanyID, "owner_id": u.UserID, "file_real_name": fh.Filename, "file_name": name, "file_path": dir, "file_type": "6", "owner_type": ""})
		if err != nil {
			return err
		}
		if strings.Contains(name, ".doc") {
			previewDir := filepath.Join(root, "convertFiles")
			preview := platform.PreviewName()
			if err = platform.RenderDocumentPDF(r.Context(), filepath.Join(dir, name), filepath.Join(previewDir, preview)); err != nil {
				return err
			}
			if _, err = platform.TxInsertNoFill(tx, "chain_common_img", map[string]any{"file_id": id, "file_name": preview, "file_path": previewDir}); err != nil {
				return err
			}
		}
		result["id"], result["fileName"] = id, name
		digest := md5.Sum(data)
		convert := "true"
		if strings.Contains(strings.ToLower(name), "pdf") {
			convert = "false"
		}
		res := platform.EsignJSON(r.Context(), "POST", "/v1/files/getUploadUrl", map[string]any{"contentMd5": base64.StdEncoding.EncodeToString(digest[:]), "contentType": "application/octet-stream", "convert2Pdf": convert, "fileName": name, "fileSize": len(data)})
		message, ok := res["data"].(map[string]any)
		if !ok {
			return fmt.Errorf("Java uploadMsg 为 null")
		}
		response := platform.EsignUpload(r.Context(), str(message, "uploadUrl"), data)
		if str(response, "errCode") == "0" {
			for i := 0; i < 100; i++ { // 最多等 5 分钟:无限轮询会一直占着事务连接(连接池仅 20 条)
				res = platform.EsignJSON(r.Context(), "GET", "/v1/files/"+platform.JnullAny(message["fileId"])+"/status", nil)
				state, ok := res["data"].(map[string]any)
				if !ok {
					return fmt.Errorf("Java 文件状态 data 为 null")
				}
				wanted := "5"
				if convert == "false" {
					wanted = "2"
				}
				if str(state, "status") == wanted {
					result["field"] = message["fileId"]
					break
				}
				select {
				case <-r.Context().Done():
					return r.Context().Err()
				case <-time.After(3 * time.Second):
				}
			}
			if result["field"] == nil {
				return fmt.Errorf("e签宝文件状态等待超时")
			}
		}
		return nil
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, result)
}
