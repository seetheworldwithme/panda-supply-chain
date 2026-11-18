package commonrest

import (
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
	"net/http"
	"strings"
)

// Register adds the source-controller route difference left by phases 1–5.
// Existing handlers keep ownership; common metadata still serves their Excel additions.
func Register(mux *http.ServeMux, s *Service) {
	reg := func(path string, h http.HandlerFunc) {
		r, _ := http.NewRequest("POST", path, nil)
		_, pattern := mux.Handler(r)
		if pattern == "" {
			mux.HandleFunc("POST "+path, authGuard(h))
		}
	}
	for name, e := range ents {
		path := "/saas-common/" + strings.ToLower(name[:1]) + name[1:]
		reg(path+"/list", func(w http.ResponseWriter, r *http.Request) { s.tplList(w, r, e) })
		reg(path+"/listByPage", func(w http.ResponseWriter, r *http.Request) {
			family := byte('A')
			if strings.HasPrefix(name, "ChainTlzf") || name == "BusMqConsumeLog" {
				family = 'C'
			}
			s.tplListByPage(w, r, e, family)
		})
		reg(path+"/getById", func(w http.ResponseWriter, r *http.Request) { s.tplGetById(w, r, e) })
		if name != "ChainBasicEnterprise" {
			reg(path+"/add", func(w http.ResponseWriter, r *http.Request) { s.tplAdd(w, r, e) })
		}
		if name != "ChainCommonFile" {
			reg(path+"/updateById", func(w http.ResponseWriter, r *http.Request) { s.tplUpdateById(w, r, e) })
			reg(path+"/deleteById", func(w http.ResponseWriter, r *http.Request) { s.tplDeleteById(w, r, e) })
			reg(path+"/delBatch", func(w http.ResponseWriter, r *http.Request) { s.tplDelBatch(w, r, e, true) })
		}
		reg(path+"/importExcel", func(w http.ResponseWriter, r *http.Request) { s.tplImportExcel(w, r, e) })
		reg(path+"/exportExcel", func(w http.ResponseWriter, r *http.Request) { s.tplExportExcel(w, r, e) })
		if strings.HasPrefix(name, "ChainTlzf") || name == "ChainBasicEnterprise" || name == "BusMqConsumeLog" {
			reg(path+"/updBatch", func(w http.ResponseWriter, r *http.Request) { s.tplUpdBatch(w, r, e) })
		}
	}
	registerBespoke(mux, s)
	registerEnterprise(mux, s)
	registerExternal(mux, s)
	registerCommonFile(mux, s)
}
func entityForTable(table string) *ent {
	for _, e := range ents {
		if e.table == table {
			return e
		}
	}
	return &ent{table: table}
}
func (s *Service) insert(r *http.Request, e *ent, vals map[string]any) (int64, error) {
	if e.fill {
		return platform.Insert(s.pick(r), e.table, vals)
	}
	return platform.InsertNoFill(s.pick(r), e.table, vals)
}
func (s *Service) update(r *http.Request, e *ent, id any, vals map[string]any) (int64, error) {
	if e.fill {
		return platform.UpdateByID(s.pick(r), e.table, id, vals)
	}
	return platform.UpdateByIDNoFill(s.pick(r), e.table, id, vals)
}
