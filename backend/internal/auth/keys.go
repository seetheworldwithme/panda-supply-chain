// Package auth 登录认证与鉴权。对照 Java hyqc-auth security 层:
// RS256 JWT(私钥签发/公钥校验,Java 版公钥存 Redis,这里改为本地文件,语义一致)、
// BCrypt 密码、Signature/AuthorizationId 免登录验签
package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// KeyPair RSA 密钥对。首次启动自动生成并持久化到 keys/,重启后存量 token 仍有效
type KeyPair struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

func LoadOrGenerateKeyPair(keysDir string) (*KeyPair, error) {
	if err := os.MkdirAll(keysDir, 0o700); err != nil {
		return nil, err
	}
	privPath := filepath.Join(keysDir, "jwt_private.pem")
	data, err := os.ReadFile(privPath)
	if err == nil {
		block, _ := pem.Decode(data)
		if block == nil {
			return nil, fmt.Errorf("keys/jwt_private.pem 不是有效 PEM")
		}
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("jwt_private.pem 不是 RSA 私钥")
		}
		return &KeyPair{Private: rsaKey, Public: &rsaKey.PublicKey}, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	pemData := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	if err := os.WriteFile(privPath, pemData, 0o600); err != nil {
		return nil, err
	}
	return &KeyPair{Private: key, Public: &key.PublicKey}, nil
}
