package auth

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAnonAllowedEnv(t *testing.T) {
	t.Setenv("PANDA_ANON_ALLOW", " /saas-warehouse/busProduct/,/x/y ")
	for path, want := range map[string]bool{
		"/saas-auth/chainAuthUser/ptlist":          false,
		"/saas-auth/chainAuthChat/getAccessToken":  false,
		"/saas-auth/chainAuthChat/getAccessTokenX": true, // 登录族前缀
		"/saas-auth/chainAuthUser/setAccount":      true,
		"/saas-warehouse/busProduct/listByPage":    true,
		"/x/y":                                     true,
		"/x/yz":                                    false,
	} {
		if got := anonAllowed(path); got != want {
			t.Fatalf("anonAllowed(%s) = %v, want %v", path, got, want)
		}
	}
}

// 公开签名(旧前端 JS 包里写死的那一对请求头)只能访问登录前接口;
// 密钥未公开的其他 AuthorizationId 保持原行为。
func TestPublicSignatureOnlyReachesPreLoginEndpoints(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	d, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/chain?parseTime=true&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	sig, err := NewService(d, nil).AnonSignature()
	if err != nil {
		t.Fatal(err)
	}
	var partner, secret string
	if err := d.QueryRow("SELECT authorization_id, secret FROM chain_auth_secret WHERE authorization_id<>? LIMIT 1", publicAuthID).Scan(&partner, &secret); err != nil {
		t.Skip("本机库没有第二方授权凭据: " + err.Error())
	}
	called := false
	h := Middleware(&KeyPair{}, d)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }))
	for _, c := range []struct {
		id, sig, path string
		want          bool
	}{
		{publicAuthID, sig, "/saas-auth/chainAuthUser/ptlist", false},
		{publicAuthID, sig, "/saas-auth/chainAuthSecret/list", false},
		{publicAuthID, sig, "/saas-auth/chainAuthChat/getAccessToken", false},
		{publicAuthID, sig, "/saas-auth/chainAuthUser/getPhoneCode", true},
		{publicAuthID, sig, "/saas-auth/chainAuthUser/setAccount", true},
		{publicAuthID, sig, "/saas-auth/chainAuthChat/loginByAccount", true},
		{publicAuthID, sig, "/saas-common/upload/imageShow", true},
		{publicAuthID, sig, "/saas-common/upload/downById?fileName=1", true}, // 页面用地址签名预览 / 下载附件
		{publicAuthID, sig, "/saas-common/upload/previewStream", true},
		{publicAuthID, sig, "/saas-warehouse/busProduct/exportExcel", true},
		{publicAuthID, sig, "/saas-warehouse/busProduct/getById", false},
		{partner, hmacSHA256Hex(partner+"hyqc"+secret, secret), "/saas-auth/chainAuthUser/ptlist", true},
		{publicAuthID, "bad", "/saas-auth/chainAuthUser/getPhoneCode", false},
	} {
		called = false
		r := httptest.NewRequest("POST", c.path, nil)
		r.Header.Set("AuthorizationId", c.id)
		r.Header.Set("Signature", c.sig)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if called != c.want {
			t.Errorf("AuthorizationId=%s Signature=%s %s: reached handler=%v, want %v (%s)", c.id, c.sig[:8], c.path, called, c.want, w.Body.String())
		}
	}
}
