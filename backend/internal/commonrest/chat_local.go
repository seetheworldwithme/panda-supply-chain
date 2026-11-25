package commonrest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// 这些入口只读本地用户和企业，不需要向微信换取授权码。
func (s *Service) registerChatLocal(mux *http.ServeMux) {
	mux.HandleFunc("POST /saas-auth/chainAuthChat/getAccessToken", s.chatGetAccessToken)
	// 服务器运行包中两个方法正常编译，直接连接同一业务实现。
	mux.HandleFunc("POST /saas-auth/chainAuthChat/loginPhone", func(w http.ResponseWriter, r *http.Request) { s.chatPhoneLogin(w, r, false) })
	mux.HandleFunc("POST /saas-auth/chainAuthChat/loginChatPhone", func(w http.ResponseWriter, r *http.Request) { s.chatPhoneLogin(w, r, true) })
	mux.HandleFunc("POST /saas-auth/wechat/updByOpenid", s.chatUpdateOfficial)
	mux.HandleFunc("POST /saas-auth/chainAuthChat/updByOpenId", s.chatUpdateProfile)

	mux.HandleFunc("POST /saas-auth/chainAuthChat/getAppInfo", func(w http.ResponseWriter, r *http.Request) {
		m := bodyMap(r)
		query, args := "SELECT * FROM chain_auth_chat_authorization WHERE app_identification=?", []any{m["appIdentification"]}
		if m["appIdentification"] == nil {
			query, args = "SELECT * FROM chain_auth_chat_authorization WHERE app_identification IS NULL", nil
		}
		rows, err := queryRows(s.pick(r), query, args...)
		if err != nil || len(rows) > 1 {
			platform.Unknown(w)
			return
		}
		if len(rows) == 0 {
			platform.Success(w, nil)
			return
		}
		info := authRestEnts["ChainAuthChatAuthorization"].shapeOne(rows[0])
		info["secret"] = nil // AppSecret 不能下发给调用方(保留键,响应形状不变)
		platform.Success(w, info)
	})
	for _, mode := range []string{"loginByAccount", "buslogin", "loginCarSystem"} {
		path := "/saas-auth/chainAuthChat/" + mode
		if mode == "loginCarSystem" {
			path = "/saas-auth/chainAuthUser/" + mode
		}
		mux.HandleFunc("POST "+path, func(w http.ResponseWriter, r *http.Request) { s.chatLocalLogin(w, r, mode) })
	}
}

func (s *Service) chatLocalLogin(w http.ResponseWriter, r *http.Request, mode string) {
	m := bodyMap(r)
	if s.Auth == nil {
		platform.Unknown(w)
		return
	}
	if mode == "loginCarSystem" && r.URL.Query().Has("phoneLogin") {
		if msg := s.Auth.ImageCodeError(r, str(m, "code")); msg != "" {
			platform.Fail(w, 602, msg)
			return
		}
	}
	query, args := "SELECT * FROM chain_auth_user WHERE account_type!=4 AND handle_result=1 AND (phone=? OR account=?)", []any{m["account"], m["account"]}
	if mode == "buslogin" {
		query, args = "SELECT * FROM chain_auth_user WHERE account_type!=4 AND level=? AND phone=? ORDER BY level DESC", []any{m["level"], m["phone"]}
	}
	rows, err := queryRows(s.pick(r), query, args...)
	if err != nil || (mode != "buslogin" && len(rows) > 1) {
		platform.Unknown(w)
		return
	}
	if mode == "loginCarSystem" && (len(rows) == 0 || bcrypt.CompareHashAndPassword([]byte(str(rows[0], "password")), []byte(str(m, "password"))) != nil || !strings.Contains(str(rows[0], "applys"), "2")) {
		platform.Fail(w, 602, "用户名或密码错误，登录失败")
		return
	}
	if len(rows) == 0 {
		platform.Unknown(w)
		return
	}
	result, err := s.chatUserInfo(r, rows[0])
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, result)
}

// 对照 UserInfoHelper，与表单登录的响应字段不同，不能复用 loginTokens。
func (s *Service) chatUserInfo(r *http.Request, row map[string]any) (map[string]any, error) {
	if row["is_agree_call"] == nil {
		return nil, nil
	} // Java Boolean 拆箱异常被 helper 捕获后返回 null。
	u := authRestEnts["ChainAuthUser"].shapeOne(row)
	var enterprise map[string]any
	if row["company_id"] != nil {
		e, err := platform.GetByID(s.pick(r), "chain_basic_enterprise", row["company_id"])
		if err != nil {
			return nil, err
		}
		enterprise = authRestEnts["ChainBasicEnterprise"].shapeOne(e)
	}
	companyID, companyName, enterpriseCode := "", "", ""
	var rotation any = false
	if enterprise != nil {
		companyID = str(u, "companyId")
		companyName = str(enterprise, "enterpriseName")
		enterpriseCode = str(enterprise, "enterpriseCode")
		rotation = enterprise["rotationFlag"]
	}
	token, err := s.Auth.KP.Generate(auth.JWTClaims{Account: str(u, "account"), Username: str(u, "username"), UserID: str(u, "id"), CompanyID: companyID, EnterpriseCode: enterpriseCode})
	if err != nil {
		return nil, err
	}
	username := u["username"]
	if num(u["level"]) > 1 {
		username = u["account"]
		companyName = str(u, "username")
	}
	if num(u["level"]) == 2 {
		username = u["username"]
		companyName = "云通大宗数字交易平台"
	}
	agree := 0
	if u["isAgreeCall"] == true {
		agree = 1
	}
	ware := u["wareIdentity"]
	if ware == nil {
		ware = 0
	}
	avatar, err := s.Auth.ChatAvatar(u["wxFaceImg"])
	if err != nil {
		return nil, err
	}
	return map[string]any{"token": token, "username": username, "nickname": u["nickname"], "account": u["account"], "companyName": companyName, "companyId": companyID, "level": str(u, "level"), "xzqhCode": u["xzqhCode"], "xzqhName": u["xzqhName"], "accountType": u["accountType"], "wxFaceImg": avatar, "isAgreeCall": agree, "userId": u["id"], "phone": u["phone"], "wareIdentity": ware, "memberType": u["memberType"], "expirationTime": u["expirationTime"], "packageId": u["packageId"], "rotationFlag": rotation, "enterprise": enterprise}, nil
}

func (s *Service) chatUpdateOfficial(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	e := authRestEnts["ChainAuthUser"]
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	if m["openid"] == nil {
		platform.Success(w, nil)
		return
	}
	err := platform.WithTx(s.pick(r), func(tx *sql.Tx) error {
		// MP update(entity,wrapper):SET=实体非 null 字段(不含主键),无自动填充;单语句。
		vals := insertVals(pick(m, e.dbcols))
		delete(vals, "id")
		if len(vals) == 0 {
			return errors.New("java empty SET clause syntax error")
		}
		cols := platform.SortedKeys(vals)
		sets := make([]string, 0, len(cols))
		args := make([]any, 0, len(cols)+1)
		for _, k := range cols {
			sets = append(sets, k+" = ?")
			args = append(args, vals[k])
		}
		args = append(args, m["openid"])
		if _, err := tx.Exec("UPDATE chain_auth_user SET "+strings.Join(sets, ",")+" WHERE openid = ? AND channel = 2", args...); err != nil {
			return err
		}
		_, err := tx.Exec("DELETE FROM chain_auth_user WHERE phone IS NULL AND openid IS NULL AND official_openid=?", m["officialOpenid"])
		return err
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

func (s *Service) chatPhoneLogin(w http.ResponseWriter, r *http.Request, chat bool) {
	m := bodyMap(r)
	if s.Auth == nil || m["code"] == nil {
		platform.Unknown(w)
		return
	}
	phone, code := str(m, "phone"), str(m, "code")
	if message := s.Auth.PhoneCodeError(phone, code); message != "" && !(chat && phone == "15012345678") {
		if chat {
			platform.Fail(w, 602, "验证码错误")
		} else {
			platform.Fail(w, 601, message)
		}
		return
	}
	rows, err := queryRows(s.pick(r), "SELECT * FROM chain_auth_user WHERE (phone=? OR account=?) AND channel=2 ORDER BY level DESC", phone, phone)
	if err != nil {
		platform.Unknown(w)
		return
	}
	fresh := len(rows) == 0
	tx, err := s.pick(r).BeginTx(r.Context(), nil)
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer tx.Rollback()
	var row map[string]any
	if fresh {
		if len(phone) < 7 {
			platform.Unknown(w)
			return
		}
		row = map[string]any{"phone": phone, "account": uuidV4(), "level": json.Number("1"), "channel": "2", "nickname": "匿名用户" + phone[7:], "is_agree_call": false}
		if chat {
			row["openid"] = m["openid"]
			row["app_client_id"] = m["pushClientid"]
		}
		id, err := platform.TxInsertNoFill(tx, "chain_auth_user", insertVals(row))
		if err != nil {
			platform.Unknown(w)
			return
		}
		row["id"] = id
	} else {
		row = rows[0]
	}
	card := map[string]any{"userId": row["id"]}
	has, err := s.callJSON(r, "/saas-warehouse/busBusinessCard/searchHasDefault", card)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if has == "false" {
		card["phone"] = phone
		if _, err = s.callJSON(r, "/saas-warehouse/busBusinessCard/add", card); err != nil {
			platform.Unknown(w)
			return
		}
	}
	if err := platform.RedisSetHash(r.Context(), "userId"+str(row, "id"), map[string]string{
		"phone": phone, "username": platform.JnullAny(row["username"]),
		"wxHeadImg": platform.JnullAny(row["wx_face_img"]), "level": platform.JnullAny(row["level"]),
	}, 0); err != nil {
		platform.Unknown(w)
		return
	}
	if chat && !fresh {
		vals := map[string]any{}
		if m["openid"] != nil {
			vals["openid"] = m["openid"]
			row["openid"] = m["openid"]
		}
		if m["pushClientid"] != nil {
			vals["app_client_id"] = m["pushClientid"]
			row["app_client_id"] = m["pushClientid"]
		}
		if err := platform.TxUpdateByIDNoFill(tx, "chain_auth_user", row["id"], vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	// 响应保留 Java 新建对象中未赋值字段的 null，不回填数据库默认值。
	info, err := s.chatUserInfo(r, row)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if chat {
		if info == nil {
			platform.Unknown(w)
			return
		}
		benefits, err := queryRows(s.pick(r), `SELECT is_show_etc_phone isShowEtcPhone,is_show_etc_name isShowEtcName,is_show_website isShowWebsite,is_see_data_briefing isSeeDataBriefing,is_see_card_visitor isSeeCardVisitor,is_see_warehouse_visitor isSeeWarehouseVisitor,is_see_testing_report isSeeTestingReport,is_see_price_trend isSeePriceTrend FROM bus_member_benefits b,chain_auth_user u WHERE b.member_type=u.member_type AND u.id=?`, row["member_type"])
		if err != nil {
			platform.Unknown(w)
			return
		}
		var benefit any = map[string]any{"isShowEtcPhone": 0, "isShowEtcName": 0}
		if len(benefits) > 0 {
			benefit = benefits[0]
		}
		care := false
		if str(row, "official_openid") != "" {
			care, err = s.chatOfficialCare(r, str(row, "official_openid"))
			if err != nil {
				platform.Unknown(w)
				return
			}
		}
		expiration := any(nil)
		if value, ok := row["expiration_time"].(string); ok && len(value) >= 10 {
			expiration = value[:10]
		}
		info["newUser"], info["openid"], info["officialOpenid"], info["appClientId"], info["isCareOfficial"], info["expirationTimeStr"], info["MemberBenefits"] = fresh, row["openid"], row["official_openid"], row["app_client_id"], care, expiration, benefit
	}
	if err := tx.Commit(); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, info)
}

// 注册积分和用户修改在同一本地事务内，失败时不能扣掉全局积分池。
func (s *Service) chatUpdateProfile(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	e := authRestEnts["ChainAuthUser"]
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	err := platform.WithTx(s.pick(r), func(tx *sql.Tx) error {
		// Java getRegistScore 的 ${} 插值:null→空串→SQL 语法错误→500 整体回滚(8765 实测)。
		if m["id"] == nil {
			return errors.New("java scoreLogs: empty ${userId} interpolation")
		}
		var count int
		if err := tx.QueryRow("SELECT count(*) FROM bus_score_log WHERE user_id=? AND score_get_type=1", m["id"]).Scan(&count); err != nil {
			return err
		}
		if count == 0 {
			// ${phone} 原样拼 SQL:空/非数字 → 语法错误 500;数字字面量才成功。
			phone := str(m, "phone")
			if m["phone"] == nil || !javaNumericLiteral(phone) {
				return errors.New("java saveScoreLog: empty/non-numeric ${phone} interpolation")
			}
			if _, err := registScoreTx(tx, m["id"], m["phone"], false); err != nil {
				return err
			}
		}
		// MP update(entity,wrapper):SET=实体非 null 字段且不含主键(实测),无自动填充;
		// allEq(Map) null2IsNull=true → 缺 openid 时 WHERE openid IS NULL(批量更新,照实复刻)。
		vals := insertVals(pick(m, e.dbcols))
		delete(vals, "id")
		if len(vals) == 0 {
			return errors.New("java empty SET clause syntax error")
		}
		cols := platform.SortedKeys(vals)
		sets := make([]string, 0, len(cols))
		args := make([]any, 0, len(cols)+1)
		for _, k := range cols {
			sets = append(sets, k+" = ?")
			args = append(args, vals[k])
		}
		where := "openid IS NULL"
		if m["openid"] != nil {
			where = "openid = ?"
			args = append(args, m["openid"])
		}
		_, err := tx.Exec("UPDATE chain_auth_user SET "+strings.Join(sets, ",")+" WHERE "+where+" AND channel = 2", args...)
		return err
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// javaNumericLiteral ${phone} 原样拼接进 SQL 时,纯数字串才是合法字面量。
func javaNumericLiteral(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
