package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT claims 对照 Java JWTInfo.java: account/username/userId/companyId/enterpriseCode 五要素
type JWTClaims struct {
	Account        string `json:"account"`
	Username       string `json:"username"`
	UserID         string `json:"userId"`
	CompanyID      string `json:"companyId"`
	EnterpriseCode string `json:"enterpriseCode"`
	jwt.RegisteredClaims
}

// Generate 签发 RS256 token,过期 7 天(Java 版 JWTHelper 同量级)
func (k *KeyPair) Generate(c JWTClaims) (string, error) {
	now := time.Now()
	c.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	return tok.SignedString(k.Private)
}

// Parse 校验并解析 token
func (k *KeyPair) Parse(tokenStr string) (*JWTClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("非预期签名算法")
		}
		return k.Public, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*JWTClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("token 无效")
	}
	return claims, nil
}
