// saas-logistics 仪表盘接口(物流量统计)。
// Java 位置:saas-logistics/.../ChainLogisticsController + ChainLogisticsBillController
package dashboard

import (
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// loCountBuses POST /saas-logistics/chainLogistics/countBuses 物流量(运费吨)累计(filter=true)。
// Java ChainLogisticsServiceImpl.countBuses:124 → mapper countBuses:35
func (s *Service) loCountBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT cast(sum(b.tonnage_paid) AS decimal(12,2)) tonnagePaid FROM chain_logistics b, chain_basic_enterprise e WHERE b.company_id=e.id"
	var args []any
	if c := str(m, "companyId"); c != "" {
		q += " AND b.company_id = ?"
		args = append(args, c)
	}
	if cond, arg := xzqhCond(m); cond != "" {
		q += cond
		args = append(args, arg)
	}
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	var tonnagePaid float64
	if len(ms) > 0 {
		tonnagePaid = num(ms[0]["tonnagePaid"])
	}
	platform.Success(w, map[string]any{"tonnagePaid": tonnagePaid})
}

// loCountMonBuses POST /saas-logistics/chainLogistics/countMonBuses 按月物流量(Java 补月,filter=true)。
// Java ChainLogisticsServiceImpl.countMonBuses:138 → mapper countMonBuses:39
func (s *Service) loCountMonBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT left(trans_time, 7) trans_time, cast(sum(b.tonnage_paid) AS decimal(12,2)) tonnagePaid FROM chain_logistics b, chain_basic_enterprise e WHERE b.company_id=e.id"
	var args []any
	if s2 := str(m, "startTime"); s2 != "" {
		q += " AND trans_time >= ?"
		args = append(args, s2)
	}
	if s2 := str(m, "endTime"); s2 != "" {
		q += " AND trans_time <= ?"
		args = append(args, s2)
	}
	if c := str(m, "companyId"); c != "" {
		q += " AND b.company_id = ?"
		args = append(args, c)
	}
	if cond, arg := xzqhCond(m); cond != "" {
		q += cond
		args = append(args, arg)
	}
	q += " GROUP BY left(trans_time, 7) ORDER BY left(trans_time, 7) DESC"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	res := map[string]any{}
	for _, row := range ms {
		res[bts(row["trans_time"])] = num(row["tonnagePaid"])
	}
	fillMonths(res, str(m, "startTime"), str(m, "endTime"))
	platform.Success(w, res)
}

// loBillCountLogist POST /saas-logistics/chainLogisticsBill/countLogist 运单量走势(月/日粒度,有租户)。
// Java ChainLogisticsBillServiceImpl.countLogist:272。注意 Java 外层 SELECT 字面量
// `1 transMethod` 恒非 null → 全部命中行进 ztlb,sdlb 只剩补月行——Java 原样行为,照实复刻
func (s *Service) loBillCountLogist(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	series := settSeries(m)
	fmtDay := monthDiff(str(m, "startTime"), str(m, "endTime")) <= 0
	format := "%Y-%m"
	if fmtDay {
		format = "%Y-%m-%d"
	}

	q := "SELECT DATE_FORMAT(trans_time, '" + format + "') t, sum(tonnage_paid) sumPrice FROM chain_logistics WHERE 1=1"
	var args []any
	q += " GROUP BY DATE_FORMAT(trans_time, '" + format + "')"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	hits := map[string]float64{}
	for _, row := range ms {
		hits[bts(row["t"])] = num(row["sumPrice"])
	}

	ztlb := []map[string]any{}
	sdlb := []map[string]any{}
	for _, key := range series {
		ztItem := map[string]any{"date": key, "sumTonnage": 0, "transMethod": 1}
		if v, has := hits[key]; has {
			ztItem["sumTonnage"] = v
		}
		ztlb = append(ztlb, ztItem)
		// Java:transMethod 恒 1 → sdlb 永远只有补月 0 行(补月行 transMethod 为 null,照 Java)
		sdlb = append(sdlb, map[string]any{"date": key, "sumTonnage": 0, "transMethod": nil})
	}
	platform.Success(w, map[string]any{"ztlb": ztlb, "sdlb": sdlb})
}
