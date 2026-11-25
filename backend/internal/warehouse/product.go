// product.go 煤种域 bespoke 端点(批1)。BusBigProdect 三端点照 Java
// BusBigProdectServiceImpl.findList/allProductType/allThirdType;
// ${} 注入语义原样保留(maps 值未包 %v%,getSql(Map) 变体按 body 键绑定)。
package warehouse

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		reg := func(path string, h http.HandlerFunc) {
			mux.HandleFunc("POST /saas-warehouse"+path, authGuard(h))
		}
		reg("/busBigProdect/list", s.bigFindList)
		reg("/busBigProdect/allProductType", s.bigAllProductType)
		reg("/busBigProdect/allThirdType", s.bigAllThirdType)
	})
}

// mapGetSql3 对照 StringHelper.getSql(Map,Class,unstr):body 键即字段名(须有
// @ExcelColumn 且不在 unstr)→ ` and col like #{maps.key}`,值原样(不包 %v%)
func (s *Service) mapGetSql3(e *ent, m map[string]any, unstr string) (string, []any, bool) {
	sqlStr := "where 1=1"
	var args []any
	for k, v := range m {
		if v == nil || strings.Contains(unstr, k) {
			continue
		}
		if !e.knownKey(k) {
			continue
		}
		if e.typeOf(k) == "x" || k == "searchStr" || k == "page" || k == "ids" {
			continue // checkFields:非实体字段不拼
		}
		bound, ok := bindVal(e.typeOf(k), v)
		if !ok {
			return "", nil, false
		}
		if bound == nil {
			continue
		}
		sqlStr += " and " + camelOf(e, k) + " like ?"
		args = append(args, bound)
	}
	return sqlStr, args, true
}

// bigFindList = findList(maps):EXISTS bus_product 的在售煤种树;
// searchType 有值切 findType(order by order_nb)
func (s *Service) bigFindList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	e := entBusBigProdect
	// getSql(maps, "unstr") 命中的是泛型 T 重载:HashMap 无 @ExcelColumn 字段 → 恒 "where 1=1"
	sqlStr, args := "where 1=1", []any{}
	secSql := " (big_pro_id = bus_big_prodect.id or three_pro_id = bus_big_prodect.id)"
	if v, has := m["parent_id"]; has && v != nil {
		sqlStr += " and parent_id like '%," + sqlEscape(str(m, "parent_id")) + ",%'"
	}
	if v, has := m["fst_pro_id"]; has && v != nil {
		sqlStr += " and fst_pro_id like '%" + sqlEscape(str(m, "fst_pro_id")) + "%'"
		secSql += " and fst_pro_id = " + str(m, "fst_pro_id")
	} else {
		secSql += " or fst_pro_id = bus_big_prodect.id"
	}
	var q string
	if _, has := m["searchType"]; has && m["searchType"] != nil {
		q = "select * from bus_big_prodect " + sqlStr + " order by order_nb asc"
	} else {
		q = "select * from bus_big_prodect " + sqlStr + " and EXISTS (SELECT 1 from bus_product where " + secSql + ") order by create_time desc"
	}
	rows, err := queryCam(s.pick(r), q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, e.shape(rows))
}

// bigAllProductType = allProductType(maps):order by order_nb asc
func (s *Service) bigAllProductType(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	e := entBusBigProdect
	sqlStr, args, ok := s.mapGetSql3(e, m, "parent_id,fst_pro_id,searchType")
	if !ok {
		platform.Unknown(w)
		return
	}
	if v, has := m["parent_id"]; has && v != nil {
		sqlStr += " and parent_id like '%," + sqlEscape(str(m, "parent_id")) + ",%'"
	}
	if v, has := m["fst_pro_id"]; has && v != nil {
		sqlStr += " and fst_pro_id like '%" + sqlEscape(str(m, "fst_pro_id")) + "%'"
	}
	rows, err := queryCam(s.pick(r), "select * from bus_big_prodect "+sqlStr+" order by order_nb asc", args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, e.shape(rows))
}

// bigAllThirdType = allThirdType(maps,jwtInfo):existWhere ${} 注入 + 空结果不加推荐行;
// 首行"推荐"是 new BusBigProdect(0,"推荐") —— 仅 id/productName 有值
func (s *Service) bigAllThirdType(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	existWhere := "SELECT 1 from bus_product p where bp.id=p.three_pro_id and status=1"
	if _, has := m["isGroupSearch"]; has && m["isGroupSearch"] != nil {
		existWhere += " and (company_id = '" + sqlEscape(str(m, "groupId")) + "' or group_id = '" + sqlEscape(str(m, "groupId")) + "')"
	}
	if _, has := m["wareIdentity"]; has && m["wareIdentity"] != nil {
		existWhere += " and ware_identity=" + str(m, "wareIdentity")
	}
	rows, err := queryCam(s.pick(r),
		"select bp.* from bus_big_prodect bp where EXISTS ("+existWhere+") order by order_nb asc")
	if err != nil {
		platform.Unknown(w)
		return
	}
	out := entBusBigProdect.shape(rows)
	if len(out) > 0 {
		rec := map[string]any{"id": jsonNumberInt64(0), "productName": "推荐"}
		for _, k := range entBusBigProdect.keys {
			if k != "id" && k != "productName" {
				rec[k] = nil
			}
		}
		out = append([]map[string]any{rec}, out...)
	}
	platform.Success(w, out)
}

// ── BusProCollection bespoke(listPro/proByPage/cancelCollection/selLikeCount/add)──

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		reg := func(path string, h http.HandlerFunc) {
			mux.HandleFunc("POST /saas-warehouse"+path, authGuard(h))
		}
		reg("/busProCollection/add", s.collAdd)
		reg("/busProCollection/listPro", s.collListPro)
		reg("/busProCollection/proByPage", s.collProByPage)
		reg("/busProCollection/cancelCollection", s.collCancel)
		reg("/busProCollection/selLikeCount", s.collSelLikeCount)
	})
}

// collMaps listPro/proByPage 共用:getSql(实体 @ExcelColumn)+ userId(jwt;body userId 优先)
func (s *Service) collMaps(w http.ResponseWriter, r *http.Request, e *ent) (map[string]any, string, []any, string, bool) {
	m := bodyMap(r)
	if !e.validateBody(m) {
		platform.Unknown(w)
		return nil, "", nil, "", false
	}
	sqlStr, args, ok := e.getSql(m)
	if !ok {
		platform.Unknown(w)
		return nil, "", nil, "", false
	}
	uid := strconv.FormatInt(currentUser(r).UserID, 10)
	if v, has := m["userId"]; has && v != nil {
		bv, ok2 := bindVal("i", v)
		if !ok2 {
			platform.Unknown(w)
			return nil, "", nil, "", false
		}
		if bv != nil {
			uid = bv.(string)
		}
	}
	return m, sqlStr, args, uid, true
}

// collListPro 收藏产品列表(连接 bus_product;DATE_FORMAT 与裸 create_time 重复列,ScanMaps 后者覆盖)
func (s *Service) collListPro(w http.ResponseWriter, r *http.Request) {
	e := entBusProCollection
	_, sqlStr, args, uid, ok := s.collMaps(w, r, e)
	if !ok {
		return
	}
	q := "select c.id as cid, DATE_FORMAT(c.create_time, '%Y-%m-%d %H:%i:%s') create_time, p.*, c.create_time from bus_pro_collection c, bus_product p " +
		sqlStr + " and c.user_id=? and c.pro_id=p.id order by c.create_time desc"
	rows, err := queryCam(s.pick(r), q, append(args, uid)...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, rows)
}

// collProByPage 收藏产品分页(CONCAT(p.id) like searchStr;searchStr 缺失绑 NULL → 空记录)
func (s *Service) collProByPage(w http.ResponseWriter, r *http.Request) {
	e := entBusProCollection
	m, sqlStr, args, uid, ok := s.collMaps(w, r, e)
	if !ok {
		return
	}
	var search any
	if v, has := m["searchStr"]; has && v != nil {
		ss, ok2 := bindVal("s", v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		search = "%" + ss.(string) + "%"
	}
	where := "select c.id as cid, p.* from bus_pro_collection c, bus_product p " + sqlStr +
		" and c.user_id=? and c.pro_id=p.id and CONCAT(p.id) like ? order by c.create_time desc"
	countQ := "SELECT COUNT(*) from bus_pro_collection c, bus_product p " + sqlStr +
		" and c.user_id=? and c.pro_id=p.id and CONCAT(p.id) like ?"
	current, size := pageOf(m)
	page, err := paged(s.pick(r), where, countQ, append(append(args, uid), search), current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, page)
}

// collCancel = cancelCollection:proId 路径(计数→删收藏→删日志→-1)/elId 路径/缺参 0
func (s *Service) collCancel(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProCollection.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	cu := currentUser(r)
	uid := strconv.FormatInt(cu.UserID, 10)
	var likeCount float64
	if v, has := m["proId"]; has && v != nil {
		bv, ok2 := bindVal("i", v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		proId := bv.(string)
		if err := d.QueryRow("SELECT count(*) FROM bus_pro_collection WHERE pro_id = ?", proId).Scan(&likeCount); err != nil {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("delete from bus_pro_collection where user_id=? and pro_id=?", uid, proId); err != nil {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("DELETE FROM bus_pro_log WHERE user_id = ? AND pro_id = ? AND operate_type = 6 AND operate_obj = 4", uid, proId); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, jsonNumberInt64(int64(likeCount)-1))
		return
	}
	if v, has := m["elId"]; has && v != nil {
		bv, ok2 := bindVal("i", v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		elId := bv.(string)
		var cnt int64
		if err := d.QueryRow("SELECT count(*) FROM bus_pro_collection WHERE el_id = ? AND user_id = ?", elId, uid).Scan(&cnt); err != nil {
			platform.Unknown(w)
			return
		}
		var selCount float64
		if err := d.QueryRow("select count(*) from bus_pro_collection where el_id=?", elId).Scan(&selCount); err != nil {
			platform.Unknown(w)
			return
		}
		if cnt == 0 {
			platform.Success(w, jsonNumberInt64(int64(selCount))) // 未收藏 → 原样返回计数
			return
		}
		var rid int64
		if err := d.QueryRow("SELECT id FROM bus_pro_collection WHERE el_id = ? AND user_id = ? ORDER BY id LIMIT 1", elId, uid).Scan(&rid); err != nil {
			platform.Unknown(w) // getOne 多行会抛异常;单行取 id
			return
		}
		if err := platform.DeleteByID(d, "bus_pro_collection", rid); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, jsonNumberInt64(int64(selCount)-1))
		return
	}
	platform.Success(w, jsonNumberInt64(0))
}

// collSelLikeCount = selLikeCount:count(el_id=?);elId 缺失 → allEq null → el_id IS NULL 计数
func (s *Service) collSelLikeCount(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	q := "SELECT count(*) FROM bus_pro_collection WHERE "
	var arg any
	has := false
	if v, ok := m["elId"]; ok && v != nil {
		has = true
		bv, ok2 := bindVal("i", v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		q += "el_id = ?"
		arg = bv.(string)
	} else {
		q += "el_id IS NULL"
	}
	var cnt int64
	var err error
	if has {
		err = s.pick(r).QueryRow(q, arg).Scan(&cnt)
	} else {
		err = s.pick(r).QueryRow(q).Scan(&cnt)
	}
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, jsonNumberInt64(cnt))
}

// collAdd = saveBusProCollection:proId 路径查重(已收藏直接回计数)、回填 targetUserId
// (产品缺 → NPE 999);elId 路径 selCount;落库后嵌套 busProLog 一并入库;响应 ++likeCount
func (s *Service) collAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProCollection.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	uid := strconv.FormatInt(currentUser(r).UserID, 10)
	var likeCount int64
	proId, hasPro := "", false
	if v, has := m["proId"]; has && v != nil {
		bv, ok2 := bindVal("i", v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		if bv != nil {
			proId, hasPro = bv.(string), true
		}
	}
	if hasPro {
		var cnt int64
		if err := d.QueryRow("SELECT count(*) FROM bus_pro_collection WHERE user_id = ? AND pro_id = ?", uid, proId).Scan(&cnt); err != nil {
			platform.Unknown(w)
			return
		}
		if err := d.QueryRow("SELECT count(*) FROM bus_pro_collection WHERE pro_id = ?", proId).Scan(&likeCount); err != nil {
			platform.Unknown(w)
			return
		}
		if cnt > 0 {
			platform.Success(w, jsonNumberInt64(likeCount)) // 已收藏 → 不重复落库
			return
		}
		var targetUser sql.NullInt64
		if err := d.QueryRow("SELECT create_user_id FROM bus_product WHERE id = ?", proId).Scan(&targetUser); err != nil || !targetUser.Valid {
			platform.Unknown(w) // pro null 或 createUserId null → NPE 999
			return
		}
		m["targetUserId"] = json.Number(strconv.FormatInt(targetUser.Int64, 10))
	}
	elId, hasEl := "", false
	if v, has := m["elId"]; has && v != nil {
		bv, ok2 := bindVal("i", v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		if bv != nil {
			elId, hasEl = bv.(string), true
		}
	}
	if hasEl {
		var cnt int64
		if err := d.QueryRow("select count(*) from bus_pro_collection where el_id=?", elId).Scan(&cnt); err != nil {
			platform.Unknown(w)
			return
		}
		likeCount = cnt
	}
	vals := pick(m, entBusProCollection.dbcols)
	delete(vals, "id")
	vals["user_id"] = uid // setUserId(jwt) 强制覆盖
	if _, err := platform.Insert(d, "bus_pro_collection", vals); err != nil {
		platform.Unknown(w)
		return
	}
	if logv, has := m["busProLog"]; has && logv != nil {
		lm, ok2 := logv.(map[string]any)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		lvals := pick(lm, entBusProLog.dbcols)
		delete(lvals, "id")
		if _, err := platform.Insert(d, "bus_pro_log", lvals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, jsonNumberInt64(likeCount+1))
}
