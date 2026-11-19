package platform

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type TencentError struct{ Code, Message string }

func (e *TencentError) Error() string { return e.Message }

// TencentJSON 复用腾讯云 API 3.0 签名协议；商户凭据仅由私有运行配置提供。
func TencentJSON(ctx context.Context, service, action, version, region string, body map[string]any) (map[string]any, error) {
	id, key := os.Getenv("TENCENT_SECRET_ID"), os.Getenv("TENCENT_SECRET_KEY")
	if id == "" {
		id = os.Getenv("secretId")
	}
	if key == "" {
		key = os.Getenv("secretKey")
	}
	if id == "" || key == "" {
		return nil, &TencentError{Code: "AuthFailure.SecretIdNotFound", Message: "腾讯云凭据未配置"}
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	host := service + ".tencentcloudapi.com"
	endpoint := os.Getenv("TENCENT_" + strings.ToUpper(service) + "_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://" + host
	}
	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	hash := func(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }
	mac := func(key []byte, text string) []byte {
		h := hmac.New(sha256.New, key)
		h.Write([]byte(text))
		return h.Sum(nil)
	}
	canonical := "POST\n/\n\ncontent-type:application/json; charset=utf-8\nhost:" + host + "\n\ncontent-type;host\n" + hash(data)
	scope := date + "/" + service + "/tc3_request"
	text := "TC3-HMAC-SHA256\n" + strconv.FormatInt(timestamp, 10) + "\n" + scope + "\n" + hash([]byte(canonical))
	signingKey := mac(mac(mac([]byte("TC3"+key), date), service), "tc3_request")
	signature := "TC3-HMAC-SHA256 Credential=" + id + "/" + scope + ", SignedHeaders=content-type;host, Signature=" + hex.EncodeToString(mac(signingKey, text))
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Host = host
	for k, v := range map[string]string{"Content-Type": "application/json; charset=utf-8", "Authorization": signature, "X-TC-Action": action, "X-TC-Version": version, "X-TC-Timestamp": strconv.FormatInt(timestamp, 10), "X-TC-Region": region} {
		req.Header.Set(k, v)
	}
	if token := os.Getenv("TENCENT_TOKEN"); token != "" {
		req.Header.Set("X-TC-Token", token)
	}
	response, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var envelope struct{ Response map[string]any }
	dec := json.NewDecoder(response.Body)
	dec.UseNumber()
	if err = dec.Decode(&envelope); err != nil {
		return nil, err
	}
	if envelope.Response == nil {
		return nil, fmt.Errorf("腾讯云响应缺少 Response")
	}
	if failure, ok := envelope.Response["Error"].(map[string]any); ok {
		return nil, &TencentError{Code: fmt.Sprint(failure["Code"]), Message: fmt.Sprint(failure["Message"])}
	}
	return envelope.Response, nil
}

// SendLegacySMS 对应 SendMessageUtil：腾讯云异常被 Java 吞掉，不阻断原业务链。
func SendLegacySMS(ctx context.Context, phone string, params []string, template, sign string) {
	SendLegacySMSMany(ctx, []string{phone}, params, template, sign)
}

func SendLegacySMSMany(ctx context.Context, phones []string, params []string, template, sign string) {
	if template == "" {
		template = "876417"
	}
	if sign == "" {
		sign = "熊猫煤仓"
	}
	appid := os.Getenv("TENCENT_SMS_APP_ID")
	if appid == "" {
		appid = "1400489027"
	}
	numbers := make([]string, 0, len(phones))
	for _, phone := range phones {
		numbers = append(numbers, "86"+phone)
	}
	_, _ = TencentJSON(ctx, "sms", "SendSms", "2019-07-11", "", map[string]any{"SmsSdkAppid": appid, "Sign": sign, "TemplateID": template, "PhoneNumberSet": numbers, "TemplateParamSet": params})
}
