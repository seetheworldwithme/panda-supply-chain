package commonrest_test

import (
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/authadmin"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/bill"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/capital"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/commonrest"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/contract"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/dashboard"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/logistics"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/warehouse"
	"net/http"
	"testing"
)

func TestAllServiceRoutesRegisterWithoutConflict(t *testing.T) {
	mux := http.NewServeMux()
	d := &db.DB{}
	authadmin.Register(mux, authadmin.NewService(d))
	dashboard.Register(mux, dashboard.NewService(d))
	contract.Register(mux, contract.NewService(d))
	logistics.Register(mux, logistics.NewService(d))
	capital.Register(mux, capital.NewService(d))
	bill.Register(mux, bill.NewService(d))
	warehouse.Register(mux, warehouse.NewService(d))
	commonrest.Register(mux, commonrest.NewService(d))
	for _, path := range []string{"/saas-common/chainMsg/list", "/saas-common/chainCommonTempmx/list", "/saas-common/chainBasicEnterprise/getById", "/saas-common/chainTlzfOrder/selSfEtc"} {
		_, pattern := mux.Handler(newRequest(path))
		if pattern == "" {
			t.Errorf("missing %s", path)
		}
	}
}
func newRequest(path string) *http.Request { r, _ := http.NewRequest("POST", path, nil); return r }
