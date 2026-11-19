package platform

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTencentRequestUsesRealProviderProtocol(t *testing.T) {
	t.Setenv("TENCENT_SECRET_ID", "fixture-id")
	t.Setenv("TENCENT_SECRET_KEY", "fixture-key")
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Header.Get("X-TC-Action") != "IDCardOCR" || !strings.HasPrefix(r.Header.Get("Authorization"), "TC3-HMAC-SHA256 Credential=fixture-id/") {
			t.Error("腾讯云请求签名或 Action 缺失")
		}
		io.WriteString(w, `{"Response":{"Name":"fixture","RequestId":"fixture-request"}}`)
	}))
	defer server.Close()
	t.Setenv("TENCENT_OCR_ENDPOINT", server.URL)
	result, err := TencentJSON(context.Background(), "ocr", "IDCardOCR", "2018-11-19", "ap-beijing", map[string]any{"ImageBase64": "fixture"})
	if err != nil || !called || result["Name"] != "fixture" {
		t.Fatalf("腾讯云请求失败: %v %v", result, err)
	}
}
