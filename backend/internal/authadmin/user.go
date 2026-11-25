// user.go 用户管理(ChainAuthUserController + ChainAuthUserListController)。
// 页面:UserList(企业授权用户三 tab)/PtUserList(平台用户)/InvitationCode(邀请码)。
// memberType 0→1 的 redis 到期缓存与短信、isAgreeCall 经纪人短信：短信异常被 Java 吞掉、
// Go 会员状态读库不依赖缓存，NPE 边界已按原顺序复刻；redis 头像缓存无契约影响。
package authadmin

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// userBoolCols chain_auth_user 的 tinyint(1) 布尔列
var userBoolCols = []string{"single_database", "msg_alarm_flag", "is_agree_call", "is_perfect",
	"is_mini_pro_user", "is_app_user", "is_web_user", "is_prompt_care"}

// userExtra 已并入 userFillKeys/userCarriers(见 shapeUser)

// userWriteCols 用户表允许写入的列(add/updateById 白名单)
var userWriteCols = []string{"id", "phone", "account", "password", "account_type", "level",
	"username", "channel", "company_id", "single_database", "single_database_name",
	"handle_result", "xzqh_code", "xzqh_name", "eqb_per_id", "wx_face_img", "openid",
	"official_openid", "app_client_id", "applys", "msg_alarm_flag", "alarm_phone",
	"ware_identity", "member_type", "expiration_time", "activate_type", "package_id",
	"nickname", "etc_name", "duty", "broker_start_date", "is_agree_call", "servicer_url",
	"is_perfect", "follow_person", "follow_status", "last_login_time", "is_mini_pro_user",
	"is_app_user", "is_web_user", "is_prompt_care",
	"remark1", "remark2", "remark3", "remark4", "remark5"}

// useRoleExtra chain_auth_use_role 实体载体键
var useRoleExtra = []string{"searchStr", "page", "ids"}

// userFillKeys 实体全键(表列驼峰 + 真载体),shapeUser 对缺键补 null——
// 对照 MP 实体 Jackson 全字段输出(supList 只查 5 列时其余键为 null)
var userFillKeys = append(append([]string{}, userExcelCols...),
	"accountType", "level", "channel", "handleResult")
var userCarriers = []string{"searchStr", "page", "ids", "newPassword", "code", "effSeconds",
	"userRoles", "roleNames", "userId", "userInfo", "orgInfo", "regType", "enterpriseName",
	"borkerYears", "startTime", "endTime", "isAcOperation", "hasPhone", "etcInfo", "isCtiExDate",
	"expirationTimeStr"}

func shapeUser(rows []map[string]any) []map[string]any {
	// 注:入参行可能已是 camel(paged 路径 queryCam)或 snake(GetByID 路径),统一在 Camelize 后处理。
	// tinyint(1) 布尔列由 driver(tinyInt1isBit)天然返回 bool,无需转换
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		// tinyint(1) 布尔列(int64 → bool,实体 Boolean 字段)
		for _, c := range []string{"singleDatabase", "msgAlarmFlag", "isAgreeCall", "isPerfect",
			"isMiniProUser", "isAppUser", "isWebUser", "isPromptCare"} {
			if v, ok := r[c]; ok && v != nil {
				r[c] = num(v) != 0
			}
		}
		// account_type 表列 int,实体 String 字段 → 字符串(Java MyBatis 映射行为)
		if v, ok := r["accountType"]; ok && v != nil {
			r["accountType"] = str(r, "accountType")
		}
		// updateTime/lastLoginTime 实体无 @JsonFormat → Jackson ISO 带 T(与 createTime 等格式不同,刻意保留差异)
		for _, c := range []string{"updateTime", "lastLoginTime"} {
			if v, ok := r[c]; ok && v != nil {
				if s := str(r, c); len(s) == 19 {
					r[c] = s[:10] + "T" + s[11:]
				}
			}
		}
		// 实体 getter 派生键:expirationTime 非空时取日期部分
		if et, ok := r["expirationTime"]; ok && et != nil {
			r["expirationTimeStr"] = str(r, "expirationTime")[:10]
		}
		r["password"] = nil // 用户列表/详情对任意登录用户开放:不回传口令散列(保留键,响应形状不变)
		fillNulls(r, userFillKeys)
		fillNulls(r, userCarriers)
	}
	return out
}

func shapeUserOne(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	return shapeUser([]map[string]any{m})[0]
}

// getSqlEq 对照 StringHelper.getSql(t, unstr):实体 ExcelColumn 字段非空、不在排除表、
// 非"0" → `and <列> like ?`(值由 getMaps 语义:数字原值、字符串 %v%)。
// fields 为 camelCase 字段名(对照 @ExcelColumn 集);返回条件串与参数。
// wrap 为 Java 非 Integer/LocalDate 类型字段集(Boolean/Long/LocalDateTime 等
// getMaps 一律包 %…%,仅布尔/长整型键会出现非字符串取值),用户组实体无此类键传 nil
func getSqlEq(m map[string]any, fields []string, unstr map[string]bool, wrap map[string]bool) (string, []any) {
	var sb strings.Builder
	var args []any
	sb.WriteString("where 1=1")
	for _, f := range fields {
		col := platform.Underscore(map[string]any{f: nil})
		c := ""
		for k := range col {
			c = k
		}
		v, ok := m[f]
		if !ok || v == nil || unstr[f] {
			continue
		}
		if s, isStr := v.(string); isStr && s == "0" {
			continue
		}
		sb.WriteString(" and " + c + " like ?")
		if _, isStr := v.(string); isStr || wrap[f] {
			args = append(args, "%"+str(m, f)+"%")
		} else {
			args = append(args, v)
		}
	}
	return sb.String(), args
}

// userExcelCols 用户实体的 @ExcelColumn 字段(getSql 范围;id/expirationTime 等 int/date 字段
// 在 getMaps 中不包 %,由值形态判断)
var userExcelCols = []string{"accountType", "level", "channel", "handleResult",
	"id", "phone", "account", "password", "username", "companyId",
	"singleDatabase", "singleDatabaseName", "xzqhCode", "xzqhName", "eqbPerId", "wxFaceImg",
	"openid", "officialOpenid", "appClientId", "applys", "msgAlarmFlag", "alarmPhone",
	"wareIdentity", "memberType", "expirationTime", "activateType", "packageId", "nickname",
	"etcName", "duty", "brokerStartDate", "isAgreeCall", "servicerUrl", "isPerfect",
	"followPerson", "followStatus", "lastLoginTime", "isMiniProUser", "isAppUser", "isWebUser",
	"isPromptCare", "createTime", "updateTime", "remark1", "remark2", "remark3", "remark4", "remark5"}

// userListByPage = ChainAuthUserServiceImpl.listByPage(mapper.listByPage):
// u.* + roleNames/enterpriseName 子查询,${searchSql} + CONCAT(IFNULL..) like searchStr
func (s *Service) userListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["page"].(map[string]any); !ok {
		platform.Unknown(w) // Java mapper requires a Page object.
		return
	}
	current, size := pageOf(m)
	searchSql, args := s.userSearchSql(r, m, "handleResult", "applys")
	q := "SELECT u.*, (SELECT GROUP_CONCAT(r.name) FROM chain_auth_use_role ur, chain_auth_role r WHERE ur.role_id=r.id AND ur.user_id=u.id GROUP BY ur.user_id) roleNames," +
		" (SELECT enterprise_name FROM chain_basic_enterprise WHERE id=company_id) enterpriseName FROM chain_auth_user u " +
		searchSql + " AND CONCAT(IFNULL(phone, ''), IFNULL(username, ''), IFNULL(account, '')) LIKE ? ORDER BY create_time DESC"
	countQ := "SELECT COUNT(*) FROM chain_auth_user u " + searchSql + " AND CONCAT(IFNULL(phone, ''), IFNULL(username, ''), IFNULL(account, '')) LIKE ?"
	args = append(args, "%"+searchStrOf(m)+"%")
	page, err := paged(s.DB.Chain, q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeUser(page["records"].([]map[string]any))
	platform.Success(w, page)
}

// userSearchSql 复刻 listByPage 的 searchSql 组装(getSql(排除 handleResult,applys) + 手工条件)
func (s *Service) userSearchSql(r *http.Request, m map[string]any, exclude ...string) (string, []any) {
	un := map[string]bool{}
	for _, e := range exclude {
		un[e] = true
	}
	searchSql, args := getSqlEq(m, userExcelCols, un, nil)
	if _, ok := m["handleResult"]; ok && m["handleResult"] != nil {
		searchSql += " and handle_result = ?"
		args = append(args, str(m, "handleResult"))
	}
	if _, ok := m["applys"]; ok && m["applys"] != nil {
		searchSql += " and applys like ?"
		args = append(args, "%"+str(m, "applys")+"%")
	}
	// Java: jwt companyId 非空且非 "0" 才加(平台账号 company_id=0 不加)
	cu := auth.From(r.Context())
	if cu.CompanyID != 0 {
		searchSql += " and company_id = ?"
		args = append(args, cu.CompanyID)
	}
	if v, ok := m["startTime"]; ok && v != nil {
		searchSql += " and DATE_FORMAT(create_time, '%Y-%m-%d') >= ?"
		args = append(args, str(m, "startTime"))
	}
	if v, ok := m["endTime"]; ok && v != nil {
		searchSql += " and DATE_FORMAT(create_time, '%Y-%m-%d') <= ?"
		args = append(args, str(m, "endTime"))
	}
	if num(m["memberType"]) == 0 && m["isCtiExDate"] == true {
		searchSql += " and expiration_time is not null"
	}
	if v, ok := m["hasPhone"]; ok && v != nil {
		if v == true {
			searchSql += " and (phone is not null and phone!='')"
		} else {
			// Java 原文 `phone==''` 为非法 SQL,执行即错 → 999,行为保持
			searchSql += " and (phone is null or phone=='')"
		}
	}
	return searchSql, args
}

// userSupList = service.supList:getSql(全字段) + 固定条件 + CONCAT(无 IFNULL,null 行被排除)
func (s *Service) userSupList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	current, size := pageOf(m)
	searchSql, args := getSqlEq(m, userExcelCols, map[string]bool{}, nil)
	searchSql += " and handle_result=1 and (account_type =3 or account_type =2)"
	q := "SELECT id, phone, account, username, level," +
		" (SELECT GROUP_CONCAT(r.name) FROM chain_auth_use_role ur, chain_auth_role r WHERE ur.role_id=r.id AND ur.user_id=u.id GROUP BY ur.user_id) roleNames," +
		" (SELECT enterprise_name FROM chain_basic_enterprise WHERE id=company_id) enterpriseName FROM chain_auth_user u " +
		searchSql + " AND CONCAT(phone, username, account) LIKE ? ORDER BY create_time DESC"
	countQ := "SELECT COUNT(*) FROM chain_auth_user u " + searchSql + " AND CONCAT(phone, username, account) LIKE ?"
	args = append(args, "%"+searchStrOf(m)+"%")
	page, err := paged(s.DB.Chain, q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeUser(page["records"].([]map[string]any))
	platform.Success(w, page)
}

// userAdd = controller.add:账号查重(602) + bcrypt + save,返回新 id
func (s *Service) userAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	var one int
	err := s.DB.Chain.QueryRow("SELECT 1 FROM chain_auth_user WHERE account = ? LIMIT 1", m["account"]).Scan(&one)
	if err != nil && err != sql.ErrNoRows {
		platform.Unknown(w)
		return
	}
	if err == nil {
		platform.FailData(w, 602, "账号信息已存在", nil)
		return
	}
	vals := pick(m, userWriteCols...)
	vals["password"] = auth.HashPassword(str(m, "password"))
	if vals["password"] == "" { // bcrypt 拒绝超过 72 字节的口令并返回空串:不能把空口令落库
		platform.Unknown(w)
		return
	}
	delete(vals, "id")
	id, err := platform.InsertNoFill(s.DB.Chain, "chain_auth_user", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, id)
}

// userUpdateById = service.updateChainAuthUser 非认证分支(认证分支见 userAuthEtc)。
// 返回文案 "修改成功"。原方法经 TxAdviceInterceptor 全量事务，后段 NPE 整体回滚。
func (s *Service) userUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if m["isAcOperation"] == true {
		s.userAuthEtc(w, r, m)
		return
	}
	old, err := platform.GetByID(s.DB.Chain, "chain_auth_user", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	vals := pick(m, userWriteCols...)
	if str(m, "handleResult") == "1" {
		vals["level"] = 1
	}
	err = platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		if err := platform.TxUpdateByIDNoFill(tx, "chain_auth_user", m["id"], vals); err != nil {
			return err
		}
		if old == nil {
			return errors.New("user missing") // Java 后段 user.getPhone() NPE
		}
		if m["memberType"] != nil && num(m["memberType"]) != 0 {
			if old["member_type"] == nil {
				return errors.New("memberType null") // 0==user.getMemberType() 拆箱 NPE
			}
			if num(old["member_type"]) == 0 && num(m["memberType"]) == 1 {
				if m["expirationTime"] == nil || m["activateType"] == nil {
					return errors.New("expiration/activateType null") // NPE
				}
				if err := auth.ScheduleMemberExpiryText(r.Context(), str(m, "id"), str(m, "expirationTime")); err != nil {
					return err
				}
				term := "一个月"
				if num(m["activateType"]) == 1 {
					term = "一年"
				}
				platform.SendLegacySMS(r.Context(), str(old, "phone"), []string{term, str(m, "expirationTime")}, "2193101", "")
			}
		}
		if m["warehouseNb"] != nil && old["is_agree_call"] != nil {
			if m["isAgreeCall"] == nil {
				return errors.New("isAgreeCall null") // 拆箱 NPE
			}
			if m["isAgreeCall"] == true {
				// 原库 bus_business_card 无 msg_type 列，Java feign findList 同样失败→999
				cards, err := queryRows(s.DB.Chain, "SELECT * FROM bus_business_card WHERE msg_type=2")
				if err != nil {
					return err
				}
				for _, c := range cards {
					if c["admin_phone"] == nil {
						return errors.New("adminPhone null") // toString() NPE
					}
					platform.SendLegacySMS(r.Context(), str(c, "admin_phone"), []string{str(old, "phone"), str(m, "warehouseNb")}, "1848436", "")
				}
				// 经纪人短信异常被 Java 吞掉，无契约影响
			}
		}
		if err := platform.RedisSetHash(r.Context(), "userId"+str(old, "id"), map[string]string{
			"phone": platform.JnullAny(old["phone"]), "username": platform.JnullAny(old["username"]),
			"wxHeadImg": platform.JnullAny(old["wx_face_img"]), "level": platform.JnullAny(old["level"]),
		}, 0); err != nil {
			return err
		}
		if v, ok := m["wxFaceImg"]; ok && v != nil && str(m, "wxFaceImg") != "" {
			// Java 该步 try-catch 吞异常,行为保持
			_, _ = tx.Exec("UPDATE chain_common_file SET effective = 1 WHERE id = ?", v)
		}
		return nil
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "修改成功")
}

// userAuthEtc = service.authEtc(企业认证通过):企业评分更新、同年认证次数限制、认证记录
func (s *Service) userAuthEtc(w http.ResponseWriter, r *http.Request, m map[string]any) {
	old, err := platform.GetByID(s.DB.Chain, "chain_auth_user", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if etcInfo, ok := m["etcInfo"].(map[string]any); ok {
		vals := pick(etcInfo, "id", "acc_opening_permit_id", "ycoal_silo", "enterprise_name",
			"credit_code", "legal_person", "phone", "address", "is_authentication")
		vals["qua_cert_score"] = 0
		if etcInfo["accOpeningPermitId"] != nil {
			vals["qua_cert_score"] = 1
		}
		vals["bas_info_score"] = 0
		if etcInfo["ycoalSilo"] != nil {
			vals["bas_info_score"] = 1
		}
		if _, err := platform.UpdateByID(s.DB.Chain, "chain_basic_enterprise", etcInfo["id"], vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	if m["companyId"] != nil && old["company_id"] != nil &&
		num(m["companyId"]) == num(old["company_id"]) {
		platform.Success(w, nil)
		return
	}
	var nLogs int
	if err := s.DB.Chain.QueryRow(
		"SELECT COUNT(*) FROM chain_auth_authen_log WHERE user_id = ? AND YEAR(create_time) = YEAR(NOW())",
		m["id"]).Scan(&nLogs); err != nil {
		platform.Unknown(w)
		return
	}
	if nLogs > 1 {
		platform.FailData(w, 602, "如需重新认证联系客服13934568608", nil)
		return
	}
	etc, err := platform.GetByID(s.DB.Chain, "chain_basic_enterprise", m["companyId"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if etc != nil && num(etc["is_authentication"]) != 0 {
		if _, err := platform.Insert(s.DB.Chain, "chain_auth_authen_log", map[string]any{
			"user_id": m["id"], "etc_id": m["companyId"],
			"etc_name": etc["enterprise_name"]}); err != nil {
			platform.Unknown(w)
			return
		}
	}
	vals := pick(m, userWriteCols...)
	if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_auth_user", m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// userDeleteById = removeById + 文案
func (s *Service) userDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if err := platform.DeleteByID(s.DB.Chain, "chain_auth_user", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

// userDelBatch = removeByIds + "删除成功"
func (s *Service) userDelBatch(w http.ResponseWriter, r *http.Request) {
	ids, _ := bodyMap(r)["ids"].([]any)
	if err := platform.DeleteByIDs(s.DB.Chain, "chain_auth_user", ids); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// userCancelBuser = mapper.cancelBuser:清经纪人字段
func (s *Service) userCancelBuser(w http.ResponseWriter, r *http.Request) {
	id := bodyMap(r)["id"]
	if _, err := s.DB.Chain.Exec(
		"UPDATE chain_auth_user SET ware_identity = NULL, servicer_url = NULL WHERE id = ?", id); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// userLogOffUser = service.logOffUser:删用户 + 5 张业务表级联删除(Java 无事务,Go 收 WithTx)
func (s *Service) userLogOffUser(w http.ResponseWriter, r *http.Request) {
	id := bodyMap(r)["id"]
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		if err := platform.TxDeleteByMap(tx, "chain_auth_user", map[string]any{"id": id}); err != nil {
			return err
		}
		for _, q := range []string{
			"DELETE FROM bus_product WHERE create_user_id = ?",
			"DELETE FROM bus_videos WHERE user_id = ?",
			"DELETE FROM bus_business_card WHERE user_id = ?",
			"DELETE FROM bus_visitor_log WHERE visitor_id = ?",
			"DELETE FROM bus_user_care WHERE user_id = ?",
		} {
			if _, err := tx.Exec(q, id); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// inviteCodes 邀请码存储(内存替代 Redis:code+"c"+companyId → 码,TTL 由 effSeconds 决定)
var inviteCodes = map[string]inviteCode{}
var inviteMu sync.Mutex

type inviteCode struct {
	value  string
	expire time.Time
}

// userCreateCode = controller.createCode:生成 code 存 TTL,返回 code+"c"+companyId。
// 降级:短信不涉及;Redis → 内存 map
func (s *Service) userCreateCode(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if m["effSeconds"] == nil || num(m["effSeconds"]) <= 0 {
		platform.Unknown(w)
		return
	}
	cu := auth.From(r.Context())
	code := auth.RandomCode()
	key := code + "c" + strconv.FormatInt(cu.CompanyID, 10)
	secs := int64(num(m["effSeconds"]))
	inviteMu.Lock()
	now := time.Now()
	for k, c := range inviteCodes { // 只增不删会被刷爆:生成时顺手清掉过期的
		if now.After(c.expire) {
			delete(inviteCodes, k)
		}
	}
	inviteCodes[key] = inviteCode{value: code, expire: now.Add(time.Duration(secs) * time.Second)}
	inviteMu.Unlock()
	platform.Success(w, key)
}

// inviteCodeGet 取邀请码(过期/不存在返回 "",false;joinCpy 等后续接口复用)
func inviteCodeGet(key string) (string, bool) {
	inviteMu.Lock()
	defer inviteMu.Unlock()
	c, ok := inviteCodes[key]
	if !ok || time.Now().After(c.expire) {
		return "", false
	}
	return c.value, true
}

// userListUpdByList = ChainAuthUserListController.updateById(变更用户状态):
// handleResult=1 时 level=1 并继承当前操作者的单库设置
func (s *Service) userListUpdByList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	vals := pick(m, userWriteCols...)
	if str(m, "handleResult") == "1" {
		vals["level"] = 1
		op, err := platform.GetByID(s.DB.Chain, "chain_auth_user", auth.From(r.Context()).UserID)
		if err != nil {
			platform.Unknown(w)
			return
		}
		if op != nil {
			if op["single_database"] != nil {
				vals["single_database"] = op["single_database"]
				vals["single_database_name"] = op["single_database_name"]
			}
		}
	}
	if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_auth_user", m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "修改成功")
}

// assignRole = ChainAuthUseRoleServiceImpl.assignRole:删旧插新(Java delByUserId+saveBatch 无事务,Go 收 WithTx)
func (s *Service) assignRole(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	roles, _ := m["userRoles"].([]any)
	err := platform.WithTx(s.DB.Chain, func(tx *sql.Tx) error {
		if _, err := tx.Exec("DELETE FROM chain_auth_use_role WHERE user_id = ?", m["userId"]); err != nil {
			return err
		}
		for _, row := range roles {
			if rm, ok := row.(map[string]any); ok {
				if _, err := platform.TxInsert(tx, "chain_auth_use_role",
					pick(rm, "id", "user_id", "role_id")); err != nil {
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
	platform.Success(w, nil)
}

// selUserRoles = MP list(allEq({user_id})):use_role 实体驼峰 + 载体 null
func (s *Service) selUserRoles(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	rows, err := platform.ListWhere(s.DB.Chain, "chain_auth_use_role", map[string]any{"user_id": m["userId"]})
	if err != nil {
		platform.Unknown(w)
		return
	}
	out := platform.CamelizeRows(rows)
	for _, row := range out {
		fillNulls(row, useRoleExtra)
	}
	platform.Success(w, out)
}

// warehouseBrokerList serves the supply page's broker selectors. The Java
// list maps whole user records; the selectors only need these three columns.
func (s *Service) warehouseBrokerList(w http.ResponseWriter, r *http.Request) {
	user := auth.From(r.Context())
	if user.UserID == 0 {
		platform.NotAuth(w)
		return
	}
	var allowed bool
	err := s.DB.Chain.QueryRowContext(r.Context(), "SELECT EXISTS(SELECT 1 FROM chain_auth_user WHERE id=? AND account_type IN (2,5))", user.UserID).Scan(&allowed)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if !allowed {
		platform.NotAuth(w)
		return
	}
	body := bodyMap(r)
	if len(body) != 1 || str(body, "ware_identity") != "8" {
		platform.Fail(w, 602, "仅支持查询经纪人列表")
		return
	}
	rows, err := s.DB.Chain.QueryContext(r.Context(), "SELECT id, username, phone FROM chain_auth_user WHERE ware_identity = ? ORDER BY id", 8)
	if err != nil {
		platform.Unknown(w)
		return
	}
	brokers, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if brokers == nil {
		brokers = []map[string]any{}
	}
	platform.Success(w, brokers)
}
