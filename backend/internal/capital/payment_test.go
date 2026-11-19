package capital

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

// 使用临时密钥和本机银行端点验证真实 HTTP 协议，不访问银行或产生真实交易。
func TestBankOrderQueryUsesJavaSignedProtocol(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, _ := x509.MarshalPKCS8PrivateKey(key)
	t.Setenv("busbank_privateKey", base64.StdEncoding.EncodeToString(der))
	t.Setenv("busbank_appId", "fixture-app")
	t.Setenv("busbank_appsecret", "fixture-secret")
	called := false
	bank := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		b, _ := io.ReadAll(r.Body)
		digest := sha256.Sum256(b)
		sign := hex.EncodeToString(digest[:])
		if r.Method != "POST" || r.URL.Path != "/order-qry" || r.Header.Get("sign") != sign || r.Header.Get("appid") != "fixture-app" || r.Header.Get("verify") != "SHA256withRSAVerify" {
			t.Error("银行请求方法、路径或摘要不符")
		}
		src := "appid=fixture-app&secret=fixture-secret&sign=" + sign + "&timestamp=" + r.Header.Get("timestamp")
		hash := sha256.Sum256([]byte(src))
		sig, _ := base64.StdEncoding.DecodeString(r.Header.Get("apisign"))
		if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, hash[:], sig); err != nil {
			t.Error("银行签名不符", err)
		}
		var req map[string]any
		if err := json.Unmarshal(b, &req); err != nil {
			t.Fatal(err)
		}
		head, _ := req["head"].(map[string]any)
		body, _ := req["body"].(map[string]any)
		dto, _ := body["orderInfoQryDto"].(map[string]any)
		if head["version"] != "1.0" || head["datetime"] == nil || dto["mchNbr"] != "fixture-merchant" || dto["mchOrdNbr"] != "fixture-order" || dto["pltNbr"] != "fixture-merchant" {
			t.Error("银行参数不符")
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"head":{"rtnCode":"SUC0000"},"body":{"status":"FIN","amount":1.20}}`)
	}))
	defer bank.Close()
	t.Setenv("busbank_host", bank.URL)
	r := httptest.NewRequest("POST", "/saas-capital/chainPayment/orderDetails", strings.NewReader(`{"mchNbr":"fixture-merchant","mchOrdNbr":"fixture-order"}`))
	w := httptest.NewRecorder()
	(&Service{}).payOrderDetails(w, r)
	if !called || w.Code != 200 || !strings.Contains(w.Body.String(), `"status":"FIN"`) || !strings.Contains(w.Body.String(), `1.20`) {
		t.Fatalf("查询未执行真实银行请求或丢失原始结果: called=%v http=%d body=%s", called, w.Code, w.Body.String())
	}
}

func TestBankPaymentPersistsJavaTransaction(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	admin, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := fmt.Sprintf("audit_bank_%d", time.Now().UnixNano())
	if _, err = admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_basic_enterprise", "chain_capital_transinfo", "chain_capital_bkcall", "chain_contract", "chain_auth_user"} {
		if _, err = d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{
		"INSERT INTO chain_basic_enterprise(id,enterprise_name,merchant_number) VALUES(1,'fixture买方','buyer'),(2,'fixture卖方','seller')",
		"INSERT INTO chain_auth_user(id,account,handle_result,level,account_type) VALUES(1,'fixture-user',1,1,1)",
		"INSERT INTO chain_contract(id,contract_code) VALUES(1,'fixture-contract')",
	} {
		if _, err = d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	der, _ := x509.MarshalPKCS8PrivateKey(key)
	t.Setenv("busbank_privateKey", base64.StdEncoding.EncodeToString(der))
	t.Setenv("busbank_appId", "fixture-app")
	t.Setenv("busbank_appsecret", "fixture-secret")
	called := false
	bank := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
		}
		body, _ := req["body"].(map[string]any)
		dto, _ := body["orderInfoPayDto"].(map[string]any)
		switch r.URL.Path {
		case "/order-pay":
			if dto["orderType"] != "PAY" || dto["mchNbr"] != "seller" {
				t.Error("付款请求不符")
			}
			io.WriteString(w, `{"head":{"rtnCode":"SUC0000"},"body":{"payUrl":"fixture-pay-url"}}`)
		case "/order-cfm":
			io.WriteString(w, `{"head":{"rtnCode":"FAIL","rtnText":"BHY2634 fixture"}}`)
		case "/order-bck":
			ref, _ := body["orderInfoBckDto"].(map[string]any)
			if ref["mchOrderNbr"] != "BKfixture-order" || ref["orderSubtype"] != "BCK" || ref["oriOrderAmt"] == nil {
				t.Error("退款参数不符")
			}
			io.WriteString(w, `{"head":{"rtnCode":"SUC0000"}}`)
		case "/bill-info":
			bill, _ := body["billQryDto"].(map[string]any)
			if bill["pageNum"] != "1" || bill["pageSize"] != "3" {
				t.Error("对账单分页类型不符")
			}
			io.WriteString(w, `{"head":{"rtnCode":"SUC0000"},"body":{"list":[1]}}`)
		default:
			t.Error("未知银行路径", r.URL.Path)
		}
	}))
	defer bank.Close()
	t.Setenv("busbank_host", bank.URL)
	s := NewService(&db.DB{Chain: d})
	kp := &auth.KeyPair{Private: key, Public: &key.PublicKey}
	token, _ := kp.Generate(auth.JWTClaims{Account: "fixture-user", UserID: "1", CompanyID: "1", EnterpriseCode: "FIX0001"})
	r := httptest.NewRequest("POST", "/saas-capital/chainPayment/payment", strings.NewReader(`{"payeeEtc":"fixture卖方","orderAmt":12.30,"orderSubtype":"GRT","contractId":1,"contractCode":"fixture-contract"}`))
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	auth.Middleware(kp, d)(http.HandlerFunc(s.payPayment)).ServeHTTP(w, r)
	if !called || !strings.Contains(w.Body.String(), `"data":"fixture-pay-url"`) {
		t.Fatalf("付款未完成: called=%v response=%s", called, w.Body.String())
	}
	var amount, subtype, payName string
	var status, isPay int
	if err = d.QueryRow("SELECT pay_money,orderSubtype,pay_unit_name,status FROM chain_capital_transinfo").Scan(&amount, &subtype, &payName, &status); err != nil {
		t.Fatal(err)
	}
	if err = d.QueryRow("SELECT is_pay FROM chain_contract WHERE id=1").Scan(&isPay); err != nil {
		t.Fatal(err)
	}
	if amount != "12.30" || subtype != "GRT" || payName != "fixture买方" || status != 0 || isPay != 1 {
		t.Fatalf("付款落库与 Java 不同: amount=%s subtype=%s status=%d isPay=%d", amount, subtype, status, isPay)
	}
	if _, err = d.Exec("UPDATE chain_capital_transinfo SET mch_order_nbr='fixture-order',mch_trans_nbr='fixture-serial',bus_status='WCF'"); err != nil {
		t.Fatal(err)
	}
	s.Auth = auth.NewService(d, kp)
	code := s.Auth.IssuePhoneCode("1")
	request := func(path, body string) string {
		r := httptest.NewRequest("POST", "/saas-capital/chainPayment/"+path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		auth.Middleware(kp, d)(http.HandlerFunc(s.handlePayment)).ServeHTTP(w, r)
		return w.Body.String()
	}
	confirm := fmt.Sprintf(`{"code":%q,"id":1,"mchNbr":"seller","mchOrderNbr":"fixture-order","mchTransNbr":"fixture-serial","transType":"CFM","transAmt":12.30}`, code)
	if result := request("confirmPay", confirm); !strings.Contains(result, "银行清算处理中，请24小时后再尝试！") {
		t.Fatalf("确认付款分支不符: %s", result)
	}
	// Java 取消接口收到银行失败仍增加成功数，同时保留 errInfos；不能修正这个计数问题。
	cancel := fmt.Sprintf(`{"code":%q,"refundList":[{"id":1,"mchNbr":"seller","mchOrderNbr":"fixture-order","mchTransNbr":"fixture-serial","transType":"CNL","transAmt":12.30,"payerBankFlag":"Y"}]}`, code)
	if result := request("cancelPay", cancel); !strings.Contains(result, "成功1笔，失败0笔") || !strings.Contains(result, "BHY2634 fixture") {
		t.Fatalf("取消付款未保留 Java 计数: %s", result)
	}
	refund := fmt.Sprintf(`{"code":%q,"refundList":[{"mchNbr":"seller","oriMchOrderNbr":"fixture-order","oriMchTransNbr":"fixture-serial","oriOrderAmt":12.30}]}`, code)
	if result := request("refund", refund); !strings.Contains(result, "操作成功") {
		t.Fatalf("退款不符: %s", result)
	}
	var bkOrder, bkSerial string
	if err = d.QueryRow("SELECT bck_order_nbr,bck_trans_nbr FROM chain_capital_transinfo").Scan(&bkOrder, &bkSerial); err != nil || bkOrder != "BKfixture-order" || bkSerial != "BKfixture-serial" {
		t.Fatalf("退款未更新订单: %v", err)
	}
	if result := request("billInfo", `{"orderDate":"20261119","pageNum":1,"pageSize":3}`); !strings.Contains(result, `"list":[1]`) {
		t.Fatalf("对账单查询不符: %s", result)
	}
	pub, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	t.Setenv("busbank_publickey", base64.StdEncoding.EncodeToString(pub))
	callback := func(message string, success bool) {
		body := `{"msgNbr":"` + message + `","mchOrderNbr":"fixture-order","mchNbr":"seller","status":"FIN","orderType":"PAY","rtnFlag":"S"}`
		hash := sha256.Sum256([]byte(body))
		sig, _ := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
		r := httptest.NewRequest("POST", "/saas-capital/chainPayment/transMsgCallBack", strings.NewReader(body))
		r.Header.Set("CMB-BodySign", base64.StdEncoding.EncodeToString(sig))
		w := httptest.NewRecorder()
		auth.Middleware(kp, d)(http.HandlerFunc(s.handlePayment)).ServeHTTP(w, r)
		if success && (w.Code != 200 || !strings.Contains(w.Body.String(), `"rtnCode":"SUC0000"`) || strings.Contains(w.Body.String(), `"data"`)) {
			t.Fatalf("招行回调未按裸 JSON 应答: %d %s", w.Code, w.Body.String())
		}
		if !success && !strings.Contains(w.Body.String(), `"code":999`) {
			t.Fatalf("SQL失败未保留错误响应: %s", w.Body.String())
		}
	}
	var count int
	if _, err = d.Exec("CREATE TRIGGER reject_callback BEFORE UPDATE ON chain_capital_transinfo FOR EACH ROW SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT='fixture rollback'"); err != nil {
		t.Fatal(err)
	}
	callback("2", false)
	if err = d.QueryRow("SELECT count(*) FROM chain_capital_bkcall").Scan(&count); err != nil || count != 0 {
		t.Fatalf("原 Java REQUIRED 回滚应撤销回调记录，实际 %d %v", count, err)
	}
	if _, err = d.Exec("DROP TRIGGER reject_callback"); err != nil {
		t.Fatal(err)
	}
	callback("2", true)
	callback("2", true)
	callback("1", true)
	if err = d.QueryRow("SELECT count(*) FROM chain_capital_bkcall").Scan(&count); err != nil || count != 1 {
		t.Fatalf("回调消息序号判断不同: count=%d err=%v", count, err)
	}
	if err = d.QueryRow("SELECT status FROM chain_capital_transinfo").Scan(&status); err != nil || status != 2 {
		t.Fatalf("回调未更新交易状态: %d %v", status, err)
	}
}
