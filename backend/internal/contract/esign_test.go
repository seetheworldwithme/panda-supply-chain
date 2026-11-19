package contract

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

func TestEsignV3EmptyFlowKeepsOriginalPendingResponse(t *testing.T) {
	s := &Service{DB: &db.DB{}}
	for _, body := range []string{`{}`, `{"eqbFlowId":"   "}`} {
		w := httptest.NewRecorder()
		s.esignExternal(w, httptest.NewRequest("POST", "/saas-common/chainEqbFlow/v3/selSignUrl", strings.NewReader(body)))
		if !strings.Contains(w.Body.String(), `"code":602`) || !strings.Contains(w.Body.String(), "请等候其他签署人进行确认") {
			t.Fatal(w.Body.String())
		}
	}
}

func TestEsignBusinessUsesServerFlow(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	admin, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := fmt.Sprintf("audit_sign_%d", time.Now().UnixNano())
	if _, err = admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_auth_user", "chain_basic_enterprise", "chain_eqb_flow", "chain_eqb_singers", "chain_eqb_config", "chain_contract", "chain_common_file", "chain_common_img", "chain_common_log"} {
		if _, err = d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{
		"INSERT INTO chain_auth_user(id,account,username,phone,level,account_type,handle_result,eqb_per_id) VALUES(1,'fixture','fixture-name','00000001121',1,1,1,'fixture-person')",
		"INSERT INTO chain_basic_enterprise(id,enterprise_name,rate_number,pub_flow_id,eqb_org_id,sign_free_count,sign_charge_count) VALUES(1,'fixture-enterprise','fixture-credit',1,'fixture-org',2,0)",
		"INSERT INTO chain_eqb_config(id,orgavailableauthtypes,personavailableauthtypes,willtypes) VALUES(1,'ORG_BANK_TRANSFER','PSN_FACEAUTH_BYURL','CODE_SMS')",
	} {
		if _, err = d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v3/signflows/createFlowOneStep":
			var m map[string]any
			if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
				t.Error(err)
			}
			if _, ok := m["signers"].([]any); !ok {
				t.Error("发起签署未传签署人")
			}
			io.WriteString(w, `{"code":0,"data":{"flowId":"fixture-flow"}}`)
		case "/v3/signflows/fixture-flow/executeUrl":
			io.WriteString(w, `{"code":0,"data":{"url":"fixture-sign-url"}}`)
		case "/v1/signflows/fixture-flow/signers":
			io.WriteString(w, `{"code":0,"data":{"signers":[{"fixture":true}]}}`)
		case "/v1/signflows/fixture-flow/documents":
			io.WriteString(w, `{"code":0,"data":{"docs":[{"fileUrl":"fixture-download"}]}}`)
		case "/v1/files/getUploadUrl":
			json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"fileId": "fixture-uploaded", "uploadUrl": server.URL + "/fixture-upload"}})
		case "/fixture-upload":
			io.WriteString(w, `{"errCode":0}`)
		case "/v1/files/fixture-uploaded/status":
			io.WriteString(w, `{"code":0,"data":{"status":2}}`)
		default:
			t.Errorf("未预期签章请求 %s", r.URL.Path)
			io.WriteString(w, `{}`)
		}
	}))
	defer server.Close()
	t.Setenv("host", server.URL)
	t.Setenv("appId", "fixture-app")
	t.Setenv("appKey", "fixture-key")
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	kp := &auth.KeyPair{Private: key, Public: &key.PublicKey}
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	mux := http.NewServeMux()
	Register(mux, s)
	token, _ := kp.Generate(auth.JWTClaims{Account: "fixture", UserID: "1", CompanyID: "1"})
	call := func(path, body string) map[string]any {
		r := httptest.NewRequest("POST", "/saas-common/chainEqbFlow/"+path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		auth.Middleware(kp, d)(mux).ServeHTTP(w, r)
		var result map[string]any
		json.Unmarshal(w.Body.Bytes(), &result)
		if w.Code != 200 || result["code"] != float64(200) {
			t.Fatalf("%s 未完成: %d %s", path, w.Code, w.Body.String())
		}
		return result
	}
	info := call("selLanucnInfo", `{}`)
	if data, _ := info["data"].(map[string]any); data["companyName"] != "fixture-enterprise" {
		t.Fatalf("发起信息不符: %v", info)
	}
	launched := call("launchFlow", `{"companyId":1,"topicName":"fixture-topic","field":"fixture-file","contractName":"fixture-contract","singers":[{"signPhone":"00000001121","signName":"fixture-name","signRole":2,"signEtcName":"fixture-enterprise"}]}`)
	if launched["data"] != "fixture-flow" {
		t.Fatalf("签署流程未创建: %v", launched)
	}
	var free, status int
	var flow string
	if err = d.QueryRow("SELECT sign_free_count FROM chain_basic_enterprise WHERE id=1").Scan(&free); err != nil || free != 1 {
		t.Fatalf("份数扣减不符: %d %v", free, err)
	}
	if err = d.QueryRow("SELECT flow_status,eqb_flow_id FROM chain_eqb_flow").Scan(&status, &flow); err != nil || status != 1 || flow != "fixture-flow" {
		t.Fatalf("签署流程未落库: %d %s %v", status, flow, err)
	}
	if result := call("v3/selSignUrl", `{"eqbFlowId":"fixture-flow"}`); result["data"] != "fixture-sign-url" {
		t.Fatal(result)
	}
	call("selSignPeople", `{"eqbFlowId":"fixture-flow"}`)
	call("getById", `{"id":1,"eqbFlowId":"fixture-flow"}`)
	if result := call("v3DownFile", `{"eqbFlowId":"fixture-flow"}`); result["data"] != "fixture-download" {
		t.Fatal(result)
	}
	t.Setenv("appSecret", "fixture-callback-key")
	if _, err = d.Exec("INSERT INTO chain_contract(id,eqb_flow_id,is_esign_complete) VALUES(1,'fixture-flow',0)"); err != nil {
		t.Fatal(err)
	}
	body := `{"action":"SIGN_FLOW_FINISH","flowStatus":2,"flowId":"fixture-flow"}`
	callback := func(valid bool) *httptest.ResponseRecorder {
		r := httptest.NewRequest("POST", "/saas-common/chainEqbFlow/getSignCallBack", strings.NewReader(body))
		r.Header.Set("X-Tsign-Open-TIMESTAMP", "fixture-time")
		sig := "bad"
		if valid {
			mac := hmac.New(sha256.New, []byte("fixture-callback-key"))
			mac.Write([]byte("fixture-time" + body))
			sig = hex.EncodeToString(mac.Sum(nil))
		}
		r.Header.Set("X-Tsign-Open-SIGNATURE", sig)
		w := httptest.NewRecorder()
		auth.Middleware(kp, d)(mux).ServeHTTP(w, r)
		return w
	}
	if w := callback(false); !strings.Contains(w.Body.String(), `"code":403`) {
		t.Fatalf("坏签回调未按网关拒绝: %d %s", w.Code, w.Body.String())
	}
	if w := callback(true); w.Code != 200 || !strings.Contains(w.Body.String(), `"code":200`) {
		t.Fatalf("真实验签回调失败: %d %s", w.Code, w.Body.String())
	}
	if err = d.QueryRow("SELECT is_esign_complete FROM chain_contract WHERE id=1").Scan(&status); err != nil || status != 1 {
		t.Fatalf("签章回调未更新合同: %d %v", status, err)
	}
	for _, q := range []string{"UPDATE chain_eqb_flow SET flow_status=0", "CREATE TRIGGER reject_callback_contract BEFORE UPDATE ON chain_contract FOR EACH ROW SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT='fixture failure'"} {
		if _, err = d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	if w := callback(true); !strings.Contains(w.Body.String(), `"code":999`) {
		t.Fatal(w.Body.String())
	}
	if err = d.QueryRow("SELECT flow_status FROM chain_eqb_flow LIMIT 1").Scan(&status); err != nil || status != 2 {
		t.Fatalf("原Java getSignCallBack不启事务，后段失败不应回滚流程状态: %d %v", status, err)
	}
	if _, err = d.Exec("DROP TRIGGER reject_callback_contract"); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	t.Setenv("PANDA_UPLOAD_DIR", dir)
	os.WriteFile(filepath.Join(dir, "fixture.pdf"), []byte("%PDF-fixture"), 0600)
	if _, err = d.Exec("INSERT INTO chain_common_file(id,file_name,file_path) VALUES(1,'fixture.pdf',?)", dir); err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("POST", "/saas-contract/chainContract/add", strings.NewReader(`{"companyId":1,"contractCode":"fixture-linked","contractType":1,"isPayOnline":false,"isEsign":true,"fileStatus":1,"contractIds":"1","eqbFlow":{"companyId":1,"topicName":"linked","singers":[{"signPhone":"00000001121","signName":"fixture","signRole":2}]}}`))
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	auth.Middleware(kp, d)(mux).ServeHTTP(w, r)
	if !strings.Contains(w.Body.String(), `"code":200`) {
		t.Fatalf("合同保存未接签章: %s", w.Body.String())
	}
	if err = d.QueryRow("SELECT eqb_flow_id FROM chain_contract WHERE contract_code='fixture-linked'").Scan(&flow); err != nil || flow != "fixture-flow" {
		t.Fatalf("合同未关联真实签章流程: %s %v", flow, err)
	}
}
