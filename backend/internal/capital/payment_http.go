package capital

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// bankPost 对照 PaymentHttpHelper.sendSignPost：请求体摘要与 RSA 签名均按实际发送字节计算。
// 沿用 Java busbank_* 配置，不在源码保存任何商户凭据。
func bankPost(r *http.Request, path string, body map[string]any) (map[string]any, error) {
	data, err := json.Marshal(map[string]any{"head": map[string]any{"version": "1.0", "datetime": time.Now().Format("2006-01-02 15:04:05")}, "body": body})
	if err != nil {
		return nil, err
	}
	der, err := base64.StdEncoding.DecodeString(os.Getenv("busbank_privateKey"))
	if err != nil {
		return nil, err
	}
	parsed, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}
	key, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("银行签名需要 RSA 私钥")
	}
	digest := sha256.Sum256(data)
	sign := hex.EncodeToString(digest[:])
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	source := "appid=" + os.Getenv("busbank_appId") + "&secret=" + os.Getenv("busbank_appsecret") + "&sign=" + sign + "&timestamp=" + timestamp
	hash := sha256.Sum256([]byte(source))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(r.Context(), "POST", os.Getenv("busbank_host")+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	for k, v := range map[string]string{"Content-Type": "application/json; charset=utf-8", "appid": os.Getenv("busbank_appId"), "timestamp": timestamp, "sign": sign, "apisign": base64.StdEncoding.EncodeToString(signature), "verify": "SHA256withRSAVerify"} {
		req.Header.Set(k, v)
	}
	response, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var result map[string]any
	dec := json.NewDecoder(response.Body)
	dec.UseNumber()
	err = dec.Decode(&result)
	return result, err
}

func bankMerchant(mch any) map[string]any {
	dto := map[string]any{"mchNbr": mch, "pltNbr": mch}
	if value, exists := os.LookupEnv("pltNbr"); exists {
		dto["pltNbr"] = value
	}
	if app, exists := os.LookupEnv("busbank_appId"); exists {
		dto["appId"] = app
	}
	return dto
}

func bankCopy(dto, body map[string]any, keys ...string) {
	for _, k := range keys {
		if body[k] != nil {
			dto[k] = body[k]
		}
	}
}
