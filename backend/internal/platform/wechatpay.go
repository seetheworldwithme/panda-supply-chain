package platform

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func WechatRSASign(text string) (string, error) {
	return wechatRSASign(LoadWxConfig(), text)
}

func wechatRSASign(cfg WxLocalConfig, text string) (string, error) {
	data, err := os.ReadFile(cfg.MerchantPrivateKeyPath)
	if err != nil {
		return "", fmt.Errorf("微信商户私钥不可用")
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return "", fmt.Errorf("微信商户私钥无效")
	}
	var key *rsa.PrivateKey
	if block.Type == "RSA PRIVATE KEY" {
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else {
		var parsed any
		parsed, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		key, _ = parsed.(*rsa.PrivateKey)
	}
	if err != nil || key == nil {
		return "", fmt.Errorf("微信商户 RSA 私钥无效")
	}
	hash := sha256.Sum256([]byte(text))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	return base64.StdEncoding.EncodeToString(signature), err
}

// WechatV3JSON 对应微信官方 Java SDK 的商户请求签名及平台证书响应验签。
func WechatV3JSON(ctx context.Context, method, path string, body map[string]any) (map[string]any, error) {
	cfg := LoadWxConfig()
	// 与原 AutoCertificateConfig 一样，商户调用前取得可用的平台证书。
	if _, err := WechatPlatformCertificate(ctx, cfg, ""); err != nil {
		return nil, err
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if method == "GET" {
		data = nil
	}
	res, err := wechatV3Request(ctx, cfg, method, path, data)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// 原微信 Java SDK WechatPay2Validator 校验响应时间与本机相差不足五分钟。
	stamp, err := strconv.ParseInt(res.Header.Get("Wechatpay-Timestamp"), 10, 64)
	if err != nil || time.Since(time.Unix(stamp, 0)).Abs() >= 5*time.Minute {
		return nil, fmt.Errorf("微信响应时间戳无效或过期")
	}
	cert, err := WechatPlatformCertificate(ctx, cfg, res.Header.Get("Wechatpay-Serial"))
	if err != nil {
		return nil, err
	}
	sig, err := base64.StdEncoding.DecodeString(res.Header.Get("Wechatpay-Signature"))
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256([]byte(res.Header.Get("Wechatpay-Timestamp") + "\n" + res.Header.Get("Wechatpay-Nonce") + "\n" + string(raw) + "\n"))
	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok || rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], sig) != nil {
		return nil, fmt.Errorf("微信响应验签失败")
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("微信支付 HTTP %d", res.StatusCode)
	}
	var out map[string]any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err = dec.Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func wechatV3Request(ctx context.Context, cfg WxLocalConfig, method, path string, data []byte) (*http.Response, error) {
	if cfg.MchID == "" || cfg.MerchantSerial == "" {
		return nil, fmt.Errorf("微信商户配置缺失")
	}
	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return nil, err
	}
	nonce, timestamp := hex.EncodeToString(nonceBytes), strconv.FormatInt(time.Now().Unix(), 10)
	signature, err := wechatRSASign(cfg, method+"\n"+path+"\n"+timestamp+"\n"+nonce+"\n"+string(data)+"\n")
	if err != nil {
		return nil, err
	}
	endpoint := os.Getenv("WECHAT_PAY_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.mch.weixin.qq.com"
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",timestamp="%s",serial_no="%s",signature="%s"`, cfg.MchID, nonce, timestamp, cfg.MerchantSerial, signature))
	return (&http.Client{Timeout: 30 * time.Second}).Do(req)
}
