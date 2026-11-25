package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// Service 登录族接口。对照 Java:
//
//	/login                    Spring Security formLogin → UserServiceImpl + LoginSuccessHandler
//	/screenLogin              ScreenAuthenticationFilter(大屏账号)
//	/code/phone               短信验证码(本地 SendMsg=false → 明文返回,不发真短信)
//	/code/image               大屏图形验证码(前端仅 DPLogin 使用,Java 版实际校验已注释,此处同不校验)
//	/chainAuthUser/updatePwd  修改密码(BCrypt 原密码校验)
type Service struct {
	DB             *sql.DB // chain 库
	KP             *KeyPair
	mu             sync.Mutex
	sms            map[string]smsCode   // account → 验证码(内存替代 Redis,60 秒过期,行为一致)
	verifiedPhones map[string]time.Time // 手机号 → 刚通过短信码校验的时间(忘记密码第二步据此限定可改的账号)
}

type smsCode struct {
	value  string
	expire time.Time
}

const phoneVerifiedTTL = 10 * time.Minute

// MarkPhoneVerified 记录"该手机号刚通过短信验证码校验"。忘记密码分两步:valCode 校验短信码并返回
// 该手机号下的账号,updPwd 再按 id 改密码——第二步原先不核对 id 与已验证手机号的关系,
// 未登录即可改任何账号(含平台管理员)的密码。
func (s *Service) MarkPhoneVerified(phone string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for p, t := range s.verifiedPhones {
		if now.Sub(t) > phoneVerifiedTTL {
			delete(s.verifiedPhones, p)
		}
	}
	s.verifiedPhones[phone] = now
}

// PhoneVerifiedRecently 该手机号最近是否通过过短信码校验
func (s *Service) PhoneVerifiedRecently(phone string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.verifiedPhones[phone]
	return ok && time.Since(t) <= phoneVerifiedTTL
}

// storeCode 写入验证码并清理过期项。键含客户端可控的 IP(X-Forwarded-For)/手机号,
// 未鉴权接口可以无限造新键:只增不删会把内存越撑越大。
func (s *Service) storeCode(key, code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for k, sc := range s.sms {
		if !now.Before(sc.expire) {
			delete(s.sms, k)
		}
	}
	s.sms[key] = smsCode{value: code, expire: now.Add(60 * time.Second)}
}

func NewService(db *sql.DB, kp *KeyPair) *Service {
	return &Service{DB: db, KP: kp, sms: map[string]smsCode{}, verifiedPhones: map[string]time.Time{}}
}

// chainAuthUser 用户行。列名与 chain_auth_user 表一致(MyBatis-Plus 驼峰映射规则)
type chainAuthUser struct {
	ID          int64
	Account     string
	Username    string
	Password    string
	CompanyID   sql.NullInt64
	Level       int
	XzqhCode    sql.NullString
	XzqhName    sql.NullString
	AccountType sql.NullString
	Phone       sql.NullString
	NickName    sql.NullString
	WxHeadImg   sql.NullString
}

const userCols = "id, account, username, password, company_id, level, xzqh_code, xzqh_name, account_type, phone, nickname, wx_face_img"

func scanUser(row interface{ Scan(...any) error }) (*chainAuthUser, error) {
	var u chainAuthUser
	err := row.Scan(&u.ID, &u.Account, &u.Username, &u.Password, &u.CompanyID, &u.Level,
		&u.XzqhCode, &u.XzqhName, &u.AccountType, &u.Phone, &u.NickName, &u.WxHeadImg)
	return &u, err
}

// selectByUserName2 对照 ChainAuthUserMapper.selectByUserName2
func (s *Service) selectByUserName2(name string) (*chainAuthUser, error) {
	row := s.DB.QueryRow("SELECT "+userCols+" FROM chain_auth_user WHERE account_type != 4 AND handle_result = 1 AND (phone = ? OR account = ?) ORDER BY level DESC LIMIT 1", name, name)
	return scanUser(row)
}

// selectEUser 兜底:绑定 e签宝的企业账号
func (s *Service) selectEUser(name string) (*chainAuthUser, error) {
	row := s.DB.QueryRow("SELECT "+userCols+" FROM chain_auth_user WHERE eqb_per_id IS NOT NULL AND handle_result = 1 AND (phone = ? OR account = ?) LIMIT 1", name, name)
	return scanUser(row)
}

// selectScreenUser 对照 UserServiceImpl.loadByScreenUsername(白名单账号不限 account_type)
func (s *Service) selectScreenUser(name string) (*chainAuthUser, error) {
	extra := " AND account_type = 4"
	if name == "gswhy" || name == "18636880777" {
		extra = ""
	}
	row := s.DB.QueryRow("SELECT "+userCols+" FROM chain_auth_user WHERE handle_result = 1 AND (phone = ? OR account = ?)"+extra+" LIMIT 1", name, name)
	return scanUser(row)
}

// enterpriseName 查企业名(LoginSuccessHandler companyName 来源)
func (s *Service) enterpriseName(companyID int64) string {
	var name sql.NullString
	_ = s.DB.QueryRow("SELECT enterprise_name FROM chain_basic_enterprise WHERE id = ?", companyID).Scan(&name)
	return name.String
}

func (s *Service) hasRole(userID int64) bool {
	var n int
	_ = s.DB.QueryRow("SELECT COUNT(*) FROM chain_auth_use_role WHERE user_id = ?", userID).Scan(&n)
	return n > 0
}

// loginTokens 构造登录成功响应 data。分支照 LoginSuccessHandler:
//
//	level>1: username=account, companyName=username(显示名)
//	level=2: username=显示名, companyName=启辰智链(applytype=zjjg 时 云账务后台系统)
func (s *Service) loginTokens(u *chainAuthUser, applyType string) (map[string]interface{}, error) {
	enterpriseCode := ""
	if u.CompanyID.Valid && u.CompanyID.Int64 > 0 {
		// enterpriseCode 五要素(四期 getSystemCode 依赖;照 Java userInfoHelper 取企业表)
		_ = s.DB.QueryRow("SELECT enterprise_code FROM chain_basic_enterprise WHERE id = ?", u.CompanyID.Int64).Scan(&enterpriseCode)
	}
	token, err := s.KP.Generate(JWTClaims{
		Account:        u.Account,
		Username:       u.Username,
		UserID:         fmt.Sprintf("%d", u.ID),
		CompanyID:      fmt.Sprintf("%d", u.CompanyID.Int64),
		EnterpriseCode: enterpriseCode,
	})
	if err != nil {
		return nil, err
	}
	username := u.Username
	companyName := ""
	level := u.Level
	if level > 1 {
		username = u.Account
		companyName = u.Username
	} else if u.CompanyID.Valid && u.CompanyID.Int64 > 0 {
		companyName = s.enterpriseName(u.CompanyID.Int64) // level=1 企业账号:公司名来自企业表
	}
	if level == 2 {
		username = u.Username
		companyName = "启辰智链"
		if applyType == "zjjg" {
			companyName = "云账务后台系统"
		}
	}
	accountType := ""
	if u.AccountType.Valid {
		accountType = u.AccountType.String
	}
	return map[string]interface{}{
		"token":       token,
		"username":    username,
		"userId":      u.ID,
		"companyName": companyName,
		"companyId":   fmt.Sprintf("%d", u.CompanyID.Int64),
		"level":       fmt.Sprintf("%d", level),
		"xzqhCode":    nullableString(u.XzqhCode),
		"xzqhName":    nullableString(u.XzqhName),
		"accountType": nullableString(u.AccountType),
		// Java UserServiceImpl 未将昵称和头像复制到 UserInfo。
		"nickName":  nil,
		"phone":     nullableString(u.Phone),
		"wxHeadImg": nil,
		"hasRole":   s.hasRole(u.ID) || accountType == "5",
	}, nil
}

func nullableString(v sql.NullString) any {
	if !v.Valid {
		return nil
	}
	return v.String
}

// Login POST /saas-auth/login (form-urlencoded: username/password/code/checked)
// Java 版 ValidateCodeFilter 的验证码校验已被注释,实际只校验账密,此处一致
func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		platform.FailMsg(w, "参数错误")
		return
	}
	name := r.FormValue("username")
	password := r.FormValue("password")
	if name == "" || password == "" {
		platform.Fail(w, 2003, "用户名或密码错误")
		return
	}
	u, err := s.selectByUserName2(name)
	if err == sql.ErrNoRows {
		u, err = s.selectEUser(name)
	}
	if err != nil {
		platform.Fail(w, 2003, "用户名或密码错误")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		platform.Fail(w, 2003, "用户名或密码错误")
		return
	}
	if u.CompanyID.Valid && u.CompanyID.Int64 > 0 {
		_ = s.enterpriseName(u.CompanyID.Int64) // 企业存在性随查(Java 版另有 eqbConfig 缓存,Go 版按需再补)
	}
	tokens, err := s.loginTokens(u, r.Header.Get("applytype"))
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, tokens)
}

// ScreenLogin POST /saas-auth/screenLogin (form: username/password/code) 大屏账号
func (s *Service) ScreenLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		platform.FailMsg(w, "参数错误")
		return
	}
	name := r.FormValue("username")
	password := r.FormValue("password")
	u, err := s.selectScreenUser(name)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
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

// PhoneCode POST /saas-auth/code/phone (JSON: {account})
// 本地 SendMsg=false 等价形态:不发短信,验证码明文返回 data,60 秒有效
func (s *Service) PhoneCode(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Account string `json:"account"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Account == "" {
		platform.Fail(w, 2002, "用户名不存在")
		return
	}
	var phone sql.NullString
	if err := s.DB.QueryRow("SELECT phone FROM chain_auth_user WHERE handle_result = 1 AND (phone = ? OR account = ?) ORDER BY level DESC LIMIT 1", body.Account, body.Account).Scan(&phone); err != nil {
		platform.Fail(w, 2002, "用户名不存在")
		return
	}
	code := s.IssuePhoneCode(body.Account)
	if platform.LoadWxConfig().SendMsg {
		platform.SendLegacySMS(r.Context(), phone.String, []string{code}, "", "")
		code = ""
	}
	platform.Success(w, code)
}

// VerifyPhoneCode 校验短信验证码(登录/改密等场景复用;Java 现行校验被注释,此处仅在显式调用时校验)
func (s *Service) VerifyPhoneCode(account, code string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	sc, ok := s.sms[account+"_PhoneYzm"]
	return ok && time.Now().Before(sc.expire) && sc.value == code
}

// imageCodeKey 与 Java 一样按客户端地址缓存验证码，去掉 IPv6 冒号。
func imageCodeKey(r *http.Request) string {
	ip := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return strings.ReplaceAll(ip, ":", "") + "_ImgYzm"
}

// ImageCode 保留 JPEG、尺寸参数、四位数字和 60 秒有效期。
func (s *Service) ImageCode(w http.ResponseWriter, r *http.Request) {
	width, height := 67, 23
	for key, dest := range map[string]*int{"width": &width, "height": &height} {
		if r.URL.Query().Has(key) {
			n, err := strconv.Atoi(r.URL.Query().Get(key))
			// 图形固定按 67×23 绘制,画布再大也只是留白;未鉴权接口要封顶,防单请求申请几十 MB
			if err != nil || n <= 0 || n > 512 {
				platform.Unknown(w)
				return
			}
			*dest = n
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{R: 240, G: 240, B: 240, A: 255})
		}
	}
	rng := func(n int) int { v, _ := rand.Int(rand.Reader, big.NewInt(int64(n))); return int(v.Int64()) }
	code := randomCode()[:4]
	// ponytail: 四位数字用内置点阵绘制，避免只为验证码引入字体依赖。
	digits := [10][5]byte{{7, 5, 5, 5, 7}, {2, 6, 2, 2, 7}, {7, 1, 7, 4, 7}, {7, 1, 7, 1, 7}, {5, 5, 7, 1, 1}, {7, 4, 7, 1, 7}, {7, 4, 7, 5, 7}, {7, 1, 1, 1, 1}, {7, 5, 7, 5, 7}, {7, 5, 7, 1, 7}}
	for i, ch := range code {
		ink := color.RGBA{R: uint8(20 + rng(110)), G: uint8(20 + rng(110)), B: uint8(20 + rng(110)), A: 255}
		for y, row := range digits[ch-'0'] {
			for x := 0; x < 3; x++ {
				if row&(1<<(2-x)) != 0 {
					for dx := 0; dx < 2; dx++ {
						for dy := 0; dy < 3; dy++ {
							img.Set(6+13*i+x*2+dx, 2+y*3+dy, ink)
						}
					}
				}
			}
		}
	}
	for i := 0; i < 155; i++ {
		img.Set(rng(width), rng(height), color.RGBA{R: 180, G: 190, B: 180, A: 255})
	}
	s.storeCode(imageCodeKey(r), code)
	w.Header().Set("Cache-Control", "no-store, no-cache")
	w.Header().Set("Content-Type", "image/jpeg")
	_ = jpeg.Encode(w, img, &jpeg.Options{Quality: 75})
}

// UpdatePwd POST /saas-auth/chainAuthUser/updatePwd (JSON: {id?, password, newPassword})
// 对照 ChainAuthUserController.updatePwd: 默认改当前登录人,带 id 则管理员改指定账号
func (s *Service) UpdatePwd(w http.ResponseWriter, r *http.Request) {
	cu := From(r.Context())
	if cu.UserID == 0 {
		platform.NotAuth(w)
		return
	}
	var body struct {
		ID          *int64 `json:"id"`
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		platform.FailMsg(w, "参数错误")
		return
	}
	if body.NewPassword == "" { // 空密码落库后该账号将永远无法登录(登录拒绝空密码)
		platform.FailMsg(w, "新密码不能为空")
		return
	}
	var id int64
	var hash string
	if body.ID != nil && *body.ID > 0 {
		id = *body.ID
		err := s.DB.QueryRow("SELECT password FROM chain_auth_user WHERE id = ?", id).Scan(&hash)
		if err != nil {
			platform.FailMsg(w, "账号不存在")
			return
		}
	} else {
		id = cu.UserID
		err := s.DB.QueryRow("SELECT password FROM chain_auth_user WHERE id = ?", id).Scan(&hash)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Password)) != nil {
		platform.FailMsg(w, "密码错误")
		return
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := s.DB.Exec("UPDATE chain_auth_user SET password = ? WHERE id = ?", string(newHash), id); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "修改成功")
}

func randomCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

// HashPassword BCrypt 加密(组织管理新增/重置用户复用登录同款编码)
func HashPassword(pw string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(h)
}

// RandomCode 6 位随机码(验证码/邀请码同源)
func RandomCode() string { return randomCode() }
