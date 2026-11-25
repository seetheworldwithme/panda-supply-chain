// bill - stats.go 统计族:countTaxRate/selZeroTax/selJxzb/countBuses/countMonBuses/
// countPtTax/countYearBuses/searchRelBill。BigDecimal 字面算术(math/big),除零 → 999
// (Java ArithmeticException)。
package bill

import (
	"database/sql"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// ── BigDecimal 字面算术(big.Int) ──

type bdec struct{ s string }

func bdecOf(v any) (bdec, bool) {
	switch t := v.(type) {
	case nil:
		return bdec{}, false
	case json.Number:
		return bdec{s: t.String()}, true
	case string:
		if strings.TrimSpace(t) == "" {
			return bdec{}, false
		}
		return bdec{s: t}, true
	default:
		return bdec{s: strings.TrimSpace(btsStr(t))}, true
	}
}

// bparts 拆 [符号]整数|小数位数
func bparts(s string) (*big.Int, int) {
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	intPart, frac := s, ""
	if i := strings.Index(s, "."); i >= 0 {
		intPart, frac = s[:i], s[i+1:]
	}
	n, _ := new(big.Int).SetString(intPart+frac, 10)
	if n == nil {
		n = big.NewInt(0)
	}
	if neg {
		n.Neg(n)
	}
	return n, len(frac)
}

func badd(a, b bdec) bdec {
	ia, sa := bparts(a.s)
	ib, sb := bparts(b.s)
	if sa < sb {
		ia.Mul(ia, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(sb-sa)), nil))
	} else if sb < sa {
		ib.Mul(ib, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(sa-sb)), nil))
	}
	return fmtDec(maxInt(sa, sb), ia.Add(ia, ib))
}

func bcmpZero(a bdec) int {
	n, _ := bparts(a.s)
	return n.Sign()
}

func bcmp(a, b bdec) int {
	ia, sa := bparts(a.s)
	ib, sb := bparts(b.s)
	if sa < sb {
		ia.Mul(ia, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(sb-sa)), nil))
	} else if sb < sa {
		ib.Mul(ib, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(sa-sb)), nil))
	}
	return ia.Cmp(ib)
}

// bdivHalfUp a/b 保留 scale 位 HALF_UP(Java divide(x, 4, ROUND_HALF_UP))
func bdivHalfUp(a, b bdec, scale int) (bdec, bool) {
	if bcmp(b, bdec{s: "0"}) == 0 {
		return bdec{}, false // ArithmeticException → 999
	}
	ia, sa := bparts(a.s)
	ib, sb := bparts(b.s)
	// 缩放对齐到 max(sa,sb) 后再乘 10^scale
	sc := sa
	if sb > sc {
		sc = sb
	}
	pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scale+sc-sa)), nil)
	num := new(big.Int).Mul(ia, pow)
	den := new(big.Int).Mul(ib, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(sc-sb)), nil))
	q, r := new(big.Int).QuoRem(num, den, new(big.Int))
	// HALF_UP:余数绝对值×2 ≥ 除数绝对值 → 远离零进位
	r2 := new(big.Int).Abs(r)
	r2.Mul(r2, big.NewInt(2))
	if r2.Cmp(new(big.Int).Abs(den)) >= 0 {
		if (num.Sign() < 0) != (den.Sign() < 0) {
			q.Sub(q, big.NewInt(1))
		} else {
			q.Add(q, big.NewInt(1))
		}
	}
	return fmtDec(scale, q), true
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func fmtDec(scale int, v *big.Int) bdec {
	neg := ""
	if v.Sign() < 0 {
		neg = "-"
		v = new(big.Int).Neg(v)
	}
	s := v.String()
	for len(s) <= scale {
		s = "0" + s
	}
	if scale == 0 {
		return bdec{s: neg + s}
	}
	return bdec{s: neg + s[:len(s)-scale] + "." + s[len(s)-scale:]}
}

// ── countTaxRate ──

// queryBillMonthType mapper:select month,type,sum(tax_amount),sum(base_money),CONCAT(month,type) mt
func (s *Service) queryBillMonthType(d *sql.DB, searchSql string, args []any) ([]map[string]any, error) {
	rows, err := d.Query("select month, type, sum(tax_amount) tax_amount, sum(base_money) base_money, CONCAT(month, type) mt from chain_bill "+searchSql, args...)
	if err != nil {
		return nil, err
	}
	return platform.ScanMaps(rows)
}

func (s *Service) billCountTaxRate(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	where := "where 1=1 "
	var args []any
	stRaw, stOk := m["startTime"]
	etRaw, etOk := m["endTime"]
	if stOk && stRaw != nil && strings.TrimSpace(str(m, "startTime")) != "" {
		where += " and concat(month, '-01') >= ?"
		args = append(args, str(m, "startTime"))
	}
	if etOk && etRaw != nil && strings.TrimSpace(str(m, "endTime")) != "" {
		where += " and concat(month, '-01') <= ?"
		args = append(args, str(m, "endTime"))
	}
	where += " GROUP BY type, month order by MONTH ASC, type DESC"
	bills, err := s.queryBillMonthType(d, where, args)
	if err != nil {
		platform.Unknown(w)
		return
	}
	res, ok := taxRateCompute(bills, str(m, "startTime"), str(m, "endTime"))
	if !ok {
		platform.Unknown(w) // 补月 NPE / 除零 ArithmeticException
		return
	}
	platform.Success(w, res)
}

// taxRateCompute 对照 countTaxRate 算法(345-417);start/end 为 body 原值(补月用,null → 999)。
func taxRateCompute(bills []map[string]any, start, end string) (map[string]any, bool) {
	rmaps := map[string]bdec{} // xxse{month} / xxbhsje{month}
	res := map[string]any{}    // TreeMap(键序不影响 sort_keys 对比)
	var sumXxse, sumJxse, sumBhsje bdec
	var ldje, xxldje bdec
	zero := bdec{s: "0"}
	var mtAll []string
	for _, bl := range bills {
		mtAll = append(mtAll, btsStr(bl["mt"])+",")
	}
	for _, bl := range bills {
		month := btsStr(bl["month"])
		tax, ok1 := bdecOf(bl["tax_amount"])
		base, ok2 := bdecOf(bl["base_money"])
		if !ok1 || !ok2 {
			return nil, false
		}
		switch btsStr(bl["type"]) {
		case "2":
			xxse := badd(tax, xxldje)
			sumXxse = badd(sumXxse, tax)
			if bcmp(base, zero) > 0 {
				sumBhsje = badd(sumBhsje, base)
			}
			if bcmp(xxse, zero) <= 0 {
				xxldje = badd(xxldje, xxse)
				res[month] = 0
			} else {
				xxldje = zero
				rmaps["xxse"+month] = tax
				rmaps["xxbhsje"+month] = base
			}
		case "1":
			jxse := badd(tax, ldje)
			ldje = zero
			sumJxse = badd(sumJxse, tax)
			if xxse, has := rmaps["xxse"+month]; has {
				xv := badd(xxse, xxldje)
				if bcmp(xv, jxse) < 0 {
					ldje = badd(ldje, bsub(jxse, xv))
					res[month] = 0
				} else if bcmp(xv, jxse) == 0 {
					res[month] = 0
				} else {
					ldje = zero
					v, ok := bdivHalfUp(bsub(xv, jxse), rmaps["xxbhsje"+month], 4)
					if !ok {
						return nil, false
					}
					res[month] = json.Number(v.s)
				}
			} else {
				ldje = badd(ldje, jxse)
				res[month] = 0
			}
		}
	}
	// 补月循环(Java 390-408):startTime/endTime null → NPE 999
	if start == "" || end == "" {
		return nil, false
	}
	st, ok1 := parseLocalDate(start)
	et, ok2 := parseLocalDate(end)
	if !ok1 || !ok2 {
		return nil, false
	}
	cur := st.Format("2006-01")
	i := 1
	for {
		if _, has := res[cur]; !has {
			if xxse, okr := rmaps["xxse"+cur]; okr {
				v, okd := bdivHalfUp(xxse, rmaps["xxbhsje"+cur], 4)
				if !okd {
					return nil, false
				}
				res[cur] = json.Number(v.s)
			} else {
				res[cur] = 0
			}
		}
		cur = time.Date(st.Year(), st.Month()+time.Month(i), 1, 0, 0, 0, 0, time.Local).Format("2006-01") // 按月首日推进:月末起点经 AddDate 会跳过下个月
		i++
		if strings.ReplaceAll(cur, "-", "") > strings.ReplaceAll(et.Format("2006-01"), "-", "") {
			break
		}
	}
	if bcmp(sumBhsje, zero) == 0 {
		res["ljsfl"] = 0
	} else {
		v, ok := bdivHalfUp(bsub(sumXxse, sumJxse), sumBhsje, 4)
		if !ok {
			return nil, false
		}
		res["ljsfl"] = json.Number(v.s)
	}
	return res, true
}

func bsub(a, b bdec) bdec {
	ia, sa := bparts(a.s)
	ib, sb := bparts(b.s)
	if sa < sb {
		ia.Mul(ia, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(sb-sa)), nil))
	} else if sb < sa {
		ib.Mul(ib, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(sa-sb)), nil))
	}
	return fmtDec(maxInt(sa, sb), ia.Sub(ia, ib))
}

// ── selZeroTax / selJxzb ──

// billSelZeroTax selMinMonth → 自最早月至今天;type=2 行且该月无 type=1 → tax/base
// (除零 → 999);type=1 且当月有销项 → (xxse-jxse)/xxbhsje。
func (s *Service) billSelZeroTax(w http.ResponseWriter, r *http.Request) {
	d := s.pick(r)
	var minMonth sql.NullString
	if err := d.QueryRow("select min(month) from chain_bill").Scan(&minMonth); err != nil {
		platform.Unknown(w)
		return
	}
	if !minMonth.Valid || minMonth.String == "" {
		platform.Success(w, []string{})
		return
	}
	where := "where 1=1  and concat(month, '-01') >= ? and concat(month, '-01') <= ? GROUP BY type, month order by MONTH ASC, type DESC"
	bills, err := s.queryBillMonthType(d, where, []any{minMonth.String + "-01", time.Now().Format("2006-01-02")})
	if err != nil {
		platform.Unknown(w)
		return
	}
	rmaps := map[string]bdec{}
	var sumXxse bdec
	var warns []string
	var mts string
	for _, bl := range bills {
		mts += btsStr(bl["mt"]) + ","
	}
	for _, bl := range bills {
		month := btsStr(bl["month"])
		tax, ok1 := bdecOf(bl["tax_amount"])
		base, ok2 := bdecOf(bl["base_money"])
		if !ok1 || !ok2 {
			platform.Unknown(w)
			return
		}
		if btsStr(bl["type"]) == "2" {
			rmaps["xxse"+month] = tax
			rmaps["xxbhsje"+month] = base
			sumXxse = badd(sumXxse, tax)
			if !strings.Contains(mts, month+"1") {
				// sfl = xxse/xxbhsje;除零 → Java ArithmeticException → 999
				v, ok := bdivHalfUp(tax, base, 4)
				if !ok {
					platform.Unknown(w)
					return
				}
				if bcmp(v, bdec{s: "0.1"}) >= 0 {
					warns = append(warns, zeroTaxWarn(month, v))
				}
			}
		}
		if btsStr(bl["type"]) == "1" {
			if xxse, has := rmaps["xxse"+month]; has {
				jxse := tax
				v, ok := bdivHalfUp(bsub(xxse, jxse), rmaps["xxbhsje"+month], 4)
				if !ok {
					platform.Unknown(w)
					return
				}
				if bcmp(v, bdec{s: "0.1"}) >= 0 {
					warns = append(warns, zeroTaxWarn(month, v))
				}
			}
		}
	}
	_ = sumXxse
	platform.Success(w, warns)
}

func zeroTaxWarn(month string, sfl bdec) string {
	openDate := strings.ReplaceAll(month, "-", "年") + "月"
	// Java:sfl.multiply(100).setScale(2)+"%"(乘 100 后 HALF_UP 两位)
	pct := bmulPow(sfl, "100", 2)
	return openDate + "税负率为" + pct.s + "%，请尽快调整！"
}

// bmulPow a×times 后按 HALF_UP 缩到 scale 位
func bmulPow(a bdec, times string, scale int) bdec {
	ia, sa := bparts(a.s)
	ib, _ := bparts(times)
	n := new(big.Int).Mul(ia, ib)
	// 当前 scale=sa;调整为目标 scale(HALF_UP)
	if sa > scale {
		div := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(sa-scale)), nil)
		q, r := new(big.Int).QuoRem(n, div, new(big.Int))
		r2 := new(big.Int).Abs(r)
		r2.Mul(r2, big.NewInt(2))
		if r2.Cmp(div) >= 0 {
			if n.Sign() < 0 {
				q.Sub(q, big.NewInt(1))
			} else {
				q.Add(q, big.NewInt(1))
			}
		}
		n = q
		sa = scale
	} else if sa < scale {
		n.Mul(n, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scale-sa)), nil))
		sa = scale
	}
	return fmtDec(sa, n)
}

// billSelJxzb 恒空数组(Java bug:startTime 误用 openDate → SQL >= NULL 恒假 → warns 恒空)
func (s *Service) billSelJxzb(w http.ResponseWriter, r *http.Request) {
	platform.Success(w, []string{})
}

// ── countBuses / countMonBuses / countPtTax / countYearBuses ──

func xzqhCond(m map[string]any) (string, any, bool) {
	if v, ok := m["xzqhCode"]; ok && v != nil && str(m, "xzqhCode") != "0" {
		return " and CONCAT(province_code,',',city_code,',',area_code, ',',county_code,',') like ?",
			"%" + str(m, "xzqhCode") + ",%", true
	}
	return "", nil, false
}

// billCountBuses sum(includ_money) 无 IFNULL:空集 → 单 NULL 列 → data null
func (s *Service) billCountBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	where := "where b.company_id=e.id and b.type=2"
	var args []any
	if v, ok := m["companyId"]; ok && v != nil {
		where += " and b.company_id = ?"
		args = append(args, v)
	}
	if cond, arg, has := xzqhCond(m); has {
		where += cond
		args = append(args, arg)
	}
	if str(m, "timeType") == "1" {
		where += " and month = DATE_FORMAT(CURDATE(), '%Y-%m')"
	} else if str(m, "timeType") == "3" {
		where += " and concat(month, '-01') between date_sub(now(),interval 6 month) and now()"
	}
	if v, ok := m["startTime"]; ok && v != nil {
		where += " and concat(month, '-01') >= ?"
		args = append(args, str(m, "startTime"))
	}
	if v, ok := m["endTime"]; ok && v != nil {
		where += " and concat(month, '-01') <= ?"
		args = append(args, str(m, "endTime"))
	}
	rows, err := d.Query("select sum(b.includ_money) baseMoney from chain_bill b, chain_basic_enterprise e "+where, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil || len(ms) == 0 {
		platform.Unknown(w)
		return
	}
	v := ms[0]["baseMoney"]
	if v == nil {
		platform.Success(w, nil) // MyBatis 单 NULL 列 → null map
		return
	}
	platform.Success(w, map[string]any{"baseMoney": json.Number(v.(string))})
}

// billCountMonBuses 补月(Java 608-620);startTime/endTime null → NPE 999
func (s *Service) billCountMonBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	where := "where b.company_id=e.id and b.type=2 "
	var args []any
	if v, ok := m["startTime"]; ok && v != nil {
		where += " and concat(month, '-01') >= ?"
		args = append(args, str(m, "startTime"))
	} else {
		platform.Unknown(w) // maps.get("startTime").toString() NPE
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
	rows, err := d.Query("select month, sum(b.base_money) baseMoney from chain_bill b, chain_basic_enterprise e "+where, args...)
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
	// 补月
	if !fillMonths(res, str(m, "startTime"), str(m, "endTime")) {
		platform.Unknown(w)
		return
	}
	platform.Success(w, res)
}

func fillMonths(res map[string]any, start, end string) bool {
	st, ok1 := parseLocalDate(start)
	et, ok2 := parseLocalDate(end)
	if !ok1 || !ok2 {
		return false
	}
	cur := st.Format("2006-01")
	i := 1
	for {
		if _, has := res[cur]; !has {
			res[cur] = 0
		}
		cur = time.Date(st.Year(), st.Month()+time.Month(i), 1, 0, 0, 0, 0, time.Local).Format("2006-01") // 按月首日推进:月末起点经 AddDate 会跳过下个月
		i++
		if strings.ReplaceAll(cur, "-", "") > strings.ReplaceAll(et.Format("2006-01"), "-", "") {
			break
		}
	}
	return true
}

// billCountPtTax controller 用 countPtTaxNew:值=zz_money.toString() 字符串;
// enterprise_name null → NPE 999
func (s *Service) billCountPtTax(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	where := "where b.company_id=e.id "
	var args []any
	if v, ok := m["companyId"]; ok && v != nil && !isEmptyStr(v) {
		where += " and b.company_id = ?"
		args = append(args, v)
	}
	if v, ok := m["startTime"]; ok && v != nil {
		where += " and concat(month, '-01') >= ?"
		args = append(args, str(m, "startTime"))
	}
	if v, ok := m["endTime"]; ok && v != nil {
		where += " and concat(month, '-01') <= ?"
		args = append(args, str(m, "endTime"))
	}
	if cond, arg, has := xzqhCond(m); has {
		where += cond
		args = append(args, arg)
		m["xzqhCode"] = "%" + str(m, "xzqhCode") + ",%"
	}
	where += " GROUP BY company_id"
	rows, err := d.Query("SELECT (SELECT simple_name from chain_basic_enterprise where id=company_id) enterprise_name, sum(IFNULL(zz_money,0))/sum(base_money) zz_money from chain_bill_rate b, chain_basic_enterprise e "+where, args...)
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
		if row["enterprise_name"] == nil || row["zz_money"] == nil {
			platform.Unknown(w) // .toString() NPE → 999
			return
		}
		res[btsStr(row["enterprise_name"])] = btsStr(row["zz_money"])
	}
	platform.Success(w, res)
}

// billCountYearBuses month like '2026%'(当年)
func (s *Service) billCountYearBuses(w http.ResponseWriter, r *http.Request) {
	year := time.Now().Format("2006")
	rows, err := s.pick(r).Query("select IFNULL(sum(includ_money), 0) baseMoney from chain_bill where month like ?", year+"%")
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil || len(ms) == 0 {
		platform.Unknown(w)
		return
	}
	platform.Success(w, map[string]any{"baseMoney": json.Number(btsStr(ms[0]["baseMoney"]))})
}

// billSearchRelBill 手工拼关联查询;b.* 16 列 CONCAT
func (s *Service) billSearchRelBill(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	where := "where 1=1 and s.relation_id=b.relation_id and (b.relation_id is not null and b.relation_id !='')"
	var args []any
	if v, ok := m["type"]; ok && v != nil && !isEmptyStr(v) {
		where += " and b.type = ?"
		args = append(args, v)
	}
	if v, ok := m["contractId"]; ok && v != nil {
		where += " and s.contract_id = ?"
		args = append(args, str(m, "contractId"))
	}
	if v, ok := m["companyName"]; ok && v != nil {
		where += " and b.company_name like ?"
		args = append(args, "%"+str(m, "companyName")+"%")
	}
	if v, ok := m["inType"]; ok && v != nil {
		where += " and b.in_type like ?"
		args = append(args, "%"+str(m, "inType")+"%")
	}
	if v, ok := m["openDate"]; ok && v != nil {
		if strings.TrimSpace(str(m, "openDate")) != "" {
			t, ok2 := parseLocalDate(str(m, "openDate"))
			if !ok2 {
				platform.Unknown(w)
				return
			}
			where += " and b.open_date = ?"
			args = append(args, t.Format("2006-01-02"))
		}
	}
	if v, ok := m["month"]; ok && v != nil && strings.TrimSpace(str(m, "month")) != "" {
		where += " and b.month like ?"
		args = append(args, "%"+str(m, "month")+"%")
	}
	if v, ok := m["relationFlag"]; ok && v != nil {
		if bv, isBool := v.(bool); isBool {
			where += " and b.relation_flag = ?"
			args = append(args, bv)
		} else if sv, isStr := v.(string); isStr && (sv == "true" || sv == "false") {
			where += " and b.relation_flag = ?"
			args = append(args, sv == "true")
		} else if strings.TrimSpace(sv) == "" {
			// 空串绑 null 不拼
		} else {
			platform.Unknown(w)
			return
		}
	}
	if v, ok := m["relationId"]; ok && v != nil && !isEmptyStr(v) {
		where += " and b.relation_id = ?"
		args = append(args, v)
	}
	full := where + " and " + relBillConcat + " like ?"
	args = append(args, searchStrLike(m))
	current, size := pageOf(m)
	q := "select b.* from chain_capital_settlement s, chain_bill b " + full + " order by b.create_time desc"
	countQ := "SELECT count(1) FROM chain_capital_settlement s, chain_bill b " + full
	ip, err := paged(s.pick(r), q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entBill.rows(recs)
	platform.Success(w, ip)
}

const relBillConcat = "CONCAT(IFNULL(b.system_code, ''), IFNULL(b.contract_name, ''), IFNULL(b.bill_code, ''), IFNULL(b.bill_number, ''), IFNULL(b.open_date, ''), IFNULL(b.duty_paragraph, ''), IFNULL(b.company_name, ''), IFNULL(b.goods_name, ''), IFNULL(b.uit, ''), IFNULL(b.count, ''), IFNULL(b.unit_price, ''), IFNULL(b.base_money, ''), IFNULL(b.tax_rate, ''), IFNULL(b.tax_amount, ''), IFNULL(b.includ_money, ''), IFNULL(b.modle, ''))"

// billCountPtTaxOld 老版 countPtTax(Cg/Xs 控制器):mapper SQL 的
// IFNULL(e.simple_name, e.enterprise_name) 列不在 GROUP BY(b.type, e.simple_name) 且非
// 函数依赖 → MySQL 8 ONLY_FULL_GROUP_BY 报 1055 → 恒 999(本库实测),照实复刻。
func (s *Service) billCountPtTaxOld(w http.ResponseWriter, r *http.Request) {
	bodyMap(r)
	platform.Unknown(w)
}
