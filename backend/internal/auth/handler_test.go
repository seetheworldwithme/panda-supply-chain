package auth

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"testing"
)

type loginDriver struct{}
type loginConn struct{}
type loginRows struct{ read bool }

func (loginDriver) Open(string) (driver.Conn, error)  { return loginConn{}, nil }
func (loginConn) Prepare(string) (driver.Stmt, error) { panic("unexpected Prepare") }
func (loginConn) Close() error                        { return nil }
func (loginConn) Begin() (driver.Tx, error)           { panic("unexpected Begin") }
func (loginConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &loginRows{}, nil
}
func (*loginRows) Columns() []string { return []string{"count"} }
func (*loginRows) Close() error      { return nil }
func (r *loginRows) Next(dest []driver.Value) error {
	if r.read {
		return io.EOF
	}
	r.read = true
	dest[0] = int64(0)
	return nil
}

func init() { sql.Register("auth-login-contract-test", loginDriver{}) }

func TestLoginPreservesNullableProfileFields(t *testing.T) {
	d, err := sql.Open("auth-login-contract-test", "")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	kp, err := LoadOrGenerateKeyPair(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := NewService(d, kp)
	u := &chainAuthUser{ID: 1, Account: "fixture", Level: 2}
	m, err := s.loginTokens(u, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"nickName", "wxHeadImg", "phone", "xzqhCode", "xzqhName", "accountType"} {
		if value, exists := m[key]; !exists || value != nil {
			t.Errorf("Java nullable field %s: got %#v, exists=%v; want null", key, value, exists)
		}
	}
	u.Phone = sql.NullString{String: "", Valid: true}
	u.NickName = sql.NullString{String: "stored nickname", Valid: true}
	u.WxHeadImg = sql.NullString{String: "stored image", Valid: true}
	m, err = s.loginTokens(u, "")
	if err != nil {
		t.Fatal(err)
	}
	if m["phone"] != "" {
		t.Errorf("empty string must remain distinct from null: %#v", m["phone"])
	}
	if m["nickName"] != nil || m["wxHeadImg"] != nil {
		t.Fatal("Java UserServiceImpl does not populate UserInfo nickname/avatar")
	}
}
