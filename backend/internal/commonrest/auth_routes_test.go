package commonrest

import (
	"net/http"
	"testing"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

func TestJavaAuthTemplateRoutes(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, NewService(&db.DB{}))
	for _, name := range []string{"chainAuthAuthenLog", "chainAuthAuthority", "chainAuthAuthorityRole", "chainAuthAuthorityUser", "chainAuthMenuRole", "chainAuthMenuUser", "chainAuthSecret", "chainAuthUseRole"} {
		for _, endpoint := range []string{"list", "listByPage", "getById", "add", "updateById", "deleteById", "delBatch", "exportExcel", "importExcel"} {
			path := "/saas-auth/" + name + "/" + endpoint
			r, _ := http.NewRequest(http.MethodPost, path, nil)
			if _, pattern := mux.Handler(r); pattern == "" {
				t.Errorf("missing Java route %s", path)
			}
		}
	}
}

func TestJavaContractSpreadsheetRoutes(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, NewService(&db.DB{}))
	for _, base := range []string{"chainContract", "chainCgContract", "chainXsContract"} {
		for _, end := range []string{"importExcelCg", "importExcelXs", "exportExcel", "importAlone", "importOnce"} {
			path := "/saas-contract/" + base + "/" + end
			r, _ := http.NewRequest("POST", path, nil)
			if _, pattern := mux.Handler(r); pattern == "" {
				t.Errorf("missing %s", path)
			}
		}
	}
}
