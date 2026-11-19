package warehouse

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func TestRedWithdrawalCallsProviderBeforeSavingOriginalRecords(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	defer admin.Close()
	schema := fmt.Sprintf("audit_red_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_auth_user", "chain_global_configuration", "bus_red_money_clog", "bus_red_envelope_record", "bus_red_envelope_batch", "bus_red_envelope_detail"} {
		if _, err := d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{
		"INSERT INTO chain_auth_user(id,account,handle_result,level,account_type,phone,openid) VALUES(1,'fixture',1,1,1,'00000001122','fixture-openid')",
		"INSERT INTO chain_global_configuration(id,red_uplimit_days) VALUES(1,1000)",
		"INSERT INTO bus_red_money_clog(id,user_id,receive_type,create_amount,is_write_off) VALUES(1,1,7,99,0)",
	} {
		if _, err := d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	der, _ := x509.MarshalPKCS8PrivateKey(key)
	path := filepath.Join(t.TempDir(), "key.pem")
	os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0600)
	cert := &x509.Certificate{SerialNumber: big.NewInt(123), NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	encoded, _ := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	old := platform.LoadWxConfig
	defer func() { platform.LoadWxConfig = old }()
	platform.LoadWxConfig = func() platform.WxLocalConfig {
		return platform.WxLocalConfig{MchID: "fixture", MerchantSerial: "fixture", MerchantPrivateKeyPath: path, PlatformCert: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: encoded}))}
	}
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if r.URL.Path != "/v3/transfer/batches" || body["total_amount"] != float64(7) {
			t.Error("未保留原 Java 客户端提现金额")
		}
		var count int
		d.QueryRow("SELECT count(*) FROM bus_red_envelope_record").Scan(&count)
		if count != 0 {
			t.Error("外部打款之前已写领取记录")
		}
		data := `{"batch_id":"fixture"}`
		stamp, nonce := fmt.Sprint(time.Now().Unix()), "fixture"
		hash := sha256.Sum256([]byte(stamp + "\n" + nonce + "\n" + data + "\n"))
		signature, _ := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
		w.Header().Set("Wechatpay-Timestamp", stamp)
		w.Header().Set("Wechatpay-Nonce", nonce)
		w.Header().Set("Wechatpay-Serial", "7b")
		w.Header().Set("Wechatpay-Signature", base64.StdEncoding.EncodeToString(signature))
		io.WriteString(w, data)
	}))
	defer server.Close()
	t.Setenv("WECHAT_PAY_ENDPOINT", server.URL)
	t.Setenv("PANDA_WAREHOUSE_RED_ENABLED", "1")
	t.Setenv("PANDA_WAREHOUSE_RED_START_TIME", "00:00:00")
	t.Setenv("PANDA_WAREHOUSE_RED_END_TIME", "23:59:59")
	kp := &auth.KeyPair{Private: key, Public: &key.PublicKey}
	token, _ := kp.Generate(auth.JWTClaims{Account: "fixture", UserID: "1"})
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	r := httptest.NewRequest("POST", "/saas-warehouse/busRedEnvelopeRecord/receiveAmount", strings.NewReader(`{"userId":1,"receiveType":7,"receiveAmount":7,"clogId":1}`))
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	auth.Middleware(kp, d)(http.HandlerFunc(s.rrReceiveAmount)).ServeHTTP(w, r)
	if !called || !strings.Contains(w.Body.String(), `"code":200`) {
		t.Fatalf("红包提现未完成: called=%v response=%s", called, w.Body.String())
	}
	for _, table := range []string{"bus_red_envelope_record", "bus_red_envelope_batch", "bus_red_envelope_detail"} {
		var count int
		if err := d.QueryRow("SELECT count(*) FROM " + table).Scan(&count); err != nil || count != 1 {
			t.Fatalf("%s 未落库: %d %v", table, count, err)
		}
	}
	var used int
	if err := d.QueryRow("SELECT is_write_off FROM bus_red_money_clog WHERE id=1").Scan(&used); err != nil || used != 1 {
		t.Fatalf("生成流水未核销: %d %v", used, err)
	}
}
