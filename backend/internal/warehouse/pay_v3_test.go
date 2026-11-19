package warehouse

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
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

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func TestMembershipV3PreservesOriginalJSAndAppParameters(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	der, _ := x509.MarshalPKCS8PrivateKey(key)
	path := filepath.Join(t.TempDir(), "key.pem")
	os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0600)
	cert := &x509.Certificate{SerialNumber: big.NewInt(123), NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	encoded, _ := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	old := platform.LoadWxConfig
	defer func() { platform.LoadWxConfig = old }()
	platform.LoadWxConfig = func() platform.WxLocalConfig {
		return platform.WxLocalConfig{MchID: "fixture-merchant", NotifyURL: "https://fixture.invalid", MerchantSerial: "fixture-serial", MerchantPrivateKeyPath: path, PlatformCert: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: encoded}))}
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m map[string]any
		json.NewDecoder(r.Body).Decode(&m)
		amount, _ := m["amount"].(map[string]any)
		if strings.HasSuffix(r.URL.Path, "/jsapi") {
			payer, _ := m["payer"].(map[string]any)
			if amount["total"] != float64(1) || payer["openid"] != "openId" || m["out_trade_no"] == "fixture-order" {
				t.Error("未保留原 JS 支付常量问题")
			}
		} else if amount["total"] != float64(1234) || m["out_trade_no"] != "fixture-order" || m["appid"] != "wx0ff59ec3274e0a53" {
			t.Error("App 支付参数不符")
		}
		body := `{"prepay_id":"fixture-prepay"}`
		stamp, nonce := fmt.Sprint(time.Now().Unix()), "fixture-nonce"
		hash := sha256.Sum256([]byte(stamp + "\n" + nonce + "\n" + body + "\n"))
		sig, _ := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
		w.Header().Set("Wechatpay-Timestamp", stamp)
		w.Header().Set("Wechatpay-Nonce", nonce)
		w.Header().Set("Wechatpay-Serial", "7b")
		w.Header().Set("Wechatpay-Signature", base64.StdEncoding.EncodeToString(sig))
		io.WriteString(w, body)
	}))
	defer server.Close()
	t.Setenv("WECHAT_PAY_ENDPOINT", server.URL)
	s := &Service{}
	r := httptest.NewRequest("POST", "/", nil)
	for _, m := range []map[string]any{{"isJsPay": true}, {"appid": "ignored-original-field"}} {
		result, err := s.v3MembershipPay(r, m, "12.34", "fixture-order")
		if err != nil || result["packageVal"] == nil {
			t.Fatalf("会员 v3 支付未完成: %v %v", result, err)
		}
	}
}
