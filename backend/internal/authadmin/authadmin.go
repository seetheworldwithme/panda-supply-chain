// Package authadmin 二期组织基础:角色/菜单/权限分配管理接口(前端 setup 目录页面)。
// Java 源:saas-auth ChainAuthRoleController/ChainAuthMenuController(网关前缀 /saas-auth)。
// selMenus/selUserMenus 已由 menu 包实现,不在此重复注册。
// 已知 Java 缺陷的行为保持(勿"修复"):菜单 getById 的 auths 设在请求体对象上被丢弃(恒 null);
// delMenu 只删菜单行不级联 menu_authy;selJxzb 类 bug 见 dashboard 包注释。
package authadmin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type Service struct {
	DB   *db.DB
	Auth *auth.Service
}

func NewService(database *db.DB) *Service { return &Service{DB: database} }

// Register 挂载角色/菜单管理路由(路径与 Java 网关后服务路径一致,全部 POST)
func Register(mux *http.ServeMux, s *Service) {
	r := func(path string, h http.HandlerFunc) {
		mux.HandleFunc("POST /saas-auth"+path, h)
	}
	s.registerUserCompat(r)
	// ChainAuthRoleController(角色与权限分配)
	r("/chainAuthRole/selEtcRoles", s.selEtcRoles)
	r("/chainAuthRole/selPtRoles", s.selPtRoles)
	r("/chainAuthRole/selPtUserRoles", s.selPtUserRoles)
	r("/chainAuthRole/selSigMenus", s.selSigMenus)
	r("/chainAuthRole/selAuthy", s.selAuthy)
	r("/chainAuthRole/selAuthyByRole", s.selAuthyByRole)
	r("/chainAuthRole/selAuthyByMenu", s.selAuthyByMenu)
	r("/chainAuthRole/addEtcRole", s.addEtcRole)
	r("/chainAuthRole/addPtRole", s.addPtRole)
	r("/chainAuthRole/saveEtcAuths", s.saveEtcAuths)
	r("/chainAuthRole/savePtAuths", s.savePtAuths)
	r("/chainAuthRole/delRoleById", s.delRoleById)
	r("/chainAuthRole/updateById", s.roleUpdateById)
	r("/chainAuthRole/updRoleForPt", s.updRoleForPt)
	r("/chainAuthRole/delBatch", s.roleDelBatch)
	r("/chainAuthRole/getByRoleId", s.getByRoleId)
	r("/chainAuthRole/getRoleForPt", s.getRoleForPt)
	r("/chainAuthRole/getByMenuId", s.getByMenuId)
	r("/chainAuthRole/addMenu", s.addMenuLoose)
	r("/chainAuthRole/updMenu", s.updMenuLoose)
	r("/chainAuthRole/delMenu", s.delMenuNoCascade)
	// ChainAuthMenuController(菜单 CRUD)
	r("/chainAuthMenu/list", s.menuList)
	r("/chainAuthMenu/listByPage", s.menuListByPage)
	r("/chainAuthMenu/getById", s.menuGetById)
	r("/chainAuthMenu/add", s.menuAdd)
	r("/chainAuthMenu/deleteById", s.menuDeleteById)
	r("/chainAuthMenu/updateById", s.menuUpdateById)
	r("/chainAuthMenu/delBatch", s.menuDelBatch)
	r("/chainAuthMenu/selAuthys", s.menuSelAuthys)
	// ChainAuthUserController(用户)
	r("/chainAuthUser/listByPage", s.userListByPage)
	r("/chainAuthUser/supList", s.userSupList)
	r("/chainAuthUser/add", s.userAdd)
	// 服务器运行包的方法正常编译；本机重编译包的缺陷不作为契约。
	r("/chainAuthUser/updateById", s.userUpdateById)
	r("/chainAuthUser/deleteById", s.userDeleteById)
	r("/chainAuthUser/delBatch", s.userDelBatch)
	r("/chainAuthUser/cancelBuser", s.userCancelBuser)
	r("/chainAuthUser/logOffUser", s.userLogOffUser)
	r("/chainAuthUser/createCode", s.userCreateCode)
	// ChainAuthUserListController(用户列表)
	r("/chainAuthUserList/listByPage", s.userListByPage)
	r("/chainAuthUserList/updateById", s.userListUpdByList)
	r("/chainAuthUserList/assignRole", s.assignRole)
	r("/chainAuthUserList/selUserRoles", s.selUserRoles)
	// 二期收尾:企业档案/全局配置/e签宝配置/上传(saas-common 服务)
	c := func(path string, h http.HandlerFunc) {
		mux.HandleFunc("POST /saas-common"+path, h)
	}
	c("/chainBasicEnterprise/getById", s.entGetById)
	c("/chainBasicEnterprise/updateById", s.entUpdateById)
	c("/chainGlobalConfiguration/list", s.glconfigList)
	c("/chainGlobalConfiguration/updateById", s.glconfigUpdateById)
	c("/chainEqbConfig/listByPage", s.eqbListByPage)
	c("/chainEqbConfig/getById", s.eqbGetById)
	c("/chainEqbConfig/add", s.eqbAdd)
	c("/chainEqbConfig/updateById", s.eqbUpdateById)
	c("/chainEqbConfig/delBatch", s.eqbDelBatch)
	c("/upload/batch", s.uploadBatch)
	// 合同供应商(saas-contract 服务)
	ct := func(path string, h http.HandlerFunc) {
		mux.HandleFunc("POST /saas-contract"+path, h)
	}
	ct("/chainContractSupplier/listByPage", s.supplierListByPage)
	ct("/chainContractSupplier/getById", s.supplierGetById)
	ct("/chainContractSupplier/add", s.supplierAdd)
	ct("/chainContractSupplier/updateById", s.supplierUpdateById)
	ct("/chainContractSupplier/delBatch", s.supplierDelBatch)
	// 红包范围与访问日志(saas-warehouse 服务)
	wr := func(path string, h http.HandlerFunc) {
		mux.HandleFunc("POST /saas-warehouse"+path, h)
	}
	wr("/busRedEnvelopeRange/list", s.brelrgeList)
	wr("/busRedEnvelopeRange/updBatch", s.brelrgeUpdBatch)
	wr("/busProLog/selMemberData", s.proLogSelMemberData)
}

// ── 共享小工具 ──

// bodyMap 解析 JSON body 为宽松 map(缺字段/类型漂移都不报错)
func bodyMap(r *http.Request) map[string]any {
	var m map[string]any
	_ = json.NewDecoder(r.Body).Decode(&m)
	if m == nil {
		return map[string]any{}
	}
	return m
}

// str 取字符串字段(JSON 数字转字符串,前端 companyId/timeType 常为数字)
func str(m map[string]any, k string) string {
	v, ok := m[k]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func parseF(s string) float64 {
	n, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return n
}

// pick 从 body 按实体字段白名单挑列(键统一转下划线),剔除 searchStr/page/ids 等非列字段
func pick(m map[string]any, cols ...string) map[string]any {
	u := platform.Underscore(m)
	out := make(map[string]any, len(cols))
	for _, c := range cols {
		if v, ok := u[c]; ok {
			out[c] = v
		}
	}
	return out
}

// queryRows 查询并原样返回(MyBatis Map 返回保留列名;空结果集返回空数组,对照 Java 空列表序列化 [])
func queryRows(d *sql.DB, q string, args ...any) ([]map[string]any, error) {
	rows, err := d.Query(q, args...)
	if err != nil {
		return nil, err
	}
	ms, err := platform.ScanMaps(rows)
	if ms == nil {
		ms = []map[string]any{}
	}
	return ms, err
}

// queryCam 查询并驼峰化(MP 实体序列化键名)
func queryCam(d *sql.DB, q string, args ...any) ([]map[string]any, error) {
	ms, err := queryRows(d, q, args...)
	if err != nil {
		return nil, err
	}
	return platform.CamelizeRows(ms), nil
}

// pageOf 从 body 取 {current,size}(MP Page 默认 1/10)
func pageOf(m map[string]any) (current, size int64) {
	current, size = 1, 10
	if p, ok := m["page"].(map[string]any); ok {
		if v := num(p["current"]); v > 0 {
			current = int64(v)
		}
		if v := num(p["size"]); v > 0 {
			size = int64(v)
		}
	}
	return
}

// num 把 driver/JSON 数值统一转 float64(DECIMAL 的 []byte、int64、float64)
func num(v any) float64 {
	switch t := v.(type) {
	case nil:
		return 0
	case []byte:
		return parseF(string(t))
	case string:
		return parseF(t)
	case int64:
		return float64(t)
	case float64:
		return t
	default:
		return 0
	}
}

// paged 按相同条件执行 count + 窗口查询,返回 MP IPage 形状。
// countSQL 由调用方给(从主查询去掉 ORDER BY 的版本),避免子查询包裹的方言坑
func paged(d *sql.DB, q, countQ string, args []any, current, size int64) (map[string]any, error) {
	var total int64
	if err := d.QueryRow(countQ, args...).Scan(&total); err != nil {
		return nil, err
	}
	off := (current - 1) * size
	rows, err := queryCam(d, q+" LIMIT ? OFFSET ?", append(append([]any{}, args...), size, off)...)
	if err != nil {
		return nil, err
	}
	return platform.IPage(rows, total, current, size), nil
}

// searchStrOf 复刻 MyBatis ${searchStr} 的 null 语义:字段缺失或 null → 字面 "null"
func searchStrOf(m map[string]any) string {
	if v, ok := m["searchStr"]; !ok || v == nil {
		return "null"
	}
	return strings.TrimSpace(str(m, "searchStr"))
}

// allEqCond 复刻 MP QueryWrapper.allEq(maps,false):键原样作列名(Java 直接拼键,大小写敏感的
// 驼峰键会触发 SQL 异常 → 全局异常处理器 999,此处用列白名单模拟:白名单外一律 999)。
// null 值键跳过(false 不生成 IS NULL)
func allEqCond(m map[string]any, cols []string) (map[string]any, bool) {
	allow := make(map[string]bool, len(cols))
	for _, c := range cols {
		allow[c] = true
	}
	cond := map[string]any{}
	for k, v := range m {
		if v == nil {
			continue
		}
		if !allow[k] {
			return nil, false
		}
		cond[k] = v
	}
	return cond, true
}

// roleCols/allMenuCols allEq 白名单(表实际列名)
var roleCols = []string{"id", "name", "role_type", "company_id",
	"create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"}
var allMenuCols = append([]string{"id", "pid"}, menuCols...)
