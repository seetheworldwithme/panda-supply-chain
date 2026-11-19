package auth

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestLegacyJavaPublicKeyAcceptsOriginalClaims(t *testing.T) {
	legacy, err := LoadOrGenerateKeyPair(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	der, err := x509.MarshalPKIXPublicKey(legacy.Public)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "jwt_legacy_public.pem"), pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), 0600); err != nil {
		t.Fatal(err)
	}
	current, err := LoadOrGenerateKeyPair(dir)
	if err != nil {
		t.Fatal(err)
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"sub": "fixture", "userName": "验收", "userId": "1", "companyId": "", "enterpriseCode": ""}).SignedString(legacy.Private)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := current.Parse(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Account != "fixture" || claims.Username != "验收" {
		t.Fatalf("Java 字段映射失败: %+v", claims)
	}
	other, err := LoadOrGenerateKeyPair(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = other.Parse(token); err == nil {
		t.Fatal("不能信任未配置的公钥")
	}
}
