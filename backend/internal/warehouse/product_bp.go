// product_bp.go BusProduct bespoke(批1核心,对照 docs/BusProductController-Go迁移规格书)。
// 本文件:简单读端点 + hisPrice + selCheckPro;listByPage 见 product_list.go;
// getById/写端点见 product_get.go。MQ 不引(发送点登记差异,仅写 bus_mq_consume_log 行为
// 由消费侧决定——本期消费侧不迁移);Redis 用进程内存等价(broadcastMsg 空)。
package warehouse

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		reg := func(path string, h http.HandlerFunc) {
			mux.HandleFunc("POST /saas-warehouse"+path, authGuard(h))
		}
		reg("/busProduct/list", s.bpList)
		reg("/busProduct/buslist", s.bpBuslist)
		reg("/busProduct/warelist", s.bpWarelist)
		reg("/busProduct/hisPrice", s.bpHisPrice)
		reg("/busProduct/selCheckPro", s.bpSelCheckPro)
		reg("/busProduct/selNewPro", s.bpSelNewPro)
		reg("/busProduct/selCnPro", s.bpSelCnPro)
		reg("/busProduct/selNowNotice", s.bpSelNowNotice)
		reg("/busProduct/selBasCount", s.bpSelBasCount)
	})
}

// bpAttachTags 为 @Results 查询的行附加 tags(@Many selTags)
func bpAttachTags(d *sql.DB, rows []map[string]any) error {
	for i, row := range rows {
		tags, err := queryCam(d, "SELECT t.* from bus_tag t, bus_pro_tag pt where t.id = pt.sign_id and pt.pro_id = ?", row["id"])
		if err != nil {
			return err
		}
		rows[i]["tags"] = entBusTag.shape(tags)
	}
	return nil
}

// bpList = cuslist:allEq(maps,false)+effective 强制 1+update_time desc;
// three_pro_id "888"/"999"(toString 比较)整体替换 kjPros/authPros(无排序)
func (s *Service) bpList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	e := entBusProduct
	cond, ok := allEqCond(m, e.dbcols)
	if !ok {
		platform.Unknown(w)
		return
	}
	cond["effective"] = 1 // maps.put 强制覆盖
	q := "SELECT * FROM " + e.table
	var args []any
	if len(cond) > 0 {
		ks := platform.SortedKeys(cond)
		var conds []string
		for _, k := range ks {
			conds = append(conds, k+" = ?")
			args = append(args, cond[k])
		}
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	q += " ORDER BY update_time DESC"
	rows := s.qEntEnt2(w, r, e, q, args...)
	if rows == nil {
		return
	}
	if v, has := m["three_pro_id"]; has && v != nil {
		ts := str(m, "three_pro_id")
		if ts == "888" {
			rows = s.qEntEnt2(w, r, e, "SELECT * from bus_product where status = 1 and effective = 1 and (company_id=420 or group_id=420)")
		} else if ts == "999" {
			rows = s.qEntEnt2(w, r, e, "SELECT * from bus_product where status = 1 and effective = 1 and is_platform_certification=1")
		}
		if rows == nil {
			return
		}
	}
	platform.Success(w, rows)
}

// bpBuslist = buslist:GROUP BY product_name,仅 productName 键(实体键其余 null)
func (s *Service) bpBuslist(w http.ResponseWriter, r *http.Request) {
	out := s.qEntEnt2(w, r, entBusProduct, "SELECT product_name from bus_product GROUP BY product_name")
	if out == nil {
		return
	}
	for _, o := range out {
		if o == nil {
			continue
		}
		for _, k := range entBusProduct.keys {
			if k != "productName" {
				o[k] = nil
			}
		}
	}
	platform.Success(w, out)
}

// bpWarelist = warelist:getSql(t,"fstProId,bigProId,threeProId")+三 pro_id 裸拼;
// GROUP BY warehouse_nb(组内任取一行),无排序,tags 恒 null
func (s *Service) bpWarelist(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	e := entBusProduct
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	sqlStr, args, ok := e.getSqlUnstr(m, "fstProId,bigProId,threeProId")
	if !ok {
		platform.Unknown(w)
		return
	}
	if v, has := m["fstProId"]; has && v != nil {
		sqlStr += " and fst_pro_id = " + str(m, "fstProId")
	}
	if v, has := m["bigProId"]; has && v != nil {
		sqlStr += " and big_pro_id = " + str(m, "bigProId")
	}
	if v, has := m["threeProId"]; has && v != nil {
		sqlStr += " and three_pro_id = " + str(m, "threeProId")
	}
	q := "select (select product_name from bus_big_prodect where id=bus_product.fst_pro_id) fstName, (select product_name from bus_big_prodect where id=bus_product.big_pro_id) secName, (select product_name from bus_big_prodect where id=bus_product.three_pro_id) thName, (select warehouse_video from bus_warehouse where id=bus_product.warehouse_id) warehouseVideo, (select warehouse_img from bus_warehouse where id=bus_product.warehouse_id) warehouseImg, bus_product.* from bus_product " +
		sqlStr + " and status=1 GROUP BY warehouse_nb"
	rows := s.qEntEnt2(w, r, e, q, args...)
	if rows == nil {
		return
	}
	platform.Success(w, rows)
}

// bpSelCnPro = selCnPro:available_tonnage 非空+status=1,order by update_time desc,page 分页
func (s *Service) bpSelCnPro(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	q := "select * from bus_product where available_tonnage is not null and status=1 ORDER BY update_time desc"
	current, size := pageOf(m)
	rows := s.qEntEnt2(w, r, entBusProduct, q+" LIMIT ? OFFSET ?", size, (current-1)*size)
	if rows == nil {
		return
	}
	var total int64
	_ = s.pick(r).QueryRow("SELECT COUNT(*) FROM (select * from bus_product where available_tonnage is not null and status=1) cnt").Scan(&total)
	platform.Success(w, ipageOf(rows, total, current, size))
}

// bpSelNewPro = selNewPro:newPro(价未更新)/pricePro(价更新)各取一行,键恒在
func (s *Service) bpSelNewPro(w http.ResponseWriter, r *http.Request) {
	base := "select (select product_name from bus_big_prodect where id=bus_product.three_pro_id) thName, bus_product.* from bus_product where status=1 and "
	r1rows := s.qEntEnt2(w, r, entBusProduct, base+"(fprice_upd_time is null and pprice_upd_time is null) and TIMESTAMPDIFF(SECOND, last_check_time, update_time) < 10 ORDER BY update_time desc LIMIT 1")
	if r1rows == nil {
		return
	}
	r2rows := s.qEntEnt2(w, r, entBusProduct, base+"(fprice_upd_time is not null or pprice_upd_time is not null) and TIMESTAMPDIFF(SECOND, last_check_time, update_time) < 10 ORDER BY update_time desc LIMIT 1")
	if r2rows == nil {
		return
	}
	var newPro, pricePro any
	if len(r1rows) > 0 {
		newPro = r1rows[0]
	}
	if len(r2rows) > 0 {
		pricePro = r2rows[0]
	}
	platform.Success(w, map[string]any{"newPro": newPro, "pricePro": pricePro})
}

// bpSelNowNotice = selNowNotice:count(status=1) 默认文案;redis broadcastMsg 内存等价(空);
// 返回 String 数组(分号切分,尾随分号出空串)
func (s *Service) bpSelNowNotice(w http.ResponseWriter, r *http.Request) {
	var n int64
	if err := s.pick(r).QueryRow("select count(1) from bus_product where status=1").Scan(&n); err != nil {
		platform.Unknown(w)
		return
	}
	msg := "今日" + strconvI64(n) + "个《炼焦煤》直供平台"
	if broadcast := bpBroadcastText(); broadcast != "" {
		msg = broadcast
	}
	parts := strings.Split(msg, ";")
	platform.Success(w, parts)
}

// bpSelBasCount = selBasCount:三计数
func (s *Service) bpSelBasCount(w http.ResponseWriter, r *http.Request) {
	d := s.pick(r)
	var proCount, userCount, visitCount int64
	if err := d.QueryRow("select count(1) from bus_product where effective=1").Scan(&proCount); err != nil {
		platform.Unknown(w)
		return
	}
	if err := d.QueryRow("select count(1) from chain_auth_user where phone is not null and channel=2").Scan(&userCount); err != nil {
		platform.Unknown(w)
		return
	}
	if err := d.QueryRow("select count(1) from bus_pro_log where operate_type=7").Scan(&visitCount); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, map[string]any{
		"proCount": jsonNumberInt64(proCount), "userCount": jsonNumberInt64(userCount), "visitCount": jsonNumberInt64(visitCount),
	})
}

// bpHisPrice = hisPrice:分支按 priceType(拆箱 NPE→999);Week/Month SQL 违反
// ONLY_FULL_GROUP_BY 恒报错→999(MySQL8 实测);HYear 可执行;pro_id 取 body.id
func (s *Service) bpHisPrice(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProduct.validateBody(m) {
		platform.Unknown(w)
		return
	}
	hpt, has := m["hisPriceType"]
	if !has || hpt == nil {
		platform.Fail(w, 802, "参数异常")
		return
	}
	hbv, ok2 := bindVal("i", hpt)
	if !ok2 || hbv == nil {
		platform.Unknown(w)
		return
	}
	hisPriceType := hbv.(string)
	ptbv, ok3 := bindVal("i", m["priceType"])
	if !ok3 || ptbv == nil {
		platform.Unknown(w) // priceType 拆箱 NPE
		return
	}
	priceType := int(num(ptbv))
	if priceType != 1 && priceType != 2 && priceType != 3 {
		platform.Unknown(w) // prices 恒 null → NPE
		return
	}
	if priceType == 1 || priceType == 2 {
		platform.Unknown(w) // Week/Month SQL 触发 only_full_group_by → 999(实测)
		return
	}
	bp, _ := platform.GetByID(s.pick(r), "bus_product", m["id"])
	if bp == nil || bp["create_time"] == nil {
		platform.Unknown(w) // bp.getCreateTime() NPE
		return
	}
	proId := strOrNull(m, "id")
	if proId == "null" {
		proId = "null"
	}
	q := "select * from bus_price_log where pro_id=" + proId + " and price_type=" + hisPriceType + " and create_time between date_sub(DATE_ADD(curdate(),interval -day(curdate())+1 day),interval 5 month) and now() HAVING 1 ORDER BY create_time"
	rows, err := queryCam(s.pick(r), q)
	if err != nil {
		platform.Unknown(w)
		return
	}
	prices := entBusPriceLog.shape(rows)
	ldate := nowSpace()[:10]
	cal := ldate
	if ct, _ := bp["create_time"].(string); len(ct) >= 10 && ct[:10] > ldate {
		cal, ldate = ct[:10], ct[:10]
	}
	pMap := map[string]json.Number{}
	xlist := []string{}
	if len(prices) > 0 {
		fPrice, fPre := decStr(prices[0]["price"]), decStr(prices[0]["prePrice"])
		if fPrice == "" || fPre == "" {
			platform.Unknown(w) // compareTo NPE
			return
		}
		if fPrice != fPre && num(fPre) != 0 {
			pMap[ldate] = json.Number(fPre)
			xlist = append(xlist, ldate)
		}
		for _, pr := range prices {
			ct, _ := pr["createTime"].(string)
			if len(ct) < 10 {
				continue
			}
			time := ct[:10]
			if !containsStr(xlist, time) {
				xlist = append(xlist, time)
			}
			pMap[time] = toNum(pr["price"])
		}
	} else {
		price := ""
		if hisPriceType == "1" {
			price = decStr(bp["exFactoryPrice"])
		} else {
			price = decStr(bp["platePrice"])
		}
		if price != "" {
			xlist = append(xlist, cal)
			pMap[ldate] = json.Number(price)
		}
	}
	keys := make([]string, 0, len(pMap))
	for k := range pMap {
		keys = append(keys, k)
	}
	sort.Strings(keys) // TreeMap 字典序
	pd := map[string]any{}
	for _, k := range keys {
		pd[k] = pMap[k]
	}
	platform.Success(w, map[string]any{"priceData": pd, "xlist": xlist})
}

func containsStr(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// bpSelCheckPro = selCheckPro:getSql(Map,Class) 精确 LIKE + 19 槽 '%${maps.searchStr}%'
func (s *Service) bpSelCheckPro(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	e := entBusProduct
	// getSql(maps, BusProduct.class):键须有 @ExcelColumn,值精确绑定(不包 %)
	sqlStr := "where 1=1"
	var args []any
	for k, v := range m {
		if v == nil || k == "searchStr" {
			continue
		}
		if e.typeOf(k) == "x" || k == "page" || k == "ids" {
			continue
		}
		bv, ok2 := bindVal(e.typeOf(k), v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		if bv == nil {
			continue
		}
		sqlStr += " and " + camelOf(e, k) + " like ?"
		args = append(args, bv)
	}
	inner := "%null%"
	if v, has := m["searchStr"]; has && v != nil {
		inner = "%" + v.(string) + "%"
	}
	q := "select (select product_name from bus_big_prodect where id=bus_product.fst_pro_id) fstName, (select product_name from bus_big_prodect where id=bus_product.big_pro_id) secName, (select product_name from bus_big_prodect where id=bus_product.three_pro_id) thName, (select warehouse_video from bus_warehouse where id=bus_product.warehouse_id) warehouseVideo, (select warehouse_img from bus_warehouse where id=bus_product.warehouse_id) warehouseImg, bus_product.* from bus_product " +
		sqlStr + " and CONCAT(IFNULL(warehouse_nb, ''),IFNULL(region, ''),IFNULL(product_name, ''),IFNULL(resource_number, ''),IFNULL(send_addres, '')) like '%" + inner + "%'  order by update_time desc"
	rows := s.qEntEnt2(w, r, e, q, args...)
	if rows == nil {
		return
	}
	platform.Success(w, rows)
}

// decStr DECIMAL 值保字面(空/NULL → "")
func decStr(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return str(map[string]any{"k": v}, "k")
}

// toNum → json.Number(DECIMAL 保小数位)
func toNum(v any) json.Number {
	s := decStr(v)
	if s == "" {
		return json.Number("0")
	}
	return json.Number(s)
}
