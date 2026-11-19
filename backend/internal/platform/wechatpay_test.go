package platform

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestWechatV3UsesSignedRequestAndVerifiesResponse(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	der, _ := x509.MarshalPKCS8PrivateKey(key)
	path := filepath.Join(t.TempDir(), "key.pem")
	os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0600)
	cert := &x509.Certificate{SerialNumber: big.NewInt(123), Subject: pkix.Name{CommonName: "fixture"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	encoded, _ := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	old := LoadWxConfig
	defer func() { LoadWxConfig = old }()
	LoadWxConfig = func() WxLocalConfig {
		return WxLocalConfig{MchID: "fixture-merchant", MerchantSerial: "fixture-serial", MerchantPrivateKeyPath: path, PlatformCert: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: encoded}))}
	}
	stamp := strconv.FormatInt(time.Now().Unix(), 10)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "WECHATPAY2-SHA256-RSA2048 ") {
			t.Error("微信商户签名缺失")
		}
		body := `{"prepay_id":"fixture-prepay"}`
		timestamp, nonce := stamp, "fixture-nonce"
		hash := sha256.Sum256([]byte(timestamp + "\n" + nonce + "\n" + body + "\n"))
		sig, _ := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
		w.Header().Set("Wechatpay-Timestamp", timestamp)
		w.Header().Set("Wechatpay-Nonce", nonce)
		w.Header().Set("Wechatpay-Serial", "7B")
		w.Header().Set("Wechatpay-Signature", base64.StdEncoding.EncodeToString(sig))
		io.WriteString(w, body)
	}))
	defer server.Close()
	t.Setenv("WECHAT_PAY_ENDPOINT", server.URL)
	response, err := WechatV3JSON(context.Background(), "POST", "/v3/pay/transactions/app", map[string]any{"fixture": true})
	if err != nil || response["prepay_id"] != "fixture-prepay" {
		t.Fatalf("微信支付链路未完成: %v %v", response, err)
	}
	stamp = strconv.FormatInt(time.Now().Add(-6*time.Minute).Unix(), 10)
	if _, err := WechatV3JSON(context.Background(), "POST", "/v3/pay/transactions/app", nil); err == nil {
		t.Fatal("原 SDK 拒绝超出五分钟的响应，Go 未拒绝")
	}
}
