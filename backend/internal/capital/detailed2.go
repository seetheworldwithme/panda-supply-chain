// capital - detailed2.go 收付款单查询/统计端点(getSystemCode/getSments/getSumPrice/
// getPtSumPrice/countDetail/getBanks/getDetails)。
package capital

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// dtGetSystemCode GET /getSystemCode:codeType=1→SFCG. 2→SFXS. 空→""
func (s *Service) dtGetSystemCode(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	u := currentUser(r)
	var idt string
	switch str(m, "codeType") {
	case "1":
		idt = "SFCG."
	case "2":
		idt = "SFXS."
	}
	d := s.pick(r)
	like := idt + u.EnterpriseCode + strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", "")
	code, ok := systemCodeNext(d, "chain_capital_detailed", like, "1")
	if !ok {
		platform.Unknown(w) // substring 越界/后缀非数字 → 999
		return
	}
	platform.Success(w, code)
}

// dtGetSments mapper:三表连接 SELECT s.*(Settlement 实体形状)
func (s *Service) dtGetSments(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w) // ${id} 绑定 Integer,body 缺 id → null → SQL "d.id = null" 恒空?
		// MyBatis ${null} 输出空串 → "d.id = " 语法错 → 999。照 999。
		return
	}
	rows, err := s.pick(r).Query("SELECT s.id,s.im_batch_number,s.system_code,s.company_id,s.contract_id,s.contract_code,s.contract_name,s.goods_name,s.set_start_time,s.set_end_time,s.settlement_count,s.unit_price,s.sum_price,s.khxje,s.sjhxje,s.type,s.supplier_id,s.supplier_name,s.rel_detail_flag,s.rel_bill_flag,s.hx_status,s.hx_code,s.month,s.relation_id,s.create_time,s.update_time,s.remark1,s.remark2,s.remark3,s.remark4,s.remark5 from chain_capital_settlement s, chain_capital_detailed d, chain_capital_detail_settlement sd where s.id=sd.statement_id and sd.detail_id = d.id and d.id = ?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if ms == nil {
		ms = []map[string]any{}
	}
	platform.Success(w, entSettlement.rows(ms))
}

// dtGetSumPrice mapper:IFNULL(sum(d.money),0) sumMoney, IFNULL(sum(d.money-d.khxje),0) sumYhxje
func (s *Service) dtGetSumPrice(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	where := " where d.hx_code=s.hx_code "
	var args []any
	if v, ok := m["contractId"]; ok && v != nil {
		where += " and (d.hx_code is not null and d.hx_code !='') and s.contract_id= ?"
		args = append(args, str(m, "contractId"))
	}
	if v, ok := m["detype"]; ok && v != nil {
		where += " and detype = ?"
		args = append(args, str(m, "detype"))
	}
	if v, ok := m["timeType"]; ok && v != nil && num(v) == 4 {
		where += " and YEAR(trans_time)=YEAR(NOW())"
	}
	rows, err := s.pick(r).Query("select IFNULL(sum(d.money), 0) sumMoney, IFNULL(sum(d.money-d.khxje), 0) sumYhxje from chain_capital_detailed d, chain_capital_settlement s"+where, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil || len(ms) == 0 {
		platform.Unknown(w)
		return
	}
	out := map[string]any{}
	for _, k := range []string{"sumMoney", "sumYhxje"} {
		if s2, ok := ms[0][k].(string); ok {
			out[k] = json.Number(s2)
		} else {
			out[k] = ms[0][k]
		}
	}
	platform.Success(w, out)
}

// dtGetPtSumPrice 平台流水合计(Map body):companyId 非空串等值;xzqhCode≠"0" like。
// IFNULL 包裹 → 空集输出 {"sumMoney":0}
func (s *Service) dtGetPtSumPrice(w http.ResponseWriter, r *http.Request) {
	where, args, ok := ptSumPriceCond(r, bodyMap(r), "IFNULL(sum(money), 0) sumMoney from chain_capital_detailed b, chain_basic_enterprise e ")
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := s.pick(r).Query("SELECT "+where, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil || len(ms) == 0 {
		platform.Unknown(w)
		return
	}
	out := map[string]any{}
	for k, v := range ms[0] {
		if s2, ok := v.(string); ok {
			out[k] = json.Number(s2)
		} else {
			out[k] = v
		}
	}
	platform.Success(w, out)
}

// ptSumPriceCond 共用 where 组装:where b.company_id=e.id + 可选条件(注入语义原样保留:
// companyId 值原样参数化;xzqhCode 拼进 maps 覆盖)。返回完整 SELECT 尾段。
func ptSumPriceCond(r *http.Request, m map[string]any, selectFrom string) (string, []any, bool) {
	where := "where b.company_id=e.id "
	var args []any
	if v, ok := m["companyId"]; ok && v != nil && str(m, "companyId") != "" {
		where += " and b.company_id = ?"
		args = append(args, v)
	}
	if v, ok := m["xzqhCode"]; ok && v != nil && str(m, "xzqhCode") != "0" {
		where += " and CONCAT(province_code,',',city_code,',',area_code, ',',county_code,',') like ?"
		args = append(args, "%"+str(m, "xzqhCode")+",%")
	}
	return selectFrom + where, args, true
}

// dtCountDetail 日期统计(Cg/Xs 路径;主路径在 dashboard)。timeType=4 或手工区间,
// 空 startTime/endTime → NPE 999。响应 {sklb, fklb},行 {date,sumPrice[,detype]}。
func (s *Service) dtCountDetail(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	series, format, ok := countSeries(m, true)
	if !ok {
		platform.Unknown(w) // startTime/endTime null NPE → 999
		return
	}
	d := s.pick(r)
	q := "SELECT DATE_FORMAT(trans_time, '" + format + "') t, detype, sum(money) sumPrice FROM chain_capital_detailed WHERE 1=1 GROUP BY DATE_FORMAT(trans_time, '" + format + "'), detype ORDER BY t, detype"
	rows, err := d.Query(q)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, errS := platform.ScanMaps(rows)
	if errS != nil { // 迭代中断连时原先返回全 0 的"成功"
		platform.Unknown(w)
		return
	}
	st, et := str(m, "startTime"), str(m, "endTime")
	sklb, fklb, ok2 := bucketSeries(series, ms, "detype", "1", "2", st, et)
	if !ok2 {
		platform.Unknown(w) // 补月循环 getStartTime() NPE → 999
		return
	}
	platform.Success(w, map[string]any{"sklb": sklb, "fklb": fklb})
}

// countSeries 对照 countDetail/countSett 序列生成(不含补月):
// timeType=4 → 本年 12 月(粒度月);否则 startTime/endTime null → 999;
// detailed 版 jg=monthDiff>0 → 月粒度,jg=0 → 日粒度 num=days+1;settlement 恒月粒度。
func countSeries(m map[string]any, detail bool) (series []string, format string, ok bool) {
	if v, ok2 := m["timeType"]; ok2 && v != nil && num(v) == 4 {
		year := time.Now().Format("2006")
		for i := 1; i <= 12; i++ {
			series = append(series, fmt.Sprintf("%s-%02d", year, i))
		}
		return series, "%Y-%m", true
	}
	sv, ok1 := m["startTime"]
	ev, ok2 := m["endTime"]
	if !ok1 || !ok2 || sv == nil || ev == nil {
		return nil, "", false
	}
	st, ok3 := parseLocalDate(str(m, "startTime"))
	et, ok4 := parseLocalDate(str(m, "endTime"))
	if !ok3 || !ok4 {
		return nil, "", false
	}
	jg := (et.Year()-st.Year())*12 + int(et.Month()) - int(st.Month())
	if !detail {
		for i := 0; i <= jg; i++ { // 按年月推进并限次:AddDate 会把"1月31日+1月"进位成 3 月,漏掉 2 月后原循环永不退出直到 OOM
			series = append(series, time.Date(st.Year(), st.Month()+time.Month(i), 1, 0, 0, 0, 0, time.Local).Format("2006-01"))
		}
		return series, "%Y-%m", true
	}
	if jg > 0 {
		for i := 0; i <= jg; i++ { // 按年月推进并限次:AddDate 会把"1月31日+1月"进位成 3 月,漏掉 2 月后原循环永不退出直到 OOM
			series = append(series, time.Date(st.Year(), st.Month()+time.Month(i), 1, 0, 0, 0, 0, time.Local).Format("2006-01"))
		}
		return series, "%Y-%m", true
	}
	days := int(et.Sub(st).Hours()/24) + 1
	for i := 0; i < days; i++ {
		series = append(series, st.AddDate(0, 0, i).Format("2006-01-02"))
	}
	return series, "%Y-%m-%d", true
}

// bucketSeries SQL 行分桶(行序=SQL 返回序)+补月循环(Java 无条件用 startTime/endTime,
// null → NPE 999)。SQL 行 3 键 {date,sumPrice,key};null 双桶(行保留 null 键);
// 非 v1/v2 值不进桶;补月行仅 {date,sumPrice:0}。
func bucketSeries(series []string, rows []map[string]any, key, v1, v2, st, et string) ([]map[string]any, []map[string]any, bool) {
	hits := map[string][]map[string]any{}
	for _, row := range rows {
		hits[btsStr(row["t"])] = append(hits[btsStr(row["t"])], row)
	}
	l1 := []map[string]any{}
	l2 := []map[string]any{}
	seen1, seen2 := "", ""
	for _, key2 := range series {
		if len(hits[key2]) == 0 {
			// 序列 LEFT JOIN 无聚合匹配 → {date, 0.00, null} 双桶
			item := map[string]any{"date": key2, "sumPrice": json.Number("0.00"), key: nil}
			l1 = append(l1, item)
			l2 = append(l2, item)
			seen1 += key2 + ","
			seen2 += key2 + ","
			continue
		}
		for _, row := range hits[key2] {
			sp, _ := row["sumPrice"].(string)
			item := map[string]any{"date": key2, "sumPrice": json.Number(sp)}
			if row[key] == nil {
				item[key] = nil
				l1 = append(l1, item)
				l2 = append(l2, item)
				seen1 += key2 + ","
				seen2 += key2 + ","
				continue
			}
			switch btsStr(row[key]) {
			case v1:
				item[key] = row[key]
				l1 = append(l1, item)
				seen1 += key2 + ","
			case v2:
				item[key] = row[key]
				l2 = append(l2, item)
				seen2 += key2 + ","
			}
		}
	}
	// 补月循环:st 月起逐月,i 从 1(先判 st 本身),数字比较越过 et 即停;null → 999
	if st == "" || et == "" {
		return nil, nil, false
	}
	stT, ok1 := parseLocalDate(st)
	etT, ok2 := parseLocalDate(et)
	if !ok1 || !ok2 {
		return nil, nil, false
	}
	for i := 1; ; i++ {
		cur := time.Date(stT.Year(), stT.Month()+time.Month(i-1), 1, 0, 0, 0, 0, time.Local).Format("2006-01") // 原先每轮都取起点月:补月循环不推进
		if cur > etT.Format("2006-01") {
			break
		}
		if !strings.Contains(seen1, cur) {
			l1 = append(l1, map[string]any{"date": cur, "sumPrice": 0})
		}
		if !strings.Contains(seen2, cur) {
			l2 = append(l2, map[string]any{"date": cur, "sumPrice": 0})
		}
	}
	return l1, l2, true
}

// dtGetBanks select bank_name group by bank_name
func (s *Service) dtGetBanks(w http.ResponseWriter, r *http.Request) {
	rows, err := s.pick(r).Query("select bank_name from chain_capital_detailed GROUP BY bank_name")
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if ms == nil {
		ms = []map[string]any{}
	}
	platform.Success(w, ms)
}

// dtGetDetails mapper:d.*, (d.money-d.khxje) sjhxje(同名列 Java MyBatis 后列覆盖=calc 值,
// Go map 覆盖语义一致)+CONCAT 8 列 like searchStr order by trans_time desc
func (s *Service) dtGetDetails(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	where := "where d.id=s.detail_id"
	var args []any
	if v, ok := m["contractId"]; ok && v != nil {
		where += " and s.contract_id = ?"
		args = append(args, str(m, "contractId"))
	}
	if v, ok := m["oppositeName"]; ok && v != nil {
		where += " and d.opposite_name like ?"
		args = append(args, "%"+str(m, "oppositeName")+"%")
	}
	if v, ok := m["transTime"]; ok && v != nil {
		if strings.TrimSpace(str(m, "transTime")) != "" {
			t, ok2 := parseLocalDate(str(m, "transTime"))
			if !ok2 {
				platform.Unknown(w)
				return
			}
			where += " and d.trans_time like ?"
			args = append(args, "%"+t.Format("2006-01-02")+"%")
		}
	}
	if v, ok := m["bankName"]; ok && v != nil {
		where += " and d.bank_name like ?"
		args = append(args, "%"+str(m, "bankName")+"%")
	}
	full := where + " and CONCAT(IFNULL(d.system_code, ''), IFNULL(d.account_name, ''), IFNULL(d.account, ''), IFNULL(d.trans_time, ''), IFNULL(d.money, ''), IFNULL(d.opposite_name, ''), IFNULL(d.opposite_account, ''), IFNULL(d.abs_content, '') ) like ?"
	args = append(args, searchStrLike(m))
	current, size := pageOf(m)
	q := "select d.*, (d.money-d.khxje) sjhxje from chain_capital_detailed d, chain_capital_detail_settlement s " + full + " order by trans_time desc"
	countQ := "SELECT count(1) FROM chain_capital_detailed d, chain_capital_detail_settlement s " + full
	ip, err := paged(s.pick(r), q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entDetailed.rows(recs)
	platform.Success(w, ip)
}

var _ = sql.ErrNoRows
var _ = json.Valid
var _ = strconv.Itoa
