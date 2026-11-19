// menu.go 菜单 CRUD(ChainAuthMenuController + Role controller 内的菜单入口)。
// 已知 Java 缺陷行为保持:menuGetById 的 auths 恒 null(见包注释);delMenuNoCascade 不级联。
package authadmin

import (
	"database/sql"
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// menuRoleCols 菜单实体允许写入的列(MP 实体非静态字段减去 searchStr/page/ids 等载体)
var menuCols = []string{"name", "route_name", "pid", "menu", "is_authy", "has_fun_child",
	"open_authy", "disabled", "menu_type", "is_son", "remark1", "remark2", "remark3", "remark4", "remark5"}

// menuList = MP list(QueryWrapper.allEq(maps,false)):键原样作列名,非法键→999, null 值键跳过
func (s *Service) menuList(w http.ResponseWriter, r *http.Request) {
	cond, ok := allEqCond(bodyMap(r), allMenuCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.DB.Chain, "chain_auth_menu", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeMenu(rows))
}

// menuListByPage = mapper.listByPage(maps 参数在 SQL 中未使用,仅 searchStr 生效):
// select *, (select name ... pid) pname where CONCAT(name, route_name) like '%ss%' order by create_time desc
func (s *Service) menuListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["page"].(map[string]any); !ok {
		platform.Unknown(w) // Java mapper requires a Page object.
		return
	}
	current, size := pageOf(m)
	ss := searchStrOf(m)
	q := "SELECT *, (SELECT name FROM chain_auth_menu WHERE id=m.pid) pname FROM chain_auth_menu m WHERE CONCAT(name, route_name) LIKE ? ORDER BY create_time DESC"
	countQ := "SELECT COUNT(*) FROM chain_auth_menu m WHERE CONCAT(name, route_name) LIKE ?"
	page, err := paged(s.DB.Chain, q, countQ, []any{"%" + ss + "%"}, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeMenu(page["records"].([]map[string]any))
	platform.Success(w, page)
}

// menuGetById = controller.getById:返回实体序列化形状 = 菜单字段 + menuAuths 键。
// Java 把 auths 设在了请求体对象上、返回的是 DB 对象 → auths 恒 null,行为保持不设置
func (s *Service) menuGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	menu, err := platform.GetByID(s.DB.Chain, "chain_auth_menu", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	menuAuths, err := queryCam(s.DB.Chain,
		"SELECT * FROM chain_auth_menu_authy WHERE menu_id = ?", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if menu == nil {
		platform.Unknown(w) // Java service dereferences the missing menu.
		return
	}
	out := shapeMenuOne(menu)
	out["menuAuths"] = shapeMenuAuthy(menuAuths)
	platform.Success(w, out)
}

// menuAdd 新增 = ChainAuthMenuServiceImpl.saveMenu:save(自增/显式 id)后按 ids[] 插 menu_authy
func (s *Service) menuAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	vals := pick(m, append([]string{"id"}, menuCols...)...)
	var newID int64
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		var err error
		newID, err = platform.TxInsert(tx, "chain_auth_menu", vals)
		if err != nil {
			return err
		}
		return insertMenuAuthys(tx, float64(newID), m["ids"])
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, m)
}

// menuDeleteById = MP removeById;文案 "删除成功，id为：N"
func (s *Service) menuDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if err := platform.DeleteByID(s.DB.Chain, "chain_auth_menu", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

// menuUpdateById 修改 = ChainAuthMenuServiceImpl.updateMenu:
// updateById → 删 menu_authy where menu_id → 按 ids[] 重插
func (s *Service) menuUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	vals := pick(m, menuCols...)
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		if err := platform.TxUpdateByID(tx, "chain_auth_menu", m["id"], vals); err != nil {
			return err
		}
		if err := platform.TxDeleteByMap(tx, "chain_auth_menu_authy", map[string]any{"menu_id": m["id"]}); err != nil {
			return err
		}
		return insertMenuAuthys(tx, num(m["id"]), m["ids"])
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, m)
}

// menuDelBatch 批量删除 = ChainAuthMenuServiceImpl.removeMenus:删菜单行 + 逐 id 删 menu_authy
func (s *Service) menuDelBatch(w http.ResponseWriter, r *http.Request) {
	ids, _ := bodyMap(r)["ids"].([]any)
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		if err := platform.TxDeleteByIDs(tx, "chain_auth_menu", ids); err != nil {
			return err
		}
		for _, id := range ids {
			if err := platform.TxDeleteByMap(tx, "chain_auth_menu_authy", map[string]any{"menu_id": id}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// menuSelAuthys 操作权限列表 = MP list() 全表
func (s *Service) menuSelAuthys(w http.ResponseWriter, r *http.Request) {
	rows, err := queryCam(s.DB.Chain, "SELECT * FROM chain_auth_authority")
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeAuthority(rows))
}

// ── Role controller 内的菜单入口(网关路径 /chainAuthRole/*) ──

// addMenuLoose = ChainAuthMenuServiceImpl.addMenu:menuAuths 非空则每行 menuId 覆盖为 body.menu.id
// 后批量插;否则 openAuthy=false;最后 save 菜单
func (s *Service) addMenuLoose(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	menu := pick(m, append([]string{"id"}, menuCols...)...)
	mas, _ := m["menuAuths"].([]any)
	if len(mas) == 0 {
		menu["open_authy"] = false
	}
	menuID := menu["id"] // Java: 新增时实体 id 为 null,menu_authy.menu_id 亦 null,行为保持
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		for _, row := range mas {
			if rm, ok := row.(map[string]any); ok {
				vals := pick(rm, "menu_id", "authy_id")
				vals["menu_id"] = menuID
				if _, err := platform.TxInsert(tx, "chain_auth_menu_authy", vals); err != nil {
					return err
				}
			}
		}
		_, err := platform.TxInsert(tx, "chain_auth_menu", menu)
		return err
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// updMenuLoose = ChainAuthMenuServiceImpl.updMenu:menuAuths 非空则 delByMenu+重插(每行 menuId 覆盖),
// 否则 openAuthy=false;最后 updateById
func (s *Service) updMenuLoose(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	menu := pick(m, append([]string{"id"}, menuCols...)...)
	mas, _ := m["menuAuths"].([]any)
	if len(mas) == 0 {
		menu["open_authy"] = false
	}
	menuID := m["id"]
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		if len(mas) > 0 {
			if err := platform.TxDeleteByMap(tx, "chain_auth_menu_authy", map[string]any{"menu_id": menuID}); err != nil {
				return err
			}
			for _, row := range mas {
				if rm, ok := row.(map[string]any); ok {
					vals := pick(rm, "menu_id", "authy_id")
					vals["menu_id"] = menuID
					if _, err := platform.TxInsert(tx, "chain_auth_menu_authy", vals); err != nil {
						return err
					}
				}
			}
		}
		return platform.TxUpdateByID(tx, "chain_auth_menu", menuID, menu)
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// delMenuNoCascade = controller.delMenu:仅 removeById,不级联 menu_authy(行为保持)
func (s *Service) delMenuNoCascade(w http.ResponseWriter, r *http.Request) {
	if err := platform.DeleteByID(s.DB.Chain, "chain_auth_menu", bodyMap(r)["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// insertMenuAuthys 按 ids[] 插 chain_auth_menu_authy(Java ids 为 null 时 NPE 500,Go 侧跳过不崩)
func insertMenuAuthys(tx *sql.Tx, menuID float64, ids any) error {
	list, _ := ids.([]any)
	for _, v := range list {
		if _, err := platform.TxInsert(tx, "chain_auth_menu_authy",
			map[string]any{"menu_id": menuID, "authy_id": v}); err != nil {
			return err
		}
	}
	return nil
}
