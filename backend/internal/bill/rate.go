// bill - rate.go /chainBillRate 缴税 + sett.go 的 Settlement 形状常驻本包;本文件为
// BillRate 12 端点(saveBatch/getSumPrice/countMonBuses 由 handleRate 分发到本文件,
// getSumPrice/countMonBuses 原 dashboard 注册移交至此精确实现)。
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

var entRate = entShape{
	table: "chain_bill_rate",
	cols: []string{"id", "im_batch_number", "company_id", "system_code", "month",
		"base_money", "tax_amount", "includ_money", "jx_money", "zz_money", "cj_money",
		"jy_money", "dfjy_money", "yh_money", "qy_money", "znj_money", "sum_money",
		"remark", "create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"base_money", "tax_amount", "includ_money", "jx_money", "zz_money",
		"cj_money", "jy_money", "dfjy_money", "yh_money", "qy_money", "znj_money", "sum_money"},
	carriers: []string{"searchStr", "page", "ids", "startTime", "endTime"},
}

var entSettShape = entShape{
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
	dateOnly: []string{"set_start_time", "set_end_time"},
	carriers: []string{"searchStr", "page", "ids", "details", "filelist", "timeType",
		"khxFlag", "sments", "bills", "startTime", "endTime"},
	primFalse: []string{"khxFlag"},
}

var entSettRow = entShape{
	table: "chain_bill_settlement",
	cols: []string{"id", "company_id", "statement_id", "bill_id", "create_time", "update_time",
		"remark1", "remark2", "remark3", "remark4", "remark5"},
	carriers: []string{"searchStr", "page", "ids"},
}

var rateEchoKeys = entRate.echoKeys()

func (s *Service) handleRate(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	switch {
	case strings.HasSuffix(p, "/list"):
		s.rateList(w, r)
	case strings.HasSuffix(p, "/listByPage"):
		s.rateListByPage(w, r)
	case strings.HasSuffix(p, "/getById"):
		s.rateGetById(w, r)
	case strings.HasSuffix(p, "/add"):
		s.rateAdd(w, r)
	case strings.HasSuffix(p, "/deleteById"):
		s.rateDeleteById(w, r)
	case strings.HasSuffix(p, "/updateById"):
		s.rateUpdateById(w, r)
	case strings.HasSuffix(p, "/importExcel"):
		s.rateImportExcel(w, r)
	case strings.HasSuffix(p, "/exportExcel"):
		s.rateExportExcel(w, r)
	case strings.HasSuffix(p, "/delBatch"):
		s.rateDelBatch(w, r)
	case strings.HasSuffix(p, "/saveBatch"):
		s.rateSaveBatch(w, r)
	case strings.HasSuffix(p, "/getSumPrice"):
		s.rateGetSumPrice(w, r)
	case strings.HasSuffix(p, "/countMonBuses"):
		s.rateCountMonBuses(w, r)
	default:
		platform.Unknown(w)
	}
}

func (s *Service) rateList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, entRate.cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_bill_rate", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, entRate.rows(rows))
}

// rateListByPage getSql 全字段+CONCAT(id) like
func (s *Service) rateListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, []string{"create_time"}, []string{}, []string{"update_time"}) {
		platform.Unknown(w)
		return
	}
	searchSql, args, ok := getSqlReflect(m, &entRate, nil)
	if !ok {
		platform.Unknown(w)
		return
	}
	where := searchSql + " and CONCAT(id) like ? order by create_time desc"
	args = append(args, searchStrLike(m))
	current, size := pageOf(m)
	q := "select * from chain_bill_rate " + where
	countQ := "SELECT count(1) FROM chain_bill_rate " + searchSql + " and CONCAT(id) like ?"
	ip, err := paged(s.pick(r), q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entRate.rows(recs)
	platform.Success(w, ip)
}

func (s *Service) rateGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	var row map[string]any
	if id, ok := int64Of(m, "id"); ok {
		var err error
		row, err = platform.GetByID(s.pick(r), "chain_bill_rate", id)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, entRate.rows([]map[string]any{row})[0])
}

// rateAdd 纯 MP save+回显(无 fill)
func (s *Service) rateAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, []string{"create_time"}, []string{}, []string{"update_time"}) {
		platform.Unknown(w)
		return
	}
	vals := writeVals(&entRate, m)
	delete(vals, "id")
	delete(vals, "create_time")
	delete(vals, "update_time")
	id, err := platform.InsertNoFill(s.pick(r), "chain_bill_rate", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, echoOf(&entRate, m, map[string]any{
		"id": json.Number(strconv.FormatInt(id, 10)),
	}))
}

func (s *Service) rateUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, []string{"create_time"}, []string{}, []string{"update_time"}) {
		platform.Unknown(w)
		return
	}
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	vals := writeVals(&entRate, m)
	delete(vals, "id")
	delete(vals, "create_time")
	delete(vals, "update_time")
	if _, err := platform.UpdateByIDNoFill(s.pick(r), "chain_bill_rate", id, vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, echoOf(&entRate, m, nil))
}

func (s *Service) rateDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	if err := platform.DeleteByID(s.pick(r), "chain_bill_rate", id); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

// rateDelBatch MP removeByIds(null → 短路 200)
func (s *Service) rateDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, ok := m["ids"].([]any)
	if !ok {
		platform.Success(w, "删除成功")
		return
	}
	if len(ids) == 0 {
		platform.Success(w, "删除成功")
		return
	}
	numIds := make([]any, 0, len(ids))
	for _, it := range ids {
		n, ok := toID(it)
		if !ok {
			platform.Unknown(w)
			return
		}
		numIds = append(numIds, json.Number(strconv.FormatInt(n, 10)))
	}
	if err := platform.DeleteByIDs(s.pick(r), "chain_bill_rate", numIds); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// rateSaveBatch saveBatchRate → saveIrt(List):PJSS. 逐条 save+日志 opType 7;响应 "新增成功"
func (s *Service) rateSaveBatch(w http.ResponseWriter, r *http.Request) {
	list, ok := bodyList(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	code, ok1 := s.rateSysCode(d, u)
	if !ok1 {
		platform.Unknown(w)
		return
	}
	for _, it := range list {
		row, ok := it.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		vals := writeVals(&entRate, row)
		delete(vals, "id")
		delete(vals, "create_time")
		delete(vals, "update_time")
		vals["system_code"] = code
		n, okb := nextSystemCodeInMem(code)
		if !okb {
			platform.Unknown(w)
			return
		}
		code = n
		newID, err := platform.InsertNoFill(d, "chain_bill_rate", vals)
		if err != nil {
			platform.Unknown(w)
			return
		}
		logOp(d, u, "新增税收数据\n"+rateChangeLog(row, newID), 7)
	}
	platform.Success(w, "新增成功")
}

func rateChangeLog(row map[string]any, id int64) string {
	return "对比结果为：\n"
}

func (s *Service) rateSysCode(d *sql.DB, u *auth.CurrentUser) (string, bool) {
	like := "PJSS." + u.EnterpriseCode + strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", "")
	return systemCodeNext(d, "chain_bill_rate", like, "1")
}

// rateGetSumPrice getSql(排除 companyId)+companyId 等值+month 比较原值(无首末日转换)
func (s *Service) rateGetSumPrice(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, []string{"create_time"}, []string{"start_time", "end_time"}, []string{"update_time"}) {
		platform.Unknown(w)
		return
	}
	searchSql, args, ok := getSqlReflect(m, &entRate, []string{"companyId"})
	if !ok {
		platform.Unknown(w)
		return
	}
	searchSql += " "
	if v, has := m["companyId"]; has && v != nil {
		searchSql += " and company_id = ?"
		args = append(args, v)
	}
	if v, has := m["startTime"]; has && v != nil && !isEmptyStr(v) {
		searchSql += " and concat(month, '-01') >= ?"
		args = append(args, str(m, "startTime"))
	}
	if v, has := m["endTime"]; has && v != nil && !isEmptyStr(v) {
		searchSql += " and concat(month, '-01') <= ?"
		args = append(args, str(m, "endTime"))
	}
	rows, err := s.pick(r).Query("select IFNULL(sum(sum_money), 0) sumMoney from chain_bill_rate "+searchSql, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil || len(ms) == 0 {
		platform.Unknown(w)
		return
	}
	platform.Success(w, map[string]any{"sumMoney": json.Number(btsStr(ms[0]["sumMoney"]))})
}

// rateCountMonBuses 补月;startTime/endTime null → NPE 999
func (s *Service) rateCountMonBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	where := "where b.company_id=e.id "
	var args []any
	if v, ok := m["startTime"]; ok && v != nil {
		where += " and concat(month, '-01') >= ?"
		args = append(args, str(m, "startTime"))
	} else {
		platform.Unknown(w)
		return
	}
	if v, ok := m["endTime"]; ok && v != nil {
		where += " and concat(month, '-01') <= ?"
		args = append(args, str(m, "endTime"))
	} else {
		platform.Unknown(w)
		return
	}
	if v, ok := m["companyId"]; ok && v != nil && !isEmptyStr(v) {
		where += " and b.company_id = ?"
		args = append(args, v)
	}
	if cond, arg, has := xzqhCond(m); has {
		where += cond
		args = append(args, arg)
		m["xzqhCode"] = "%" + str(m, "xzqhCode") + ",%"
	}
	where += " GROUP BY month order by month desc"
	rows, err := d.Query("select month, sum(b.sum_money) baseMoney from chain_bill_rate b, chain_basic_enterprise e "+where, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	res := map[string]any{}
	for _, row := range ms {
		res[btsStr(row["month"])] = json.Number(btsStr(row["baseMoney"]))
	}
	if !fillMonths(res, str(m, "startTime"), str(m, "endTime")) {
		platform.Unknown(w)
		return
	}
	platform.Success(w, res)
}
