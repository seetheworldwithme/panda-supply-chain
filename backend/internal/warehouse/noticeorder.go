package warehouse

// 六期批3 bespoke:busNotice / busNoticeMx / busAdvertisement / busEnterpriseLeasing / busOrder 域
// 全部语义对照 Java 8765 活体实测(每 handler 注释附实测要点);注册经 addExtra 挂 extraRegs。

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// javaNowMillis 对照 System.currentTimeMillis()(订单号)
func javaNowMillis() string { return strconv.FormatInt(time.Now().UnixMilli(), 10) }

func init() {
	for i := range tplCtls {
		if tplCtls[i].Path == "/busOrder" {
			tplCtls[i].Skip = append(tplCtls[i].Skip, "/list")
		}
	}
	addExtra(func(mux *http.ServeMux, s *Service) {
		reg := func(path string, h http.HandlerFunc) {
			mux.HandleFunc("POST /saas-warehouse"+path, authGuard(h))
		}
		// busNotice
		reg("/busNotice/selAllGrounding", s.noticeSelAllGrounding)
		reg("/busNotice/listByPage", s.noticeLbListByPage)
		reg("/busNotice/getById", s.noticeGetById)
		// busNoticeMx
		reg("/busNoticeMx/listMx", s.mxListMx)
		reg("/busNoticeMx/getCard", s.mxGetCard)
		// busAdvertisement
		reg("/busAdvertisement/list", s.advFindList)
		reg("/busAdvertisement/selAllGrounding", s.advSelAllGrounding)
		reg("/busAdvertisement/getById", s.advGetById)
		reg("/busAdvertisement/deleteById", s.advDeleteById)
		// busEnterpriseLeasing
		reg("/busEnterpriseLeasing/listByPage", s.elListByPage)
		reg("/busEnterpriseLeasing/like", s.elLike)
		// busVideos(listByPage/getById/deleteById)
		reg("/busVideos/listByPage", s.vidLbListByPage)
		reg("/busVideos/getById", s.vidGetById)
		reg("/busVideos/deleteById", s.vidDeleteById)
		// busOrder
		reg("/busOrder/list", s.orderList)
		reg("/busOrder/listByPage", s.orderListByPage)
		reg("/busOrder/getById", s.orderGetById)
		reg("/busOrder/placeOrder", s.placeOrder)
	})
}

// ── busNotice ──────────────────────────────────────────────────────────────

// noticeSelAllGrounding 实测{}→999:mapper 只收 Map,SQL 内 ${searchSql} 从 map 取键,
// 缺键 → MyBatis 绑定异常 → 999;仅当 body 显式带 searchSql 字符串时才拼 SQL(200)。
func (s *Service) noticeSelAllGrounding(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ss, _ := m["searchSql"].(string)
	if ss == "" {
		platform.Unknown(w)
		return
	}
	search := mapSearchArg(m)
	q := "select * from bus_notice " + ss + " and CONCAT(id) like ? and end_date>=CURDATE() order by create_time desc"
	rows, cols, err := s.queryEnt(r, q, search)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, entBusNotice.shapeRows(rows, cols))
}

// noticeGetById 实测{"id":1}(缺行)→999:getBusNotice 先 getById 再 busNotice.setFormList
// 缺行 NPE;行存在时 formList = listMx({notice_id:id}).formList(JSONArray,分组结构)。
func (s *Service) noticeGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	rows, cols, err := s.queryEnt(r, "SELECT * FROM bus_notice WHERE id=?", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	shaped := entBusNotice.shapeRows(rows, cols)
	if len(shaped) == 0 {
		platform.Unknown(w) // 缺行 → NPE → 999
		return
	}
	fl, ok := s.listMxRows(r, map[string]any{"notice_id": m["id"]})
	if !ok {
		platform.Unknown(w)
		return
	}
	// Java:busNotice.setFormList(...getJSONArray("formList")) → 取数组本体
	if fm, ism := fl.(map[string]any); ism {
		shaped[0]["formList"] = fm["formList"]
	} else {
		shaped[0]["formList"] = fl
	}
	platform.Success(w, shaped[0])
}

// ── busNoticeMx ────────────────────────────────────────────────────────────

// listMxRows 对照 BusNoticeMxServiceImpl.listMx:allEq(map) 等值条件 + ORDER BY id DESC,
// 按 etcName+uniqueId 分组(对象跨行复用:formList 首现入列,contentList 持续追加,
// id/enterpriseName/uniqueId 逐行覆盖 → 最终为组内最小 id 行的值)。
func (s *Service) listMxRows(r *http.Request, conds map[string]any) (any, bool) {
	q := "SELECT * FROM bus_notice_mx"
	var args []any
	if len(conds) > 0 {
		var cs []string
		for _, k := range platform.SortedKeys(conds) {
			cs = append(cs, k+" = ?")
			args = append(args, conds[k])
		}
		q += " WHERE " + strings.Join(cs, " AND ")
	}
	q += " ORDER BY id DESC"
	rows, cols, err := s.queryEnt(r, q, args...)
	if err != nil {
		return nil, false
	}
	shaped := entBusNoticeMx.shapeRows(rows, cols)
	formList := []any{}
	groups := map[string]map[string]any{}
	for _, row := range shaped {
		key := str(row, "etcName") + str(row, "uniqueId")
		g := groups[key]
		if g == nil {
			g = map[string]any{
				"enterpriseName": row["etcName"],
				"contentList":    []any{},
				"id":             row["id"],
				"uniqueId":       row["uniqueId"],
			}
			groups[key] = g
			formList = append(formList, g)
		}
		// 逐行覆盖(Java 同款:mxObj.put 每行执行,最终值为组内最后一行)
		g["enterpriseName"] = row["etcName"]
		g["id"] = row["id"]
		g["uniqueId"] = row["uniqueId"]
		g["contentList"] = append(g["contentList"].([]any), map[string]any{
			"noticeTitle":   row["titleDescription"],
			"noticeContent": row["titleContent"],
		})
	}
	return map[string]any{"formList": formList}, true
}

// mxListMx 实测{}→200:formList 全表分组;body 键须为 bus_notice_mx 列
// (QueryWrapper.allEq 非列键 → SQL 报错 → 999)。
func (s *Service) mxListMx(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	conds := map[string]any{}
	for _, k := range platform.SortedKeys(m) {
		if k == "page" || k == "ids" || k == "searchStr" {
			continue
		}
		if !isDbcol(entBusNoticeMx, k) {
			platform.Unknown(w) // 非列键 → SQL 报错
			return
		}
		conds[camelCol2snake(k)] = m[k]
	}
	fl, ok := s.listMxRows(r, conds)
	if !ok {
		platform.Unknown(w)
		return
	}
	platform.Success(w, fl)
}

// mxGetCard 实测{}→999:notice=selectById(noticeId) 缺行/缺参 → getNoticeTheme NPE;
// 返回 {noticeTheme,noticeImgId,businessName,businessPhone,mxs:listMx({notice_id,unique_id})}。
func (s *Service) mxGetCard(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	rows, cols, err := s.queryEnt(r, "SELECT * FROM bus_notice WHERE id=?", m["noticeId"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	shaped := entBusNotice.shapeRows(rows, cols)
	if len(shaped) == 0 {
		platform.Unknown(w)
		return
	}
	fl, ok := s.listMxRows(r, map[string]any{"notice_id": m["noticeId"], "unique_id": m["uniqueId"]})
	if !ok {
		platform.Unknown(w)
		return
	}
	platform.Success(w, map[string]any{
		"noticeTheme":   shaped[0]["noticeTheme"],
		"noticeImgId":   shaped[0]["noticeImgId"],
		"businessName":  shaped[0]["businessName"],
		"businessPhone": shaped[0]["businessPhone"],
		"mxs":           fl,
	})
}

// ── busAdvertisement ───────────────────────────────────────────────────────

// mapConds 对照 StringHelper.getSql(Map):body 键即列名,字符串字面注入、数字原样;
// 非列键 → SQL 报错。{} → 空(调用方自带 where 1=1 语义时须自行补)。
func (s *Service) mapConds(m map[string]any, e *ent) (string, bool) {
	var cs []string
	for _, k := range platform.SortedKeys(m) {
		if k == "page" || k == "ids" || k == "searchStr" {
			continue
		}
		if !isDbcol(e, k) {
			return "", false
		}
		switch v := m[k].(type) {
		case string:
			cs = append(cs, k+" = '"+v+"'")
		case nil:
			cs = append(cs, k+" = null")
		default:
			cs = append(cs, k+" = "+javaNumStr(v))
		}
	}
	if len(cs) == 0 {
		return " where 1=1", true // getSql(Map) 无条件时返回 where 1=1(实测:advertisement selAllGrounding {} → SQL 合法 → 200 [])
	}
	return " where 1=1 and " + strings.Join(cs, " and "), true
}

// javaNumStr 数字字面(整数无小数点,浮点保 Java toString 形态)
func javaNumStr(v any) string {
	switch n := v.(type) {
	case json.Number:
		return n.String()
	case int64:
		return strconv.FormatInt(n, 10)
	case float64:
		if n == float64(int64(n)) {
			return strconv.FormatInt(int64(n), 10)
		}
		return strconv.FormatFloat(n, 'f', -1, 64)
	case bool:
		if n {
			return "true"
		}
		return "false"
	}
	return ""
}

func isDbcol(e *ent, k string) bool {
	for _, c := range e.dbcols {
		if c == k {
			return true
		}
	}
	return false
}

func camelCol2snake(k string) string {
	var b strings.Builder
	for i, c := range k {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(c + 32)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

// mapSearchArg Map 变体 getMaps:searchStr 原样(不包 %v%),缺失为 nil(#{maps.searchStr}→NULL)
func mapSearchArg(m map[string]any) any {
	if v, has := m["searchStr"]; has {
		return v
	}
	return nil
}

// advFindList 实测{}→200 全量:mapper findList = select *,(notice_type 子查询)
// ${searchSql} order by create_time desc(无 CONCAT 条款)。
func (s *Service) advFindList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	// Java calls generic getSql(T) with a Map: it inspects HashMap fields,
	// finds no ExcelColumn annotations and produces no predicates.
	for _, v := range m {
		if v == nil {
			platform.Unknown(w)
			return
		}
	}
	cond := " where 1=1"
	q := "select *, (select notice_type from bus_notice where id=notice_id) noticeType from bus_advertisement" + cond + " order by create_time desc"
	rows, cols, err := s.queryEnt(r, q)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, entBusAdvertisement.shapeRows(rows, cols))
}

// advSelAllGrounding 实测{}→200 []:select * ${searchSql} and CONCAT(id) like #{maps.searchStr}
// and end_date>=CURDATE() order by create_time desc;{}→search 为 NULL 恒空 + end_date 已过期。
func (s *Service) advSelAllGrounding(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := s.mapConds(m, entBusAdvertisement)
	if !ok {
		platform.Unknown(w)
		return
	}
	q := "select * from bus_advertisement" + cond + " and CONCAT(id) like ? and end_date>=CURDATE() order by create_time desc"
	rows, cols, err := s.queryEnt(r, q, mapSearchArg(m))
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, entBusAdvertisement.shapeRows(rows, cols))
}

// advGetById 实测{"id":1}(缺行)→999:getBusAdvertisement 链式取 notice.isDetailShow,
// 行缺/notice 缺 → NPE。
func (s *Service) advGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	rows, cols, err := s.queryEnt(r, "SELECT * FROM bus_advertisement WHERE id=?", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	shaped := entBusAdvertisement.shapeRows(rows, cols)
	if len(shaped) == 0 {
		platform.Unknown(w)
		return
	}
	nrows, ncols, err := s.queryEnt(r, "SELECT * FROM bus_notice WHERE id=?", shaped[0]["noticeId"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	nshaped := entBusNotice.shapeRows(nrows, ncols)
	if len(nshaped) == 0 {
		platform.Unknown(w)
		return
	}
	shaped[0]["isDetailShow"] = nshaped[0]["isDetailShow"]
	platform.Success(w, shaped[0])
}

// advDeleteById 实测{"id":1}(缺行)→999:removeBusAdvertisement 先 getById 再
// adv.getNoticeId() 缺行 NPE;存在则删广告行 + 删关联 notice,响应字面"删除成功，id为：N"。
func (s *Service) advDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	rows, _, err := s.queryEnt(r, "SELECT * FROM bus_advertisement WHERE id=?", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(rows) == 0 {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	if _, err := d.Exec("DELETE FROM bus_advertisement WHERE id=?", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("DELETE FROM bus_notice WHERE id=?", rows[0]["notice_id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

// ── busEnterpriseLeasing ───────────────────────────────────────────────────

// fileEnt chain_common_file 实体形状(@Many selfiles 行;对照 ChainCommonFile 字段)
var fileEnt = &ent{
	table: "chain_common_file",
	keys: []string{"id", "companyId", "fileRealName", "fileName", "fileDescription", "filePath", "fType", "fileType",
		"ownerType", "ownerId", "relationId", "relationType", "effective", "transStatus", "associationId",
		"isCompressedFile", "sourceId", "compressedId", "createTime", "updateTime",
		"remark1", "remark2", "remark3", "remark4", "remark5", "searchStr", "page", "ids", "configFilePath"},
	dbcols: []string{"id", "company_id", "file_real_name", "file_name", "file_description", "file_path", "f_type",
		"file_type", "owner_type", "owner_id", "relation_id", "relation_type", "effective", "trans_status",
		"association_id", "is_compressed_file", "source_id", "compressed_id", "create_time", "update_time",
		"remark1", "remark2", "remark3", "remark4", "remark5"},
	q:     nil,
	excel: nil,
	types: []jf{
		{"id", "i"}, {"companyId", "i"}, {"fileRealName", "s"}, {"fileName", "s"}, {"fileDescription", "s"},
		{"filePath", "s"}, {"fType", "i"}, {"fileType", "s"}, {"ownerType", "s"}, {"ownerId", "i"},
		{"relationId", "i"}, {"relationType", "i"}, {"effective", "b"}, {"transStatus", "i"},
		{"associationId", "s"}, {"isCompressedFile", "b"}, {"sourceId", "i"}, {"compressedId", "i"},
		{"createTime", "t"}, {"updateTime", "t"},
		{"remark1", "s"}, {"remark2", "s"}, {"remark3", "s"}, {"remark4", "s"}, {"remark5", "s"},
		{"searchStr", "s"}, {"page", "page"}, {"ids", "list"}, {"configFilePath", "s"},
	},
}

// elListByPage 实测:Java 走 QueryWrapper selectPage(非 mapper listSQL!):
// getWrapper 恒 ORDER BY update_time DESC + service 追加 area ASC;子查询见下;photoList=@Many selfiles。
func (s *Service) elListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusEnterpriseLeasing.validateBody(m) {
		platform.Unknown(w)
		return
	}
	uid := currentUser(r).UserID
	cond, args, ok := entBusEnterpriseLeasing.getSql(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	// family C:startTime/endTime DATE_FORMAT 字面注入(service 95-100 行)
	for _, kv := range []struct{ k, op string }{{"startTime", ">= '"}, {"endTime", "<= '"}} {
		if v, has := m[kv.k]; has && v != nil {
			d, ok2 := bindVal("D", v)
			if !ok2 || d == nil {
				continue
			}
			cond += " and DATE_FORMAT(create_time, '%Y-%m-%d')" + kv.op + d.(string) + "'"
		}
	}
	// Java:select el.*, (wxFaceImg 子查询), (collCount by maps.userId), (collTotCount),
	// (callCount) from bus_enterprise_leasing el ${searchSql} and CONCAT(ifnull(enterprise_name,''),
	// ifnull(cinfo,''), ifnull(address,''), ifnull(coal_rel_name,''), ifnull(coal_rel_phone,''))
	// like #{maps.searchStr} order by create_time desc
	q := "select el.*, (select wx_face_img from chain_auth_user where id=el.create_user_id) wxFaceImg, " +
		"(SELECT IFNULL(count(1), 0) FROM bus_pro_collection WHERE el_id = el.id and user_id=?) collCount, " +
		"(SELECT IFNULL(count(1), 0) FROM bus_pro_collection WHERE el_id = el.id) collTotCount, " +
		"(select IFNULL(count(1),0) from bus_phone_log where lease_id=el.id) callCount " +
		"from bus_enterprise_leasing el " + cond +
		" and CONCAT( ifnull(enterprise_name,''), ifnull(cinfo,''), ifnull(address,''), ifnull(coal_rel_name,''), ifnull(coal_rel_phone,'')) like ? " +
		"order by create_time desc"
	search := mapSearchArg(m)
	searchStr := "%null%"
	if t, has := m["searchStr"]; has && t != nil {
		if sv, ok2 := bindVal("s", t); ok2 {
			searchStr = "%" + sv.(string) + "%"
			search = searchStr
		}
	}
	allArgs := append([]any{uid}, args...)
	allArgs = append(allArgs, search)
	countQ := "SELECT COUNT(*) FROM (" + stripOrderBy(q) + ") cnt"
	current, size := int64(1), int64(10)
	pageObj, hasPage := m["page"].(map[string]any)
	if hasPage {
		if v, ok := pageObj["current"]; ok && v != nil {
			current = int64(num(v))
		}
		if v, ok := pageObj["size"]; ok && v != nil {
			size = int64(num(v))
		}
	}
	d := s.pick(r)
	if !hasPage { // listNoPage:无 photoList/子查询,普通列 + CONCAT(id) like
		q2 := "select * from bus_enterprise_leasing " + cond + " and CONCAT(id) like ? order by create_time desc"
		rows, cols, err := s.queryEnt(r, q2, append(args, search)...)
		if err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, entBusEnterpriseLeasing.shapeRows(rows, cols))
		return
	}
	var total int64
	if err := d.QueryRow(countQ, allArgs...).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	off := (current - 1) * size
	rows, cols, err := s.queryEnt(r, q+" LIMIT ? OFFSET ?", append(allArgs, size, off)...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	shaped := entBusEnterpriseLeasing.shapeRows(rows, cols)
	// photoList = @Many selfiles(photo_association_id):恒为数组(空匹配 → [])
	for _, row := range shaped {
		frows, fcols, err := s.queryEnt(r, "SELECT t.* from chain_common_file t where association_id = ?", row["photoAssociationId"])
		if err != nil {
			platform.Unknown(w)
			return
		}
		row["photoList"] = fileEnt.shapeRows(frows, fcols)
	}
	platform.Success(w, ipageOf(shaped, total, current, size))
}

// elLike 实测{}→999:like() 内 getLikeCount()+1 拆箱 NPE;有值则非空子集 updateById。
func (s *Service) elLike(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusEnterpriseLeasing.validateBody(m) {
		platform.Unknown(w)
		return
	}
	lc := m["likeCount"]
	if lc == nil {
		platform.Unknown(w) // null+1 → NPE
		return
	}
	vals := pick(m, entBusEnterpriseLeasing.dbcols)
	delete(vals, "id")
	vals["like_count"] = num(lc) + 1
	if _, err := platform.UpdateByID(s.pick(r), entBusEnterpriseLeasing.table, m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功") // 控制器 success("删除成功")(实测)
}

// ── busOrder ───────────────────────────────────────────────────────────────

// orderListByPage 实测{page,searchStr}→records 空列表:getSql(entity,"status") 排除 status,
// status 非空再原样注入;isPtSearch 空/false → and bus_order.create_phone = '<user.phone>'。
func (s *Service) orderListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["page"].(map[string]any); !ok {
		platform.Unknown(w) // Java mapper requires a Page object.
		return
	}
	if !entBusOrder.validateBody(m) {
		platform.Unknown(w)
		return
	}
	cm := map[string]any{}
	for k, v := range m {
		if k == "status" {
			continue // getSql(entity,"status") 排除
		}
		cm[k] = v
	}
	cond, args, ok := entBusOrder.getSql(cm)
	if !ok {
		platform.Unknown(w)
		return
	}
	if m["status"] != nil {
		cond += " and bus_order.status = " + javaNumStr(num(m["status"]))
	}
	pt := false
	if b, isb := m["isPtSearch"].(bool); isb {
		pt = b
	}
	if !pt {
		phone, ok2 := s.userPhone(r)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		cond += " and bus_order.create_phone = '" + phone + "'"
	}
	q := strings.Replace(entBusOrder.listSQL, "${searchSql}", cond, 1)
	search := "%null%"
	if v, has := m["searchStr"]; has && v != nil {
		if sv, ok2 := bindVal("s", v); ok2 {
			search = sv.(string) // searchLit:SQL 已含 '%...%' 包裹
		}
	}
	q = strings.Replace(q, "${maps.searchStr}", search, 1)
	countQ := "SELECT COUNT(*) FROM (" + stripOrderBy(q) + ") cnt"
	current, size := int64(1), int64(10)
	pageObj, hasPage := m["page"].(map[string]any)
	if hasPage {
		if v, ok := pageObj["current"]; ok && v != nil {
			current = int64(num(v))
		}
		if v, ok := pageObj["size"]; ok && v != nil {
			size = int64(num(v))
		}
	}
	d := s.pick(r)
	var total int64
	if err := d.QueryRow(countQ, args...).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	if !hasPage {
		rows, cols, err := s.queryEnt(r, q, args...)
		if err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, entBusOrder.shapeRows(rows, cols))
		return
	}
	off := (current - 1) * size
	rows, cols, err := s.queryEnt(r, q+" LIMIT ? OFFSET ?", append(args, size, off)...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, ipageOf(entBusOrder.shapeRows(rows, cols), total, current, size))
}

// userPhone 对照 Java user.getPhone()(JWT):取 chain_auth_user.phone
func (s *Service) userPhone(r *http.Request) (string, bool) {
	uid := currentUser(r).UserID
	var phone string
	if err := s.pick(r).QueryRow("SELECT phone FROM chain_auth_user WHERE id=?", uid).Scan(&phone); err != nil {
		return "", false
	}
	return phone, true
}

// orderGetById 实测{"id":1}(缺行)→999:getBusOrder 链式取 product.warehouseNb,缺行 NPE。
func (s *Service) orderGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	rows, cols, err := s.queryEnt(r, "SELECT * FROM bus_order WHERE id=?", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	shaped := entBusOrder.shapeRows(rows, cols)
	if len(shaped) == 0 {
		platform.Unknown(w)
		return
	}
	prows, pcols, err := s.queryEnt(r, "SELECT * FROM bus_product WHERE id=?", shaped[0]["proId"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	pshaped := entBusProduct.shapeRows(prows, pcols)
	if len(pshaped) == 0 {
		platform.Unknown(w)
		return
	}
	shaped[0]["warehouseNb"] = pshaped[0]["warehouseNb"]
	platform.Success(w, shaped[0])
}

// placeOrder 实测{}→999 且无新行:save() 先插行,随后 pro=getById(proId) 缺失 →
// params 组装 NPE → 全局事务回滚；有产品时按原配置发送通知。
func (s *Service) placeOrder(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusOrder.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	// 短信配置存在性决定是否走到产品校验(Java:cfgs.size()>0 才取 pro)
	var cfgs int
	if err := d.QueryRow("SELECT COUNT(*) FROM bus_config WHERE msg_type=1").Scan(&cfgs); err != nil {
		platform.Unknown(w)
		return
	}
	phone, ok := s.userPhone(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	vals := pick(m, entBusOrder.dbcols)
	vals["create_user_id"] = currentUser(r).UserID
	vals["create_phone"] = phone
	vals["status"] = int64(0)
	vals["order_nb"] = javaNowMillis()
	now := time.Now().Format("2006-01-02 15:04:05")
	vals["create_time"] = now // @TableField(fill=INSERT)
	vals["update_time"] = now // @TableField(fill=INSERT_UPDATE)
	tx, err := d.Begin()
	if err != nil {
		platform.Unknown(w)
		return
	}
	cols := platform.SortedKeys(vals)
	ph := strings.TrimSuffix(strings.Repeat("?,", len(cols)), ",")
	var args []any
	for _, k := range cols {
		args = append(args, vals[k])
	}
	if _, err := tx.Exec("INSERT INTO bus_order ("+strings.Join(cols, ",")+") VALUES ("+ph+")", args...); err != nil {
		tx.Rollback()
		platform.Unknown(w)
		return
	}
	if cfgs > 0 {
		// Java:pro=getById(proId) 缺失 → pro.getWarehouseNb() NPE → 回滚 → 999 无行
		var pid any
		if err := tx.QueryRow("SELECT id FROM bus_product WHERE id=?", vals["pro_id"]).Scan(&pid); err != nil {
			tx.Rollback()
			platform.Unknown(w)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		platform.Unknown(w)
		return
	}
	if cfgs > 0 {
		product, _ := platform.GetByID(d, "bus_product", vals["pro_id"])
		configs, _ := queryRows(d, "SELECT admin_phone FROM bus_config WHERE msg_type=1")
		for _, config := range configs {
			platform.SendLegacySMS(r.Context(), javaStr(config["admin_phone"]), []string{phone, javaStr(product["warehouse_nb"])}, "1990680", "")
		}
	}
	platform.Success(w, nil)
}

// ── busNotice/listByPage( advertisements @Many 关联) ─────────────────────────

// noticeLbListByPage = busNotice/listByPage:模板管道 + 每行 advertisements =
// selAdvers(@Many):SELECT t.* from bus_advertisement t where notice_id = ${id}(id 字面注入)。
func (s *Service) noticeLbListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["page"].(map[string]any); !ok {
		platform.Unknown(w) // Java mapper requires a Page object.
		return
	}
	if !entBusNotice.validateBody(m) {
		platform.Unknown(w)
		return
	}
	cond, args, ok := entBusNotice.getSql(m)
	if !ok {
		platform.Unknown(w)
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
	q := strings.Replace(entBusNotice.listSQL, "${searchSql}", cond, 1)
	q = strings.Replace(q, "#{maps.searchStr}", "?", 1)
	baseArgs := append([]any{}, args...)
	baseArgs = append(baseArgs, search)
	countQ := "SELECT COUNT(*) FROM (" + stripOrderBy(q) + ") cnt"
	current, size := int64(1), int64(10)
	pageObj, hasPage := m["page"].(map[string]any)
	if hasPage {
		if v, ok := pageObj["current"]; ok && v != nil {
			current = int64(num(v))
		}
		if v, ok := pageObj["size"]; ok && v != nil {
			size = int64(num(v))
		}
	}
	d := s.pick(r)
	var total int64
	if err := d.QueryRow(countQ, baseArgs...).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	off := (current - 1) * size
	rows, cols, err := s.queryEnt(r, q+" LIMIT ? OFFSET ?", append(baseArgs, size, off)...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	shaped := entBusNotice.shapeRows(rows, cols)
	for _, row := range shaped {
		arows, acols, err := s.queryEnt(r, "SELECT t.* from bus_advertisement t where notice_id = ?", row["id"])
		if err != nil {
			platform.Unknown(w)
			return
		}
		row["advertisements"] = entBusAdvertisement.shapeRows(arows, acols)
	}
	platform.Success(w, ipageOf(shaped, total, current, size))
}

// ── busVideos getById / deleteById / listByPage ────────────────────────────

// vidGetById = BusVideosServiceImpl.getBusVideos:getById 缺行 → getProId() NPE → 999;
// proId 非空挂 busProduct;红包分支需 shareUserId+当天新访客({}体不触发,触发路径登记差异)。
func (s *Service) vidGetById(w http.ResponseWriter, r *http.Request) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	rows, cols, err := s.queryEnt(r, "SELECT * FROM bus_videos WHERE id=?", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	shaped := entBusVideos.shapeRows(rows, cols)
	if len(shaped) == 0 {
		platform.Unknown(w) // 缺行 → NPE
		return
	}
	if shaped[0]["proId"] != nil {
		prows, pcols, err := s.queryEnt(r, "SELECT * FROM bus_product WHERE id=?", shaped[0]["proId"])
		if err != nil {
			platform.Unknown(w)
			return
		}
		pshaped := entBusProduct.shapeRows(prows, pcols)
		if len(pshaped) > 0 {
			shaped[0]["busProduct"] = pshaped[0]
		}
	}
	if err := s.videoVisit(r, m); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shaped[0])
}

// vidDeleteById = removeBusVideos:getById 缺行 → vdo.getVideoId() NPE → 999;存在则
// removeById + removeFile 链(与 vDelByVideoId 的 ② 段同构,file id = 行 video_id)。
func (s *Service) vidDeleteById(w http.ResponseWriter, r *http.Request) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	rows, err := queryRows(d, "SELECT * FROM bus_videos WHERE id=?", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(rows) == 0 {
		platform.Unknown(w) // 缺行 → NPE
		return
	}
	vid := rows[0]["video_id"]
	if _, err := d.Exec("DELETE FROM bus_videos WHERE id=?", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("DELETE FROM chain_common_file WHERE id=?", vid); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("UPDATE chain_basic_enterprise SET advertising_video_id=NULL WHERE advertising_video_id=?", vid); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("UPDATE chain_basic_enterprise SET website_img_id=NULL, website_ced_img_id=NULL WHERE website_img_id=? OR website_ced_img_id=?", vid, vid); err != nil {
		platform.Unknown(w)
		return
	}
	frow, err := platform.GetByID(d, "chain_common_file", vid)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if frow != nil {
		fp, fn := str(frow, "file_path"), str(frow, "file_name")
		if fp != "" && fn != "" {
			dot := strings.Index(fn, ".")
			if dot < 0 {
				platform.Unknown(w) // Java substring 越界
				return
			}
			_ = os.Remove(fp + "/" + fn)
			_ = os.RemoveAll(fp + "/" + fn[:dot] + "_img")
		}
	}
	irows, err := queryRows(d, "SELECT file_path, file_name FROM chain_common_img WHERE file_id=?", vid)
	if err != nil {
		platform.Unknown(w)
		return
	}
	for _, ir := range irows {
		_ = os.Remove(str(ir, "file_path") + "/" + str(ir, "file_name"))
	}
	if _, err := d.Exec("DELETE FROM chain_common_file WHERE source_id=?", vid); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

// vidLbListByPage = BusVideosMapper.listByPage/listByPage2(巨型别名 JOIN,Map 行原样):
// seeUserId 非空 → listByPage(if* 子查询,userId=seeUserId)+ 访客记录副作用;否则 listByPage2。
func (s *Service) vidLbListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["page"].(map[string]any); !ok {
		platform.Unknown(w) // Java mapper requires a Page object.
		return
	}
	if !entBusVideos.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	uid := currentUser(r).UserID
	// 访客记录副作用(seeUserId 非空且 ≠ 当前用户)
	if sv, has := m["seeUserId"]; has && sv != nil {
		seeUID := int64(num(sv))
		if seeUID != uid {
			vrows, err := queryRows(d, "SELECT * FROM bus_visitor_log WHERE visitor_id=?", seeUID)
			if err != nil {
				platform.Unknown(w)
				return
			}
			if len(vrows) > 0 {
				if _, err := d.Exec("UPDATE bus_visitor_log SET update_time=NOW() WHERE id=?", vrows[0]["id"]); err != nil {
					platform.Unknown(w)
					return
				}
			} else {
				urow, err := queryRows(d, "SELECT wx_face_img, nickname FROM chain_auth_user WHERE id=?", seeUID)
				if err != nil {
					platform.Unknown(w)
					return
				}
				if len(urow) == 0 || urow[0]["wx_face_img"] == nil || urow[0]["nickname"] == nil {
					platform.Unknown(w) // Java userMap.get(...).toString() NPE
					return
				}
				if _, err := d.Exec("INSERT INTO bus_visitor_log (visitor_id, target_user_id, visitor_headimg, visitor_nickname) VALUES (?,?,?,?)",
					uid, seeUID, urow[0]["wx_face_img"], urow[0]["nickname"]); err != nil {
					platform.Unknown(w)
					return
				}
			}
		}
	}
	cm := map[string]any{}
	for k, v := range m {
		if k == "status" {
			continue // getSql(entity,"status")
		}
		cm[k] = v
	}
	cond, args, ok := entBusVideos.getSql(cm)
	if !ok {
		platform.Unknown(w)
		return
	}
	if m["status"] != nil {
		cond += " and v.status=" + javaNumStr(num(m["status"]))
	}
	search := "%null%"
	if v, has := m["searchStr"]; has && v != nil {
		if sv, ok2 := bindVal("s", v); ok2 {
			search = "%" + sv.(string) + "%"
		}
	}
	// seeUserId 非空走 listByPage(if* 子查询 user_id=${maps.userId} 字面);否则 listByPage2
	var q string
	var largs []any
	if sv, has := m["seeUserId"]; has && sv != nil {
		q = strings.Replace(videosLbSQL(true), "${maps.userId}", strconv.FormatInt(int64(num(sv)), 10), 3)
	} else {
		q = videosLbSQL(false)
	}
	q = strings.Replace(q, "${searchSql}", cond, 1)
	q = strings.Replace(q, "#{maps.searchStr}", "?", 1)
	largs = append(largs, args...)
	largs = append(largs, search)
	countQ := "SELECT COUNT(*) FROM (" + stripOrderBy(q) + ") cnt"
	current, size := int64(1), int64(10)
	pageObj, hasPage := m["page"].(map[string]any)
	if hasPage {
		if v, ok := pageObj["current"]; ok && v != nil {
			current = int64(num(v))
		}
		if v, ok := pageObj["size"]; ok && v != nil {
			size = int64(num(v))
		}
	}
	var total int64
	if err := d.QueryRow(countQ, largs...).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	off := (current - 1) * size
	rows, err := queryRows(d, q+" LIMIT ? OFFSET ?", append(largs, size, off)...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, ipageOf(vidShapeRows(d, rows), total, current, size))
}

// videosLbSQL busVideos/listByPage 原文(withIf=true → 含 ifLike/ifCollect/ifCare 子查询)
func videosLbSQL(withIf bool) string {
	base := "SELECT v.id as id, v.user_id userId, video_id as videoId,video_name as videoName, video_title as videoTitle,like_num as likeNum,collect_num as collectNum,u.phone as phone,v.status as status,v.check_remark as checkRemark,pro_id as proId,v.create_time as createTime,v.update_time as updateTime,pro_describe as proDescribe,u.nickname as nickname,u.wx_face_img as userHeadImg,forward_num as forwardNum,play_num as playNum, warehouse_id as warehouseId,etc_pro_id as etcProId,warehouse_nb as warehouseNb,enterprise_name as enterpriseName,fst_pro_id as fstProId,big_pro_id as bigProId,three_pro_id as threeProId,region as region,product_name as productName,ble_tonnage as bleTonnage,available_tonnage as availableTonnage,min_granularity as minGranularity,max_granularity as maxGranularity,min_ash_content as minAshContent,max_ash_content as maxAshContent,is_show_ash_content as isShowAshContent,min_coal_ash as minCoalAsh,max_coal_ash as maxCoalAsh,min_volatiles as minVolatiles,max_volatiles as maxVolatiles,is_show_volatiles as isShowVolatiles,min_std as minStd,max_std as maxStd,is_show_std as isShowStd,min_internal_sulfur as minInternalSulfur,max_internal_sulfur as maxInternalSulfur,min_water_content as minWaterContent,max_water_content as maxWaterContent,is_show_water_content as isShowWaterContent,min_bond_index as minBondIndex,max_bond_index as maxBondIndex,is_show_bond_index as isShowBondIndex,min_calorific_value as minCalorificValue,max_calorific_value as maxCalorificValue,min_csr as minCsr,max_csr as maxCsr,is_show_csr as isShowCsr,min_y as minY,max_y as maxY,is_show_y as isShowY,min_hgi as minHgi,max_hgi as maxHgi,is_show_hgi as isShowHgi,recovery_analysis as recoveryAnalysis,is_show_recovery_analysis as isShowRecoveryAnalysis,lithofacies as lithofacies,is_show_lithofacies as isShowLithofacies,calorific_value as calorificValue,is_show_calorific_value as isShowCalorificValue,remark as remark,resource_number as resourceNumber,laboratory as laboratory,sampling as sampling,test_number as testNumber,type_of_shipping as typeOfShipping,send_addres as sendAddres,pre_ex_factory_price as preExFactoryPrice,ex_factory_price as exFactoryPrice,plate_price as platePrice,last_test_date as lastTestDate,last_transaction_time as lastTransactionTime,last_tonnage_traded as lastTonnageTraded,last_transaction_price as lastTransactionPrice,ckhy_file_name as ckhyFileName,data_provider as dataProvider,resource_type as resourceType,business_name as businessName,business_phone as businessPhone,notice_time as noticeTime,is_recommende as isRecommende,create_user_id as createUserId,create_phone as createPhone,sharing_times as sharingTimes,view_times as viewTimes,association_id as associationId,img_ids as imgIds,msg_ids as msgIds,is_platform_certification as isPlatformCertification, coal_rel_name coalRelName, coal_rel_phone coalRelPhone, (select product_name from bus_big_prodect where id=p.three_pro_id) thName, "
	if withIf {
		base += "(SELECT count(1) from bus_user_video_lc where video_id=v.id and type=1 and user_id=${maps.userId}) ifLike, (SELECT count(1) from bus_user_video_lc where video_id=v.id and type=2 and user_id=${maps.userId}) ifCollect, (SELECT count(1) from bus_user_care where care_user_id=v.user_id and user_id=${maps.userId}) ifCare, "
	}
	base += "(select is_show_etc_name  from bus_member_benefits b, chain_auth_user uu where b.member_type=uu.member_type and uu.id=p.create_user_id) isShowEtcName"
	if withIf {
		// listByPage(seeUserId 非空)才带 isShowEtcPhone;listByPage2 实测无此列
		base += ", (select is_show_etc_phone  from bus_member_benefits b, chain_auth_user uu where b.member_type=uu.member_type and uu.id=p.create_user_id) isShowEtcPhone"
	}
	base += " from bus_videos v LEFT JOIN bus_product p on v.pro_id=p.id left join chain_auth_user u on v.user_id=u.id ${searchSql} and CONCAT(IFNULL(p.region, ''),IFNULL(p.product_name, ''),IFNULL(p.resource_number, ''),IFNULL(p.send_addres, '')) like #{maps.searchStr} order by v.create_time desc"
	return base
}

// orderList preserves the Map branch: any nonnull isPtSearch selects the mapper,
// while its absence uses MP allEq with null-to-IS-NULL semantics.
func (s *Service) orderList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	q := "SELECT * FROM bus_order WHERE 1=1"
	var args []any
	if m["isPtSearch"] != nil {
		fields := entFieldSet(entBusOrder)
		for _, k := range platform.SortedKeys(m) {
			if fields[k] {
				q += " AND " + camelCol2snake(k) + " LIKE ?"
				args = append(args, m[k])
			}
		}
		search := "null"
		if m["searchStr"] != nil {
			search = javaStr(m["searchStr"])
		}
		q += " AND CONCAT(bus_order.id) LIKE '%" + search + "%' ORDER BY create_time DESC"
	} else {
		for _, k := range platform.SortedKeys(m) {
			if !isDbcol(entBusOrder, k) {
				platform.Unknown(w)
				return
			}
			if m[k] == nil {
				q += " AND " + k + " IS NULL"
			} else {
				q += " AND " + k + " = ?"
				args = append(args, m[k])
			}
		}
	}
	rows, cols, err := s.queryEnt(r, q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, entBusOrder.shapeRows(rows, cols))
}
