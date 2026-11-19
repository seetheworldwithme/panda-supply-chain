package platform

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type bodyContract struct {
	kind     string
	required bool
}

// ValidateJavaBodies 按 JAR 中的 @RequestBody 校验，避免畸形 JSON 被当成空查询。
// 与 Spring MVC 一样先鉴权，再绑定请求体；默认只读取第一个 JSON 值。
func ValidateJavaBodies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contract, ok := javaBodyContracts[r.Method+" "+r.URL.Path]
		if !ok {
			contract, ok = javaBodyContracts["ANY "+r.URL.Path]
		}
		if !ok || contract.kind == "raw" {
			next.ServeHTTP(w, r)
			return
		}
		var raw []byte
		var err error
		if r.Body != nil {
			raw, err = io.ReadAll(r.Body)
		}
		if err != nil {
			Unknown(w)
			return
		}
		trimmed := bytes.TrimSpace(raw)
		if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
			if contract.required {
				Unknown(w)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(raw))
			next.ServeHTTP(w, r)
			return
		}
		var first json.RawMessage
		if json.NewDecoder(bytes.NewReader(raw)).Decode(&first) != nil || (contract.required && bytes.Equal(first, []byte("null"))) || (contract.kind == "object" && first[0] != '{') || (contract.kind == "array" && first[0] != '[') {
			Unknown(w)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(raw))
		next.ServeHTTP(w, r)
	})
}
