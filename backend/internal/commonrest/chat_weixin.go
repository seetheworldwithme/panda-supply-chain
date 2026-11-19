package commonrest

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// 微信授权族入口，以服务器实际运行 JAR 为准；本机重编译包的编译缺陷不属于兼容契约。

func (s *Service) registerChatWeixin(mux *http.ServeMux) {
	mux.HandleFunc("POST /saas-auth/chainAuthChat/login", s.chatLogin)
	mux.HandleFunc("POST /saas-auth/chainAuthChat/selPhoneMsg", s.chatSelPhoneMsg)
	mux.HandleFunc("POST /saas-auth/chainAuthChat/loginChat", s.chatLoginChat)
	mux.HandleFunc("POST /saas-auth/chainAuthChat/loginChatApp", s.chatLoginChatApp)
	mux.HandleFunc("POST /saas-auth/chainAuthChat/getWxConfig", s.chatWxConfig)
	mux.HandleFunc("POST /saas-auth/chainAuthUser/reg", s.chatReg)
	mux.HandleFunc("POST /saas-auth/chainAuthUser/updMerExpire", s.chatUpdMerExpire)
	mux.HandleFunc("GET /saas-auth/wechat/auth", s.wechatAuth)
	mux.HandleFunc("GET /saas-auth/wechat/authorize", s.wechatAuthorize)
	mux.HandleFunc("POST /saas-auth/wechat/userInfo", s.wechatUserInfo)
	mux.HandleFunc("POST /saas-auth/code/refundMoney", s.codeRefundMoney)
	mux.HandleFunc("/saas-common/start", func(w http.ResponseWriter, r *http.Request) { s.startTLTask(); w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("/saas-common/stopById", func(w http.ResponseWriter, r *http.Request) {
		s.stopTLTasks(r.FormValue("taskId"), false)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/saas-common/stopAll", func(w http.ResponseWriter, r *http.Request) { s.stopTLTasks("", true); w.WriteHeader(http.StatusOK) })
}

// jstr Java Map.get 后字符串拼接语义：缺键/null → "null"。
func jstr(m map[string]any, k string) string {
	if v, ok := m[k]; ok && v != nil {
		return str(m, k)
	}
	return "null"
}

func (s *Service) jscode2session(r *http.Request, appid, secret, code string) (map[string]any, error) {
	q := url.Values{"appid": {appid}, "js_code": {code}, "grant_type": {"authorization_code"}, "secret": {secret}}
	return s.wxJSON(r, "GET", "https://api.weixin.qq.com/sns/jscode2session?"+q.Encode(), nil)
}

// chatAppRows getOne(app_identification=?)：null 键按 IS NULL，0 行或多行都令 Java getOne 失败。
func chatAppRows(d *sql.DB, app any) ([]map[string]any, error) {
	if app == nil {
		return queryRows(d, "SELECT * FROM chain_auth_chat_authorization WHERE app_identification IS NULL")
	}
	return queryRows(d, "SELECT * FROM chain_auth_chat_authorization WHERE app_identification=?", app)
}

// chatResultMap JSONObject.toJSON(ChatResult)：五字段全键，null 保留（实测 selPhoneMsg err 形状）。
func chatResultMap(m map[string]any) map[string]any {
	out := make(map[string]any, 5)
	for _, k := range []string{"openid", "session_key", "unionid", "errcode", "errmsg"} {
		out[k] = m[k]
	}
	return out
}

// chatResultJSON JSONObject.toJSONString(ChatResult)：fastjson 序列化省略 null 字段。
func chatResultJSON(m map[string]any) string {
	out := make(map[string]any, 5)
	for _, k := range []string{"openid", "session_key", "unionid", "errcode", "errmsg"} {
		if m[k] != nil {
			out[k] = m[k]
		}
	}
	b, _ := json.Marshal(out)
	return string(b)
}

// wxDecryptData 对照 AesUtils.decryptData：任何失败返回 ""，不抛异常。
// 填充校验按 BouncyCastle AES/CBC/PKCS7Padding（块 16，pad 1..16）。
func wxDecryptData(encryptedData, sessionKey, iv string) string {
	key, err1 := base64.StdEncoding.DecodeString(sessionKey)
	ivb, err2 := base64.StdEncoding.DecodeString(iv)
	data, err3 := base64.StdEncoding.DecodeString(encryptedData)
	if err1 != nil || err2 != nil || err3 != nil || len(key) == 0 || len(ivb) != aes.BlockSize || len(data) == 0 || len(data)%aes.BlockSize != 0 {
		return ""
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	out := make([]byte, len(data))
	cipher.NewCBCDecrypter(block, ivb).CryptBlocks(out, data)
	pad := int(out[len(out)-1])
	if pad < 1 || pad > aes.BlockSize {
		return ""
	}
	for _, b := range out[len(out)-pad:] {
		if int(b) != pad {
			return ""
		}
	}
	return string(out[:len(out)-pad])
}

// parseJSONObject fastjson parseObject：空串→null，非法 JSON→异常（调用方回 999）。
func parseJSONObject(text string) (map[string]any, error) {
	if strings.TrimSpace(text) == "" {
		return nil, nil
	}
	var m map[string]any
	dec := json.NewDecoder(strings.NewReader(text))
	dec.UseNumber()
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) chatLogin(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	rows, err := chatAppRows(d, "chain")
	if err != nil || len(rows) != 1 {
		platform.Unknown(w)
		return
	}
	res, err := s.jscode2session(r, str(rows[0], "appid"), str(rows[0], "secret"), jstr(m, "code"))
	if err != nil {
		platform.Unknown(w)
		return
	}
	obj, err := parseJSONObject(wxDecryptData(jstr(m, "encryptedData"), str(res, "session_key"), jstr(m, "iv")))
	if err != nil || obj == nil {
		platform.Unknown(w) // Java jsonObj.getString NPE
		return
	}
	users, err := queryRows(d, "SELECT * FROM chain_auth_user WHERE account_type!=4 AND handle_result=1 AND (phone=? OR account=?) ORDER BY level DESC", obj["phoneNumber"], obj["phoneNumber"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	switch len(users) {
	case 0:
		platform.Success(w, map[string]any{"accounts": []any{}})
	case 1:
		info, err := s.chatUserInfo(r, users[0])
		if err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, info) // assembleUserInfo 异常被捕获返回 null 时 success(null)
	default:
		phoneUsers, err := queryRows(d, "SELECT account, (SELECT simple_name FROM chain_basic_enterprise WHERE id=company_id) remark1 FROM chain_auth_user WHERE account_type!=4 AND handle_result=1 AND (phone=? OR account=?) ORDER BY level DESC", obj["phoneNumber"], obj["phoneNumber"])
		if err != nil {
			platform.Unknown(w)
			return
		}
		accounts := make([]string, 0, len(phoneUsers))
		for _, u := range phoneUsers {
			remark := str(u, "remark1")
			if remark == "" {
				remark = "平台"
			}
			accounts = append(accounts, str(u, "account")+"-"+remark)
		}
		platform.Success(w, map[string]any{"accounts": accounts})
	}
}

func (s *Service) chatSelPhoneMsg(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	const appid = "wx2c78c1ed3e6c6167" // Java 硬编码 appid（非密钥），secret 走私有运行配置
	res, err := s.jscode2session(r, appid, platform.LoadWxConfig().Secrets[appid], jstr(m, "code"))
	if err != nil {
		platform.Unknown(w)
		return
	}
	if res["session_key"] == nil {
		platform.Success(w, map[string]any{"err": chatResultMap(res)})
		return
	}
	obj, err := parseJSONObject(wxDecryptData(jstr(m, "encryptedData"), str(res, "session_key"), jstr(m, "iv")))
	if err != nil {
		platform.Unknown(w)
		return
	}
	if obj == nil {
		platform.Success(w, map[string]any{"err": chatResultMap(res)})
		return
	}
	phone := obj["phoneNumber"]
	users, err := queryRows(d, "SELECT max(company_id) companyId, `level`, count(1) FROM chain_auth_user WHERE phone=? AND (`level`=1 OR `level`=2) GROUP BY `level`", phone)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(users) == 0 {
		if _, err := platform.InsertNoFill(d, "chain_auth_user", insertVals(map[string]any{"phone": phone, "account": uuidV4(), "level": 1, "channel": "2"})); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, map[string]any{"etcName": "", "phone": phone, "identityType": 1, "wxFaceImg": ""})
		return
	}
	// 原 Java "1".equals(Integer level) 恒 false，etcName 永远空串；wxFaceImg 未被该查询选出恒 null。
	var identity any = users[0]["level"]
	if len(users) > 1 {
		identity = "1,2"
	}
	platform.Success(w, map[string]any{"etcName": "", "phone": phone, "identityType": identity, "wxFaceImg": nil})
}

// chatTouchIdentity 对照 updUserIdentity + setLastLoginTime/setUpdateTime；UA 缺失 Java NPE。
func chatTouchIdentity(r *http.Request, row map[string]any) error {
	values := r.Header.Values("User-Agent")
	if len(values) == 0 {
		return errors.New("User-Agent missing")
	}
	ua := strings.ToLower(values[0])
	now := time.Now().Format("2006-01-02 15:04:05")
	row["last_login_time"] = now
	row["update_time"] = now
	if strings.Contains(ua, "micromessenger") || strings.Contains(ua, "miniprogram") {
		row["is_mini_pro_user"] = true
	} else if strings.Contains(ua, "uni-app") {
		row["is_app_user"] = true
	} else if strings.Contains(ua, "windows") {
		row["is_web_user"] = true
	}
	return nil
}

func chatUpdateUserTx(tx *sql.Tx, row map[string]any) error {
	vals := insertVals(pick(row, authRestEnts["ChainAuthUser"].dbcols))
	delete(vals, "id")
	return platform.TxUpdateByIDNoFill(tx, "chain_auth_user", row["id"], vals)
}

// registScoreTx 对照 getRegistScore：newUser 或无注册积分日志时发放，返回发放额。
// 原服务方法经 TxAdviceInterceptor("*"→REQUIRED) 包裹，积分链与调用方同事务。
func registScoreTx(tx *sql.Tx, userID, phone any, newUser bool) (int, error) {
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM bus_score_log WHERE user_id=? AND score_get_type=1", userID).Scan(&count); err != nil {
		return 0, err
	}
	if !newUser && count > 0 {
		return 0, nil
	}
	var score, pool int
	if err := tx.QueryRow("SELECT regist_score,mo_get_sum_score FROM bus_score_rule LIMIT 1").Scan(&score, &pool); err != nil {
		return 0, err
	}
	award := score
	if pool < award {
		award = pool
	}
	if err := tx.QueryRow("SELECT count(*) FROM bus_user_score WHERE user_id=?", userID).Scan(&count); err != nil {
		return 0, err
	}
	if count == 0 {
		if _, err := tx.Exec("INSERT INTO bus_user_score(user_id,hjj_num,dj_hjj_num) VALUES (?,?,0)", userID, award); err != nil {
			return 0, err
		}
	} else if _, err := tx.Exec("UPDATE bus_user_score SET hjj_num=? WHERE user_id=?", award, userID); err != nil {
		return 0, err
	}
	if _, err := tx.Exec("INSERT INTO bus_score_log(user_id,user_phone,content,score,final_available_score,score_get_type) VALUES (?,?,'新用户注册',?,?,1)", userID, phone, award, award); err != nil {
		return 0, err
	}
	if _, err := tx.Exec("UPDATE bus_score_rule SET mo_get_sum_score=? WHERE id=1", pool-award); err != nil {
		return 0, err
	}
	return award, nil
}

func memberBenefits(d *sql.DB, id any) (map[string]any, error) {
	rows, err := queryRows(d, `SELECT is_show_etc_phone isShowEtcPhone,is_show_etc_name isShowEtcName,is_show_website isShowWebsite,is_see_data_briefing isSeeDataBriefing,is_see_card_visitor isSeeCardVisitor,is_see_warehouse_visitor isSeeWarehouseVisitor,is_see_testing_report isSeeTestingReport,is_see_price_trend isSeePriceTrend FROM bus_member_benefits b,chain_auth_user u WHERE b.member_type=u.member_type AND u.id=?`, id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return map[string]any{"isShowEtcPhone": 0, "isShowEtcName": 0}, nil
	}
	return rows[0], nil
}

// chatCardDefault 对照 searchHasDefault/add Feign 链。
func (s *Service) chatCardDefault(r *http.Request, userID any, phone string) error {
	has, err := s.callJSON(r, "/saas-warehouse/busBusinessCard/searchHasDefault", map[string]any{"userId": userID})
	if err != nil {
		return err
	}
	if has == "false" {
		_, err = s.callJSON(r, "/saas-warehouse/busBusinessCard/add", map[string]any{"userId": userID, "phone": phone})
	}
	return err
}

func (s *Service) chatLoginChat(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	rows, err := chatAppRows(d, m["appIdentification"])
	if err != nil || len(rows) != 1 {
		platform.Unknown(w)
		return
	}
	res, err := s.jscode2session(r, str(rows[0], "appid"), str(rows[0], "secret"), jstr(m, "code"))
	if err != nil {
		platform.Unknown(w)
		return
	}
	if res["session_key"] == nil {
		platform.Fail(w, 601, chatResultJSON(res))
		return
	}
	openid := str(res, "openid")
	openUsers, err := queryRows(d, "SELECT * FROM chain_auth_user WHERE openid=? AND last_login_time IS NOT NULL ORDER BY last_login_time DESC", openid)
	if err != nil {
		platform.Unknown(w)
		return
	}
	// 原服务方法经 TxAdviceInterceptor 全量事务：后段依赖失败时前段落库整体回滚。
	tx, err := d.BeginTx(r.Context(), nil)
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer tx.Rollback()
	newUser := false
	var row map[string]any
	phone := ""
	if m["searchUser"] != nil {
		if len(openUsers) > 0 {
			if openUsers[0]["phone"] == nil {
				platform.FailData(w, 603, "", authRestEnts["ChainAuthUser"].shapeOne(openUsers[0]))
				return
			}
			row = openUsers[0]
			if err := chatTouchIdentity(r, row); err != nil {
				platform.Unknown(w)
				return
			}
			if err := chatUpdateUserTx(tx, row); err != nil {
				platform.Unknown(w)
				return
			}
		} else {
			ouser := map[string]any{"account": uuidV4(), "level": 1, "channel": "2", "is_agree_call": false, "openid": openid}
			if err := chatTouchIdentity(r, ouser); err != nil {
				platform.Unknown(w)
				return
			}
			id, err := platform.TxInsertNoFill(tx, "chain_auth_user", insertVals(ouser))
			if err != nil {
				platform.Unknown(w)
				return
			}
			ouser["id"] = id
			// 605 是正常返回（无异常），Java 事务提交，新用户已落库。
			if err := tx.Commit(); err != nil {
				platform.Unknown(w)
				return
			}
			platform.FailData(w, 605, "", authRestEnts["ChainAuthUser"].shapeOne(ouser))
			return
		}
		phone = str(row, "phone")
	} else {
		obj, err := parseJSONObject(wxDecryptData(jstr(m, "encryptedData"), str(res, "session_key"), jstr(m, "iv")))
		if err != nil {
			platform.Unknown(w)
			return
		}
		if obj == nil {
			platform.Fail(w, 602, "")
			return
		}
		phone = str(obj, "phoneNumber")
		phoneUsers, err := queryRows(d, "SELECT * FROM chain_auth_user WHERE (phone=? OR account=?) AND channel=2 ORDER BY level DESC", obj["phoneNumber"], obj["phoneNumber"])
		if err != nil {
			platform.Unknown(w)
			return
		}
		switch {
		case len(phoneUsers) > 0:
			row = phoneUsers[0]
			row["openid"] = openid
		case len(openUsers) > 0 && openUsers[0]["phone"] == nil:
			row = openUsers[0]
			if len(phone) < 7 {
				platform.Unknown(w)
				return
			}
			row["nickname"] = "匿名用户" + phone[7:]
			row["phone"] = phone
			newUser = true
		default:
			if len(phone) < 7 {
				platform.Unknown(w)
				return
			}
			row = map[string]any{"phone": phone, "account": uuidV4(), "level": 1, "channel": "2", "nickname": "匿名用户" + phone[7:], "is_agree_call": false, "openid": openid}
			newUser = true
		}
		if err := chatTouchIdentity(r, row); err != nil {
			platform.Unknown(w)
			return
		}
		if newUser && row["id"] == nil {
			id, err := platform.TxInsertNoFill(tx, "chain_auth_user", insertVals(row))
			if err != nil {
				platform.Unknown(w)
				return
			}
			row["id"] = id
		} else if err := chatUpdateUserTx(tx, row); err != nil {
			platform.Unknown(w)
			return
		}
	}
	score, err := registScoreTx(tx, row["id"], row["phone"], newUser)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if err := s.chatCardDefault(r, row["id"], phone); err != nil {
		platform.Unknown(w)
		return
	}
	info, err := s.chatUserInfo(r, row)
	if err != nil || info == nil {
		platform.Unknown(w) // assembleUserInfo null 后 rMap.put NPE
		return
	}
	benefits, err := memberBenefits(d, row["id"]) // loginChat 传 user.id（loginChatPhone 传 memberType，二者不同）
	if err != nil {
		platform.Unknown(w)
		return
	}
	care := false
	if str(row, "official_openid") != "" {
		care, err = s.chatOfficialCare(r, str(row, "official_openid"))
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	var expiration any
	if value, ok := row["expiration_time"].(string); ok && len(value) >= 10 {
		expiration = value[:10]
	}
	shaped := authRestEnts["ChainAuthUser"].shapeOne(row)
	info["newUser"] = newUser
	info["registScore"] = score
	info["openid"] = openid
	info["officialOpenid"] = row["official_openid"]
	info["isCareOfficial"] = care
	info["MemberBenefits"] = benefits
	info["isPromptCare"] = shaped["isPromptCare"]
	info["expirationTimeStr"] = expiration
	if err := tx.Commit(); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, info)
}

func (s *Service) chatLoginChatApp(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if m["phone"] != nil {
		s.chatHasPhoneLogin(w, r, m)
	} else {
		s.chatNoPhoneLogin(w, r, m)
	}
}

func (s *Service) chatHasPhoneLogin(w http.ResponseWriter, r *http.Request, m map[string]any) {
	d := s.pick(r)
	phone := str(m, "phone")
	rows, err := queryRows(d, "SELECT * FROM chain_auth_user WHERE phone=? AND channel=2", m["phone"])
	if err != nil || len(rows) > 1 {
		platform.Unknown(w)
		return
	}
	tx, err := d.BeginTx(r.Context(), nil)
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer tx.Rollback()
	var row map[string]any
	if len(rows) == 0 {
		if len(phone) < 7 {
			platform.Unknown(w)
			return
		}
		row = map[string]any{"phone": phone, "account": uuidV4(), "level": 1, "channel": "2", "openid": m["openid"], "app_client_id": m["pushClientid"], "nickname": "匿名用户" + phone[7:]}
		id, err := platform.TxInsert(tx, "chain_auth_user", insertVals(row))
		if err != nil {
			platform.Unknown(w)
			return
		}
		row["id"] = id
	} else {
		row = rows[0]
	}
	if err := s.chatCardDefault(r, row["id"], phone); err != nil {
		platform.Unknown(w)
		return
	}
	if len(rows) > 0 {
		if m["openid"] != nil {
			row["openid"] = m["openid"]
		}
		if m["pushClientid"] != nil {
			row["app_client_id"] = m["pushClientid"]
		}
		if err := chatUpdateUserTx(tx, row); err != nil {
			platform.Unknown(w)
			return
		}
	}
	info, err := s.chatUserInfo(r, row)
	if err != nil || info == nil {
		platform.Unknown(w) // 新用户 isAgreeCall null → assembleUserInfo null → rMap.put NPE → 事务回滚
		return
	}
	benefits, err := memberBenefits(d, row["member_type"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	info["MemberBenefits"] = benefits
	if err := tx.Commit(); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, info)
}

func (s *Service) chatNoPhoneLogin(w http.ResponseWriter, r *http.Request, m map[string]any) {
	// Java 整体 try/catch(Exception) → null → HTTP 200 空体。
	empty := func() { w.WriteHeader(http.StatusOK) }
	d := s.pick(r)
	const appid = "wx0ff59ec3274e0a53"
	secret := platform.LoadWxConfig().Secrets[appid]
	res, err := s.wxJSON(r, "GET", "https://api.weixin.qq.com/sns/oauth2/access_token?appid="+appid+"&secret="+secret+"&code="+jstr(m, "code")+"&grant_type=authorization_code", nil)
	if err != nil {
		empty()
		return
	}
	var rows []map[string]any
	if res["openid"] == nil {
		rows, err = queryRows(d, "SELECT * FROM chain_auth_user WHERE openid IS NULL")
	} else {
		rows, err = queryRows(d, "SELECT * FROM chain_auth_user WHERE openid=?", res["openid"])
	}
	if err != nil || len(rows) > 1 {
		empty()
		return
	}
	if len(rows) == 0 || rows[0]["phone"] == nil {
		platform.FailData(w, 602, "请先绑定手机号", res)
		return
	}
	row := rows[0]
	if err := s.chatCardDefault(r, row["id"], str(row, "phone")); err != nil {
		empty()
		return
	}
	info, err := s.chatUserInfo(r, row)
	if err != nil {
		empty()
		return
	}
	platform.Success(w, info) // info 可为 null → success(null)
}

func (s *Service) chatWxConfig(w http.ResponseWriter, r *http.Request) {
	const appid = "wxa9a9e7a161a0e989"
	const requestURL = "https://zsc.coalcloudtech.com/wxindex"
	secret := platform.LoadWxConfig().Secrets[appid]
	var accessToken any = ""
	var jsapiTicket any = ""
	// httpRequest 任何异常→null：token 失败则 ticket 跳过，签名串保留空串。
	if tokenRes, err := s.wxJSON(r, "GET", "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="+appid+"&secret="+secret, nil); err == nil {
		accessToken = tokenRes["access_token"]
		if ticketRes, err := s.wxJSON(r, "GET", "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token="+platform.JnullAny(accessToken)+"&type=jsapi", nil); err == nil {
			jsapiTicket = ticketRes["ticket"]
		}
	}
	nonce := uuidV4()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	sign := "jsapi_ticket=" + platform.JnullAny(jsapiTicket) + "&noncestr=" + nonce + "&timestamp=" + timestamp + "&url=" + requestURL
	sum := sha1.Sum([]byte(sign))
	platform.Success(w, map[string]any{"appId": appid, "timestamp": timestamp, "nonceStr": nonce, "signature": hex.EncodeToString(sum[:])})
}

func (s *Service) wechatAuth(w http.ResponseWriter, r *http.Request) {
	echo := "没有获取到"
	if q := r.URL.Query(); q.Has("echostr") {
		echo = q.Get("echostr")
	}
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	_, _ = w.Write([]byte(echo))
}

// javaURLEncode 对照 URLEncoder.encode(s,"UTF-8")：保留 . - * _，空格→+。
func javaURLEncode(s string) string {
	var b strings.Builder
	for _, c := range []byte(s) {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '.', c == '-', c == '*', c == '_':
			b.WriteByte(c)
		case c == ' ':
			b.WriteByte('+')
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

func (s *Service) wechatAuthorize(w http.ResponseWriter, r *http.Request) {
	returnURL := "https://qkongtao.cn/"
	if q := r.URL.Query(); q.Has("returnUrl") {
		returnURL = q.Get("returnUrl")
	}
	sig, err := s.Auth.AnonSignature()
	if err != nil {
		platform.Unknown(w)
		return
	}
	redirectURI := "https://test.coalcloudtech.com/saas/saas-warehouse/wechat/userInfo?Signature=" + sig + "&AuthorizationId=wxda0b1c468b83bc99"
	cfg := platform.LoadWxConfig()
	u := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + platform.Jnull(cfg.PushOfficialAppid) +
		"&redirect_uri=" + javaURLEncode(redirectURI) + "&response_type=code&scope=snsapi_userinfo&state=" + javaURLEncode(returnURL) + "#wechat_redirect"
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	_, _ = w.Write([]byte("redirect:" + u)) // @RestController 下 "redirect:" 前缀是字面响应体，不是跳转
}

func (s *Service) wechatUserInfo(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if m["code"] == nil || m["state"] == nil {
		platform.Unknown(w) // maps.get(...).toString() NPE
		return
	}
	d := s.pick(r)
	cfg := platform.LoadWxConfig()
	tokenRes, err := s.wxJSON(r, "GET", "https://api.weixin.qq.com/sns/oauth2/access_token?appid="+platform.Jnull(cfg.PushOfficialAppid)+"&secret="+platform.Jnull(cfg.PushOfficialSecret)+"&code="+str(m, "code")+"&grant_type=authorization_code", nil)
	openID := ""
	var accessToken any
	if err == nil && tokenRes["errcode"] == nil {
		openID = str(tokenRes, "openid")
		accessToken = tokenRes["access_token"]
	}
	if openID != "" {
		users, err := queryRows(d, "SELECT * FROM chain_auth_user WHERE official_openid=? ORDER BY last_login_time DESC", openID)
		if err != nil {
			platform.Unknown(w)
			return
		}
		if len(users) == 0 {
			if _, err := platform.InsertNoFill(d, "chain_auth_user", map[string]any{"official_openid": openID}); err != nil {
				platform.Unknown(w)
				return
			}
		}
	}
	// oauth2getUserInfo 抛 WxErrorException 未被捕获 → 999。
	infoRes, err := s.wxJSON(r, "GET", "https://api.weixin.qq.com/sns/userinfo?access_token="+platform.JnullAny(accessToken)+"&openid="+openID+"&lang=zh_CN", nil)
	if err != nil || infoRes["errcode"] != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, openID)
}

// esignCreatePerAccount 对照 SignHelper.createPerAccount：env appId/appKey/host 缺失或响应无 data.accountId → 失败（Java NPE→999）。
func (s *Service) esignCreatePerAccount(r *http.Request, userInfo any) (string, error) {
	if userInfo == nil {
		return "", errors.New("userInfo null")
	}
	host := os.Getenv("host")
	if host == "" {
		return "", errors.New("esign host 未配置")
	}
	body, err := json.Marshal(userInfo)
	if err != nil {
		return "", err
	}
	const api = "/v1/accounts/createByThirdPartyUserId"
	sum := md5.Sum(body)
	contentMD5 := base64.StdEncoding.EncodeToString(sum[:])
	plaintext := "POST\n*/*\n" + contentMD5 + "\napplication/json; charset=UTF-8\n\n" + api
	mac := hmac.New(sha256.New, []byte(os.Getenv("appKey")))
	mac.Write([]byte(plaintext))
	req, err := http.NewRequestWithContext(r.Context(), "POST", host+api, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Tsign-Open-App-Id", os.Getenv("appId"))
	req.Header.Set("X-Tsign-Open-Auth-Mode", "Signature")
	req.Header.Set("X-Tsign-Open-Ca-Timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Tsign-Open-Ca-Signature", base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	req.Header.Set("Content-MD5", contentMD5)
	client := s.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result map[string]any
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	if err := dec.Decode(&result); err != nil {
		return "", err
	}
	data, _ := result["data"].(map[string]any)
	if data == nil || data["accountId"] == nil {
		return "", errors.New("esign accountId missing")
	}
	return str(data, "accountId"), nil
}

func (s *Service) chatReg(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	e := authRestEnts["ChainAuthUser"]
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	if m["code"] == nil {
		platform.Unknown(w) // getCode().equals NPE
		return
	}
	if msg := s.Auth.PhoneCodeError(jstr(m, "phone"), str(m, "code")); msg != "" {
		platform.Fail(w, 3004, "验证码不匹配")
		return
	}
	users, err := queryRows(d, "SELECT * FROM chain_auth_user WHERE account_type!=4 AND handle_result=1 AND (phone=? OR account=?) ORDER BY level DESC", m["phone"], m["phone"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(users) == 0 {
		perID, err := s.esignCreatePerAccount(r, m["userInfo"])
		if err != nil {
			platform.Unknown(w)
			return
		}
		vals := insertVals(pick(m, e.dbcols))
		vals["eqb_per_id"] = perID
		if _, err := platform.InsertNoFill(d, "chain_auth_user", vals); err != nil {
			platform.Unknown(w)
			return
		}
	} else {
		perID := str(users[0], "eqb_per_id")
		if perID == "" {
			perID, err := s.esignCreatePerAccount(r, m["userInfo"])
			if err != nil {
				platform.Unknown(w)
				return
			}
			if _, err := d.Exec("UPDATE chain_auth_user SET eqb_per_id=? WHERE phone=?", perID, m["phone"]); err != nil {
				platform.Unknown(w)
				return
			}
		} else {
			if m["regType"] == nil {
				platform.Unknown(w) // 1 == getRegType() 拆箱 NPE
				return
			}
			if num(m["regType"]) == 1 {
				platform.Fail(w, 2005, "用户信息已存在")
				return
			}
		}
	}
	// 原 Java "2".equals(Integer regType) 恒 false，企业注册分支为死代码，不复刻。
	platform.Success(w, nil)
}

func (s *Service) chatUpdMerExpire(w http.ResponseWriter, r *http.Request) {
	rows, err := queryRows(s.pick(r), "SELECT * FROM chain_auth_user WHERE member_type=1")
	if err != nil {
		platform.Unknown(w)
		return
	}
	for _, row := range rows {
		if row["expiration_time"] == nil {
			platform.Unknown(w) // Java expirationTime.atZone NPE
			return
		}
		if err := auth.ScheduleMemberExpiryText(r.Context(), str(row, "id"), str(row, "expiration_time")); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "更新成功")
}

func (s *Service) codeRefundMoney(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	u := currentUser(r)
	if u.UserID == 0 {
		platform.Unknown(w) // jwtInfo null NPE
		return
	}
	if platform.LoadWxConfig().SendMsg {
		user, err := platform.GetByID(s.pick(r), "chain_auth_user", u.UserID)
		if err != nil || user == nil || m["templateId"] == nil {
			platform.Unknown(w)
			return
		}
		code := s.Auth.IssuePhoneCode(fmt.Sprintf("%d", u.UserID))
		platform.SendLegacySMS(r.Context(), str(user, "phone"), []string{code}, str(m, "templateId"), "")
		platform.Success(w, "")
		return
	}
	platform.Success(w, s.Auth.IssuePhoneCode(fmt.Sprintf("%d", u.UserID)))
}
