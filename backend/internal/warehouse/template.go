// template.go 模板十件套(list/listByPage/getById/add/deleteById/updateById/importExcel/
// exportExcel/delBatch/updBatch)的通用实现,行为对照 BusTag 家族 Java 模板。
// 三族 listByPage 差异见包注释;响应整形 shapeEnt 对照 MP 实体 Jackson 序列化。
package warehouse

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type ent struct {
	table    string
	keys     []string   // 实体全键(含载体 searchStr/page/ids 等),响应回填 null 用
	dbcols   []string   // 实体映射 DB 列(snake)
	q        []sqlf     // getSql 候选(@ExcelColumn 字段,声明序)
	excel    []impField // Excel 导入导出列(@ExcelColumn 非空头)
	decimals []string   // BigDecimal 列(保字面小数位)
	dates    []string   // LocalDate 列(输出 date-only)
	isoTimes []string   // 无 @JsonFormat 的 LocalDateTime(输出 ISO T)
	bools    []string   // Boolean 列(tinyint → true/false)
	types    []jf       // 全字段绑定类型(Jackson 校验):s/i/l/d/t/T/D/b/list/page/x
	syscol   bool       // 表有 system_code 列(bus_* 全部没有)
	listSQL  string     // mapper listByPage 原文(逐实体,含子查询列/别名)
	listNoSQL string    // mapper listNoPage 原文
	searchLit bool      // true → '%${maps.searchStr}%' 字面注入(缺失为字面 null)
}

type jf struct{ Name, Kind string }

// typeOf 字段绑定类型(生成清单线性查,实体字段 ≤217,无热点)
func (e *ent) typeOf(name string) string {
	for _, f := range e.types {
		if f.Name == name {
			return f.Kind
		}
	}
	return "s"
}

// javaStr 值的 Java toString 形态(数字、true/false、LocalDateTime 空格→T、BigDecimal 保字面)
func javaStr(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	default:
		b, _ := json.Marshal(v)
		s := string(b)
		if len(s) == 21 && s[10] == ' ' { // "2023-07-04 10:28:13" → LocalDateTime.toString
			return s[:10] + "T" + s[11:]
		}
		return s
	}
}

// shape 兼容入口(无列序,重复驼峰键随机——新代码用 shapeRows)
func (e *ent) shape(rows []map[string]any) []map[string]any {
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		cam := platform.Camelize(row)
		m := e.baseMap(cam)
		out = appendRow(out, m, row)
	}
	return out
}

func appendRow(out []map[string]any, m map[string]any, row map[string]any) []map[string]any {
	empty := true
	for _, v := range row {
		if v != nil {
			empty = false
			break
		}
	}
	if empty {
		return append(out, nil)
	}
	return append(out, m)
}

// shapeRows 有序整形:同列名(快照列 vs 子查询别名)按 SELECT 顺序覆盖
// (后列含 NULL 也覆盖,Java callSettersOnNulls=true + foundValues 行为实证),
// 全空行 → null 元素(returnInstanceForEmptyRow=false 实证)
func (e *ent) shapeRows(rows []map[string]any, cols []string) []map[string]any {
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		m := make(map[string]any, len(e.keys))
		for _, c := range cols {
			m[camelCol(c)] = row[c]
		}
		for _, k := range e.keys {
			if _, ok := m[k]; !ok {
				m[k] = nil
			}
		}
		e.decorate(m)
		out = appendRow(out, m, row)
	}
	return out
}

// baseMap 无序版键回填(遗留)
func (e *ent) baseMap(cam map[string]any) map[string]any {
	m := make(map[string]any, len(e.keys))
	for _, k := range e.keys {
		if v, ok := cam[k]; ok {
			m[k] = v
		} else {
			m[k] = nil
		}
	}
	e.decorate(m)
	return m
}

// decorate 按实体字段 Java 类型整形(实体字段类型驱动:String←int 列转字符串、
// Boolean←tinyint、BigDecimal 保字面、LocalDate 去时刻、无 jf LocalDateTime 转 T)
func (e *ent) decorate(m map[string]any) {
	for k, v := range m {
		if v == nil {
			continue
		}
		switch e.typeOf(k) {
		case "s":
			switch t := v.(type) {
			case string:
			case float64:
				m[k] = strconv.FormatFloat(t, 'f', -1, 64)
			case int64:
				m[k] = strconv.FormatInt(t, 10)
			case json.Number:
			case bool:
				if t {
					m[k] = "true"
				} else {
					m[k] = "false"
				}
			}
		case "i", "l":
			if _, ok := v.(json.Number); !ok {
				m[k] = jsonNumberInt64(int64(num(v)))
			}
		case "d":
			if s, ok := v.(string); ok && s != "" {
				m[k] = json.Number(s)
			}
		case "b":
			m[k] = num(v) != 0
		case "D":
			if s, ok := v.(string); ok && len(s) == 19 {
				m[k] = s[:10]
			}
		case "T":
			if s, ok := v.(string); ok && len(s) == 19 {
				m[k] = s[:10] + "T" + s[11:]
			}
		}
	}
}

// shapeOne 单行(getById);nil 行 → nil
func (e *ent) shapeOne(row map[string]any) map[string]any {
	if row == nil {
		return nil
	}
	return e.shape([]map[string]any{row})[0]
}

// camelCol 列标签→驼峰(platform.Camelize 同款,含连续下划线逐段大写)
func camelCol(s string) string {
	out := platform.Camelize(map[string]any{s: nil})
	for k := range out {
		return k
	}
	return s
}

// qEntEnt 有序查询+实体整形(含别名列/JOIN 重复列的查询统一走这里;失败回 999)
func (s *Service) qEntEnt(w http.ResponseWriter, d interface{ Query(string, ...any) (*sql.Rows, error); QueryRow(string, ...any) *sql.Row }, e *ent, q string, args ...any) []map[string]any {
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return nil
	}
	ms, cols, err := platform.ScanMapsCols(rows)
	if err != nil {
		platform.Unknown(w)
		return nil
	}
	return e.shapeRows(ms, cols)
}

// getSqlUnstr 同 getSql 但带 unstr(跳过指定字段)
func (e *ent) getSqlUnstr(m map[string]any, unstr string) (string, []any, bool) {
	sqlStr := "where 1=1"
	var args []any
	for _, f := range e.q {
		if strings.Contains(unstr, f.Name) {
			continue
		}
		v, ok := m[f.Name]
		if !ok || v == nil {
			continue
		}
		if s, isStr := v.(string); isStr && s == "0" {
			continue
		}
		val := v
		if f.Kind == "w" {
			val = "%" + javaStr(v) + "%"
		}
		sqlStr += " and " + camelOf(e, f.Name) + " like ?"
		args = append(args, val)
	}
	return sqlStr, args, true
}

// camelOf 实体字段名 → DB 列名
func camelOf(e *ent, field string) string {
	u := platform.Underscore(map[string]any{field: nil})
	for k := range u {
		return k
	}
	return field
}

// getSql 对照 StringHelper.getSql+getMaps:@ExcelColumn 字段值非 null 且 ≠ 字符串"0"
// 才拼条件;绑定值 Integer/LocalDate 原值、其余 %v% 包裹(Long 也包裹)。
func (e *ent) getSql(m map[string]any) (string, []any, bool) {
	sqlStr := "where 1=1"
	var args []any
	for _, f := range e.q {
		v, ok := m[f.Name]
		if !ok || v == nil {
			continue
		}
		if s, isStr := v.(string); isStr && s == "0" {
			continue // "0".equals(值):仅字符串 "0" 跳过,Integer 0 不跳
		}
		val := v
		if f.Kind == "w" {
			val = "%" + javaStr(v) + "%"
		}
		sqlStr += " and " + camelOf(e, f.Name) + " like ?"
		args = append(args, val)
	}
	return sqlStr, args, true
}

// ── Jackson 绑定校验(body 值类型不符 → 999)──

// bindVal 对照 Jackson 对各 Java 类型的绑定接受面;返回 (绑定值, ok)。
// "" 对非 String 一律绑 null;失败面:数字字段非数字串、日期字段坏格式、
// Boolean 收数字/数组、标量字段收对象等。
func bindVal(kind string, v any) (any, bool) {
	switch kind {
	case "s":
		switch t := v.(type) {
		case string:
			return t, true
		case float64:
			return strconv.FormatFloat(t, 'f', -1, 64), true
		case bool:
			if t {
				return "true", true
			}
			return "false", true
		}
		return nil, false
	case "i", "l":
		switch t := v.(type) {
		case string:
			if t == "" {
				return nil, true
			}
			if _, err := strconv.ParseInt(t, 10, 64); err != nil {
				return nil, false
			}
			return t, true
		case float64:
			if t != float64(int64(t)) {
				return nil, false // Jackson 整数字段拒小数(8765 实测布尔/小数都绑失败 999)
			}
			return strconv.FormatFloat(t, 'f', -1, 64), true
		}
		return nil, false // 布尔/对象/数组不绑定整数字段
	case "d":
		switch t := v.(type) {
		case string:
			if t == "" {
				return nil, true
			}
			if _, err := strconv.ParseFloat(t, 64); err != nil {
				return nil, false
			}
			return t, true
		case float64:
			return strconv.FormatFloat(t, 'f', -1, 64), true
		}
		return nil, false
	case "D":
		s, ok := v.(string)
		if !ok {
			return nil, false
		}
		if s == "" {
			return nil, true
		}
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return nil, false
		}
		return s, true
	case "t": // @JsonFormat 空格格式:只收 "yyyy-MM-dd HH:mm:ss"
		s, ok := v.(string)
		if !ok {
			return nil, false
		}
		if s == "" {
			return nil, true
		}
		if _, err := time.Parse("2006-01-02 15:04:05", s); err != nil {
			return nil, false
		}
		return s, true
	case "T": // 无注解 LocalDateTime:只收 ISO "yyyy-MM-ddTHH:mm:ss"
		s, ok := v.(string)
		if !ok {
			return nil, false
		}
		if s == "" {
			return nil, true
		}
		if _, err := time.Parse("2006-01-02T15:04:05", s); err != nil {
			return nil, false
		}
		return s, true
	case "b":
		switch t := v.(type) {
		case bool:
			return t, true
		case string:
			if t == "" {
				return nil, true
			}
			if t == "true" {
				return true, true
			}
			if t == "false" {
				return false, true
			}
		case float64:
			return t != 0, true // Jackson 数字→Boolean:0=false 非零=true(8765 实测 isRealCheck:1)
		}
		return nil, false
	}
	// list/page/x 载体:getById 等标量路径不绑;调用方单独处理
	return v, true
}

// validateBody 对照 Jackson 对实体的绑定:body 键按实体类型校验,失败 → 999。
// 实体外键 Jackson 忽略;page 键必须对象;ids 必须整数数组
func (e *ent) validateBody(m map[string]any) bool {
	for k, v := range m {
		if v == nil {
			continue
		}
		switch k {
		case "page":
			if _, ok := v.(map[string]any); !ok {
				return false
			}
			continue
		case "ids":
			arr, ok := v.([]any)
			if !ok {
				return false
			}
			for _, it := range arr {
				if _, ok2 := bindVal("i", it); !ok2 {
					return false
				}
			}
			continue
		}
		kind := e.typeOf(k)
		if kind == "list" || kind == "page" || kind == "x" {
			continue // 载体对象键:标量校验路径不涉及(listByPage body 传对象即忽略)
		}
		if _, ok := bindVal(kind, v); !ok {
			return false
		}
	}
	return true
}

// ── 模板端点实现 ──

// tplList 对照 /list:QueryWrapper.allEq(maps,false) 全部等值,返回实体列表
func (s *Service) tplList(w http.ResponseWriter, r *http.Request, e *ent) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, e.dbcols)
	if !ok {
		platform.Unknown(w) // 驼峰键等非列键 → SQL 报错
		return
	}
	q := "SELECT * FROM " + e.table
	var args []any
	if len(cond) > 0 {
		ks := platform.SortedKeys(cond)
		var conds []string
		for _, k := range ks {
			conds = append(conds, k+" = ?")
			args = append(args, cond[k])
		}
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	rows := s.qEntEnt2(w, r, e, q, args...)
	if rows == nil {
		return
	}
	platform.Success(w, rows)
}

// ipageOf MP Page JSON 形状(size=0 时 pages=0,防除零)
func ipageOf(records []map[string]any, total, current, size int64) map[string]any {
	if records == nil {
		records = []map[string]any{}
	}
	pages := int64(0)
	if size > 0 {
		pages = (total + size - 1) / size
	}
	return map[string]any{
		"records": records, "total": total, "size": size, "current": current,
		"orders": []any{}, "optimizeCountSql": true, "hitCount": false,
		"countId": nil, "maxLimit": nil, "searchCount": true, "pages": pages,
	}
}

// tplListByPage 对照三族 listByPage:
//   - 'A' 族(Config 等):controller 包装 service 返回的 IPage;page 缺失 → mapper 收 null
//     Page → 999(8765 实测);size<0 → MP 不分页;current=0 → OFFSET 负数 → 999
//   - 'B' 族(Tag 等):page 缺失 → listNoPage 全量 List
//   - 'C' 族(AdvOperaLog 等):B 族 + startTime/endTime DATE_FORMAT 字面注入
func (s *Service) tplListByPage(w http.ResponseWriter, r *http.Request, e *ent, family byte) {
	m := bodyMap(r)
	if !e.validateBody(m) {
		platform.Unknown(w) // Jackson 绑定失败
		return
	}
	sqlStr, args, ok := e.getSql(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	if family == 'C' {
		for _, kv := range []struct {
			k, op string
		}{{"startTime", ">= '"}, {"endTime", "<= '"}} {
			if v, has := m[kv.k]; has && v != nil {
				d, ok2 := bindVal("D", v)
				if !ok2 {
					platform.Unknown(w)
					return
				}
				if d == nil {
					continue
				}
				sqlStr += " and DATE_FORMAT(create_time, '%Y-%m-%d')" + kv.op + d.(string) + "'"
			}
		}
	}
	var search any
	if v, has := m["searchStr"]; has && v != nil {
		ss, ok2 := bindVal("s", v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		search = "%" + ss.(string) + "%"
	}
	// 逐实体 mapper SQL(生成清单):${searchSql} 替换、#{maps.searchStr} 参数化、
	// ${maps.searchStr} 字面注入(getMaps 已包 %v%,模板再包一层 → '%v%' 语义;缺失为字面 null)
	q := strings.Replace(e.listSQL, "${searchSql}", sqlStr, 1)
	if e.searchLit {
		inner := "%null%"
		if s, ok2 := search.(string); ok2 {
			inner = s
		}
		q = strings.Replace(q, "${maps.searchStr}", inner, 1)
	}
	baseArgs := args
	if !e.searchLit {
		q = strings.Replace(q, "#{maps.searchStr}", "?", 1)
		baseArgs = append(args, search)
	}
	countQ := "SELECT COUNT(*) FROM (" + stripOrderBy(q) + ") cnt"
	where := q

	current, size := int64(1), int64(10)
	pageObj, hasPage := m["page"].(map[string]any)
	if hasPage {
		if v, ok := pageObj["current"]; ok && v != nil {
			current = int64(num(v))
		}
		if v, ok := pageObj["size"]; ok && v != nil {
			size = int64(num(v))
		}
	}
	d := s.pick(r)
	if family == 'A' && !hasPage {
		platform.Unknown(w) // 实测:family A 无 page → mapper 收 null Page → 999
		return
	}
	if (family == 'B' || family == 'C') && !hasPage {
		rows := s.qEntEnt2(w, r, e, where, baseArgs...)
		if rows == nil {
			return
		}
		platform.Success(w, rows)
		return
	}
	if size < 0 { // MP:size<0 不分页
		rows := s.qEntEnt2(w, r, e, where, baseArgs...)
		if rows == nil {
			return
		}
		var total int64
		_ = d.QueryRow(countQ, baseArgs...).Scan(&total)
		platform.Success(w, ipageOf(rows, total, current, size))
		return
	}
	var total int64
	if err := d.QueryRow(countQ, baseArgs...).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	off := (current - 1) * size
	rows := s.qEntEnt2(w, r, e, where+" LIMIT ? OFFSET ?", append(baseArgs, size, off)...)
	if rows == nil {
		return
	}
	platform.Success(w, ipageOf(rows, total, current, size))
}

// queryEnt 有序扫描(原始行+列序)
func (s *Service) queryEnt(r *http.Request, q string, args ...any) ([]map[string]any, []string, error) {
	rows, err := s.pick(r).Query(q, args...)
	if err != nil {
		return nil, nil, err
	}
	return platform.ScanMapsCols(rows)
}

// qEntEnt2 查询+实体整形一体化(失败写 999 并返 nil)
func (s *Service) qEntEnt2(w http.ResponseWriter, r *http.Request, e *ent, q string, args ...any) []map[string]any {
	rows, cols, err := s.queryEnt(r, q, args...)
	if err != nil {
		platform.Unknown(w)
		return nil
	}
	return e.shapeRows(rows, cols)
}

// stripOrderBy 去 ORDER BY(MP count 优化近似:模板 SQL 的 order by 恒在尾部)
func stripOrderBy(q string) string {
	i := strings.LastIndex(strings.ToLower(q), " order by ")
	if i < 0 {
		return q
	}
	return q[:i]
}

// tplGetById 对照 /getById:MP getById,body 实体只取 id
func (s *Service) tplGetById(w http.ResponseWriter, r *http.Request, e *ent) {
	m := bodyMap(r)
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	row, err := platform.GetByID(s.pick(r), e.table, m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, e.shapeOne(row))
}

// tplAdd 对照 /add:MP save(IdType.AUTO 忽略 body id)+ fill 时间戳;模板响应 success() 无 data
func (s *Service) tplAdd(w http.ResponseWriter, r *http.Request, e *ent) {
	m := bodyMap(r)
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	vals := pick(m, e.dbcols)
	delete(vals, "id")
	if _, err := platform.Insert(s.pick(r), e.table, vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// tplUpdateById 对照 /updateById:MP updateById(非空子集 + fill update_time);
// warehouse 实体全带 fill → SET 恒有 update_time,永不空 SET
func (s *Service) tplUpdateById(w http.ResponseWriter, r *http.Request, e *ent) {
	m := bodyMap(r)
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	vals := pick(m, e.dbcols)
	delete(vals, "id")
	if _, err := platform.UpdateByID(s.pick(r), e.table, m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// tplDeleteById 对照 /deleteById:"删除成功，id为：X"(id 缺失 → 字面 "null")
func (s *Service) tplDeleteById(w http.ResponseWriter, r *http.Request, e *ent) {
	m := bodyMap(r)
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	if err := platform.DeleteByID(s.pick(r), e.table, m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	idStr := "null"
	if m["id"] != nil {
		idStr = str(m, "id")
	}
	platform.Success(w, "删除成功，id为："+idStr)
}

// tplDelBatch 对照 /delBatch:removeByIds(null/空 → MP 短路);模板响应 data "删除成功"
func (s *Service) tplDelBatch(w http.ResponseWriter, r *http.Request, e *ent, withData bool) {
	m := bodyMap(r)
	ids, _ := m["ids"].([]any)
	if len(ids) > 0 {
		if err := platform.DeleteByIDs(s.pick(r), e.table, ids); err != nil {
			platform.Unknown(w)
			return
		}
	}
	if withData {
		platform.Success(w, "删除成功")
		return
	}
	platform.Success(w, nil)
}

// tplUpdBatch 对照 /updBatch:MP updateBatchById(逐行非空子集+fill);空数组短路成功
func (s *Service) tplUpdBatch(w http.ResponseWriter, r *http.Request, e *ent) {
	var items []any
	_ = json.NewDecoder(r.Body).Decode(&items)
	for _, it := range items {
		em, _ := it.(map[string]any)
		if em == nil {
			platform.Unknown(w)
			return
		}
		if !e.validateBody(em) {
			platform.Unknown(w)
			return
		}
		vals := pick(em, e.dbcols)
		delete(vals, "id")
		if _, err := platform.UpdateByID(s.pick(r), e.table, em["id"], vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, nil)
}

// tplImportExcel 对照模板 importExcel:首行 selectSystemCode 探测(bus_* 无 system_code
// 列 → 恒 999,实测);探测过后 4 参 readExcel(全 sheet,编号传空)→ 错走 fail(errInfos)、
// 成功 saveBatch + "导入成功，共N条"
func (s *Service) tplImportExcel(w http.ResponseWriter, r *http.Request, e *ent) {
	if !e.syscol {
		platform.Unknown(w)
		return
	}
	left := currentUser(r).EnterpriseCode + time.Now().Format("20060102") // idt 恒 ""
	if !sysProbe(s.pick(r), e.table, left) {
		platform.Unknown(w)
		return
	}
	s.importTemplate(w, r, e.table, e.excel)
}

// tplExportExcel 对照 /exportExcel:list() 全表 → writeExcel(@ExcelColumn 声明序)
func (s *Service) tplExportExcel(w http.ResponseWriter, r *http.Request, e *ent) {
	s.exportHandler(w, r, e.table, e.excel)
}

// nowSpace fill 时间(空格秒级,对照 @JsonFormat 序列化)
func nowSpace() string { return time.Now().Format("2006-01-02 15:04:05") }

// decorateRow 无列序版单行整形(键已驼峰;按实体字段类型驱动)
func (e *ent) decorateRow(row map[string]any) map[string]any {
	m := make(map[string]any, len(e.keys))
	for k, v := range row {
		m[camelCol(k)] = v
	}
	for _, k := range e.keys {
		if _, ok := m[k]; !ok {
			m[k] = nil
		}
	}
	e.decorate(m)
	return m
}
