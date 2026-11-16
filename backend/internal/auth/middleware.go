package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
)

// CurrentUser 请求级用户身份。登录用户由 JWT 解出,免登录验签请求为匿名(SignedAnon=true)
type CurrentUser struct {
	SignedAnon bool // Signature/AuthorizationId 验签通过的匿名请求(内部/前端登录前调用)
	UserID     int64
	Account    string
	Username   string
	CompanyID  int64
	Level       int
	AccountType string
	EnterpriseCode string // 对照 Java JWTInfo.enterpriseCode(企业账号才有)
}

type ctxKey int

const userKey ctxKey = 1

// Middleware 鉴权中间件:有 Bearer token 则解析注入用户;否则尝试 Signature/AuthorizationId 验签;
// 两者皆无则匿名(SecurityContext 为空,受保护接口返回 403)。
// 放行清单(permitAll)由 main.go 路由挂载顺序控制,中间件本身总是执行、不直接拦截
func Middleware(kp *KeyPair, secretDB *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := &CurrentUser{}
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				if claims, err := kp.Parse(strings.TrimPrefix(authHeader, "Bearer ")); err == nil {
					user.UserID = parseInt64(claims.UserID)
					user.Account = claims.Account
					user.Username = claims.Username
					user.CompanyID = parseInt64(claims.CompanyID)
					user.EnterpriseCode = claims.EnterpriseCode
				}
				// 解析失败按无 token 处理,继续走验签/匿名(与 Java 过滤器行为一致)
			}
			if user.UserID == 0 && !user.SignedAnon {
				user.SignedAnon = checkSignature(r, secretDB)
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey, user)))
		})
	}
}

// From 取当前请求用户
func From(ctx context.Context) *CurrentUser {
	if v, ok := ctx.Value(userKey).(*CurrentUser); ok {
		return v
	}
	return &CurrentUser{}
}

// checkSignature 对照 SafeVerify.checkSignature:
//   mySig = HmacSHA256(key=secret, msg=authorizationId+batchNumber+secret)
//   inSig = HmacSHA256(key=secret, msg=authorizationId+"hyqc"+secret)
// 传入 Signature 命中任一即通过。secret 存于 chain_auth_secret 表
func checkSignature(r *http.Request, secretDB *sql.DB) bool {
	authorizationId := r.Header.Get("AuthorizationId")
	signature := r.Header.Get("Signature")
	if signature == "" {
		signature = r.URL.Query().Get("Signature")
		authorizationId = firstNonEmpty(r.Header.Get("AuthorizationId"), r.URL.Query().Get("AuthorizationId"))
	}
	if authorizationId == "" || signature == "" {
		return false
	}
	var secret string
	err := secretDB.QueryRow("SELECT secret FROM chain_auth_secret WHERE authorization_id = ?", authorizationId).Scan(&secret)
	if err != nil {
		return false
	}
	batchNumber := r.Header.Get("BatchNumber")
	if hmacEqual(signature, hmacSHA256Hex(authorizationId+batchNumber+secret, secret)) ||
		hmacEqual(signature, hmacSHA256Hex(authorizationId+"hyqc"+secret, secret)) {
		return true
	}
	return false
}

func hmacSHA256Hex(msg, key string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(msg))
	return hex.EncodeToString(m.Sum(nil))
}

func hmacEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func parseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
