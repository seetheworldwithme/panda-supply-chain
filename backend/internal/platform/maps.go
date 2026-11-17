// platform - maps.go 通用行扫描与分页响应(对照 MP IPage)
package platform

import (
	"database/sql"
	"sort"
	"strings"
	"time"
)

// ScanMaps 通用行扫描:列名为键,[]byte 转字符串(nil 保持 null)。
// DATETIME(parseTime=true 时为 time.Time)统一格式化 "yyyy-MM-dd HH:mm:ss",
// 对照 Java 实体 @JsonFormat(pattern="yyyy-MM-dd HH:mm:ss")
func ScanMaps(rows *sql.Rows) ([]map[string]any, error) {
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var out []map[string]any
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		m := make(map[string]any, len(cols))
		for i, c := range cols {
			switch t := vals[i].(type) {
			case []byte:
				m[c] = string(t)
			case time.Time:
				m[c] = t.Format("2006-01-02 15:04:05")
			default:
				m[c] = vals[i]
			}
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// SortedKeys map 键排序(稳定 SQL 列序)
func SortedKeys(m map[string]any) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// IPage 对照 MyBatis-Plus Page 完整 Jackson 序列化形状(orders/optimizeCountSql 等固定字段照 Java 输出)
func IPage(records []map[string]any, total, current, size int64) map[string]any {
	if records == nil {
		records = []map[string]any{}
	}
	pages := (total + size - 1) / size
	return map[string]any{
		"records":          records,
		"total":            total,
		"size":             size,
		"current":          current,
		"orders":           []any{},
		"optimizeCountSql": true,
		"hitCount":         false,
		"countId":          nil,
		"maxLimit":         nil,
		"searchCount":      true,
		"pages":            pages,
	}
}

// Camelize 列名下划线转驼峰(role_type→roleType),对照 MP 实体 Jackson 序列化键名。
// remark1 等无下划线键原样;连续下划线逐段驼峰(pro_big_type→proBigType)
func Camelize(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[toCamel(k)] = v
	}
	return out
}

// CamelizeRows 整批驼峰化
func CamelizeRows(rows []map[string]any) []map[string]any {
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		out = append(out, Camelize(r))
	}
	return out
}

func toCamel(s string) string {
	if !strings.Contains(s, "_") {
		return s
	}
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// Underscore 驼峰转下划线(roleType→role_type),用于前端回传实体的写入
func Underscore(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		var b strings.Builder
		for i, c := range k {
			if c >= 'A' && c <= 'Z' {
				if i > 0 {
					b.WriteByte('_')
				}
				b.WriteRune(c - 'A' + 'a')
			} else {
				b.WriteRune(c)
			}
		}
		out[b.String()] = v
	}
	return out
}
