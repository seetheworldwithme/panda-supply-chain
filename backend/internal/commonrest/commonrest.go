package commonrest

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type Service struct {
	DB         *db.DB
	Auth       *auth.Service
	HTTPClient *http.Client // 微信接口测试可注入本地传输，不使用真实商户配置。
	BatchMux   http.Handler // in-process replacements for the original Feign batch calls
	Jobs       *platform.Jobs
	tlTasksMu  sync.Mutex
	tlTasks    map[string]context.CancelFunc
	tlSyncMu   sync.Mutex
	tlCache    map[string]time.Time
}

func NewService(database *db.DB) *Service { return &Service{DB: database} }

// ── 共享小工具(与 authadmin 同款,warehouse 包独立一份防跨包私用)──

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
	case float64:
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

func parseF(s string) float64 {
	n, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return n
}

// num 把 driver/JSON 数值统一转 float64(DECIMAL 的 []byte、int64、float64、json.Number)
func num(v any) float64 {
	switch t := v.(type) {
	case nil:
		return 0
	case []byte:
		return parseF(string(t))
	case string:
		return parseF(t)
	case json.Number:
		return parseF(string(t))
	case int64:
		return float64(t)
	case float64:
		return t
	default:
		return 0
	}
}

func queryRows(d *sql.DB, q string, args ...any) ([]map[string]any, error) {
	rows, err := d.Query(q, args...)
	if err != nil {
		return nil, err
	}
	ms, err := platform.ScanMaps(rows)
	if ms == nil {
		ms = []map[string]any{}
	}
	return ms, err
}

func queryCam(d *sql.DB, q string, args ...any) ([]map[string]any, error) {
	ms, err := queryRows(d, q, args...)
	if err != nil {
		return nil, err
	}
	return platform.CamelizeRows(ms), nil
}

// pageOf 从 body 取 {current,size}(MP Page 默认 1/10)
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

// paged 相同条件 count + 窗口查询,返回 MP IPage 形状
func paged(d *sql.DB, q, countQ string, args []any, current, size int64) (map[string]any, error) {
	var total int64
	if err := d.QueryRow(countQ, args...).Scan(&total); err != nil {
		return nil, err
	}
	off := (current - 1) * size
	rows, err := queryCam(d, q+" LIMIT ? OFFSET ?", append(append([]any{}, args...), size, off)...)
	if err != nil {
		return nil, err
	}
	return platform.IPage(rows, total, current, size), nil
}

// allEqCond 复刻 MP QueryWrapper.allEq(maps,false):键原样作列名,白名单外键(如驼峰)
// 会 SQL 报错 → 999;null 值键跳过(false 不生成 IS NULL)
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

// pick 从 body 按实体列白名单挑写入列(键转下划线,载体键自然被排除)
func pick(m map[string]any, cols []string) map[string]any {
	u := platform.Underscore(m)
	out := make(map[string]any, len(cols))
	for _, c := range cols {
		if v, ok := u[c]; ok {
			out[c] = v
		}
	}
	return out
}

// currentUser 当前请求身份
func currentUser(r *http.Request) *auth.CurrentUser { return auth.From(r.Context()) }

// authGuard 对照 Java 网关 JwtAuthorizationTokenFilter:无 token/解析失败 → 403。
// warehouse 全部端点在 Java 侧都要求登录(除网关白名单,本包无白名单端点)
func authGuard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := currentUser(r)
		if u.UserID == 0 && !u.SignedAnon {
			platform.NotAuth(w)
			return
		}
		next(w, r)
	}
}

func (s *Service) pick(r *http.Request) *sql.DB {
	if r.Header.Get("applytype") == "zjjg" {
		return s.DB.Zjjg
	}
	return s.DB.Chain
}

// jsonNumberInt64 整数字面(JSON 序列化为无小数点数字)
func jsonNumberInt64(n int64) json.Number { return json.Number(strconv.FormatInt(n, 10)) }

// ── Map 查询整形(SELECT * 进 Map 的 MyBatis 行)──
// 列类型驱动:tinyint(1)→bool、decimal→Double 最短表示(整数补 .0,Java Double.toString)、
// datetime→ISO T、date→日期、整数→数字、其余字符串。information_schema 懒加载缓存。
var mapTypeCache = syncMapOf{} // table → []colType

type colType struct{ name, dataType, columnType string }

type syncMapOf struct{ m map[string][]colType }

func (s *syncMapOf) load(d *sql.DB, table string) []colType {
	if s.m == nil {
		s.m = map[string][]colType{}
	}
	if v, ok := s.m[table]; ok {
		return v
	}
	rows, err := d.Query("SELECT column_name, data_type, column_type FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name=? ORDER BY ordinal_position", table)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []colType
	for rows.Next() {
		var c colType
		_ = rows.Scan(&c.name, &c.dataType, &c.columnType)
		out = append(out, c)
	}
	s.m[table] = out
	return out
}

// mapShape 按列类型整形一行(键保留蛇形标签,对照 MyBatis Map + Jackson)
func mapShape(row map[string]any, cols []colType) map[string]any {
	out := make(map[string]any, len(row))
	for _, c := range cols {
		v := row[c.name]
		switch {
		case v == nil:
			out[c.name] = nil
		case c.dataType == "tinyint" && c.columnType == "tinyint(1)": // unsigned zerofill 的 tinyint(1) Java 落 Integer
			out[c.name] = num(v) != 0
		case c.dataType == "decimal":
			// BigDecimal 保字面刻度(MyBatis Map → Jackson 原样输出)
			out[c.name] = json.Number(v.(string))
		case c.dataType == "datetime" || c.dataType == "timestamp":
			if s, ok := v.(string); ok && len(s) == 19 {
				out[c.name] = s[:10] + "T" + s[11:]
			} else {
				out[c.name] = v
			}
		case c.dataType == "date":
			if s, ok := v.(string); ok && len(s) == 19 {
				out[c.name] = s[:10]
			} else {
				out[c.name] = v
			}
		case c.dataType == "int" || c.dataType == "bigint" || c.dataType == "smallint" || c.dataType == "mediumint":
			out[c.name] = json.Number(strconv.FormatInt(int64(num(v)), 10))
		default:
			out[c.name] = v
		}
	}
	return out
}

func mapShapeRows(d *sql.DB, table string, rows []map[string]any) []map[string]any {
	cols := mapTypeCache.load(d, table)
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapShape(r, cols))
	}
	return out
}
