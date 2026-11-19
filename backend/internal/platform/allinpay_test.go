package platform

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAllinpayUsesServerSDKSM2Protocol(t *testing.T) {
	openssl, err := exec.LookPath("openssl")
	if err != nil {
		t.Fatal("通联 SM2 需要原生 OpenSSL", err)
	}
	dir := t.TempDir()
	key, cert, pfx := filepath.Join(dir, "key.pem"), filepath.Join(dir, "cert.pem"), filepath.Join(dir, "merchant.pfx")
	for _, args := range [][]string{
		{"genpkey", "-algorithm", "SM2", "-out", key},
		{"req", "-new", "-x509", "-key", key, "-sm3", "-subj", "/CN=fixture", "-days", "1", "-out", cert},
		{"pkcs12", "-legacy", "-export", "-inkey", key, "-in", cert, "-passout", "pass:fixture", "-out", pfx},
	} {
		if output, err := exec.Command(openssl, args...).CombinedOutput(); err != nil {
			t.Fatalf("临时证书生成失败: %v %s", err, output)
		}
	}
	t.Setenv("tlPfxPath", pfx)
	t.Setenv("tlCertPath", cert)
	t.Setenv("TL_SM2_PFX_PASSWORD", "fixture")
	rawServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusBadGateway)
		}
		io.WriteString(w, "provider-unavailable")
	}))
	defer rawServer.Close()
	for _, path := range []string{"/plain", "/error"} {
		_, _, raw, err := AllinpayPost(context.Background(), rawServer.URL+path, "fixture-inst", "2089", nil, false)
		if err != nil || !strings.Contains(raw, "provider-unavailable") {
			t.Fatalf("原 Java 2089 不解析响应，HTTP/坏 XML 也应原样返回: %q %v", raw, err)
		}
		if path == "/error" && !strings.HasPrefix(raw, "status:HTTP/1.1 502 Bad Gateway@@headers:") {
			t.Fatalf("未保留 SDK HTTP 错误文本包装: %q", raw)
		}
	}
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		xml := r.FormValue("msgPlain")
		if r.Header.Get("allinpay-signType") != "SM2" || r.Header.Get("allinpay-certSn") == "" || !strings.Contains(xml, "<processing_code>3001</processing_code>") {
			t.Error("通联头或参数不符")
		}
		if err := VerifyAllinpay(context.Background(), xml); err != nil {
			t.Error("请求签名不符", err)
		}
		plain := `<?xml version="1.0" encoding="GBK"?><transaction><head><processing_code>3001</processing_code></head><response><resp_code>0000</resp_code><qur_rst>1122</qur_rst></response></transaction>`
		cmd := exec.Command(openssl, "dgst", "-sm3", "-sign", key, "-sigopt", "distid:1234567812345678")
		cmd.Stdin = strings.NewReader(plain)
		sig, err := cmd.Output()
		if err != nil {
			t.Error(err)
			return
		}
		io.WriteString(w, strings.Replace(plain, "</head>", "<sign_code>"+base64.StdEncoding.EncodeToString(sig)+"</sign_code></head>", 1))
	}))
	defer server.Close()
	response, _, raw, err := AllinpayPost(context.Background(), server.URL, "fixture-inst", "3001", map[string]any{"org_order_num": "fixture-order"}, true)
	if err != nil || !called || response["qur_rst"] != "1122" || !strings.Contains(raw, "<resp_code>0000</resp_code>") {
		t.Fatalf("通联请求未完成: called=%v response=%v err=%v", called, response, err)
	}
	if err = VerifyAllinpay(context.Background(), strings.Replace(raw, "1122", "0000", 1)); err == nil {
		t.Fatal("篡改的通联响应未拒绝")
	}
	if _, err = os.Stat(pfx); err != nil {
		t.Fatal("不应修改输入证书")
	}
}
