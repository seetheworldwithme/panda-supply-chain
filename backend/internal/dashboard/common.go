// saas-common 仪表盘接口(企业简介/业务量统计/企业列表/区划/图片)。
// Java 位置:saas-common/.../ChainBasicMainPageController + ChainBasicEnterpriseController + ChainAreaInfoController + UploadController
package dashboard

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/commonrest"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// selEtcMsg POST /saas-common/chainBasicMainPage/selEtcMsg 企业简介(名称/经营煤种/承运平台/地址)。
// Java ChainBasicBussinessServiceImpl.selEtcMsg:322;pro/selLogistics 的拦截器租户在本环境被
// ignore-tables 全量跳过 → 全表;企业名/地址按 jwt companyId 手工查询(照 Java isNoneBlank 分支)
func (s *Service) selEtcMsg(w http.ResponseWriter, r *http.Request) {
	cu := auth.From(r.Context())
	d := s.pick(r)
	res := map[string]any{"pro": "", "carriagePlatform": "暂无数据"}

	// 煤种:合同大类的字典名去重拼接(selPro)
	rows, err := d.Query(`SELECT (SELECT d_key FROM chain_common_dict WHERE d_type='pro_type' AND d_value=pro_big_type) AS d_key
		FROM chain_contract GROUP BY pro_big_type`)
	if err == nil {
		if pros, err := scanMaps(rows); err == nil {
			var ps []string
			for _, p := range pros {
				if v := bts(p["d_key"]); strings.TrimSpace(v) != "" {
					ps = append(ps, v)
				}
			}
			res["pro"] = strings.Join(ps, ",")
		}
	}
	// 承运平台:最新一条运单的平台字典名(selLogistics)
	var cp any
	_ = d.QueryRow(`SELECT (SELECT d_key FROM chain_common_dict WHERE d_type='logistics_platform' AND d_value=carriage_platform)
		FROM chain_logistics ORDER BY trans_time LIMIT 0,1`).Scan(&cp)
	if s2, ok := cp.([]byte); ok && len(s2) > 0 {
		res["carriagePlatform"] = string(s2)
	} else if s2, ok := cp.(string); ok && s2 != "" {
		res["carriagePlatform"] = s2
	}
	// 企业名与地址(jwt companyId)
	if cu.CompanyID > 0 {
		var name, addr any
		_ = d.QueryRow("SELECT enterprise_name, address FROM chain_basic_enterprise WHERE id = ?", cu.CompanyID).Scan(&name, &addr)
		res["name"] = bts(name)
		res["address"] = bts(addr)
	}
	platform.Success(w, res)
}

// bts sql.NullString 简化版:driver []byte/nil → 字符串
func bts(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case []byte:
		return string(t)
	case string:
		return t
	default:
		return fmt.Sprint(t)
	}
}

// countActBus POST /saas-common/chainBasicMainPage/countActBus 累计实际业务量。
// Java ChainBasicBussinessServiceImpl.countActBus:88 → mapper countActBus:30
func (s *Service) countActBus(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	sqlStr := "SELECT IFNULL(sum(bus_count),0) busCount, IFNULL(sum(logistics_count),0) logisticsCount, IFNULL(sum(tax_payment),0) taxCount, IFNULL(sum(trans_count),0) transCount" +
		" FROM chain_basic_bussiness b, chain_basic_enterprise e WHERE b.company_id=e.id"
	var args []any
	if c := f64Ptr(m, "companyId"); c != nil {
		sqlStr += " AND b.company_id = ?"
		args = append(args, *c)
	}
	if s := str(m, "startTime"); s != "" {
		sqlStr += " AND concat(month, '-01') >= ?"
		args = append(args, s)
	}
	if s := str(m, "endTime"); s != "" {
		sqlStr += " AND concat(month, '-01') <= ?"
		args = append(args, s)
	}
	if cond, arg := xzqhCond(m); cond != "" {
		sqlStr += cond
		args = append(args, arg)
	}
	rows, err := d.Query(sqlStr, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := scanMaps(rows)
	if err != nil || len(ms) == 0 {
		platform.Unknown(w)
		return
	}
	row := ms[0]
	platform.Success(w, map[string]any{
		"busCount":       num(row["busCount"]),
		"logisticsCount": num(row["logisticsCount"]),
		"taxCount":       num(row["taxCount"]),
		"transCount":     num(row["transCount"]),
	})
}

// countMonActBus POST /saas-common/chainBasicMainPage/countMonActBus 按月业务量(Java 补月)。
// Java ChainBasicBussinessServiceImpl.countMonActBus:109 → mapper countMonActBus:33
func (s *Service) countMonActBus(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	sqlStr := "SELECT month, IFNULL(sum(bus_count),0) busCount, IFNULL(sum(logistics_count),0) logisticsCount, IFNULL(sum(tax_payment),0) taxPayment" +
		" FROM chain_basic_bussiness b, chain_basic_enterprise e WHERE b.company_id=e.id"
	var args []any
	if c := f64Ptr(m, "companyId"); c != nil {
		sqlStr += " AND b.company_id = ?"
		args = append(args, *c)
	}
	if s := str(m, "startTime"); s != "" {
		sqlStr += " AND concat(month, '-01') >= ?"
		args = append(args, s)
	}
	if s := str(m, "endTime"); s != "" {
		sqlStr += " AND concat(month, '-01') <= ?"
		args = append(args, s)
	}
	if cond, arg := xzqhCond(m); cond != "" {
		sqlStr += cond
		args = append(args, arg)
	}
	sqlStr += " GROUP BY month ORDER BY month DESC"
	rows, err := d.Query(sqlStr, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	resBus := map[string]any{}
	resLog := map[string]any{}
	resTax := map[string]any{}
	for _, row := range ms {
		mo, _ := row["month"].(string)
		resBus[mo] = num(row["busCount"])
		resLog[mo] = num(row["logisticsCount"])
		resTax[mo] = num(row["taxPayment"])
	}
	fillMonths(resBus, resLog, resTax, str(m, "startTime"), str(m, "endTime"))
	platform.Success(w, map[string]any{"busCount": resBus, "logisticsCount": resLog, "taxPayment": resTax})
}

// fillMonths 月序列补零(等价 Java 补足月份循环)
func fillMonths(maps ...any) {
	// 最后两个参数是 start/end,前面是各序列 map
	if len(maps) < 3 {
		return
	}
	endIdx := len(maps) - 1
	end, _ := maps[endIdx].(string)
	start, _ := maps[endIdx-1].(string)
	for _, mo := range monthSeries(start, end) {
		for _, mv := range maps[:endIdx-1] {
			if mm, ok := mv.(map[string]any); ok {
				if _, has := mm[mo]; !has {
					mm[mo] = 0
				}
			}
		}
	}
}

// countMonData POST /saas-common/chainBasicMainPage/countMonData 开票企业数/资料完成企业数按月。
// Java countMonData:170 → mapper countMonOpen:40(无租户 filter=true) + countMonImport:44
func (s *Service) countMonData(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	cond, args := timeCond(m, "month")

	resOpen := map[string]any{}
	if rows, err := d.Query("SELECT month, count(1) openCount FROM chain_basic_bussiness WHERE 1=1"+cond+" AND bus_count>0 GROUP BY month ORDER BY month DESC", args...); err == nil {
		if ms, _ := scanMaps(rows); ms != nil {
			for _, row := range ms {
				resOpen[bts(row["month"])] = num(row["openCount"])
			}
		}
	}
	resImport := map[string]any{}
	if rows, err := d.Query("SELECT month, count(DISTINCT company_id) importCount FROM chain_bill WHERE 1=1"+cond+" GROUP BY month ORDER BY month DESC", args...); err == nil {
		if ms, _ := scanMaps(rows); ms != nil {
			for _, row := range ms {
				resImport[bts(row["month"])] = num(row["importCount"])
			}
		}
	}
	fillMonths(resOpen, resImport, str(m, "startTime"), str(m, "endTime"))
	platform.Success(w, map[string]any{"openCount": resOpen, "importCount": resImport})
}

// countNetLogist POST /saas-common/chainBasicMainPage/countNetLogist 网络货运占比。
// Java countNetLogist:216:带 companyId 走 buses/netlogists 月序列,否则 netlist/translist 占比(filter=true 无租户)
func (s *Service) countNetLogist(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	res := map[string]any{}

	tcond, targs := timeCond(m, "month")

	if c := f64Ptr(m, "companyId"); c != nil {
		// 月序列分支(mapper countByMonth:72 + countLoByMonth:76,均 filter=true)
		buses := []map[string]any{}
		if rows, err := d.Query("SELECT month, sum(IFNULL(trans_count,0)) transCount, sum(IFNULL(logistics_count,0)) logisticsCount FROM chain_basic_bussiness WHERE 1=1"+tcond+" AND company_id = ? GROUP BY month", append(targs, *c)...); err == nil {
			if ms, _ := scanMaps(rows); ms != nil {
				for _, row := range ms {
					buses = append(buses, map[string]any{
						"month":          bts(row["month"]),
						"transCount":     num(row["transCount"]),
						"logisticsCount": num(row["logisticsCount"]),
					})
				}
			}
		}
		netlogists := []map[string]any{}
		q := "SELECT DATE_FORMAT(trans_time, '%Y-%m') transTime, sum(IFNULL(tonnage_paid,0)) tonnagePaid FROM chain_logistics WHERE trans_time IS NOT NULL"
		qArgs := []any{*c}
		if s := str(m, "startTime"); s != "" {
			q += " AND DATE_FORMAT(trans_time, '%Y-%m-%d') >= ?"
			qArgs = append(qArgs, s)
		}
		if s := str(m, "endTime"); s != "" {
			q += " AND DATE_FORMAT(trans_time, '%Y-%m-%d') <= ?"
			qArgs = append(qArgs, s)
		}
		q += " AND company_id = ? GROUP BY DATE_FORMAT(trans_time, '%Y-%m')"
		if rows, err := d.Query(q, qArgs...); err == nil {
			if ms, _ := scanMaps(rows); ms != nil {
				for _, row := range ms {
					netlogists = append(netlogists, map[string]any{
						"transTime":   bts(row["transTime"]),
						"tonnagePaid": num(row["tonnagePaid"]),
					})
				}
			}
		}
		res["buses"] = buses
		res["netlogists"] = netlogists
	} else {
		// 占比分支(mapper countNetLogist:48 + countTrans:52,均 filter=true)
		netlist := []map[string]any{}
		nq := `SELECT (SELECT simple_name FROM chain_basic_enterprise e WHERE e.id=l.company_id) simpleName,
			sum(tonnage_paid)/(SELECT sum(logistics_count) FROM chain_basic_bussiness WHERE company_id=l.company_id AND (concat(month, '-01') >= ? AND concat(month, '-01') <= ?)) lrate
			FROM chain_logistics l WHERE 1=1`
		nargs := []any{str(m, "startTime"), str(m, "endTime")}
		if s := str(m, "startTime"); s != "" {
			nq += " AND trans_time >= ?"
			nargs = append(nargs, s)
		}
		if s := str(m, "endTime"); s != "" {
			nq += " AND trans_time <= ?"
			nargs = append(nargs, s)
		}
		nq += " GROUP BY company_id"
		if rows, err := d.Query(nq, nargs...); err == nil {
			if ms, _ := scanMaps(rows); ms != nil {
				for _, row := range ms {
					netlist = append(netlist, map[string]any{"simpleName": bts(row["simpleName"]), "lrate": num(row["lrate"])})
				}
			}
		}
		translist := []map[string]any{}
		tq := `SELECT (SELECT simple_name FROM chain_basic_enterprise e WHERE e.id=b.company_id) simpleName, IFNULL(sum(trans_count)/sum(logistics_count),0) lorate
			FROM chain_basic_bussiness b WHERE 1=1` + tcond + " GROUP BY company_id"
		if rows, err := d.Query(tq, targs...); err == nil {
			if ms, _ := scanMaps(rows); ms != nil {
				for _, row := range ms {
					translist = append(translist, map[string]any{"simpleName": bts(row["simpleName"]), "lorate": num(row["lorate"])})
				}
			}
		}
		res["netlist"] = netlist
		res["translist"] = translist
	}
	// lastedMonth:业务量表最近月份(QueryWrapper orderByDesc month 取一)
	var lm any
	_ = d.QueryRow("SELECT month FROM chain_basic_bussiness ORDER BY month DESC LIMIT 1").Scan(&lm)
	if lm != nil {
		res["lastedMonth"] = bts(lm)
	}
	platform.Success(w, res)
}

// countEtcTax POST /saas-common/chainBasicMainPage/countEtcTax 企业所得税占比(filter=true 无租户)。
// Java countEtcTax:251 → mapper countEtcTax:56
func (s *Service) countEtcTax(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	cond, args := timeCond(m, "month")
	rows, err := d.Query(`SELECT (SELECT simple_name FROM chain_basic_enterprise e WHERE e.id=l.company_id) simpleName,
		sum(qy_money)/sum(base_money) trate FROM chain_bill_rate l WHERE 1=1`+cond+" GROUP BY company_id", args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	out := []map[string]any{}
	for _, row := range ms {
		out = append(out, map[string]any{"simpleName": bts(row["simpleName"]), "trate": num(row["trate"])})
	}
	platform.Success(w, out)
}

// countLogistArea POST /saas-common/chainBasicMainPage/countLogistArea 收发货地轨迹。
// Java countLogistArea:267 → mapper countLogistArea:59(拦截器在本环境被 ignore-tables 跳过,全表)
func (s *Service) countLogistArea(w http.ResponseWriter, r *http.Request) {
	d := s.pick(r)
	q := `SELECT loading_area loadingArea, max(loading_lnglog) loadingLnglog, unloading_area unloadingArea,
		max(unloading_lnglog) unloadingLnglog, count(1) logCount FROM chain_logistics
		WHERE loading_area IS NOT NULL AND unloading_area IS NOT NULL`
	var args []any
	q += " GROUP BY loading_area, unloading_area ORDER BY loading_area DESC"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	out := []map[string]any{}
	for _, row := range ms {
		out = append(out, map[string]any{
			"loadingArea":     bts(row["loadingArea"]),
			"loadingLnglog":   bts(row["loadingLnglog"]),
			"unloadingArea":   bts(row["unloadingArea"]),
			"unloadingLnglog": bts(row["unloadingLnglog"]),
			"logCount":        num(row["logCount"]),
		})
	}
	platform.Success(w, out)
}

// countPriceTrend POST /saas-common/chainBasicMainPage/countPriceTrend 价格走势。
// Java countPriceTrend:275 → mapper countPriceTrend:62(月均价 AVG(DISTINCT unit_price),不补月)
func (s *Service) countPriceTrend(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT type, month, AVG(DISTINCT unit_price) unitPrice FROM chain_capital_settlement WHERE 1=1"
	var args []any
	if s := str(m, "startTime"); s != "" {
		q += " AND concat(month, '-01') >= ?"
		args = append(args, s)
	}
	if s := str(m, "endTime"); s != "" {
		q += " AND concat(month, '-01') <= ?"
		args = append(args, s)
	}
	q += " GROUP BY type, `month` ORDER BY type ASC, `month` DESC"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	cg := map[string]any{}
	xs := map[string]any{}
	for _, row := range ms {
		mo, _ := row["month"].(string)
		if bts(row["type"]) == "1" {
			cg[mo] = num(row["unitPrice"])
		} else if bts(row["type"]) == "2" {
			xs[mo] = num(row["unitPrice"])
		}
	}
	platform.Success(w, map[string]any{"cgprice": cg, "xxprice": xs})
}

// ptlist POST /saas-common/chainBasicEnterprise/ptlist 平台企业下拉列表(filter=true 无租户)。
// Java ChainBasicEnterpriseServiceImpl.ptlist:201 → mapper ptlist:40(is_chain_etc=1 且备案号非空)。
// 前端仪表盘只传 {countyCode, searchStr},getSql 对 countyCode 生成无 % 的 like(等值)
func (s *Service) ptlist(w http.ResponseWriter, r *http.Request) {
	commonrest.NewService(s.DB).EnterpriseList(w, r)
}

// entSelectCount POST /saas-common/chainBasicEnterprise/selectCount 企业数量(filter=true)。
// Java ChainBasicEnterpriseServiceImpl.selectCount:173;xzqhCode 缺失时 Java 拼 '%null,%' 恒假 → 全 0,照实复刻
func (s *Service) entSelectCount(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT count(1) count FROM chain_basic_enterprise WHERE 1=1"
	var args []any
	if x := str(m, "xzqhCode"); x != "" && x != "0" {
		q += " AND CONCAT(province_code,',',city_code,',',area_code,',',county_code,',') LIKE ?"
		args = append(args, "%"+x+",%")
	} else {
		// 复刻 Java:键缺失/null 时 like '%null,%' 恒假(前端 sessionStorage 无 xzqhCode 才触发)
		q += " AND CONCAT(province_code,',',city_code,',',area_code,',',county_code,',') LIKE ?"
		args = append(args, "%null,%")
	}
	var sumCount, monCount, yearCount int64
	_ = d.QueryRow(q, args...).Scan(&sumCount)
	_ = d.QueryRow(q+" AND DATE_FORMAT(create_time, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')", args...).Scan(&monCount)
	_ = d.QueryRow(q+" AND DATE_FORMAT(create_time, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m') AND DATE_FORMAT(create_time, '%Y') = DATE_FORMAT(CURDATE(), '%Y')", args...).Scan(&yearCount)
	platform.Success(w, map[string]any{"sumCount": sumCount, "monCount": monCount, "yearCount": yearCount})
}

// selByParent POST /saas-common/chainAreaInfo/selByParent 下级区划(parentCode=0 只查区划,否则限"有企业入驻")。
// Java ChainAreaInfoServiceImpl.selByParent:66;区划表为公共表无租户
func (s *Service) selByParent(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	parent := str(m, "parentCode")
	var q string
	if parent == "0" {
		q = "SELECT * FROM chain_area_info WHERE parent_code = ?"
	} else {
		q = `SELECT * FROM chain_area_info WHERE EXISTS (
			SELECT 1 FROM chain_basic_enterprise WHERE parent_code = ? AND
			CONCAT(province_code, ',', city_code, ',', area_code, ',', county_code, ',') LIKE CONCAT('%', chain_area_info.code, ',', '%'))`
	}
	rows, err := d.Query(q, parent)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	if ms == nil {
		ms = []map[string]any{}
	}
	platform.Success(w, ms)
}
