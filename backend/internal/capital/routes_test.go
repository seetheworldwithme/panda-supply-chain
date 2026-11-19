package capital

import (
	"net/http"
	"testing"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

func TestDetailedAliasesExposeJavaBalancingAccounts(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, NewService(&db.DB{}))
	for _, base := range []string{"chainCapitalDetailed", "chainCapitalCgDetailed", "chainCapitalXsDetailed"} {
		path := "/saas-capital/" + base + "/balancingAccounts"
		r, _ := http.NewRequest(http.MethodPost, path, nil)
		if _, pattern := mux.Handler(r); pattern == "" {
			t.Errorf("Java mapping missing: POST %s", path)
		}
	}
}
