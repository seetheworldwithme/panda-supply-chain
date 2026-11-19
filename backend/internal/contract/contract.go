// Package contract 三期合同模块:saas-contract 全服务 + saas-common esign(EMain 电子签章)DB 层。
// Java 源:saas-contract ChainContractController/ChainXs/ChainCg/ChainExtend/ChainSupplier,
// saas-common esign 包 ChainEqbFlowController/ChainEqbSingersController/ChainEqbFilesController。
//
// 行为保持要点(勿"修复"):
//   - add 缺 isPayOnline/isEsign → Java 直接解引用 NPE → catch → 602 "保存合同异常"(不落库)
//   - isEsign=true 且 fileStatus≠0 → 上传附件并发起签章；份数不足返回 602 →
//     res.getData().toString() NPE → catch → 602 "保存合同异常"(合同不落库);
//     updateById 同分支忽略 launchFlow 返回值 → 更新已落库,响应 success("保存成功")
//   - updateById 异常 catch → ResultTool.success("保存失败")(不是 fail)
//   - removeContracts:DB 行 is_pay_online 为 NULL → getIsPayOnline() NPE → 999
//   - chainContract 实体 updateTime 无 @JsonFormat(T 格式),createTime 有(空格)——各实体不同,shape 分开处理
//   - toPdfFile 读硬编码 D:\gsht 路径,Linux 必失败 catch 后返回纯文本 "This is to pdf"
package contract

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type Service struct {
	DB *db.DB
}

func NewService(database *db.DB) *Service { return &Service{DB: database} }

// Register 挂载三期合同路由(全部 POST,路径与 Java 网关后服务路径一致)。
// selectQuantity/countPro 主控制器路径已由 dashboard 包注册,此处只挂 Xs/Cg 别名路径。
func Register(mux *http.ServeMux, s *Service) {
	// zjjgGuard 对照 Java 网关 JwtAuthorizationTokenFilter:applytype=zjjg 时在监管库
	// chain_auth_user 按 account 查行,无行(企业/政府账号/匿名)→ 鉴权失败 403;
	// 实测 8765:仅 ptadmin(平台,监管库有行)通过,其余账号一律 403
	zjjgGuard := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("applytype") == "zjjg" {
				u := currentUser(r)
				var cnt int64
				_ = s.DB.Zjjg.QueryRow("SELECT count(1) FROM chain_auth_user WHERE account = ?", u.Account).Scan(&cnt)
				if cnt == 0 {
					platform.NotAuth(w)
					return
				}
			}
			next(w, r)
		}
	}
	ct := func(path string, h http.HandlerFunc) {
		mux.HandleFunc("POST /saas-contract"+path, zjjgGuard(h))
	}
	// ── chainContract 主控制器(Excel 端点债务未注册)──
	// selectQuantity/countPro 主路径已由 dashboard 包注册(同语义实现),此处只挂 Xs/Cg 别名
	for _, p := range []string{"/chainContract", "/chainXsContract", "/chainCgContract"} {
		ct(p+"/list", s.list)
		ct(p+"/listByPage", s.listByPage)
		ct(p+"/getById", s.getById)
		ct(p+"/deleteById", s.deleteById)
		ct(p+"/updateById", s.updateById)
		ct(p+"/delBatch", s.delBatch)
		ct(p+"/getSystemCode", s.getSystemCode)
		ct(p+"/getSumPrice", s.getSumPrice)
		ct(p+"/searchFour", s.searchFour)
	}
	for _, p := range []string{"/chainXsContract", "/chainCgContract"} {
		ct(p+"/selectQuantity", s.cttSelectQuantity)
		ct(p+"/countPro", s.cttCountPro)
	}
	// 主控制器独有
	ct("/chainContract/countBigPro", s.countBigPro)
	ct("/chainContract/countSumContract", s.countSumContract)
	ct("/chainContract/selZjjgEtc", s.selZjjgEtc)
	ct("/chainContract/getLogistics", s.getLogistics)
	ct("/chainContract/delOnce", s.delOnce)
	// Java 此入口没有 method 限制。
	mux.HandleFunc("/saas-contract/chainContract/toPdfFile", zjjgGuard(s.toPdfFile))
	// add:主/Xs(强制 type=2)/Cg(强制 type=1)
	ct("/chainContract/add", s.add)
	ct("/chainXsContract/add", func(w http.ResponseWriter, r *http.Request) { s.addWith(w, r, 2) })
	ct("/chainCgContract/add", func(w http.ResponseWriter, r *http.Request) { s.addWith(w, r, 1) })
	// Xs/Cg 独有(主控制器没有 searchFour 等,见 Java 控制器差异;
	// countBigPro/countSumContract/selZjjgEtc/getLogistics 仅主控制器有)

	// ── chainContractExtend(Excel 2 端点债务)──
	ct("/chainContractExtend/list", s.extendList)
	ct("/chainContractExtend/listByPage", s.extendListByPage)
	ct("/chainContractExtend/getById", s.extendGetById)
	ct("/chainContractExtend/add", s.extendAdd)
	ct("/chainContractExtend/deleteById", s.extendDeleteById)
	ct("/chainContractExtend/updateById", s.extendUpdateById)
	ct("/chainContractExtend/delBatch", s.extendDelBatch)

	// ── chainContractSupplier 增量(listByPage/getById/add/updateById/delBatch 二期已注册)──
	ct("/chainContractSupplier/list", s.supplierList)
	ct("/chainContractSupplier/deleteById", s.supplierDeleteById)
	ct("/chainContractSupplier/increaseMoney", s.supplierIncreaseMoney)
	ct("/chainContractSupplier/decreaseMoney", s.supplierDecreaseMoney)

	// ── EMain 电子签章 esign 包(/saas-common,Excel 4 端点债务)──
	eq := func(path string, h http.HandlerFunc) {
		mux.HandleFunc("POST /saas-common"+path, zjjgGuard(h))
	}
	eq("/chainEqbFlow/list", s.flowList)
	eq("/chainEqbFlow/listByPage", s.flowListByPage)
	eq("/chainEqbFlow/getById", s.flowGetById)
	eq("/chainEqbFlow/add", s.flowAdd)
	eq("/chainEqbFlow/deleteById", s.flowDeleteById)
	eq("/chainEqbFlow/updateById", s.flowUpdateById)
	eq("/chainEqbFlow/delBatch", s.flowDelBatch)
	// 外部 e签宝 openapi 集成端点:本地无密钥,保 Java 失败语义
	eq("/chainEqbFlow/launchFlow", s.flowLaunchFlow)
	eq("/chainEqbFlow/confirmFlow", s.eqbExternalUnavail)
	eq("/chainEqbFlow/v3/confirmFlow", s.eqbExternalUnavail)
	eq("/chainEqbFlow/selSignUrl", s.eqbExternalUnavail)
	eq("/chainEqbFlow/v3/selSignUrl", s.eqbExternalUnavail)
	eq("/chainEqbFlow/selSignPeople", s.eqbExternalUnavail)
	eq("/chainEqbFlow/uploadFile", s.eqbExternalUnavail)
	eq("/chainEqbFlow/selLanucnInfo", s.eqbExternalUnavail)
	eq("/chainEqbFlow/getSignCallBack", s.eqbExternalUnavail)
	eq("/chainEqbFlow/v3DownFile", s.eqbExternalUnavail)

	eq("/chainEqbSingers/list", s.singersList)
	eq("/chainEqbSingers/listByPage", s.singersListByPage)
	eq("/chainEqbSingers/getById", s.singersGetById)
	eq("/chainEqbSingers/add", s.singersAdd)
	eq("/chainEqbSingers/deleteById", s.singersDeleteById)
	eq("/chainEqbSingers/updateById", s.singersUpdateById)
	eq("/chainEqbSingers/delBatch", s.singersDelBatch)

	eq("/chainEqbFiles/list", s.filesList)
	eq("/chainEqbFiles/listByPage", s.filesListByPage)
	eq("/chainEqbFiles/getById", s.filesGetById)
	eq("/chainEqbFiles/add", s.filesAdd)
	eq("/chainEqbFiles/deleteById", s.filesDeleteById)
	eq("/chainEqbFiles/updateById", s.filesUpdateById)
	eq("/chainEqbFiles/delBatch", s.filesDelBatch)
	eq("/chainEqbFiles/updBatch", s.filesUpdBatch)

	// 基础维护字典(合同页面 filters 依赖)
	s.registerDict(eq)
}

// ── 共享小工具(与 authadmin/dashboard 同款,包内私有)──

func bodyMap(r *http.Request) map[string]any {
	var m map[string]any
	_ = json.NewDecoder(r.Body).Decode(&m)
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
	case []byte:
		n, _ := strconv.ParseFloat(string(t), 64)
		return n
	case string:
		n, _ := strconv.ParseFloat(t, 64)
		return n
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
	if f, ok := v.(float64); ok {
		return int64(f), true
	}
	if n, err := strconv.ParseInt(str(m, k), 10, 64); err == nil {
		return n, true
	}
	return 0, false
}

// pick 按 DB 下划线列白名单挑写入列(前端键为驼峰)
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

// pageOf MP Page 默认 1/10
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

// searchStrOf 对照 ${searchStr} 缺失/null → 字面 "null"
func searchStrOf(m map[string]any) string {
	if v, ok := m["searchStr"]; !ok || v == nil {
		return "null"
	}
	return strings.TrimSpace(str(m, "searchStr"))
}

// paged count + 窗口查询,MP IPage 形状
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

// getSqlEq 对照 StringHelper.getSql(t, unstr):fields(@ExcelColumn camel 集合)非空、
// 不在排除表、值非"0" → `and <列> like ?`;值 Integer/LocalDate 语义原样,其余 %v%。
// ints 为不包 % 的字段集(Java 泛型 Integer;LocalDate 同样原样)
func getSqlEq(m map[string]any, fields []string, unstr map[string]bool, ints map[string]bool) (string, []any) {
	var sb strings.Builder
	var args []any
	sb.WriteString("where 1=1")
	for _, f := range fields {
		if unstr[f] {
			continue
		}
		v, ok := m[f]
		if !ok || v == nil {
			continue
		}
		if s, isStr := v.(string); isStr && s == "0" {
			continue
		}
		col := snake(f)
		sb.WriteString(" and " + col + " like ?")
		if ints[f] {
			args = append(args, v)
		} else {
			args = append(args, "%"+str(m, f)+"%")
		}
	}
	return sb.String(), args
}

func snake(s string) string {
	var b strings.Builder
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(c - 'A' + 'a')
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

// autoGenericCode 对照 StringHelper.autoGenericCode:parseInt 后 %0Nd 补零
func autoGenericCode(code string, width int) string {
	n, err := strconv.Atoi(strings.TrimSpace(code))
	if err != nil {
		n = 0
	}
	s := strconv.Itoa(n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

// fillNulls 补载体键 null(实体 null 字段 Jackson 输出)
func fillNulls(m map[string]any, keys []string) {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			m[k] = nil
		}
	}
}

// currentUser 当前请求身份(fill user id / companyId / enterpriseCode)
func currentUser(r *http.Request) *auth.CurrentUser { return auth.From(r.Context()) }

// ── chainContract 主表 CRUD ──

// contractCols 实体 DB 字段列(无 rec_lnglog/del_lnglog——实体没有这两列,MP 不写)
var contractCols = []string{"id", "im_batch_number", "company_id", "order_id", "template_id",
	"system_code", "contract_code", "contract_resource", "contract_type", "signing_time",
	"contract_name", "sales_unit", "buyer_unit", "sales_unit_id", "buyer_unit_id",
	"receipt_place", "delivery_place", "collaborator", "pro_big_type", "pro_sma_type",
	"pro_remark", "contract_num", "unit_price", "tax_included", "sum_price", "sulfur",
	"volatilization", "water", "remark", "price_type", "pay_method", "pay_type", "carrier",
	"transport_mode", "contract_ids", "agreement_ids", "supplier_id", "cust_id", "month",
	"cust_type", "is_pay_online", "is_esign", "is_esign_complete", "file_status", "sys_flow_id",
	"eqb_flow_id", "is_pay", "is_payer_confirms", "is_payee_confirms", "is_platform_confirms",
	"sum_freeze_amt", "sum_thaw_amt", "sum_refund_amt", "residue_freeze_amt", "sum_trace_fees",
	"pay_rate", "create_user_id", "update_user_id", "contract_status", "min_granularity",
	"max_granularity", "min_ash_content", "max_ash_content", "min_coal_ash", "max_coal_ash",
	"min_volatiles", "max_volatiles", "min_std", "max_std", "min_internal_sulfur",
	"max_internal_sulfur", "min_water_content", "max_water_content", "min_bond_index",
	"max_bond_index", "min_calorific_value", "max_calorific_value", "min_csr", "max_csr",
	"min_y", "max_y", "min_hgi", "max_hgi", "recovery_analysis", "lithofacies",
	"calorific_value", "clean_coal_recovery_rate", "middling_coal", "gangue", "slurry",
	"fst_pro_id", "big_pro_id", "three_pro_id", "fst_name", "sec_name", "th_name",
	"is_effect", "create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"}

var contractDecimalCols = []string{"contract_num", "unit_price", "sum_price", "sulfur",
	"volatilization", "water", "sum_freeze_amt", "sum_thaw_amt", "sum_refund_amt",
	"residue_freeze_amt", "sum_trace_fees", "min_granularity", "max_granularity",
	"min_ash_content", "max_ash_content", "min_coal_ash", "max_coal_ash", "min_volatiles",
	"max_volatiles", "min_std", "max_std", "min_internal_sulfur", "max_internal_sulfur",
	"min_water_content", "max_water_content", "min_bond_index", "max_bond_index",
	"min_calorific_value", "max_calorific_value", "min_csr", "max_csr", "min_y", "max_y",
	"min_hgi", "max_hgi", "lithofacies", "calorific_value", "clean_coal_recovery_rate",
	"middling_coal", "gangue", "slurry"}

var contractBoolCols = []string{"tax_included", "is_pay_online", "is_esign",
	"is_esign_complete", "is_pay", "is_payer_confirms", "is_payee_confirms",
	"is_platform_confirms", "is_effect"}

// contractCarrier 载体键(exist=false,Jackson 恒输出)
var contractCarrier = []string{"list", "filelist", "searchStr", "ext", "page", "ids",
	"startTime", "endTime", "ctStartTime", "ctEndTime", "timeType", "eqbFlow", "singers",
	"notBindOrder", "isPtSearch"}

// contractEchoKeys 回显键 = 实体列驼峰 + 载体键(controller 恒回显请求实体)
var contractEchoKeys = append(camelKeys(contractCols), contractCarrier...)

// shapeContract 实体 JSON 形状:DECIMAL 保小数位、tinyint(1)→布尔、
// updateTime 无 @JsonFormat(T 格式)、signing_time DATE 列输出 yyyy-MM-dd;
// rec_lnglog/del_lnglog 是表列但实体无此字段(select * 带出后剔除)
func shapeContract(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		for _, c := range contractDecimalCols {
			if s, ok := r[camel(c)].(string); ok && s != "" {
				r[camel(c)] = json.Number(s)
			}
		}
		for _, c := range contractBoolCols {
			if v, ok := r[camel(c)]; ok && v != nil {
				r[camel(c)] = num(v) != 0
			}
		}
		if s, ok := r["updateTime"].(string); ok && len(s) == 19 {
			r["updateTime"] = s[:10] + "T" + s[11:]
		}
		if s, ok := r["signingTime"].(string); ok && len(s) >= 10 {
			r["signingTime"] = s[:10]
		}
		delete(r, "recLnglog")
		delete(r, "delLnglog")
		fillNulls(r, contractCarrier)
	}
	return out
}

// camel 下划线列名转驼峰(形状层键名)
func camel(s string) string {
	out := platform.Camelize(map[string]any{s: nil})
	for k := range out {
		return k
	}
	return s
}

// list = findList:QueryWrapper.allEq(maps, false)(键原样作列名,null 跳过);
// 键必须落在实体列白名单内,否则 Java SQL 异常 → 999
func (s *Service) list(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	cond, ok := allEqCond(m, contractCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(d, "chain_contract", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, shapeContract(rows))
}

// allEqCond 同 authadmin 版:null 跳过(false),白名单外 → false(999)
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

// getById = selectById:ckById + selEqbFlow + selEqbSingers;查不到 → Java NPE → 999
func (s *Service) getById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w) // Java ckById(null) 空指针 → 999
		return
	}
	rows, err := d.Query("SELECT * FROM chain_contract WHERE id = ?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, _ := platform.ScanMaps(rows)
	if len(ms) == 0 {
		platform.Unknown(w) // cta.setId(id) NPE
		return
	}
	out := shapeContract(ms)[0]
	out["eqbFlow"] = s.eqbFlowOf(d, id)
	out["singers"] = s.eqbSingersOf(d, id)
	// ckById @Results@Many → chain_contract_extend where contract_id(空 → [])
	extRows, err := d.Query("SELECT * FROM chain_contract_extend WHERE contract_id = ?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	extMs, _ := platform.ScanMaps(extRows)
	if extMs == nil {
		extMs = []map[string]any{}
	}
	out["list"] = shapeExtend(extMs)
	platform.Success(w, out)
}

// eqbFlowOf selEqbFlow:chain_eqb_flow where contract_id(无行 → Java null)
func (s *Service) eqbFlowOf(d *sql.DB, id int64) map[string]any {
	rows, err := d.Query("SELECT * FROM chain_eqb_flow WHERE contract_id = ?", id)
	if err != nil {
		return nil
	}
	ms, _ := platform.ScanMaps(rows)
	if len(ms) == 0 {
		return nil
	}
	return shapeEqbFlow([]map[string]any{ms[0]})[0]
}

// eqbSingersOf selEqbSingers:chain_eqb_singers join chain_eqb_flow(空 → [])
func (s *Service) eqbSingersOf(d *sql.DB, id int64) []map[string]any {
	rows, err := d.Query("SELECT s.* FROM chain_eqb_singers s, chain_eqb_flow f WHERE s.sys_flow_id=f.id AND f.contract_id = ?", id)
	if err != nil {
		return []map[string]any{}
	}
	ms, _ := platform.ScanMaps(rows)
	if ms == nil {
		ms = []map[string]any{}
	}
	return shapeEqbSingers(ms)
}

// add = saveContract。行为保持:缺 isPayOnline/isEsign → 602 "保存合同异常";
// isPayOnline=true 对方未注册 → 602 "...未注册,请联系对方注册系统";
// isEsign=true 且 fileStatus≠0 → 签章返回空数据时 getData() NPE → 602(不落库)。
// 落库:chain_contract(fill:时间+create/update_user_id)+ chain_logistics_contract(ids)+
// chain_common_file(filelist,经 saas-common updBthById)+ chain_contract_extend(list)+
// chain_common_log(日志,getChangeLog 文本细节为债务,仅记前缀)
func (s *Service) add(w http.ResponseWriter, r *http.Request) {
	s.addWith(w, r, 0)
}

func (s *Service) addWith(w http.ResponseWriter, r *http.Request, t int) {
	m := bodyMap(r)
	if t != 0 {
		m["contractType"] = t
	}
	d := s.pick(r)
	u := currentUser(r)

	if _, ok := m["isPayOnline"]; !ok || m["isPayOnline"] == nil {
		platform.Fail(w, 602, "保存合同异常") // getIsPayOnline() NPE
		return
	}
	if m["isEsign"] == nil { // getIsEsign() 解引用(注意:与 isPayOnline 判断方式一致,Java 均直接解引用)
		platform.Fail(w, 602, "保存合同异常")
		return
	}
	isPayOnline := m["isPayOnline"] == true
	isEsign := m["isEsign"] == true
	fileStatus := num(m["fileStatus"])

	// 在线支付:对方单位必须已注册
	if isPayOnline {
		etcName := str(m, "buyerUnit")
		if num(m["contractType"]) == 1 {
			etcName = str(m, "salesUnit")
		}
		var cnt int64
		_ = d.QueryRow("SELECT count(1) FROM chain_basic_enterprise WHERE enterprise_name = ?", etcName).Scan(&cnt)
		if cnt == 0 {
			platform.Fail(w, 602, etcName+"未注册，请联系对方注册系统")
			return
		}
	}
	m["eqbFlowId"] = ""
	if isEsign && fileStatus != 0 {
		flowID, err := s.contractSign(r, m, false)
		if err != nil {
			platform.Fail(w, 602, "保存合同异常")
			return
		}
		m["eqbFlowId"] = flowID
	}

	vals := pick(m, contractCols...)
	delete(vals, "id") // @TableId(IdType.AUTO):body id 被忽略,自增
	vals["create_user_id"] = u.UserID
	vals["update_user_id"] = u.UserID
	id, err := platform.Insert(d, "chain_contract", vals) // fill 实体:自动 create_time/update_time
	if err != nil {
		platform.Unknown(w)
		return
	}
	s.saveChildren(r, d, m, id, u)
	platform.Success(w, "保存成功")
}

// saveChildren add/updateById 共用的子表维护(ids 运单关联/filelist 文件/list 扩展)+ 日志
func (s *Service) saveChildren(r *http.Request, d *sql.DB, m map[string]any, id int64, u *auth.CurrentUser) {
	// ids → chain_logistics_contract(先删后插 = Java add 只插/update 先删后插;add 时表无旧行)
	if ids, ok := m["ids"].([]any); ok {
		if _, err := d.Exec("DELETE FROM chain_logistics_contract WHERE contract_id = ?", id); err == nil {
			for _, it := range ids {
				loId, ok := it.(float64)
				if !ok {
					continue
				}
				_, _ = d.Exec("INSERT INTO chain_logistics_contract (contract_id, logistics_id) VALUES (?, ?)", id, int64(loId))
			}
		}
	}
	// filelist → chain_common_file(对照 saas-common updBthById:按 id 更新非空字段)
	if files, ok := m["filelist"].([]any); ok {
		for _, it := range files {
			f, ok := it.(map[string]any)
			if !ok || f["id"] == nil {
				continue
			}
			if _, err := platform.UpdateByIDNoFill(d, "chain_common_file", f["id"], pick(f, "file_type", "owner_type", "owner_id", "relation_id", "relation_type", "effective")); err != nil {
				continue
			}
		}
	}
	// list → chain_contract_extend:add=saveBatch(插),update 分支见 updateById
	if exts, ok := m["list"].([]any); ok {
		for _, it := range exts {
			e, ok := it.(map[string]any)
			if !ok {
				continue
			}
			if e["id"] == nil {
				em := map[string]any(e)
				em["contractId"] = id
				// extend 键名规则同 add:dtype/dkey/dvalue 才被 Jackson 绑定
				if _, err := platform.InsertNoFill(d, "chain_contract_extend", bodyVals(em, []string{"company_id", "contract_id", "dict_id", "d_type", "d_key", "ext_type", "d_value", "remark1", "remark2", "remark3", "remark4", "remark5"}, extendRename)); err != nil {
					continue
				}
			}
		}
	}
	// 日志(chain_common_log,无 fill:create_time 库默认;operation 文本细节债务)
	content := "新增合同\n对比结果为：\n"
	_, _ = d.Exec("INSERT INTO chain_common_log (company_id, user_id, user_name, operation, operation_type) VALUES (?, ?, ?, ?, ?)",
		strconv.FormatInt(u.CompanyID, 10), strconv.FormatInt(u.UserID, 10), u.Username, content, 1)
}

// updateById = updateContractById + controller:controller 忽略 service 返回值,恒
// ResultTool.success(请求实体)——602/"保存失败"只阻断落库,响应一律 200 实体回显。
// service 校验链:`body.isPayOnline!=null && body.isPayOnline && cc.isPayOnline && cc.createUserId!=uid`
// (请求体 false/null 短路,DB isPayOnline NULL 不触发拆箱);getById 无行 NPE 也不落库
func (s *Service) updateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	u := currentUser(r)
	echo := func() {
		em := make(map[string]any, len(m))
		for k, v := range m {
			em[k] = v
		}
		// 载体对象按实体全键回显(Java 泛型实体 Jackson 序列化;body 未绑定键输出 null)
		if ls, ok := em["list"].([]any); ok {
			out := make([]any, 0, len(ls))
			for _, it := range ls {
				if e, ok := it.(map[string]any); ok {
					out = append(out, platform.EntityEcho(e, extendEchoKeys, nil))
				} else {
					out = append(out, it)
				}
			}
			em["list"] = out
		}
		if ef, ok := em["eqbFlow"].(map[string]any); ok {
			em["eqbFlow"] = platform.EntityEcho(ef, flowEchoKeys, nil)
		}
		if ss, ok := em["singers"].([]any); ok {
			out := make([]any, 0, len(ss))
			for _, it := range ss {
				if e, ok := it.(map[string]any); ok {
					out = append(out, platform.EntityEcho(e, singersEchoKeys, nil))
				} else {
					out = append(out, it)
				}
			}
			em["singers"] = out
		}
		platform.Success(w, platform.EntityEcho(em, contractEchoKeys, nil))
	}
	id, ok := int64Of(m, "id")
	if !ok {
		echo() // Java ckById(null) NPE → catch,controller 仍回显
		return
	}
	blocked := false
	if v, ok := m["isPayOnline"]; ok && v != nil && v == true {
		etcName := str(m, "buyerUnit")
		if num(m["contractType"]) == 1 {
			etcName = str(m, "salesUnit")
			m["isPayerConfirms"] = true
		} else {
			m["isPayeeConfirms"] = true
		}
		var cnt int64
		_ = d.QueryRow("SELECT count(1) FROM chain_basic_enterprise WHERE enterprise_name = ?", etcName).Scan(&cnt)
		if cnt == 0 {
			blocked = true // 602 被 controller 丢弃
		}
	}
	if !blocked {
		cc, err := platform.GetByID(d, "chain_contract", id)
		if err != nil || cc == nil {
			blocked = true // getById NPE → catch
		} else if v, ok := m["isPayOnline"]; ok && v != nil && v == true {
			if cc["is_pay_online"] == nil {
				blocked = true // cc.getIsPayOnline() 拆箱 NPE
			} else if num(cc["is_pay_online"]) != 0 {
				if cc["create_user_id"] == nil {
					blocked = true // createUserId 拆箱 NPE
				} else if int64(num(cc["create_user_id"])) != u.UserID {
					blocked = true // 602 只有创建者本人能编辑
				}
			}
		}
	}
	if blocked {
		echo()
		return
	}
	vals := pick(m, contractCols...)
	delete(vals, "id")
	vals["update_user_id"] = u.UserID
	if _, err := platform.UpdateByID(d, "chain_contract", id, vals); err != nil {
		echo()
		return
	}
	// ids:先删后插
	if ids, ok := m["ids"].([]any); ok {
		if _, err := d.Exec("DELETE FROM chain_logistics_contract WHERE contract_id = ?", id); err == nil {
			for _, it := range ids {
				loId, ok := it.(float64)
				if !ok {
					continue
				}
				_, _ = d.Exec("INSERT INTO chain_logistics_contract (contract_id, logistics_id) VALUES (?, ?)", id, int64(loId))
			}
		}
	}
	// filelist 非空:更新文件 + (list 非空才)批更扩展 —— Java 的 list 批更嵌在 filelist 分支内
	if files, ok := m["filelist"].([]any); ok {
		for _, it := range files {
			f, ok := it.(map[string]any)
			if !ok || f["id"] == nil {
				continue
			}
			_, _ = platform.UpdateByIDNoFill(d, "chain_common_file", f["id"], pick(f, "file_type", "owner_type", "owner_id", "relation_id", "relation_type", "effective"))
		}
		if exts, ok := m["list"].([]any); ok {
			for _, it := range exts {
				e, ok := it.(map[string]any)
				if !ok || e["id"] == nil {
					continue
				}
				ev := bodyVals(e, []string{"company_id", "contract_id", "dict_id", "d_type", "d_key", "ext_type", "d_value", "remark1", "remark2", "remark3", "remark4", "remark5"}, extendRename)
				if len(ev) == 0 {
					continue // SET 空 → Java updateBatchById 异常 → catch(后续日志也不写)
				}
				_, _ = platform.UpdateByIDNoFill(d, "chain_contract_extend", e["id"], ev)
			}
		}
	}
	if m["isEsign"] == true && num(m["fileStatus"]) != 0 {
		if _, err := s.contractSign(r, m, true); err != nil {
			echo()
			return
		}
	}
	{
		content := "修改合同\n对比结果为：\n"
		_, _ = d.Exec("INSERT INTO chain_common_log (company_id, user_id, user_name, operation, operation_type) VALUES (?, ?, ?, ?, ?)",
			strconv.FormatInt(u.CompanyID, 10), strconv.FormatInt(u.UserID, 10), u.Username, content, 1)
	}
	echo()
}

// deleteById = controller:removeContract 后恒 success("删除成功,id为:X")
func (s *Service) deleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w)
		return
	}
	s.removeContract(s.pick(r), id)
	platform.Success(w, "删除成功，id为："+strconv.FormatInt(id, 10))
}

// removeContract 5 表清理:主行、运单关联、扩展、文件(relation_id)、结算/发票解绑
func (s *Service) removeContract(d *sql.DB, id int64) {
	_, _ = d.Exec("DELETE FROM chain_contract WHERE id = ?", id)
	_, _ = d.Exec("DELETE FROM chain_logistics_contract WHERE contract_id = ?", id)
	_, _ = d.Exec("DELETE FROM chain_contract_extend WHERE contract_id = ?", id)
	// delByRelId:文件行删除(saas-common removeFile 同时清文件盘;盘文件清理为债务)
	_, _ = d.Exec("DELETE FROM chain_common_file WHERE relation_id = ?", id)
	// 解除结算单/发票对应关系(updSment/updBill)
	_, _ = d.Exec("UPDATE chain_capital_settlement SET contract_id=null, contract_code=null, contract_name=null WHERE contract_id = ?", id)
	_, _ = d.Exec("UPDATE chain_bill SET contract_id=null, contract_code=null, contract_name=null WHERE contract_id = ?", id)
}

// delBatch = removeContracts:Java `cc.getIsPayOnline() && cc.getCreateUserId() != uid`
// 拆箱链——is_pay_online false 短路跳过创建者校验;true 时 create_user_id NULL 拆箱 NPE 999;
// is_pay NULL 同理 NPE;listByIds 缺行不进校验循环(删除空删仍成功)
func (s *Service) delBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	u := currentUser(r)
	ids, ok := m["ids"].([]any)
	if !ok || len(ids) == 0 {
		platform.Success(w, "删除成功") // listByIds 空集合 → 循环零次直接成功
		return
	}
	var numIds []int64
	for _, it := range ids {
		f, ok := it.(float64)
		if !ok {
			platform.Unknown(w)
			return
		}
		numIds = append(numIds, int64(f))
	}
	for _, id := range numIds {
		cc, err := platform.GetByID(d, "chain_contract", id)
		if err != nil {
			platform.Unknown(w)
			return
		}
		if cc == nil {
			continue
		}
		if cc["is_pay_online"] == nil {
			platform.Unknown(w) // Boolean 拆箱 NPE
			return
		}
		if num(cc["is_pay_online"]) != 0 {
			if cc["create_user_id"] == nil {
				platform.Unknown(w) // Integer 拆箱 NPE
				return
			}
			if int64(num(cc["create_user_id"])) != u.UserID {
				platform.Fail(w, 602, "操作失败，包含非该用户的合同")
				return
			}
		}
		if cc["is_pay"] == nil {
			platform.Unknown(w) // getIsPay() 拆箱 NPE
			return
		}
		if num(cc["is_pay"]) != 0 {
			platform.Fail(w, 602, "操作失败，包含已付款的合同")
			return
		}
	}
	for _, id := range numIds {
		s.removeContract(d, id)
	}
	platform.Success(w, "删除成功")
}

// ShapeContractRows 导出 ChainContract 行形状(四期物流 bill.getContracts/selRelContracts 复用,
// 与主控制器 getContracts 完全同源,避免形状逻辑复制)
func ShapeContractRows(rows []map[string]any) []map[string]any {
	if rows == nil {
		rows = []map[string]any{}
	}
	return shapeContract(rows)
}

// delOnce = 按条件一键删:先 allEq(maps)(null→IS NULL)校验 chain_contract 有行,
// 再同条件删 7 张表。键超出各表公共列 → Java SQL 异常 → 999
func (s *Service) delOnce(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	// 7 张表的公共安全键(实际前端只用 im_batch_number;其余键走白名单校验)
	common := map[string]bool{"im_batch_number": true, "month": true, "create_time": true,
		"update_time": true, "company_id": true, "remark1": true, "remark2": true,
		"remark3": true, "remark4": true, "remark5": true}
	var conds []string
	var args []any
	for k, v := range m {
		if !common[k] {
			platform.Unknown(w) // 不在公共列 → 某表 SQL 异常
			return
		}
		if v == nil {
			conds = append(conds, k+" IS NULL")
			continue
		}
		conds = append(conds, k+" = ?")
		args = append(args, v)
	}
	where := ""
	if len(conds) > 0 {
		where = " WHERE " + strings.Join(conds, " AND ")
	}
	var cnt int64
	if err := d.QueryRow("SELECT count(1) FROM chain_contract"+where, args...).Scan(&cnt); err != nil || cnt == 0 {
		platform.Fail(w, 602, "批次号不存在")
		return
	}
	for _, t := range []string{"chain_contract", "chain_capital_settlement", "chain_capital_detailed",
		"chain_bill", "chain_bill_rate", "chain_logistics", "chain_logistics_weighing_list"} {
		if _, err := d.Exec("DELETE FROM "+t+where, args...); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, nil)
}
