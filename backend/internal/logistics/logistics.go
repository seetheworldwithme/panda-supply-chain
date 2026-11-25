// Package logistics 四期物流模块:saas-logistics 全量端点。
// Java 源:saas-logistics/.../business/controller/ChainLogisticsController、
// ChainLogisticsBillController、ChainLogisticsContractController、
// ChainLogisticsThirdWaybillController、ChainLogisticsWeighingListController。
//
// 行为保持要点(勿"修复"):
//   - fill 实情(与任务前提相反,以 Java 源为准,2026-11-17 实测):ChainLogistics(createTime=INSERT,
//     updateTime=INSERT_UPDATE)与 ChainLogisticsWeighingList(同)带 fill;ChainLogisticsBill 仅
//     updateTime 带 fill(插库 create_time 走列默认,update_time 自动);ChainLogisticsContract/
//     ChainLogisticsThirdWaybill 无 fill。写库按此分 Insert/InsertNoFill。
//   - 第三方对接/过磅单控制器的 add/updateById 返回 success() 不带实体回显;运单/订单/运单合同
//     对应表回显内存实体(EntityEcho,body 未传键 null)。
//   - 统计端点 countBuses/countMonBuses/countLogist 已由一期 dashboard 包注册,本期只校准语义
//     (见 dashboard/logistics.go),不重注册。
//   - bill.selectQuat 多行结果 Java 抛 TooManyResultsException → 999,照实复刻。
//   - countLogist timeType=4 分支:补月循环 getStartTime() NPE → 恒 999,照实复刻。
//   - Excel 导入按 ExcelUtils.readExcel 语义逐条复刻(表头包含匹配、NotBankColumn 必填、
//     DictColumn 字典翻译、notNullMap 首行后耗尽的原版 bug、4 参版本遍历全部 sheet)。
package logistics

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

// Register 挂载四期物流路由(全部 POST)。countBuses/countMonBuses(主表)与
// chainLogisticsBill/countLogist 由 dashboard 包注册,此处不重复。
func Register(mux *http.ServeMux, s *Service) {
	// ── /chainLogistics 运单(11 条,统计 2 条在 dashboard) ──
	for _, p := range []string{"/chainLogistics/list", "/chainLogistics/listByPage", "/chainLogistics/getById",
		"/chainLogistics/add", "/chainLogistics/deleteById", "/chainLogistics/updateById",
		"/chainLogistics/importExcel", "/chainLogistics/exportExcel", "/chainLogistics/delBatch",
		"/chainLogistics/updTransLocation", "/chainLogistics/saveBatchLogistics"} {
		mux.HandleFunc("POST /saas-logistics"+p, s.handleLo)
	}
	// ── /chainLogisticsBill 订单(17 条,统计 countLogist 在 dashboard) ──
	for _, p := range []string{"/chainLogisticsBill/list", "/chainLogisticsBill/listByPage", "/chainLogisticsBill/getById",
		"/chainLogisticsBill/add", "/chainLogisticsBill/deleteById", "/chainLogisticsBill/updateById",
		"/chainLogisticsBill/importExcel", "/chainLogisticsBill/exportExcel", "/chainLogisticsBill/delBatch",
		"/chainLogisticsBill/getSystemCode", "/chainLogisticsBill/selectQuat", "/chainLogisticsBill/getContracts",
		"/chainLogisticsBill/getSumPrice", "/chainLogisticsBill/countBuses", "/chainLogisticsBill/countMonBuses",
		"/chainLogisticsBill/selOrders", "/chainLogisticsBill/selAllOrders"} {
		mux.HandleFunc("POST /saas-logistics"+p, s.handleBill)
	}
	// ── /chainLogisticsContract 运单/合同对应表(9 条) ──
	for _, p := range []string{"/chainLogisticsContract/list", "/chainLogisticsContract/listByPage", "/chainLogisticsContract/getById",
		"/chainLogisticsContract/add", "/chainLogisticsContract/deleteById", "/chainLogisticsContract/updateById",
		"/chainLogisticsContract/importExcel", "/chainLogisticsContract/exportExcel", "/chainLogisticsContract/delBatch"} {
		mux.HandleFunc("POST /saas-logistics"+p, s.handleLc)
	}
	// ── /chainLogisticsThirdWaybill 对接运单(12 条) ──
	for _, p := range []string{"/chainLogisticsThirdWaybill/list", "/chainLogisticsThirdWaybill/listByPage", "/chainLogisticsThirdWaybill/getById",
		"/chainLogisticsThirdWaybill/add", "/chainLogisticsThirdWaybill/addBatch", "/chainLogisticsThirdWaybill/deleteById",
		"/chainLogisticsThirdWaybill/updateById", "/chainLogisticsThirdWaybill/updateBatch",
		"/chainLogisticsThirdWaybill/importExcel", "/chainLogisticsThirdWaybill/exportExcel",
		"/chainLogisticsThirdWaybill/delBatch", "/chainLogisticsThirdWaybill/clearData"} {
		mux.HandleFunc("POST /saas-logistics"+p, s.handleTw)
	}
	// ── 物流页依赖的 saas-common 字典端点(四期) ──
	mux.HandleFunc("POST /saas-common/chainCommonTempmx/formatlist", s.tempmxFormatlist)

	// ── /chainLogisticsWeighingList 过磅单(10 条) ──
	for _, p := range []string{"/chainLogisticsWeighingList/list", "/chainLogisticsWeighingList/listByPage", "/chainLogisticsWeighingList/getById",
		"/chainLogisticsWeighingList/add", "/chainLogisticsWeighingList/deleteById", "/chainLogisticsWeighingList/updateById",
		"/chainLogisticsWeighingList/importExcel", "/chainLogisticsWeighingList/exportExcel",
		"/chainLogisticsWeighingList/delBatch", "/chainLogisticsWeighingList/saveBatchWeighingList"} {
		mux.HandleFunc("POST /saas-logistics"+p, s.handleWl)
	}
}

// pick 按请求头 applytype 选库(物流服务链路同 dashboard;Java 侧物流无监管库路由,保留同款以防 applytype)
func (s *Service) pick(r *http.Request) *sql.DB {
	if r.Header.Get("applytype") == "zjjg" {
		return s.DB.Zjjg
	}
	return s.DB.Chain
}

// ── 共享小工具(与 contract/dashboard 同款,包内私有;UseNumber 保 DECIMAL 字面) ──

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

// pickDB 白名单挑写入列(键为下划线列)
func pick(m map[string]any, cols ...string) map[string]any {
	u := platform.Underscore(m)
	out := make(map[string]any, len(cols))
	// 列名大小写不敏感匹配:Invoicing_time 这类驼峰列名原先永远匹配不上,接口传入的 invoicingTime 被静默丢弃
	lower := make(map[string]any, len(u))
	for k, v := range u {
		lower[strings.ToLower(k)] = v
	}
	for _, c := range cols {
		if v, ok := lower[strings.ToLower(c)]; ok {
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

func searchStrOf(m map[string]any) string {
	if v, ok := m["searchStr"]; !ok || v == nil {
		return "" // MyBatis ${searchStr} 对 null 输出空串 → like '%%'(8765 实测;非 "null" 字面)
	}
	return sqlEscape(strings.TrimSpace(str(m, "searchStr"))) // 调用方把它字面拼进 like '...':不转义即 SQL 注入
}

// sqlEscape 转义要字面拼进 SQL 字符串的值(反斜杠、单引号)。Java mapper 的 ${} 是字面注入,
// 照搬就是 SQL 注入;良性输入转义后拼出的 SQL 与原先逐字相同。
func sqlEscape(v string) string {
	return strings.ReplaceAll(strings.ReplaceAll(v, "\\", "\\\\"), "'", "''")
}

// sqlString 值 → 带引号的 MySQL 字符串字面量(用于原先裸拼接的等值条件)。
func sqlString(v string) string { return "'" + sqlEscape(v) + "'" }

// paged count + 窗口查询(MP IPage 形状;countQ 不带 order by)
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

func fillNulls(m map[string]any, keys []string) {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			m[k] = nil
		}
	}
}

func currentUser(r *http.Request) *auth.CurrentUser { return auth.From(r.Context()) }

// fileKeys ChainCommonFile 实体全键(回显 filelist 元素时按实体整形)
var fileKeys = []string{"id", "companyId", "fileRealName", "fileName", "fileDescription", "filePath",
	"ftype", "fileType", "ownerType", "ownerId", "relationId", "relationType", "effective",
	"transStatus", "associationId", "isCompressedFile", "sourceId", "compressedId", "createTime",
	"updateTime", "remark1", "remark2", "remark3", "remark4", "remark5", "searchStr", "configFilePath",
	"page", "ids"}

// shapeFilelist 回显体 filelist 元素整形(Java 反序列化的 ChainCommonFile 对象全键序列化)
func shapeFilelist(m map[string]any) {
	if files, ok := m["filelist"].([]any); ok {
		out := make([]any, 0, len(files))
		for _, it := range files {
			if f, ok := it.(map[string]any); ok {
				out = append(out, platform.EntityEcho(f, fileKeys, nil))
			} else {
				out = append(out, it)
			}
		}
		m["filelist"] = out
	}
}

// shapeFilelistSvc service 已给元素设置 6 固定字段后的回显(Java 内存实体含这些值)
func shapeFilelistSvc(m map[string]any, relationID int64, ownerUserID int64) {
	if files, ok := m["filelist"].([]any); ok {
		out := make([]any, 0, len(files))
		for _, it := range files {
			if f, ok := it.(map[string]any); ok {
				extra := map[string]any{"fileType": "1", "ownerType": "1",
					"ownerId":      json.Number(strconv.FormatInt(ownerUserID, 10)),
					"relationId":   json.Number(strconv.FormatInt(relationID, 10)),
					"relationType": 1, "effective": true}
				out = append(out, platform.EntityEcho(f, fileKeys, extra))
			} else {
				out = append(out, it)
			}
		}
		m["filelist"] = out
	}
}

// bodyList 解析 JSON 数组 body(批量端点;UseNumber 保 DECIMAL 字面)
func bodyList(r *http.Request) ([]any, bool) {
	var list []any
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	_ = dec.Decode(&list)
	return list, list != nil
}

// parseLocalDateTime Jackson LocalDateTime 默认 ISO(yyyy-MM-ddTHH:mm:ss);空格格式绑失败 → 999
func parseLocalDateTime(s string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02T15:04:05.000", "2006-01-02T15:04"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("bad datetime: %s", s)
}

// parseLocalDate Jackson LocalDate 默认 ISO(yyyy-MM-dd)
func parseLocalDate(s string) (time.Time, bool) {
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	return t, err == nil
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

func urlQueryEscape(s string) string { return url.QueryEscape(s) }

// validateDT Jackson 时间绑定校验:@JsonFormat(空格 pattern) 列只收空格格式,
// 无注解 LocalDateTime 列只收 ISO-T;不匹配 → HttpMessageNotReadableException → 999
func validateDT(m map[string]any, spaceCols, tCols []string) bool {
	for _, c := range spaceCols {
		v, ok := m[camel(c)]
		if !ok || v == nil {
			continue
		}
		if _, err := time.ParseInLocation("2006-01-02 15:04:05", str(m, camel(c)), time.Local); err != nil {
			return false
		}
	}
	for _, c := range tCols {
		v, ok := m[camel(c)]
		if !ok || v == nil {
			continue
		}
		if _, err := parseLocalDateTime(str(m, camel(c))); err != nil {
			return false
		}
	}
	return true
}

// autoGenericCode 对照 StringHelper.autoGenericCode:parseInt 后 %0Nd 补零
func autoGenericCode(code string, width int) string {
	n, err := strconv.Atoi(strings.TrimSpace(code))
	if err != nil {
		return code // Java NumberFormatException 由调用方先拦,此处防御
	}
	s := strconv.Itoa(n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

// nextSystemCode 对照 selectSystemCode:IFNULL(max(system_code),”) like 'prefix%'。
// 空表 → 前缀+start(补零 7 位);已有 → 前 20 位截断后整数 +1。
// 后缀非数字 → Java NumberFormatException → 999(ok=false)
func nextSystemCode(d *sql.DB, table, prefix, start string) (string, bool) {
	var cur sql.NullString
	err := d.QueryRow("SELECT IFNULL(max(system_code),'') FROM "+table+" WHERE system_code LIKE ?", prefix+"%").Scan(&cur)
	if err != nil {
		return "", false
	}
	if cur.String == "" {
		return prefix + autoGenericCode(start, 7), true
	}
	if len(cur.String) < 20 {
		return "", false // substring(0,20) 越界 → Java StringIndexOutOfBounds → 999
	}
	s1, s2 := cur.String[:20], cur.String[20:]
	n, err := strconv.Atoi(s2)
	if err != nil {
		return "", false
	}
	return s1 + autoGenericCode(strconv.Itoa(n+1), 7), true
}

// bumpSystemCode 对照 getNextSystemCode:前 20 位 + 后缀整数 +1(补零 7 位)
func bumpSystemCode(code string) (string, bool) {
	if len(code) < 20 {
		return "", false
	}
	n, err := strconv.Atoi(code[20:])
	if err != nil {
		return "", false
	}
	return code[:20] + autoGenericCode(strconv.Itoa(n+1), 7), true
}

// logOp 写 chain_common_log(ChainCommonLog 无 fill:create_time 走列默认;companyId 构造器未带,不写)
func logOp(d *sql.DB, u *auth.CurrentUser, operation string, opType int) {
	_, _ = d.Exec("INSERT INTO chain_common_log (user_id, user_name, operation, operation_type) VALUES (?, ?, ?, ?)",
		strconv.FormatInt(u.UserID, 10), u.Username, operation, opType)
}

// ── 实体形状元数据 ──

// entShape 实体 JSON 形状:DECIMAL 保字面小数位、tinyint(1)→布尔、指定列转 T 格式时间、
// 列名改写(Invoicing_time→invoicingTime)、载体键(null 回填)
type entShape struct {
	table      string
	cols       []string
	decimals   []string
	bools      []string
	tTimes     []string // 以 T 分隔时间格式输出的 datetime 列(snake;默认格式为空格)
	renames    map[string]string
	carriers   []string
	fillInsert bool // 实体带 createTime fill(insert 自动填 create+update 时间;仅 updateTime fill 的 bill 单独处理)
}

func (e entShape) has(list []string, c string) bool {
	for _, x := range list {
		if x == c {
			return true
		}
	}
	return false
}

// rows 数据行整形(SELECT * 行 → 实体 Jackson 形状)
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
		}
		for from, to := range e.renames {
			if v, ok := r[from]; ok {
				r[to] = v
				delete(r, from)
			}
		}
		fillNulls(r, e.carriers)
	}
	return out
}

// echoKeys 实体全键(列驼峰 + 载体键,renames 改写后)——add/updateById 内存实体回显
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

func camel(s string) string {
	out := platform.Camelize(map[string]any{s: nil})
	for k := range out {
		return k
	}
	return s
}

// nowSpace / nowTNanos 内存 fill 时间(Java MateMetaObjectHandler LocalDateTime.now():
// createTime @JsonFormat 空格秒级;updateTime 无 @JsonFormat 时 ISO T 带微秒)
func nowSpace() string { return time.Now().Format("2006-01-02 15:04:05") }

func nowTNanos() string { return time.Now().Format("2006-01-02T15:04:05.000000") }

// ── 链路路由分发(路径 → 各控制器 handler) ──

func (s *Service) handleLo(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/saas-logistics/chainLogistics/list":
		s.loList(w, r)
	case "/saas-logistics/chainLogistics/listByPage":
		s.loListByPage(w, r)
	case "/saas-logistics/chainLogistics/getById":
		s.loGetById(w, r)
	case "/saas-logistics/chainLogistics/add":
		s.loAdd(w, r)
	case "/saas-logistics/chainLogistics/deleteById":
		s.loDeleteById(w, r)
	case "/saas-logistics/chainLogistics/updateById":
		s.loUpdateById(w, r)
	case "/saas-logistics/chainLogistics/importExcel":
		s.loImportExcel(w, r)
	case "/saas-logistics/chainLogistics/exportExcel":
		s.loExportExcel(w, r)
	case "/saas-logistics/chainLogistics/delBatch":
		s.loDelBatch(w, r)
	case "/saas-logistics/chainLogistics/updTransLocation":
		s.loUpdTransLocation(w, r)
	case "/saas-logistics/chainLogistics/saveBatchLogistics":
		s.loSaveBatch(w, r)
	default:
		platform.Unknown(w)
	}
}

func (s *Service) handleBill(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/saas-logistics/chainLogisticsBill/list":
		s.billList(w, r)
	case "/saas-logistics/chainLogisticsBill/listByPage":
		s.billListByPage(w, r)
	case "/saas-logistics/chainLogisticsBill/getById":
		s.billGetById(w, r)
	case "/saas-logistics/chainLogisticsBill/add":
		s.billAdd(w, r)
	case "/saas-logistics/chainLogisticsBill/deleteById":
		s.billDeleteById(w, r)
	case "/saas-logistics/chainLogisticsBill/updateById":
		s.billUpdateById(w, r)
	case "/saas-logistics/chainLogisticsBill/importExcel":
		s.billImportExcel(w, r)
	case "/saas-logistics/chainLogisticsBill/exportExcel":
		s.billExportExcel(w, r)
	case "/saas-logistics/chainLogisticsBill/delBatch":
		s.billDelBatch(w, r)
	case "/saas-logistics/chainLogisticsBill/getSystemCode":
		s.billGetSystemCode(w, r)
	case "/saas-logistics/chainLogisticsBill/selectQuat":
		s.billSelectQuat(w, r)
	case "/saas-logistics/chainLogisticsBill/getContracts":
		s.billGetContracts(w, r)
	case "/saas-logistics/chainLogisticsBill/getSumPrice":
		s.billGetSumPrice(w, r)
	case "/saas-logistics/chainLogisticsBill/countBuses":
		s.billCountBuses(w, r)
	case "/saas-logistics/chainLogisticsBill/countMonBuses":
		s.billCountMonBuses(w, r)
	case "/saas-logistics/chainLogisticsBill/selOrders":
		s.billSelOrders(w, r)
	case "/saas-logistics/chainLogisticsBill/selAllOrders":
		s.billSelAllOrders(w, r)
	default:
		platform.Unknown(w)
	}
}

func (s *Service) handleLc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/saas-logistics/chainLogisticsContract/list":
		s.lcList(w, r)
	case "/saas-logistics/chainLogisticsContract/listByPage":
		s.lcListByPage(w, r)
	case "/saas-logistics/chainLogisticsContract/getById":
		s.lcGetById(w, r)
	case "/saas-logistics/chainLogisticsContract/add":
		s.lcAdd(w, r)
	case "/saas-logistics/chainLogisticsContract/deleteById":
		s.lcDeleteById(w, r)
	case "/saas-logistics/chainLogisticsContract/updateById":
		s.lcUpdateById(w, r)
	case "/saas-logistics/chainLogisticsContract/importExcel":
		s.lcImportExcel(w, r)
	case "/saas-logistics/chainLogisticsContract/exportExcel":
		s.lcExportExcel(w, r)
	case "/saas-logistics/chainLogisticsContract/delBatch":
		s.lcDelBatch(w, r)
	default:
		platform.Unknown(w)
	}
}

func (s *Service) handleTw(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/saas-logistics/chainLogisticsThirdWaybill/list":
		s.twList(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/listByPage":
		s.twListByPage(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/getById":
		s.twGetById(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/add":
		s.twAdd(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/addBatch":
		s.twAddBatch(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/deleteById":
		s.twDeleteById(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/updateById":
		s.twUpdateById(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/updateBatch":
		s.twUpdateBatch(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/importExcel":
		s.twImportExcel(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/exportExcel":
		s.twExportExcel(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/delBatch":
		s.twDelBatch(w, r)
	case "/saas-logistics/chainLogisticsThirdWaybill/clearData":
		s.twClearData(w, r)
	default:
		platform.Unknown(w)
	}
}

func (s *Service) handleWl(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/saas-logistics/chainLogisticsWeighingList/list":
		s.wlList(w, r)
	case "/saas-logistics/chainLogisticsWeighingList/listByPage":
		s.wlListByPage(w, r)
	case "/saas-logistics/chainLogisticsWeighingList/getById":
		s.wlGetById(w, r)
	case "/saas-logistics/chainLogisticsWeighingList/add":
		s.wlAdd(w, r)
	case "/saas-logistics/chainLogisticsWeighingList/deleteById":
		s.wlDeleteById(w, r)
	case "/saas-logistics/chainLogisticsWeighingList/updateById":
		s.wlUpdateById(w, r)
	case "/saas-logistics/chainLogisticsWeighingList/importExcel":
		s.wlImportExcel(w, r)
	case "/saas-logistics/chainLogisticsWeighingList/exportExcel":
		s.wlExportExcel(w, r)
	case "/saas-logistics/chainLogisticsWeighingList/delBatch":
		s.wlDelBatch(w, r)
	case "/saas-logistics/chainLogisticsWeighingList/saveBatchWeighingList":
		s.wlSaveBatch(w, r)
	default:
		platform.Unknown(w)
	}
}
