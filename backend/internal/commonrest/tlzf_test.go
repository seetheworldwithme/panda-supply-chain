package commonrest

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func TestAllinpayBusinessUsesServerProtocol(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	admin, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := fmt.Sprintf("audit_tl_%d", time.Now().UnixNano())
	if _, err = admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_auth_user", "chain_basic_enterprise", "chain_tlzf_order", "chain_tlzf_log", "chain_tlzf_config", "chain_common_log", "chain_contract", "chain_tlzf_contract_log"} {
		source := "chain"
		if table == "chain_tlzf_contract_log" {
			source = "zjjg"
		}
		if _, err = d.Exec("CREATE TABLE " + table + " LIKE " + source + "." + table); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{
		"INSERT INTO chain_auth_user(id,account,level,account_type,handle_result) VALUES(1,'fixture',1,1,1)",
		"INSERT INTO chain_basic_enterprise(id,enterprise_name,tl_sign_num,tl_mcht_no,tl_acct_sub_no) VALUES(1,'fixture买方','sign1','merchant1','account1'),(2,'fixture卖方','sign2','merchant2','account2')",
		"INSERT INTO chain_contract(id,contract_code,sales_unit_id,sum_price,contract_status,sum_freeze_amt,sum_thaw_amt,sum_refund_amt,residue_freeze_amt,sum_trace_fees) VALUES(1,'fixture-contract',2,100,1,0,0,0,0,0)",
	} {
		if _, err = d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	dir := t.TempDir()
	key, cert, pfx := filepath.Join(dir, "key.pem"), filepath.Join(dir, "cert.pem"), filepath.Join(dir, "merchant.pfx")
	for _, args := range [][]string{
		{"genpkey", "-algorithm", "SM2", "-out", key},
		{"req", "-new", "-x509", "-key", key, "-sm3", "-subj", "/CN=fixture", "-days", "1", "-out", cert},
		{"pkcs12", "-export", "-inkey", key, "-in", cert, "-passout", "pass:fixture", "-out", pfx},
	} {
		if output, err := exec.Command("openssl", args...).CombinedOutput(); err != nil {
			t.Fatalf("临时证书生成失败: %v %s", err, output)
		}
	}
	t.Setenv("tlPfxPath", pfx)
	t.Setenv("tlCertPath", cert)
	t.Setenv("TL_SM2_PFX_PASSWORD", "fixture")
	t.Setenv("TL_INST_ID", "fixture-inst")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		xml := r.FormValue("msgPlain")
		payload := `<resp_code>0000</resp_code>`
		if strings.Contains(xml, "<processing_code>3002</processing_code>") {
			payload += `<qur_rst>[{"amount":12.30}]</qur_rst>`
		}
		if strings.Contains(xml, "<processing_code>3001</processing_code>") {
			payload += `<qur_rst>1122</qur_rst>`
		}
		plain := `<?xml version="1.0" encoding="GBK"?><transaction><head></head><response>` + payload + `</response></transaction>`
		cmd := exec.Command("openssl", "dgst", "-sm3", "-sign", key, "-sigopt", "distid:1234567812345678")
		cmd.Stdin = strings.NewReader(plain)
		sig, err := cmd.Output()
		if err != nil {
			t.Error(err)
			return
		}
		io.WriteString(w, strings.Replace(plain, "</head>", "<sign_code>"+base64.StdEncoding.EncodeToString(sig)+"</sign_code></head>", 1))
	}))
	defer server.Close()
	t.Setenv("TL_QUERY_REQUEST_URL", server.URL)
	t.Setenv("TL_REQUEST_URL", server.URL)
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	kp := &auth.KeyPair{Private: rsaKey, Public: &rsaKey.PublicKey}
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	defer s.stopTLTasks("", true)
	s.Auth = auth.NewService(d, kp)
	mux := http.NewServeMux()
	registerBespoke(mux, s)
	registerEnterprise(mux, s)
	token, _ := kp.Generate(auth.JWTClaims{Account: "fixture", UserID: "1", CompanyID: "1", Username: "fixture"})
	call := func(path, body string) map[string]any {
		r := httptest.NewRequest("POST", "/saas-common/"+path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		auth.Middleware(kp, d)(mux).ServeHTTP(w, r)
		var result map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("响应无法解析: %s", w.Body.String())
		}
		if result["code"] != float64(200) {
			t.Fatalf("%s 未完成: %s", path, w.Body.String())
		}
		return result
	}
	result := call("chainTlzfOrder/selAccAmt", `{}`)
	if list, ok := result["data"].([]any); !ok || len(list) != 1 {
		t.Fatalf("余额结果不符: %v", result)
	}
	call("chainTlzfOrder/payMoney", `{"buyCompanyId":1,"saleCompanyId":2,"contractId":1,"contractCode":"fixture-contract","amtTran":100,"saleAcctNum":"fixture-account","saleBnkId":"fixture-bank"}`)
	s.stopTLTasks("", true)
	var count int
	if err = d.QueryRow("SELECT count(*) FROM chain_tlzf_order WHERE trans_type=1 AND amt_tran=100 AND trace_fee=1.00").Scan(&count); err != nil || count != 1 {
		t.Fatalf("付款未按 Java 落库: %d %v", count, err)
	}
	call("chainTlzfOrder/selTransResult", `{"id":1}`)
	row, err := platform.GetByID(d, "chain_tlzf_order", 1)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.tlSyncStatus(tlRequest(context.Background(), ""), d, row, false); err != nil {
		t.Fatal(err)
	}
	var frozen string
	if err = d.QueryRow("SELECT sum_freeze_amt FROM chain_contract WHERE id=1").Scan(&frozen); err != nil || frozen != "100.00" {
		t.Fatalf("付款同步不符: %s %v", frozen, err)
	}
	code := s.Auth.IssuePhoneCode("13994243151")
	call("chainTlzfOrder/confirmOrCancel", fmt.Sprintf(`{"code":%q,"reqTimeMinners":"fixture-confirm","contractId":1,"contractCode":"fixture-contract","transType":2,"confirmType":1,"amtJf":20,"actTraceFee":0.50}`, code))
	s.stopTLTasks("", true)
	var amount, fee string
	if err = d.QueryRow("SELECT amt_jf,act_trace_fee FROM chain_tlzf_order WHERE trans_type=2").Scan(&amount, &fee); err != nil || amount != "19.50" || fee != "0.50" {
		t.Fatalf("解冻扣费与 Java 不符: %s %s %v", amount, fee, err)
	}
}
