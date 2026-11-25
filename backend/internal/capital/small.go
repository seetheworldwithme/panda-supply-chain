// capital - small.go 7 个小控制器标准套件:账单结算记录(relieve 特有)/资金账户/
// 银行回单(fill)/交易信息(fill)/核销记录/核销明细/明细结算对应。
// add/updateById 响应差异:Bkcall/Transinfo 恒 data=null(success() 无参),其余回显实体。
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

// ── 实体形状 ──

var entBillSettRec = entShape{
	table: "chain_bill_sett_rec",
	cols: []string{"id", "company_id", "sett_num", "bill_num", "operater", "type",
		"create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"},
	carriers: []string{"searchStr", "page", "ids", "relDate"},
}

var entAccounts = entShape{
	table: "chain_capital_accounts",
	cols: []string{"id", "company_id", "account_name", "account_number", "owner_type",
		"create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"},
	carriers: []string{"searchStr", "page", "ids"},
}

// entBkcall 招行回调记录(fill 实体:insert 自动 create+update 时间,update 自动 update_time);
// 列 msgText 驼峰,JSON 键 msgtext(renames 由 rows/echoKeys 处理)
var entBkcall = entShape{
	table: "chain_capital_bkcall",
	cols: []string{"id", "msg_nbr", "msg_type", "msg_key", "plt_nbr", "mch_nbr", "order_nbr",
		"order_sub_nbr", "mch_order_nbr", "mch_trans_nbr", "order_type", "order_subtype",
		"order_post_time", "order_amt", "trans_amt", "order_ccy_nbr", "status", "rtn_flag",
		"bck_flag", "bck_amt", "order_str_time", "order_end_time",
		"payer_acc_nbr", "payer_acc_name", "payer_bank_nbr", "payer_bank_name",
		"payee_acc_nbr", "payee_acc_name", "payee_bank_nbr", "payee_bank_name",
		"msgText", "create_time", "update_time",
		"remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"order_amt", "trans_amt", "bck_amt"},
	renames:  map[string]string{"msgText": "msgtext"}, // JSON 键全小写(lombok getter)
	carriers: []string{"searchStr", "page", "ids"},
	fillIns:  true,
}

// entTransinfo 交易记录(fill 实体);列 orderSubtype 驼峰,JSON 键 ordersubtype
var entTransinfo = entShape{
	table: "chain_capital_transinfo",
	cols: []string{"id", "system_code", "mch_order_nbr", "mch_trans_nbr", "bck_order_nbr",
		"bck_trans_nbr", "company_id", "pay_unit_id", "payee_unit_id", "pay_unit_name",
		"payee_unit_name", "orderSubtype", "contract_id", "contract_code",
		"pay_account_name", "pay_account", "pay_money", "payee_account_name",
		"payee_account", "mch_nbr", "status", "approve_status", "bus_status",
		"payer_bank_flag", "create_time", "update_time",
		"remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"pay_money"},
	renames:  map[string]string{"orderSubtype": "ordersubtype"}, // JSON 键全小写
	carriers: []string{"searchStr", "page", "ids", "transType", "startTime", "endTime"},
	fillIns:  true,
}

var entCount = entShape{
	table: "chain_capital_count",
	cols: []string{"id", "system_code", "company_id", "user_id", "user_name", "sment_num",
		"detail_num", "type", "money", "create_time", "update_time",
		"remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"money"},
	carriers: []string{"searchStr", "page", "ids"},
}

var entCountMx = entShape{
	table: "chain_capital_count_mx",
	cols: []string{"id", "company_id", "statement_id", "sett_money", "detail_id",
		"detail_money", "hx_code", "create_time", "update_time",
		"remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"sett_money", "detail_money"},
	carriers: []string{"searchStr", "page", "ids"},
}

var entDetailSettlement = entShape{
	table: "chain_capital_detail_settlement",
	cols: []string{"id", "company_id", "contract_id", "statement_id", "detail_id", "money",
		"hx_code", "create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"money"},
	carriers: []string{"searchStr", "page", "ids"},
}

// ── 分发 ──

func (s *Service) handleSmall(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	switch {
	case strings.HasPrefix(p, "/saas-capital/chainBillSettRec"):
		s.smallCtl(w, r, smallDef{name: "chainBillSettRec", ent: &entBillSettRec, sysTable: "chain_bill_sett_rec", noSysCol: true, relieve: true})
	case strings.HasPrefix(p, "/saas-capital/chainCapitalAccounts"):
		s.smallCtl(w, r, smallDef{name: "chainCapitalAccounts", ent: &entAccounts, sysTable: "chain_capital_accounts", noSysCol: true})
	case strings.HasPrefix(p, "/saas-capital/chainCapitalBkcall"):
		s.smallCtl(w, r, smallDef{name: "chainCapitalBkcall", ent: &entBkcall, sysTable: "chain_capital_bkcall", noSysCol: true, noEcho: true, fill: true})
	case strings.HasPrefix(p, "/saas-capital/chainCapitalTransinfo"):
		s.smallCtl(w, r, smallDef{name: "chainCapitalTransinfo", ent: &entTransinfo, sysTable: "chain_capital_transinfo", noEcho: true, fill: true, trans: true})
	case strings.HasPrefix(p, "/saas-capital/chainCapitalCount/"):
		s.smallCtl(w, r, smallDef{name: "chainCapitalCount", ent: &entCount, sysTable: "chain_capital_count"})
	case strings.HasPrefix(p, "/saas-capital/chainCapitalCountMx"):
		s.smallCtl(w, r, smallDef{name: "chainCapitalCountMx", ent: &entCountMx, sysTable: "chain_capital_count_mx", noSysCol: true, bindFail: true})
	case strings.HasPrefix(p, "/saas-capital/chainCapitalDetailSettlement"):
		s.smallCtl(w, r, smallDef{name: "chainCapitalDetailSettlement", ent: &entDetailSettlement, sysTable: "chain_capital_detail_settlement", noSysCol: true, rawSearch: true})
	default:
		platform.Unknown(w)
	}
}

type smallDef struct {
	name      string
	ent       *entShape
	sysTable  string // selectSystemCode 查询表
	noSysCol  bool   // 表无 system_code 列 → importExcel 恒 999
	noEcho    bool   // add/updateById 恒 data=null(Bkcall/Transinfo)
	fill      bool   // fill 实体 → Insert/UpdateByID
	trans     bool   // Transinfo:listByPage 特殊(jwt transType + getSql)
	relieve   bool   // BillSettRec:/relieve 端点
	bindFail  bool   // CountMx:mapper #{searchStr} BindingException 恒 999
	rawSearch bool   // DetailSettlement:%${searchStr}% 原生内插
}

// smallCtl 标准骨架 9+1 端点
func (s *Service) smallCtl(w http.ResponseWriter, r *http.Request, def smallDef) {
	p := r.URL.Path
	d := s.pick(r)
	switch {
	case strings.HasSuffix(p, "/list"):
		m := bodyMap(r)
		cond, ok := allEqCond(m, def.ent.cols)
		if !ok {
			platform.Unknown(w)
			return
		}
		rows, err := platform.ListWhere(d, def.ent.table, cond)
		if err != nil {
			platform.Unknown(w)
			return
		}
		if rows == nil {
			rows = []map[string]any{}
		}
		platform.Success(w, def.ent.rows(rows))
	case strings.HasSuffix(p, "/listByPage"):
		s.smallListByPage(w, r, d, def)
	case strings.HasSuffix(p, "/getById"):
		m := bodyMap(r)
		var row map[string]any
		if id, ok := int64Of(m, "id"); ok {
			var err error
			row, err = platform.GetByID(d, def.ent.table, id)
			if err != nil {
				platform.Unknown(w)
				return
			}
		}
		if row == nil {
			platform.Success(w, nil)
			return
		}
		platform.Success(w, def.ent.rows([]map[string]any{row})[0])
	case strings.HasSuffix(p, "/add"):
		s.smallAdd(w, r, d, def)
	case strings.HasSuffix(p, "/deleteById"):
		m := bodyMap(r)
		id, ok := int64Of(m, "id")
		if !ok {
			platform.Unknown(w)
			return
		}
		if err := platform.DeleteByID(d, def.ent.table, id); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, "删除成功，id为："+str(m, "id"))
	case strings.HasSuffix(p, "/updateById"):
		s.smallUpdateById(w, r, d, def)
	case strings.HasSuffix(p, "/importExcel"):
		// 链路第一步 selectSystemCode:无 system_code 列的表 SQL 直接报错 → 恒 999
		if def.noSysCol {
			platform.Unknown(w)
			return
		}
		s.smallImportExcel(w, r, d, def)
	case strings.HasSuffix(p, "/exportExcel"):
		s.smallExportExcel(w, r, d, def)
	case strings.HasSuffix(p, "/delBatch"):
		m := bodyMap(r)
		ids, ok := m["ids"].([]any)
		if !ok {
			// MP removeByIds:null/缺失 → 内部短路 0 行 → 200 "删除成功"(8765 实测)
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
		if err := platform.DeleteByIDs(d, def.ent.table, numIds); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, "删除成功")
	case strings.HasSuffix(p, "/relieve") && def.relieve:
		s.settRecRelieve(w, r, d)
	default:
		platform.Unknown(w)
	}
}

// smallListByPage 各小控制器 listByPage(实现各异)
func (s *Service) smallListByPage(w http.ResponseWriter, r *http.Request, d *sql.DB, def smallDef) {
	m := bodyMap(r)
	if _, ok := m["page"].(map[string]any); !ok {
		platform.Unknown(w) // Java mapper requires a Page object.
		return
	}
	current, size := pageOf(m)
	var q, countQ string
	var args []any
	var shape func(rows []map[string]any) []map[string]any
	switch def.name {
	case "chainBillSettRec":
		// relDate 等值 DATE_FORMAT(create_time);searchStr "%null%"/"%v%"
		where := "where 1=1 "
		if v, ok := m["relDate"]; ok && v != nil && strings.TrimSpace(str(m, "relDate")) != "" {
			t, ok2 := parseLocalDate(str(m, "relDate"))
			if !ok2 {
				platform.Unknown(w)
				return
			}
			where += " and DATE_FORMAT(create_time, '%Y-%m-%d') = ?"
			args = append(args, t.Format("2006-01-02"))
		}
		q = "select * from chain_bill_sett_rec " + where + " and CONCAT(sett_num, bill_num, operater) like ? order by create_time desc"
		countQ = "SELECT count(1) FROM chain_bill_sett_rec " + where + " and CONCAT(sett_num, bill_num, operater) like ?"
		args = append(args, searchStrLike(m))
		shape = def.ent.rows
	case "chainCapitalAccounts":
		// searchSql 基础串自带尾随 "and ":ownerType null 时为 "where 1=1 and  CONCAT..."(合法)
		where := "where 1=1 and "
		if v, has := m["ownerType"]; has && v != nil {
			where += "owner_type = ? and "
			args = append(args, v)
		}
		like := searchStrLike(m)
		q = "select * from chain_capital_accounts " + where + "CONCAT(IFNULL(account_name, ''), IFNULL(account_number, '')) like ? order by create_time desc"
		countQ = "SELECT count(1) FROM chain_capital_accounts " + where + "CONCAT(IFNULL(account_name, ''), IFNULL(account_number, '')) like ?"
		args = append(args, like)
		shape = def.ent.rows
	case "chainCapitalBkcall":
		// StringHelper.getSql 反射通用:全部业务字段 @ExcelColumn,值非 null 非 "0" → like
		sqlCond, margs, ok := getSqlReflect(m, def.ent, nil)
		if !ok {
			platform.Unknown(w)
			return
		}
		q = "select * from chain_capital_bkcall " + sqlCond + " and CONCAT(id) like ? order by create_time desc"
		countQ = "SELECT count(1) FROM chain_capital_bkcall " + sqlCond + " and CONCAT(id) like ?"
		args = append(margs, searchStrLike(m))
		shape = def.ent.rows
	case "chainCapitalTransinfo":
		// getSql + transType 分支(jwt) + startTime/endTime
		sqlCond, margs, ok := getSqlReflect(m, def.ent, nil)
		if !ok {
			platform.Unknown(w)
			return
		}
		u := currentUser(r)
		if str(m, "transType") == "1" {
			sqlCond += " and pay_unit_id = ?"
			margs = append(margs, json.Number(strconv.FormatInt(u.CompanyID, 10)))
		} else if str(m, "transType") == "2" {
			sqlCond += " and payee_unit_id = ?"
			margs = append(margs, json.Number(strconv.FormatInt(u.CompanyID, 10)))
		}
		if v, ok2 := m["startTime"]; ok2 && v != nil && strings.TrimSpace(str(m, "startTime")) != "" {
			t, ok3 := parseLocalDate(str(m, "startTime"))
			if !ok3 {
				platform.Unknown(w)
				return
			}
			sqlCond += " and DATE_FORMAT(create_time, '%Y-%m-%d') >= ?"
			margs = append(margs, t.Format("2006-01-02"))
		}
		if v, ok2 := m["endTime"]; ok2 && v != nil && strings.TrimSpace(str(m, "endTime")) != "" {
			t, ok3 := parseLocalDate(str(m, "endTime"))
			if !ok3 {
				platform.Unknown(w)
				return
			}
			sqlCond += " and DATE_FORMAT(create_time, '%Y-%m-%d') <= ?"
			margs = append(margs, t.Format("2006-01-02"))
		}
		q = "select * from chain_capital_transinfo " + sqlCond + " and CONCAT(IFNULL(mch_order_nbr, ''), IFNULL(mch_trans_nbr, ''), IFNULL(payee_unit_name, ''), IFNULL(contract_code, ''), IFNULL(pay_account_name, ''), IFNULL(pay_account, ''), IFNULL(pay_money, ''), IFNULL(payee_account_name, ''), IFNULL(payee_account, '')) like ? order by create_time desc"
		countQ = "SELECT count(1) FROM chain_capital_transinfo " + sqlCond + " and CONCAT(IFNULL(mch_order_nbr, ''), IFNULL(mch_trans_nbr, ''), IFNULL(payee_unit_name, ''), IFNULL(contract_code, ''), IFNULL(pay_account_name, ''), IFNULL(pay_account, ''), IFNULL(pay_money, ''), IFNULL(payee_account_name, ''), IFNULL(payee_account, '')) like ?"
		args = append(margs, searchStrLike(m))
		shape = def.ent.rows
	case "chainCapitalCount":
		where := "where 1=1 "
		if v, ok := m["createTime"]; ok && v != nil && strings.TrimSpace(str(m, "createTime")) != "" {
			// @JsonFormat 空格 pattern:只收 "yyyy-MM-dd HH:mm:ss"
			t, err := time.ParseInLocation("2006-01-02 15:04:05", str(m, "createTime"), time.Local)
			if err != nil {
				platform.Unknown(w) // 绑定失败 → 999
				return
			}
			where += " and create_time like ?"
			args = append(args, "%"+t.Format("2006-01-02")+"%")
		}
		q = "select * from chain_capital_count " + where + " and CONCAT(sment_num, detail_num, money) like ? order by create_time desc"
		countQ = "SELECT count(1) FROM chain_capital_count " + where + " and CONCAT(sment_num, detail_num, money) like ?"
		args = append(args, searchStrLike(m))
		shape = def.ent.rows
	case "chainCapitalCountMx":
		// mapper #{searchStr} 参数错配 → BindingException 恒 999
		platform.Unknown(w)
		return
	case "chainCapitalDetailSettlement":
		// %${searchStr}% 原生内插:null → 空串(四期实证 MyBatis ${null} 输出 "")
		ss := "%%"
		if v, ok := m["searchStr"]; ok && v != nil {
			ss = "%" + str(m, "searchStr") + "%"
		}
		q = "select * from chain_capital_detail_settlement where CONCAT(id) like ? order by create_time desc"
		countQ = "SELECT count(1) FROM chain_capital_detail_settlement where CONCAT(id) like ?"
		args = []any{ss} // 原先字面内插,可注入
		shape = def.ent.rows
	default:
		platform.Unknown(w)
		return
	}
	ip, err := paged(d, q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = shape(recs)
	platform.Success(w, ip)
}

// getSqlReflect 对照 StringHelper.getSql/getMaps:@ExcelColumn 字段值非 null 非 "0" →
// " and <snake> like ?"(值 Integer/LocalDate 原样,其余 %v%)。可排除字段 excluded。
func getSqlReflect(m map[string]any, ent *entShape, excluded map[string]bool) (string, []any, bool) {
	where := "where 1=1"
	var args []any
	for _, c := range ent.cols {
		k := camel(c)
		if excluded != nil && excluded[k] {
			continue
		}
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		if s, isStr := v.(string); isStr && s == "0" {
			continue
		}
		where += " and " + c + " like ?"
		switch t := v.(type) {
		case json.Number:
			args = append(args, t.String())
		case bool:
			args = append(args, "%"+str(m, k)+"%")
		default:
			args = append(args, "%"+str(m, k)+"%")
		}
	}
	return where, args, true
}

// smallAdd 标准骨架 add:Bkcall/Transinfo 恒 data=null,其余回显
func (s *Service) smallAdd(w http.ResponseWriter, r *http.Request, d *sql.DB, def smallDef) {
	m := bodyMap(r)
	vals := writeVals(def.ent, m)
	delete(vals, "id")
	delete(vals, "create_time")
	delete(vals, "update_time")
	var id int64
	var err error
	if def.fill {
		id, err = platform.Insert(d, def.ent.table, vals)
	} else {
		id, err = platform.InsertNoFill(d, def.ent.table, vals)
	}
	if err != nil {
		platform.Unknown(w)
		return
	}
	if def.noEcho {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, echoOf(def.ent, m, map[string]any{
		"id": json.Number(strconv.FormatInt(id, 10)),
	}))
}

// smallUpdateById:Bkcall/Transinfo data=null;fill 实体 update_time 自动
func (s *Service) smallUpdateById(w http.ResponseWriter, r *http.Request, d *sql.DB, def smallDef) {
	m := bodyMap(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	vals := writeVals(def.ent, m)
	delete(vals, "id")
	delete(vals, "create_time")
	delete(vals, "update_time")
	var err error
	if def.fill {
		_, err = platform.UpdateByID(d, def.ent.table, id, vals)
	} else {
		_, err = platform.UpdateByIDNoFill(d, def.ent.table, id, vals)
	}
	if err != nil {
		platform.Unknown(w)
		return
	}
	if def.noEcho {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, echoOf(def.ent, m, nil))
}

// smallImportExcel 标准模板(仅 Transinfo 有 system_code 列可走通;其余 noSysCol 已 999)
func (s *Service) smallImportExcel(w http.ResponseWriter, r *http.Request, d *sql.DB, def smallDef) {
	fh, ok := formFile(r)
	if !ok {
		irt := importResult{ErrInfos: []string{"上传文件为空"}}
		platform.FailData(w, 999, "失败", irt.ErrInfos)
		return
	}
	sysCode, ok := s.transSysCode(d, currentUser(r))
	if !ok {
		platform.Unknown(w)
		return
	}
	irt := readExcelFile(fh, def.ent.fields, sysCode, false, emptyDict())
	if len(irt.ErrInfos) > 0 {
		platform.FailData(w, 999, "失败", irt.ErrInfos)
		return
	}
	for _, row := range irt.DataList {
		vals := pick(row, def.ent.cols...)
		delete(vals, "id")
		delete(vals, "create_time")
		delete(vals, "update_time")
		if _, err := platform.Insert(d, def.ent.table, vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "导入成功，共"+strconv.Itoa(len(irt.DataList))+"条")
}

func (s *Service) transSysCode(d *sql.DB, u *auth.CurrentUser) (string, bool) {
	like := "JYXX." + u.EnterpriseCode + strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", "")
	return systemCodeNext(d, "chain_capital_transinfo", like, "1")
}

// smallExportExcel list() 全表 → writeExcel
func (s *Service) smallExportExcel(w http.ResponseWriter, r *http.Request, d *sql.DB, def smallDef) {
	s.exportHandler(w, r, def.ent.table, def.ent.fields)
}

// settRecRelieve 解绑:updSett/updBill(关系清空)+自删;ids null → NPE 999
func (s *Service) settRecRelieve(w http.ResponseWriter, r *http.Request, d *sql.DB) {
	m := bodyMap(r)
	ids, ok := m["ids"].([]any)
	if !ok {
		platform.Unknown(w)
		return
	}
	for _, it := range ids {
		id, ok := toID(it)
		if !ok {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("update chain_capital_settlement set relation_id=null, rel_bill_flag=false where relation_id = ?", id); err != nil {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("update chain_bill set relation_id=null, relation_flag=false where relation_id = ?", id); err != nil {
			platform.Unknown(w)
			return
		}
		if err := platform.DeleteByID(d, "chain_bill_sett_rec", id); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "成功")
}
