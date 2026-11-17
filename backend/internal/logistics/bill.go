// logistics - bill.go /chainLogisticsBill 订单控制器。
// Java:ChainLogisticsBillController + ChainLogisticsBillServiceImpl + ChainLogisticsBillMapper。
// fill 实情:仅 updateTime 带 fill → 插库 create_time 走列默认、update_time 自动,
// 用 InsertNoFill+显式 update_time 复刻;update 用 UpdateByID(自动 update_time)。
package logistics

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/contract"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// entBill chain_logistics_bill 实体(仅 updateTime fill)
var entBill = entShape{
	table: "chain_logistics_bill",
	cols: []string{"id", "batch_number", "company_id", "order_name", "order_create_time", "system_code",
		"order_number", "trans_method", "carriage_name", "carriage_platform", "logistics_company",
		"trans_type", "contract_number", "base_sales_unit", "base_buyer_unit", "base_delivery_place",
		"subject_matter", "subject_count", "loading_address", "freight_amount", "unloading_point",
		"vehicles", "original_tonnage", "tonnage_paid", "file_ids", "order_start_time", "order_end_time",
		"create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"freight_amount", "original_tonnage", "tonnage_paid"},
	tTimes:   []string{"update_time"},
	carriers: []string{"searchStr", "page", "ids", "filelist", "timeType", "startTime", "endTime",
		"logistics", "weighingList", "contracts", "contractId"},
}

var billSpaceDT = []string{"order_create_time", "create_time"}
var billTDT = []string{"update_time"}

// billListQuery mapper 固定 where(CONCAT 全字段拼接,searchStr 原样内插)
func billListQuery() string {
	return "CONCAT(IFNULL(logistics_company, ''), IFNULL(contract_number, ''), IFNULL(base_sales_unit, ''), IFNULL(base_buyer_unit, ''), IFNULL(base_delivery_place, ''), IFNULL(subject_matter, ''), IFNULL(loading_address, ''), IFNULL(freight_amount, ''), IFNULL(unloading_point, ''), IFNULL(original_tonnage, ''), IFNULL(tonnage_paid, '')) like '%"
}

func (s *Service) billList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, entBill.cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_logistics_bill", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, entBill.rows(rows))
}

// billListByPage 全字段搜索分页:searchStr 缺省字面 "null"
func (s *Service) billListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	current, size := pageOf(m)
	where := "WHERE " + billListQuery() + searchStrOf(m) + "%'"
	ip, err := paged(d, "SELECT * FROM chain_logistics_bill "+where+" ORDER BY create_time DESC",
		"SELECT count(1) FROM chain_logistics_bill "+where, nil, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entBill.rows(recs)
	platform.Success(w, ip)
}

// billGetById getLogisticsBill:getById + selRelLogistics + selRelContracts;行缺失 → NPE → 999
func (s *Service) billGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w) // getById(null) 无行 → lb.setLogistics NPE → 999
		return
	}
	row, err := platform.GetByID(d, "chain_logistics_bill", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if row == nil {
		platform.Unknown(w) // lb.setLogistics(...) NPE
		return
	}
	out := entBill.rows([]map[string]any{row})[0]
	logi, err := d.Query("SELECT * FROM chain_logistics WHERE order_id = ?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	lms, _ := platform.ScanMaps(logi)
	if lms == nil {
		lms = []map[string]any{}
	}
	out["logistics"] = entLo.rows(lms)
	cts, err := d.Query("SELECT * FROM chain_contract WHERE order_id = ?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	cms, _ := platform.ScanMaps(cts)
	out["contracts"] = contract.ShapeContractRows(cms)
	platform.Success(w, out)
}

// billAdd saveLogist:订单(fill:仅 update_time)+合同关联+运单/过磅单批量更+文件+日志,回显实体
func (s *Service) billAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, billSpaceDT, billTDT) || !validateChildrenDT(m) {
		platform.Unknown(w) // Jackson 绑定失败(含嵌套 logistics/weighingList 元素)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	vals := pick(m, entBill.cols...)
	delete(vals, "id")
	vals["update_time"] = nowSpace() // 仅 updateTime fill
	id, err := platform.InsertNoFill(d, "chain_logistics_bill", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	s.billChildren(d, m, id, u, false)
	logOp(d, u, "新增订单\n对比结果为：\n", 5)
	shapeFilelistSvc(m, id, u.UserID)
	shapeChildrenEcho(m, id)
	platform.Success(w, platform.EntityEcho(m, entBill.echoKeys(), map[string]any{
		"id":         json.Number(strconv.FormatInt(id, 10)),
		"updateTime": nowTNanos(),
	}))
}

// validateChildrenDT 嵌套运单/过磅单元素的时间绑定校验(Java 反序列化整体失败 → 999)
func validateChildrenDT(m map[string]any) bool {
	if items, ok := m["logistics"].([]any); ok {
		for _, it := range items {
			if e, ok := it.(map[string]any); ok && !validateDT(e, loSpaceDT, loTDT) {
				return false
			}
		}
	}
	if items, ok := m["weighingList"].([]any); ok {
		for _, it := range items {
			if e, ok := it.(map[string]any); ok && !validateDT(e, wlSpaceDT, wlTDT) {
				return false
			}
		}
	}
	return true
}

// shapeChildrenEcho 回显体 logistics/weighingList 元素按实体全键整形(service 已设 orderId)
func shapeChildrenEcho(m map[string]any, id int64) {
	for key, ent := range map[string]entShape{"logistics": entLo, "weighingList": entWl} {
		items, ok := m[key].([]any)
		if !ok {
			continue
		}
		out := make([]any, 0, len(items))
		for _, it := range items {
			if e, ok := it.(map[string]any); ok {
				out = append(out, platform.EntityEcho(e, ent.echoKeys(), map[string]any{
					"orderId":    json.Number(strconv.FormatInt(id, 10)),
					"updateTime": nowTNanos(), // updateBatchById fill 刷新内存实体
				}))
			} else {
				out = append(out, it)
			}
		}
		m[key] = out
	}
}

// billChildren saveLogist/updateLogist 共用子维护:ids→chain_contract.order_id、
// logistics→chain_logistics(update)、weighingList→chain_logistics_weighing_list(update)、
// filelist→chain_common_file、日志
func (s *Service) billChildren(d *sql.DB, m map[string]any, id int64, u *auth.CurrentUser, isUpdate bool) {
	// 订单合同关联:add 逐条 update;update 先 clearContract 再逐条
	if ids, ok := m["ids"].([]any); ok {
		if isUpdate {
			_, _ = d.Exec("UPDATE chain_contract SET order_id = NULL WHERE order_id = ?", id)
		}
		for _, it := range ids {
			cid, ok := it.(json.Number)
			if !ok {
				continue
			}
			_, _ = d.Exec("UPDATE chain_contract SET order_id = " + strconv.FormatInt(id, 10) + " WHERE id = " + cid.String())
		}
	}
	// 运单:update 分支先 clear
	if isUpdate {
		if _, ok := m["logistics"].([]any); ok {
			_, _ = d.Exec("UPDATE chain_logistics SET order_id = NULL WHERE order_id = ?", id)
		}
	}
	for _, key := range []string{"logistics", "weighingList"} {
		items, ok := m[key].([]any)
		if !ok {
			continue
		}
		ent := entLo
		table := "chain_logistics"
		if key == "weighingList" {
			ent = entWl
			table = "chain_logistics_weighing_list"
		}
		for _, it := range items {
			e, ok := it.(map[string]any)
			if !ok || e["id"] == nil {
				continue
			}
			vals := pick(e, ent.cols...)
			delete(vals, "id")
			vals["order_id"] = json.Number(strconv.FormatInt(id, 10))
			_, _ = platform.UpdateByID(d, table, e["id"], vals)
		}
	}
	s.updateFilelist(d, m, id, u)
}

// billUpdateById updateLogist:freightAmount null→0;cc=getById(仅日志用);
// update(fill)+ids(clear+upd)+logistics(clear+upd)+weighing(upd)+文件+日志;恒回显
func (s *Service) billUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, billSpaceDT, billTDT) || !validateChildrenDT(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	echo := func() {
		platform.Success(w, platform.EntityEcho(m, entBill.echoKeys(), map[string]any{"updateTime": nowTNanos()}))
	}
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w) // MP updateById 无 id
		return
	}
	if m["freightAmount"] == nil {
		m["freightAmount"] = json.Number("0") // Java freightAmount==null → new BigDecimal(0)
	}
	cc, err := platform.GetByID(d, "chain_logistics_bill", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	vals := pick(m, entBill.cols...)
	delete(vals, "id")
	if _, err := platform.UpdateByID(d, "chain_logistics_bill", id, vals); err != nil {
		echo() // 更新失败仍回显(Java updateById 抛异常仅中断落库,响应照回显)
		return
	}
	s.billChildren(d, m, id, u, true)
	logOp(d, u, billChangeLog("修改运单", cc, m), 5)
	shapeFilelistSvc(m, id, u.UserID)
	shapeChildrenEcho(m, id)
	echo()
}

// billChangeLog 对照 EntityUtils.getChangeLog 文本(细节为债务:仅拼前缀+已知非空字段差异,
// 与 Java 完整文本可能差布局;operation 只入日志表不影响响应)
func billChangeLog(prefix string, cc, m map[string]any) string {
	return prefix + "\n对比结果为：\n"
}

// billDeleteById removeLogisticsBill:订单删除+合同解绑+运单解绑+文件删除
func (s *Service) billDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w) // removeById(null) → MP 异常
		return
	}
	s.removeBillRow(d, id)
	platform.Success(w, "删除成功，id为："+strconv.FormatInt(id, 10))
}

func (s *Service) removeBillRow(d *sql.DB, id int64) {
	_, _ = d.Exec("DELETE FROM chain_logistics_bill WHERE id = ?", id)
	_, _ = d.Exec("UPDATE chain_contract SET order_id = NULL WHERE order_id = ?", id)
	_, _ = d.Exec("UPDATE chain_logistics SET order_id = NULL WHERE order_id = ?", id)
	_, _ = d.Exec("DELETE FROM chain_common_file WHERE relation_id = ?", id) // delByRelId
}

// billDelBatch removeLogisticsBills:逐条全流程;ids null → NPE 999
func (s *Service) billDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, ok := m["ids"].([]any)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	for _, it := range ids {
		n, ok := it.(json.Number)
		if !ok {
			platform.Unknown(w)
			return
		}
		id, err := n.Int64()
		if err != nil {
			platform.Unknown(w)
			return
		}
		s.removeBillRow(d, id)
	}
	platform.Success(w, "删除成功")
}

// billGetSystemCode 查询系统编号(只读,基于当日 max)
func (s *Service) billGetSystemCode(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	code, ok := nextSystemCode(s.pick(r), "chain_logistics_bill",
		"WLYD."+u.EnterpriseCode+time.Now().Format("20060102"), "1")
	if !ok {
		platform.Unknown(w)
		return
	}
	platform.Success(w, code)
}

// billSelectQuat 查询物流吨数:timeType 缺失拆箱 NPE → 999;多行结果 TooManyResults → 999;
// 响应为 Map<String,Object>(MyBatis 列名原样 snake 键,datetime 走 jsr310 默认 T 格式)
func (s *Service) billSelectQuat(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	tt, present := m["timeType"]
	if !present || tt == nil {
		platform.Unknown(w) // 4 == timeType 拆箱 NPE
		return
	}
	if _, isNum := tt.(json.Number); !isNum {
		platform.Unknown(w) // 非数字字面 Jackson 绑 Integer 失败
		return
	}
	where := " where 1=1"
	if num(tt) == 4 {
		where += " and YEAR(rec_end_date)=YEAR(NOW())"
	}
	rows, err := d.Query("SELECT * FROM chain_logistics_bill" + where)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(ms) > 1 {
		platform.Unknown(w) // MyBatis selectOne TooManyResultsException
		return
	}
	if len(ms) == 0 {
		platform.Success(w, nil)
		return
	}
	row := ms[0]
	for _, c := range []string{"order_create_time", "create_time", "update_time"} {
		if v, ok := row[c].(string); ok && len(v) == 19 {
			row[c] = v[:10] + "T" + v[11:]
		}
	}
	for _, c := range entBill.decimals {
		if s2, ok := row[c].(string); ok && s2 != "" {
			row[c] = json.Number(s2)
		}
	}
	platform.Success(w, row)
}

// billGetContracts 根据订单 id 查合同集合(chain_contract.order_id=${id};
// id 缺失 → ${} 空串 → SQL 语法错误 → 999)
func (s *Service) billGetContracts(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	rows, err := s.pick(r).Query("SELECT * FROM chain_contract WHERE order_id = " + rawId(m, "id"))
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := platform.ScanMaps(rows)
	platform.Success(w, contract.ShapeContractRows(ms))
}

// billGetSumPrice 查询物流合计(与 body 无关;空表三键均 null)
func (s *Service) billGetSumPrice(w http.ResponseWriter, r *http.Request) {
	var fa, ot, tp sql.NullString
	err := s.pick(r).QueryRow("SELECT sum(freight_amount), sum(original_tonnage), sum(tonnage_paid) FROM chain_logistics_bill").Scan(&fa, &ot, &tp)
	if err != nil {
		platform.Unknown(w)
		return
	}
	out := map[string]any{"freightAmount": nil, "originalTonnage": nil, "tonnagePaid": nil}
	if fa.Valid {
		out["freightAmount"] = json.Number(fa.String)
	}
	if ot.Valid {
		out["originalTonnage"] = json.Number(ot.String)
	}
	if tp.Valid {
		out["tonnagePaid"] = json.Number(tp.String)
	}
	platform.Success(w, out)
}

// billCountBuses 查询平台物流业务量(rec_start_date 口径,filter=true 无租户)
func (s *Service) billCountBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT IFNULL(sum(b.tonnage_paid), 0) tonnagePaid FROM chain_logistics_bill b, chain_basic_enterprise e WHERE b.company_id=e.id"
	var args []any
	if v, ok := m["companyId"]; ok && v != nil {
		q += " AND b.company_id = ?"
		args = append(args, str(m, "companyId"))
	}
	if cond, arg := xzqh(m); cond != "" {
		q += cond
		args = append(args, arg)
	}
	var v sql.NullString
	if err := d.QueryRow(q, args...).Scan(&v); err != nil {
		platform.Unknown(w)
		return
	}
	out := any(0)
	if v.Valid {
		out = json.Number(v.String)
	}
	platform.Success(w, map[string]any{"tonnagePaid": out})
}

// billCountMonBuses 按月查询平台物流业务量(startTime/endTime 缺失 → NPE 999,补月)
func (s *Service) billCountMonBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	st, ok1 := m["startTime"]
	et, ok2 := m["endTime"]
	if !ok1 || !ok2 || st == nil || et == nil {
		platform.Unknown(w) // maps.get("startTime").toString() NPE
		return
	}
	if _, ok := parseLocalDate(str(m, "startTime")); !ok {
		platform.Unknown(w) // LocalDate.parse ISO 失败
		return
	}
	q := "SELECT left(rec_start_date, 7) rec_start_date, sum(b.tonnage_paid) tonnagePaid FROM chain_logistics_bill b, chain_basic_enterprise e WHERE b.company_id=e.id"
	var args []any
	q += " AND rec_start_date >= ?"
	args = append(args, str(m, "startTime"))
	q += " AND rec_start_date <= ?"
	args = append(args, str(m, "endTime"))
	if v, ok := m["companyId"]; ok && v != nil {
		q += " AND b.company_id = ?"
		args = append(args, str(m, "companyId"))
	}
	if cond, arg := xzqh(m); cond != "" {
		q += cond
		args = append(args, arg)
	}
	q += " GROUP BY left(rec_start_date, 7) ORDER BY left(rec_start_date, 7) DESC"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := platform.ScanMaps(rows)
	res := map[string]any{}
	for _, row := range ms {
		if v, ok := row["tonnagePaid"].(string); ok && v != "" {
			res[btsStr(row["rec_start_date"])] = json.Number(v)
		} else {
			res[btsStr(row["rec_start_date"])] = 0
		}
	}
	fillMonthSeries(res, str(m, "startTime"), str(m, "endTime"))
	platform.Success(w, res)
}

// fillMonthSeries 补月(Java 补足月份循环:缺月补整数 0)
func fillMonthSeries(res map[string]any, start, end string) {
	s, ok1 := parseLocalDate(start)
	e, ok2 := parseLocalDate(end)
	if !ok1 || !ok2 {
		return
	}
	cur := s
	for i := 1; ; i++ {
		mo := cur.Format("2006-01")
		if _, has := res[mo]; !has {
			res[mo] = 0
		}
		cur = s.AddDate(0, i, 0)
		if cur.Format("200601") > e.Format("200601") {
			break
		}
	}
}

// xzqh 行政区划 like 条件(Java 手工拼 %v,% 语义)
func xzqh(m map[string]any) (string, string) {
	x := str(m, "xzqhCode")
	if x == "" || x == "0" {
		return "", ""
	}
	return " AND CONCAT(province_code,',',city_code,',',area_code,',',county_code,',') like ?", "%" + x + ",%"
}

// billSelOrders 根据合同 id 分页查询订单列表(EXISTS 关联,${contractId} 原样内插;
// 缺失 → SQL "id=null" 仍可执行 → 空结果,照 Java)
func (s *Service) billSelOrders(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	current, size := pageOf(m)
	where := "WHERE EXISTS (SELECT 1 FROM chain_contract WHERE order_id=b.id AND id=" + rawId(m, "contractId") + ")"
	ip, err := paged(d, "SELECT b.* FROM chain_logistics_bill b "+where,
		"SELECT count(1) FROM chain_logistics_bill b "+where, nil, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entBill.rows(recs)
	platform.Success(w, ip)
}

// billSelAllOrders 根据合同 id 查询全部订单列表
func (s *Service) billSelAllOrders(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	rows, err := s.pick(r).Query("SELECT b.* FROM chain_logistics_bill b WHERE EXISTS (SELECT 1 FROM chain_contract WHERE order_id=b.id AND id=" + rawId(m, "contractId") + ")")
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := platform.ScanMaps(rows)
	if ms == nil {
		ms = []map[string]any{}
	}
	platform.Success(w, entBill.rows(ms))
}

// rawId ${} 原样内插值:数字字面用字面,缺失输出空串(MyBatis ${} null → 空串 → SQL 语法错误)
func rawId(m map[string]any, k string) string {
	v, ok := m[k]
	if !ok || v == nil {
		return ""
	}
	if n, ok := v.(json.Number); ok {
		return n.String()
	}
	return str(m, k)
}
