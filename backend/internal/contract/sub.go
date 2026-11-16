// sub.go 三期剩余:chainContractExtend、chainContractSupplier 增量端点、
// EMain 电子签章(esign 包)DB 层 flow/singers/files 与外部调用失败语义。
// esign 三实体时间列均有 @JsonFormat(空格格式);files 有 fill(createTime/updateTime),
// flow/singers/extend 无 fill(与二期 chain_contract_supplier 同语义)
package contract

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// ── chainContractExtend ──

var extendCols = []string{"id", "company_id", "contract_id", "dict_id", "d_type", "d_key",
	"ext_type", "d_value", "create_time", "update_time", "remark1", "remark2", "remark3",
	"remark4", "remark5"}

var extendCarrier = []string{"searchStr", "page", "ids"}

func shapeExtend(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		// Java Beans getter getDType/getDKey/getDValue → Jackson 输出 dtype/dkey/dvalue(全小写,实测)
		for _, pair := range [][2]string{{"dType", "dtype"}, {"dKey", "dkey"}, {"dValue", "dvalue"}} {
			if v, ok := r[pair[0]]; ok {
				delete(r, pair[0])
				r[pair[1]] = v
			}
		}
		fillNulls(r, extendCarrier)
	}
	return out
}

func (s *Service) extendList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, extendCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_contract_extend", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, shapeExtend(rows))
}

// extendListByPage = mapper:`CONCAT(id) like '%${searchStr}%' order by create_time desc`
func (s *Service) extendListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	current, size := pageOf(m)
	cond := "WHERE CONCAT(id) like '%" + searchStrOf(m) + "%'"
	page, err := paged(s.pick(r), "SELECT * FROM chain_contract_extend "+cond+" ORDER BY create_time DESC",
		"SELECT COUNT(*) FROM chain_contract_extend "+cond, nil, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeExtend(page["records"].([]map[string]any))
	platform.Success(w, page)
}

func (s *Service) extendGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	row, err := platform.GetByID(s.pick(r), "chain_contract_extend", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if row == nil {
		platform.Success(w, nil) // MP getById 无行 → null
		return
	}
	platform.Success(w, shapeExtend([]map[string]any{row})[0])
}

func (s *Service) extendAdd(w http.ResponseWriter, r *http.Request) {
	s.writeAdd(w, r, "chain_contract_extend", bodyMap(r), extendCols, extendRename, extendEchoKeys, false, nil)
}

func (s *Service) extendUpdateById(w http.ResponseWriter, r *http.Request) {
	s.writeUpdate(w, r, "chain_contract_extend", bodyMap(r), extendCols, extendRename, extendEchoKeys, false)
}

func (s *Service) extendDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, err := s.pick(r).Exec("DELETE FROM chain_contract_extend WHERE id = ?", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

func (s *Service) extendDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, _ := m["ids"].([]any)
	if err := platform.DeleteByIDs(s.pick(r), "chain_contract_extend", ids); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// echoKeys 实体驼峰全键(实体列驼峰 + 载体键),配合 platform.EntityEcho 用
func camelKeys(cols []string) []string {
	out := make([]string, len(cols))
	for i, c := range cols {
		out[i] = camel(c)
	}
	return out
}

// bodyVals body 驼峰键 → DB 下划线列值表(仅非 nil);rename 覆盖列→JSON 键映射
// (extend 实体 getDType 类 getter → Jackson 只绑定 dtype/dkey/dvalue 全小写键,实测)
func bodyVals(m map[string]any, cols []string, rename map[string]string) map[string]any {
	vals := map[string]any{}
	for _, c := range cols {
		k := camel(c)
		if rename != nil {
			if rk, ok := rename[c]; ok {
				k = rk
			}
		}
		if v, ok := m[k]; ok && v != nil {
			vals[c] = v
		}
	}
	return vals
}

var extendRename = map[string]string{"d_type": "dtype", "d_key": "dkey", "d_value": "dvalue"}

// 各实体回显键 = 实体列驼峰(extend 用小写特例;flow 剔除表有列实体没有的 custom_biz_num)+ 载体键
var extendEchoKeys = append([]string{"id", "companyId", "contractId", "dictId", "dtype",
	"dkey", "extType", "dvalue", "createTime", "updateTime", "remark1", "remark2",
	"remark3", "remark4", "remark5"}, extendCarrier...)

var flowEchoKeys = func() []string {
	ks := camelKeys(flowCols)
	out := make([]string, 0, len(ks)+len(flowCarrier))
	for _, k := range ks {
		if k != "customBizNum" { // 表有列实体没有 → JSON 无此键
			out = append(out, k)
		}
	}
	return append(out, flowCarrier...)
}()

var singersEchoKeys = append(camelKeys(singersCols), singersCarrier...)

var filesEchoKeys = append(camelKeys(filesCols), filesCarrier...)

// writeUpdate 通用 updateById:MP updateById 语义——SET 空(仅 id)→ MybatisPlusException 999;
// SET 非空 0 行也 200;响应 = 请求实体全键回显
func (s *Service) writeUpdate(w http.ResponseWriter, r *http.Request, table string, m map[string]any,
	cols []string, rename map[string]string, echoKeys []string, fill bool) {
	vals := bodyVals(m, cols, rename)
	delete(vals, "id")
	if len(vals) == 0 {
		platform.Unknown(w)
		return
	}
	var err error
	if fill {
		_, err = platform.UpdateByID(s.pick(r), table, m["id"], vals)
	} else {
		_, err = platform.UpdateByIDNoFill(s.pick(r), table, m["id"], vals)
	}
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, platform.EntityEcho(m, echoKeys, nil))
}

// writeAdd 通用 add:实体 @TableId(AUTO)——忽略 body id 自增;响应 = 内存实体全键回显 + 回填 id
func (s *Service) writeAdd(w http.ResponseWriter, r *http.Request, table string, m map[string]any,
	cols []string, rename map[string]string, echoKeys []string, fill bool, extra map[string]any) {
	vals := bodyVals(m, cols, rename)
	delete(vals, "id")
	var id int64
	var err error
	if fill {
		id, err = platform.Insert(s.pick(r), table, vals)
	} else {
		id, err = platform.InsertNoFill(s.pick(r), table, vals)
	}
	if err != nil {
		platform.Unknown(w)
		return
	}
	if extra == nil {
		extra = map[string]any{}
	}
	extra["id"] = id
	platform.Success(w, platform.EntityEcho(m, echoKeys, extra))
}

// ── chainContractSupplier 增量(listByPage/getById/add/updateById/delBatch 在 authadmin 包)──

var supplierAllCols = []string{"id", "company_id", "enterprise_name", "enterprise_nature",
	"credit_code", "business_address", "coal_mine", "annual_output", "coal_yard", "pro1",
	"pro2", "sulfur", "volatilization", "water", "calorific_value", "partner_type", "cust_type",
	"amount_receivable", "amount_payable", "pro_remark", "create_time", "update_time",
	"remark1", "remark2", "remark3", "remark4", "remark5"}

func (s *Service) supplierList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, supplierAllCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_contract_supplier", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, shapeSupplierList(rows))
}

// shapeSupplierList = Java 实体直出:DECIMAL 列数字输出 + 实体载体键(searchStr/page/ids/increaseType)
var supplierCarrier = []string{"searchStr", "page", "ids", "increaseType"}

func shapeSupplierList(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		for _, c := range []string{"sulfur", "volatilization", "water", "calorificValue",
			"amountReceivable", "amountPayable"} {
			if s, ok := r[c].(string); ok && s != "" {
				r[c] = json.Number(s)
			}
		}
		fillNulls(r, supplierCarrier)
	}
	return out
}

func (s *Service) supplierDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, err := s.pick(r).Exec("DELETE FROM chain_contract_supplier WHERE id = ?", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

// supplierMoney increaseMoney/decreaseMoney:Java 以 body 的 amountReceivable/amountPayable
// (按 increaseType)作为金额源做 BigDecimal.add——body 缺该键或 DB 原值为 NULL → NPE → 999;
// increaseType 无效 → 仅原值回写(行为保持)
func (s *Service) supplierMoney(w http.ResponseWriter, r *http.Request, sign float64) {
	m := bodyMap(r)
	d := s.pick(r)
	row, err := platform.GetByID(d, "chain_contract_supplier", m["id"])
	if err != nil || row == nil {
		platform.Unknown(w) // getById null NPE
		return
	}
	col := "amount_payable"
	key := "amountPayable"
	if str(m, "increaseType") == "1" {
		col = "amount_receivable"
		key = "amountReceivable"
	}
	if row[col] == nil || m[key] == nil {
		platform.Unknown(w) // BigDecimal null NPE
		return
	}
	cur := num(row[col])
	amt := num(m[key])
	val := cur + sign*amt
	if _, err := d.Exec("UPDATE chain_contract_supplier SET "+col+" = ? WHERE id = ?", val, m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "修改成功")
}

func (s *Service) supplierIncreaseMoney(w http.ResponseWriter, r *http.Request) {
	s.supplierMoney(w, r, 1)
}

func (s *Service) supplierDecreaseMoney(w http.ResponseWriter, r *http.Request) {
	s.supplierMoney(w, r, -1)
}

// ── EMain 电子签章(esgN 包 DB 层)──

var flowCols = []string{"id", "company_id", "contract_id", "topic_name", "contract_name",
	"field", "launch_type", "etc_name", "participants", "phone", "launch_per_id",
	"launch_pub_id", "eqb_flow_id", "flow_status", "custom_biz_num", "create_time",
	"update_time", "remark1", "remark2", "remark3", "remark4", "remark5"}

var flowCarrier = []string{"searchStr", "page", "ids", "userInfo", "orgInfo", "singers",
	"signers", "searchType", "signSurplusCount"}

// shapeEqbFlow flow 实体 JSON(custom_biz_num 表有列实体没有 → JSON 无此键,先剔除)
func shapeEqbFlow(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		delete(r, "customBizNum")
		fillNulls(r, flowCarrier)
	}
	return out
}

var singersCols = []string{"id", "company_id", "sign_etc_name", "sign_name", "sign_phone",
	"sign_etc_id", "sign_per_id", "sys_flow_id", "confirm_launch", "sign_role",
	"is_launch_role", "create_time", "update_time", "remark1", "remark2", "remark3",
	"remark4", "remark5"}

var singersCarrier = []string{"searchStr", "page", "ids"}

var singersBoolCols = []string{"confirmLaunch", "isLaunchRole"}

func shapeEqbSingers(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		for _, c := range singersBoolCols {
			if v, ok := r[c]; ok && v != nil {
				r[c] = num(v) != 0
			}
		}
		fillNulls(r, singersCarrier)
	}
	return out
}

var filesCols = []string{"id", "name", "field", "type", "sys_flow_id", "eqb_flow_id",
	"create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"}

var filesCarrier = []string{"searchStr", "page", "ids", "startTime", "endTime"}

func shapeEqbFiles(rows []map[string]any) []map[string]any {
	out := platform.CamelizeRows(rows)
	for _, r := range out {
		fillNulls(r, filesCarrier)
	}
	return out
}

// ── flow CRUD ──

func (s *Service) flowList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, flowCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_eqb_flow", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, shapeEqbFlow(rows))
}

// flowListByPage = ChainEqbFlowServiceImpl.listByPage:searchType 解引用(缺 → 999),
// 1 我发起(phone)/2 我参与(confirm_launch=1)/3 待确认(confirm_launch=0)
func (s *Service) flowListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	u := currentUser(r)
	current, size := pageOf(m)
	st, ok := m["searchType"]
	if !ok || st == nil {
		platform.Unknown(w) // getSearchType() 解引用 NPE
		return
	}
	searchSql, args := getSqlEq(m, []string{"id", "companyId", "contractId", "topicName",
		"contractName", "field", "launchType", "etcName", "participants", "phone",
		"launchPerId", "launchPubId", "eqbFlowId", "flowStatus", "createTime", "updateTime",
		"remark1", "remark2", "remark3", "remark4", "remark5"}, nil,
		map[string]bool{"id": true, "companyId": true, "contractId": true, "flowStatus": true})
	switch num(st) {
	case 1:
		var phone string
		_ = d.QueryRow("SELECT IFNULL(phone,'') FROM chain_auth_user WHERE id = ?", u.UserID).Scan(&phone)
		searchSql += " and phone = ?"
		args = append(args, phone)
	case 2:
		var phone string
		_ = d.QueryRow("SELECT IFNULL(phone,'') FROM chain_auth_user WHERE id = ?", u.UserID).Scan(&phone)
		searchSql += " and EXISTS (SELECT 1 from chain_eqb_singers where sys_flow_id=f.id and sign_phone=? and confirm_launch=1)"
		args = append(args, phone)
	case 3:
		var phone string
		_ = d.QueryRow("SELECT IFNULL(phone,'') FROM chain_auth_user WHERE id = ?", u.UserID).Scan(&phone)
		searchSql += " and EXISTS (SELECT 1 from chain_eqb_singers where sys_flow_id=f.id and sign_phone=? and confirm_launch=0)"
		args = append(args, phone)
	}
	var ss any
	if v, ok := m["searchStr"]; ok && v != nil {
		ss = "%" + str(m, "searchStr") + "%"
	}
	args = append(args, ss)
	q := "SELECT * FROM chain_eqb_flow f " + searchSql + " and CONCAT(id) like ? order by create_time desc"
	countQ := "SELECT count(*) FROM chain_eqb_flow f " + searchSql + " and CONCAT(id) like ?"
	page, err := paged(d, q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeEqbFlow(page["records"].([]map[string]any))
	platform.Success(w, page)
}

// flowGetById = Java controller 在查询后调 SignHelper.selSigners(请求体的 eqbFlowId):
// 请求体不带 eqbFlowId → null NPE;带 eqbFlowId → Emain 外部签章 API 调用失败。
// 两路在本环境均 999,实测恒 999(行为保持)
func (s *Service) flowGetById(w http.ResponseWriter, r *http.Request) {
	platform.Unknown(w)
}

func (s *Service) flowAdd(w http.ResponseWriter, r *http.Request) {
	s.writeAdd(w, r, "chain_eqb_flow", bodyMap(r), flowCols, nil, flowEchoKeys, false, nil)
}

func (s *Service) flowUpdateById(w http.ResponseWriter, r *http.Request) {
	s.writeUpdate(w, r, "chain_eqb_flow", bodyMap(r), flowCols, nil, flowEchoKeys, false)
}

func (s *Service) flowDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, err := s.pick(r).Exec("DELETE FROM chain_eqb_flow WHERE id = ?", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

func (s *Service) flowDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, _ := m["ids"].([]any)
	if err := platform.DeleteByIDs(s.pick(r), "chain_eqb_flow", ids); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// flowLaunchFlow = ChainEqbFlowServiceImpl.launchFlow 前置:企业签章份数 0 → 602;
// 有份数走外部 e签宝 openapi(本地不可达 → 异常 999)
func (s *Service) flowLaunchFlow(w http.ResponseWriter, r *http.Request) {
	d := s.pick(r)
	u := currentUser(r)
	var free, charge sql.NullInt64
	err := d.QueryRow("SELECT sign_free_count, sign_charge_count FROM chain_basic_enterprise WHERE id = ?", u.CompanyID).Scan(&free, &charge)
	if err != nil {
		platform.Unknown(w) // etc NPE / SQL 异常
		return
	}
	if free.Int64+charge.Int64 == 0 {
		platform.Fail(w, 602, "签章剩余份数不足，请联系平台管理员")
		return
	}
	platform.Unknown(w) // ponytail: e签宝 openapi 未接入(SignHelper 外部调用),本地必失败;接密钥后补
}

// eqbExternalUnavail 其余外部集成端点:confirmFlow/selSignUrl/selSignPeople/uploadFile/
// selLanucnInfo/getSignCallBack/v3DownFile——本地无密钥必异常,保 999 语义
func (s *Service) eqbExternalUnavail(w http.ResponseWriter, r *http.Request) {
	platform.Unknown(w) // ponytail: e签宝 openapi 未接入;接真实密钥后逐端点复刻
}

// ── singers CRUD ──

// singersList = allList:带 contractId → selByConId 三表链,否则 allEq
func (s *Service) singersList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	if cid, ok := m["contractId"]; ok && cid != nil {
		rows, err := d.Query("SELECT s.* FROM chain_eqb_singers s, chain_eqb_flow f, chain_contract c WHERE s.sys_flow_id=f.id AND f.eqb_flow_id=c.eqb_flow_id AND c.eqb_flow_id is not null AND c.eqb_flow_id !='' AND c.id = ?", cid)
		if err != nil {
			platform.Unknown(w)
			return
		}
		ms, _ := platform.ScanMaps(rows)
		if ms == nil {
			ms = []map[string]any{}
		}
		platform.Success(w, shapeEqbSingers(ms))
		return
	}
	cond, ok := allEqCond(m, singersCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(d, "chain_eqb_singers", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, shapeEqbSingers(rows))
}

// singersListByPage = mapper:`${searchSql} and CONCAT(id) like #{searchStr} order by create_time desc`
// singersListByPage = Java mapper `and CONCAT(id) like #{searchStr}` 引用了不存在的顶层参数
// (方法参数名是 searchSql/maps/page),MyBatis 绑定异常 → 实测恒 999,与 body 无关(行为保持)
func (s *Service) singersListByPage(w http.ResponseWriter, r *http.Request) {
	platform.Unknown(w)
}

func (s *Service) singersGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	row, err := platform.GetByID(s.pick(r), "chain_eqb_singers", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, shapeEqbSingers([]map[string]any{row})[0])
}

func (s *Service) singersAdd(w http.ResponseWriter, r *http.Request) {
	s.writeAdd(w, r, "chain_eqb_singers", bodyMap(r), singersCols, nil, singersEchoKeys, false, nil)
}

func (s *Service) singersUpdateById(w http.ResponseWriter, r *http.Request) {
	s.writeUpdate(w, r, "chain_eqb_singers", bodyMap(r), singersCols, nil, singersEchoKeys, false)
}

func (s *Service) singersDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, err := s.pick(r).Exec("DELETE FROM chain_eqb_singers WHERE id = ?", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

func (s *Service) singersDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, _ := m["ids"].([]any)
	if err := platform.DeleteByIDs(s.pick(r), "chain_eqb_singers", ids); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// ── files CRUD(fill 实体:Insert/UpdateByID 自动时间)──

func (s *Service) filesList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, filesCols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_eqb_files", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, shapeEqbFiles(rows))
}

// filesListByPage = ChainEqbFilesServiceImpl.listByPage:getSql + startTime/endTime 日期直拼;
// page 缺失 → listNoPage 返回数组(不分页)
func (s *Service) filesListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	searchSql, args := getSqlEq(m, []string{"id", "name", "field", "type", "sysFlowId",
		"eqbFlowId", "createTime", "updateTime", "remark1", "remark2", "remark3", "remark4",
		"remark5"}, nil, map[string]bool{"id": true, "sysFlowId": true, "eqbFlowId": true})
	if v, ok := m["startTime"]; ok && v != nil {
		searchSql += " and DATE_FORMAT(create_time, '%Y-%m-%d')>= '" + str(m, "startTime") + "'"
	}
	if v, ok := m["endTime"]; ok && v != nil {
		searchSql += " and DATE_FORMAT(create_time, '%Y-%m-%d')<= '" + str(m, "endTime") + "'"
	}
	var ss any
	if v, ok := m["searchStr"]; ok && v != nil {
		ss = "%" + str(m, "searchStr") + "%"
	}
	args = append(args, ss)
	if _, ok := m["page"]; !ok {
		rows, err := d.Query("SELECT * FROM chain_eqb_files "+searchSql+" and CONCAT(id) like ? order by create_time desc", args...)
		if err != nil {
			platform.Unknown(w)
			return
		}
		ms, _ := platform.ScanMaps(rows)
		if ms == nil {
			ms = []map[string]any{}
		}
		platform.Success(w, shapeEqbFiles(ms))
		return
	}
	current, size := pageOf(m)
	q := "SELECT * FROM chain_eqb_files " + searchSql + " and CONCAT(id) like ? order by create_time desc"
	countQ := "SELECT count(*) FROM chain_eqb_files " + searchSql + " and CONCAT(id) like ?"
	page, err := paged(d, q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	page["records"] = shapeEqbFiles(page["records"].([]map[string]any))
	platform.Success(w, page)
}

func (s *Service) filesGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	row, err := platform.GetByID(s.pick(r), "chain_eqb_files", m["id"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, shapeEqbFiles([]map[string]any{row})[0])
}

func (s *Service) filesAdd(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format("2006-01-02 15:04:05")
	s.writeAdd(w, r, "chain_eqb_files", bodyMap(r), filesCols, nil, filesEchoKeys, true,
		map[string]any{"createTime": now, "updateTime": now}) // fill 实体:内存时间进响应
}

func (s *Service) filesUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	vals := bodyVals(m, filesCols, nil)
	delete(vals, "id")
	if len(vals) == 0 {
		platform.Unknown(w) // SET 空 → MybatisPlusException
		return
	}
	if _, err := platform.UpdateByID(s.pick(r), "chain_eqb_files", m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil) // Java 该端点 ResultTool.success() 无 data
}

func (s *Service) filesDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, err := s.pick(r).Exec("DELETE FROM chain_eqb_files WHERE id = ?", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

func (s *Service) filesDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, _ := m["ids"].([]any)
	if err := platform.DeleteByIDs(s.pick(r), "chain_eqb_files", ids); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// filesUpdBatch = updateBatchById(fill 实体:逐行刷 update_time)
func (s *Service) filesUpdBatch(w http.ResponseWriter, r *http.Request) {
	var items []any
	_ = json.NewDecoder(r.Body).Decode(&items)
	d := s.pick(r)
	for _, it := range items {
		em, ok := it.(map[string]any)
		if !ok || em["id"] == nil {
			platform.Unknown(w) // MP updateById 空 id 抛异常
			return
		}
		if _, err := platform.UpdateByID(d, "chain_eqb_files", em["id"], pick(em, filesCols...)); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, nil)
}
