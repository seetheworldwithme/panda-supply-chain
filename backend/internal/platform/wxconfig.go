package platform

import (
	"encoding/json"
	"os"
)

// WxLocalConfig 对应旧 Java 微信凭据与 Nacos config.*，由运行时私有 JSON 提供。
// WECHAT_LOCAL_CONFIG 指向配置文件，每次读取以支持凭据和证书轮换，凭据不入仓库。
type WxLocalConfig struct {
	Secrets                map[string]string `json:"secrets"` // appid→secret（Java 源码硬编码的三组）
	PushOfficialAppid      string            `json:"pushOfficialAppid"`
	PushOfficialSecret     string            `json:"pushOfficialSecret"`
	MchID                  string            `json:"mchId"`
	MchSecret              string            `json:"mchSecret"`
	NotifyURL              string            `json:"notifyUrl"`
	APIV3Key               string            `json:"apiV3Key"`
	MerchantSerial         string            `json:"merchantSerial"`
	MerchantPrivateKeyPath string            `json:"merchantPrivateKeyPath"`
	PlatformCert           string            `json:"platformCert"`  // 微信支付平台证书 PEM
	MiniAppPayKey          string            `json:"miniAppPayKey"` // 小程序调起支付签名密钥（Java 硬编码）
	SendMsg                bool              `json:"sendMsg"`
}

var LoadWxConfig = func() WxLocalConfig {
	var c WxLocalConfig
	if p := os.Getenv("WECHAT_LOCAL_CONFIG"); p != "" {
		if data, err := os.ReadFile(p); err == nil {
			_ = json.Unmarshal(data, &c)
		}
	}
	return c
}

// Jnull Java 字符串拼接语义：null 拼成 "null"（本地 Nacos 缺省配置走该分支）。
func Jnull(s string) string {
	if s == "" {
		return "null"
	}
	return s
}

// JnullAny 任意值按 Java 拼接：nil→"null"，其余按字面。
func JnullAny(v any) string {
	if v == nil {
		return "null"
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}
