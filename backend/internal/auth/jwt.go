package auth

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT claims 对照 Java JWTInfo.java: account/username/userId/companyId/enterpriseCode 五要素
type JWTClaims struct {
	Account        string `json:"account"`
	Username       string `json:"username"`
	JavaUsername   string `json:"userName,omitempty"`
	UserID         string `json:"userId"`
	CompanyID      string `json:"companyId"`
	EnterpriseCode string `json:"enterpriseCode"`
	jwt.RegisteredClaims
}

// Generate 签发 RS256 token,过期 7 天(Java 版 JWTHelper 同量级)
func (k *KeyPair) Generate(c JWTClaims) (string, error) {
	return k.generateWithLifetime(c, 7*24*time.Hour)
}

func (k *KeyPair) generateWithLifetime(c JWTClaims, lifetime time.Duration) (string, error) {
	now := time.Now()
	c.JavaUsername = c.Username
	c.RegisteredClaims = jwt.RegisteredClaims{
		Subject:   c.Account,
		ExpiresAt: jwt.NewNumericDate(now.Add(lifetime)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	return tok.SignedString(k.Private)
}

// Parse 校验并解析 token
func (k *KeyPair) Parse(tokenStr string) (*JWTClaims, error) {
	parse := func(key *rsa.PublicKey) (*jwt.Token, error) {
		return jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodRS256 {
				return nil, errors.New("非预期签名算法")
			}
			return key, nil
		})
	}
	tok, err := parse(k.Public)
	if errors.Is(err, jwt.ErrTokenSignatureInvalid) && k.LegacyPublic != nil {
		tok, err = parse(k.LegacyPublic)
	}
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*JWTClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("token 无效")
	}
	if claims.Account == "" {
		claims.Account = claims.Subject
	}
	if claims.Username == "" {
		claims.Username = claims.JavaUsername
	}
	return claims, nil
}
