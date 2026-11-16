// saas-capital 仪表盘接口(结算量/收付款走势)。
// Java 位置:saas-capital/.../ChainCapitalSettlementController + ChainCapitalDetailedController
package dashboard

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// timeTypeCond timeType 1本月 2上月 3近半年 4本年(无条件),对照各 Java service selectQuat
func timeTypeCond(m map[string]any) string {
	switch str(m, "timeType") {
	case "1":
		return " AND month = DATE_FORMAT(CURDATE(), '%Y-%m')"
	case "2":
		return " AND month = DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m')"
	case "3":
		return " AND concat(month, '-01') BETWEEN date_sub(now(), INTERVAL 6 MONTH) AND now()"
	}
	return ""
}

// settSelectQuantity POST /saas-capital/chainCapitalSettlement/selectQuantity 结算量按类型(有租户)。
// Java ChainCapitalSettlementServiceImpl.selectQuat:322 → mapper selectQuat:42
func (s *Service) settSelectQuantity(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT type, sum(settlement_count) sumNum FROM chain_capital_settlement WHERE 1=1" + timeTypeCond(m)
	var args []any
	q += " GROUP BY type"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	out := []map[string]any{}
	for _, row := range ms {
		out = append(out, map[string]any{"type": num(row["type"]), "sumNum": num(row["sumNum"])})
	}
	platform.Success(w, out)
}

// settSeries 日期序列(timeType=4 本年 12 个月;跨月→月序列;同月→日序列)。
// 对照 Java countSett/countDetail/countLogist 的 dt/dtType/num 推导
func settSeries(m map[string]any) []string {
	if str(m, "timeType") == "4" {
		year := time.Now().Format("2006")
		out := make([]string, 0, 12)
		for i := 1; i <= 12; i++ {
			out = append(out, fmt.Sprintf("%s-%02d", year, i))
		}
		return out
	}
	start, end := str(m, "startTime"), str(m, "endTime")
	if monthDiff(start, end) > 0 {
		return monthSeries(start, end)
	}
	return daySeries(start, end)
}

// settCountSett POST /saas-capital/chainCapitalSettlement/countSett 结算额收/付月走势(有租户)。
// Java ChainCapitalSettlementServiceImpl.countSett:338-447:日期序列 LEFT JOIN 月聚合,
// type=1 → cglb(采购/收款),type=2 → xslb;无数据月补 0 行。响应 {cglb, xslb}
func (s *Service) settCountSett(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	series := settSeries(m)

	q := "SELECT month t, type, sum(sum_price) sumPrice FROM chain_capital_settlement WHERE 1=1"
	var args []any
	q += " GROUP BY month, type"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	hits := map[string][]map[string]any{}
	for _, row := range ms {
		hits[bts(row["t"])] = append(hits[bts(row["t"])], row)
	}

	cglb := []map[string]any{}
	xslb := []map[string]any{}
	skdate, fkdate := "", ""
	for _, key := range series {
		for _, row := range hits[key] {
			item := map[string]any{"date": key, "sumPrice": num(row["sumPrice"])}
			switch bts(row["type"]) {
			case "1":
				item["type"] = 1
				cglb = append(cglb, item)
				skdate += key + ","
			case "2":
				item["type"] = 2
				xslb = append(xslb, item)
				fkdate += key + ","
			}
		}
		if !strings.Contains(skdate, key) {
			cglb = append(cglb, map[string]any{"date": key, "sumPrice": 0, "type": nil})
		}
		if !strings.Contains(fkdate, key) {
			xslb = append(xslb, map[string]any{"date": key, "sumPrice": 0, "type": nil})
		}
	}
	platform.Success(w, map[string]any{"cglb": cglb, "xslb": xslb})
}

// detailCountDetail POST /saas-capital/chainCapitalDetailed/countDetail 收付款流水走势(有租户)。
// Java ChainCapitalDetailedServiceImpl.countDetail:319:月/日粒度,detype=1 → sklb(收款),2 → fklb(付款)
func (s *Service) detailCountDetail(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	series := settSeries(m)
	fmtDay := monthDiff(str(m, "startTime"), str(m, "endTime")) <= 0 // 同月/日序列 → 按天
	format := "%Y-%m"
	if fmtDay {
		format = "%Y-%m-%d"
	}

	q := "SELECT DATE_FORMAT(trans_time, '" + format + "') t, detype, sum(money) sumPrice FROM chain_capital_detailed WHERE 1=1"
	var args []any
	q += " GROUP BY DATE_FORMAT(trans_time, '" + format + "'), detype"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	hits := map[string][]map[string]any{}
	for _, row := range ms {
		hits[bts(row["t"])] = append(hits[bts(row["t"])], row)
	}

	sklb := []map[string]any{}
	fklb := []map[string]any{}
	skdate, fkdate := "", ""
	for _, key := range series {
		for _, row := range hits[key] {
			item := map[string]any{"date": key, "sumPrice": num(row["sumPrice"])}
			switch bts(row["detype"]) {
			case "1":
				item["detype"] = 1
				sklb = append(sklb, item)
				skdate += key + ","
			case "2":
				item["detype"] = 2
				fklb = append(fklb, item)
				fkdate += key + ","
			}
		}
		if !strings.Contains(skdate, key) {
			sklb = append(sklb, map[string]any{"date": key, "sumPrice": 0, "detype": nil})
		}
		if !strings.Contains(fkdate, key) {
			fklb = append(fklb, map[string]any{"date": key, "sumPrice": 0, "detype": nil})
		}
	}
	platform.Success(w, map[string]any{"sklb": sklb, "fklb": fklb})
}
