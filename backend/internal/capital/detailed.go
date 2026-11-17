// capital - detailed.go /chainCapitalDetailed 收付款单(主+Cg/Xs 变体,共用 service)。
// Java:ChainCapitalDetailedController + ChainCapitalDetailedServiceImpl + ChainCapitalDetailedMapper(core)。
// 实体无 fill:insert 走库默认 create_time,update_time 恒不刷新;响应=add/updateById 内存实体回显。
package capital

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

// entDetailed chain_capital_detailed 实体全列(载体键附后)
var entDetailed = entShape{
	table: "chain_capital_detailed",
	cols: []string{"id", "im_batch_number", "system_code", "company_id", "account_name", "account",
		"trans_time", "money", "opposite_name", "opposite_account", "abs_content", "purpose",
		"detype", "khxje", "sjhxje", "hx_status", "relation_flag", "bank_name", "hx_code",
		"create_time", "update_time", "remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"money", "khxje", "sjhxje"},
	bools:    []string{"relation_flag"},
	tTimes:   []string{"update_time"}, // updateTime 无 @JsonFormat → ISO T
	dateOnly: []string{"trans_time"},  // Java LocalDate
	carriers: []string{"searchStr", "page", "ids", "sments", "filelist", "timeType", "khxFlag",
		"startTime", "endTime", "contractId"},
	primFalse: []string{"khxFlag"}, // 原生 boolean
}

// 时间绑定列:createTime @JsonFormat 空格;transTime/startTime/endTime LocalDate(ISO);
// updateTime LocalDateTime 无注解(ISO-T)
var dtDetailedSpace = []string{"create_time"}
var dtDetailedDate = []string{"trans_time", "start_time", "end_time"}
var dtDetailedT = []string{"update_time"}

// searchStrLike listByPage 族 searchStr 条件值:service 拼 "%"+getSearchStr()+"%",
// null → 字面 "%null%"(字符串拼接);空串 → "%%"。返回 (参数值, nullLike 与否统一为参数化)。
func searchStrLike(m map[string]any) string {
	if v, ok := m["searchStr"]; !ok || v == nil {
		return "%null%"
	}
	return "%" + str(m, "searchStr") + "%"
}

// ── 路由分发 ──

func (s *Service) handleDetailed(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	// 主 + Cg/Xs 共用:去前缀后按动作分发(Cg/Xs 无 cancelBalancing/getPtSumPrice,但路径注册表已区分)
	switch {
	case strings.HasSuffix(p, "/list"):
		s.dtList(w, r)
	case strings.HasSuffix(p, "/listByPage"):
		s.dtListByPage(w, r)
	case strings.HasSuffix(p, "/getById"):
		s.dtGetById(w, r)
	case strings.HasSuffix(p, "/add"):
		s.dtAdd(w, r)
	case strings.HasSuffix(p, "/deleteById"):
		s.dtDeleteById(w, r)
	case strings.HasSuffix(p, "/updateById"):
		s.dtUpdateById(w, r)
	case strings.HasSuffix(p, "/saveBatchDetails"):
		s.dtSaveBatch(w, r)
	case strings.HasSuffix(p, "/importAlone"):
		s.dtImportAlone(w, r)
	case strings.HasSuffix(p, "/importExcelSk"), strings.HasSuffix(p, "/importExcelFk"):
		s.dtImportExcel(w, r, p)
	case strings.HasSuffix(p, "/exportExcel"):
		s.dtExportExcel(w, r)
	case strings.HasSuffix(p, "/delBatch"):
		s.dtDelBatch(w, r)
	case strings.HasSuffix(p, "/balancingAccounts"):
		s.dtBalancingAccounts(w, r)
	case strings.HasSuffix(p, "/cancelBalancing"):
		s.dtCancelBalancing(w, r)
	case strings.HasSuffix(p, "/getSystemCode"):
		s.dtGetSystemCode(w, r)
	case strings.HasSuffix(p, "/getSments"):
		s.dtGetSments(w, r)
	case strings.HasSuffix(p, "/getSumPrice"):
		s.dtGetSumPrice(w, r)
	case strings.HasSuffix(p, "/getPtSumPrice"):
		s.dtGetPtSumPrice(w, r)
	case strings.HasSuffix(p, "/countDetail"):
		s.dtCountDetail(w, r)
	case strings.HasSuffix(p, "/getBanks"):
		s.dtGetBanks(w, r)
	case strings.HasSuffix(p, "/getDetails"):
		s.dtGetDetails(w, r)
	case strings.HasSuffix(p, "/countDetail"):
		s.dtCountDetail(w, r)
	default:
		platform.Unknown(w)
	}
}

// ── list:allEq(maps,false) + orderByDesc(trans_time) ──

func (s *Service) dtList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, entDetailed.cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	q := "SELECT " + entDetailed.selectList() + " FROM chain_capital_detailed"
	var args []any
	if len(cond) > 0 {
		var conds []string
		for _, k := range platform.SortedKeys(cond) {
			conds = append(conds, k+" = ?")
			args = append(args, cond[k])
		}
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	q += " ORDER BY trans_time DESC"
	rows, err := d.Query(q, args...)
	if err != nil {
		platform.Unknown(w)
		return
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if ms == nil {
		ms = []map[string]any{}
	}
	platform.Success(w, entDetailed.rows(ms))
}

// dtListByPageSearch 手工 searchSql(Java #{maps.x} 参数化 → Go ? 占位,值保真引号不破 SQL)
func dtListByPageSearch(m map[string]any) (string, []any, bool) {
	where := "where 1=1 "
	var args []any
	if v, ok := m["detype"]; ok && v != nil {
		where += " and detype = ?"
		args = append(args, str(m, "detype"))
	}
	if bv, ok := m["khxFlag"]; ok && bv == true {
		where += " and (hx_status = 1 or hx_status = 2)"
	}
	if v, ok := m["hxCode"]; ok && v != nil {
		where += " and hx_code like ?"
		args = append(args, "%"+str(m, "hxCode")+"%")
	}
	if v, ok := m["oppositeName"]; ok && v != nil {
		where += " and opposite_name like ?"
		args = append(args, "%"+str(m, "oppositeName")+"%")
	}
	if v, ok := m["transTime"]; ok && v != nil {
		if strings.TrimSpace(str(m, "transTime")) == "" {
			// 空串绑 null,不拼
		} else {
			t, ok2 := parseLocalDate(str(m, "transTime"))
			if !ok2 {
				return "", nil, false // Jackson LocalDate 绑定失败 → 999
			}
			where += " and trans_time like ?"
			args = append(args, "%"+t.Format("2006-01-02")+"%")
		}
	}
	if v, ok := m["oppositeAccount"]; ok && v != nil {
		where += " and opposite_account like ?"
		args = append(args, "%"+str(m, "oppositeAccount")+"%")
	}
	if v, ok := m["bankName"]; ok && v != nil {
		where += " and bank_name like ?"
		args = append(args, "%"+str(m, "bankName")+"%")
	}
	if v, ok := m["startTime"]; ok && v != nil && strings.TrimSpace(str(m, "startTime")) != "" {
		t, ok2 := parseLocalDate(str(m, "startTime"))
		if !ok2 {
			return "", nil, false
		}
		where += " and DATE_FORMAT(trans_time, '%Y-%m-%d') >= ?"
		args = append(args, t.Format("2006-01-02"))
	}
	if v, ok := m["endTime"]; ok && v != nil && strings.TrimSpace(str(m, "endTime")) != "" {
		t, ok2 := parseLocalDate(str(m, "endTime"))
		if !ok2 {
			return "", nil, false
		}
		where += " and DATE_FORMAT(trans_time, '%Y-%m-%d') <= ?"
		args = append(args, t.Format("2006-01-02"))
	}
	return where, args, true
}

// dtListByPage mapper:显式列(IFNULL(khxje,money) 双列)+CONCAT 8 列 like,order by trans_time desc
func (s *Service) dtListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, dtDetailedSpace, dtDetailedDate, dtDetailedT) {
		platform.Unknown(w)
		return
	}
	searchSql, args, ok := dtListByPageSearch(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	current, size := pageOf(m)
	where := searchSql + " and CONCAT(IFNULL(system_code, ''), IFNULL(account_name, ''), IFNULL(account, ''), IFNULL(trans_time, ''), IFNULL(money, ''), IFNULL(opposite_name, ''), IFNULL(opposite_account, ''), IFNULL(abs_content, '') ) like ? order by trans_time desc"
	args = append(args, searchStrLike(m))
	q := "select id,system_code,company_id,account_name,account,trans_time,money,opposite_name,opposite_account,abs_content,purpose,detype,IFNULL(khxje, money) khxje,IFNULL(khxje, money) sjhxje,hx_status,relation_flag,bank_name,hx_code,create_time,update_time,remark1,remark2,remark3,remark4,remark5 from chain_capital_detailed " + where
	countQ := "SELECT count(1) FROM chain_capital_detailed " + searchSql + " and CONCAT(IFNULL(system_code, ''), IFNULL(account_name, ''), IFNULL(account, ''), IFNULL(trans_time, ''), IFNULL(money, ''), IFNULL(opposite_name, ''), IFNULL(opposite_account, ''), IFNULL(abs_content, '') ) like ?"
	ip, err := paged(s.pick(r), q, countQ, args, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entDetailed.rows(recs)
	platform.Success(w, ip)
}

func (s *Service) dtGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	var row map[string]any
	if id, ok := int64Of(m, "id"); ok {
		var err error
		row, err = platform.GetByID(s.pick(r), "chain_capital_detailed", id)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, entDetailed.rows([]map[string]any{row})[0])
}

// dtAdd saveDetailed:save(NoFill)+filelist+日志;响应内存实体回显
func (s *Service) dtAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, dtDetailedSpace, dtDetailedDate, dtDetailedT) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	vals := pick(m, entDetailed.cols...)
	delete(vals, "id") // AUTO:忽略 body id
	delete(vals, "create_time")
	delete(vals, "update_time")
	id, err := platform.InsertNoFill(d, "chain_capital_detailed", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	updateFilelist(d, m, id, u)
	logOp(d, u, "新增流水\n"+changeLog(nil, m, dtEchoKeys), 3)
	platform.Success(w, echoOf(&entDetailed, m, map[string]any{
		"id": json.Number(strconv.FormatInt(id, 10)),
	}))
}

var dtEchoKeys = entDetailed.echoKeys()

// dtUpdateById updateDetail:getById(日志前值)+updateById(NoFill)+filelist+日志;恒回显
func (s *Service) dtUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, dtDetailedSpace, dtDetailedDate, dtDetailedT) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	var before map[string]any
	id, idok := int64Of(m, "id")
	if idok {
		before, _ = platform.GetByID(d, "chain_capital_detailed", id)
	}
	vals := pick(m, entDetailed.cols...)
	delete(vals, "id")
	delete(vals, "create_time")
	delete(vals, "update_time")
	if !idok {
		platform.Unknown(w) // MP updateById null id → 999
		return
	}
	if _, err := platform.UpdateByIDNoFill(d, "chain_capital_detailed", id, vals); err != nil {
		platform.Unknown(w)
		return
	}
	updateFilelist(d, m, id, u)
	logOp(d, u, "修改流水\n"+changeLog(before, m, dtEchoKeys), 3)
	platform.Success(w, echoOf(&entDetailed, m, nil))
}

// dtDeleteById removeChainCapitalDetailed:removeById+detail_settlement(detail_id)+文件(relation_id)
func (s *Service) dtDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w) // removeById(null) → 999
		return
	}
	d := s.pick(r)
	if err := platform.DeleteByID(d, "chain_capital_detailed", id); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("DELETE FROM chain_capital_detail_settlement WHERE detail_id = ?", id); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("DELETE FROM chain_common_file WHERE relation_id = ?", id); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为："+str(m, "id"))
}

// dtDelBatch removeChainCapitalDetaileds:手工 for 逐条(同 deleteById 副作用);
// ids 缺失/null → Java for-null NPE → 999;空数组 → "删除成功"
func (s *Service) dtDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, ok := m["ids"].([]any)
	if !ok {
		platform.Unknown(w) // ids null/非数组 → NPE/ClassCastException → 999
		return
	}
	d := s.pick(r)
	for _, it := range ids {
		id, ok := toID(it)
		if !ok {
			platform.Unknown(w)
			return
		}
		if err := platform.DeleteByID(d, "chain_capital_detailed", id); err != nil {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("DELETE FROM chain_capital_detail_settlement WHERE detail_id = ?", id); err != nil {
			platform.Unknown(w)
			return
		}
		if _, err := d.Exec("DELETE FROM chain_common_file WHERE relation_id = ?", id); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "删除成功")
}

func toID(v any) (int64, bool) {
	if n, ok := v.(json.Number); ok {
		i, err := n.Int64()
		return i, err == nil
	}
	if s, ok := v.(string); ok {
		i, err := strconv.ParseInt(s, 10, 64)
		return i, err == nil
	}
	return 0, false
}

// dtSaveBatchDetails saveIrt(List):detype=1→SFXS. xscode、2→SFCG. cgcode 分配;响应 "新增成功"
func (s *Service) dtSaveBatch(w http.ResponseWriter, r *http.Request) {
	list, ok := bodyList(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	cg, ok1 := s.detailedSysCode(d, u, "SFCG.")
	xs, ok2 := s.detailedSysCode(d, u, "SFXS.")
	if !ok1 || !ok2 {
		platform.Unknown(w)
		return
	}
	for _, it := range list {
		row, ok := it.(map[string]any)
		if !ok {
			platform.Unknown(w)
			return
		}
		vals := pick(row, entDetailed.cols...)
		delete(vals, "id")
		delete(vals, "create_time")
		delete(vals, "update_time")
		switch str(row, "detype") {
		case "1":
			vals["system_code"] = xs
			xs, _ = nextSystemCodeInMem(xs)
		case "2":
			vals["system_code"] = cg
			cg, _ = nextSystemCodeInMem(cg)
		}
		if _, err := platform.InsertNoFill(d, "chain_capital_detailed", vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "新增成功")
}

// detailedSysCode 对照 DetailedService.getSystemCode:前缀+企业码+日期定位编号
func (s *Service) detailedSysCode(d *sql.DB, u *auth.CurrentUser, idt string) (string, bool) {
	like := idt + u.EnterpriseCode + strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", "")
	return systemCodeNext(d, "chain_capital_detailed", like, "1")
}
