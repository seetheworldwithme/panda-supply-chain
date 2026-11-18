package commonrest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCommonFileBespokeRequiresAuth(t *testing.T) {
	mux := http.NewServeMux()
	registerCommonFile(mux, &Service{})
	for _, p := range []string{"deleteById", "updateById", "delBatch", "updBthById", "delByRelId", "compressedFile"} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest("POST", "/saas-common/chainCommonFile/"+p, strings.NewReader(`{}`)))
		if !strings.Contains(w.Body.String(), `"code":403`) {
			t.Fatalf("%s %d %s", p, w.Code, w.Body.String())
		}
	}
}
