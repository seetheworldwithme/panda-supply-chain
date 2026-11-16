// stats2.go 合同统计补遗:四流合一(searchFour 多分支)、运单分页(getLogistics)、
// 系统编号(getSystemCode)、Xs/Cg 复用的数量/份额端点、toPdfFile 行为保持。
// searchFour 的 fromBill 分支复刻 Java bug:deMap 与 billMap 都取 selBillByBid(detMoney=billMoney)
package contract

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// pick 按请求头 applytype 选库(默认 chain,zjjg 切监管库),对照 Java DynamicDataSource
func (s *Service) pick(r *http.Request) *sql.DB {
	if r.Header.Get("applytype") == "zjjg" {
		return s.DB.Zjjg
	}
	return s.DB.Chain
}

// getSystemCode = ChainContractServiceImpl.getSystemCode:
// idt+enterpriseCode+yyyyMMdd 前缀查 max(system_code),尾号 +1 补零 7 位。
// enterpriseCode 运行时取自主库 chain_basic_enterprise(JWT claims 里恒空,实测 8765);
// 前缀串 ${} 直拼;systemCode 长度<20 或尾号非数字 → Java 异常 999
func (s *Service) getSystemCode(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	u := currentUser(r)
	idt := ""
	if str(m, "codeType") == "1" {
		idt = "HTCG."
	} else if str(m, "codeType") == "2" {
		idt = "HTXS."
	}
	var etpCode string
	_ = s.DB.Chain.QueryRow("SELECT IFNULL(enterprise_code,'') FROM chain_basic_enterprise WHERE id = ?", u.CompanyID).Scan(&etpCode)
	leftDate := idt + etpCode + time.Now().Format("20060102")
	var systemCode string
	err := d.QueryRow("SELECT IFNULL(max(system_code),'') FROM chain_contract WHERE system_code LIKE ?", leftDate+"%").Scan(&systemCode)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if systemCode == "" {
		platform.Success(w, leftDate+autoGenericCode("1", 7))
		return
	}
	if len(systemCode) < 20 {
		platform.Unknown(w) // substring(0,20) 越界
		return
	}
	s3, err := strconv.Atoi(systemCode[20:])
	if err != nil {
		platform.Unknown(w) // 尾号 parseInt 失败
		return
	}
	platform.Success(w, systemCode[:20]+autoGenericCode(strconv.Itoa(s3), 7))
}

// cttSelectQuantity / cttCountPro:Xs/Cg 控制器路径的合同量与产品份额
// (语义 = dashboard 包同名 handler,此处为 Xs/Cg 路径注册)
func (s *Service) cttSelectQuantity(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	cond := timeTypeCond(m)
	htllist := []map[string]any{}
	if rows, err := d.Query("SELECT contract_type, sum(contract_num) sumNum FROM chain_contract WHERE 1=1"+cond+" GROUP BY contract_type"); err == nil {
		if ms, _ := platform.ScanMaps(rows); ms != nil {
			for _, row := range ms {
				htllist = append(htllist, map[string]any{"contract_type": num(row["contract_type"]), "sumNum": num(row["sumNum"])})
			}
		}
	}
	fslist := []map[string]any{}
	if rows, err := d.Query("SELECT contract_type, count(1) contractNum FROM chain_contract WHERE 1=1"+cond+" GROUP BY contract_type"); err == nil {
		if ms, _ := platform.ScanMaps(rows); ms != nil {
			for _, row := range ms {
				fslist = append(fslist, map[string]any{"contract_type": num(row["contract_type"]), "contractNum": num(row["contractNum"])})
			}
		}
	}
	platform.Success(w, map[string]any{"htllist": htllist, "fslist": fslist})
}

// timeTypeCond 对照 selectQuat/selectCount 的 timeType 分支(1 本月/2 上月/3 近半年/4 空)
func timeTypeCond(m map[string]any) string {
	switch num(m["timeType"]) {
	case 1:
		return " and month = DATE_FORMAT(CURDATE(), '%Y-%m')"
	case 2:
		return " and month = date_format(DATE_SUB(curdate(), INTERVAL 1 MONTH),'%Y-%m')"
	case 3:
		return " and concat(month, '-01') between date_sub(now(),interval 6 month) and now()"
	}
	return ""
}

func (s *Service) cttCountPro(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT d.d_key, count(c.id) count FROM chain_contract c, chain_common_dict d WHERE 1=1"
	var args []any
	if v, ok := m["companyId"]; ok && v != nil {
		q += " AND c.company_id LIKE ?"
		args = append(args, str(m, "companyId"))
	}
	q += " AND d.d_type='pro__middle_type' AND c.pro_sma_type=d.d_value GROUP BY d.d_key"
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

// ── searchFour 四流合一 ──

// searchFour = ChainContractServiceImpl.searchFour(796-870)。
// 空 body → contractId=null → scCta(null) 无行 → NPE → 999;settId 无行 → 零响应写死值
func (s *Service) searchFour(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	settId, hasSett := m["settId"]
	detId, hasDet := m["detId"]
	billId, hasBill := m["billId"]
	loId, hasLo := m["loId"]
	orderId, hasOrder := m["orderId"]
	var mp map[string]any = map[string]any{} // Java 初始 new HashMap():空 body 也非 null → scCta(null) → 999
	if hasSett && settId != nil {
		mp = oneRow(d, "SELECT c.id contractId, c.order_id orderId FROM chain_contract c, chain_capital_settlement s WHERE c.id=s.contract_id AND s.id = ?", settId) // selSett(无行 → null)
	} else if hasDet && detId != nil {
		s.searchFourReply(w, d, fromDetail(d, detId))
		return
	} else if hasBill && billId != nil {
		s.searchFourReply(w, d, fromBill(d, billId))
		return
	} else if hasLo && loId != nil {
		s.searchFourReply(w, d, fromlo(d, loId))
		return
	}
	contractId, hasCid := m["contractId"]
	if !hasCid || contractId == nil {
		if mp == nil {
			// Java 写死的零响应(mp 为 null:selSett 无行)
			platform.Success(w, zeroFour())
			return
		}
		contractId = mp["contractId"] // 空 mp(空 body)→ null
		orderId = mp["orderId"]
		hasOrder = orderId != nil
	}
	if contractId == nil {
		// Java scCta:where id = ${contractId} 直拼,null 直拼致 SQL 报错 → 999(行为保持)
		platform.Unknown(w)
		return
	}
	out := map[string]any{}
	ct := oneRow(d, "SELECT IFNULL(sum(sum_price), 0) sumPrice, IFNULL(sum(contract_num), 0) contractNum FROM chain_contract WHERE id = ?", contractId) // scCta
	if ct == nil {
		platform.Unknown(w) // ct null → getSumPrice() NPE(空 body 路径)
		return
	}
	setMap := oneRow(d, "SELECT count(1) count, IFNULL(sum(sum_price), 0) sumPrice, IFNULL(sum(settlement_count), 0) settlementCount FROM chain_capital_settlement WHERE contract_id = ?", contractId)
	deMap := oneRow(d, "SELECT count(d.id) count, sum(s.money) money FROM chain_capital_detail_settlement s, chain_capital_detailed d WHERE s.detail_id=d.id AND s.contract_id = ?", contractId)
	billMap := oneRow(d, "SELECT count(b.id) count, IFNULL(sum(count), 0) billCount, IFNULL(sum(includ_money), 0) includMoney FROM chain_capital_settlement s, chain_bill b WHERE (s.relation_id is not null and s.relation_id !='') and s.relation_id=b.relation_id AND s.contract_id = ?", contractId)
	if !hasOrder || orderId == nil {
		out["loNum"] = 0
		out["loMoney"] = "---"
		out["loCount"] = 0
		out["orderNum"] = 0
	} else {
		loMap := oneRow(d, "SELECT count(1) count, IFNULL(sum(tonnage_paid), 0) tonnagePaid FROM chain_logistics WHERE order_id = ?", orderId)
		out["loNum"] = numOf(loMap, "count")
		out["loMoney"] = "---"
		out["loCount"] = numOf(loMap, "tonnagePaid")
		out["orderNum"] = 1
	}
	out["contractId"] = contractId
	out["contractNum"] = 1
	out["contractMoney"] = jsonNum(ct["sumPrice"])
	out["contractCount"] = jsonNum(ct["contractNum"])
	out["settNum"] = numOf(setMap, "count")
	out["settMoney"] = jsonNumOf(setMap, "sumPrice")
	out["settCount"] = jsonNumOf(setMap, "settlementCount")
	out["detNum"] = numOf(deMap, "count")
	out["detMoney"] = jsonNumOf(deMap, "money")
	out["detCount"] = "---"
	out["billNum"] = numOf(billMap, "count")
	out["billMoney"] = jsonNumOf(billMap, "includMoney")
	out["billCount"] = jsonNumOf(billMap, "billCount")
	s.searchFourReply(w, d, out)
}

func (s *Service) searchFourReply(w http.ResponseWriter, d *sql.DB, out map[string]any) {
	platform.Success(w, out)
}

func zeroFour() map[string]any {
	return map[string]any{"contractId": 0, "contractNum": 0, "contractMoney": 0, "contractCount": 0,
		"settNum": 0, "settMoney": 0, "settCount": 0, "detNum": 0, "detMoney": 0, "detCount": "---",
		"billNum": 0, "billMoney": 0, "billCount": 0, "loNum": 0, "loMoney": "---", "loCount": 0}
}

// numOf Map 标量取数(null → 0;Java count 键由 IFNULL/聚合保证或 NPE 已在上游拦截)
func numOf(m map[string]any, k string) any {
	if m == nil {
		return nil
	}
	return m[k]
}

func jsonNumOf(m map[string]any, k string) any {
	if m == nil {
		return nil
	}
	if v := m[k]; v != nil {
		return jsonNum(v)
	}
	return nil
}

func oneRow(d *sql.DB, q string, args ...any) map[string]any {
	rows, err := d.Query(q, args...)
	if err != nil {
		return nil
	}
	ms, _ := platform.ScanMaps(rows)
	if len(ms) == 0 {
		return nil
	}
	return ms[0]
}

// fromDetail detId 分支:流水 → hxCode 拼模糊条件 + 五路标量
func fromDetail(d *sql.DB, detId any) map[string]any {
	det := oneRow(d, "SELECT * FROM chain_capital_detailed WHERE id = ?", detId)
	if det == nil || det["hx_code"] == nil || str(det, "hx_code") == "" {
		return zeroFour()
	}
	sqlCond := "where"
	idx := 0
	for _, s := range splitCSV(str(det, "hx_code")) {
		if s == "" {
			continue
		}
		if idx == 0 {
			sqlCond += " (hx_code like '%" + s + "%'"
		} else {
			sqlCond += " or hx_code like '%" + s + "%'"
		}
		idx++
	}
	sqlCond += ")"
	deMap := oneRow(d, "SELECT count(id) count, IFNULL(sum(money-khxje), 0) money FROM chain_capital_detailed "+sqlCond)
	ctaMap := oneRow(d, "SELECT count(cc.id) count, IFNULL(sum(cc.sum_price), 0) sumPrice, IFNULL(sum(cc.contract_num), 0) contractNum FROM (select DISTINCT c.id, c.sum_price, c.contract_num from chain_capital_detail_settlement s, chain_contract c where s.contract_id=c.id and s.detail_id = ?) cc", detId)
	setMap := oneRow(d, "SELECT count(c.id) count, IFNULL(sum(c.sum_price), 0) sumPrice, IFNULL(sum(c.settlement_count), 0) settlementCount FROM chain_capital_detail_settlement s, chain_capital_settlement c WHERE s.statement_id=c.id AND s.detail_id = ?", detId)
	billMap := oneRow(d, "SELECT count(b.id) count, IFNULL(sum(b.includ_money), 0) includMoney, IFNULL(sum(b.count), 0) billCount FROM chain_capital_detail_settlement s, chain_capital_settlement c, chain_bill b WHERE s.statement_id=c.id AND c.relation_id=b.relation_id AND s.detail_id = ?", detId)
	loMap := oneRow(d, "SELECT count(l.id) count, IFNULL(sum(l.tonnage_paid), 0) tonnagePaid FROM chain_capital_detail_settlement ds, chain_contract c, chain_logistics l WHERE ds.contract_id=c.id AND c.order_id=l.order_id AND ds.detail_id = ?", detId)
	ctas := idList(d, "SELECT c.id FROM chain_capital_detail_settlement s, chain_contract c WHERE s.contract_id=c.id AND s.detail_id = ?", detId)
	return assembleFour(ctas, ctaMap, setMap, deMap, billMap, loMap)
}

// fromBill billId 分支:行为保持——deMap 误用 selBillByBid(Java 930 行,与 billMap 同源)
func fromBill(d *sql.DB, billId any) map[string]any {
	bill := oneRow(d, "SELECT * FROM chain_bill WHERE id = ?", billId)
	if bill == nil || bill["relation_id"] == nil || str(bill, "relation_id") == "" {
		return zeroFour()
	}
	ctaMap := oneRow(d, "SELECT count(cc.id) count, IFNULL(sum(cc.sum_price), 0) sumPrice, IFNULL(sum(cc.contract_num), 0) contractNum FROM (SELECT DISTINCT c.id, c.sum_price, c.contract_num from chain_contract c, chain_capital_settlement s, chain_bill b where c.id=s.contract_id and s.relation_id=b.relation_id and b.id=?) cc", billId)
	setMap := oneRow(d, "SELECT count(s.id) count, IFNULL(sum(s.sum_price), 0) sumPrice, IFNULL(sum(s.settlement_count), 0) settlementCount FROM chain_capital_settlement s, chain_bill b WHERE s.relation_id=b.relation_id AND b.id = ?", billId)
	deMap := oneRow(d, "SELECT count(m.id) count, IFNULL(sum(m.includ_money), 0) includMoney, IFNULL(sum(m.count), 0) billCount FROM chain_bill m, chain_bill n WHERE m.relation_id = n.relation_id AND n.id = ?", billId) // Java bug:selBillByBid
	billMap := oneRow(d, "SELECT count(m.id) count, IFNULL(sum(m.includ_money), 0) includMoney, IFNULL(sum(m.count), 0) billCount FROM chain_bill m, chain_bill n WHERE m.relation_id = n.relation_id AND n.id = ?", billId)
	loMap := oneRow(d, "SELECT count(l.id) count, IFNULL(sum(l.tonnage_paid), 0) tonnagePaid FROM chain_capital_settlement s, chain_bill b, chain_contract c, chain_logistics l WHERE s.relation_id=b.relation_id AND s.contract_id=c.id AND c.order_id=l.order_id AND b.id = ?", billId)
	ctas := idList(d, "SELECT DISTINCT s.contract_id FROM chain_capital_settlement s, chain_bill b WHERE s.relation_id=b.relation_id AND b.id = ?", billId)
	return assembleFour(ctas, ctaMap, setMap, deMap, billMap, loMap)
}

// fromlo loId 分支
func fromlo(d *sql.DB, loId any) map[string]any {
	lo := oneRow(d, "SELECT * FROM chain_logistics WHERE id = ?", loId)
	if lo == nil || lo["order_id"] == nil || str(lo, "order_id") == "" {
		return zeroFour()
	}
	ctaMap := oneRow(d, "SELECT count(c.id) count, IFNULL(sum(c.sum_price), 0) sumPrice, IFNULL(sum(c.contract_num), 0) contractNum FROM chain_contract c, chain_logistics l WHERE c.order_id=l.order_id AND l.id = ?", loId)
	setMap := oneRow(d, "SELECT count(s.id) count, IFNULL(sum(s.sum_price), 0) sumPrice, IFNULL(sum(s.settlement_count), 0) settlementCount FROM chain_capital_settlement s, chain_contract c, chain_logistics l WHERE s.contract_id=c.id AND c.order_id=l.order_id AND l.id = ?", loId)
	deMap := oneRow(d, "SELECT count(ds.detail_id) count, IFNULL(sum(ds.money), 0) money FROM chain_capital_detail_settlement ds, chain_capital_settlement s, chain_contract c, chain_logistics l WHERE ds.statement_id=s.id AND s.contract_id=c.id AND c.order_id=l.order_id AND l.id = ?", loId)
	billMap := oneRow(d, "SELECT count(b.id) count, IFNULL(sum(b.includ_money), 0) includMoney, IFNULL(sum(b.count), 0) billCount FROM chain_bill b, chain_capital_settlement s, chain_contract c, chain_logistics l WHERE b.relation_id=s.relation_id AND s.contract_id=c.id AND c.order_id=l.order_id AND l.id = ?", loId)
	loMap := oneRow(d, "SELECT count(m.id) count, IFNULL(sum(m.tonnage_paid), 0) tonnagePaid FROM chain_logistics m, chain_logistics n WHERE m.order_id=n.order_id AND n.id = ?", loId)
	ctas := idList(d, "SELECT c.id FROM chain_contract c, chain_logistics l WHERE c.order_id=l.order_id AND l.id = ?", loId)
	return assembleFour(ctas, ctaMap, setMap, deMap, billMap, loMap)
}

func assembleFour(ctas []map[string]any, ctaMap, setMap, deMap, billMap, loMap map[string]any) map[string]any {
	if ctas == nil {
		ctas = []map[string]any{}
	}
	out := map[string]any{
		"contractIds":   ctas, // Java mapper select 返回 Map 直出,保留 raw 列名(如 contract_id)
		"contractNum":   numOf(ctaMap, "count"),
		"contractMoney": jsonNumOf(ctaMap, "sumPrice"),
		"contractCount": jsonNumOf(ctaMap, "contractNum"),
		"settNum":       numOf(setMap, "count"),
		"settMoney":     jsonNumOf(setMap, "sumPrice"),
		"settCount":     jsonNumOf(setMap, "settlementCount"),
		"detNum":        numOf(deMap, "count"),
		"detMoney":      jsonNumOf(deMap, "money"),
		"detCount":      "---",
		"billNum":       numOf(billMap, "count"),
		"billMoney":     jsonNumOf(billMap, "includMoney"),
		"billCount":     jsonNumOf(billMap, "billCount"),
		"loNum":         numOf(loMap, "count"),
		"loMoney":       "---",
		"loCount":       jsonNumOf(loMap, "tonnagePaid"),
	}
	return out
}

func idList(d *sql.DB, q string, args ...any) []map[string]any {
	rows, err := d.Query(q, args...)
	if err != nil {
		return nil
	}
	ms, _ := platform.ScanMaps(rows)
	return ms
}

func splitCSV(s string) []string {
	return strings.Split(s, ",")
}

// ── getLogistics ──

// logisticsExcelCols ChainLogistics @ExcelColumn 字段(getSql(t,"id") 范围)
var logisticsExcelCols = []string{"id", "imBatchNumber", "batchNumber", "companyId", "orderId",
	"systemCode", "billNumber", "transTime", "carNumber", "loadingAddress", "unloadingPoint",
	"subjectMatter", "freightUnit", "freightAmount", "originalTonnage", "tonnagePaid",
	"payMoney", "carriagePlatform", "fileIds", "loCreateTime", "loadingProvince",
	"unloadingProvince", "loadingArea", "unloadingArea", "loadingLnglog", "unloadingLnglog",
	"loadingProLnglog", "unloadingProLnglog", "loadingLng", "unloadingLng", "invoicingStatus",
	"invoicingAmount", "invoicingTonnage", "updateTime", "createTime", "remark1", "remark2",
	"remark3", "remark4", "remark5"}

var logisticsIntFields = map[string]bool{"id": true, "companyId": true, "orderId": true,
	"carriagePlatform": true, "invoicingStatus": true}

var logisticsDecimalCols = []string{"freight_unit", "freight_amount", "original_tonnage",
	"tonnage_paid", "pay_money", "invoicing_amount", "invoicing_tonnage"}

var logisticsCarrier = []string{"searchStr", "page", "startTime", "endTime", "ids", "list",
	"filelist", "notBindOrder"}

// shapeLogistics 实体形状:createTime 有 @JsonFormat(空格),transTime/loCreateTime/updateTime 无(T)
func shapeLogistics(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		for _, c := range logisticsDecimalCols {
			if s, ok := r[camel(c)].(string); ok && s != "" {
				r[camel(c)] = json.Number(s)
			}
		}
		for _, c := range []string{"transTime", "loCreateTime", "updateTime"} {
			if s, ok := r[c].(string); ok && len(s) == 19 {
				r[c] = s[:10] + "T" + s[11:]
			}
		}
		fillNulls(r, logisticsCarrier)
	}
	return out
}

// getLogistics = mapper getLogistics + getLoSumPrices:
// SELECT l.* 双表链 c.id=${} 直拼;searchStr 由 getMaps %包裹;响应 {logistics: IPage, sumTonnagePaid}
func (s *Service) getLogistics(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	current, size := pageOf(m)
	searchSql, sargs := getSqlEq(m, logisticsExcelCols, map[string]bool{"id": true}, logisticsIntFields)
	cid := "null"
	if v, ok := m["id"]; ok && v != nil {
		cid = str(m, "id")
	}
	// searchStr:getMaps 值 %v%,缺失 → 参数 null(like null 永假)
	var ss any
	if v, ok := m["searchStr"]; ok && v != nil {
		ss = "%" + str(m, "searchStr") + "%"
	}
	concat := "CONCAT(IFNULL(bill_number, ''), IFNULL(trans_time,''), IFNULL(car_number, ''), IFNULL(loading_address, ''), IFNULL(unloading_point, ''), IFNULL(subject_matter, ''), IFNULL(freight_unit, ''), IFNULL(freight_amount, 0), IFNULL(original_tonnage, ''), IFNULL(tonnage_paid, ''), IFNULL(pay_money, ''), IFNULL(carriage_platform, ''))"
	args := append(append([]any{}, sargs...), ss)
	q := "SELECT l.* FROM chain_logistics l, chain_contract c " + searchSql +
		" and l.order_id=c.order_id and c.id = " + cid + " and " + concat + " like ?"
	countQ := "SELECT count(*) FROM chain_logistics l, chain_contract c " + searchSql +
		" and l.order_id=c.order_id and c.id = " + cid + " and " + concat + " like ?"
	page, err := paged(d, q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeLogistics(page["records"].([]map[string]any))
	sumRow := oneRow(d, "SELECT IFNULL(sum(l.tonnage_paid), 0) tonnage_paid FROM chain_logistics l, chain_contract c "+searchSql+" and l.order_id=c.order_id and c.id = "+cid+" and "+concat+" like ?", args...)
	platform.Success(w, map[string]any{"logistics": page, "sumTonnagePaid": jsonNumOf(sumRow, "tonnage_paid")})
}

// toPdfFile 读硬编码 D:\gsht\...\HTMB.pdf(运行环境必失败),catch 后返回纯文本
func (s *Service) toPdfFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	_, _ = w.Write([]byte("This is to pdf"))
}

// _ 保持 time 导入(firstDay/lastDay 在 stats.go)
var _ = time.Now
