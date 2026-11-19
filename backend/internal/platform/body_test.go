package platform

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJavaRequiredBodyRejectsMalformedAndWrongShape(t *testing.T) {
	for _, body := range []string{"", "null", "{bad", "[]", "1", "true"} {
		called := false
		handler := ValidateJavaBodies(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest("POST", "/saas-auth/chainAuthUser/list", strings.NewReader(body)))
		if called || w.Code != 500 {
			t.Errorf("body %q: handler=%v status=%d", body, called, w.Code)
		}
	}
}

func TestBodyValidationMatchesJavaFirstJSONValue(t *testing.T) {
	called := false
	h := ValidateJavaBodies(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/saas-auth/chainAuthUser/list", strings.NewReader(`{} {}`)))
	if !called {
		t.Fatal("Java 默认只读取第一个 JSON 值")
	}
}
func TestBodyValidationLeavesValidBodyReadable(t *testing.T) {
	w := httptest.NewRecorder()
	h := ValidateJavaBodies(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m map[string]any
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			t.Fatal(err)
		}
		if m["id"] != float64(1) {
			t.Fatal(m)
		}
	}))
	h.ServeHTTP(w, httptest.NewRequest("POST", "/saas-auth/chainAuthUser/list", strings.NewReader(`{"id":1}`)))
}
func TestBodyValidationDoesNotInventExportRequestBody(t *testing.T) {
	called := false
	h := ValidateJavaBodies(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/saas-auth/chainAuthUser/exportExcel", nil))
	if !called {
		t.Fatal("Java export has no @RequestBody")
	}
}
