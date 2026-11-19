// stats.go 合同列表(listByPage 大 searchSql 组装)与统计/四流合一查询。
// 逐条对照 ChainContractServiceImpl + ChainContractMapper 注解 SQL(行号见各 handler 注释)
package contract

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// contractExcelCols chain_contract 实体 @ExcelColumn 字段(getSql/getMaps 范围;DB 字段全集)
var contractExcelCols = []string{"id", "imBatchNumber", "companyId", "orderId", "templateId",
	"systemCode", "contractCode", "contractResource", "contractType", "signingTime",
	"contractName", "salesUnit", "buyerUnit", "salesUnitId", "buyerUnitId", "receiptPlace",
	"deliveryPlace", "collaborator", "proBigType", "proSmaType", "proRemark", "contractNum",
	"unitPrice", "taxIncluded", "sumPrice", "sulfur", "volatilization", "water", "remark",
	"priceType", "payMethod", "payType", "carrier", "transportMode", "contractIds",
	"agreementIds", "supplierId", "custId", "month", "custType", "isPayOnline", "isEsign",
	"isEsignComplete", "fileStatus", "sysFlowId", "eqbFlowId", "isPay", "isPayerConfirms",
	"isPayeeConfirms", "isPlatformConfirms", "contractStatus", "sumFreezeAmt", "sumThawAmt",
	"sumRefundAmt", "residueFreezeAmt", "sumTraceFees", "payRate", "createUserId",
	"updateUserId", "minGranularity", "maxGranularity", "minAshContent", "maxAshContent",
	"minCoalAsh", "maxCoalAsh", "minVolatiles", "maxVolatiles", "minStd", "maxStd",
	"minInternalSulfur", "maxInternalSulfur", "minWaterContent", "maxWaterContent",
	"minBondIndex", "maxBondIndex", "minCalorificValue", "maxCalorificValue", "minCsr",
	"maxCsr", "minY", "maxY", "minHgi", "maxHgi", "recoveryAnalysis", "lithofacies",
	"calorificValue", "cleanCoalRecoveryRate", "middlingCoal", "gangue", "slurry",
	"fstProId", "bigProId", "threeProId", "fstName", "secName", "thName", "isEffect",
	"createTime", "updateTime", "remark1", "remark2", "remark3", "remark4", "remark5"}

// contractIntFields Java 泛型 Integer/LocalDate 字段(getMaps 值原样,不包 %)
var contractIntFields = func() map[string]bool {
	m := map[string]bool{}
	for _, f := range []string{"id", "companyId", "orderId", "templateId", "contractResource",
		"contractType", "salesUnitId", "buyerUnitId", "proBigType", "proSmaType", "priceType",
		"payMethod", "payType", "carrier", "transportMode", "supplierId", "custId", "custType",
		"fileStatus", "sysFlowId", "contractStatus", "payRate", "createUserId", "updateUserId",
		"fstProId", "bigProId", "threeProId", "signingTime"} {
		m[f] = true
	}
	return m
}()

// getMapsOf 对照 StringHelper.getMaps(t):全部非 null 字段,值 Integer/LocalDate 原样、其余 %v%
func getMapsOf(m map[string]any) map[string]any {
	out := map[string]any{}
	for _, f := range contractExcelCols {
		if v, ok := m[f]; ok && v != nil {
			if contractIntFields[f] {
				out[f] = v
			} else {
				out[f] = "%" + str(m, f) + "%"
			}
		}
	}
	return out
}

// listByPage = ChainContractServiceImpl.listByPage(175-284)+ mapper selectContainExt/selectDrafts。
// company_id/is_esign 等直拼(与 Java 相同的注入面),maps 值参数化
func (s *Service) listByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["page"].(map[string]any); !ok {
		platform.Unknown(w) // Java mapper requires a Page object.
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	current, size := pageOf(m)
	searchSql := "where 1=1"
	var args []any

	// contractType 非 null:重写为 company_id 直拼 + contract_type(非 1 一律按 2)
	if v, ok := m["contractType"]; ok && v != nil {
		t := "2"
		if num(v) == 1 {
			t = "1"
		}
		searchSql = "where 1=1 and ((company_id = " + strconvI64(u.CompanyID) + " and contract_type =" + t + ") )"
	}
	// 非 ptSearch:zjjg 头 → 双方单位条件,否则 company_id(Java 会重复拼,行为保持)
	isPtSearch, _ := m["isPtSearch"].(bool)
	if !isPtSearch {
		if r.Header.Get("applytype") == "zjjg" {
			searchSql += " and (c.sales_unit_id = " + strconvI64(u.CompanyID) + " or c.buyer_unit_id = " + strconvI64(u.CompanyID) + ")"
		} else {
			searchSql += " and company_id = " + strconvI64(u.CompanyID)
		}
	}
	if v, ok := m["id"]; ok && v != nil {
		searchSql += " and c.id = ?"
		args = append(args, v)
	}
	if v, ok := m["fileStatus"]; ok && v != nil {
		searchSql += " and c.file_status = ?"
		args = append(args, v)
	}
	if v, ok := m["isEsign"]; ok && v != nil {
		searchSql += " and c.is_esign is " + boolLiteral(v)
	}
	if v, ok := m["isEsignComplete"]; ok && v != nil {
		searchSql += " and c.is_esign_complete is " + boolLiteral(v)
	}
	if v, ok := m["isPayOnline"]; ok && v != nil {
		searchSql += " and c.is_pay_online is " + boolLiteral(v)
	}
	if v, ok := m["startTime"]; ok && v != nil {
		searchSql += " and concat(month, '-01') >= ?"
		args = append(args, firstDay(str(m, "startTime")))
	}
	if v, ok := m["endTime"]; ok && v != nil {
		searchSql += " and concat(month, '-01') <= ?"
		args = append(args, lastDay(str(m, "endTime")))
	}
	if v, ok := m["contractStatus"]; ok && v != nil {
		searchSql += " and c.contract_status = ?"
		args = append(args, v)
	}
	if v, ok := m["salesUnitId"]; ok && v != nil {
		searchSql += " and c.sales_unit_id = ?"
		args = append(args, v)
	}
	if v, ok := m["buyerUnitId"]; ok && v != nil {
		searchSql += " and c.buyer_unit_id = ?"
		args = append(args, v)
	}
	if v, ok := m["ctStartTime"]; ok && v != nil {
		searchSql += " and DATE_FORMAT(c.create_time, '%Y-%m-%d')>= '" + str(m, "ctStartTime") + "'"
	}
	if v, ok := m["ctEndTime"]; ok && v != nil {
		searchSql += " and DATE_FORMAT(c.create_time, '%Y-%m-%d')<= '" + str(m, "ctEndTime") + "'"
	}
	if v, ok := m["isEffect"]; ok && v != nil {
		searchSql += " and is_effect is " + boolLiteral(v)
	}
	if ids, ok := m["ids"].([]any); ok && len(ids) > 0 {
		var parts []string
		for _, it := range ids {
			parts = append(parts, "c.id like "+str(map[string]any{"x": it}, "x"))
		}
		searchSql += " and (" + strings.Join(parts, " or ") + ")"
	}
	if ext, ok := m["ext"].(map[string]any); ok && ext != nil {
		searchSql += " and EXISTS (select 1 from chain_contract_extend where contract_id=c.id"
		if v, ok := ext["dictId"]; ok && v != nil {
			searchSql += " and dict_id = ?"
			args = append(args, v)
		}
		if v, ok := ext["dType"]; ok && v != nil && str(ext, "dType") != "" {
			searchSql += " and dType = ?"
			args = append(args, str(ext, "dType"))
		}
		if v, ok := ext["dValue"]; ok && v != nil && str(ext, "dValue") != "" {
			searchSql += " and dValue = ?"
			args = append(args, str(ext, "dValue"))
		}
		searchSql += ") "
	}
	// searchStr:Java "%"+getSearchStr()+"%",null → "%null%"
	searchSql += " and CONCAT(IFNULL(system_code, ''), IFNULL(contract_code, ''), IFNULL(signing_time, ''), IFNULL(contract_name, ''), IFNULL(sales_unit, ''), IFNULL(buyer_unit, ''), IFNULL(receipt_place, ''), IFNULL(delivery_place, ''), IFNULL(collaborator, ''), IFNULL(pro_remark, ''), IFNULL(contract_num, ''), IFNULL(unit_price, ''), IFNULL(sum_price, ''), IFNULL(sulfur, ''), IFNULL(volatilization, ''), IFNULL(water, ''), IFNULL(remark, '')) like ?"
	args = append(args, "%"+searchStrOf(m)+"%")

	fileStatus := num(m["fileStatus"])
	hasFileStatus := m["fileStatus"] != nil
	base := "SELECT " + strings.Join(dbCols(contractCols), ",") + " FROM chain_contract c " + searchSql
	var q, countQ string
	if hasFileStatus && fileStatus == 0 {
		// 草稿:order by update_time desc
		q = base + " order by update_time desc"
		countQ = "SELECT count(*) FROM chain_contract c " + searchSql
	} else {
		if v, ok := m["notBindOrder"]; ok && v != nil {
			searchSql += " and (c.order_id is null or c.order_id='' or c.order_id=" + orNull(str(m, "orderId")) + ")"
			q = base + " order by create_time desc"
			countQ = "SELECT count(*) FROM chain_contract c " + searchSql
		} else {
			q = base + " order by create_time desc"
			countQ = "SELECT count(*) FROM chain_contract c " + searchSql
		}
	}
	page, err := paged(d, q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeContract(page["records"].([]map[string]any))
	platform.Success(w, page)
}

// orNull notBindOrder 分支的 order_id 直拼(null → 字面 null)
func orNull(s string) string {
	if s == "" {
		return "null"
	}
	return s
}

func boolLiteral(v any) string {
	if v == true {
		return "true"
	}
	return "false"
}

func dbCols(cols []string) []string {
	out := make([]string, 0, len(cols))
	for _, c := range cols {
		out = append(out, "c."+c)
	}
	return out
}

func firstDay(day string) string {
	if t, ok := parseDay(day); ok {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}
	return day
}

func lastDay(day string) string {
	if t, ok := parseDay(day); ok {
		return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}
	return day
}

func parseDay(s string) (time.Time, bool) {
	if len(s) > 10 {
		s = s[:10]
	}
	t, err := time.Parse("2006-01-02", s)
	return t, err == nil
}

func strconvI64(n int64) string {
	return strconv.FormatInt(n, 10)
}

// getSumPrice = mapper getSumPrice:contract_type = ${contractType} 直拼,
// 实测 Java 侧 body 无 contractType 时 ${} 求值致 SQL 报错 → 999(行为保持)
func (s *Service) getSumPrice(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	v, ok := m["contractType"]
	if !ok || v == nil {
		platform.Unknown(w)
		return
	}
	var args []any
	q := "SELECT sum(IFNULL(contract_num, 0)) contractNum, sum(IFNULL(sum_price, 0)) sumPrice FROM chain_contract WHERE contract_type = ?"
	args = append(args, v)
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := platform.ScanMaps(rows)
	out := map[string]any{"contractNum": nil, "sumPrice": nil}
	if len(ms) > 0 {
		if v := ms[0]["contractNum"]; v != nil {
			out["contractNum"] = jsonNum(v)
		}
		if v := ms[0]["sumPrice"]; v != nil {
			out["sumPrice"] = jsonNum(v)
		}
	}
	platform.Success(w, out)
}

func jsonNum(v any) any {
	if s, ok := v.(string); ok && s != "" {
		return json.Number(s)
	}
	return v
}

// countBigPro = mapper countBigPro:d_type='pro_type'(countPro 主路径在 dashboard 包)
func (s *Service) countBigPro(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	searchSql := "where 1=1"
	var args []any
	if v, ok := m["companyId"]; ok && v != nil {
		searchSql += " and c.company_id like ?"
		args = append(args, v) // Integer 原样(无 %)
	}
	q := "SELECT d.d_key, count(c.id) count FROM chain_contract c, chain_common_dict d " + searchSql +
		" and d.d_type='pro_type' and c.pro_sma_type=d.d_value GROUP BY d.d_key"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := platform.ScanMaps(rows)
	out := []map[string]any{}
	for _, row := range ms {
		out = append(out, map[string]any{"d_key": bts(row["d_key"]), "count": num(row["count"])})
	}
	platform.Success(w, out)
}

func bts(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	if v == nil {
		return ""
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// countSumContract = mapper countSumContract:count(1) join 企业(行政区划 like)
func (s *Service) countSumContract(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	searchSql := "where b.company_id=e.id "
	var args []any
	if v, ok := m["companyId"]; ok && v != nil && str(m, "companyId") != "" {
		searchSql += " and b.company_id = ?"
		args = append(args, v)
	}
	if v, ok := m["xzqhCode"]; ok && v != nil && str(m, "xzqhCode") != "0" {
		searchSql += " and CONCAT(province_code,',',city_code,',',area_code, ',',county_code,',') like ?"
		args = append(args, "%"+str(m, "xzqhCode")+",%")
	}
	var cnt int64
	if err := d.QueryRow("SELECT count(1) count FROM chain_contract b, chain_basic_enterprise e "+searchSql, args...).Scan(&cnt); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, map[string]any{"count": cnt})
}

// selZjjgEtc = Java mapper buyEtcs/saleEtcs:`select distinct ... order by create_time`
// 的 ORDER BY 列不在 SELECT 列表,MySQL 8 严格模式直接 SQL 报错 → Java 实测恒 999(行为保持)
func (s *Service) selZjjgEtc(w http.ResponseWriter, r *http.Request) {
	platform.Unknown(w)
}

func camelizeAll(ms []map[string]any) []map[string]any {
	out := []map[string]any{}
	for _, m := range ms {
		out = append(out, platform.Camelize(m))
	}
	return out
}
