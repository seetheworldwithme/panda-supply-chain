// bill - billctl.go /chainBill 发票(主+Cg/Xs 变体,共用 service)。
// Java:ChainBillController + ChainBillServiceImpl + ChainBillMapper(core)。
// 实体无 fill;add/saveBill 含 ids(chain_bill_settlement 对应行)+filelist+日志 opType 2。
package bill

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
	dateOnly: []string{"open_date"}, // Java LocalDate
	carriers: []string{"searchStr", "page", "ids", "filelist", "timeType", "startTime",
		"endTime", "relSearch"},
	primFalse: []string{"relSearch"}, // 原生 boolean 载体键
}

var dtBillSpace = []string{"create_time"}
var dtBillDate = []string{"open_date", "start_time", "end_time"}
var dtBillT = []string{"update_time"}

var billEchoKeys = entBill.echoKeys()

func (s *Service) handleBill(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	switch {
	case strings.HasSuffix(p, "/list"):
		s.billList(w, r)
	case strings.HasSuffix(p, "/listByPage"):
		s.billListByPage(w, r)
	case strings.HasSuffix(p, "/getById"):
		s.billGetById(w, r)
	case strings.HasSuffix(p, "/add"):
		s.billAdd(w, r)
	case strings.HasSuffix(p, "/deleteById"):
		s.billDeleteById(w, r)
	case strings.HasSuffix(p, "/updateById"):
		s.billUpdateById(w, r)
	case strings.HasSuffix(p, "/saveBatch"):
		s.billSaveBatch(w, r, p)
	case strings.HasSuffix(p, "/importExcelJx"), strings.HasSuffix(p, "/importExcelXx"):
		s.billImportExcel(w, r, p)
	case strings.HasSuffix(p, "/exportExcel"):
		s.billExportExcel(w, r)
	case strings.HasSuffix(p, "/delBatch"):
		s.billDelBatch(w, r)
	case strings.HasSuffix(p, "/getSystemCode"):
		s.billGetSystemCode(w, r)
	case strings.HasSuffix(p, "/getSments"):
		s.billGetSments(w, r)
	case strings.HasSuffix(p, "/getSumPrice"):
		s.billGetSumPrice(w, r)
	case strings.HasSuffix(p, "/selectBusCount"):
		s.billSelectBusCount(w, r)
	case strings.HasSuffix(p, "/countBill"):
		s.billCountBill(w, r)
	case strings.HasSuffix(p, "/countTaxRate"):
		s.billCountTaxRate(w, r)
	case strings.HasSuffix(p, "/selZeroTax"):
		s.billSelZeroTax(w, r)
	case strings.HasSuffix(p, "/selJxzb"):
		s.billSelJxzb(w, r)
	case strings.HasSuffix(p, "/searchRelBill"):
		s.billSearchRelBill(w, r)
	case strings.HasSuffix(p, "/countBuses"):
		s.billCountBuses(w, r)
	case strings.HasSuffix(p, "/countMonBuses"):
		s.billCountMonBuses(w, r)
	case strings.HasSuffix(p, "/countPtTax"):
		if strings.HasPrefix(p, "/saas-bill/chainBill/") {
			s.billCountPtTax(w, r) // 主控制器用 countPtTaxNew
		} else {
			s.billCountPtTaxOld(w, r) // Cg/Xs 用老版(ONLY_FULL_GROUP_BY 恒 999)
		}
	case strings.HasSuffix(p, "/countYearBuses"):
		s.billCountYearBuses(w, r)
	default:
		platform.Unknown(w)
	}
}

// ── list:allEq(maps,false)(无 orderBy) ──

func (s *Service) billList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, entBill.cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	q := "SELECT " + entBill.selectList() + " FROM chain_bill"
	var args []any
	if len(cond) > 0 {
		var conds []string
		for _, k := range platform.SortedKeys(cond) {
			conds = append(conds, k+" = ?")
			args = append(args, cond[k])
		}
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	rows, err := d.Query(q, args...)
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

// billListByPageSearch getSql(排除 openDate/relationFlag)+手工四条
func billListByPageSearch(m map[string]any) (string, []any, bool) {
	where, args, ok := getSqlReflect(m, &entBill, []string{"openDate", "relationFlag"})
	if !ok {
		return "", nil, false
	}
	where += " " // Java 直接 += " and ..." → "where 1=1 and x like ? and open_date = ?"
	if v, ok2 := m["openDate"]; ok2 && v != nil {
		if s, isStr := v.(string); isStr && strings.TrimSpace(s) == "" {
			// 空串绑 null → getOpenDate()==null → 不拼条件
		} else {
			t, ok3 := parseLocalDate(str(m, "openDate"))
			if !ok3 {
				return "", nil, false
			}
			where += " and open_date = ?"
			args = append(args, t.Format("2006-01-02"))
		}
	}
	if v, ok2 := m["startTime"]; ok2 && v != nil && strings.TrimSpace(str(m, "startTime")) != "" {
		t, ok3 := parseLocalDate(str(m, "startTime"))
		if !ok3 {
			return "", nil, false
		}
		first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
		where += " and concat(month, '-01') >= ?"
		args = append(args, first.Format("2006-01-02"))
	}
	if v, ok2 := m["endTime"]; ok2 && v != nil && strings.TrimSpace(str(m, "endTime")) != "" {
		t, ok3 := parseLocalDate(str(m, "endTime"))
		if !ok3 {
			return "", nil, false
		}
		last := t.AddDate(0, 1, -1)
		where += " and concat(month, '-01') <= ?"
		args = append(args, last.Format("2006-01-02"))
	}
	if v, ok2 := m["relationFlag"]; ok2 && v != nil {
		if s, isStr := v.(string); isStr && strings.TrimSpace(s) == "" {
			// 空串绑 null 不拼
		} else if bv, isBool := v.(bool); isBool {
			where += " and relation_flag = ?"
			args = append(args, bv)
		} else if sv, isStr := v.(string); isStr && (sv == "true" || sv == "false") {
			where += " and relation_flag = ?"
			args = append(args, sv == "true")
		} else {
			return "", nil, false
		}
	}
	return where, args, true
}

const billConcat = "CONCAT(IFNULL(system_code, ''), IFNULL(contract_name, ''), IFNULL(bill_code, ''), IFNULL(bill_number, ''), IFNULL(open_date, ''), IFNULL(duty_paragraph, ''), IFNULL(company_name, ''), IFNULL(goods_name, ''), IFNULL(uit, ''), IFNULL(count, ''), IFNULL(unit_price, ''), IFNULL(base_money, ''), IFNULL(tax_rate, ''), IFNULL(tax_amount, ''), IFNULL(includ_money, ''), IFNULL(modle, ''))"

// searchStrLike getMaps(searchStr) 是 %v%(searchStr 是 String → % 包裹);缺失 → "%null%"
func searchStrLike(m map[string]any) string {
	if v, ok := m["searchStr"]; !ok || v == nil {
		return "%null%"
	}
	return "%" + str(m, "searchStr") + "%"
}

func (s *Service) billListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, dtBillSpace, dtBillDate, dtBillT) {
		platform.Unknown(w)
		return
	}
	searchSql, args, ok := billListByPageSearch(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	where := searchSql + " and " + billConcat + " like ? order by create_time desc"
	args = append(args, searchStrLike(m))
	current, size := pageOf(m)
	q := "select * from chain_bill " + where
	countQ := "SELECT count(1) FROM chain_bill " + searchSql + " and " + billConcat + " like ?"
	ip, err := paged(s.pick(r), q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entBill.rows(recs)
	platform.Success(w, ip)
}

func (s *Service) billGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	var row map[string]any
	if id, ok := int64Of(m, "id"); ok {
		var err error
		row, err = platform.GetByID(s.pick(r), "chain_bill", id)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, entBill.rows([]map[string]any{row})[0])
}

// updateFilelist 同 capital 语义
func updateFilelist(d *sql.DB, m map[string]any, id int64, u *auth.CurrentUser) {
	files, ok := m["filelist"].([]any)
	if !ok {
		return
	}
	for _, it := range files {
		f, ok := it.(map[string]any)
		if !ok || f["id"] == nil {
			continue
		}
		_, _ = platform.UpdateByIDNoFill(d, "chain_common_file", f["id"], map[string]any{
			"file_type": "1", "owner_type": "1", "owner_id": strconv.FormatInt(u.UserID, 10),
			"relation_id": strconv.FormatInt(id, 10), "relation_type": 1, "effective": true,
		})
	}
}

// billIds 子表维护:ids → chain_bill_settlement(bill_id, statement_id=ids[i])
func billIds(d *sql.DB, m map[string]any, billID int64) bool {
	ids, ok := m["ids"].([]any)
	if !ok {
		return false
	}
	if _, err := d.Exec("DELETE FROM chain_bill_settlement WHERE bill_id = ?", billID); err != nil {
		return false
	}
	for _, it := range ids {
		id, ok := it.(json.Number)
		if !ok {
			return false
		}
		if _, err := d.Exec("INSERT INTO chain_bill_settlement (statement_id, bill_id) VALUES (?, ?)", id.String(), billID); err != nil {
			return false
		}
	}
	return true
}

// billAdd saveBill:save(NoFill)+ids 子表+filelist+日志 opType 2;响应回显
func (s *Service) billAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, dtBillSpace, dtBillDate, dtBillT) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	vals := writeVals(&entBill, m)
	delete(vals, "id")
	delete(vals, "create_time")
	delete(vals, "update_time")
	id, err := platform.InsertNoFill(d, "chain_bill", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	billIds(d, m, id)
	updateFilelist(d, m, id, u)
	logOp(d, u, "新增发票\n"+changeLog(nil), 2)
	platform.Success(w, echoOf(&entBill, m, map[string]any{
		"id": json.Number(strconv.FormatInt(id, 10)),
	}))
}

func changeLog(before map[string]any) string {
	return "对比结果为：\n"
}

// billUpdateById updateBill:getById+updateById+ids 重挂+filelist+日志;恒回显
func (s *Service) billUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, dtBillSpace, dtBillDate, dtBillT) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	before, _ := platform.GetByID(d, "chain_bill", id)
	vals := writeVals(&entBill, m)
	delete(vals, "id")
	delete(vals, "create_time")
	delete(vals, "update_time")
	if _, err := platform.UpdateByIDNoFill(d, "chain_bill", id, vals); err != nil {
		platform.Unknown(w)
		return
	}
	billIds(d, m, id)
	updateFilelist(d, m, id, u)
	logOp(d, u, "修改发票\n"+changeLog(before), 2)
	platform.Success(w, echoOf(&entBill, m, nil))
}

// billDeleteById removeChainBill:removeById+bill_settlement(bill_id)+文件(relation_id)
func (s *Service) billDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	if err := platform.DeleteByID(d, "chain_bill", id); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("DELETE FROM chain_bill_settlement WHERE bill_id = ?", id); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("DELETE FROM chain_common_file WHERE relation_id = ?", id); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

// billDelBatch 手工 for(同 deleteById 副作用);ids null → 999
func (s *Service) billDelBatch(w http.ResponseWriter, r *http.Request) {
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
		if err := platform.DeleteByID(d, "chain_bill", id); err != nil {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("DELETE FROM chain_bill_settlement WHERE bill_id = ?", id); err != nil {
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

func toID(v any) (int64, bool) {
	if n, ok := v.(json.Number); ok {
		i, err := n.Int64()
		return i, err == nil
	}
	if s, ok := v.(string); ok {
		i, err := strconv.ParseInt(s, 10, 64)
		return i, err == nil
	}
	return 0, false
}

// billSaveBatch saveIrt(List):selByNumber 查重+selectByContractCode 回填+type 编号;
// errInfos 非空 → 不落库。主 Bill 响应 errInfos 数组;Cg/Xs 恒 "新增成功"。
func (s *Service) billSaveBatch(w http.ResponseWriter, r *http.Request, p string) {
	list, ok := bodyList(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	errInfos, ok2 := s.billSaveIrtRows(w, r, list, false)
	if !ok2 {
		return
	}
	if p == "/saas-bill/chainBill/saveBatch" {
		if errInfos == nil {
			errInfos = []string{}
		}
		platform.Success(w, errInfos)
	} else {
		platform.Success(w, "新增成功")
	}
}

// billSaveIrtRows 对照 saveIrt(List):返回 errInfos(带查重/回填/编号)。withBatch=true
// 时为 ImportResult 场景(批次号+日志由调用方处理)。rows 已写入即成功。
func (s *Service) billSaveIrtRows(w http.ResponseWriter, r *http.Request, list []any, importMode bool) ([]string, bool) {
	d := s.pick(r)
	u := currentUser(r)
	cg, ok1 := s.billSysCode(d, u, "PJJX.")
	xs, ok2 := s.billSysCode(d, u, "PJXX.")
	if !ok1 || !ok2 {
		platform.Unknown(w)
		return nil, false
	}
	var errInfos []string
	type prepared struct {
		vals map[string]any
	}
	var rows []prepared
	for _, it := range list {
		row, ok := it.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return nil, false
		}
		// selByNumber(billCode, billNumber):参数化 = 比较,任一为 null 恒假(Java 同形)
		var dup int
		_ = d.QueryRow("SELECT count(*) FROM chain_bill WHERE bill_code = ? AND bill_number = ?",
			nullOrNil(row["billCode"]), nullOrNil(row["billNumber"])).Scan(&dup)
		if dup > 0 {
			errInfos = append(errInfos, "代码为"+str(row, "billCode")+"，号码为"+str(row, "billNumber")+"的发票已存在，保存中止")
			continue
		}
		vals := writeVals(&entBill, row)
		delete(vals, "id")
		delete(vals, "create_time")
		delete(vals, "update_time")
		// selectByContractCode 回填
		if v, has := row["contractCode"]; has && v != nil {
			var cid sql.NullInt64
			var cname sql.NullString
			if err := d.QueryRow("SELECT id, contract_name FROM chain_contract WHERE contract_code = ?",
				str(row, "contractCode")).Scan(&cid, &cname); err == nil {
				vals["contract_id"] = json.Number(strconv.FormatInt(cid.Int64, 10))
				vals["contract_name"] = cname.String
			} else if err != sql.ErrNoRows {
				platform.Unknown(w)
				return nil, false
			}
		}
		switch str(row, "type") {
		case "1":
			vals["system_code"] = cg
			if n, okb := nextSystemCodeInMem(cg); okb {
				cg = n
			}
		case "2":
			vals["system_code"] = xs
			if n, okb := nextSystemCodeInMem(xs); okb {
				xs = n
			}
		}
		rows = append(rows, prepared{vals: vals})
	}
	if len(errInfos) == 0 {
		for _, pr := range rows {
			if _, err := platform.InsertNoFill(d, "chain_bill", pr.vals); err != nil {
				platform.Unknown(w)
				return nil, false
			}
		}
	}
	return errInfos, true
}

func (s *Service) billSysCode(d *sql.DB, u *auth.CurrentUser, idt string) (string, bool) {
	like := idt + u.EnterpriseCode + strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", "")
	return systemCodeNext(d, "chain_bill", like, "1")
}

// billGetSystemCode codeType=1→PJJX. 2→PJXX.
func (s *Service) billGetSystemCode(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	u := currentUser(r)
	var idt string
	switch str(m, "codeType") {
	case "1":
		idt = "PJJX."
	case "2":
		idt = "PJXX."
	}
	code, ok := s.billSysCode(s.pick(r), u, idt)
	if !ok {
		platform.Unknown(w)
		return
	}
	platform.Success(w, code)
}

// billGetSments mapper:三表连接 SELECT s.*(Settlement 形状)
func (s *Service) billGetSments(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := s.pick(r).Query("SELECT s.* from chain_capital_settlement s, chain_bill b, chain_bill_settlement bs where s.id=bs.statement_id and bs.bill_id = b.id and b.id = ?", id)
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
	platform.Success(w, entSettShape.rows(ms))
}

// billGetSumPrice contractId 有 → 表切换 s,b
func (s *Service) billGetSumPrice(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	where := "from chain_bill b where 1=1 "
	var args []any
	if v, ok := m["contractId"]; ok && v != nil {
		where = "from chain_capital_settlement s, chain_bill b where 1=1 and s.relation_id=b.relation_id "
		where += " and (b.relation_id is not null and b.relation_id !='') and s.contract_id = ?"
		args = append(args, str(m, "contractId"))
	}
	if v, ok := m["type"]; ok && v != nil && !isEmptyStr(v) {
		where += " and b.type = ?"
		args = append(args, v)
	}
	if v, ok := m["relationFlag"]; ok && v != nil {
		if bv, isBool := v.(bool); isBool {
			where += " and b.relation_flag=?"
			args = append(args, bv)
		} else if sv, isStr := v.(string); isStr && (sv == "true" || sv == "false") {
			where += " and b.relation_flag=?"
			args = append(args, sv == "true")
		} else {
			platform.Unknown(w)
			return
		}
	}
	if v, ok := m["relationId"]; ok && v != nil && !isEmptyStr(v) {
		where += " and b.relation_id = ?"
		args = append(args, v)
	}
	rows, err := s.pick(r).Query("select IFNULL(sum(b.includ_money), 0) sumMoney, IFNULL(sum(b.base_money), 0) baseMoney, IFNULL(sum(b.tax_amount), 0) taxAmount, IFNULL(sum(b.count), 0) count "+where, args...)
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
	for _, k := range []string{"sumMoney", "baseMoney", "taxAmount", "count"} {
		if s2, ok := ms[0][k].(string); ok {
			out[k] = json.Number(s2)
		} else {
			out[k] = ms[0][k]
		}
	}
	platform.Success(w, out)
}

// billSelectBusCount jwt companyId!="0" 强制;="0" 用 body companyId;响应 DECIMAL 保真
func (s *Service) billSelectBusCount(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	u := currentUser(r)
	d := s.pick(r)
	where := " where 1=1"
	var args []any
	if v, ok := m["startTime"]; ok && v != nil && !isEmptyStr(v) {
		where += " and concat(month, '-01') >= ?"
		args = append(args, str(m, "startTime"))
	}
	if v, ok := m["endTime"]; ok && v != nil && !isEmptyStr(v) {
		where += " and concat(month, '-01') <= ?"
		args = append(args, str(m, "endTime"))
	}
	if u.CompanyID != 0 {
		where += " and company_id = ?"
		args = append(args, json.Number(strconv.FormatInt(u.CompanyID, 10)))
	} else if v, ok := m["companyId"]; ok && v != nil {
		where += " and company_id = ?"
		args = append(args, v)
	}
	where += " group by type"
	rows, err := d.Query("select type, IFNULL(sum(tax_amount), 0) sumTaxAmount, IFNULL(sum(includ_money), 0) sumIncludMoney from chain_bill"+where, args...)
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

// billCountBill timeType 拆箱 999;List<ChainBill> 全实体
func (s *Service) billCountBill(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	tt, ok := m["timeType"]
	if !ok || tt == nil {
		platform.Unknown(w) // 1 == timeType 拆箱 NPE
		return
	}
	where := " where 1=1"
	switch int(num(tt)) {
	case 1:
		where += " and DATE_FORMAT(open_date, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')"
	case 2:
		where += " and DATE_FORMAT(open_date, '%Y-%m') = DATE_FORMAT(DATE_SUB(curdate(), INTERVAL 1 MONTH),'%Y-%m')"
	case 3:
		where += " and open_date between date_sub(now(),interval 6 month) and now()"
	case 4:
		where += " and YEAR(open_date)=YEAR(NOW())"
	}
	rows, err := s.pick(r).Query("select * from chain_bill" + where)
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

func nullOrNil(v any) any {
	if v == nil {
		return nil
	}
	return str(map[string]any{"k": v}, "k")
}
