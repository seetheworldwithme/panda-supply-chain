package platform

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// ponytail: 当前为单商户串行缓存；支持多商户时再按商户拆分缓存及下载锁。
var wechatCertificates struct {
	sync.Mutex
	identity  [32]byte
	certs     map[string]*x509.Certificate
	refreshed time.Time
	attempted time.Time // 最近一次因未知序列号触发下载的时间
}

// WechatPlatformCertificate 复刻原 SDK 的商户签名下载、APIv3 解密及财付通 CA 校验。
// 证书在使用时每小时刷新；未知序列号立即刷新，回调与下单共用同一证书集合。
func WechatPlatformCertificate(ctx context.Context, cfg WxLocalConfig, serial string) (*x509.Certificate, error) {
	wechatCertificates.Lock()
	defer wechatCertificates.Unlock()
	identity := sha256.Sum256([]byte(cfg.MchID + "\x00" + cfg.MerchantSerial + "\x00" + cfg.MerchantPrivateKeyPath + "\x00" + cfg.APIV3Key + "\x00" + cfg.PlatformCert + "\x00" + os.Getenv("WECHAT_PAY_ENDPOINT")))
	if wechatCertificates.identity != identity {
		wechatCertificates.identity = identity
		wechatCertificates.certs = map[string]*x509.Certificate{}
		wechatCertificates.refreshed = time.Time{}
		// 保留已有配置的手动证书，同时允许以 APIv3 证书自动轮换。
		if cfg.PlatformCert != "" {
			block, _ := pem.Decode([]byte(cfg.PlatformCert))
			if block == nil {
				return nil, fmt.Errorf("微信平台证书无效")
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			wechatCertificates.certs[strings.ToUpper(cert.SerialNumber.Text(16))] = cert
			wechatCertificates.refreshed = time.Now()
		}
	}
	selectCert := func() *x509.Certificate {
		if serial != "" {
			return wechatCertificates.certs[strings.ToUpper(serial)]
		}
		var latest *x509.Certificate
		for _, cert := range wechatCertificates.certs {
			if time.Now().Before(cert.NotAfter) && (latest == nil || cert.NotAfter.After(latest.NotAfter)) {
				latest = cert
			}
		}
		return latest
	}
	cert := selectCert()
	if cert == nil && serial != "" {
		// 回调里的序列号来自未鉴权请求:未知序列号一分钟内只下载一次,
		// 否则随机序列号就能让下单/回调全部堵在这把全局锁上(下载超时 30s)。
		if time.Since(wechatCertificates.attempted) < time.Minute {
			return nil, fmt.Errorf("微信平台证书序号不存在")
		}
		wechatCertificates.attempted = time.Now()
	}
	if cert != nil && (cfg.APIV3Key == "" || time.Since(wechatCertificates.refreshed) < time.Hour) {
		return cert, nil
	}
	certs, err := downloadWechatCertificates(ctx, cfg)
	if err != nil {
		if cert != nil {
			wechatCertificates.refreshed = time.Now()
			return cert, nil // 原 SDK 定时刷新失败时继续使用已有证书。
		}
		return nil, err
	}
	wechatCertificates.certs, wechatCertificates.refreshed = certs, time.Now()
	if cert := selectCert(); cert != nil {
		return cert, nil
	}
	return nil, fmt.Errorf("微信平台证书序号不存在")
}

func downloadWechatCertificates(ctx context.Context, cfg WxLocalConfig) (map[string]*x509.Certificate, error) {
	if len(cfg.APIV3Key) != 32 {
		return nil, fmt.Errorf("微信 APIv3 密钥长度无效")
	}
	res, err := wechatV3Request(ctx, cfg, "GET", "/v3/certificates?algorithm_type=RSA", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("微信平台证书 HTTP %d", res.StatusCode)
	}
	var payload struct {
		Data []struct {
			Encrypt struct {
				Algorithm  string `json:"algorithm"`
				Nonce      string `json:"nonce"`
				Associated string `json:"associated_data"`
				Ciphertext string `json:"ciphertext"`
			} `json:"encrypt_certificate"`
		} `json:"data"`
	}
	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	block, err := aes.NewCipher([]byte(cfg.APIV3Key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM([]byte(wechatCertificateRoot)) {
		return nil, fmt.Errorf("财付通 CA 无效")
	}
	certs := map[string]*x509.Certificate{}
	for _, item := range payload.Data {
		encrypted := item.Encrypt
		if encrypted.Algorithm != "AEAD_AES_256_GCM" || len(encrypted.Nonce) != gcm.NonceSize() {
			return nil, fmt.Errorf("微信平台证书加密参数无效")
		}
		data, err := base64.StdEncoding.DecodeString(encrypted.Ciphertext)
		if err != nil {
			return nil, err
		}
		plain, err := gcm.Open(nil, []byte(encrypted.Nonce), data, []byte(encrypted.Associated))
		if err != nil {
			return nil, fmt.Errorf("微信平台证书解密失败")
		}
		pemBlock, _ := pem.Decode(plain)
		if pemBlock == nil {
			return nil, fmt.Errorf("微信平台证书格式无效")
		}
		cert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, err
		}
		if _, err = cert.Verify(x509.VerifyOptions{Roots: roots, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny}}); err != nil {
			return nil, fmt.Errorf("微信平台证书 CA 校验失败: %w", err)
		}
		certs[strings.ToUpper(cert.SerialNumber.Text(16))] = cert
	}
	if len(certs) == 0 {
		return nil, fmt.Errorf("微信未返回平台证书")
	}
	return certs, nil
}

// 财付通公开 CA，来源：原服务器 wechatpay-java-core 0.2.11 的 RSACertificateHandler。
var wechatCertificateRoot = `-----BEGIN CERTIFICATE-----
MIIEcDCCA1igAwIBAgIUG9QiDlDbwEsGrTl1SYRsAcPo69IwDQYJKoZIhvcNAQEL
BQAwcDELMAkGA1UEBhMCQ04xEzARBgNVBAoMCmlUcnVzQ2hpbmExHDAaBgNVBAsM
E0NoaW5hIFRydXN0IE5ldHdvcmsxLjAsBgNVBAMMJWlUcnVzQ2hpbmEgQ2xhc3Mg
MiBFbnRlcnByaXNlIENBIC0gRzMwHhcNMTcwODA5MDkxNTU1WhcNMzIwODA5MDkx
NTU1WjBeMQswCQYDVQQGEwJDTjETMBEGA1UEChMKVGVucGF5LmNvbTEdMBsGA1UE
CxMUVGVucGF5LmNvbSBDQSBDZW50ZXIxGzAZBgNVBAMTElRlbnBheS5jb20gUm9v
dCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALvnPD6k39BdPYAH
+6lnWPjuHH+2pcmZUf2E8cNFQFNr+ECRZylYV2iKyItCQt3I2/7VIDZl6aR9TE7n
sZrtSmOXCw635QOrq2yF9LTSDotAhf3ER0+216w3age/VzGcNVQpTf6gRCHCuQIk
8pe/oh06JagGvX0wERa+I6NfuG58ZHQY9d6RqLXKQl0Up95v73HDsG487z8k6jcn
qpGngmHQxdWiWRJugqxNRUD+awv2/DUsqGOffPX4jzJ6rLSJSlQXvuniDYxmaiaD
cK0bUbB5aM+1zMwogoHSYxWj/6B+vgcnHQCUrwGdiQR5+F+yRWzy5bO09IzaFgeO
PNPLPOsCAwEAAaOCARIwggEOMBIGA1UdEwEB/wQIMAYBAf8CAQAwDgYDVR0PAQH/
BAQDAgEGMCAGA1UdEQQZMBekFTATMREwDwYDVQQDDAhzd2JlLTI2NjAdBgNVHQ4E
FgQUTFo4GLdm9oHX52HcWnzuL4tui2gwHwYDVR0jBBgwFoAUK1vVxWgI69vN5LA5
MqJf/8dPmEUwRgYDVR0gBD8wPTA7BgoqgRyG7xcBAQECMC0wKwYIKwYBBQUHAgEW
H2h0dHBzOi8vd3d3Lml0cnVzLmNvbS5jbi9jdG5jcHMwPgYDVR0fBDcwNTAzoDGg
L4YtaHR0cDovL3RvcGNhLml0cnVzLmNvbS5jbi9jcmwvaXRydXNjMmNhZzMuY3Js
MA0GCSqGSIb3DQEBCwUAA4IBAQBwZhL/eiOQmMyo1D0IR9mu1DPWl5J3XXhjc4R6
mFgsN/FCeVP9M4U9y2FJH6i5Ha5YCecKGw5pwhA0rjZr/6okWwo22GF+nzI/gQiz
6ugAKs5VjFbeiEb04Ncz4HT8FP1idK3tyCjqCUTkLNt0U3tR7wy26hgOqlT2wCZ9
X4MfT8dUMdt9nCZx4ujN5yZOzaLOCHmzoGDGxgKg91bbu0TG2Yzd2ylhrxxRtFH9
aZ/J1x5UoF7uwhTM8P92DuAldWC1/bX1kciOtQvQEZeAy+9y/1BtFxoBnmDxnqkX
+lirIUYTLDaL7HaLrOLECUlaxZCU/Nkwm3tmqQxtCh+XQBdd
-----END CERTIFICATE-----`
