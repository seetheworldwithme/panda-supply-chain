// role.go 角色与权限分配(ChainAuthRoleController + 权限 service/mapper)。
// 每个 handler 注释标注 Java 源行为与复刻点。
package authadmin

import (
	"database/sql"
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// selEtcRoles 企业角色列表 = MP list(QueryWrapper.allEq(maps,false)):键原样作列名,
// 驼峰键(如 companyId)在 Java 直接触发 SQL 异常→999;null 值键跳过。响应实体驼峰
func (s *Service) selEtcRoles(w http.ResponseWriter, r *http.Request) {
	cond, ok := allEqCond(bodyMap(r), roleCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.DB.Chain, "chain_auth_role", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeRole(rows))
}

// selPtRoles 平台角色 = mapper.selPtRoles:role_type=1 按创建时间倒序
func (s *Service) selPtRoles(w http.ResponseWriter, r *http.Request) {
	rows, err := queryCam(s.DB.Chain,
		"SELECT * FROM chain_auth_role WHERE role_type=1 ORDER BY create_time DESC")
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeRole(rows))
}

// selPtUserRoles = mapper.selPtUserRoles(三表联查,Map 返回保留列名 username/name)
func (s *Service) selPtUserRoles(w http.ResponseWriter, r *http.Request) {
	rows, err := queryRows(s.DB.Chain,
		"SELECT username, r.name FROM chain_auth_user u, chain_auth_use_role ur, chain_auth_role r WHERE u.id=ur.user_id AND ur.role_id=r.id")
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, rows)
}

// selSigMenus 单个菜单+其操作权限 id 集(ChainAuthRoleController.selSigMenus):
// {menu: MP getById, auths: [{authy_id}]}(selByMenu 的 Map 列名原样)
func (s *Service) selSigMenus(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	menu, err := platform.GetByID(s.DB.Chain, "chain_auth_menu", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	auths, err := queryRows(s.DB.Chain,
		"SELECT authy_id FROM chain_auth_menu_authy WHERE menu_id = ?", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, map[string]any{
		"menu":  shapeMenuOne(menu),
		"auths": auths,
	})
}

// selAuthy 操作权限列表 = mapper.selAuthy:select id, zhname, enname,实体全键输出
func (s *Service) selAuthy(w http.ResponseWriter, r *http.Request) {
	rows, err := queryRows(s.DB.Chain, "SELECT id, zhname, enname FROM chain_auth_authority")
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeAuthority(rows))
}

// selAuthyByRole 角色已有权限 = ChainAuthAuthorityServiceImpl.selAuthyByRole:
// {menuRoles: selectMenuRoles, authRoles: selAuthRoles},两查询列名已驼峰,原样返回
func (s *Service) selAuthyByRole(w http.ResponseWriter, r *http.Request) {
	roleID := bodyMap(r)["roleId"]
	menuRoles, err := queryRows(s.DB.Chain,
		"SELECT menu_id id, role_id roleId, m.is_son isSon FROM chain_auth_menu_role mr, chain_auth_menu m WHERE mr.menu_id=m.id AND (m.menu IS NOT NULL AND m.menu !='') AND role_id = ?", roleID)
	if err != nil {
		platform.Unknown(w)
		return
	}
	authRoles, err := queryRows(s.DB.Chain,
		"SELECT menu_id menuId, menu, authoritys FROM chain_auth_authority_role r WHERE menu !='' AND role_id = ?", roleID)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, map[string]any{"menuRoles": menuRoles, "authRoles": authRoles})
}

// selAuthyByMenu 菜单已有操作权限 = mapper.selAuthyByMenu(Map 列名 zhname/enname 原样)
func (s *Service) selAuthyByMenu(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	rows, err := queryRows(s.DB.Chain,
		"SELECT zhname, enname FROM chain_auth_authority a, chain_auth_menu_authy ma WHERE a.id=ma.authy_id AND ma.menu_id = ?", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, rows)
}

// addEtcRole 新增企业角色:roleType 缺省 2,MP save(自动填充 create/update_time)。
// 无返回 data(ResultTool.success() 无参)
func (s *Service) addEtcRole(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	vals := pick(m, "name", "role_type", "company_id", "remark1", "remark2", "remark3", "remark4", "remark5")
	if cu := auth.From(r.Context()); !isPlatformAdmin(cu) { // 非平台管理员只能建本企业的企业角色,不能带出平台角色
		if !s.isEnterpriseAdmin(r, cu) {
			deny(w)
			return
		}
		vals["role_type"], vals["company_id"] = 2, cu.CompanyID
	}
	if _, ok := vals["role_type"]; !ok {
		vals["role_type"] = 2
	}
	if _, err := platform.Insert(s.DB.Chain, "chain_auth_role", vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// addPtRole 新增平台角色:roleType 强制 1
func (s *Service) addPtRole(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.requirePlatform(w, r) {
		return
	}
	vals := pick(m, "name", "company_id", "remark1", "remark2", "remark3", "remark4", "remark5")
	vals["role_type"] = 1
	if _, err := platform.Insert(s.DB.Chain, "chain_auth_role", vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// saveEtcAuths 保存企业权限分配 = ChainAuthAuthorityRoleServiceImpl.saveAuths:
// 先删后插 authority_role 与 menu_role 两张表(Java 无事务,Go 收进 WithTx)
func (s *Service) saveEtcAuths(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	roleID := m["roleId"]
	if !s.canManageRole(r, roleID) { // 原先任何登录用户可改任何角色的权限
		deny(w)
		return
	}
	authRoles, _ := m["authRoles"].([]any)
	menuRoles, _ := m["menuRoles"].([]any)
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		if err := platform.TxDeleteByMap(tx, "chain_auth_authority_role", map[string]any{"role_id": roleID}); err != nil {
			return err
		}
		for _, row := range authRoles {
			if rm, ok := row.(map[string]any); ok {
				if _, err := platform.TxInsert(tx, "chain_auth_authority_role",
					pick(rm, "id", "role_id", "menu_id", "menu", "authoritys")); err != nil {
					return err
				}
			}
		}
		if err := platform.TxDeleteByMap(tx, "chain_auth_menu_role", map[string]any{"role_id": roleID}); err != nil {
			return err
		}
		for _, row := range menuRoles {
			if rm, ok := row.(map[string]any); ok {
				if _, err := platform.TxInsert(tx, "chain_auth_menu_role",
					pick(rm, "id", "menu_id", "role_id")); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// savePtAuths 保存平台权限分配 = ChainAuthAuthorityRoleServiceImpl.savePtAuths:
// 重插 menu_user;auth_user 逐行先删后插(Java 循环里每行删一次,幂等)
func (s *Service) savePtAuths(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.requirePlatform(w, r) { // 平台用户的菜单 / 权限分配
		return
	}
	userID := m["userId"]
	menuUsers, _ := m["menuUsers"].([]any)
	authUsers, _ := m["authUsers"].([]any)
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		if err := platform.TxDeleteByMap(tx, "chain_auth_menu_user", map[string]any{"user_id": userID}); err != nil {
			return err
		}
		for _, row := range menuUsers {
			if rm, ok := row.(map[string]any); ok {
				if _, err := platform.TxInsert(tx, "chain_auth_menu_user",
					pick(rm, "id", "user_id", "menu_id")); err != nil {
					return err
				}
			}
		}
		for _, row := range authUsers {
			rm, ok := row.(map[string]any)
			if !ok {
				continue
			}
			if err := platform.TxDeleteByMap(tx, "chain_auth_authority_user",
				map[string]any{"user_id": rm["userId"]}); err != nil {
				return err
			}
			if _, err := platform.TxInsert(tx, "chain_auth_authority_user",
				pick(rm, "id", "user_id", "menu", "authoritys")); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// delRoleById 删除单个角色 = ChainAuthRoleServiceImpl.removeRole:
// 删角色行 + 删 use_role where role_id(Java 两步无事务,Go 收 WithTx)
func (s *Service) delRoleById(w http.ResponseWriter, r *http.Request) {
	id := bodyMap(r)["id"]
	if !s.canManageRole(r, id) {
		deny(w)
		return
	}
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		if err := platform.TxDeleteByMap(tx, "chain_auth_role", map[string]any{"id": id}); err != nil {
			return err
		}
		return platform.TxDeleteByMap(tx, "chain_auth_use_role", map[string]any{"role_id": id})
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// roleUpdateById 修改角色 = MP updateById(非空字段更新),返回请求体实体
func (s *Service) roleUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	vals := pick(m, "name", "role_type", "company_id", "remark1", "remark2", "remark3", "remark4", "remark5")
	if !s.canManageRole(r, m["id"]) {
		deny(w)
		return
	}
	if !isPlatformAdmin(auth.From(r.Context())) { // 不能把企业角色改成平台角色或换归属
		delete(vals, "role_type")
		delete(vals, "company_id")
	}
	if _, err := platform.UpdateByID(s.DB.Chain, "chain_auth_role", m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, m)
}

// updRoleForPt 平台改角色 = mapper.updRoleForPt 命名列 UPDATE,返回请求体实体
func (s *Service) updRoleForPt(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.requirePlatform(w, r) {
		return
	}
	if _, err := s.DB.Chain.Exec(
		"UPDATE chain_auth_role SET name=?, role_type=? WHERE id=?",
		m["name"], m["roleType"], m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, m)
}

// roleDelBatch 批量删除 = MP removeByIds;返回文案 "删除成功"
func (s *Service) roleDelBatch(w http.ResponseWriter, r *http.Request) {
	ids, _ := bodyMap(r)["ids"].([]any)
	for _, id := range ids {
		if !s.canManageRole(r, id) {
			deny(w)
			return
		}
	}
	if err := platform.DeleteByIDs(s.DB.Chain, "chain_auth_role", ids); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// getByRoleId = MP getById;无记录时 data null
func (s *Service) getByRoleId(w http.ResponseWriter, r *http.Request) {
	menu, err := platform.GetByID(s.DB.Chain, "chain_auth_role", bodyMap(r)["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeRoleOne(menu))
}

// getRoleForPt 平台查角色 = mapper.getRoleForPt(SQL 与 getById 相同)
func (s *Service) getRoleForPt(w http.ResponseWriter, r *http.Request) {
	menu, err := platform.GetByID(s.DB.Chain, "chain_auth_role", bodyMap(r)["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeRoleOne(menu))
}

// getByMenuId = MP getById(chain_auth_menu)
func (s *Service) getByMenuId(w http.ResponseWriter, r *http.Request) {
	menu, err := platform.GetByID(s.DB.Chain, "chain_auth_menu", bodyMap(r)["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeMenuOne(menu))
}

