package platform

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// 通联协议来自服务器运行包的 allinpay-sdk：UTF-8 字符串 SM3withSM2 签名插入 head，
// 表单字段 msgPlain，不是加密 JSON。复用平台 OpenSSL，不增加 Go 加密依赖。
func tlOpenSSL(ctx context.Context, input []byte, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "openssl", args...)
	cmd.Stdin = bytes.NewReader(input)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("通联证书或 SM2 操作失败: %w", err)
	}
	return output, nil
}

func allinpayXML(code, inst string, request map[string]any) string {
	now := time.Now()
	head := map[string]any{"processing_code": code, "inst_id": inst, "trans_date": now.Format("20060102"), "trans_time": now.Format("150405"), "ver_num": "1.00"}
	if date := request["_trans_date"]; date != nil {
		head["trans_date"] = date
		delete(request, "_trans_date")
	}
	return `<?xml version="1.0" encoding="GBK"?>` + tlXMLMap("transaction", map[string]any{"head": head, "request": request})
}

func tlXMLMap(name string, values map[string]any) string {
	var b strings.Builder
	b.WriteString("<" + name + ">")
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := values[k]
		if child, ok := v.(map[string]any); ok {
			b.WriteString(tlXMLMap(k, child))
			continue
		}
		b.WriteString("<" + k + ">")
		if v != nil {
			text := JnullAny(v)
			if _, ok := v.([]any); ok {
				data, _ := json.Marshal(v)
				text = string(data)
			}
			_ = xml.EscapeText(&b, []byte(text))
		}
		b.WriteString("</" + k + ">")
	}
	b.WriteString("</" + name + ">")
	return b.String()
}

func AllinpayPost(ctx context.Context, endpoint, inst, code string, request map[string]any, verify bool) (map[string]any, string, string, error) {
	if os.Getenv("tlPfxPath") == "" || os.Getenv("tlCertPath") == "" {
		return nil, "", "", fmt.Errorf("通联证书未配置")
	}
	dir, err := os.MkdirTemp("", "panda-tl-")
	if err != nil {
		return nil, "", "", err
	}
	defer os.RemoveAll(dir)
	key, cert := filepath.Join(dir, "merchant.pem"), filepath.Join(dir, "merchant-cert.pem")
	// 密码只通过环境传给子进程，私钥仅写入权限为 0700 的临时目录。
	for _, args := range [][]string{
		{"pkcs12", "-in", os.Getenv("tlPfxPath"), "-passin", "env:TL_SM2_PFX_PASSWORD", "-nocerts", "-nodes", "-out", key},
		{"pkcs12", "-in", os.Getenv("tlPfxPath"), "-passin", "env:TL_SM2_PFX_PASSWORD", "-clcerts", "-nokeys", "-out", cert},
	} {
		if _, err = tlOpenSSL(ctx, nil, args...); err != nil {
			// 旧 Java 商户 pfx 可能使用 RC2，OpenSSL 3 需启用 legacy provider。
			if _, err = tlOpenSSL(ctx, nil, append(args, "-legacy")...); err != nil {
				return nil, "", "", err
			}
		}
	}
	serial, err := tlOpenSSL(ctx, nil, "x509", "-in", cert, "-noout", "-serial")
	if err != nil {
		return nil, "", "", err
	}
	sn, ok := new(big.Int).SetString(strings.TrimSpace(strings.TrimPrefix(string(serial), "serial=")), 16)
	if !ok {
		return nil, "", "", fmt.Errorf("通联证书序号无效")
	}
	plain := allinpayXML(code, inst, request)
	sig, err := tlOpenSSL(ctx, []byte(plain), "dgst", "-sm3", "-sign", key, "-sigopt", "distid:1234567812345678")
	if err != nil {
		return nil, "", "", err
	}
	signed := strings.Replace(plain, "</head>", "<sign_code>"+base64.StdEncoding.EncodeToString(sig)+"</sign_code></head>", 1)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(url.Values{"msgPlain": {signed}}.Encode()))
	if err != nil {
		if code == "2089" && !verify {
			return nil, signed, err.Error(), nil
		}
		return nil, "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("allinpay-signType", "SM2")
	req.Header.Set("allinpay-certSn", sn.String())
	res, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		if code == "2089" && !verify {
			return nil, signed, err.Error(), nil
		}
		return nil, "", "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	raw := string(body)
	if err != nil {
		return nil, signed, raw, err
	}
	if strings.Contains(strings.ToLower(res.Header.Get("Content-Type")), "charset=gb") {
		body, err = execInput(ctx, body, "iconv", "-f", "GB18030", "-t", "UTF-8")
		raw = string(body)
		if err != nil {
			return nil, signed, raw, err
		}
	}
	if res.StatusCode != http.StatusOK {
		// 原 SDK readHttpResponse 将非 200 响应包装为文本，2089 仍保存订单。
		var headers strings.Builder
		names := make([]string, 0, len(res.Header))
		for name := range res.Header {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			for _, value := range res.Header[name] {
				headers.WriteString("\t" + name + ": " + value)
			}
		}
		raw = fmt.Sprintf("status:%s %s@@headers:%s@@response length:%d@@response content:%s", res.Proto, res.Status, headers.String(), len([]rune(raw)), strings.ReplaceAll(raw, "\r\n", ""))
	}
	if code == "2089" && !verify {
		return nil, signed, raw, nil
	}
	if verify {
		if err = VerifyAllinpay(ctx, raw); err != nil {
			return nil, signed, raw, err
		}
	}
	response, err := AllinpayResponse(raw)
	return response, signed, raw, err
}

func execInput(ctx context.Context, input []byte, command string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdin = bytes.NewReader(input)
	return cmd.Output()
}

func VerifyAllinpay(ctx context.Context, raw string) error {
	start, end := strings.Index(raw, "<sign_code>"), strings.Index(raw, "</sign_code>")
	if start < 0 || end < start {
		return fmt.Errorf("通联响应签名缺失")
	}
	sig, err := base64.StdEncoding.DecodeString(raw[start+len("<sign_code>") : end])
	if err != nil {
		return err
	}
	plain := raw[:start] + raw[end+len("</sign_code>"):]
	dir, err := os.MkdirTemp("", "panda-tl-verify-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	pub, err := tlOpenSSL(ctx, nil, "x509", "-in", os.Getenv("tlCertPath"), "-pubkey", "-noout")
	if err != nil {
		return err
	}
	pubPath, sigPath := filepath.Join(dir, "public.pem"), filepath.Join(dir, "signature")
	if err = os.WriteFile(pubPath, pub, 0600); err != nil {
		return err
	}
	if err = os.WriteFile(sigPath, sig, 0600); err != nil {
		return err
	}
	_, err = tlOpenSSL(ctx, []byte(plain), "dgst", "-sm3", "-verify", pubPath, "-signature", sigPath, "-sigopt", "distid:1234567812345678")
	return err
}

type tlXMLNode struct {
	XMLName  xml.Name
	Text     string      `xml:",chardata"`
	Children []tlXMLNode `xml:",any"`
}

// 查询接口的 Java XML.toJSONObject 与 SDK DataRow 都保留层次、重复元素及空值。
func AllinpayResponse(raw string) (map[string]any, error) {
	var root tlXMLNode
	dec := xml.NewDecoder(strings.NewReader(raw))
	dec.CharsetReader = func(_ string, in io.Reader) (io.Reader, error) { return in, nil } // HTTP 层已将正文转成 UTF-8。
	if err := dec.Decode(&root); err != nil {
		return nil, err
	}
	if root.XMLName.Local != "transaction" {
		return nil, fmt.Errorf("通联响应根节点错误")
	}
	for _, node := range root.Children {
		if node.XMLName.Local == "response" {
			out, _ := tlXMLValue(node).(map[string]any)
			return out, nil
		}
	}
	return nil, fmt.Errorf("通联响应缺少 response")
}

func tlXMLValue(node tlXMLNode) any {
	if len(node.Children) == 0 {
		return strings.TrimSpace(node.Text)
	}
	out := map[string]any{}
	for _, child := range node.Children {
		key := child.XMLName.Local
		value := tlXMLValue(child)
		if old, ok := out[key]; ok {
			if list, ok := old.([]any); ok {
				out[key] = append(list, value)
			} else {
				out[key] = []any{old, value}
			}
		} else {
			out[key] = value
		}
	}
	return out
}
