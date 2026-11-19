package commonrest

import (
	"net/http"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// registerAuthRest 补齐权限模板接口；已有专用处理保留，用户和菜单仅补 Excel 入口。
func registerAuthRest(mux *http.ServeMux, s *Service) {
	s.registerChatLocal(mux)
	s.registerChatWeixin(mux)
	reg := func(path string, h http.HandlerFunc) {
		r, _ := http.NewRequest(http.MethodPost, path, nil)
		if _, pattern := mux.Handler(r); pattern == "" {
			mux.HandleFunc("POST "+path, authGuard(h))
		}
	}
	for name, e := range authRestEnts {
		if name == "ChainAuthChatAuthorization" {
			continue
		}
		if !strings.HasPrefix(name, "ChainAuth") && !strings.HasPrefix(name, "ChainContract") {
			continue
		}
		base := "/saas-auth/" + strings.ToLower(name[:1]) + name[1:]
		if strings.HasPrefix(name, "ChainContract") {
			base = "/saas-contract/" + strings.ToLower(name[:1]) + name[1:]
		}
		reg(base+"/exportExcel", func(w http.ResponseWriter, r *http.Request) { s.tplExportExcel(w, r, e) })
		if name == "ChainContract" {
			for _, alias := range []string{"chainContract", "chainCgContract", "chainXsContract"} {
				reg("/saas-contract/"+alias+"/importAlone", s.contractImportAll)
				if alias == "chainContract" {
					reg("/saas-contract/"+alias+"/importOnce", s.contractImportAll)
				} else {
					reg("/saas-contract/"+alias+"/importOnce", s.contractImportAliasOnce)
				}
				reg("/saas-contract/"+alias+"/importExcelCg", func(w http.ResponseWriter, r *http.Request) { s.contractImport(w, r, "HTCG.") })
				reg("/saas-contract/"+alias+"/importExcelXs", func(w http.ResponseWriter, r *http.Request) { s.contractImport(w, r, "HTXS.") })
				reg("/saas-contract/"+alias+"/exportExcel", func(w http.ResponseWriter, r *http.Request) { s.tplExportExcel(w, r, e) })
			}
			reg("/saas-contract/chainContract/importStep", s.contractImportStep)
			continue // 合同导入涉及多表，不能使用单表模板。
		}
		reg(base+"/importExcel", func(w http.ResponseWriter, r *http.Request) {
			// 原服务导入前均查询 max(system_code)，即使表中没有该列也执行原查询。
			left := currentUser(r).EnterpriseCode + strings.ReplaceAll(nowSpace()[:10], "-", "")
			if !sysProbe(s.pick(r), e.table, left) {
				platform.Unknown(w)
				return
			}
			s.importTemplate(w, r, e.table, e.excel)
		})
		if name == "ChainAuthUser" {
			reg(base+"/list", func(w http.ResponseWriter, r *http.Request) { s.tplList(w, r, e) })
			continue
		}
		if name == "ChainAuthMenu" || strings.HasPrefix(name, "ChainContract") {
			continue
		}
		reg(base+"/list", func(w http.ResponseWriter, r *http.Request) { s.tplList(w, r, e) })
		reg(base+"/listByPage", func(w http.ResponseWriter, r *http.Request) {
			family := byte('A')
			if name == "ChainAuthAuthenLog" {
				family = 'B'
			}
			s.tplListByPage(w, r, e, family)
		})
		reg(base+"/getById", func(w http.ResponseWriter, r *http.Request) { s.tplGetById(w, r, e) })
		reg(base+"/add", func(w http.ResponseWriter, r *http.Request) { s.tplAdd(w, r, e) })
		reg(base+"/updateById", func(w http.ResponseWriter, r *http.Request) { s.tplUpdateById(w, r, e) })
		reg(base+"/deleteById", func(w http.ResponseWriter, r *http.Request) { s.tplDeleteById(w, r, e) })
		reg(base+"/delBatch", func(w http.ResponseWriter, r *http.Request) { s.tplDelBatch(w, r, e, true) })
		if name == "ChainAuthAuthenLog" {
			reg(base+"/updBatch", func(w http.ResponseWriter, r *http.Request) { s.tplUpdBatch(w, r, e) })
		}
	}
}
