// Package menu 动态菜单。对照 Java ChainAuthMenuServiceImpl.selUserMenus/selMenus,
// 响应结构逐字段对齐(前端 store 依赖 routeName/routelist/menuAuthy 的键名)
package menu

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service { return &Service{DB: db} }

// MenuRow 菜单行。占位字段 add/edit/... 空串起步,有权限时置同名值(照 Java 原样,前端 if(m.add) 判断)
type MenuRow struct {
	ID           int64          `json:"id"`
	PID          int64          `json:"pid"`
	Name         string         `json:"name"`
	RouteName    sql.NullString `json:"routeName"`
	Menu         sql.NullString `json:"menu"`
	HasFunChild  sql.NullString `json:"hasFunChild"`
	Title        *string        `json:"title,omitempty"` // 仅 selByPid/selMenus 分支返回
	Disabled     *int           `json:"disabled,omitempty"`
	Add          string         `json:"add"`
	Edit         string         `json:"edit"`
	Check        string         `json:"check"`
	Delete       string         `json:"delete"`
	Import       string         `json:"import"`
	Export       string         `json:"export"`
	FileUpload   string         `json:"fileUpload"`
	FileDelete   string         `json:"fileDelete"`
}

const roleMenuCols = "m.id, m.pid, m.name, m.route_name routeName, m.menu, m.has_fun_child hasFunChild"

// SelUserMenus POST /saas-auth/chainAuthMenu/selUserMenus
// 响应 {menus, routelist},routelist 驱动前端动态路由 redirect(见 store/index.js asynSetplatformRoutes)
func (s *Service) SelUserMenus(w http.ResponseWriter, r *http.Request) {
	cu := auth.From(r.Context())
	if cu.UserID == 0 {
		platform.NotAuth(w)
		return
	}
	// 用户 account_type(决定走超管分支还是角色菜单分支)
	var accountType sql.NullString
	if err := s.DB.QueryRow("SELECT account_type FROM chain_auth_user WHERE id = ?", cu.UserID).Scan(&accountType); err != nil {
		platform.Unknown(w)
		return
	}

	var menus []map[string]interface{}
	routelist := []map[string]interface{}{
		{"name": "", "chinname": ""},
		{"name": "", "chinname": ""},
	}

	if accountType.String == "5" {
		// 超管分支:selSuMenus + 固定 routelist(照 Java)
		rows, err := s.DB.Query("SELECT " + roleMenuCols + ", '' as 'add', '' as 'edit', '' as 'check', '' as 'delete', '' as 'import', '' as 'export' FROM chain_auth_menu m WHERE m.route_name != 'menuPower' AND (m.menu_type = 2 OR m.menu_type = 3)")
		if err != nil {
			platform.Unknown(w)
			return
		}
		menus = collectRoleMenus(rows, false)
		routelist[0] = map[string]interface{}{"name": "Mindex", "chinname": "PtIndex"}
		routelist[1] = map[string]interface{}{"name": "", "chinname": "PtUserList"}
	} else {
		// 角色菜单 = 角色授权菜单 + 用户直配菜单(Java 端 addAll 不去重,照原样)
		rows, err := s.DB.Query(`SELECT DISTINCT ` + roleMenuCols + `, '' as 'add', '' as 'edit', '' as 'check', '' as 'delete', '' as 'import', '' as 'export', '' as 'fileUpload', '' as 'fileDelete'
			FROM chain_auth_use_role ur, chain_auth_menu_role mr, chain_auth_menu m
			WHERE ur.role_id = mr.role_id AND mr.menu_id = m.id AND ur.user_id = ?`, cu.UserID)
		if err != nil {
			platform.Unknown(w)
			return
		}
		menus = collectRoleMenus(rows, true)
		rows2, err := s.DB.Query(`SELECT DISTINCT ` + roleMenuCols + `, '' as 'add', '' as 'edit', '' as 'check', '' as 'delete', '' as 'import', '' as 'export', '' as 'fileUpload', '' as 'fileDelete'
			FROM chain_auth_menu_user mu, chain_auth_menu m
			WHERE mu.menu_id = m.id AND mu.user_id = ?`, cu.UserID)
		if err == nil {
			menus = append(menus, collectRoleMenus(rows2, true)...)
		}
		s.applyAuthorityMarks(cu.UserID, menus)

		// chinname/name 采集与功能子菜单(照 Java 遍历逻辑)
		menuMap := map[string]bool{}
		mchinDone, cchinDone := false, false
		var childMenus []map[string]interface{}
		for _, m := range menus {
			if rn, ok := m["routeName"].(string); ok {
				menuMap[rn] = true
			}
			if !mchinDone && m["pid"] == int64(1) {
				if rn, ok := m["routeName"].(string); ok {
					routelist[0]["chinname"] = rn
					mchinDone = true
				}
			} else if !cchinDone && m["pid"] == int64(2) {
				if rn, ok := m["routeName"].(string); ok {
					routelist[1]["chinname"] = rn
					cchinDone = true
				}
			}
			if hfc, ok := m["hasFunChild"].(string); ok && hfc == "true" {
				id, _ := m["id"].(int64)
				childMenus = append(childMenus, s.selByPid(id)...)
			}
		}
		if menuMap["Mindex"] {
			routelist[0]["name"] = "Mindex"
			routelist[1]["name"] = ""
		} else if menuMap["Setup"] {
			routelist[0]["name"] = ""
			routelist[1]["name"] = "Setup"
		}
		menus = append(menus, childMenus...)
	}
	if menus == nil {
		menus = []map[string]interface{}{}
	}
	platform.Success(w, map[string]interface{}{"menus": menus, "routelist": routelist})
}

// selByPid 对照 selByPid: is_authy=0 的功能按钮子菜单
func (s *Service) selByPid(pid int64) []map[string]interface{} {
	rows, err := s.DB.Query("SELECT id, pid, route_name routeName, name title, menu FROM chain_auth_menu WHERE pid = ? AND is_authy = 0", pid)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var id, p int64
		var rn, menu sql.NullString
		var title sql.NullString
		if err := rows.Scan(&id, &p, &rn, &title, &menu); err != nil {
			continue
		}
		out = append(out, map[string]interface{}{"id": id, "pid": p, "routeName": rn.String, "title": title.String, "menu": menu.String})
	}
	return out
}

// applyAuthorityMarks 对照 selMenuAuths + 遍历标记:GROUP_CONCAT authoritys 命中即置同名值
func (s *Service) applyAuthorityMarks(userID int64, menus []map[string]interface{}) {
	rows, err := s.DB.Query(`SELECT DISTINCT ar.menu, GROUP_CONCAT(ar.authoritys) authoritys
		FROM chain_auth_use_role ur, chain_auth_authority_role ar
		WHERE ur.role_id = ar.role_id AND ur.user_id = ? GROUP BY menu`, userID)
	if err != nil {
		return
	}
	defer rows.Close()
	lsmap := map[string]string{}
	for rows.Next() {
		var menu, authoritys sql.NullString
		if _ = rows.Scan(&menu, &authoritys); menu.Valid && authoritys.Valid {
			lsmap[menu.String] = authoritys.String
		}
	}
	marks := []string{"add", "edit", "delete", "check", "import", "export", "fileUpload", "fileDelete"}
	for _, m := range menus {
		menuKey, _ := m["menu"].(string)
		authoritys, ok := lsmap[menuKey]
		if !ok {
			continue
		}
		for _, mk := range marks {
			if contains(authoritys, mk) {
				m[mk] = mk
			}
		}
	}
}

// SelMenus POST /saas-auth/chainAuthRole/selMenus (JSON: {searchType})
// 前端角色权限页 searchType=3 → selOwnerMenus/selOwnerMas;响应 {menus, menuAuthy}
func (s *Service) SelMenus(w http.ResponseWriter, r *http.Request) {
	cu := auth.From(r.Context())
	if cu.UserID == 0 {
		platform.NotAuth(w)
		return
	}
	rows, err := s.DB.Query(`SELECT m.id, m.pid, m.route_name routeName, m.name title, m.menu, m.disabled
		FROM chain_auth_menu m, chain_auth_menu_role mr, chain_auth_use_role ur
		WHERE m.id = mr.menu_id AND mr.role_id = ur.role_id AND ur.user_id = ?`, cu.UserID)
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer rows.Close()
	menus := []map[string]interface{}{}
	for rows.Next() {
		var id, pid int64
		var rn, menu sql.NullString
		var title sql.NullString
		var disabled sql.NullInt64
		if err := rows.Scan(&id, &pid, &rn, &title, &menu, &disabled); err != nil {
			continue
		}
		row := map[string]interface{}{"id": id, "pid": pid, "routeName": rn.String, "title": title.String, "menu": menu.String}
		if disabled.Valid {
			row["disabled"] = int(disabled.Int64)
		} else {
			row["disabled"] = nil
		}
		menus = append(menus, row)
	}

	masRows, err := s.DB.Query(`SELECT m.id, m.menu, a.zhname, a.enname
		FROM chain_auth_menu m, chain_auth_menu_authy ma, chain_auth_authority a, chain_auth_menu_role mr, chain_auth_use_role ur
		WHERE m.id = ma.menu_id AND ma.authy_id = a.id AND ma.menu_id = mr.menu_id AND mr.role_id = ur.role_id AND ur.user_id = ?`, cu.UserID)
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer masRows.Close()
	menuAuthy := map[string]map[string]string{}
	for masRows.Next() {
		var id int64
		var menuID, zhname, enname sql.NullString
		if err := masRows.Scan(&id, &menuID, &zhname, &enname); err != nil {
			continue
		}
		key := menuID.String
		if key == "" {
			key = strconv.FormatInt(id, 10)
		}
		if menuAuthy[key] == nil {
			menuAuthy[key] = map[string]string{}
		}
		if enname.Valid {
			menuAuthy[key][enname.String] = enname.String
		}
	}
	platform.Success(w, map[string]interface{}{"menus": menus, "menuAuthy": menuAuthy})
}

// collectRoleMenus 扫描角色菜单行。withFileOps=true 时含 fileUpload/fileDelete 两列
func collectRoleMenus(rows *sql.Rows, withFileOps bool) []map[string]interface{} {
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var id, pid int64
		var name string
		var rn, menu, hfc sql.NullString
		var add, edit, check, del, imp, exp string
		if withFileOps {
			var fu, fd string
			if err := rows.Scan(&id, &pid, &name, &rn, &menu, &hfc, &add, &edit, &check, &del, &imp, &exp, &fu, &fd); err != nil {
				continue
			}
			out = append(out, map[string]interface{}{
				"id": id, "pid": pid, "name": name, "routeName": rn.String, "menu": menu.String, "hasFunChild": hfc.String,
				"add": add, "edit": edit, "check": check, "delete": del, "import": imp, "export": exp,
				"fileUpload": fu, "fileDelete": fd,
			})
		} else {
			if err := rows.Scan(&id, &pid, &name, &rn, &menu, &hfc, &add, &edit, &check, &del, &imp, &exp); err != nil {
				continue
			}
			out = append(out, map[string]interface{}{
				"id": id, "pid": pid, "name": name, "routeName": rn.String, "menu": menu.String, "hasFunChild": hfc.String,
				"add": add, "edit": edit, "check": check, "delete": del, "import": imp, "export": exp,
			})
		}
	}
	return out
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
