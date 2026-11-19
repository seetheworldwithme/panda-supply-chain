package authadmin

import (
	"net/http"
	"testing"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

func TestLegacyLocalAuthRoutes(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, NewService(&db.DB{}))
	for _, path := range []string{
		"/chainAuthUser/list", "/chainAuthUser/getById", "/chainAuthUser/ptlist", "/chainAuthUser/updHeadImg",
		"/chainAuthUser/listCarAccount", "/chainAuthUser/addCarAccount", "/chainAuthUser/updCarAccount", "/chainAuthUser/resetPwd",
		"/chainAuthUser/selByIds", "/chainAuthUser/isPerfect", "/chainAuthUser/updBatch", "/chainAuthUser/updPwd", "/chainAuthUser/joinCpy",
		"/chainAuthUser/valIdentity", "/chainAuthUser/valCode", "/chainAuthUser/setAccount",
		"/chainAuthUserList/deleteById", "/chainAuthUserList/selEtcRoles", "/chainAuthUserList/selPtRoles", "/chainAuthAuthorityRole/selAuthy",
	} {
		r, _ := http.NewRequest("POST", "/saas-auth"+path, nil)
		if _, pattern := mux.Handler(r); pattern == "" {
			t.Errorf("missing Java mapping %s", path)
		}
	}
}
