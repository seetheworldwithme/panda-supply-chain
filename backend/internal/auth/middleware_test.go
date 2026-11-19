package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddlewareRejectsAnonymousProtectedRequests(t *testing.T) {
	called := false
	h := Middleware(&KeyPair{}, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("POST", "/saas-capital/chainCapitalDetailed/list", strings.NewReader(`{}`)))
	if called || w.Code != 200 || !strings.Contains(w.Body.String(), `"code":403`) {
		t.Fatalf("protected handler called=%v, response=%d %s", called, w.Code, w.Body.String())
	}
}

func TestMiddlewareAllowsJavaPublicLoginAndCodePaths(t *testing.T) {
	for _, path := range []string{"/saas-auth/login", "/saas-auth/screenLogin", "/saas-auth/code/image", "/healthz"} {
		called := false
		h := Middleware(&KeyPair{}, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", path, nil))
		if !called {
			t.Errorf("public request rejected: %s", path)
		}
	}
}

func TestMiddlewareRejectsTokenWhenAccountCannotBeValidated(t *testing.T) {
	kp, err := LoadOrGenerateKeyPair(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	token, err := kp.Generate(JWTClaims{Account: "deleted-fixture", UserID: "1"})
	if err != nil {
		t.Fatal(err)
	}
	called := false
	h := Middleware(kp, nil)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }))
	r := httptest.NewRequest("POST", "/saas-auth/chainAuthMenu/list", strings.NewReader(`{}`))
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if called || !strings.Contains(w.Body.String(), `"code":403`) {
		t.Fatal("无法复核账号时仍接受了 token")
	}
}
