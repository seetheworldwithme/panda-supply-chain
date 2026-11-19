package platform

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEsignPreservesJavaGetAndPostSignature(t *testing.T) {
	t.Setenv("appId", "fixture-app")
	t.Setenv("appKey", "fixture-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		sum := "{}"
		if r.Method == "POST" {
			hash := md5.Sum(body)
			sum = base64.StdEncoding.EncodeToString(hash[:])
		}
		text := r.Method + "\n*/*\n" + sum + "\napplication/json; charset=UTF-8\n\n" + r.URL.RequestURI()
		mac := hmac.New(sha256.New, []byte("fixture-key"))
		mac.Write([]byte(text))
		if r.Header.Get("Content-MD5") != sum || r.Header.Get("X-Tsign-Open-Ca-Signature") != base64.StdEncoding.EncodeToString(mac.Sum(nil)) {
			t.Error("e签宝签名协议不符")
		}
		io.WriteString(w, `{"code":0,"data":{"flowId":"fixture-flow"}}`)
	}))
	defer server.Close()
	t.Setenv("host", server.URL)
	for _, method := range []string{"GET", "POST"} {
		result := EsignJSON(context.Background(), method, "/fixture?a=b", map[string]any{"fixture": true})
		if result["data"] == nil {
			t.Fatalf("%s 未返回真实签章结果", method)
		}
	}
}
