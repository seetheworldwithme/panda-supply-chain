// platform - javamap.go Java 集合迭代序复刻(Excel notNullMap 错误输出序)。
package platform

import "strings"

// JavaFieldName MP 驼峰映射:列名→实体字段名(notNullMap 的键)。
func JavaFieldName(col string) string {
	parts := strings.Split(col, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// JavaHashedMapOrder 复刻 commons-collections4 4.4 HashedMap 的 keySet 迭代序
// (JAR 5 参 readExcel 的 notNullMap,经 jshell+反射实测校准):
// hash 混合 → 桶内头插 → size>=threshold 时容量翻倍(重哈希按旧桶降序、链头→尾、头插)
// → 迭代按桶降序、链头→尾。all=完整插入序(声明序),keep=移除后幸存键。
func JavaHashedMapOrder(all []string, keep map[string]bool) []string {
	hash := func(s string) int32 {
		var h int32
		for _, c := range s {
			h = 31*h + int32(c)
		}
		h += ^(h << 9)
		h ^= int32(uint32(h) >> 14)
		h += h << 4
		h ^= int32(uint32(h) >> 10)
		return h
	}
	type entry struct {
		key string
		h   int32
	}
	capacity, threshold, size := 16, 12, 0
	buckets := make([][]entry, capacity)
	for _, k := range all {
		h := hash(k)
		idx := int(h) & (capacity - 1)
		buckets[idx] = append([]entry{{k, h}}, buckets[idx]...)
		size++
		if size >= threshold {
			nc := capacity * 2
			nb := make([][]entry, nc)
			for i := capacity - 1; i >= 0; i-- {
				for _, e := range buckets[i] {
					ni := int(e.h) & (nc - 1)
					nb[ni] = append([]entry{e}, nb[ni]...)
				}
			}
			buckets, capacity, threshold = nb, nc, nc*3/4
		}
	}
	var out []string
	for i := capacity - 1; i >= 0; i-- {
		for _, e := range buckets[i] {
			if keep[e.key] {
				out = append(out, e.key)
			}
		}
	}
	return out
}