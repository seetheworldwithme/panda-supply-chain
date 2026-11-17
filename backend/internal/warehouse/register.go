// register.go 模板十件套路由挂载。tplCtls(routes_gen.go)逐控制器配置;
// bespoke 端点在各业务文件里经 extraRegs 注册,注册前已 grep 全 internal 防冲突
// (/busRedEnvelopeRange/list、/busRedEnvelopeRange/updBatch、/busProLog/selMemberData
//  二期 authadmin 已实现,路由表已标 skip)。
package warehouse

import (
	"encoding/json"
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// tplCtl 单控制器的模板端点配置(见 routes_gen.go)
type tplCtl struct {
	Path         string
	Ent          *ent
	Family       byte   // 'A' 无 page → 999;'B' 无 page → listNoPage;'C' B + startTime/endTime 过滤
	DelBatchData bool   // delBatch 响应 data "删除成功"(true)或无 data(false)
	ExportGET    bool   // exportExcel 用 GET(仅 BusProduct,十件套 skip 另走 bespoke)
	Wrap2        bool   // listByPage 把 JsonResult 再包一层(仅 BusGlobalConfig,8765 实测双层 data)
	HasUpdBatch  bool   // 控制器有 /updBatch 附加端点
	Skip         []string // 不按模板注册的端点路径(bespoke/重复)
}

// extraRegs bespoke 端点注册钩子(各业务文件在 init 或 Register 前填充)
var extraRegs []func(mux *http.ServeMux, s *Service)

func addExtra(f func(mux *http.ServeMux, s *Service)) { extraRegs = append(extraRegs, f) }

// Register 挂载全部 warehouse 路由
func Register(mux *http.ServeMux, s *Service) {
	for i := range tplCtls {
		c := &tplCtls[i]
		skip := map[string]bool{}
		for _, p := range c.Skip {
			skip[p] = true
		}
		e := c.Ent
		if !skip["/list"] {
			mux.HandleFunc("POST /saas-warehouse"+c.Path+"/list", authGuard(func(w http.ResponseWriter, r *http.Request) {
				s.tplList(w, r, e)
			}))
		}
		if !skip["/listByPage"] {
			mux.HandleFunc("POST /saas-warehouse"+c.Path+"/listByPage", authGuard(func(w http.ResponseWriter, r *http.Request) {
				s.tplListByPageWrap(w, r, e, c.Family, c.Wrap2)
			}))
		}
		if !skip["/getById"] {
			mux.HandleFunc("POST /saas-warehouse"+c.Path+"/getById", authGuard(func(w http.ResponseWriter, r *http.Request) {
				s.tplGetById(w, r, e)
			}))
		}
		if !skip["/add"] {
			mux.HandleFunc("POST /saas-warehouse"+c.Path+"/add", authGuard(func(w http.ResponseWriter, r *http.Request) {
				s.tplAdd(w, r, e)
			}))
		}
		if !skip["/deleteById"] {
			mux.HandleFunc("POST /saas-warehouse"+c.Path+"/deleteById", authGuard(func(w http.ResponseWriter, r *http.Request) {
				s.tplDeleteById(w, r, e)
			}))
		}
		if !skip["/updateById"] {
			mux.HandleFunc("POST /saas-warehouse"+c.Path+"/updateById", authGuard(func(w http.ResponseWriter, r *http.Request) {
				s.tplUpdateById(w, r, e)
			}))
		}
		if !skip["/importExcel"] {
			mux.HandleFunc("POST /saas-warehouse"+c.Path+"/importExcel", authGuard(func(w http.ResponseWriter, r *http.Request) {
				s.tplImportExcel(w, r, e)
			}))
		}
		if !skip["/exportExcel"] {
			method := "POST /saas-warehouse"
			if c.ExportGET {
				method = "GET /saas-warehouse"
			}
			mux.HandleFunc(method+c.Path+"/exportExcel", authGuard(func(w http.ResponseWriter, r *http.Request) {
				s.tplExportExcel(w, r, e)
			}))
		}
		if !skip["/delBatch"] {
			withData := c.DelBatchData
			mux.HandleFunc("POST /saas-warehouse"+c.Path+"/delBatch", authGuard(func(w http.ResponseWriter, r *http.Request) {
				s.tplDelBatch(w, r, e, withData)
			}))
		}
		if c.HasUpdBatch && !skip["/updBatch"] {
			mux.HandleFunc("POST /saas-warehouse"+c.Path+"/updBatch", authGuard(func(w http.ResponseWriter, r *http.Request) {
				s.tplUpdBatch(w, r, e)
			}))
		}
	}
	for _, f := range extraRegs {
		f(mux, s)
	}
}

// tplListByPageWrap family 分发 + Wrap2 双层包裹(仅 BusGlobalConfig)
func (s *Service) tplListByPageWrap(w http.ResponseWriter, r *http.Request, e *ent, family byte, wrap2 bool) {
	if !wrap2 {
		s.tplListByPage(w, r, e, family)
		return
	}
	rec := &capWriter{header: http.Header{}}
	s.tplListByPage(rec, r, e, family)
	var inner any
	_ = json.Unmarshal(rec.buf, &inner)
	platform.Success(w, inner)
}

// capWriter 捕获内层响应体(双层包裹怪癖专用)
type capWriter struct {
	header http.Header
	buf    []byte
	code   int
}

func (c *capWriter) Header() http.Header         { return c.header }
func (c *capWriter) Write(b []byte) (int, error) { c.buf = append(c.buf, b...); return len(b), nil }
func (c *capWriter) WriteHeader(code int)        { c.code = code }
