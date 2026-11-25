package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) PhoneCodeError(phone, code string) string {
	if strings.TrimSpace(code) == "" {
		return "验证码不能为空"
	}
	s.mu.Lock()
	sc, ok := s.sms[phone+"_PhoneYzm"]
	s.mu.Unlock()
	if !ok || !time.Now().Before(sc.expire) {
		return "未发送或验证码已过期"
	}
	if sc.value != code {
		return "验证码不匹配"
	}
	return ""
}

func (s *Service) ImageCodeError(r *http.Request, code string) string {
	if strings.TrimSpace(code) == "" {
		return "验证码不能为空"
	}
	s.mu.Lock()
	sc, ok := s.sms[imageCodeKey(r)]
	s.mu.Unlock()
	if !ok || !time.Now().Before(sc.expire) {
		return "未发送或验证码已过期"
	}
	if sc.value != code {
		return "验证码不匹配"
	}
	return ""
}

// AnonSignature 免登录预览协议签名，从本地授权表计算，源码不保存密钥。
func (s *Service) AnonSignature() (string, error) {
	const id = "wxda0b1c468b83bc99"
	var secret string
	if err := s.DB.QueryRow("SELECT secret FROM chain_auth_secret WHERE authorization_id=?", id).Scan(&secret); err != nil {
		return "", err
	}
	return hmacSHA256Hex(id+"hyqc"+secret, secret), nil
}

// 数字头像仍使用原签名预览协议；签名从本地授权表计算，源码不保存密钥。
func (s *Service) ChatAvatar(value any) (any, error) {
	text, ok := value.(string)
	if !ok {
		return value, nil
	}
	if _, err := strconv.ParseInt(text, 10, 32); err != nil {
		return value, nil
	}
	prefix := strings.TrimRight(os.Getenv("LEGACY_FRONT_PREFIX"), "/")
	if prefix == "" {
		return nil, fmt.Errorf("数字头像需要配置 LEGACY_FRONT_PREFIX")
	}
	sig, err := s.AnonSignature()
	if err != nil {
		return nil, err
	}
	return prefix + "/saas/saas-common/upload/imageShow?fileName=" + url.QueryEscape(text) + "&Signature=" + sig + "&AuthorizationId=wxda0b1c468b83bc99", nil
}

// MobileLogin 保留原表单协议：Java ValidateCodeFilter 的短信校验已注释，按原行为不校验。
func (s *Service) MobileLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		platform.Unknown(w)
		return
	}
	phone := strings.TrimSpace(r.FormValue("phone"))
	if phone == "" {
		platform.Fail(w, 2004, "用户名不能为空")
		return
	}
	u, err := s.selectByUserName2(phone)
	if err != nil {
		u, err = s.selectEUser(phone)
	}
	if err != nil {
		platform.Fail(w, 2003, "用户名或密码错误")
		return
	}
	tokens, err := s.loginTokens(u, r.Header.Get("applytype"))
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, tokens)
}

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	// The Java logout handler's Redis write is commented out. It returns this
	// Fastjson object with no data key, leaving stateless tokens valid until expiry.
	w.Header().Set("Content-Type", "text/json;charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "code": 200, "msg": "成功"})
}

func (s *Service) IssuePhoneCode(phone string) string {
	code := randomCode()
	s.storeCode(phone+"_PhoneYzm", code)
	return code
}

func compatBody(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	var m map[string]any
	d := json.NewDecoder(r.Body)
	d.UseNumber()
	if err := d.Decode(&m); err != nil || m == nil {
		platform.Unknown(w)
		return nil, false
	}
	return m, true
}
func compatString(v any) string {
	if v == nil {
		return "null"
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// UserPhoneCode differs from /code/phone: Java generates a code for any supplied
// phone and returns the request entity's id, not the code itself.
func (s *Service) UserPhoneCode(w http.ResponseWriter, r *http.Request) {
	m, ok := compatBody(w, r)
	if !ok {
		return
	}
	phone := compatString(m["phone"])
	code := s.IssuePhoneCode(phone)
	if platform.LoadWxConfig().SendMsg {
		platform.SendLegacySMS(r.Context(), phone, []string{code}, "", "")
	}
	platform.Success(w, m["id"])
}
func (s *Service) ChatPhoneCode(w http.ResponseWriter, r *http.Request) {
	m, ok := compatBody(w, r)
	if !ok {
		return
	}
	code := s.IssuePhoneCode(compatString(m["account"]))
	if platform.LoadWxConfig().SendMsg {
		platform.SendLegacySMS(r.Context(), compatString(m["account"]), []string{code}, "", "")
		code = ""
	}
	platform.Success(w, code)
}

func (s *Service) CreateExpiredToken(w http.ResponseWriter, r *http.Request) {
	u := From(r.Context())
	token, err := s.KP.generateWithLifetime(JWTClaims{Account: u.Account, Username: u.Username, UserID: compatString(u.UserID), CompanyID: compatString(u.CompanyID), EnterpriseCode: u.EnterpriseCode}, -time.Second)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, token)
}
