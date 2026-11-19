package platform

import (
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
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
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestWechatPlatformCertificateBootstrapsAndRotates(t *testing.T) {
	merchant, _ := rsa.GenerateKey(rand.Reader, 2048)
	der, _ := x509.MarshalPKCS8PrivateKey(merchant)
	path := filepath.Join(t.TempDir(), "merchant.pem")
	if err := os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0600); err != nil {
		t.Fatal(err)
	}
	caKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	ca := &x509.Certificate{SerialNumber: big.NewInt(1), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	caDER, _ := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
	oldRoot := wechatCertificateRoot
	wechatCertificateRoot = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER}))
	defer func() { wechatCertificateRoot = oldRoot }()
	platformKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	certs := map[string]string{}
	for _, serial := range []int64{101, 102} {
		cert := &x509.Certificate{SerialNumber: big.NewInt(serial), KeyUsage: x509.KeyUsageDigitalSignature, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
		encoded, err := x509.CreateCertificate(rand.Reader, cert, ca, &platformKey.PublicKey, caKey)
		if err != nil {
			t.Fatal(err)
		}
		certs[fmt.Sprintf("%X", serial)] = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: encoded}))
	}
	cfg := WxLocalConfig{MchID: "fixture", MerchantSerial: "fixture", MerchantPrivateKeyPath: path, APIV3Key: "01234567890123456789012345678901"}
	oldConfig := LoadWxConfig
	LoadWxConfig = func() WxLocalConfig { return cfg }
	defer func() { LoadWxConfig = oldConfig }()
	serial := "65"
	downloads := 0
	tamper := false
	unavailable := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "WECHATPAY2-SHA256-RSA2048 ") {
			t.Error("缺少商户签名")
		}
		attributes := map[string]string{}
		for _, match := range regexp.MustCompile(`([a-z_]+)="([^"]*)"`).FindAllStringSubmatch(r.Header.Get("Authorization"), -1) {
			attributes[match[1]] = match[2]
		}
		requestBody, _ := io.ReadAll(r.Body)
		requestDigest := sha256.Sum256([]byte(r.Method + "\n" + r.URL.RequestURI() + "\n" + attributes["timestamp"] + "\n" + attributes["nonce_str"] + "\n" + string(requestBody) + "\n"))
		requestSignature, _ := base64.StdEncoding.DecodeString(attributes["signature"])
		if err := rsa.VerifyPKCS1v15(&merchant.PublicKey, crypto.SHA256, requestDigest[:], requestSignature); err != nil {
			t.Error("商户请求签名不符", err)
		}
		if r.URL.Path == "/v3/certificates" {
			downloads++
			if unavailable {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			if r.URL.Query().Get("algorithm_type") != "RSA" {
				t.Error("未沿用 SDK RSA 证书接口")
			}
			block, _ := aes.NewCipher([]byte(cfg.APIV3Key))
			gcm, _ := cipher.NewGCM(block)
			nonce, associated := "fixture00001", "certificate"
			data := gcm.Seal(nil, []byte(nonce), []byte(certs[serial]), []byte(associated))
			if tamper {
				data[len(data)-1] ^= 1
			}
			json.NewEncoder(w).Encode(map[string]any{"data": []any{map[string]any{"serial_no": serial, "encrypt_certificate": map[string]any{"algorithm": "AEAD_AES_256_GCM", "nonce": nonce, "associated_data": associated, "ciphertext": base64.StdEncoding.EncodeToString(data)}}}})
			return
		}
		body := `{"prepay_id":"fixture-prepay"}`
		stamp, nonce := fmt.Sprint(time.Now().Unix()), "fixture"
		digest := sha256.Sum256([]byte(stamp + "\n" + nonce + "\n" + body + "\n"))
		signature, _ := rsa.SignPKCS1v15(rand.Reader, platformKey, crypto.SHA256, digest[:])
		w.Header().Set("Wechatpay-Timestamp", stamp)
		w.Header().Set("Wechatpay-Nonce", nonce)
		w.Header().Set("Wechatpay-Serial", serial)
		w.Header().Set("Wechatpay-Signature", base64.StdEncoding.EncodeToString(signature))
		io.WriteString(w, body)
	}))
	defer server.Close()
	t.Setenv("WECHAT_PAY_ENDPOINT", server.URL)
	for _, next := range []string{"65", "66", "66"} {
		serial = next
		out, err := WechatV3JSON(context.Background(), "POST", "/v3/pay/transactions/app", nil)
		if err != nil || out["prepay_id"] != "fixture-prepay" {
			t.Fatalf("证书 %s 轮换失败: %v %v", next, out, err)
		}
	}
	if downloads != 2 {
		t.Fatalf("应初始化及轮换各下载一次，实际 %d", downloads)
	}
	unavailable = true
	wechatCertificates.refreshed = time.Now().Add(-2 * time.Hour)
	if _, err := WechatPlatformCertificate(context.Background(), cfg, "66"); err != nil {
		t.Fatal("原 SDK 刷新失败会保留已有证书", err)
	}
	unavailable = false
	if _, err := WechatPlatformCertificate(context.Background(), cfg, "66"); err != nil {
		t.Fatal("回调证书选择不可复用", err)
	}
	if _, err := WechatPlatformCertificate(context.Background(), cfg, "FFFF"); err == nil {
		t.Fatal("未知序列号未拒绝")
	}
	tamper = true
	if _, err := downloadWechatCertificates(context.Background(), cfg); err == nil {
		t.Fatal("篡改的证书密文未拒绝")
	}
	tamper = false
	wechatCertificateRoot = oldRoot
	if _, err := downloadWechatCertificates(context.Background(), cfg); err == nil {
		t.Fatal("非财付通 CA 签发的证书未拒绝")
	}
}

func TestWxConfigReloadsPrivateFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wechat.json")
	t.Setenv("WECHAT_LOCAL_CONFIG", path)
	for _, serial := range []string{"one", "two"} {
		if err := os.WriteFile(path, []byte(`{"merchantSerial":"`+serial+`"}`), 0600); err != nil {
			t.Fatal(err)
		}
		if got := LoadWxConfig().MerchantSerial; got != serial {
			t.Fatalf("运行配置未重读: %q", got)
		}
	}
}
