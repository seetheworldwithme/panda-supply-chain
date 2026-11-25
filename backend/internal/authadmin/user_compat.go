package authadmin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) registerUserCompat(reg func(string, http.HandlerFunc)) {
	reg("/chainAuthUser/list", s.userAllEqList)
	reg("/chainAuthUser/getById", s.userGetDetail)
	reg("/chainAuthUser/ptlist", func(w http.ResponseWriter, r *http.Request) {
		rows, err := queryRows(s.DB.Chain, "SELECT id,account,username FROM chain_auth_user WHERE level>1")
		if err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, shapeUser(rows))
	})
	reg("/chainAuthUser/updHeadImg", func(w http.ResponseWriter, r *http.Request) {
		m := bodyMap(r)
		// 只能改自己的头像:原先按 body 的 phone 直写
		if _, err := s.DB.Chain.Exec("UPDATE chain_auth_user SET wx_face_img=? WHERE phone=? AND id=?", m["wxFaceImg"], m["phone"], auth.From(r.Context()).UserID); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, nil)
	})
	reg("/chainAuthUser/listCarAccount", s.userCarList)
	reg("/chainAuthUser/addCarAccount", func(w http.ResponseWriter, r *http.Request) { s.userCarWrite(w, r, true) })
	reg("/chainAuthUser/updCarAccount", func(w http.ResponseWriter, r *http.Request) { s.userCarWrite(w, r, false) })
	reg("/chainAuthUser/resetPwd", s.userResetPwd)
	reg("/chainAuthUser/updPwd", s.userSetPwd)
	reg("/chainAuthUser/selByIds", s.usersByIDs)
	reg("/chainAuthUser/isPerfect", s.userIsPerfect)
	reg("/chainAuthUser/updBatch", s.userUpdateBatch)
	reg("/chainAuthUser/joinCpy", s.userJoinCompany)
	reg("/chainAuthUser/valIdentity", s.userValidateIdentity)
	reg("/chainAuthUser/valCode", s.userValidateCode)
	reg("/chainAuthUser/setAccount", s.userSetAccount)
	reg("/chainAuthUserList/deleteById", s.userDeleteById)
	reg("/chainAuthUserList/selEtcRoles", s.selEtcRoles)
	reg("/chainAuthUserList/selPtRoles", s.selPtRoles)
	reg("/chainAuthAuthorityRole/selAuthy", s.authorityRoleForUser)
}

func (s *Service) userValidateIdentity(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cols := append(append([]string{}, userWriteCols...), "create_time", "update_time")
	cond, ok := allEqCond(m, cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.DB.Chain, "chain_auth_user", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(rows) == 0 {
		platform.Success(w, "身份验证失败")
		return
	}
	if s.Auth == nil {
		platform.Unknown(w)
		return
	}
	phone := str(m, "phone")
	if m["phone"] == nil {
		phone = "null"
	}
	s.Auth.IssuePhoneCode(phone)
	platform.Success(w, nil)
}

func (s *Service) userValidateCode(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if m["code"] == nil {
		platform.Unknown(w)
		return
	}
	if s.Auth == nil || !s.Auth.VerifyPhoneCode(str(m, "phone"), str(m, "code")) {
		platform.FailMsg(w, "验证失败")
		return
	}
	s.Auth.MarkPhoneVerified(str(m, "phone")) // 忘记密码第二步(updPwd)只允许改这个手机号下的账号
	where, args := getSqlEq(m, userExcelCols, nil, map[string]bool{"singleDatabase": true, "isAgreeCall": true, "memberType": false})
	rows, err := queryRows(s.DB.Chain, "SELECT id,account FROM chain_auth_user "+where+" AND handle_result=1 ORDER BY create_time DESC", args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, rows)
}

func (s *Service) userSetAccount(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if m["code"] == nil {
		platform.Unknown(w)
		return
	}
	if s.Auth == nil || !s.Auth.VerifyPhoneCode(str(m, "phone"), str(m, "code")) {
		platform.Fail(w, 999, "验证码错误或已过期")
		return
	}
	vals := pick(m, userWriteCols...)
	delete(vals, "id")
	// 注册接口未登录可调:不能靠 body 带出 accountType=5 / level=2 的"管理员"。
	scrubUserVals(&auth.CurrentUser{}, vals, false)
	// Java chainAuthUserService.save(user) 原样入库(明文密码),1:1 复刻,不做加密。
	id, err := platform.InsertNoFill(s.DB.Chain, "chain_auth_user", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, map[string]any{"userId": id})
}

func (s *Service) userAllEqList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cols := append(append([]string{}, userWriteCols...), "create_time", "update_time")
	cond, ok := allEqCond(m, cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.DB.Chain, "chain_auth_user", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeUser(rows))
}

func (s *Service) userGetDetail(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	row, err := platform.GetByID(s.DB.Chain, "chain_auth_user", m["id"])
	if err != nil || row == nil {
		platform.Unknown(w)
		return
	}
	terminals, err := queryRows(s.DB.Chain, "SELECT id FROM risk_terminal WHERE user_id=?", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	ids := []any{}
	for _, terminal := range terminals {
		ids = append(ids, terminal["id"])
	}
	if row["nickname"] == nil {
		phone, ok := row["phone"].(string)
		if !ok || len(phone) < 7 {
			platform.Unknown(w)
			return
		}
		nickname := strconv.FormatInt(time.Now().UnixMilli(), 10) + phone[7:]
		if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_auth_user", m["id"], map[string]any{"nickname": nickname}); err != nil {
			platform.Unknown(w)
			return
		}
		row["nickname"] = nickname
	}
	out := shapeUserOne(row)
	out["ids"] = ids
	platform.Success(w, out)
}

func (s *Service) userCarList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	where := " WHERE 1=1"
	args := []any{}
	fields := map[string]bool{}
	for _, f := range userExcelCols {
		fields[f] = true
	}
	for _, k := range platform.SortedKeys(m) {
		if k == "applys" || !fields[k] {
			continue
		}
		if m[k] == nil {
			platform.Unknown(w)
			return
		} // Java getMaps(Map) dereferences values
		col := platform.Underscore(map[string]any{k: nil})
		for c := range col {
			where += " AND " + c + " LIKE ?"
		}
		if k == "companyId" {
			args = append(args, m[k])
		} else {
			args = append(args, "%"+str(m, k)+"%")
		}
	}
	if m["applys"] != nil {
		where += " AND applys LIKE ?"
		args = append(args, "%"+str(m, "applys")+"%")
	}
	var search any
	if m["searchStr"] != nil {
		search = "%" + str(m, "searchStr") + "%"
	}
	args = append(args, search)
	q := "SELECT id,phone,account,username,(SELECT GROUP_CONCAT(r.name) FROM chain_auth_use_role ur,chain_auth_role r WHERE ur.role_id=r.id AND ur.user_id=u.id GROUP BY ur.user_id) roleNames,(SELECT enterprise_name FROM chain_basic_enterprise WHERE id=company_id) enterpriseName FROM chain_auth_user u" + where + " AND CONCAT(IFNULL(phone,''),IFNULL(username,''),IFNULL(account,'')) LIKE ? ORDER BY create_time DESC"
	rows, err := queryRows(s.DB.Chain, q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeUser(rows))
}

// Vehicle account and terminal assignments share the Java service transaction.
func (s *Service) userCarWrite(w http.ResponseWriter, r *http.Request, create bool) {
	m := bodyMap(r)
	if !s.requirePlatform(w, r) { // 车辆账号与终端分配是平台操作
		return
	}
	vals := pick(m, userWriteCols...)
	delete(vals, "id")
	id := m["id"]
	if create {
		password, ok := m["password"].(string)
		if !ok {
			platform.Unknown(w)
			return
		}
		vals["password"] = auth.HashPassword(password)
		if vals["password"] == "" { // bcrypt 拒绝超过 72 字节的口令:不能把空口令落库
			platform.Unknown(w)
			return
		}
		vals["applys"] = "2"
		vals["handle_result"] = "1"
		vals["level"] = 1
	}
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		if create {
			n, err := platform.TxInsertNoFill(tx, "chain_auth_user", vals)
			if err != nil {
				return err
			}
			id = n
		} else {
			if err := platform.TxUpdateByIDNoFill(tx, "chain_auth_user", id, vals); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE risk_terminal SET user_id=NULL,account=NULL WHERE user_id=?", id); err != nil {
				return err
			}
		}
		if ids, ok := m["ids"].([]any); ok {
			for _, terminalID := range ids {
				if _, err := tx.Exec("UPDATE risk_terminal SET user_id=?,account=? WHERE id=?", id, m["account"], terminalID); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, id)
}

func (s *Service) userSetPwd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	// 原先请求体里写谁的 id 就改谁的密码。现在只允许:改自己、平台管理员 / 企业管理员改其管理的账号、
	// 忘记密码流程(未登录)且目标账号的手机号刚在 valCode 通过短信码校验。
	cu := auth.From(r.Context())
	verifiedOwner := func() bool {
		if s.Auth == nil {
			return false
		}
		target, err := platform.GetByID(s.DB.Chain, "chain_auth_user", m["id"])
		return err == nil && target != nil && s.Auth.PhoneVerifiedRecently(str(target, "phone"))
	}
	if !(cu.UserID != 0 && str(m, "id") == strconv.FormatInt(cu.UserID, 10)) && !s.canManageID(r, m["id"]) && !verifiedOwner() {
		platform.Fail(w, 602, "请先完成身份验证")
		return
	}
	password, ok := m["newPassword"].(string)
	if !ok {
		platform.Unknown(w)
		return
	}
	vals := pick(m, userWriteCols...)
	vals["password"] = auth.HashPassword(password)
	if vals["password"] == "" { // bcrypt 拒绝超过 72 字节的口令并返回空串:不能把空口令落库
		platform.Unknown(w)
		return
	}
	scrubUserVals(auth.From(r.Context()), vals, s.isEnterpriseAdmin(r, auth.From(r.Context())))
	if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_auth_user", m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "修改成功")
}
func (s *Service) userResetPwd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.canManageID(r, m["id"]) { // 原先任何登录用户可把任何账号(含平台管理员)的密码重置为 123456
		deny(w)
		return
	}
	vals := pick(m, userWriteCols...)
	vals["password"] = auth.HashPassword("123456") // original Java resetPwd contract
	scrubUserVals(auth.From(r.Context()), vals, s.isEnterpriseAdmin(r, auth.From(r.Context())))
	if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_auth_user", m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "重置成功")
}
func (s *Service) usersByIDs(w http.ResponseWriter, r *http.Request) {
	var ids []any
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil || len(ids) == 0 {
		platform.Unknown(w)
		return
	}
	marks := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	rows, err := queryRows(s.DB.Chain, "SELECT * FROM chain_auth_user WHERE id IN ("+marks+")", ids...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, shapeUser(rows))
}
func (s *Service) userIsPerfect(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	row, err := platform.GetByID(s.DB.Chain, "chain_auth_user", m["id"])
	if err != nil || row == nil {
		platform.Unknown(w)
		return
	}
	perfect := true
	for _, key := range []string{"ware_identity", "nickname", "etc_name", "username", "duty"} {
		if row[key] == nil {
			perfect = false
		}
	}
	platform.Success(w, perfect)
}
func (s *Service) userUpdateBatch(w http.ResponseWriter, r *http.Request) {
	var users []map[string]any
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		platform.Unknown(w)
		return
	}
	if len(users) > 0 {
		for _, m := range users {
			if !s.canManageID(r, m["id"]) {
				deny(w)
				return
			}
		}
		cu, ent := auth.From(r.Context()), false
		if cu.UserID != 0 {
			ent = s.isEnterpriseAdmin(r, cu)
		}
		err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
			for _, m := range users {
				vals := pick(m, userWriteCols...)
				delete(vals, "id")
				scrubUserVals(cu, vals, ent)
				if err := platform.TxUpdateByIDNoFill(tx, "chain_auth_user", m["id"], vals); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "删除成功")
}

func (s *Service) userJoinCompany(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	var count int
	q := "SELECT count(*) FROM chain_auth_user WHERE account=?"
	args := []any{m["account"]}
	if m["account"] == nil {
		q = "SELECT count(*) FROM chain_auth_user"
		args = nil
	}
	if err := s.DB.Chain.QueryRow(q, args...).Scan(&count); err != nil || count > 1 {
		platform.Unknown(w)
		return
	}
	if count == 1 {
		platform.Fail(w, 999, "账号信息已存在，请重新设置")
		return
	}
	code := str(m, "code")
	if _, ok := inviteCodeGet(code); !ok {
		platform.Fail(w, 999, "邀请码不存在或已过期")
		return
	}
	parts := strings.Split(code, "c")
	if len(parts) != 2 {
		platform.Unknown(w)
		return
	}
	cid, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || m["password"] == nil {
		platform.Unknown(w)
		return
	}
	// 只能认领"刚注册、尚未归属"的账号:原先 body 的 id 任意,任何登录用户拿着邀请码
	// 就能改掉别人(含平台管理员)的密码、账号和归属。
	if target, err := platform.GetByID(s.DB.Chain, "chain_auth_user", m["id"]); err != nil || target == nil ||
		num(target["company_id"]) != 0 || num(target["level"]) != 0 || num(target["account_type"]) != 1 || str(target, "handle_result") != "" {
		platform.Fail(w, 602, "无权限操作")
		return
	}
	vals := pick(m, userWriteCols...)
	for _, c := range []string{"account_type", "level", "applys", "member_type", "expiration_time", "activate_type", "package_id", "broker_start_date", "servicer_url", "single_database", "single_database_name"} {
		delete(vals, c) // 加入企业只写账号、密码和归属,带不出平台属性或会员
	}
	vals["company_id"] = cid
	vals["handle_result"] = "0"
	vals["password"] = auth.HashPassword(str(m, "password"))
	if vals["password"] == "" { // bcrypt 拒绝超过 72 字节的口令并返回空串:不能把空口令落库
		platform.Unknown(w)
		return
	}
	if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_auth_user", m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}
func (s *Service) authorityRoleForUser(w http.ResponseWriter, r *http.Request) {
	rows, err := queryRows(s.DB.Chain, "SELECT ar.menu,GROUP_CONCAT(ar.authoritys) FROM chain_auth_use_role ur,chain_auth_authority_role ar WHERE ur.role_id=ar.role_id AND ur.user_id=? GROUP BY ar.menu", auth.From(r.Context()).UserID)
	if err != nil || len(rows) > 1 {
		platform.Unknown(w)
		return
	}
	if len(rows) == 0 {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, rows[0])
}
