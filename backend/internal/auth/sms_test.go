package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 验证码缓存只增不删会被未鉴权请求撑大:写入新码时必须清掉过期的。
func TestStoreCodePrunesExpiredCodes(t *testing.T) {
	s := NewService(nil, nil)
	s.storeCode("gone", "1")
	s.sms["gone"] = smsCode{value: "1", expire: time.Now().Add(-time.Second)}
	s.storeCode("fresh", "2")
	if _, ok := s.sms["gone"]; ok {
		t.Fatal("过期验证码未清理")
	}
	if s.sms["fresh"].value != "2" {
		t.Fatal("新验证码未写入")
	}
}

// 图形验证码画布封顶:未鉴权接口不能靠 width/height 申请大内存。
func TestImageCodeRejectsOversizedCanvas(t *testing.T) {
	s := NewService(nil, nil)
	w := httptest.NewRecorder()
	s.ImageCode(w, httptest.NewRequest("GET", "/saas-auth/code/imageGet?width=4096&height=4096", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("超大画布未被拒绝: %d", w.Code)
	}
}
