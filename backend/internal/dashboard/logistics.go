// saas-logistics 仪表盘接口(物流量统计)——四期物流按 8765 实测校准。
// Java 位置:saas-logistics/.../ChainLogisticsController + ChainLogisticsBillController
// 路由为一期注册,四期只校准语义(不重注册):DECIMAL 保字面小数位、countMonBuses 缺参 999、
// countLogist timeType=4 因补月循环 getStartTime() NPE 恒 999(原版 bug 照实复刻)。
package dashboard

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// loCountBuses POST /saas-logistics/chainLogistics/countBuses 物流量(运费吨)累计(filter=true)。
// Java ChainLogisticsServiceImpl.countBuses:124 → mapper countBuses:35。
// DECIMAL(12,2) 保字面("348.79"/"0.00"),空表 sum→NULL→data null(照 Java)
func (s *Service) loCountBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	q := "SELECT cast(sum(b.tonnage_paid) AS decimal(12,2)) tonnagePaid FROM chain_logistics b, chain_basic_enterprise e WHERE b.company_id=e.id"
	var args []any
	if v, ok := m["companyId"]; ok && v != nil {
		q += " and b.company_id = ?"
		args = append(args, str(m, "companyId"))
	}
	if cond, arg := xzqhCond(m); cond != "" {
		q += cond
		args = append(args, arg)
	}
	var v sql.NullString
	if err := d.QueryRow(q, args...).Scan(&v); err != nil {
		platform.Unknown(w)
		return
	}
	// 照 8765 实测:sum 为 NULL(条件无行)时 MyBatis 单 NULL 列行 → map 为 null → data null
	if !v.Valid {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, map[string]any{"tonnagePaid": json.Number(v.String)})
}

// loCountMonBuses POST /saas-logistics/chainLogistics/countMonBuses 按月物流量(filter=true)。
// Java ChainLogisticsServiceImpl.countMonBuses:138:startTime/endTime 缺失 →
// maps.get("startTime").toString() NPE → 999;ISO 解析失败 → DateTimeParseException → 999。
// 行值 decimal(12,2) 保字面;补月为整数 0
func (s *Service) loCountMonBuses(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	st, ok1 := m["startTime"]
	et, ok2 := m["endTime"]
	if !ok1 || !ok2 || st == nil || et == nil {
		platform.Unknown(w) // maps.get("startTime").toString() NPE
		return
	}
	if _, ok := parseDay(str(m, "startTime")); !ok {
		platform.Unknown(w) // LocalDate.parse ISO 失败
		return
	}
	if _, ok := parseDay(str(m, "endTime")); !ok {
		platform.Unknown(w) // endTime 非法日期 → MySQL 1292 → 999(8765 实测)
		return
	}
	q := "SELECT left(trans_time, 7) trans_time, cast(sum(b.tonnage_paid) AS decimal(12,2)) tonnagePaid FROM chain_logistics b, chain_basic_enterprise e WHERE b.company_id=e.id"
	var args []any
	q += " and trans_time >= ?"
	args = append(args, str(m, "startTime"))
	q += " and trans_time <= ?"
	args = append(args, str(m, "endTime"))
	if v, ok := m["companyId"]; ok && v != nil {
		q += " and b.company_id = ?"
		args = append(args, str(m, "companyId"))
	}
	if cond, arg := xzqhCond(m); cond != "" {
		q += cond
		args = append(args, arg)
	}
	q += " GROUP BY left(trans_time, 7) order by left(trans_time, 7) desc"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := scanMaps(rows)
	res := map[string]any{}
	for _, row := range ms {
		key := bts(row["trans_time"])
		if v, ok := row["tonnagePaid"].(string); ok && v != "" {
			res[key] = json.Number(v)
		} else {
			res[key] = 0
		}
	}
	// 补足月份(Java 同款循环:缺失月整数 0)
	fillMonths(res, str(m, "startTime"), str(m, "endTime"))
	platform.Success(w, res)
}

// loBillCountLogist POST /saas-logistics/chainLogisticsBill/countLogist 日期统计物流。
// Java ChainLogisticsBillServiceImpl.countLogist:272:
//   - timeType=4:先建 @cdate 月序列(SQL 可执行),随后补月循环 getStartTime() NPE → 恒 999
//   - 其余:startTime/endTime 必填(缺失 monthDiff 拆箱 NPE → 999);跨月 → 月粒度 num=月差+1,
//     同月 → 日粒度 num=天数+1;end<start → LIMIT 负数 → SQL 错误 → 999
//   - SQL 照 Java 生成(@cdate 用户变量序列 + LEFT JOIN 聚合);transMethod 恒 1 →
//     ztlb=全部序列行,sdlb=补月 0 行(仅 date/sumTonnage 两键,照 Java)
func (s *Service) loBillCountLogist(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)

	tt, hasTT := m["timeType"]
	if hasTT && tt != nil && num(tt) == 4 {
		platform.Unknown(w) // 补月循环 getStartTime() NPE,原版恒 999
		return
	}
	st, ok1 := m["startTime"]
	et, ok2 := m["endTime"]
	if !ok1 || !ok2 || st == nil || et == nil {
		platform.Unknown(w) // monthDiff(start,end) 拆箱 NPE
		return
	}
	startD, okS := parseDay(str(m, "startTime"))
	endD, okE := parseDay(str(m, "endTime"))
	if !okS || !okE {
		platform.Unknown(w) // LocalDate 绑定失败
		return
	}
	jg := monthDiff(str(m, "startTime"), str(m, "endTime"))
	dtType, format := "DAY", "%Y-%m-%d"
	if jg > 0 {
		dtType, format = "MONTH", "%Y-%m"
	} else if jg == 0 {
		jg = int(endD.Sub(startD).Hours() / 24)
	}
	// num = jg+1(Java (jg==0?days:jg)+1;jg<0 → LIMIT 负数 → SQL 错误 999)
	num := jg + 1
	endLit := "'" + endD.Format("2006-01-02") + "'"
	sqlStr := "SELECT" +
		"	date_c.*, IFNULL(sumPrice, 0) sumTonnage, 1 transMethod" +
		" FROM" +
		"	(" +
		"		SELECT" +
		"			DATE_FORMAT(" +
		"				@cdate := date_add(@cdate, INTERVAL - 1 " + dtType + ")," +
		"				'" + format + "'" +
		"			) AS date" +
		"		FROM" +
		"			(" +
		"				SELECT" +
		"					@cdate := date_add(" + endLit + ", INTERVAL 1 " + dtType + ")" +
		"				FROM" +
		"					chain_common_dict" +
		"				LIMIT " + itoa(num) +
		"			) d" +
		"		ORDER BY" +
		"			date" +
		"	) date_c" +
		" LEFT JOIN (" +
		"	SELECT" +
		"		sum(tonnage_paid) AS sumPrice," +
		"		DATE_FORMAT(trans_time, '" + format + "') AS t, 1" +
		"	FROM" +
		"		chain_logistics" +
		"	WHERE" +
		"		1 = 1" +
		"	GROUP BY" +
		"		DATE_FORMAT(trans_time, '" + format + "')" +
		") tab ON t = date"
	rows, err := d.Query(sqlStr)
	if err != nil {
		platform.Unknown(w) // LIMIT 负数/表缺失等 SQL 异常
		return
	}
	ms, err := scanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ztlb := []map[string]any{}
	sdlb := []map[string]any{}
	skdate := ""
	for _, row := range ms {
		item := map[string]any{"date": bts(row["date"]), "sumTonnage": 0, "transMethod": 1}
		if v, ok := row["sumTonnage"].(string); ok && v != "" && v != "0" {
			item["sumTonnage"] = json.Number(v)
		}
		ztlb = append(ztlb, item) // transMethod 恒 1 → 全部进 ztlb
		skdate += bts(row["date"]) + ","
	}
	// 补足月份(Java 同款循环):ztlb 按 skdate 补;sdlb 按 fkdate 补——fkdate 只累计
	// transMethod=2 行,本 SQL 恒无 → sdlb 恒补全部月份(8765 实测)
	fkdate := ""
	start := startD
	endMo := endD.Format("2006-01")
	for i := 1; ; i++ {
		cur := start.AddDate(0, i-1, 0)
		mo := cur.Format("2006-01")
		if !strings.Contains(skdate, mo) {
			ztlb = append(ztlb, map[string]any{"date": mo, "sumTonnage": 0})
		}
		if !strings.Contains(fkdate, mo) {
			sdlb = append(sdlb, map[string]any{"date": mo, "sumTonnage": 0})
		}
		if mo >= endMo {
			break
		}
	}
	platform.Success(w, map[string]any{"ztlb": ztlb, "sdlb": sdlb})
}

func itoa(n int) string { return strings.TrimSpace(strings.Repeat("", 0) + itoaImpl(n)) }

func itoaImpl(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
