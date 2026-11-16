package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// Service 登录族接口。对照 Java:
//   /login                    Spring Security formLogin → UserServiceImpl + LoginSuccessHandler
//   /screenLogin              ScreenAuthenticationFilter(大屏账号)
//   /code/phone               短信验证码(本地 SendMsg=false → 明文返回,不发真短信)
//   /code/image               大屏图形验证码(前端仅 DPLogin 使用,Java 版实际校验已注释,此处同不校验)
//   /chainAuthUser/updatePwd  修改密码(BCrypt 原密码校验)
type Service struct {
	DB  *sql.DB // chain 库
	KP  *KeyPair
	mu  sync.Mutex
	sms map[string]smsCode // account → 验证码(内存替代 Redis,60 秒过期,行为一致)
}

type smsCode struct {
	value  string
	expire time.Time
}

func NewService(db *sql.DB, kp *KeyPair) *Service {
	return &Service{DB: db, KP: kp, sms: map[string]smsCode{}}
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
//   level>1: username=account, companyName=username(显示名)
//   level=2: username=显示名, companyName=启辰智链(applytype=zjjg 时 云账务后台系统)
func (s *Service) loginTokens(u *chainAuthUser, applyType string) (map[string]interface{}, error) {
	token, err := s.KP.Generate(JWTClaims{
		Account:   u.Account,
		Username:  u.Username,
		UserID:    fmt.Sprintf("%d", u.ID),
		CompanyID: fmt.Sprintf("%d", u.CompanyID.Int64),
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
		"xzqhCode":    u.XzqhCode.String,
		"xzqhName":    u.XzqhName.String,
		"accountType": accountType,
		"nickName":    u.NickName.String,
		"phone":       u.Phone.String,
		"wxHeadImg":   u.WxHeadImg.String,
		"hasRole":     s.hasRole(u.ID) || accountType == "5",
	}, nil
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
	var exists int
	if err := s.DB.QueryRow("SELECT COUNT(*) FROM chain_auth_user WHERE handle_result = 1 AND (phone = ? OR account = ?)", body.Account, body.Account).Scan(&exists); err != nil || exists == 0 {
		platform.Fail(w, 2002, "用户名不存在")
		return
	}
	code := randomCode()
	s.mu.Lock()
	s.sms[body.Account+"_PhoneYzm"] = smsCode{value: code, expire: time.Now().Add(60 * time.Second)}
	s.mu.Unlock()
	platform.Success(w, code)
}

// VerifyPhoneCode 校验短信验证码(登录/改密等场景复用;Java 现行校验被注释,此处仅在显式调用时校验)
func (s *Service) VerifyPhoneCode(account, code string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	sc, ok := s.sms[account+"_PhoneYzm"]
	return ok && time.Now().Before(sc.expire) && sc.value == code
}

// ImageCode POST /saas-auth/code/image 返回简单 PNG(Java 版校验已注释,前端主流程不使用)
func (s *Service) ImageCode(w http.ResponseWriter, r *http.Request) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 40))
	for x := 0; x < 100; x++ {
		for y := 0; y < 40; y++ {
			img.Set(x, y, color.RGBA{R: 240, G: 240, B: 240, A: 255})
		}
	}
	rng := func(n int) int { v, _ := rand.Int(rand.Reader, big.NewInt(int64(n))); return int(v.Int64()) }
	for i := 0; i < 60; i++ {
		img.Set(rng(100), rng(40), color.RGBA{R: uint8(rng(256)), G: uint8(rng(256)), B: uint8(rng(256)), A: 255})
	}
	w.Header().Set("Content-Type", "image/png")
	_ = png.Encode(w, img)
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
