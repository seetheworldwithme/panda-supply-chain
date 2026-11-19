package commonrest

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEnterpriseSearchCallsOriginalQichachaAPI(t *testing.T) {
	t.Setenv("QCC_APP_KEY", "fixture-app")
	t.Setenv("QCC_SECRET_KEY", "fixture-secret")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sum := md5.Sum([]byte("fixture-app" + r.Header.Get("Timespan") + "fixture-secret"))
		if r.Header.Get("Token") != strings.ToUpper(hex.EncodeToString(sum[:])) || r.URL.Query().Get("searchName") != "fixture" {
			t.Error("企查查请求不符")
		}
		io.WriteString(w, `{"Status":"200","Result":{"Data":[{"Name":"fixture"}]}}`)
	}))
	defer server.Close()
	t.Setenv("QCC_ENDPOINT", server.URL)
	w := httptest.NewRecorder()
	(&Service{}).enterpriseSearch(w, httptest.NewRequest("POST", "/saas-common/chainBasicEnterprise/selEtcName", strings.NewReader(`{"etcName":"fixture"}`)))
	if !strings.Contains(w.Body.String(), `"Name":"fixture"`) {
		t.Fatalf("企业搜索未返回真实数据: %s", w.Body.String())
	}
}
