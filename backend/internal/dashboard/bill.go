// saas-bill 仪表盘接口(发票/缴税统计与税负率提示)。
// Java 位置:saas-bill/.../ChainBillController + ChainBillRateController → ChainBillServiceImpl/ChainBillRateServiceImpl
package dashboard

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// round4 对照 BigDecimal.ROUND_HALF_UP 保留 4 位
func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}

func today() string { return time.Now().Format("2006-01-02") }

// selJxzb POST /saas-bill/chainBill/selJxzb 进项占比提示。
// Java ChainBillServiceImpl.selJxzb:480 用 selMinMonth 后把 startTime 设为 bill.getOpenDate()
// ——selMinMonth 只查 month 列,openDate 恒 null → SQL `>= NULL` 恒假 → bills 恒空 → 永远返回空数组。
// Go 忠实复刻该线上行为(修复属业务决策,留二期与用户确认)
func (s *Service) selJxzb(w http.ResponseWriter, r *http.Request) {
	platform.Success(w, []string{})
}

// selZeroTax POST /saas-bill/chainBill/selZeroTax 税负率提示。
// Java ChainBillServiceImpl.selZeroTax:423 → countTaxRate SQL(自最早月份至今天,全表)。
// 对每月:无进项(月内无 type=1)时 sfl=销项税额/不含税额;有进项时 sfl=(销项-进项)/销项不含税额;≥0.1 出文案
func (s *Service) selZeroTax(w http.ResponseWriter, r *http.Request) {
	d := s.pick(r)
	var minMonth any
	q := "SELECT min(month) FROM chain_bill WHERE 1=1"
	args := []any{}
	if err := d.QueryRow(q, args...).Scan(&minMonth); err != nil || minMonth == nil {
		platform.Success(w, []string{})
		return
	}
	start := bts(minMonth) + "-01"
	end := today()
	bills, err := s.queryBillMonthType(d, start, end)
	if err != nil {
		platform.Unknown(w)
		return
	}
	warns := []string{}
	xxse := map[string]float64{}    // 月 → 销项税额
	xxbhsje := map[string]float64{} // 月 → 销项不含税金额
	hasType1 := map[string]bool{}
	for _, bl := range bills {
		mo := bts(bl["month"])
		if bts(bl["type"]) == "2" {
			xxse[mo] = num(bl["tax_amount"])
			xxbhsje[mo] = num(bl["base_money"])
		} else if bts(bl["type"]) == "1" {
			hasType1[mo] = true
		}
	}
	for _, bl := range bills {
		mo := bts(bl["month"])
		if bts(bl["type"]) == "2" && !hasType1[mo] {
			if xxbhsje[mo] == 0 {
				continue // Java 此处除零抛异常致 500;Go 跳过该月避免崩溃
			}
			sfl := round4(xxse[mo] / xxbhsje[mo])
			if sfl >= 0.1 {
				warns = append(warns, taxWarn(mo, sfl))
			}
		}
		if bts(bl["type"]) == "1" {
			if _, has := xxse[mo]; has {
				if xxbhsje[mo] == 0 {
					continue
				}
				sfl := round4((xxse[mo] - num(bl["tax_amount"])) / xxbhsje[mo])
				if sfl >= 0.1 {
					warns = append(warns, taxWarn(mo, sfl))
				}
			}
		}
	}
	platform.Success(w, warns)
}

func taxWarn(month string, sfl float64) string {
	openDate := strings.ReplaceAll(month, "-", "年") + "月"
	return openDate + "税负率为" + fmt.Sprintf("%.2f", math.Round(sfl*10000)/100) + "%，请尽快调整！"
}

// countTaxRate POST /saas-bill/chainBill/countTaxRate 月度税负率(留抵税额递推,有租户)。
// Java ChainBillServiceImpl.countTaxRate:325-416,算法逐步复刻:
// 按月排序后,type=2(销项)累计留抵 xxldje;type=1(进项)与当月销项比较,进项>销项的部分转留抵,
// 销项>进项时税负率=(销项税额+留抵-进项税额)/销项不含税金额,否则 0;缺月补 0(或仅有销项时直接算);
// ljsfl=(累计销项税额-累计进项税额)/累计不含税销项金额
func (s *Service) countTaxRate(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	start := str(m, "startTime")
	end := str(m, "endTime")
	if start == "" || end == "" {
		platform.FailMsg(w, "参数错误")
		return
	}
	bills, err := s.queryBillMonthType(d, start, end)
	if err != nil {
		platform.Unknown(w)
		return
	}

	res := map[string]any{}   // month → 税负率
	keep := map[string]float64{} // "xxse{month}" / "xxbhsje{month}"(对照 Java rmaps)
	var sumXxse, sumJxse, sumBhsje float64
	var ldje, xxldje float64 // 进项/销项留抵金额

	for _, bl := range bills {
		mo := bts(bl["month"])
		typ := bts(bl["type"])
		tax := num(bl["tax_amount"])
		base := num(bl["base_money"])
		switch typ {
		case "2": // 销项
			xxse := tax + xxldje
			sumXxse += tax
			if base > 0 {
				sumBhsje += base
			}
			if xxse <= 0 {
				xxldje += xxse
				res[mo] = 0
			} else {
				xxldje = 0
				keep["xxse"+mo] = tax
				keep["xxbhsje"+mo] = base
			}
		case "1": // 进项
			jxse := tax + ldje
			ldje = 0
			sumJxse += tax
			if xxbhs, has := keep["xxse"+mo]; has {
				xxse := xxbhs + xxldje
				if xxse < jxse {
					ldje += jxse - xxse
					res[mo] = 0
				} else if xxse == jxse {
					res[mo] = 0
				} else {
					ldje = 0
					if bhs := keep["xxbhsje"+mo]; bhs != 0 {
						res[mo] = round4((xxse - jxse) / bhs)
					} else {
						res[mo] = 0
					}
				}
			} else {
				ldje += jxse
				res[mo] = 0
			}
		}
	}
	// 补足月份(照 Java 389-408:缺月有销项记录则直接销项/不含税额,否则 0)
	for _, mo := range monthSeries(start, end) {
		if _, has := res[mo]; !has {
			if xxbhs, ok := keep["xxse"+mo]; ok {
				if bhs := keep["xxbhsje"+mo]; bhs != 0 {
					res[mo] = round4(xxbhs / bhs)
				} else {
					res[mo] = 0
				}
			} else {
				res[mo] = 0
			}
		}
	}
	if sumBhsje == 0 {
		res["ljsfl"] = 0
	} else {
		res["ljsfl"] = round4((sumXxse - sumJxse) / sumBhsje)
	}
	platform.Success(w, res)
}

// queryBillMonthType countTaxRate mapper:48(拦截器租户在本环境被 ignore-tables 跳过,全表)
func (s *Service) queryBillMonthType(d *sql.DB, start, end string) ([]map[string]any, error) {
	q := "SELECT month, type, sum(tax_amount) tax_amount, sum(base_money) base_money, CONCAT(month, type) mt FROM chain_bill WHERE 1=1"
	args := []any{}
	if start != "" {
		q += " AND concat(month, '-01') >= ?"
		args = append(args, start)
	}
	if end != "" {
		q += " AND concat(month, '-01') <= ?"
		args = append(args, end)
	}
	q += " GROUP BY type, month ORDER BY month ASC, type DESC"
	rows, err := d.Query(q, args...)
	if err != nil {
		return nil, err
	}
	return scanMaps(rows)
}

// selectBusCount POST /saas-bill/chainBill/selectBusCount 进销项税额统计(filter=true 无租户,companyId 手工拼)。
// Java ChainBillServiceImpl.selectBusCount:211 → mapper selectBusCount:39;jwt companyId!="0" 时用 jwt 值
func (s *Service) selectBusCount(w http.ResponseWriter, r *http.Request) {
	cu := auth.From(r.Context())
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT type, IFNULL(sum(tax_amount),0) sumTaxAmount, IFNULL(sum(includ_money),0) sumIncludMoney FROM chain_bill WHERE 1=1"
	var args []any
	if s := str(m, "startTime"); s != "" {
		q += " AND concat(month, '-01') >= ?"
		args = append(args, s)
	}
	if s := str(m, "endTime"); s != "" {
		q += " AND concat(month, '-01') <= ?"
		args = append(args, s)
	}
	if cu.CompanyID != 0 {
		q += " AND company_id = ?"
		args = append(args, cu.CompanyID)
	} else if c := f64Ptr(m, "companyId"); c != nil {
		q += " AND company_id = ?"
		args = append(args, *c)
	}
	q += " GROUP BY type"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	out := []map[string]any{}
	for _, row := range ms {
		out = append(out, map[string]any{
			"type":           num(row["type"]),
			"sumTaxAmount":   num(row["sumTaxAmount"]),
			"sumIncludMoney": num(row["sumIncludMoney"]),
		})
	}
	platform.Success(w, out)
}

// billCountBuses POST /saas-bill/chainBill/countBuses 开票业务额累计(filter=true)。
// Java ChainBillServiceImpl.countBuses:561 → mapper countBuses:60;timeType 1本月 3近半年
func (s *Service) billCountBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q, args := billBusesCond(m, false)
	q = "SELECT sum(b.includ_money) baseMoney FROM chain_bill b, chain_basic_enterprise e " + q
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	var baseMoney float64
	if len(ms) > 0 {
		baseMoney = num(ms[0]["baseMoney"])
	}
	platform.Success(w, map[string]any{"baseMoney": baseMoney})
}

// billCountMonBuses POST /saas-bill/chainBill/countMonBuses 按月开票额(Java 补月,filter=true)。
// Java ChainBillServiceImpl.countMonBuses:586 → mapper countMonBuses:64
func (s *Service) billCountMonBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q, args := billBusesCond(m, true)
	q = "SELECT month, sum(b.base_money) baseMoney FROM chain_bill b, chain_basic_enterprise e " + q + " GROUP BY month ORDER BY month DESC"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	res := map[string]any{}
	for _, row := range ms {
		res[bts(row["month"])] = num(row["baseMoney"])
	}
	fillMonths(res, str(m, "startTime"), str(m, "endTime"))
	platform.Success(w, res)
}

// billBusesCond countBuses/countMonBuses 共用条件(type=2 销项 + companyId/xzqhCode/时间段)
func billBusesCond(m map[string]any, withTime bool) (string, []any) {
	q := "WHERE b.company_id=e.id AND b.type=2"
	var args []any
	if c := f64Ptr(m, "companyId"); c != nil && *c != 0 {
		q += " AND b.company_id = ?"
		args = append(args, *c)
	}
	if cond, arg := xzqhCond(m); cond != "" {
		q += cond
		args = append(args, arg)
	}
	if withTime {
		if s := str(m, "startTime"); s != "" {
			q += " AND concat(month, '-01') >= ?"
			args = append(args, s)
		}
		if s := str(m, "endTime"); s != "" {
			q += " AND concat(month, '-01') <= ?"
			args = append(args, s)
		}
	} else {
		switch str(m, "timeType") {
		case "1":
			q += " AND month = DATE_FORMAT(CURDATE(), '%Y-%m')"
		case "3":
			q += " AND concat(month, '-01') BETWEEN date_sub(now(), INTERVAL 6 MONTH) AND now()"
		}
		if s := str(m, "startTime"); s != "" {
			q += " AND concat(month, '-01') >= ?"
			args = append(args, s)
		}
		if s := str(m, "endTime"); s != "" {
			q += " AND concat(month, '-01') <= ?"
			args = append(args, s)
		}
	}
	return q, args
}

// countPtTax POST /saas-bill/chainBill/countPtTax 平台各企业税负率(值域 "企业名":"0.xxxx",filter=true)。
// Java ChainBillServiceImpl.countPtTaxNew:685 → mapper countPtTaxNew:83;值是 BigDecimal.toString 字符串
func (s *Service) countPtTax(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT (SELECT simple_name FROM chain_basic_enterprise WHERE id=company_id) enterprise_name," +
		" sum(IFNULL(zz_money,0))/sum(base_money) zz_money FROM chain_bill_rate b, chain_basic_enterprise e WHERE b.company_id=e.id"
	var args []any
	if c := f64Ptr(m, "companyId"); c != nil && *c != 0 {
		q += " AND b.company_id = ?"
		args = append(args, *c)
	}
	if s := str(m, "startTime"); s != "" {
		q += " AND concat(month, '-01') >= ?"
		args = append(args, s)
	}
	if s := str(m, "endTime"); s != "" {
		q += " AND concat(month, '-01') <= ?"
		args = append(args, s)
	}
	if cond, arg := xzqhCond(m); cond != "" {
		q += cond
		args = append(args, arg)
	}
	q += " GROUP BY company_id"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	res := map[string]any{}
	for _, row := range ms {
		if row["enterprise_name"] == nil || row["zz_money"] == nil {
			continue // Java 此处 null.toString 抛异常;Go 跳过空行
		}
		res[bts(row["enterprise_name"])] = numStr(row["zz_money"])
	}
	platform.Success(w, res)
}

// rateGetSumPrice POST /saas-bill/chainBillRate/getSumPrice 缴税金额累计(filter=true)。
// Java ChainBillRateServiceImpl.getSumPrice:126 → mapper getSumPrice:32;xzqhCode 被 Java 忽略,照实
func (s *Service) rateGetSumPrice(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT IFNULL(sum(sum_money),0) sumMoney FROM chain_bill_rate WHERE 1=1"
	var args []any
	if c := f64Ptr(m, "companyId"); c != nil && *c != 0 {
		q += " AND company_id = ?"
		args = append(args, *c)
	}
	if s := str(m, "startTime"); s != "" {
		q += " AND concat(month, '-01') >= ?"
		args = append(args, s)
	}
	if s := str(m, "endTime"); s != "" {
		q += " AND concat(month, '-01') <= ?"
		args = append(args, s)
	}
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	var sumMoney float64
	if len(ms) > 0 {
		sumMoney = num(ms[0]["sumMoney"])
	}
	platform.Success(w, map[string]any{"sumMoney": sumMoney})
}

// rateCountMonBuses POST /saas-bill/chainBillRate/countMonBuses 按月缴税额(Java 补月,filter=true)。
// Java ChainBillRateServiceImpl.countMonBuses:146 → mapper countMonBuses:40
func (s *Service) rateCountMonBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT month, sum(b.sum_money) baseMoney FROM chain_bill_rate b, chain_basic_enterprise e WHERE b.company_id=e.id"
	var args []any
	if c := f64Ptr(m, "companyId"); c != nil && *c != 0 {
		q += " AND b.company_id = ?"
		args = append(args, *c)
	}
	if s := str(m, "startTime"); s != "" {
		q += " AND concat(month, '-01') >= ?"
		args = append(args, s)
	}
	if s := str(m, "endTime"); s != "" {
		q += " AND concat(month, '-01') <= ?"
		args = append(args, s)
	}
	if cond, arg := xzqhCond(m); cond != "" {
		q += cond
		args = append(args, arg)
	}
	q += " GROUP BY month ORDER BY month DESC"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	res := map[string]any{}
	for _, row := range ms {
		res[bts(row["month"])] = num(row["baseMoney"])
	}
	fillMonths(res, str(m, "startTime"), str(m, "endTime"))
	platform.Success(w, res)
}
