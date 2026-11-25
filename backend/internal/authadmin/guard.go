package authadmin

import (
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// 账号边界。原 Java 的用户 / 角色 / 平台配置写接口只按 body 的 id 直写、不看调用者:
// 任何登录用户都能 resetPwd 改管理员密码、add 出 accountType=5 的账号、改全局支付开关。
// 按库内实际的账号类型收口(account_type:1 普通/个人、2 平台管理员、3 企业管理员(企业注册时建)、4 大屏、5 超管):
//
//	平台管理员(2/5)                        → 不限制;
//	企业管理员(3,或持有用户管理菜单的账号) → 只能管本企业的非平台账号、非平台角色;
//	其他登录用户                           → 只能改自己的资料,平台属性等字段一律不收。

func isPlatformAdmin(cu *auth.CurrentUser) bool {
	return cu.UserID != 0 && (cu.AccountType == "2" || cu.AccountType == "5")
}

func isPlatformRow(u map[string]any) bool {
	return str(u, "account_type") == "2" || str(u, "account_type") == "5"
}

func deny(w http.ResponseWriter) { platform.Fail(w, 602, "无权限操作") }

// requirePlatform 平台管理员专属接口的统一入口
func (s *Service) requirePlatform(w http.ResponseWriter, r *http.Request) bool {
	if !isPlatformAdmin(auth.From(r.Context())) {
		deny(w)
		return false
	}
	return true
}

// platformOnly 把 handler 包成平台管理员专属(菜单是全平台共用的一张表,增删改只能由平台来做)。
func (s *Service) platformOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.requirePlatform(w, r) {
			return
		}
		h(w, r)
	}
}

// isEnterpriseAdmin 企业管理员:企业注册时建的主账号(account_type=3),或被授予用户管理菜单的账号。
func (s *Service) isEnterpriseAdmin(r *http.Request, cu *auth.CurrentUser) bool {
	if cu.UserID == 0 || cu.CompanyID == 0 {
		return false
	}
	if cu.AccountType == "3" {
		return true
	}
	var ok bool
	err := s.DB.Chain.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM chain_auth_user u WHERE u.id=? AND u.company_id=? AND (
		EXISTS(SELECT 1 FROM chain_auth_use_role ur, chain_auth_menu_role mr, chain_auth_menu m
			WHERE ur.role_id=mr.role_id AND mr.menu_id=m.id AND ur.user_id=u.id AND m.route_name IN ('UserList','PtUserList'))
		OR EXISTS(SELECT 1 FROM chain_auth_menu_user mu, chain_auth_menu m
			WHERE mu.menu_id=m.id AND mu.user_id=u.id AND m.route_name IN ('UserList','PtUserList'))))`,
		cu.UserID, cu.CompanyID).Scan(&ok)
	return err == nil && ok
}

// canManage 调用者能否管理目标账号(target 为目标行,不存在时为 nil)。
func (s *Service) canManage(r *http.Request, target map[string]any) bool {
	cu := auth.From(r.Context())
	if isPlatformAdmin(cu) {
		return true
	}
	if cu.UserID == 0 || target == nil {
		return false
	}
	if num(target["id"]) == float64(cu.UserID) {
		return true // 自己
	}
	return s.isEnterpriseAdmin(r, cu) && !isPlatformRow(target) && num(target["company_id"]) == float64(cu.CompanyID)
}

// canManageID 同 canManage,按 id 取目标行。
func (s *Service) canManageID(r *http.Request, id any) bool {
	if isPlatformAdmin(auth.From(r.Context())) {
		return true
	}
	target, err := platform.GetByID(s.DB.Chain, "chain_auth_user", id)
	if err != nil {
		return false
	}
	return s.canManage(r, target)
}

// canManageRole 非平台管理员只能动本企业的非平台角色(role_type!=1)。
func (s *Service) canManageRole(r *http.Request, id any) bool {
	cu := auth.From(r.Context())
	if isPlatformAdmin(cu) {
		return true
	}
	if !s.isEnterpriseAdmin(r, cu) {
		return false
	}
	role, err := platform.GetByID(s.DB.Chain, "chain_auth_role", id)
	if err != nil || role == nil || num(role["role_type"]) == 1 {
		return false
	}
	return num(role["company_id"]) == 0 || num(role["company_id"]) == float64(cu.CompanyID)
}

// scrubUserVals 非平台管理员写账号时的字段约束:平台属性、会员、经纪人属性一律不收;
// 企业管理员可变更本企业账号的 handle_result(用户列表页的启用 / 停用),层级只能是 1。
func scrubUserVals(cu *auth.CurrentUser, vals map[string]any, enterpriseAdmin bool) {
	if isPlatformAdmin(cu) {
		return
	}
	for _, c := range []string{"account_type", "company_id", "applys", "member_type", "expiration_time",
		"activate_type", "package_id", "broker_start_date", "servicer_url"} {
		delete(vals, c)
	}
	if !enterpriseAdmin {
		for _, c := range []string{"handle_result", "level", "single_database", "single_database_name"} {
			delete(vals, c)
		}
	} else if v, ok := vals["level"]; !ok || num(v) != 1 {
		delete(vals, "level")
	}
}

// scrubUserBody 非平台管理员不能靠请求体触发会员开通的副作用(到期排期、开通短信)。
func scrubUserBody(cu *auth.CurrentUser, m map[string]any) {
	if isPlatformAdmin(cu) {
		return
	}
	for _, k := range []string{"memberType", "expirationTime", "activateType"} {
		delete(m, k)
	}
}

// entPlatformOnly 企业档案里只有平台管理员能改的列(认证状态、签章份数、服务期、各项评分、商户与通联配置)。
var entPlatformOnly = []string{"is_authentication", "sign_free_count", "sign_surplus_count", "sign_charge_count",
	"service_expiration_time", "qua_cert_score", "bas_info_score", "coal_verify_score", "upd_share_score",
	"report_score", "qualified_score", "upd_price_score", "enterprise_code", "merchant_number",
	"merchant_account", "rate_number", "eqb_org_id", "pub_flow_id", "per_flow_id", "tl_sign_num",
	"tl_acct_sub_no", "tl_mcht_no", "tl_open_time", "status_mark", "is_chain_etc", "is_business_etc",
	"is_sign_etc", "is_member_etc", "group_id"}
