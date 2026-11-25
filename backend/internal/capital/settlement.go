// capital - settlement.go /chainCapitalSettlement 结算单(主+Cg/Xs 变体,共用 service)。
// Java:ChainCapitalSettlementController + ChainCapitalSettlementServiceImpl +
// ChainCapitalSettlementMapper(core)。无 fill 实体;add 的 khxje 被 sumPrice 覆盖。
package capital

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// entBill ChainBill 实体(capital 域视图:getBills 回显;bill 包另有完整版)
var entBill = entShape{
	table: "chain_bill",
	cols: []string{"id", "im_batch_number", "company_id", "system_code", "contract_id",
		"contract_code", "contract_name", "in_type", "bill_code", "bill_number",
		"open_date", "duty_paragraph", "company_name", "goods_name", "uit", "count",
		"unit_price", "base_money", "tax_rate", "tax_amount", "includ_money", "modle",
		"month", "type", "relation_flag", "relation_id", "create_time", "update_time",
		"remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"count", "unit_price", "base_money", "tax_amount", "includ_money"},
	bools:    []string{"relation_flag"},
	tTimes:   []string{"update_time"},
	carriers: []string{"searchStr", "page", "ids", "filelist", "timeType", "startTime",
		"endTime", "relSearch"},
}

var entSettlement = entShape{
	table: "chain_capital_settlement",
	cols: []string{"id", "im_batch_number", "system_code", "company_id", "contract_id",
		"contract_code", "contract_name", "goods_name", "set_start_time", "set_end_time",
		"settlement_count", "unit_price", "sum_price", "khxje", "sjhxje", "type",
		"supplier_id", "supplier_name", "rel_bill_flag", "rel_detail_flag", "hx_status",
		"hx_code", "month", "relation_id", "create_time", "update_time",
		"remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"settlement_count", "unit_price", "sum_price", "khxje", "sjhxje"},
	bools:    []string{"rel_bill_flag", "rel_detail_flag"},
	tTimes:   []string{"update_time"},
	dateOnly: []string{"set_start_time", "set_end_time"}, // Java LocalDate
	carriers: []string{"searchStr", "page", "ids", "details", "filelist", "timeType",
		"khxFlag", "sments", "bills", "startTime", "endTime"},
	primFalse: []string{"khxFlag"},
}

var dtSettSpace = []string{"create_time"}
var dtSettDate = []string{"set_start_time", "set_end_time", "start_time", "end_time"}
var dtSettT = []string{"update_time"}

var stEchoKeys = entSettlement.echoKeys()

func (s *Service) handleSettlement(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	switch {
	case strings.HasSuffix(p, "/list"):
		s.stList(w, r)
	case strings.HasSuffix(p, "/listByPage"):
		s.stListByPage(w, r)
	case strings.HasSuffix(p, "/getById"):
		s.stGetById(w, r)
	case strings.HasSuffix(p, "/add"):
		s.stAdd(w, r)
	case strings.HasSuffix(p, "/deleteById"):
		s.stDeleteById(w, r)
	case strings.HasSuffix(p, "/updateById"):
		s.stUpdateById(w, r)
	case strings.HasSuffix(p, "/saveBatchSments"):
		s.stSaveBatch(w, r)
	case strings.HasSuffix(p, "/importAlone"):
		s.stImportAlone(w, r)
	case strings.HasSuffix(p, "/importExcel"):
		s.stImportExcel(w, r)
	case strings.HasSuffix(p, "/exportExcel"):
		s.stExportExcel(w, r)
	case strings.HasSuffix(p, "/delBatch"):
		s.stDelBatch(w, r)
	case strings.HasSuffix(p, "/getSystemCode"):
		s.stGetSystemCode(w, r)
	case strings.HasSuffix(p, "/getDetails"):
		s.stGetDetails(w, r)
	case strings.HasSuffix(p, "/getBills"):
		s.stGetBills(w, r)
	case strings.HasSuffix(p, "/selectQuantity"):
		s.stSelectQuantity(w, r)
	case strings.HasSuffix(p, "/countSett"):
		s.stCountSett(w, r)
	case strings.HasSuffix(p, "/getSumPrice"):
		s.stGetSumPrice(w, r)
	case strings.HasSuffix(p, "/getPtSumPrice"):
		s.stGetPtSumPrice(w, r)
	case strings.HasSuffix(p, "/relSmentsBills"):
		s.stRelSmentsBills(w, r)
	default:
		platform.Unknown(w)
	}
}

// stList allEq(maps,false)(Settlement 无 orderBy)
func (s *Service) stList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, entSettlement.cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	q := "SELECT " + entSettlement.selectList() + " FROM chain_capital_settlement"
	var args []any
	if len(cond) > 0 {
		var conds []string
		for _, k := range platform.SortedKeys(cond) {
			conds = append(conds, k+" = ?")
			args = append(args, cond[k])
		}
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	rows0, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ScanMaps(rows0)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, entSettlement.rows(rows))
}

// stListByPageSearch 手工 searchSql(参数化)
func stListByPageSearch(m map[string]any) (string, []any, bool) {
	where := "where 1=1 "
	var args []any
	if v, ok := m["type"]; ok && v != nil && !isEmptyStr(v) {
		where += " and type = ?"
		args = append(args, v)
	}
	if bv, ok := m["khxFlag"]; ok && bv == true {
		where += " and (hx_status = 1 or hx_status = 2)"
	}
	if v, ok := m["hxCode"]; ok && v != nil {
		where += " and hx_code like ?"
		args = append(args, "%"+str(m, "hxCode")+"%")
	}
	if v, ok := m["supplierId"]; ok && v != nil && !isEmptyStr(v) {
		where += " and supplier_id = ?"
		args = append(args, v)
	}
	if v, ok := m["supplierName"]; ok && v != nil {
		where += " and supplier_name like ?"
		args = append(args, "%"+str(m, "supplierName")+"%")
	}
	if v, ok := m["contractId"]; ok && v != nil && !isEmptyStr(v) {
		where += " and contract_id = ?"
		args = append(args, v)
	}
	if v, ok := m["setStartTime"]; ok && v != nil && strings.TrimSpace(str(m, "setStartTime")) != "" {
		t, ok2 := parseLocalDate(str(m, "setStartTime"))
		if !ok2 {
			return "", nil, false
		}
		where += " and set_start_time >= ?"
		args = append(args, t.Format("2006-01-02"))
	}
	if v, ok := m["setEndTime"]; ok && v != nil && strings.TrimSpace(str(m, "setEndTime")) != "" {
		t, ok2 := parseLocalDate(str(m, "setEndTime"))
		if !ok2 {
			return "", nil, false
		}
		where += " and set_end_time <= ?"
		args = append(args, t.Format("2006-01-02"))
	}
	if v, ok := m["relBillFlag"]; ok && v != nil {
		if s, isStr := v.(string); isStr && strings.TrimSpace(s) == "" {
			// 空串绑 null(Jackson 对 Boolean ""→null),不拼条件
		} else if bv, isBool := v.(bool); isBool {
			where += " and rel_bill_flag = ?"
			args = append(args, bv)
		} else if sv, isStr := v.(string); isStr && (sv == "true" || sv == "false") {
			where += " and rel_bill_flag = ?"
			args = append(args, sv == "true")
		} else {
			return "", nil, false // Jackson Boolean 绑定失败 → 999
		}
	}
	if v, ok := m["relationId"]; ok && v != nil && !isEmptyStr(v) {
		where += " and relation_id = ?"
		args = append(args, v)
	}
	if v, ok := m["month"]; ok && v != nil && strings.TrimSpace(str(m, "month")) != "" {
		where += " and month = ?"
		args = append(args, str(m, "month"))
	}
	if v, ok := m["startTime"]; ok && v != nil && strings.TrimSpace(str(m, "startTime")) != "" {
		t, ok2 := parseLocalDate(str(m, "startTime"))
		if !ok2 {
			return "", nil, false
		}
		first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
		where += " and concat(month, '-01') >= ?"
		args = append(args, first.Format("2006-01-02"))
	}
	if v, ok := m["endTime"]; ok && v != nil && strings.TrimSpace(str(m, "endTime")) != "" {
		t, ok2 := parseLocalDate(str(m, "endTime"))
		if !ok2 {
			return "", nil, false
		}
		last := t.AddDate(0, 1, -1)
		where += " and concat(month, '-01') <= ?"
		args = append(args, last.Format("2006-01-02"))
	}
	return where, args, true
}

// stListByPage mapper:显式列 IFNULL(khxje,sum_price)+CONCAT 10 列,order by set_start_time desc
func (s *Service) stListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["page"].(map[string]any); !ok {
		platform.Unknown(w) // Java mapper requires a Page object.
		return
	}
	if !validateDT(m, dtSettSpace, dtSettDate, dtSettT) {
		platform.Unknown(w)
		return
	}
	searchSql, args, ok := stListByPageSearch(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	where := searchSql + " and CONCAT(IFNULL(system_code, ''), IFNULL(contract_code, ''), IFNULL(contract_name, ''), IFNULL(goods_name, ''), IFNULL(set_start_time, ''), IFNULL(set_end_time, ''), IFNULL(settlement_count, ''), IFNULL(unit_price, ''), IFNULL(sum_price, ''), IFNULL(supplier_name, '') ) like ? order by set_start_time desc"
	args = append(args, searchStrLike(m))
	current, size := pageOf(m)
	q := "select id,system_code,company_id,contract_id,contract_code,contract_name,goods_name,set_start_time,set_end_time,settlement_count,unit_price,sum_price,IFNULL(khxje, sum_price) khxje,IFNULL(khxje, sum_price) sjhxje,type,supplier_id,supplier_name,rel_bill_flag,rel_detail_flag,hx_status,hx_code,create_time,update_time,remark1,remark2,remark3,remark4,remark5,month from chain_capital_settlement " + where
	countQ := "SELECT count(1) FROM chain_capital_settlement " + searchSql + " and CONCAT(IFNULL(system_code, ''), IFNULL(contract_code, ''), IFNULL(contract_name, ''), IFNULL(goods_name, ''), IFNULL(set_start_time, ''), IFNULL(set_end_time, ''), IFNULL(settlement_count, ''), IFNULL(unit_price, ''), IFNULL(sum_price, ''), IFNULL(supplier_name, '') ) like ?"
	ip, err := paged(s.pick(r), q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entSettlement.rows(recs)
	platform.Success(w, ip)
}

func (s *Service) stGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	var row map[string]any
	if id, ok := int64Of(m, "id"); ok {
		var err error
		row, err = platform.GetByID(s.pick(r), "chain_capital_settlement", id)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, entSettlement.rows([]map[string]any{row})[0])
}

// stAdd saveSettlement:khxje=sumPrice 覆盖(为 null 则 SET 不含该列)→save→filelist→日志 type4
func (s *Service) stAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, dtSettSpace, dtSettDate, dtSettT) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	vals := pick(m, entSettlement.cols...)
	delete(vals, "id")
	delete(vals, "create_time")
	delete(vals, "update_time")
	delete(vals, "khxje")
	if v, ok := m["sumPrice"]; ok && v != nil {
		vals["khxje"] = v
	}
	id, err := platform.InsertNoFill(d, "chain_capital_settlement", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	updateFilelist(d, m, id, u)
	logOp(d, u, "新增结算单\n"+changeLog(nil, m, stEchoKeys), 4)
	// 回显:Java 内存实体 khxje 已被 set 成 sumPrice
	echo := echoOf(&entSettlement, m, map[string]any{
		"id": json.Number(strconv.FormatInt(id, 10)),
	})
	if v, ok := m["sumPrice"]; ok {
		echo["khxje"] = v
	}
	platform.Success(w, echo)
}

// stUpdateById updateSettlement:getById(日志前值)+updateById+filelist+日志;恒回显
func (s *Service) stUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, dtSettSpace, dtSettDate, dtSettT) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	id, idok := int64Of(m, "id")
	if !idok {
		platform.Unknown(w)
		return
	}
	before, _ := platform.GetByID(d, "chain_capital_settlement", id)
	vals := pick(m, entSettlement.cols...)
	delete(vals, "id")
	delete(vals, "create_time")
	delete(vals, "update_time")
	if _, err := platform.UpdateByIDNoFill(d, "chain_capital_settlement", id, vals); err != nil {
		platform.Unknown(w)
		return
	}
	updateFilelist(d, m, id, u)
	logOp(d, u, "修改结算单\n"+changeLog(before, m, stEchoKeys), 4)
	platform.Success(w, echoOf(&entSettlement, m, nil))
}

// stDeleteById removeChainCapitalSettlement:getById(未用)+removeById+
// detail_settlement/bill_settlement deleteByMap(statement_id)+文件
func (s *Service) stDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	if err := platform.DeleteByID(d, "chain_capital_settlement", id); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("DELETE FROM chain_capital_detail_settlement WHERE statement_id = ?", id); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("DELETE FROM chain_bill_settlement WHERE statement_id = ?", id); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("DELETE FROM chain_common_file WHERE relation_id = ?", id); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

// stDelBatch 手工 for 逐条(同 deleteById 副作用);ids null → 999
func (s *Service) stDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, ok := m["ids"].([]any)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	for _, it := range ids {
		id, ok := toID(it)
		if !ok {
			platform.Unknown(w)
			return
		}
		if err := platform.DeleteByID(d, "chain_capital_settlement", id); err != nil {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("DELETE FROM chain_capital_detail_settlement WHERE statement_id = ?", id); err != nil {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("DELETE FROM chain_bill_settlement WHERE statement_id = ?", id); err != nil {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("DELETE FROM chain_common_file WHERE relation_id = ?", id); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "删除成功")
}

// stSaveBatch saveIrt(List):type=1→JSCG. 2→JSXS. 分配;响应 errInfos(恒空列表)
func (s *Service) stSaveBatch(w http.ResponseWriter, r *http.Request) {
	list, ok := bodyList(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	cg, ok1 := s.settSysCode(d, u, "JSCG.")
	xs, ok2 := s.settSysCode(d, u, "JSXS.")
	if !ok1 || !ok2 {
		platform.Unknown(w)
		return
	}
	for _, it := range list {
		row, ok := it.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		vals := pick(row, entSettlement.cols...)
		delete(vals, "id")
		if n, err := intType(row, "type"); err == nil {
			switch n {
			case 1:
				vals["system_code"] = cg
				cg, _ = nextSystemCodeInMem(cg)
			case 2:
				vals["system_code"] = xs
				xs, _ = nextSystemCodeInMem(xs)
			}
		}
		if _, err := platform.InsertNoFill(d, "chain_capital_settlement", vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, []any{})
}

func intType(m map[string]any, k string) (int64, error) {
	v, ok := m[k]
	if !ok || v == nil {
		return 0, errNotFound
	}
	if n, ok := v.(json.Number); ok {
		return n.Int64()
	}
	return strconv.ParseInt(str(m, k), 10, 64)
}

type simpleErr struct{ msg string }

func (e *simpleErr) Error() string { return e.msg }

var errNotFound = &simpleErr{"not found"}

func (s *Service) settSysCode(d *sql.DB, u *auth.CurrentUser, idt string) (string, bool) {
	like := idt + u.EnterpriseCode + strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", "")
	return systemCodeNext(d, "chain_capital_settlement", like, "1")
}

// stGetSystemCode codeType=1→JSCG. 2→JSXS. 空→""
func (s *Service) stGetSystemCode(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	u := currentUser(r)
	var idt string
	switch str(m, "codeType") {
	case "1":
		idt = "JSCG."
	case "2":
		idt = "JSXS."
	}
	code, ok := s.settSysCode(s.pick(r), u, idt)
	if !ok {
		platform.Unknown(w)
		return
	}
	platform.Success(w, code)
}

// stGetDetails mapper:SELECT d.*(Detailed 实体) 三表连接
func (s *Service) stGetDetails(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := s.pick(r).Query("SELECT d.id,d.im_batch_number,d.system_code,d.company_id,d.account_name,d.account,d.trans_time,d.money,d.opposite_name,d.opposite_account,d.abs_content,d.purpose,d.detype,d.khxje,d.sjhxje,d.hx_status,d.relation_flag,d.bank_name,d.hx_code,d.create_time,d.update_time,d.remark1,d.remark2,d.remark3,d.remark4,d.remark5 from chain_capital_settlement s, chain_capital_detailed d, chain_capital_detail_settlement sd where s.id=sd.statement_id and sd.detail_id = d.id and s.id = ?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if ms == nil {
		ms = []map[string]any{}
	}
	platform.Success(w, entDetailed.rows(ms))
}

// stGetBills mapper:SELECT b.*(ChainBill 实体) 三表连接
func (s *Service) stGetBills(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := s.pick(r).Query("SELECT b.* from chain_capital_settlement s, chain_bill b, chain_bill_settlement bs where s.id=bs.statement_id and bs.bill_id = b.id and s.id = ?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if ms == nil {
		ms = []map[string]any{}
	}
	platform.Success(w, entBill.rows(ms))
}

// stSelectQuantity selectQuat:Integer 拆箱,空 timeType → 999(Cg/Xs 路径;主路径 dashboard)
func (s *Service) stSelectQuantity(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	tt, ok := m["timeType"]
	if !ok || tt == nil {
		platform.Unknown(w) // 1 == timeType 拆箱 NPE
		return
	}
	n := int(num(tt))
	where := " where 1=1"
	switch n {
	case 1:
		where += " and month = DATE_FORMAT(CURDATE(), '%Y-%m')"
	case 2:
		where += " and month = DATE_FORMAT(DATE_SUB(curdate(), INTERVAL 1 MONTH),'%Y-%m')"
	case 3:
		where += " and concat(month, '-01') between date_sub(now(),interval 6 month) and now()"
	case 4:
	}
	where += " group by type"
	rows, err := s.pick(r).Query("select type, sum(settlement_count) sumNum from chain_capital_settlement" + where)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if ms == nil {
		ms = []map[string]any{}
	}
	platform.Success(w, ms)
}

// stCountSett 日期统计(Cg/Xs 路径;主路径 dashboard)。月粒度恒定,响应 {cglb, xslb}。
func (s *Service) stCountSett(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	series, _, ok := countSeries(m, false)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	q := "SELECT month t, type, sum(sum_price) sumPrice FROM chain_capital_settlement WHERE 1=1 GROUP BY month, type ORDER BY t, type"
	rows, err := d.Query(q)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, errS := platform.ScanMaps(rows)
	if errS != nil { // 迭代中断连时原先返回全 0 的"成功"
		platform.Unknown(w)
		return
	}
	st, et := str(m, "startTime"), str(m, "endTime")
	cglb, xslb, ok2 := bucketSeries(series, ms, "type", "1", "2", st, et)
	if !ok2 {
		platform.Unknown(w) // 补月循环 getStartTime() NPE → 999
		return
	}
	platform.Success(w, map[string]any{"cglb": cglb, "xslb": xslb})
}

// stGetSumPrice mapper:IFNULL(sum(sum_price),0) sumMoney, IFNULL(sum(settlement_count),0) settlementCount
func (s *Service) stGetSumPrice(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	where := " where 1=1"
	var args []any
	if v, ok := m["type"]; ok && v != nil && !isEmptyStr(v) {
		where += " and type = ?"
		args = append(args, v)
	}
	if v, ok := m["relBillFlag"]; ok && v != nil {
		bv, isBool := v.(bool)
		sv, isStr := v.(string)
		if !isBool && !(isStr && (sv == "true" || sv == "false")) {
			platform.Unknown(w)
			return
		}
		if !isBool {
			bv = sv == "true"
		}
		where += " and rel_bill_flag = ?"
		args = append(args, bv)
	}
	if v, ok := m["contractId"]; ok && v != nil && !isEmptyStr(v) {
		where += " and contract_id = ?"
		args = append(args, v)
	}
	if v, ok := m["relationId"]; ok && v != nil && !isEmptyStr(v) {
		where += " and relation_id = ?"
		args = append(args, v)
	}
	rows, err := s.pick(r).Query("select IFNULL(sum(sum_price), 0) sumMoney, IFNULL(sum(settlement_count), 0) settlementCount from chain_capital_settlement"+where, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil || len(ms) == 0 {
		platform.Unknown(w)
		return
	}
	out := map[string]any{}
	for _, k := range []string{"sumMoney", "settlementCount"} {
		if s2, ok := ms[0][k].(string); ok {
			out[k] = json.Number(s2)
		} else {
			out[k] = ms[0][k]
		}
	}
	platform.Success(w, out)
}

// stGetPtSumPrice 平台结算合计:sum(sum_price) 无 IFNULL → 空集单 NULL 列 → data null
func (s *Service) stGetPtSumPrice(w http.ResponseWriter, r *http.Request) {
	where, args, ok := ptSumPriceCond(r, bodyMap(r), "sum(sum_price) sumPrice from chain_capital_settlement b, chain_basic_enterprise e ")
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := s.pick(r).Query("SELECT "+where, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil || len(ms) == 0 {
		platform.Unknown(w)
		return
	}
	v := ms[0]["sumPrice"]
	if v == nil {
		platform.Success(w, nil) // MyBatis 单 NULL 列 → null map → data null
		return
	}
	platform.Success(w, map[string]any{"sumPrice": json.Number(v.(string))})
}

// stRelSmentsBills 关联发票:save 结算记录+结算单批更+bills 批更。
// sments null/空 → NPE/IndexOutOfBounds → 999。
func (s *Service) stRelSmentsBills(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	sments, _ := m["sments"].([]any)
	bills, _ := m["bills"].([]any)
	if sments == nil || len(sments) == 0 {
		platform.Unknown(w) // get(0) IndexOutOfBounds → 999
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	first, ok := sments[0].(map[string]any)
	if !ok {
		platform.Unknown(w)
		return
	}
	stType, ok1 := int64Of(first, "type")
	if !ok1 {
		platform.Unknown(w) // 构造器 Integer null → 999
		return
	}
	res, err := platform.InsertNoFill(d, "chain_bill_sett_rec", map[string]any{
		"sett_num": json.Number(strconv.Itoa(len(sments))),
		"bill_num": json.Number(strconv.Itoa(len(bills))),
		"operater": u.Username,
		"type":     json.Number(strconv.FormatInt(stType, 10)),
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	for _, it := range sments {
		row, ok := it.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		id, ok := int64Of(row, "id")
		if !ok {
			platform.Unknown(w)
			return
		}
		if _, err := platform.UpdateByIDNoFill(d, "chain_capital_settlement", id, map[string]any{
			"relation_id":   json.Number(strconv.FormatInt(res, 10)),
			"rel_bill_flag": true,
		}); err != nil {
			platform.Unknown(w)
			return
		}
	}
	for _, it := range bills {
		row, ok := it.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		id, ok := int64Of(row, "id")
		if !ok {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("update chain_bill set relation_id=?, relation_flag=true where id = ?", res, id); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "成功")
}
