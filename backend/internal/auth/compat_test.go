package auth

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMobileLoginPreservesJavaDisabledCodeValidation(t *testing.T) {
	d, _ := sql.Open("auth-mobile-contract-test", "")
	defer d.Close()
	kp, err := LoadOrGenerateKeyPair(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := NewService(d, kp)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/saas-auth/authentication/mobile", strings.NewReader("phone=fixture"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.MobileLogin(w, r)
	if !strings.Contains(w.Body.String(), `"code":200`) || !strings.Contains(w.Body.String(), `"token"`) {
		t.Fatalf("Java 未启用表单短信校验，Go 应保留原行为: %s", w.Body.String())
	}
}

type mobileDriver struct{}
type mobileConn struct{ loginConn }
type mobileRows struct{ read bool }

func (mobileDriver) Open(string) (driver.Conn, error) { return mobileConn{}, nil }
func (mobileConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &mobileRows{}, nil
}
func (*mobileRows) Columns() []string { return strings.Split(userCols, ", ") }
func (*mobileRows) Close() error      { return nil }
func (r *mobileRows) Next(dest []driver.Value) error {
	if r.read {
		return io.EOF
	}
	r.read = true
	copy(dest, []driver.Value{int64(1), "fixture", "fixture", "", nil, int64(1), nil, nil, "1", "fixture", nil, nil})
	return nil
}
func init() { sql.Register("auth-mobile-contract-test", mobileDriver{}) }

func TestLogoutOmitsDataLikeJavaSecurityHandler(t *testing.T) {
	w := httptest.NewRecorder()
	NewService(nil, nil).Logout(w, httptest.NewRequest("POST", "/saas-auth/logout", nil))
	if strings.Contains(w.Body.String(), `"data"`) || !strings.Contains(w.Body.String(), `"code":200`) {
		t.Fatalf("unexpected Java logout shape: %s", w.Body.String())
	}
	if w.Header().Get("Content-Type") != "text/json;charset=utf-8" {
		t.Fatal(w.Header())
	}
}
