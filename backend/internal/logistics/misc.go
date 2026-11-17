// logistics - misc.go 运单/合同对应表 + 对接运单 + 过磅单三控制器。
// Java:ChainLogisticsContractController / ChainLogisticsThirdWaybillController /
// ChainLogisticsWeighingListController 及各自 ServiceImpl。
package logistics

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// ── 实体元数据 ──

// entLc chain_logistics_contract(无 fill;createTime @JsonFormat 空格,updateTime T)
var entLc = entShape{
	table: "chain_logistics_contract",
	cols: []string{"id", "company_id", "contract_id", "logistics_id", "create_time", "update_time",
		"remark1", "remark2", "remark3", "remark4", "remark5"},
	carriers: []string{"searchStr", "page", "ids"},
}

// entTw chain_logistics_third_waybill(无 fill;全部 datetime @JsonFormat 空格;clean_flag 布尔)
var entTw = entShape{
	table: "chain_logistics_third_waybill",
	cols: []string{"id", "batch_number", "credit_code", "order_name", "order_create_time", "order_number",
		"order_loading_point", "order_unloading_point", "order_start_time", "order_end_time", "bill_number",
		"trans_time", "wb_loading_address", "wb_unloading_address", "car_number", "subject_matter",
		"freight_unit", "freight_amount", "original_tonnage", "tonnage_paid", "pay_money", "carriage_platform",
		"loading_pro_lnglog", "unloading_pro_lnglog", "invoicing_status", "invoicing_amount",
		"invoicing_tonnage", "clean_flag", "create_time", "update_time", "remark1", "remark2", "remark3",
		"remark4", "remark5"},
	decimals: []string{"freight_unit", "freight_amount", "original_tonnage", "tonnage_paid", "pay_money",
		"invoicing_amount", "invoicing_tonnage"},
	bools:    []string{"clean_flag"},
	carriers: []string{"searchStr", "page", "ids"},
}

// entWl chain_logistics_weighing_list(fill:createTime INSERT + updateTime INSERT_UPDATE;
// Invoicing_time 列驼峰成 InvoicingTime,实体键 invoicingTime 需改写)
var entWl = entShape{
	table: "chain_logistics_weighing_list",
	cols: []string{"id", "im_batch_number", "company_id", "order_id", "system_code", "bill_number",
		"issuing_office", "raw_coal_supplier", "sstd_company", "license_plate", "in_time", "ss_in_time",
		"contract_no", "purchaser", "coal_type", "tare", "ss_tare", "tare_time", "rough_weight",
		"ss_rough_weight", "rough_weight_time", "coal_loading_driver", "coal_loading_time", "net_weight",
		"ss_net_weight", "Invoicing_time", "sstd_time", "upload_time", "ss_status", "print_status",
		"time_interval_minute", "add_or_not", "repatriator", "weigher", "ss_repatriator", "ss_weigher",
		"reason_cancellation", "reason_supplementary", "create_time", "update_time",
		"remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"tare", "ss_tare", "rough_weight", "ss_rough_weight", "net_weight", "ss_net_weight"},
	tTimes: []string{"in_time", "ss_in_time", "tare_time", "rough_weight_time", "coal_loading_time",
		"Invoicing_time", "sstd_time", "upload_time"},
	renames:    map[string]string{"InvoicingTime": "invoicingTime"},
	carriers:   []string{"searchStr", "page", "ids", "list", "filelist", "notBindOrder"},
	fillInsert: true,
}

var lcSpaceDT = []string{"create_time"}
var lcTDT = []string{"update_time"}
var twSpaceDT = []string{"order_create_time", "trans_time", "create_time", "update_time"}
var wlSpaceDT = []string{"create_time", "update_time"}
var wlTDT = []string{"in_time", "ss_in_time", "tare_time", "rough_weight_time",
	"coal_loading_time", "Invoicing_time", "sstd_time", "upload_time"}

// ── /chainLogisticsContract ──

func (s *Service) lcList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, entLc.cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_logistics_contract", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, entLc.rows(rows))
}

// lcListByPage where CONCAT(id) like '%${searchStr}%' order by create_time desc
func (s *Service) lcListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	current, size := pageOf(m)
	where := "WHERE CONCAT(id) like '%" + searchStrOf(m) + "%'"
	ip, err := paged(d, "SELECT * FROM chain_logistics_contract "+where+" ORDER BY create_time DESC",
		"SELECT count(1) FROM chain_logistics_contract "+where, nil, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entLc.rows(recs)
	platform.Success(w, ip)
}

func (s *Service) lcGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	var row map[string]any
	if id, ok := int64Of(m, "id"); ok {
		var err error
		row, err = platform.GetByID(s.pick(r), "chain_logistics_contract", id)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, entLc.rows([]map[string]any{row})[0])
}

// lcAdd save(无 fill):回显内存实体
func (s *Service) lcAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, lcSpaceDT, lcTDT) {
		platform.Unknown(w)
		return
	}
	vals := pick(m, entLc.cols...)
	delete(vals, "id")
	id, err := platform.InsertNoFill(s.pick(r), "chain_logistics_contract", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, platform.EntityEcho(m, entLc.echoKeys(), map[string]any{
		"id": json.Number(strconv.FormatInt(id, 10)),
	}))
}

func (s *Service) lcDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	id, ok := int64Of(m, "id")
	if ok {
		if err := platform.DeleteByID(d, "chain_logistics_contract", id); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, "删除成功，id为："+strconv.FormatInt(id, 10))
		return
	}
	if _, err := d.Exec("DELETE FROM chain_logistics_contract WHERE id = NULL"); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为：null")
}

// lcUpdateById updateById(无 fill):恒回显
func (s *Service) lcUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, lcSpaceDT, lcTDT) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	vals := pick(m, entLc.cols...)
	delete(vals, "id")
	if _, err := platform.UpdateByIDNoFill(d, "chain_logistics_contract", id, vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, platform.EntityEcho(m, entLc.echoKeys(), nil))
}

func (s *Service) lcDelBatch(w http.ResponseWriter, r *http.Request) {
	s.delBatchByIDs(w, r, entLc.table)
}

// delBatchByIDs removeByIds 语义:null/空集合 → MP 短路返回 0 行,照 8765 实测 200 "删除成功"
func (s *Service) delBatchByIDs(w http.ResponseWriter, r *http.Request, table string) {
	m := bodyMap(r)
	ids, ok := m["ids"].([]any)
	if !ok || len(ids) == 0 {
		platform.Success(w, "删除成功")
		return
	}
	numIds := make([]any, 0, len(ids))
	for _, it := range ids {
		n, ok := it.(json.Number)
		if !ok {
			platform.Unknown(w)
			return
		}
		numIds = append(numIds, n.String())
	}
	if err := platform.DeleteByIDs(s.pick(r), table, numIds); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// ── /chainLogisticsThirdWaybill ──

func (s *Service) twList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, entTw.cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_logistics_third_waybill", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, entTw.rows(rows))
}

// twGetSql 对照 StringHelper.getSql(t,"cleanFlag"):@ExcelColumn 字段声明序,
// 非 Integer/LocalDate → %v%,cleanFlag 由 unstr 排除
var twSqlFields = []string{"id", "batchNumber", "creditCode", "orderName", "orderCreateTime",
	"orderNumber", "orderLoadingPoint", "orderUnloadingPoint", "orderStartTime", "orderEndTime",
	"billNumber", "transTime", "wbLoadingAddress", "wbUnloadingAddress", "carNumber", "subjectMatter",
	"freightUnit", "freightAmount", "originalTonnage", "tonnagePaid", "payMoney", "carriagePlatform",
	"loadingProLnglog", "unloadingProLnglog", "invoicingStatus", "invoicingAmount", "invoicingTonnage",
	"cleanFlag", "createTime", "updateTime", "remark1", "remark2", "remark3", "remark4", "remark5"}

var twIntFields = map[string]bool{"id": true, "carriagePlatform": true, "invoicingStatus": true}

func (s *Service) twListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	u := currentUser(r)
	// jwt 公司信用代码前置(实体缺行 → Java NPE → 999)
	var rate sql.NullString
	if err := d.QueryRow("SELECT rate_number FROM chain_basic_enterprise WHERE id = ?", u.CompanyID).Scan(&rate); err != nil || !rate.Valid {
		platform.Unknown(w) // etc.getRateNumber() NPE / selectById 无行
		return
	}
	m["creditCode"] = rate.String

	cond := "where 1=1"
	var args []any
	for _, f := range twSqlFields {
		if f == "cleanFlag" { // unstr 排除
			continue
		}
		v, ok := m[f]
		if !ok || v == nil || str(m, f) == "0" {
			continue
		}
		cond += " and " + snake(f) + " like ?"
		if twIntFields[f] {
			args = append(args, v)
		} else {
			args = append(args, "%"+str(m, f)+"%")
		}
	}
	// cleanFlag 手工分支(值原样绑定,Java like #{maps.cleanFlag})
	if v, ok := m["cleanFlag"]; ok && v != nil {
		cond += " and clean_flag like ?"
		args = append(args, v)
	}
	// CONCAT 全字段 like #{maps.searchStr}:getMaps 包 %;searchStr 缺失 → 参数 NULL → 恒不匹配
	if v, ok := m["searchStr"]; ok && v != nil {
		cond += " and CONCAT(IFNULL(batch_number, ''), IFNULL(credit_code, ''), IFNULL(order_name, ''), IFNULL(order_create_time, ''), IFNULL(order_number, ''), IFNULL(order_loading_point, ''), IFNULL(order_unloading_point, ''), IFNULL(order_start_time, ''), IFNULL(order_end_time, ''), IFNULL(bill_number, ''), IFNULL(trans_time, ''), IFNULL(wb_loading_address, ''), IFNULL(wb_unloading_address, ''), IFNULL(car_number, ''), IFNULL(subject_matter, ''), IFNULL(freight_unit, ''), IFNULL(freight_amount, ''), IFNULL(original_tonnage, ''), IFNULL(tonnage_paid, ''), IFNULL(pay_money, ''), IFNULL(loading_pro_lnglog, ''), IFNULL(unloading_pro_lnglog, ''), IFNULL(invoicing_amount, ''), IFNULL(invoicing_tonnage, '') ) like ?"
		args = append(args, "%"+str(m, "searchStr")+"%")
	} else {
		cond += " and CONCAT(IFNULL(batch_number, ''), IFNULL(credit_code, ''), IFNULL(order_name, ''), IFNULL(order_create_time, ''), IFNULL(order_number, ''), IFNULL(order_loading_point, ''), IFNULL(order_unloading_point, ''), IFNULL(order_start_time, ''), IFNULL(order_end_time, ''), IFNULL(bill_number, ''), IFNULL(trans_time, ''), IFNULL(wb_loading_address, ''), IFNULL(wb_unloading_address, ''), IFNULL(car_number, ''), IFNULL(subject_matter, ''), IFNULL(freight_unit, ''), IFNULL(freight_amount, ''), IFNULL(original_tonnage, ''), IFNULL(tonnage_paid, ''), IFNULL(pay_money, ''), IFNULL(loading_pro_lnglog, ''), IFNULL(unloading_pro_lnglog, ''), IFNULL(invoicing_amount, ''), IFNULL(invoicing_tonnage, '') ) like ?"
		args = append(args, nil)
	}
	current, size := pageOf(m)
	ip, err := paged(d, "SELECT * FROM chain_logistics_third_waybill "+cond+" ORDER BY trans_time DESC",
		"SELECT count(1) FROM chain_logistics_third_waybill "+cond, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entTw.rows(recs)
	platform.Success(w, ip)
}

func (s *Service) twGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	var row map[string]any
	if id, ok := int64Of(m, "id"); ok {
		var err error
		row, err = platform.GetByID(s.pick(r), "chain_logistics_third_waybill", id)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, entTw.rows([]map[string]any{row})[0])
}

// twAdd save 后 success() 不带回显(与运单/订单控制器不同)
func (s *Service) twAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, twSpaceDT, nil) {
		platform.Unknown(w)
		return
	}
	vals := pick(m, entTw.cols...)
	delete(vals, "id")
	if _, err := platform.InsertNoFill(s.pick(r), "chain_logistics_third_waybill", vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// twAddBatch 批量新增:运单号重复(HashSet 全字段 equals)→ fail;批次号已存在 → fail(带条数);
// BatchCount 头缺失/不符 → fail;成功 saveBatch 后 success()
func (s *Service) twAddBatch(w http.ResponseWriter, r *http.Request) {
	items, _ := bodyList(r)
	d := s.pick(r)
	if hasDupEntities(items) {
		platform.FailMsg(w, "当前批次运单号重复，请核查")
		return
	}
	batchNumber := r.Header.Get("BatchNumber")
	batchCount := r.Header.Get("BatchCount")
	var cnt int64
	q := "SELECT count(1) FROM chain_logistics_third_waybill WHERE "
	var arg any
	if batchNumber == "" {
		q += "batch_number IS NULL" // allEq 单键 null → IS NULL
		arg = nil
	} else {
		q += "batch_number = ?"
		arg = batchNumber
	}
	if err := d.QueryRow(q, arg).Scan(&cnt); err != nil && arg != nil {
		platform.Unknown(w)
		return
	}
	if cnt > 0 {
		platform.FailMsg(w, "当前批次号重复，请核查，该批次已有数据"+strconv.FormatInt(cnt, 10)+"条")
		return
	}
	if batchCount == "" {
		platform.FailMsg(w, "当前批次预计数量不能为空")
		return
	}
	bc, err := strconv.Atoi(batchCount)
	if err != nil {
		platform.Unknown(w) // Integer.parseInt NumberFormatException
		return
	}
	if bc != len(items) {
		platform.FailMsg(w, "当前批次预计数量与传输数量不一致，应收"+batchCount+"条，实收"+strconv.Itoa(len(items))+"条")
		return
	}
	for _, it := range items {
		e, ok := it.(map[string]any)
		if !ok {
			continue
		}
		if !validateDT(e, twSpaceDT, nil) {
			platform.Unknown(w)
			return
		}
		e["batchNumber"] = batchNumber
		vals := pick(e, entTw.cols...)
		delete(vals, "id")
		if _, err := platform.InsertNoFill(d, "chain_logistics_third_waybill", vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, nil)
}

// hasDupEntities 对照 new HashSet<>(list).size() != list.size():全字段等价去重
func hasDupEntities(items []any) bool {
	seen := map[string]bool{}
	for _, it := range items {
		b, err := json.Marshal(it)
		if err != nil {
			continue
		}
		if seen[string(b)] {
			return true
		}
		seen[string(b)] = true
	}
	return false
}

func (s *Service) twDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	id, ok := int64Of(m, "id")
	if ok {
		if err := platform.DeleteByID(d, "chain_logistics_third_waybill", id); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, "删除成功，id为："+strconv.FormatInt(id, 10))
		return
	}
	if _, err := d.Exec("DELETE FROM chain_logistics_third_waybill WHERE id = NULL"); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为：null")
}

// twUpdateById updateById(无 fill),success() 不带回显
func (s *Service) twUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, twSpaceDT, nil) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	vals := pick(m, entTw.cols...)
	delete(vals, "id")
	if _, err := platform.UpdateByIDNoFill(d, "chain_logistics_third_waybill", id, vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// twUpdateBatch 批量修改:查重 → BatchCount 头校验 → 逐条按 bill_number 定位
// (有行取其 id,无行先 insert)→ updateBatchById(无 fill)
func (s *Service) twUpdateBatch(w http.ResponseWriter, r *http.Request) {
	items, _ := bodyList(r)
	d := s.pick(r)
	if hasDupEntities(items) {
		platform.FailMsg(w, "当前批次运单号重复，请核查")
		return
	}
	batchCount := r.Header.Get("BatchCount")
	if batchCount == "" {
		platform.FailMsg(w, "预计数量不能为空")
		return
	}
	bc, err := strconv.Atoi(batchCount)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if bc != len(items) {
		platform.FailMsg(w, "预计数量与传输数量不一致")
		return
	}
	for _, it := range items {
		e, ok := it.(map[string]any)
		if !ok {
			continue
		}
		if !validateDT(e, twSpaceDT, nil) {
			platform.Unknown(w)
			return
		}
		var id sql.NullInt64
		_ = d.QueryRow("SELECT id FROM chain_logistics_third_waybill WHERE bill_number LIKE ?", str(e, "billNumber")).Scan(&id)
		if !id.Valid {
			vals := pick(e, entTw.cols...)
			delete(vals, "id")
			newID, err := platform.InsertNoFill(d, "chain_logistics_third_waybill", vals)
			if err != nil {
				platform.Unknown(w)
				return
			}
			e["id"] = json.Number(strconv.FormatInt(newID, 10))
		} else {
			e["id"] = json.Number(strconv.FormatInt(id.Int64, 10))
		}
		vals := pick(e, entTw.cols...)
		delete(vals, "id")
		if len(vals) == 0 {
			continue
		}
		if _, err := platform.UpdateByIDNoFill(d, "chain_logistics_third_waybill", e["id"], vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, nil)
}

func (s *Service) twDelBatch(w http.ResponseWriter, r *http.Request) {
	s.delBatchByIDs(w, r, entTw.table)
}

// twClearData 手动数据清洗:按信用代码+clean_flag=0(+body 条件)取数,
// 逐条生成/更新订单(WLDD.)与运单(WLYD.)、标记已清洗,再按订单号汇总更新订单统计。
// 金额字段 NULL → BigDecimal.add NPE → 999;订单号重复多行 getOne → TooManyResults → 999
func (s *Service) twClearData(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	u := currentUser(r)
	companyID := u.CompanyID

	// allEq(params)(null2IsNull=true):body 键原样作列名,null → IS NULL
	cond := "WHERE clean_flag = 0"
	var args []any
	badCol := false
	for k, v := range m {
		if !validCol(k, entTw.cols) {
			badCol = true
			break
		}
		if v == nil {
			cond += " AND " + k + " IS NULL"
			continue
		}
		cond += " AND " + k + " = ?"
		args = append(args, v)
	}
	if badCol {
		platform.Unknown(w) // 非法列 → SQL 异常
		return
	}
	// 企业行存在才带 credit_code(值为 null → IS NULL 条件,allEq null2IsNull 语义)
	var etcRow sql.NullString
	if err := d.QueryRow("SELECT rate_number FROM chain_basic_enterprise WHERE id = ?", companyID).Scan(&etcRow); err == nil {
		if etcRow.Valid && etcRow.String != "" {
			cond += " AND credit_code = ?"
			args = append(args, etcRow.String)
		} else {
			cond += " AND credit_code IS NULL"
		}
	}
	rows, err := d.Query("SELECT * FROM chain_logistics_third_waybill "+cond, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	bills, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}

	etpCode := ""
	if len(bills) > 0 {
		var etcID sql.NullInt64
		_ = d.QueryRow("SELECT id FROM chain_basic_enterprise WHERE rate_number = ?", btsStr(bills[0]["credit_code"])).Scan(&etcID)
		if etcID.Valid {
			companyID = etcID.Int64
			var code sql.NullString
			_ = d.QueryRow("SELECT enterprise_code FROM chain_basic_enterprise WHERE id = ?", etcID.Int64).Scan(&code)
			etpCode = code.String
		}
	}
	today := time.Now().Format("20060102")
	bscode, ok1 := nextSystemCode(d, "chain_logistics_bill", "WLDD."+etpCode+today, "1")
	lcode, ok2 := nextSystemCode(d, "chain_logistics", "WLYD."+etpCode+today, "1")
	if !ok1 || !ok2 {
		platform.Unknown(w)
		return
	}

	type agg struct {
		vehicles      int64
		freightAmount float64
		originalTon   float64
		tonnagePaid   float64
	}
	bmap := map[string]*agg{}
	var order []string

	for _, bl := range bills {
		orderNumber := btsStr(bl["order_number"])
		// 订单:getOne(order_number) 多行 → 999
		brows, err := d.Query("SELECT * FROM chain_logistics_bill WHERE order_number = ?", orderNumber)
		if err != nil {
			platform.Unknown(w)
			return
		}
		lobs, _ := platform.ScanMaps(brows)
		if len(lobs) > 1 {
			platform.Unknown(w) // getOne TooManyResultsException
			return
		}
		billVals := map[string]any{
			"batch_number":      nilIfEmpty(btsStr(bl["batch_number"])),
			"company_id":        json.Number(strconv.FormatInt(companyID, 10)),
			"order_name":        nilIfEmpty(btsStr(bl["order_name"])),
			"order_create_time": tFormatToSpace(btsStr(bl["order_create_time"])),
			"system_code":       bscode,
			"order_number":      orderNumber,
			"loading_address":   nilIfEmpty(btsStr(bl["order_loading_point"])),
			"unloading_point":   nilIfEmpty(btsStr(bl["order_unloading_point"])),
			"order_start_time":  nilIfEmpty(btsStr(bl["order_start_time"])),
			"order_end_time":    nilIfEmpty(btsStr(bl["order_end_time"])),
			"trans_method":      1,
			"carriage_platform": intPtr(bl["carriage_platform"]),
			"trans_type":        1,
			"carriage_name":     2,
			"update_time":       nowSpace(), // bill 仅 updateTime fill → 插/更都自动 update_time
		}
		var lbID int64
		if len(lobs) == 0 {
			newID, err := platform.InsertNoFill(d, "chain_logistics_bill", billVals)
			if err != nil {
				platform.Unknown(w)
				return
			}
			lbID = newID
		} else {
			lbID = lobs[0]["id"].(int64)
			if _, err := platform.UpdateByID(d, "chain_logistics_bill", lbID, dropNil(billVals)); err != nil {
				platform.Unknown(w)
				return
			}
		}
		// 运单
		lgVals := map[string]any{
			"batch_number":      nilIfEmpty(btsStr(bl["batch_number"])),
			"company_id":        json.Number(strconv.FormatInt(companyID, 10)),
			"order_id":          json.Number(strconv.FormatInt(lbID, 10)),
			"system_code":       lcode,
			"bill_number":       nilIfEmpty(btsStr(bl["bill_number"])),
			"trans_time":        tFormatToSpace(btsStr(bl["trans_time"])),
			"car_number":        nilIfEmpty(btsStr(bl["car_number"])),
			"loading_address":   nilIfEmpty(btsStr(bl["wb_loading_address"])),
			"unloading_point":   nilIfEmpty(btsStr(bl["wb_unloading_address"])),
			"subject_matter":    nilIfEmpty(btsStr(bl["subject_matter"])),
			"freight_unit":      decOrNull(bl["freight_unit"]),
			"freight_amount":    decOrNull(bl["freight_amount"]),
			"original_tonnage":  decOrNull(bl["original_tonnage"]),
			"tonnage_paid":      decOrNull(bl["tonnage_paid"]),
			"pay_money":         decOrNull(bl["pay_money"]),
			"carriage_platform": intPtr(bl["carriage_platform"]),
			"lo_create_time":    nowSpace(),
			"loading_lng":       nilIfEmpty(btsStr(bl["loading_pro_lnglog"])),
			"unloading_lng":     nilIfEmpty(btsStr(bl["unloading_pro_lnglog"])),
			"invoicing_status":  intPtr(bl["invoicing_status"]),
			"invoicing_amount":  decOrNull(bl["invoicing_amount"]),
			"invoicing_tonnage": decOrNull(bl["invoicing_tonnage"]),
		}
		{
			// 运单按 bill_number 查(service 用 list(allEq)取首条)
			lrows, err := d.Query("SELECT id FROM chain_logistics WHERE bill_number = ? LIMIT 2", btsStr(bl["bill_number"]))
			if err != nil {
				platform.Unknown(w)
				return
			}
			loms, _ := platform.ScanMaps(lrows)
			if len(loms) == 0 {
				if _, err := platform.Insert(d, "chain_logistics", dropNil(lgVals)); err != nil {
					platform.Unknown(w)
					return
				}
			} else {
				if _, err := platform.UpdateByID(d, "chain_logistics", loms[0]["id"], dropNil(lgVals)); err != nil {
					platform.Unknown(w)
					return
				}
			}
		}
		// 标记已清洗(无 fill,updateById 全量非空列回写 → 等价 clean_flag 置位)
		if _, err := d.Exec("UPDATE chain_logistics_third_waybill SET clean_flag = 1 WHERE id = ?", bl["id"]); err != nil {
			platform.Unknown(w)
			return
		}
		next1, okb := bumpSystemCode(bscode)
		next2, okc := bumpSystemCode(lcode)
		if !okb || !okc {
			platform.Unknown(w)
			return
		}
		bscode, lcode = next1, next2
		// 汇总(金额 NULL → Java NPE → 999)
		fa, okf := bl["freight_amount"]
		ot, oko := bl["original_tonnage"]
		tp, okt := bl["tonnage_paid"]
		if !okf || !oko || !okt || fa == nil || ot == nil || tp == nil {
			platform.Unknown(w)
			return
		}
		a := bmap[orderNumber]
		if a == nil {
			a = &agg{}
			bmap[orderNumber] = a
			order = append(order, orderNumber)
		}
		a.vehicles++
		a.freightAmount += num(fa)
		a.originalTon += num(ot)
		a.tonnagePaid += num(tp)
	}
	for _, on := range order {
		a := bmap[on]
		_, err := d.Exec("UPDATE chain_logistics_bill SET vehicles=?, freight_amount=?, original_tonnage=?, tonnage_paid=?, update_time=? WHERE order_number=?",
			a.vehicles, dec2(a.freightAmount), dec2(a.originalTon), dec2(a.tonnagePaid), nowSpace(), on)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "清理数据完毕，条数为："+strconv.Itoa(len(bills)))
}

// ── /chainLogisticsWeighingList ──

func (s *Service) wlList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, entWl.cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_logistics_weighing_list", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, entWl.rows(rows))
}

// wlSqlFields @ExcelColumn 声明序(unstr="orderId")
var wlSqlFields = []string{"id", "imBatchNumber", "companyId", "systemCode", "billNumber",
	"issuingOffice", "rawCoalSupplier", "sstdCompany", "licensePlate", "inTime", "ssInTime",
	"contractNo", "purchaser", "coalType", "tare", "ssTare", "tareTime", "roughWeight", "ssRoughWeight",
	"roughWeightTime", "coalLoadingDriver", "coalLoadingTime", "netWeight", "ssNetWeight",
	"invoicingTime", "sstdTime", "uploadTime", "ssStatus", "printStatus", "timeIntervalMinute",
	"addOrNot", "repatriator", "weigher", "ssRepatriator", "ssWeigher", "reasonCancellation",
	"reasonSupplementary", "createTime", "updateTime", "remark1", "remark2", "remark3", "remark4", "remark5"}

var wlRawFields = map[string]bool{"id": true, "companyId": true, "timeIntervalMinute": true,
	"inTime": true, "ssInTime": true, "tareTime": true, "roughWeightTime": true, "coalLoadingTime": true,
	"invoicingTime": true, "sstdTime": true, "uploadTime": true}

const wlConcat = "CONCAT(IFNULL(bill_number, ''), IFNULL(issuing_office, ''), IFNULL(raw_coal_supplier, ''), IFNULL(sstd_company, ''), IFNULL(license_plate, ''), IFNULL(in_time, ''), IFNULL(ss_in_time, ''), IFNULL(contract_no, ''), IFNULL(purchaser, ''), IFNULL(coal_type, ''), IFNULL(tare, ''), IFNULL(ss_tare, ''), IFNULL(tare_time, ''), IFNULL(rough_weight, ''), IFNULL(ss_rough_weight, ''), IFNULL(rough_weight_time, ''), IFNULL(coal_loading_driver, ''), IFNULL(coal_loading_time, ''), IFNULL(net_weight, ''), IFNULL(ss_net_weight, ''), IFNULL(Invoicing_time, ''), IFNULL(sstd_time, ''), IFNULL(upload_time, ''), IFNULL(ss_status, ''), IFNULL(print_status, ''), IFNULL(repatriator, ''), IFNULL(weigher, ''), IFNULL(ss_repatriator, ''), IFNULL(ss_weigher, ''))"

func (s *Service) wlListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	cond := "where 1=1"
	var args []any
	for _, f := range wlSqlFields {
		if f == "orderId" { // unstr 排除
			continue
		}
		v, ok := m[f]
		if !ok || v == nil || str(m, f) == "0" {
			continue
		}
		cond += " and " + snake(f) + " like ?"
		if wlRawFields[f] {
			args = append(args, v)
		} else {
			args = append(args, "%"+str(m, f)+"%")
		}
	}
	// notBindOrder/orderId 手工分支(值原样内插)
	if v, ok := m["notBindOrder"]; ok && v != nil {
		cond += " and (order_id is null or order_id='' or order_id=" + str(m, "orderId") + ")"
	} else if v, ok := m["orderId"]; ok && v != nil {
		cond += " and order_id =" + str(m, "orderId")
	}
	if v, ok := m["searchStr"]; ok && v != nil {
		cond += " and " + wlConcat + " like ?"
		args = append(args, "%"+str(m, "searchStr")+"%")
	} else {
		cond += " and " + wlConcat + " like ?"
		args = append(args, nil)
	}
	current, size := pageOf(m)
	ip, err := paged(d, "SELECT * FROM chain_logistics_weighing_list "+cond+" ORDER BY create_time DESC",
		"SELECT count(1) FROM chain_logistics_weighing_list "+cond, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entWl.rows(recs)
	platform.Success(w, ip)
}

func (s *Service) wlGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	var row map[string]any
	if id, ok := int64Of(m, "id"); ok {
		var err error
		row, err = platform.GetByID(s.pick(r), "chain_logistics_weighing_list", id)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, entWl.rows([]map[string]any{row})[0])
}

// wlAdd save(fill)+filelist,success() 不带回显
func (s *Service) wlAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, wlSpaceDT, wlTDT) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	vals := pick(m, entWl.cols...)
	delete(vals, "id")
	id, err := platform.Insert(d, "chain_logistics_weighing_list", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	s.updateFilelist(d, m, id, u)
	platform.Success(w, nil)
}

func (s *Service) wlDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	id, ok := int64Of(m, "id")
	if ok {
		if err := platform.DeleteByID(d, "chain_logistics_weighing_list", id); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, "删除成功，id为："+strconv.FormatInt(id, 10))
		return
	}
	if _, err := d.Exec("DELETE FROM chain_logistics_weighing_list WHERE id = NULL"); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为：null")
}

// wlUpdateById updateById(fill:update_time),success() 不带回显
func (s *Service) wlUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, wlSpaceDT, wlTDT) {
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
	vals := pick(m, entWl.cols...)
	delete(vals, "id")
	if _, err := platform.UpdateByID(d, "chain_logistics_weighing_list", id, vals); err != nil {
		platform.Unknown(w)
		return
	}
	s.updateFilelist(d, m, id, u)
	platform.Success(w, nil)
}

func (s *Service) wlDelBatch(w http.ResponseWriter, r *http.Request) {
	s.delBatchByIDs(w, r, entWl.table)
}

// wlSaveBatch 批量新增过磅单(WLGBD. 编号顺序赋值),响应固定 "新增成功"
func (s *Service) wlSaveBatch(w http.ResponseWriter, r *http.Request) {
	s.saveBatchWithCode(w, r, "chain_logistics_weighing_list", entWl, "WLGBD.")
}

// ── 小工具 ──

func snake(s string) string {
	var b strings.Builder
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(c - 'A' + 'a')
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func validCol(k string, cols []string) bool {
	for _, c := range cols {
		if strings.EqualFold(c, k) {
			return true
		}
	}
	return false
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func intPtr(v any) any {
	if v == nil {
		return nil
	}
	return json.Number(btsStr(v))
}

func decOrNull(v any) any {
	if v == nil {
		return nil
	}
	return json.Number(btsStr(v))
}

func dec2(f float64) string { return strconv.FormatFloat(f, 'f', 2, 64) }

// tFormatToSpace ISO T 行值转空格格式(入库;空串保持空)
func tFormatToSpace(s string) any {
	if s == "" {
		return nil
	}
	if len(s) == 19 && s[10] == 'T' {
		return s[:10] + " " + s[11:]
	}
	return s
}

func dropNil(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if v != nil {
			out[k] = v
		}
	}
	return out
}

// ── /saas-common/chainCommonTempmx/formatlist(物流页基础维护字典依赖) ──
// Java:ChainCommonTempmxController.formatlist → alllist(template_id, jwt companyId),
// 无行回退 company 0,按 d_type 分组(HashMap);实体 createTime @JsonFormat 空格
var tempmxCols = []string{"id", "company_id", "template_id", "dict_id", "d_parent_id", "d_type",
	"d_type_name", "child_key", "d_key", "d_value", "d_sel_type", "d_base_default",
	"d_template_default", "d_extend", "create_time", "update_time",
	"remark1", "remark2", "remark3", "remark4", "remark5"}
var tempmxCarriers = []string{"searchStr", "page", "ids"}

func (s *Service) tempmxFormatlist(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	u := currentUser(r)
	tid := str(m, "template_id")
	if _, err := strconv.Atoi(tid); err != nil {
		platform.Unknown(w) // Integer.parseInt(params.get("template_id"))
		return
	}
	cid := strconv.FormatInt(u.CompanyID, 10)
	rows, err := platform.ListWhere(d, "chain_common_tempmx", map[string]any{"template_id": tid, "company_id": cid})
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(rows) == 0 {
		rows, err = platform.ListWhere(d, "chain_common_tempmx", map[string]any{"template_id": tid, "company_id": "0"})
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	out := map[string][]map[string]any{}
	for _, r2 := range rows {
		dt := btsStr(r2["d_type"])
		out[dt] = append(out[dt], r2)
	}
	// 行整形:驼峰+布尔+时间格式+载体键
	res := map[string]any{}
	for dt, list := range out {
		shaped := entShape{
			cols: tempmxCols, carriers: tempmxCarriers,
			bools:  []string{"d_base_default", "d_template_default", "d_extend"},
			tTimes: []string{"update_time"},
		}.rows(list)
		for _, rr := range shaped {
			if s2, ok := rr["createTime"].(string); ok && len(s2) == 19 {
				// createTime 保持空格(默认即空格),update_time 已转 T
			}
		}
		res[dt] = shaped
	}
	platform.Success(w, res)
}
