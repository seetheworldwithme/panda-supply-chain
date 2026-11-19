package auth

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// 对照网关 handleCMBBodySignHeader；签名覆盖原始正文，读取后还原给业务入口。
func checkBankSignature(r *http.Request) bool {
	if r.Header.Get("CMB-BodySign") == "" {
		return false
	}
	body, err := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		return false
	}
	der, err := base64.StdEncoding.DecodeString(os.Getenv("busbank_publickey"))
	if err != nil {
		return false
	}
	parsed, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return false
	}
	key, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return false
	}
	sig, err := base64.StdEncoding.DecodeString(r.Header.Get("CMB-BodySign"))
	if err != nil {
		return false
	}
	hash := sha256.Sum256(body)
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], sig) == nil
}

func checkEsignSignature(r *http.Request) bool {
	if !strings.Contains(r.URL.Path, "getSignCallBack") {
		return false
	}
	body, err := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body))
	return err == nil && platform.VerifyEsignCallback(r, body)
}
