// Package bill 五期票据模块:saas-bill 全量端点(chainBill 主+Cg/Xs、chainBillRate、
// chainBillSettlement)。
// Java 源:saas-bill/.../business/controller/ChainBillController、ChainCgBillController、
// ChainXsBillController、ChainBillRateController、ChainBillSettlementController。
//
// 行为保持要点(勿"修复"):
//   - selJxzb 恒空数组(Java bug:selMinMonth 只查 month,startTime 误用 openDate 恒 null
//     → SQL >= NULL 恒假);selZeroTax/countTaxRate 除零 → ArithmeticException 999。
//   - 一期 dashboard 注册的 chainBill 7 条统计与 chainBillRate 2 条由本期精确实现接管,
//     dashboard 侧删除(float 近似/空行跳过等非保真实现)。
//   - importExcelJx/Xx 无 fail 分支:errInfos 进 ImportResult 响应(saveIrt)。
//   - saveBatch:主 Bill 响应 errInfos 数组,Cg/Xs 恒 "新增成功"(控制器差异)。
package bill

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type Service struct {
	DB *db.DB
}

func NewService(database *db.DB) *Service { return &Service{DB: database} }

// Register 挂载五期票据路由(全部 POST;主 Bill 统计 7 条与 BillRate 2 条原 dashboard 注册移交至此)
func Register(mux *http.ServeMux, s *Service) {
	// ── /chainBill 主(25 条) ──
	for _, p := range []string{"/chainBill/list", "/chainBill/listByPage", "/chainBill/getById",
		"/chainBill/add", "/chainBill/deleteById", "/chainBill/updateById", "/chainBill/saveBatch",
		"/chainBill/importExcelJx", "/chainBill/importExcelXx", "/chainBill/exportExcel",
		"/chainBill/delBatch", "/chainBill/getSystemCode", "/chainBill/getSments",
		"/chainBill/getSumPrice", "/chainBill/selectBusCount", "/chainBill/countBill",
		"/chainBill/countTaxRate", "/chainBill/selZeroTax", "/chainBill/selJxzb",
		"/chainBill/searchRelBill", "/chainBill/countBuses", "/chainBill/countMonBuses",
		"/chainBill/countPtTax", "/chainBill/countYearBuses"} {
		mux.HandleFunc("POST /saas-bill"+p, s.handleBill)
	}
	// ── /chainCgBill /chainXsBill(进项/销项发票,共用主 service;saveBatch 响应"新增成功";
	//    少 countYearBuses;selectBusCount/countTaxRate/selJxzb/selZeroTax/countBuses/countMonBuses/countPtTax 在此注册) ──
	for _, base := range []string{"/chainCgBill", "/chainXsBill"} {
		for _, p := range []string{"/list", "/listByPage", "/getById", "/add", "/deleteById", "/updateById",
			"/saveBatch", "/importExcelJx", "/importExcelXx", "/exportExcel", "/delBatch",
			"/getSystemCode", "/getSments", "/getSumPrice", "/selectBusCount", "/countBill",
			"/countTaxRate", "/selZeroTax", "/selJxzb", "/searchRelBill", "/countBuses",
			"/countMonBuses", "/countPtTax"} {
			mux.HandleFunc("POST /saas-bill"+base+p, s.handleBill)
		}
	}
	// ── /chainBillRate(13 条) ──
	for _, p := range []string{"/chainBillRate/list", "/chainBillRate/listByPage", "/chainBillRate/getById",
		"/chainBillRate/add", "/chainBillRate/deleteById", "/chainBillRate/updateById",
		"/chainBillRate/importExcel", "/chainBillRate/exportExcel", "/chainBillRate/delBatch",
		"/chainBillRate/saveBatch", "/chainBillRate/getSumPrice", "/chainBillRate/countMonBuses"} {
		mux.HandleFunc("POST /saas-bill"+p, s.handleRate)
	}
	// ── /chainBillSettlement(9 条) ──
	for _, p := range []string{"/chainBillSettlement/list", "/chainBillSettlement/listByPage",
		"/chainBillSettlement/getById", "/chainBillSettlement/add", "/chainBillSettlement/deleteById",
		"/chainBillSettlement/updateById", "/chainBillSettlement/importExcel",
		"/chainBillSettlement/exportExcel", "/chainBillSettlement/delBatch"} {
		mux.HandleFunc("POST /saas-bill"+p, s.handleSett)
	}
}

// pick 按请求头 applytype 选库
func (s *Service) pick(r *http.Request) *sql.DB {
	if r.Header.Get("applytype") == "zjjg" {
		return s.DB.Zjjg
	}
	return s.DB.Chain
}

// ── 共享小工具(与 capital 包同款语义,包内私有) ──

func bodyMap(r *http.Request) map[string]any {
	var m map[string]any
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	_ = dec.Decode(&m)
	if m == nil {
		return map[string]any{}
	}
	return m
}

func bodyList(r *http.Request) ([]any, bool) {
	var list []any
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	_ = dec.Decode(&list)
	return list, list != nil
}

func str(m map[string]any, k string) string {
	v, ok := m[k]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func num(v any) float64 {
	switch t := v.(type) {
	case nil:
		return 0
	case json.Number:
		f, _ := t.Float64()
		return f
	case []byte:
		f, _ := strconv.ParseFloat(string(t), 64)
		return f
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case int64:
		return float64(t)
	case float64:
		return t
	case bool:
		if t {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func int64Of(m map[string]any, k string) (int64, bool) {
	v, ok := m[k]
	if !ok || v == nil {
		return 0, false
	}
	if f, ok := v.(json.Number); ok {
		n, err := f.Int64()
		return n, err == nil
	}
	if f, ok := v.(float64); ok {
		return int64(f), true
	}
	n, err := strconv.ParseInt(str(m, k), 10, 64)
	return n, err == nil
}

func btsStr(v any) string {
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

func camel(s string) string {
	out := platform.Camelize(map[string]any{s: nil})
	for k := range out {
		return k
	}
	return s
}

func isEmptyStr(v any) bool {
	s, ok := v.(string)
	return ok && strings.TrimSpace(s) == ""
}

func pageOf(m map[string]any) (current, size int64) {
	current, size = 1, 10
	if p, ok := m["page"].(map[string]any); ok {
		if v := num(p["current"]); v > 0 {
			current = int64(v)
		}
		if v := num(p["size"]); v > 0 {
			size = int64(v)
		}
	}
	return
}

func parseLocalDate(s string) (time.Time, bool) {
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	return t, err == nil
}

func parseLocalDateTime(s string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02T15:04:05.000", "2006-01-02T15:04"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("bad datetime: %s", s)
}

// validateDT 空 pattern 列只收空格格式,LocalDate 列只收 ISO;空串绑 null 不算失败
func validateDT(m map[string]any, spaceCols, dateCols, tCols []string) bool {
	for _, c := range spaceCols {
		k := camel(c)
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		if s, isStr := v.(string); isStr && strings.TrimSpace(s) == "" {
			continue
		}
		if _, err := time.ParseInLocation("2006-01-02 15:04:05", str(m, k), time.Local); err != nil {
			return false
		}
	}
	for _, c := range dateCols {
		k := camel(c)
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		if s, isStr := v.(string); isStr && strings.TrimSpace(s) == "" {
			continue
		}
		if _, ok2 := parseLocalDate(str(m, k)); !ok2 {
			return false
		}
	}
	for _, c := range tCols {
		k := camel(c)
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		if s, isStr := v.(string); isStr && strings.TrimSpace(s) == "" {
			continue
		}
		if _, err := parseLocalDateTime(str(m, k)); err != nil {
			return false
		}
	}
	return true
}

func autoGenericCode(code string, width int) string {
	n, err := strconv.Atoi(strings.TrimSpace(code))
	if err != nil {
		return code
	}
	s := strconv.Itoa(n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

func systemCodeNext(d *sql.DB, table, likePrefix, start string) (string, bool) {
	var cur sql.NullString
	err := d.QueryRow("SELECT IFNULL(max(system_code),'') FROM "+table+" WHERE system_code LIKE ?", likePrefix+"%").Scan(&cur)
	if err != nil {
		return "", false
	}
	if cur.String == "" {
		return likePrefix + autoGenericCode(start, 7), true
	}
	if len(cur.String) < 20 {
		return "", false
	}
	n, err := strconv.Atoi(cur.String[20:])
	if err != nil {
		return "", false
	}
	return cur.String[:20] + autoGenericCode(strconv.Itoa(n+1), 7), true
}

func nextSystemCodeInMem(code string) (string, bool) {
	if len(code) < 20 {
		return "", false
	}
	n, err := strconv.Atoi(code[20:])
	if err != nil {
		return "", false
	}
	return code[:20] + autoGenericCode(strconv.Itoa(n+1), 7), true
}

func logOp(d *sql.DB, u *auth.CurrentUser, operation string, opType int) {
	_, _ = d.Exec("INSERT INTO chain_common_log (user_id, user_name, operation, operation_type) VALUES (?, ?, ?, ?)",
		strconv.FormatInt(u.UserID, 10), u.Username, operation, opType)
}

func currentUser(r *http.Request) *auth.CurrentUser { return auth.From(r.Context()) }

// allEqCond 键须落在列白名单(null 跳过,白名单外 → 999)
func allEqCond(m map[string]any, cols []string) (map[string]any, bool) {
	allow := make(map[string]bool, len(cols))
	for _, c := range cols {
		allow[c] = true
	}
	cond := map[string]any{}
	for k, v := range m {
		if v == nil {
			continue
		}
		if !allow[k] {
			return nil, false
		}
		cond[k] = v
	}
	return cond, true
}

func paged(d *sql.DB, q, countQ string, args []any, current, size int64) (map[string]any, error) {
	var total int64
	if err := d.QueryRow(countQ, args...).Scan(&total); err != nil {
		return nil, err
	}
	off := (current - 1) * size
	rows, err := d.Query(q+" LIMIT ? OFFSET ?", append(append([]any{}, args...), size, off)...)
	if err != nil {
		return nil, err
	}
	ms, err := platform.ScanMaps(rows)
	if ms == nil {
		ms = []map[string]any{}
	}
	return platform.IPage(platform.CamelizeRows(ms), total, current, size), nil
}

// ── 实体形状 ──

type entShape struct {
	table     string
	cols      []string
	decimals  []string
	bools     []string
	tTimes    []string
	dateOnly  []string
	spaceT    []string
	renames   map[string]string
	carriers  []string
	primFalse []string
	fields    []impField
}

func (e entShape) has(list []string, c string) bool {
	for _, x := range list {
		if x == c {
			return true
		}
	}
	return false
}

func (e entShape) selectList() string { return strings.Join(e.cols, ",") }

func (e entShape) rows(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		for _, c := range e.cols {
			key := camel(c)
			if e.has(e.decimals, c) {
				if s, ok := r[key].(string); ok && s != "" {
					r[key] = json.Number(s)
				}
			}
			if e.has(e.bools, c) {
				if v, ok := r[key]; ok && v != nil {
					r[key] = num(v) != 0
				}
			}
			if e.has(e.tTimes, c) {
				if s, ok := r[key].(string); ok && len(s) == 19 {
					r[key] = s[:10] + "T" + s[11:]
				}
			}
			if e.has(e.dateOnly, c) {
				if s, ok := r[key].(string); ok && len(s) >= 10 {
					r[key] = s[:10]
				}
			}
		}
		for from, to := range e.renames {
			if v, ok := r[from]; ok {
				r[to] = v
				delete(r, from)
			}
		}
		keep := map[string]bool{}
		for _, c := range e.cols {
			k := camel(c)
			if to, has := e.renames[k]; has {
				k = to
			}
			keep[k] = true
		}
		for _, k := range e.carriers {
			keep[k] = true
		}
		for _, k := range e.primFalse {
			keep[k] = true
		}
		for k := range r {
			if !keep[k] {
				delete(r, k)
			}
		}
		for _, c := range e.cols {
			k := camel(c)
			if to, has := e.renames[k]; has {
				k = to
			}
			if _, ok := r[k]; !ok {
				r[k] = nil
			}
		}
		fillNulls(r, e.carriers)
		for _, k := range e.primFalse {
			if r[k] == nil {
				r[k] = false
			}
		}
	}
	return out
}

func (e entShape) echoKeys() []string {
	keys := make([]string, 0, len(e.cols)+len(e.carriers))
	for _, c := range e.cols {
		k := camel(c)
		if to, has := e.renames[k]; has {
			k = to
		}
		keys = append(keys, k)
	}
	return append(keys, e.carriers...)
}

func fillNulls(m map[string]any, keys []string) {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			m[k] = nil
		}
	}
}

// echoOf 内存实体回显 + 原生 boolean 载体键缺省 false
func echoOf(ent *entShape, m map[string]any, extra map[string]any) map[string]any {
	out := platform.EntityEcho(m, ent.echoKeys(), extra)
	for _, k := range ent.primFalse {
		if out[k] == nil {
			out[k] = false
		}
	}
	return out
}

var colNameRe = regexp.MustCompile(`^[a-z_][a-zA-Z0-9_]*$`)

// writeVals body 实体 → 入库列(Underscore + renames 反向)
func writeVals(ent *entShape, m map[string]any) map[string]any {
	u := platform.Underscore(m)
	out := make(map[string]any, len(u)+len(ent.renames))
	// 载体键(page/searchStr/relSearch…)和非列名键不入库:Java 忽略非实体字段;
	// 照单全收轻则 1054 报错,重则非法键名触发 CheckCols panic
	for _, k := range append(ent.carriers, "rel_search", "relSearch") {
		delete(u, k)
	}
	for k, v := range u {
		if !colNameRe.MatchString(k) {
			continue
		}
		out[k] = v
	}
	for camelCol, jsonKey := range ent.renames {
		snake := platform.Underscore(map[string]any{camelCol: nil})
		for k := range snake {
			delete(out, k)
		}
		delete(out, camelCol)
		delete(out, strings.ToLower(camelCol))
		if v, ok := m[jsonKey]; ok {
			out[camelCol] = v
		}
	}
	return out
}

// getSqlReflect 对照 StringHelper.getSql(t, unstr):@ExcelColumn 字段值非 null 非 "0"
// 且不在 unstr 排除表 → " and <snake> like ?"(值 Integer/LocalDate 原样,其余 %v%)
func getSqlReflect(m map[string]any, ent *entShape, unstr []string) (string, []any, bool) {
	where := "where 1=1"
	var args []any
	for _, c := range ent.cols {
		k := camel(c)
		if c == "id" {
			continue // id 是 @ExcelColumn 但 StringHelper 拼条件同样包含——保留注入语义
		}
		excluded := false
		for _, u := range unstr {
			if u == k {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		if s, isStr := v.(string); isStr && s == "0" {
			continue
		}
		where += " and " + c + " like ?"
		switch v.(type) {
		case json.Number:
			// Integer/BigDecimal 原值(getMaps 对 Integer 原样;BigDecimal %v% 包裹——
			// ChainBill 的 count/unitPrice/baseMoney/taxAmount/includMoney 是 BigDecimal → %v%)
			if isIntCol(c) {
				args = append(args, str(m, k))
			} else {
				args = append(args, "%"+str(m, k)+"%")
			}
		default:
			args = append(args, "%"+str(m, k)+"%")
		}
	}
	return where, args, true
}

func isIntCol(c string) bool {
	switch c {
	case "id", "company_id", "contract_id", "type", "relation_id":
		return true
	}
	return false
}

func nowSpace() string  { return time.Now().Format("2006-01-02 15:04:05") }
func nowTNanos() string { return time.Now().Format("2006-01-02T15:04:05.000000") }
