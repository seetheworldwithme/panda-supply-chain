// entity.go 实体 JSON 形状层:对照 MP 实体 Jackson 序列化。
// 差异三源——DATETIME 格式(ScanMaps 已修)、tinyint(1)→布尔、实体载体字段(searchStr/page/ids
// 等非列字段)恒输出 null——此处统一处理。新增实体时仿写 shape 函数。
package authadmin

import "github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"

// menuBoolCols 已并入 shapeMenu 的 camel 键字面量(历史备注)

// 载体键(null)
var roleExtra = []string{"searchStr", "page", "ids"}
var menuExtra = []string{"searchStr", "page", "ids", "menuAuths", "auths", "pname"}
var authyAll = []string{"id", "urls", "zhname", "enname", "createTime", "updateTime",
	"remark1", "remark2", "remark3", "remark4", "remark5", "searchStr", "page", "ids"}
var menuAuthyAll = []string{"id", "menuId", "authyId", "createTime", "updateTime",
	"remark1", "remark2", "remark3", "remark4"}

// shapeRole 角色行输出(select * 行 + 载体 null)
func shapeRole(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		fillNulls(r, roleExtra)
	}
	return out
}

// shapeMenu 菜单行输出(布尔化 + 载体 null + searchType int 0)。
// 入参可能 snake(ListWhere/GetByID)或已 camel(paged),统一在 Camelize 后按 camel 键处理
func shapeMenu(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		for _, c := range []string{"isAuthy", "hasFunChild", "openAuthy", "disabled", "isSon"} {
			if v, ok := r[c]; ok && v != nil {
				r[c] = num(v) != 0
			}
		}
		fillNulls(r, menuExtra)
		if _, ok := r["searchType"]; !ok {
			r["searchType"] = 0
		}
	}
	return out
}

// shapeMenuOne 单行菜单(selSigMenus/getByMenuId)
func shapeMenuOne(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	return shapeMenu([]map[string]any{m})[0]
}

// shapeRoleOne 单行角色
func shapeRoleOne(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	return shapeRole([]map[string]any{m})[0]
}

// shapeAuthority 操作权限行(SQL 只取部分列,实体全键输出)
func shapeAuthority(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		fillNulls(r, authyAll)
	}
	return out
}

// shapeMenuAuthy 菜单权限关联行(实体全键)
func shapeMenuAuthy(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		fillNulls(r, menuAuthyAll)
	}
	return out
}

// fillNulls 补缺失键为 null(对照实体 null 字段的 Jackson 输出)
func fillNulls(m map[string]any, keys []string) {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			m[k] = nil
		}
	}
}
