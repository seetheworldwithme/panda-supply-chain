// Package capital 五期资金结算模块:saas-capital 全量端点。
// Java 源:saas-capital/.../business/controller/ChainCapitalDetailedController、
// ChainCapitalCgDetailedController、ChainCapitalXsDetailedController、
// ChainCapitalSettlementController、ChainCapitalCgSettlementController、
// ChainCapitalXsSettlementController 及 payment/controller/ChainPaymentController;
// Cg/Xs 变体与主控制器共用同一 service,仅路径与个别端点差异。
//
// 行为保持要点(勿"修复"):
//   - fill 实情(2026-11-17 源码确认):仅 ChainCapitalBkcall/ChainCapitalTransinfo 带
//     createTime(INSERT)+updateTime(INSERT_UPDATE) 双 fill;Detailed/Settlement 等全部
//     无 fill → InsertNoFill/UpdateByIDNoFill,update_time 不自动刷新。
//   - selectQuantity/countSett/countDetail 与 /saas-bill 统计 9 条已由一期 dashboard
//     注册,此处不重复;语义校准在 dashboard 包(五期)。
//   - listByPage 的 searchSql 在 Java 是 #{maps.x} 参数化,Go 侧用 ? 占位参数化保真
//     (值含引号不破坏 SQL);mapper ${} 内插点(leftDate/searchStr)照原样内插。
//   - 空参 NPE/拆箱 999 逐个复刻:selectQuat 空 timeType、countSett/countTaxRate 空
//     startTime/endTime、delBatch 空 ids、relSmentsBills 空 sments 等。
package capital

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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

// Register 挂载五期资金路由(全部 POST)。selectQuantity/countSett/countDetail 已由
// dashboard 包注册,此处不重复。
func Register(mux *http.ServeMux, s *Service) {
	// ── /chainCapitalDetailed 主收付款单(21 条;countDetail 一期 dashboard 版已由本期精确实现接管) ──
	for _, p := range []string{"/chainCapitalDetailed/list", "/chainCapitalDetailed/listByPage",
		"/chainCapitalDetailed/getById", "/chainCapitalDetailed/add", "/chainCapitalDetailed/deleteById",
		"/chainCapitalDetailed/updateById", "/chainCapitalDetailed/saveBatchDetails",
		"/chainCapitalDetailed/importAlone", "/chainCapitalDetailed/importExcelSk",
		"/chainCapitalDetailed/importExcelFk", "/chainCapitalDetailed/exportExcel",
		"/chainCapitalDetailed/delBatch", "/chainCapitalDetailed/balancingAccounts",
		"/chainCapitalDetailed/cancelBalancing", "/chainCapitalDetailed/getSystemCode",
		"/chainCapitalDetailed/getSments", "/chainCapitalDetailed/getSumPrice",
		"/chainCapitalDetailed/getPtSumPrice", "/chainCapitalDetailed/getBanks",
		"/chainCapitalDetailed/getDetails", "/chainCapitalDetailed/countDetail"} {
		mux.HandleFunc("POST /saas-capital"+p, s.handleDetailed)
	}
	// ── /chainCapitalCgDetailed /chainCapitalXsDetailed(采购付款单/销售收款单,共用主 service,少 2 端点) ──
	for _, base := range []string{"/chainCapitalCgDetailed", "/chainCapitalXsDetailed"} {
		for _, p := range []string{"/list", "/listByPage", "/getById", "/add", "/deleteById", "/updateById",
			"/saveBatchDetails", "/importAlone", "/importExcelSk", "/importExcelFk", "/exportExcel",
			"/delBatch", "/getSystemCode", "/getSments", "/getSumPrice", "/countDetail", "/getBanks", "/getDetails"} {
			mux.HandleFunc("POST /saas-capital"+base+p, s.handleDetailed)
		}
	}
	// ── /chainCapitalSettlement 主结算单(19 条;selectQuantity/countSett 一期 dashboard 版已由本期精确实现接管) ──
	for _, p := range []string{"/chainCapitalSettlement/list", "/chainCapitalSettlement/listByPage",
		"/chainCapitalSettlement/getById", "/chainCapitalSettlement/add", "/chainCapitalSettlement/deleteById",
		"/chainCapitalSettlement/updateById", "/chainCapitalSettlement/saveBatchSments",
		"/chainCapitalSettlement/importAlone", "/chainCapitalSettlement/importExcel",
		"/chainCapitalSettlement/exportExcel", "/chainCapitalSettlement/delBatch",
		"/chainCapitalSettlement/getSystemCode", "/chainCapitalSettlement/getDetails",
		"/chainCapitalSettlement/getBills", "/chainCapitalSettlement/getSumPrice",
		"/chainCapitalSettlement/getPtSumPrice", "/chainCapitalSettlement/relSmentsBills",
		"/chainCapitalSettlement/selectQuantity", "/chainCapitalSettlement/countSett"} {
		mux.HandleFunc("POST /saas-capital"+p, s.handleSettlement)
	}
	// ── /chainCapitalCgSettlement /chainCapitalXsSettlement(采购/销售结算单,共用主 service,少 getPtSumPrice) ──
	for _, base := range []string{"/chainCapitalCgSettlement", "/chainCapitalXsSettlement"} {
		for _, p := range []string{"/list", "/listByPage", "/getById", "/add", "/deleteById", "/updateById",
			"/saveBatchSments", "/importAlone", "/importExcel", "/exportExcel", "/delBatch",
			"/getSystemCode", "/getDetails", "/getBills", "/selectQuantity", "/getSumPrice",
			"/relSmentsBills", "/countSett"} {
			mux.HandleFunc("POST /saas-capital"+base+p, s.handleSettlement)
		}
	}
	// ── 小控制器套件(账户/银行回单/交易/核销记录/核销明细/明细结算对应/账单结算记录) ──
	for _, base := range []string{"/chainCapitalAccounts", "/chainCapitalBkcall", "/chainCapitalTransinfo",
		"/chainCapitalCount", "/chainCapitalCountMx", "/chainCapitalDetailSettlement", "/chainBillSettRec"} {
		for _, p := range []string{"/list", "/listByPage", "/getById", "/add", "/deleteById", "/updateById",
			"/importExcel", "/exportExcel", "/delBatch"} {
			mux.HandleFunc("POST /saas-capital"+base+p, s.handleSmall)
		}
	}
	mux.HandleFunc("POST /saas-capital/chainBillSettRec/relieve", s.handleSmall)
	// ── 支付(招行 openbank) ──
	for _, p := range []string{"/payment", "/confirmPay", "/cancelPay", "/refund",
		"/orderDetails", "/billInfo", "/transMsgCallBack"} {
		mux.HandleFunc("POST /saas-capital/chainPayment"+p, s.handlePayment)
	}
}

// pick 按请求头 applytype 选库(资金链路同 dashboard/logistics)
func (s *Service) pick(r *http.Request) *sql.DB {
	if r.Header.Get("applytype") == "zjjg" {
		return s.DB.Zjjg
	}
	return s.DB.Chain
}

// ── 共享小工具(与 logistics 包同款语义;包内私有保持各期独立) ──

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

func pick(m map[string]any, cols ...string) map[string]any {
	u := platform.Underscore(m)
	out := make(map[string]any, len(cols))
	for _, c := range cols {
		if v, ok := u[c]; ok {
			out[c] = v
		}
	}
	return out
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

// searchStrOf Java 实体绑定:body 无 searchStr 键 → getSearchStr() null;
// mapper #{maps.searchStr} 参数化输出 SQL NULL → like NULL 恒假(整查空);有键(含空串)→ %v%。
// 注意与四期 ${searchStr} 内插形态(缺失→"null"/空串)不同,此处全部是 #{} 参数化绑定。
func searchStrOf(m map[string]any) (string, bool) {
	v, ok := m["searchStr"]
	if !ok || v == nil {
		return "", false
	}
	return str(m, "searchStr"), true
}

func currentUser(r *http.Request) *auth.CurrentUser { return auth.From(r.Context()) }

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

// camel 列名→驼峰(借 platform.Camelize 单键变换)
func camel(s string) string {
	out := platform.Camelize(map[string]any{s: nil})
	for k := range out {
		return k
	}
	return s
}

// parseLocalDate Jackson LocalDate 默认只收 ISO yyyy-MM-dd(空串绑 null 由调用方处理)
func parseLocalDate(s string) (time.Time, bool) {
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	return t, err == nil
}

// parseLocalDateTime Jackson LocalDateTime 默认 ISO-T
func parseLocalDateTime(s string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02T15:04:05.000", "2006-01-02T15:04"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("bad datetime: %s", s)
}

// validateDT 时间绑定校验:空格 pattern 列只收空格格式,无注解 LocalDateTime 列只收 ISO-T;
// LocalDate 列只收 yyyy-MM-dd。不匹配 → HttpMessageNotReadableException → 999。
// 空串对 Integer/LocalDate(Time)/Boolean 字段绑 null(四期实证⑥),不算失败。
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

// autoGenericCode 对照 StringHelper.autoGenericCode:parseInt 后 %0Nd 补零
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

// systemCodeNext 对照 getSystemCode 族:空表→前缀+start 补零 7 位;已有→前 20 位截断
// 后整数 +1(substring(0,20) 越界或后缀非数字 → Java 异常 → 999)。
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

// nextSystemCodeInMem 对照 getNextSystemCode:内存编号自增(长度/数字校验失败 → 999)
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

// logOp 写 chain_common_log(ChainCommonLog 无 fill:create_time 走列默认;companyId 构造器未带)
func logOp(d *sql.DB, u *auth.CurrentUser, operation string, opType int) {
	_, _ = d.Exec("INSERT INTO chain_common_log (user_id, user_name, operation, operation_type) VALUES (?, ?, ?, ?)",
		strconv.FormatInt(u.UserID, 10), u.Username, operation, opType)
}

// ── 实体形状元数据(同 logistics entShape) ──

type entShape struct {
	table     string
	cols      []string
	decimals  []string
	bools     []string
	tTimes    []string
	dateOnly  []string // Java LocalDate 字段(datetime 列 → 只输出 yyyy-MM-dd)
	spaceT    []string // 空格 pattern 已是 ScanMaps 默认格式,仅占位说明
	renames   map[string]string
	carriers  []string
	primFalse []string   // Java 原生 boolean 载体键(缺省 false 而非 null,如 khxFlag)
	fillIns   bool       // 带双 fill 的实体(bkcall/transinfo):insert 自动 create+update
	fields    []impField // Excel 导入导出列(Java @ExcelColumn 声明序)
}

// selectList MP list() 风格的显式列清单(SELECT * 行序与列清单不同,并列 trans_time 行序随执行计划)
func (e entShape) selectList() string {
	return strings.Join(e.cols, ",")
}

func (e entShape) has(list []string, c string) bool {
	for _, x := range list {
		if x == c {
			return true
		}
	}
	return false
}

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
		// 键白名单:只保留实体列(renames 后)+载体键,剔除 SELECT * 的多余列
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
		// mapper 显式列外的实体键补 null(Java 实体全键序列化);键走 renames 后形态
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

func nowSpace() string { return time.Now().Format("2006-01-02 15:04:05") }

// nowTNanos 无 @JsonFormat 的 updateTime 序列化为 ISO T 带微秒(MateMetaObjectHandler LocalDateTime.now())
func nowTNanos() string { return time.Now().Format("2006-01-02T15:04:05.000000") }

// allEqCond 同 contract 版:键须落在列白名单(null 跳过,白名单外 → 999)
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

// paged count + 窗口查询(MP IPage 形状)
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
