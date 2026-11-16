// biz.go 二期收尾的合同供应商(saas-contract)与仓库侧红包范围/访问日志(saas-warehouse)。
// Java 源:ChainContractSupplierController/BusRedEnvelopeRangeController/BusProLogController.selMemberData。
// 行为保持:updateBatchRanges 无事务;空列表/空配置表时 Java 直接抛异常走 999
package authadmin

import (
	"encoding/json"
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// ── 供应商 chain_contract_supplier ──

var supplierCols = []string{"id", "company_id", "enterprise_name", "enterprise_nature",
	"credit_code", "business_address", "coal_mine", "annual_output", "coal_yard", "pro1",
	"pro2", "sulfur", "volatilization", "water", "calorific_value", "partner_type", "cust_type",
	"amount_receivable", "amount_payable", "pro_remark", "create_time", "update_time",
	"remark1", "remark2", "remark3", "remark4", "remark5"}

var supplierFill = []string{"id", "companyId", "enterpriseName",
	"enterpriseNature", "creditCode", "businessAddress", "coalMine", "annualOutput", "coalYard",
	"pro1", "pro2", "sulfur", "volatilization", "water", "calorificValue", "partnerType",
	"custType", "amountReceivable", "amountPayable", "proRemark", "createTime", "updateTime",
	"remark1", "remark2", "remark3", "remark4", "remark5", "searchStr", "page", "ids", "increaseType"}

// shapeSupplier:DECIMAL 列(10,0)/(10,2) 保持原样小数位输出
func shapeSupplier(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		for _, c := range []string{"sulfur", "volatilization", "water", "calorificValue",
			"amountReceivable", "amountPayable"} {
			if s, ok := r[c].(string); ok && s != "" {
				r[c] = json.Number(s)
			}
		}
		fillNulls(r, supplierFill)
	}
	return out
}

// supplierListByPage = ChainContractSupplierServiceImpl.listByPage:
// where 1=1 [and partner_type = N] and CONCAT(enterprise_name) like '%${searchStr}%';
// ${} 原样拼接(含注入语义),searchStr 缺失为字面 "null"
func (s *Service) supplierListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	current, size := pageOf(m)
	sql := "where 1=1 "
	if v, ok := m["partnerType"]; ok && v != nil {
		sql += " and partner_type = " + str(m, "partnerType")
	}
	cond := sql + " and CONCAT(enterprise_name) like '%" + searchStrOf(m) + "%'"
	q := "SELECT * FROM chain_contract_supplier " + cond + " ORDER BY create_time DESC"
	page, err := paged(s.DB.Chain, q, "SELECT COUNT(*) FROM chain_contract_supplier "+cond, nil, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeSupplier(page["records"].([]map[string]any))
	platform.Success(w, page)
}

func (s *Service) supplierGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	row, err := platform.GetByID(s.DB.Chain, "chain_contract_supplier", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, shapeSupplier([]map[string]any{row})[0])
}

func (s *Service) supplierAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	vals := pick(m, supplierCols...)
	delete(vals, "id") // @TableId(IdType.AUTO):body id 被忽略,自增
	id, err := platform.InsertNoFill(s.DB.Chain, "chain_contract_supplier", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, platform.EntityEcho(m, supplierFill, map[string]any{"id": id}))
}

func (s *Service) supplierUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	vals := pick(m, supplierCols...)
	delete(vals, "id")
	if len(vals) == 0 {
		platform.Unknown(w) // MP updateById SET 空 → MybatisPlusException
		return
	}
	if _, err := platform.UpdateByIDNoFill(s.DB.Chain, "chain_contract_supplier", m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, platform.EntityEcho(m, supplierFill, nil))
}

func (s *Service) supplierDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, _ := m["ids"].([]any)
	if err := platform.DeleteByIDs(s.DB.Chain, "chain_contract_supplier", ids); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// ── 红包范围 bus_red_envelope_range ──

var brelrgeCols = []string{"id", "min_amount", "max_amount", "receive_type",
	"create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"}

var brelrgeFill = []string{"id", "minAmount", "maxAmount", "receiveType",
	"createTime", "updateTime", "remark1", "remark2", "remark3", "remark4", "remark5",
	"searchStr", "page", "ids", "redUplimitDays"}

// amountConfig 对照 selAmountConfig:chain_global_configuration 首行;空表时 Java NPE → 999
func (s *Service) amountConfig() (map[string]any, bool) {
	rows, err := queryCam(s.DB.Chain, "SELECT * FROM chain_global_configuration LIMIT 1")
	if err != nil || len(rows) == 0 {
		return nil, false
	}
	return rows[0], true
}

// brelrgeList = listRange:allEq + 每行注入 redUplimitDays(取自全局配置)
func (s *Service) brelrgeList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, brelrgeCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.DB.Chain, "bus_red_envelope_range", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	config, ok := s.amountConfig()
	if !ok {
		platform.Unknown(w) // configMap.get("red_uplimit_days").toString() NPE
		return
	}
	out := platform.CamelizeRows(rows)
	days := int64(num(config["redUplimitDays"])) // amountConfig 走 queryCam,键已驼峰化
	for _, r := range out {
		r["redUplimitDays"] = days
		fillNulls(r, brelrgeFill)
	}
	platform.Success(w, out)
}

// brelrgeUpdBatch = updateBatchRanges:逐行更新 + 全局配置 red_uplimit_days 取首行值。
// 空列表时 Java get(0) 越界 → 999;redUplimitDays 缺失时 ${} 拼出字面 NULL
func (s *Service) brelrgeUpdBatch(w http.ResponseWriter, r *http.Request) {
	var items []any
	_ = json.NewDecoder(r.Body).Decode(&items)
	if len(items) == 0 {
		platform.Unknown(w)
		return
	}
	for _, it := range items {
		em, _ := it.(map[string]any)
		if em == nil || em["id"] == nil {
			platform.Unknown(w) // MP updateById 空 id 抛异常
			return
		}
		if _, err := platform.UpdateByID(s.DB.Chain, "bus_red_envelope_range", em["id"], pick(em, brelrgeCols...)); err != nil {
			platform.Unknown(w)
			return
		}
	}
	config, ok := s.amountConfig()
	if !ok {
		platform.Unknown(w) // configMap.get("id").toString() NPE
		return
	}
	first, _ := items[0].(map[string]any)
	q := "UPDATE chain_global_configuration SET red_uplimit_days="
	args := []any{}
	if v, ok := first["redUplimitDays"]; ok && v != nil {
		q += "?"
		args = append(args, str(first, "redUplimitDays"))
	} else {
		q += "NULL"
	}
	q += " WHERE id=?"
	args = append(args, str(config, "id"))
	if _, err := s.DB.Chain.Exec(q, args...); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// ── 访问日志 bus_pro_log.selMemberData ──

// proLogSelMemberData = BusProLogServiceImpl.selMemberData:三标量(最后访问时间/拨打客服次数/点击支付次数)
func (s *Service) proLogSelMemberData(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	uid := str(m, "userId")
	res := map[string]any{"lastVisTime": nil, "callTimes": int64(0), "clickTimes": int64(0)}
	if vis, err := queryCam(s.DB.Chain,
		"SELECT create_time FROM bus_pro_log WHERE operate_obj=6 AND operate_type=7 AND user_id=? ORDER BY create_time DESC LIMIT 1", uid); err == nil && len(vis) > 0 {
		if sv, ok := vis[0]["createTime"].(string); ok && len(sv) == 19 {
			// LocalDateTime 无 @JsonFormat → ISO T 格式
			res["lastVisTime"] = sv[:10] + "T" + sv[11:]
		}
	}
	var callTimes, clickTimes int64
	_ = s.DB.Chain.QueryRow(
		"SELECT IFNULL(count(1),0) FROM bus_pro_log WHERE operate_obj=7 AND operate_type=8 AND user_id=?", uid).Scan(&callTimes)
	_ = s.DB.Chain.QueryRow(
		"SELECT IFNULL(count(1),0) FROM bus_pro_log WHERE operate_obj=8 AND operate_type=10 AND user_id=?", uid).Scan(&clickTimes)
	res["callTimes"], res["clickTimes"] = callTimes, clickTimes
	platform.Success(w, res)
}
