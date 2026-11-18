package commonrest

import (
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func registerEnterprise(mux *http.ServeMux, s *Service) {
	reg := func(p string, h http.HandlerFunc) {
		mux.HandleFunc("POST /saas-common/chainBasicEnterprise/"+p, authGuard(h))
	}
	reg("add", s.enterpriseAdd)
	reg("selProvinces", s.selProvinces)
	reg("selEtcName", func(w http.ResponseWriter, r *http.Request) {
		// The local Qichacha channel has no successful response. Java's catch returns an empty JSONArray.
		platform.Success(w, []any{})
	})
	reg("selEtcInfo", func(w http.ResponseWriter, r *http.Request) {
		m := bodyMap(r)
		rows, err := queryRows(s.pick(r), "SELECT * FROM chain_basic_enterprise WHERE enterprise_name <=> ? AND is_authentication=1", m["enterpriseName"])
		if err != nil || len(rows) > 1 {
			platform.Fail(w, 603, "")
			return
		}
		if len(rows) == 1 {
			platform.Success(w, entChainBasicEnterprise.shapeOne(rows[0]))
			return
		}
		platform.Fail(w, 603, "")
	})
	for _, p := range []string{"tlOpenBind", "tlChangeAcc"} {
		mux.HandleFunc("POST /saas-common/chainTlzfConfig/"+p, authGuard(func(w http.ResponseWriter, r *http.Request) {
			m := bodyMap(r)
			if !entChainTlzfConfig.validateBody(m) {
				platform.Unknown(w)
				return
			}
			platform.Fail(w, 603, "请求失败")
		}))
	}
}

// enterpriseAdd retains the Java nontransactional sequence: a user-update exception
// can follow a successfully inserted enterprise; saveEnterprise's returned failure is ignored.
func (s *Service) enterpriseAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	e := entChainBasicEnterprise
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	var n int
	if err := d.QueryRow("SELECT count(*) FROM chain_auth_user WHERE account=?", m["account"]).Scan(&n); err != nil {
		platform.Unknown(w)
		return
	}
	if n > 0 {
		platform.Fail(w, 999, "账号信息已存在")
		return
	}
	cond := map[string]any{}
	for _, p := range []struct{ json, col string }{{"enterpriseName", "enterprise_name"}, {"creditCode", "credit_code"}} {
		if m[p.json] != nil {
			cond[p.col] = m[p.json]
		}
	}
	rows, err := platform.ListWhere(d, "chain_basic_enterprise", cond)
	if err != nil || len(rows) > 1 {
		platform.Unknown(w)
		return
	}
	res := map[string]any{"hasAdmin": false, "result": "操作成功"}
	var cid any
	if len(rows) == 0 {
		vals := pick(m, e.dbcols)
		delete(vals, "id")
		vals["is_authentication"] = true
		vals["sign_free_count"] = 100
		vals["sign_charge_count"] = 0
		vals["status_mark"] = 1
		configs, err := queryRows(d, "SELECT id,flow_type FROM chain_eqb_config")
		if err != nil {
			platform.Unknown(w)
			return
		}
		for _, c := range configs {
			if javaStr(c["flow_type"]) == "1" || javaStr(c["flow_type"]) == "2" {
				vals["pub_flow_id"] = c["id"]
			}
		}
		var max *string
		if err := d.QueryRow("SELECT max(enterprise_code) FROM chain_basic_enterprise").Scan(&max); err != nil {
			platform.Unknown(w)
			return
		}
		code := "0000001"
		if max != nil {
			code = *max
		}
		n, err := strconv.Atoi(code)
		if err != nil {
			platform.Unknown(w)
			return
		}
		vals["enterprise_code"] = autoGenericCode(strconv.Itoa(n+1), 7)
		id, err := platform.InsertNoFill(d, "chain_basic_enterprise", vals)
		if err != nil {
			platform.Success(w, res)
			return
		}
		cid = id
		if list, ok := m["filelist"].([]any); ok {
			for _, v := range list {
				fm, ok := v.(map[string]any)
				if !ok {
					platform.Unknown(w)
					return
				}
				fv := pick(fm, entChainCommonFile.dbcols)
				delete(fv, "id")
				fv["company_id"] = id
				fv["file_type"] = "3"
				fv["owner_type"] = "1"
				fv["owner_id"] = currentUser(r).UserID
				fv["relation_id"] = id
				fv["relation_type"] = 6
				fv["effective"] = true
				if _, err := platform.UpdateByIDNoFill(d, "chain_common_file", fm["id"], fv); err != nil {
					platform.Unknown(w)
					return
				}
			}
		}
	} else {
		cid = rows[0]["id"]
		if err := d.QueryRow("SELECT count(*) FROM chain_auth_user WHERE company_id=? AND account_type=3 AND level=1", cid).Scan(&n); err != nil {
			platform.Unknown(w)
			return
		}
		if n > 0 {
			res["hasAdmin"] = true
			platform.Success(w, res)
			return
		}
	}
	password, ok := m["password"].(string)
	if !ok {
		platform.Unknown(w)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		platform.Unknown(w)
		return
	}
	vals := map[string]any{"company_id": cid, "account_type": "3", "account": m["account"], "handle_result": "0", "level": 1, "username": m["adminName"], "password": string(hashed)}
	if _, err := platform.UpdateByID(d, "chain_auth_user", m["userId"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, res)
}
