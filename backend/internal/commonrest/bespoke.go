package commonrest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func vidDoubleJson(f float64) any {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return json.Number(s)
}

func registerBespoke(mux *http.ServeMux, s *Service) {
	reg := func(path string, h http.HandlerFunc) { mux.HandleFunc("POST /saas-common"+path, authGuard(h)) }
	reg("/chainAreaInfo/selByPet", func(w http.ResponseWriter, r *http.Request) {
		m := bodyMap(r)
		rows := s.qEntEnt2(w, r, ents["ChainAreaInfo"], "SELECT * FROM chain_area_info WHERE parent_code=?", m["parentCode"])
		if rows != nil {
			platform.Success(w, rows)
		}
	})
	for from, to := range map[string]string{
		"/chainBasicMainPage/selectCount": "/chainBasicEnterprise/selectCount",
		"/chainBasicMainPage/ptlist":      "/chainBasicEnterprise/ptlist",
		"/chainBasicBussiness/selEcts":    "/chainBasicEnterprise/ptlist",
	} {
		reg(from, func(w http.ResponseWriter, r *http.Request) {
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/saas-common" + to
			mux.ServeHTTP(w, r2)
		})
	}
	reg("/chainBasicMainPage/countPtLogistArea", s.countPtLogistArea)
	reg("/chainBasicMainPage/selProvinces", s.selProvinces)
	// Controller delegates to the exact same service method.
	reg("/chainCommonLog/selUserExists", s.selUserExists)
	reg("/chainTlzfOrder/selSfEtc", s.selSfEtc)
	reg("/chainTlzfOrder/selPtAmt", s.selPtAmt)
	reg("/chainTlzfOrder/selRelTraceFee", s.selRelTraceFee)
	reg("/chainTlzfOrder/asyncTlzf", s.tlAsync)
	reg("/chainTlzfOrder/confirmOrCancel", s.tlConfirm)
	reg("/chainTlzfOrder/payMoney", s.tlPayMoney)
	reg("/chainTlzfOrder/selAccAmt", s.tlBalance)
	reg("/chainTlzfOrder/selTransResult", s.tlTransResult)
	reg("/chainTlzfOrder/synchronousStatus", s.tlSynchronous)
}
func (s *Service) selSfEtc(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entChainTlzfOrder.validateBody(m) {
		platform.Unknown(w)
		return
	}
	// The Java service uses getSql except company-side special conditions.
	cm := map[string]any{}
	for k, v := range m {
		if k != "saleCompanyId" {
			cm[k] = v
		}
	}
	cond, args, ok := entChainTlzfOrder.getSql(cm)
	if !ok {
		platform.Unknown(w)
		return
	}
	for _, k := range []string{"saleCompanyId"} {
		if m[k] != nil {
			cond += " and " + snake(k) + "=" + javaStr(m[k])
		}
	}
	b, err := queryRows(s.pick(r), "SELECT DISTINCT buy_company_id buyCompanyId,buy_ect_name buyEctName FROM chain_tlzf_order "+cond, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	a, err := queryRows(s.pick(r), "SELECT DISTINCT sale_company_id saleCompanyId,sale_ect_name saleEctName FROM chain_tlzf_order "+cond, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, map[string]any{"buyEtcs": b, "saleEtcs": a})
}
func snake(k string) string {
	for c := range platform.Underscore(map[string]any{k: nil}) {
		return c
	}
	return k
}
func (s *Service) selPtAmt(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id := m["companyId"]
	if id == nil {
		id = currentUser(r).CompanyID
	}
	var v *string
	if err := s.pick(r).QueryRow("SELECT sum(residue_freeze_amt) FROM chain_contract WHERE sales_unit_id=?", id).Scan(&v); err != nil {
		platform.Unknown(w)
		return
	}
	if v == nil {
		platform.Success(w, nil)
	} else {
		platform.Success(w, json.Number(*v))
	}
}
func (s *Service) selRelTraceFee(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entChainTlzfOrder.validateBody(m) {
		platform.Unknown(w)
		return
	}
	id := ""
	if m["contractId"] != nil {
		id = javaStr(m["contractId"])
	}
	var v *string
	if err := s.pick(r).QueryRow("SELECT sum(rel_trace_fee) FROM chain_tlzf_order WHERE trans_type=1 AND trans_status=1 AND rel_operate_amt>0 AND contract_id=" + id).Scan(&v); err != nil {
		platform.Unknown(w)
		return
	}
	if v == nil {
		platform.Success(w, nil)
	} else {
		platform.Success(w, json.Number(*v))
	}
}
func (s *Service) countPtLogistArea(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond := " WHERE loading_area IS NOT NULL AND unloading_area IS NOT NULL"
	if m["xzqhCode"] == "140000" {
		cond += " AND loading_address LIKE '%山西省%' AND unloading_point LIKE '%山西省%'"
	}
	rows, err := queryRows(s.pick(r), "SELECT loading_area loadingArea,max(loading_lnglog) loadingLnglog,unloading_area unloadingArea,max(unloading_lnglog) unloadingLnglog,count(1) logCount FROM chain_logistics"+cond+" GROUP BY loading_area,unloading_area ORDER BY loading_area DESC")
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, rows)
}
func (s *Service) selProvinces(w http.ResponseWriter, r *http.Request) {
	rows, err := queryRows(s.pick(r), "select a.code, a.name from chain_basic_enterprise b, chain_area_info a where b.province_code=a.code GROUP BY province_code")
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, rows)
}
func (s *Service) selUserExists(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	row, err := platform.GetByID(s.pick(r), "chain_auth_user", u.UserID)
	if err != nil {
		platform.Unknown(w)
		return
	}
	op := "用户id为" + strconv.FormatInt(u.UserID, 10) + "，查询失败"
	if row != nil {
		op = "用户id为" + strconv.FormatInt(u.UserID, 10) + "，手机号为" + javaStr(row["phone"]) + "，查询成功"
	}
	_, err = platform.InsertNoFill(s.pick(r), "chain_common_log", map[string]any{"user_id": u.UserID, "user_name": u.Account, "operation": op, "operation_type": 12})
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// EnterpriseList is shared with the phase-one dashboard owner to preserve entity serialization.
func (s *Service) EnterpriseList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	e := entChainBasicEnterprise
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	cond, args, ok := e.getSql(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	// getMaps(entity,"provinceCode,cityCode,areaCode,countyCode") keeps these Strings unwrapped.
	argi := 0
	for _, f := range e.q {
		v, has := m[f.Name]
		if !has || v == nil || v == "0" {
			continue
		}
		if strings.Contains("provinceCode,cityCode,areaCode,countyCode", f.Name) {
			args[argi] = v
		}
		argi++
	}
	search := "%null%"
	if m["searchStr"] != nil {
		search = "%" + javaStr(m["searchStr"]) + "%"
	}
	q := "SELECT * FROM chain_basic_enterprise " + cond + " and CONCAT(IFNULL(enterprise_name, ''), IFNULL(credit_code, ''), IFNULL(address, ''), IFNULL(legal_person, ''), IFNULL(controller_person, ''), IFNULL(contacts_person, ''), IFNULL(enterprise_type, ''), IFNULL(type, ''), IFNULL(reg_time, ''), IFNULL(rate_number, ''), IFNULL(sh_time, ''), IFNULL(phone, ''), IFNULL(leg_identity, ''), IFNULL(reg_money, ''), IFNULL(reg_adress, ''), IFNULL(businesser, ''), IFNULL(county_name, ''), IFNULL(boss, '')) like ? and is_chain_etc=1 and rate_number is not null and rate_number!='' order by create_time desc"
	rows := s.qEntEnt2(w, r, e, q, append(args, search)...)
	if rows != nil {
		platform.Success(w, rows)
	}
}
