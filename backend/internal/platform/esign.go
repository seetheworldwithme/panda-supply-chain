package platform

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// EsignJSON 对照服务器 SignHttpHelper：GET 的 Content-MD5 是字面 {}，保留原错误。
// Java 捕获网络/解析错误后返回空 JSONObject；具体业务按原先的 null 解引用或 code 判断处理。
func EsignJSON(ctx context.Context, method, path string, body any) map[string]any {
	data := []byte(nil)
	sum := "{}"
	if method != "GET" {
		var err error
		data, err = json.Marshal(body)
		if err != nil {
			return map[string]any{}
		}
		hash := md5.Sum(data)
		sum = base64.StdEncoding.EncodeToString(hash[:])
	}
	return esignRequest(ctx, method, os.Getenv("host")+path, path, data, sum, "application/json; charset=UTF-8")
}

func VerifyEsignCallback(r *http.Request, raw []byte) bool {
	if os.Getenv("appSecret") == "" {
		return false
	}
	// 签名只防篡改不防重放:时间戳偏离本机超过 10 分钟的回调直接拒绝(服务器时钟需保持同步)。
	const window = 10 * 60 * 1000
	now := time.Now().UnixMilli()
	if ts, e := strconv.ParseInt(r.Header.Get("X-Tsign-Open-TIMESTAMP"), 10, 64); e != nil || ts > now+window || ts < now-window {
		return false
	}
	query := r.URL.Query()
	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var data strings.Builder
	data.WriteString(r.Header.Get("X-Tsign-Open-TIMESTAMP"))
	for _, key := range keys {
		data.WriteString(query.Get(key))
	}
	data.Write(raw)
	mac := hmac.New(sha256.New, []byte(os.Getenv("appSecret")))
	mac.Write([]byte(data.String()))
	return hmac.Equal([]byte(hex.EncodeToString(mac.Sum(nil))), []byte(r.Header.Get("X-Tsign-Open-SIGNATURE")))
}

func EsignUpload(ctx context.Context, endpoint string, data []byte) map[string]any {
	sum := md5.Sum(data)
	return esignRequest(ctx, "PUT", endpoint, endpoint, data, base64.StdEncoding.EncodeToString(sum[:]), "application/octet-stream")
}

func esignRequest(ctx context.Context, method, endpoint, path string, data []byte, sum, media string) map[string]any {
	empty := map[string]any{}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(data))
	if err != nil {
		return empty
	}
	mac := hmac.New(sha256.New, []byte(os.Getenv("appKey")))
	mac.Write([]byte(method + "\n*/*\n" + sum + "\n" + media + "\n\n" + path))
	for k, v := range map[string]string{"X-Tsign-Open-App-Id": os.Getenv("appId"), "X-Tsign-Open-Auth-Mode": "Signature", "X-Tsign-Open-Ca-Timestamp": strconv.FormatInt(time.Now().UnixMilli(), 10), "Accept": "*/*", "Content-Type": media, "Content-MD5": sum, "X-Tsign-Open-Ca-Signature": base64.StdEncoding.EncodeToString(mac.Sum(nil))} {
		req.Header.Set(k, v)
	}
	res, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return empty
	}
	defer res.Body.Close()
	data, err = io.ReadAll(io.LimitReader(res.Body, 16<<20))
	if err != nil {
		return empty
	}
	var out map[string]any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if dec.Decode(&out) != nil || out == nil {
		return empty
	}
	return out
}
