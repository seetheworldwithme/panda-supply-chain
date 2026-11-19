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

	"github.com/golang-jwt/jwt/v5"
)

// KeyPair RSA 密钥对。首次启动自动生成并持久化到 keys/,重启后存量 token 仍有效
type KeyPair struct {
	Private      *rsa.PrivateKey
	Public       *rsa.PublicKey
	LegacyPublic *rsa.PublicKey // 迁移期间仅加载旧 Java 公钥，不复制旧私钥。
}

func legacyKey(k *KeyPair, keysDir string) (*KeyPair, error) {
	data, err := os.ReadFile(filepath.Join(keysDir, "jwt_legacy_public.pem"))
	if os.IsNotExist(err) {
		return k, nil
	}
	if err != nil {
		return nil, err
	}
	k.LegacyPublic, err = jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return nil, fmt.Errorf("旧 Java JWT 公钥无效: %w", err)
	}
	return k, nil
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
		return legacyKey(&KeyPair{Private: rsaKey, Public: &rsaKey.PublicKey}, keysDir)
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
	return legacyKey(&KeyPair{Private: key, Public: &key.PublicKey}, keysDir)
}
