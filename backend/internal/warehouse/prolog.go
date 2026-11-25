// prolog.go BusProLog bespoke 端点(批1;selMemberData 二期 authadmin 已注册不重复)。
// 全部计数 SQL 对照 BusProLogMapper/BusProductMapper/BusPhoneLogMapper/BusVisitorRecordMapper,
// ${} 注入点(wareIdentity/startTime/endTime/proId/operateObj 等)原样保留;
// redis userId{uid} hash 的 phone 在 Go 单体为空等价(无登录缓存)→ myPhoneCalledCount 按空串计数。
package warehouse

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		reg := func(path string, h http.HandlerFunc) {
			mux.HandleFunc("POST /saas-warehouse"+path, authGuard(h))
		}
		reg("/busProLog/add", s.plAdd)
		reg("/busProLog/selBuserData", s.plSelBuserData)
		reg("/busProLog/selBuserTeamData", s.plSelBuserTeamData)
		reg("/busProLog/selProShare", s.plSelProShare)
		reg("/busProLog/selProShareTimes", s.plSelProShareTimes)
		reg("/busProLog/selVisisted", s.plSelVisisted)
		reg("/busProLog/selCountByUser", s.plSelCountByUser)
		reg("/busProLog/selVisitedPop", s.plSelVisitedPop)
		reg("/busProLog/selVisitedData", s.plSelVisitedData)
		reg("/busProLog/selCallData", s.plSelCallData)
	})
}

func counts(w http.ResponseWriter, d *sql.DB, qs map[string]string) map[string]any {
	out := map[string]any{}
	for k, q := range qs {
		var s sql.NullString
		if err := d.QueryRow(q).Scan(&s); err != nil {
			platform.Unknown(w)
			return nil
		}
		// Java mapper 返回 int:DECIMAL 求和按 BigDecimal.intValue 截断,NULL → 0
		v := 0.0
		if s.Valid && s.String != "" {
			v = num(s.String)
		}
		out[k] = jsonNumberInt64(int64(v))
	}
	return out
}

// pagedRows 同 paged 但不驼峰化(MyBatis Map 查询保留列标签原样)
func pagedRows(d *sql.DB, q, countQ string, args []any, current, size int64) (map[string]any, error) {
	var total int64
	if err := d.QueryRow(countQ, args...).Scan(&total); err != nil {
		return nil, err
	}
	off := (current - 1) * size
	rows, err := queryRows(d, q+" LIMIT ? OFFSET ?", append(append([]any{}, args...), size, off)...)
	if err != nil {
		return nil, err
	}
	return platform.IPage(rows, total, current, size), nil
}

// plBody 解析 body 并校验(载体 startTime/endTime/selType 为 LocalDate/Integer)
func plBody(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	m := bodyMap(r)
	if !entBusProLog.validateBody(m) {
		platform.Unknown(w)
		return nil, false
	}
	return m, true
}

// plSelBuserData = selBuserData:我发布/我管理/我操作 三组计数;redis 空→phone="" 
func (s *Service) plSelBuserData(w http.ResponseWriter, r *http.Request) {
	m, ok := plBody(w, r)
	if !ok {
		return
	}
	uid := str(m, "userId")
	if uid == "" {
		uid = strconvI64(currentUser(r).UserID)
	}
	d := s.pick(r)
	pub := counts(w, d, map[string]string{
		"myProCount":        "select count(1) from bus_product where create_user_id=" + uid + " and effective=1",
		"myProVisitedCount": "select count(1) from bus_pro_log where target_user_id=" + uid + " and operate_type=7 and operate_obj=4",
		"myProForwardCount": "select count(1) from bus_pro_log where target_user_id=" + uid + " and operate_type=5 and operate_obj=4",
		"myProShareCount":   "select count(1) from bus_pro_log where operate_type=4 and operate_obj=4 and target_user_id=" + uid,
		"myPhoneCalledCount": "select count(1) from bus_phone_log where target_user_phone=''", // redis 空 → phone=""
		"myCardVisitedCount": "select count(1) from bus_pro_log where target_user_id=" + uid + " and operate_type=7 and operate_obj=2",
		"myProCollectCount": "select count(1) from bus_pro_collection where target_user_id=" + uid,
	})
	if pub == nil {
		return
	}
	pop, err := queryRows(d, "select user_id, count(1) from bus_pro_log where target_user_id="+uid+" and operate_type=7 and operate_obj=4 group by user_id")
	if err != nil {
		platform.Unknown(w)
		return
	}
	pub["myProPopVisitedCount"] = jsonNumberInt64(int64(len(pop)))
	mng := counts(w, d, map[string]string{
		"myEtcCount":        "select IFNULL(sum(c.count), 0) from ( SELECT count(1) count from bus_product where business_id=" + uid + " or area_buser_id=" + uid + " GROUP BY company_id ) c",
		"myCustProCount":    "SELECT IFNULL(count(1), 0) count from bus_product where business_id=" + uid + " or area_buser_id=" + uid,
		"myCustCjCount":     "SELECT IFNULL(sum(sum_tonnage_traded), 0) count from bus_product where business_id=" + uid + " or area_buser_id=" + uid,
		"myCustForwardCount": "SELECT IFNULL(sum(forward_times), 0) count from bus_product where business_id=" + uid + " or area_buser_id=" + uid,
	})
	if mng == nil {
		return
	}
	ope := counts(w, d, map[string]string{
		"myIndexShareCount": "select count(1) from bus_pro_log where user_id=" + uid + " and operate_type=4 and operate_obj=1",
		"myCardShareCount":  "select count(1) from bus_pro_log where user_id=" + uid + " and operate_type=4 and operate_obj=2",
		"myVideoShareCount": "select count(1) from bus_pro_log where user_id=" + uid + " and operate_type=4 and operate_obj=3",
		"myCoalShareCount":  "select count(1) from bus_pro_log where user_id=" + uid + " and operate_type=4 and operate_obj=3", // Java 同 myVideoShareCount(复用)
		"myForwardCount":    "select count(1) from bus_pro_log where user_id=" + uid + " and operate_type=5",
		"myShareCount":      "select count(1) from bus_pro_log where user_id=" + uid + " and operate_type=4",
	})
	if ope == nil {
		return
	}
	platform.Success(w, map[string]any{"myPublish": pub, "myManage": mng, "myOperate": ope})
}

// plSelBuserTeamData = selBuserTeamData:where/proWhere 字面注入(Java endTime 分支
// 复制的是 >= startTime——bug 照抄);brokerList 恒空 → brokerObj {}
func (s *Service) plSelBuserTeamData(w http.ResponseWriter, r *http.Request) {
	m, ok := plBody(w, r)
	if !ok {
		return
	}
	where, proWhere := "", ""
	if v, has := m["wareIdentity"]; has && v != nil {
		where += " and ware_identity = " + str(m, "wareIdentity")
	}
	// 空串日期:Java 绑定得 null 后跳过该条件;Go 的 bindVal 返回 (nil, true),不判空会对 nil 断言 panic
	if v, has := m["startTime"]; has && v != nil {
		if st, okD := bindVal("D", v); okD && st != nil {
			proWhere += " and DATE_FORMAT(create_time, '%Y-%m-%d')>= " + st.(string)
			where += " and DATE_FORMAT(create_time, '%Y-%m-%d')>= " + st.(string)
		}
	}
	if v, has := m["endTime"]; has && v != nil {
		if et, okD := bindVal("D", v); okD && et != nil {
			proWhere += " and DATE_FORMAT(create_time, '%Y-%m-%d')<= " + et.(string)
			where += " and DATE_FORMAT(create_time, '%Y-%m-%d')>= " + strOr(m, "startTime", "") // 原 bug:用 startTime
		}
	}
	d := s.pick(r)
	res := counts(w, d, map[string]string{
		"sumVisitedCount":    "SELECT IFNULL(count(1), 0) from bus_visitor_record where 1=1 " + proWhere,
		"sumProCount":        "SELECT IFNULL(count(1), 0) from bus_product where 1=1 " + proWhere,
		"sumForwardCount":    "select count(1) from bus_pro_log where operate_type=5 " + where,
		"sumCoalForwardCount": "select count(1) from bus_pro_log where operate_type=5 and operate_obj=4 " + where,
		"sumCoalShareCount":  "select count(1) from bus_pro_log where operate_type=4 and operate_obj=4 " + where,
		"sumIdxForwardCount": "select count(1) from bus_pro_log where operate_type=5 and operate_obj=1 " + where,
		"sumCardForwardCount": "select count(1) from bus_pro_log where operate_type=5 and operate_obj=2 " + where,
	})
	if res == nil {
		return
	}
	res["brokerObj"] = map[string]any{}
	platform.Success(w, res)
}

// plSelProShare = selProShare:share SQL 忽略筛选条件,仅 ${maps.proId} 生效;
// 无 page → List;有 page → IPage
func (s *Service) plSelProShare(w http.ResponseWriter, r *http.Request) {
	m, ok := plBody(w, r)
	if !ok {
		return
	}
	// ${maps.proId}:getMaps 无该键 → MyBatis 缺键异常 999(实测)
	bv, ok2 := bindVal("i", m["proId"])
	if !ok2 || bv == nil {
		platform.Unknown(w) // ${maps.proId} 缺键/类型错 → MyBatis 异常 999(实测)
		return
	}
	proId := bv.(string)
	q := "SELECT (SELECT wx_face_img from chain_auth_user where phone=user_phone and channel=2) wxFaceImg,user_phone userPhone, count(1) shareTimes from bus_pro_log where operate_obj=4 and operate_type=4 and pro_id=" + proId + " GROUP BY user_phone"
	_, hasPage := m["page"].(map[string]any)
	if !hasPage {
		rows, err := queryRows(s.pick(r), q)
		if err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, rows)
		return
	}
	current, size := pageOf(m)
	countQ := "SELECT COUNT(*) FROM (SELECT user_phone FROM bus_pro_log where operate_obj=4 and operate_type=4 and pro_id=" + proId + " GROUP BY user_phone) cnt"
	page, err := paged(s.pick(r), q, countQ, nil, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, page)
}

// plSelProShareTimes = selProShareTimes:allEq 三键计数(proId 缺失 → pro_id IS NULL)
func (s *Service) plSelProShareTimes(w http.ResponseWriter, r *http.Request) {
	m, ok := plBody(w, r)
	if !ok {
		return
	}
	q := "select count(1) from bus_pro_log where operate_obj=4 and operate_type=4"
	if v, has := m["proId"]; has && v != nil {
		bv, ok2 := bindVal("i", v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		if bv != nil {
			q += " and pro_id=" + bv.(string)
		} else {
			q += " and pro_id IS NULL"
		}
	} else {
		q += " and pro_id IS NULL"
	}
	var n int64
	if err := s.pick(r).QueryRow(q).Scan(&n); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, jsonNumberInt64(n))
}

// plSelVisisted = selVisisted:getSql + DATE_FORMAT(pl.create_time) 无引号注入;
// 逐行 phoneLogs(body userId;缺失 → call_user_id=null 空行)
func (s *Service) plSelVisisted(w http.ResponseWriter, r *http.Request) {
	m, ok := plBody(w, r)
	if !ok {
		return
	}
	e := entBusProLog
	sqlStr, args, ok2 := e.getSql(m)
	if !ok2 {
		platform.Unknown(w)
		return
	}
	if v, has := m["startTime"]; has && v != nil { // 空串 → nil,跳过(同上)
		if st, okD := bindVal("D", v); okD && st != nil {
			sqlStr += " and DATE_FORMAT(pl.create_time, '%Y-%m-%d')>= " + st.(string)
		}
	}
	if v, has := m["endTime"]; has && v != nil {
		if et, okD := bindVal("D", v); okD && et != nil {
			sqlStr += " and DATE_FORMAT(pl.create_time, '%Y-%m-%d')<= " + et.(string)
		}
	}
	uid := str(m, "userId")
	if uid == "" {
		uid = "null"
	}
	q := "select enterprise_name enterpriseName, warehouse_nb warehouseNb, max(pl.create_time) createTime, pl.pro_id proId from bus_pro_log pl, bus_product p " +
		sqlStr + " and operate_type=7 and operate_obj=4 and pl.pro_id=p.id GROUP BY pl.user_id, pl.pro_id, enterprise_name, warehouse_nb ORDER BY pl.create_time desc"
	visits, err := queryRows(s.pick(r), q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	for _, row := range visits {
		phoneLogs, err := queryRows(s.pick(r),
			"select target_user_phone phone, target_user_name username, count(1) callCount from bus_phone_log where call_user_id="+uid+" and pro_id="+strOf(row["proId"])+" group by call_user_id, pro_id, target_user_phone, target_user_name")
		if err != nil {
			platform.Unknown(w)
			return
		}
		row["phoneLogs"] = phoneLogs
	}
	platform.Success(w, visits)
}

// plSelCountByUser = selCountByUser:count(DISTINCT user_id);缺参 → =null 恒 0
func (s *Service) plSelCountByUser(w http.ResponseWriter, r *http.Request) {
	m, ok := plBody(w, r)
	if !ok {
		return
	}
	// ${param} null → MyBatis 异常 999(实测任一缺参即 999)
	obj := strOrNull(m, "operateObj")
	typ := strOrNull(m, "operateType")
	if obj == "null" || typ == "null" {
		platform.Unknown(w)
		return
	}
	q := "select IFNULL(count(DISTINCT user_id), 0) from bus_pro_log where operate_obj=" + obj + " and operate_type =" + typ
	var n int64
	if err := s.pick(r).QueryRow(q).Scan(&n); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, jsonNumberInt64(n))
}

// plSelVisitedPop = selVisitedPop:getSql + group by user_phone 聚合页
func (s *Service) plSelVisitedPop(w http.ResponseWriter, r *http.Request) {
	m, ok := plBody(w, r)
	if !ok {
		return
	}
	e := entBusProLog
	sqlStr, args, ok2 := e.getSql(m)
	if !ok2 {
		platform.Unknown(w)
		return
	}
	var search any
	if v, has := m["searchStr"]; has && v != nil {
		ss, ok3 := bindVal("s", v)
		if !ok3 {
			platform.Unknown(w)
			return
		}
		search = "%" + ss.(string) + "%"
	}
	q := "select (select wx_face_img from chain_auth_user where id=user_id) visitorHeadimg, user_phone userPhone, max(create_time) createTime from bus_pro_log " +
		sqlStr + " and CONCAT(id) like ? group by user_phone order by max(create_time) desc"
	countQ := "SELECT COUNT(*) FROM (SELECT user_phone FROM bus_pro_log " + sqlStr + " and CONCAT(id) like ? group by user_phone) cnt"
	current, size := pageOf(m)
	page, err := paged(s.pick(r), q, countQ, append(args, search), current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, page)
}

// plSelVisitedData = selVisitedData:selType 1/2/3 三聚合页;selType 缺失 → 拆箱 NPE 999
func (s *Service) plSelVisitedData(w http.ResponseWriter, r *http.Request) {
	m, ok := plBody(w, r)
	if !ok {
		return
	}
	st, hasSt := m["startTime"], false
	et := any(nil)
	if v, has := m["startTime"]; has && v != nil {
		st, _ = bindVal("D", v)
		hasSt = st != nil
	}
	if v, has := m["endTime"]; has && v != nil {
		et, _ = bindVal("D", v)
	}
	whereSql := ""
	if hasSt {
		whereSql += " and DATE_FORMAT(pl.create_time, '%Y-%m-%d')>= '" + st.(string) + "'"
	}
	if ev, has := m["endTime"]; has && ev != nil && et != nil {
		whereSql += " and DATE_FORMAT(pl.create_time, '%Y-%m-%d')<= '" + et.(string) + "'"
	}
	selType := 0
	if v, has := m["selType"]; has && v != nil {
		bv, ok2 := bindVal("i", v)
		if !ok2 || bv == nil {
			platform.Unknown(w)
			return
		}
		selType = int(num(bv))
	} else {
		platform.Unknown(w) // 1==null 拆箱 NPE
		return
	}
	var q, countQ string
	switch selType {
	case 1:
		q = "SELECT pl.user_phone, count(1) visitCount, (SELECT count(1) clickCount from bus_pro_log where operate_obj=8 and operate_type=10 and user_phone=pl.user_phone) clickCount, DATE_FORMAT(max(pl.create_time), '%Y-%m-%d %H:%i:%S') visitTime from bus_pro_log pl WHERE pl.operate_obj=6 and pl.operate_type=7 " + whereSql + " GROUP BY pl.user_phone order by visitTime desc"
		countQ = "SELECT COUNT(*) FROM (SELECT user_phone FROM bus_pro_log pl WHERE pl.operate_obj=6 and pl.operate_type=7 " + whereSql + " GROUP BY pl.user_phone) cnt"
	case 2:
		q = "SELECT pl.user_phone, p.enterprise_name, p.id proId, p.warehouse_nb, count(1) visitCount, (SELECT count(1) from bus_phone_log WHERE pro_id=pl.pro_id and call_user_phone=pl.user_phone) callCount, DATE_FORMAT(max(pl.create_time), '%Y-%m-%d %H:%i:%S') visitTime from bus_pro_log pl, bus_product p WHERE pl.operate_obj=4 and pl.operate_type=7 and pl.user_phone is not null and user_phone!='' and pl.pro_id=p.id " + whereSql + " GROUP BY pl.user_phone, pro_id order by visitTime desc"
		countQ = "SELECT COUNT(*) FROM (SELECT pl.user_phone FROM bus_pro_log pl, bus_product p WHERE pl.operate_obj=4 and pl.operate_type=7 and pl.user_phone is not null and user_phone!='' and pl.pro_id=p.id " + whereSql + " GROUP BY pl.user_phone, pro_id) cnt"
	case 3:
		q = "SELECT pl.user_phone, count(1) visitCount, (SELECT count(1) from bus_phone_log WHERE lease_id is not null and call_user_phone=pl.user_phone) callCount, DATE_FORMAT(max(pl.create_time), '%Y-%m-%d %H:%i:%S') visitTime from bus_pro_log pl WHERE pl.operate_obj=5 and pl.operate_type=7 and pl.user_phone is not null and user_phone!='' " + whereSql + " GROUP BY pl.user_phone order by visitTime desc"
		countQ = "SELECT COUNT(*) FROM (SELECT pl.user_phone FROM bus_pro_log pl WHERE pl.operate_obj=5 and pl.operate_type=7 and pl.user_phone is not null and user_phone!='' " + whereSql + " GROUP BY pl.user_phone) cnt"
	default:
		platform.Success(w, nil) // result=null
		return
	}
	current, size := pageOf(m)
	page, err := pagedRows(s.pick(r), q, countQ, nil, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, page)
}

// plSelCallData = selCallData:按拨打方手机号聚合目标(isService/isPtjjr 硬编码号码)
func (s *Service) plSelCallData(w http.ResponseWriter, r *http.Request) {
	m, ok := plBody(w, r)
	if !ok {
		return
	}
	phone := "null"
	if v, has := m["userPhone"]; has && v != nil {
		bv, ok2 := bindVal("s", v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		if bv != nil {
			phone = "'" + sqlEscape(bv.(string)) + "'"
		}
	}
	q := "SELECT target_user_phone, target_user_name, (target_user_phone= '13934568608' or target_user_phone='15536587788') isService, (target_user_phone= '18636842628' or target_user_phone='19135022011' or target_user_phone='13383516668') isPtjjr, count(1) callCount from bus_phone_log where call_user_phone=" + phone + " and pro_id =" + strOrNull(m, "proId") + " GROUP BY target_user_phone"
	rows, err := queryRows(s.pick(r), q)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, rows)
}

// ── 小工具 ──

func strconvI64(n int64) string {
	return strconv.FormatInt(n, 10)
}

func strOf(v any) string {
	if v == nil {
		return "null"
	}
	return str(map[string]any{"k": v}, "k")
}

func strOrNull(m map[string]any, k string) string {
	v, has := m[k]
	if !has || v == nil {
		return "null"
	}
	bv, ok := bindVal("i", v)
	if !ok || bv == nil {
		return "null"
	}
	return bv.(string)
}

func strOr(m map[string]any, k, def string) string {
	if v, has := m[k]; has && v != nil {
		return str(m, k)
	}
	return def
}

// plAdd = saveBusProLog(操作日志 + 分享/转发计数与积分):
// 落 bus_pro_log;operate_obj=4 且 pro_id 非空 → sharingTimes/forwardTimes+1、
// 本人分享 → updUserScore(4,1);响应 score
func (s *Service) plAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProLog.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	vals := pick(m, entBusProLog.dbcols)
	delete(vals, "id")
	if uid := strconvI64(currentUser(r).UserID); uid != "" {
		if _, has := vals["user_id"]; !has {
			vals["user_id"] = uid
		}
	}
	if _, err := platform.Insert(d, "bus_pro_log", vals); err != nil {
		platform.Unknown(w)
		return
	}
	score := 0
	if str(m, "operateObj") == "4" {
		proId := str(m, "proId")
		if proId != "" {
			bp, err := platform.GetByID(d, "bus_product", proId)
			if err != nil || bp == nil {
				platform.Unknown(w) // getById(proId) null → NPE
				return
			}
			ot := str(m, "operateType")
			if ot == "4" {
				st := int(num(bp["sharing_times"])) + 1
				_, _ = d.Exec("UPDATE bus_product SET sharing_times=?, last_share_time=? WHERE id=?", strconv.Itoa(st), nowSpace(), proId)
			} else if ot == "5" {
				ft := int(num(bp["forward_times"])) + 1
				_, _ = d.Exec("UPDATE bus_product SET forward_times=?, last_share_time=? WHERE id=?", strconv.Itoa(ft), nowSpace(), proId)
			}
			if ot == "4" || ot == "5" {
				if uid := str(m, "userId"); uid != "" && num(uid) == num(bp["create_user_id"]) {
					// updEtcScore(proId, companyId) 单产品分享分(MQ 消费侧等价——登记差异,
					// 这里只做直接可见的 updUserScore(4,1) 分享积分)
					sc := s.bpUpdUserScore(w, d, bp, 4, 1)
					if sc < 0 {
						return
					}
					score = sc
				}
			}
			if _, err := d.Exec("UPDATE bus_product SET sharing_times=IFNULL(sharing_times,0) WHERE id=?", proId); err != nil {
				platform.Unknown(w)
				return
			}
		}
	}
	platform.Success(w, jsonNumberInt64(int64(score)))
}
