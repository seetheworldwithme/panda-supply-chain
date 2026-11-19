package auth

import (
	"image/jpeg"
	"net/http/httptest"
	"testing"
)

func TestImageCodeJPEGDimensionsAndCache(t *testing.T) {
	s := NewService(nil, nil)
	for _, query := range []string{"", "?width=90&height=30"} {
		r := httptest.NewRequest("GET", "/saas-auth/code/imageGet"+query, nil)
		w := httptest.NewRecorder()
		s.ImageCode(w, r)
		cfg, err := jpeg.DecodeConfig(w.Body)
		if err != nil || w.Header().Get("Content-Type") != "image/jpeg" {
			t.Fatalf("验证码必须为 JPEG: %v", err)
		}
		width, height := 67, 23
		if query != "" {
			width, height = 90, 30
		}
		if cfg.Width != width || cfg.Height != height {
			t.Fatalf("尺寸不一致: %+v", cfg)
		}
		if w.Header().Get("Cache-Control") != "no-store, no-cache" {
			t.Fatal("缺少禁止缓存头")
		}
		if sc := s.sms[imageCodeKey(r)]; len(sc.value) != 4 {
			t.Fatal("缺少四位验证码缓存")
		}
	}
}
