// logistics - waybill.go /chainLogistics 运单控制器(统计 countBuses/countMonBuses 在 dashboard)。
// Java:ChainLogisticsController + ChainLogisticsServiceImpl + ChainLogisticsMapper(core)
package logistics

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// entLo chain_logistics 实体(fill:createTime INSERT + updateTime INSERT_UPDATE → platform.Insert)
var entLo = entShape{
	table: "chain_logistics",
	cols: []string{"id", "im_batch_number", "batch_number", "company_id", "order_id", "system_code",
		"bill_number", "trans_time", "car_number", "loading_address", "unloading_point", "subject_matter",
		"freight_unit", "freight_amount", "original_tonnage", "tonnage_paid", "pay_money", "carriage_platform",
		"file_ids", "lo_create_time", "loading_province", "unloading_province", "loading_area", "unloading_area",
		"loading_pro_lnglog", "unloading_pro_lnglog", "loading_lnglog", "unloading_lnglog", "loading_lng",
		"unloading_lng", "invoicing_status", "invoicing_amount", "invoicing_tonnage", "create_time",
		"update_time", "remark1", "remark2", "remark3", "remark4", "remark5"},
	decimals: []string{"freight_unit", "freight_amount", "original_tonnage", "tonnage_paid", "pay_money",
		"invoicing_amount", "invoicing_tonnage"},
	tTimes:     []string{"trans_time", "lo_create_time", "update_time"},
	carriers:   []string{"searchStr", "page", "startTime", "endTime", "ids", "list", "filelist", "notBindOrder"},
	fillInsert: true,
}

// 时间绑定列(@JsonFormat 空格列 / 无注解 T 列)
var loSpaceDT = []string{"create_time"}
var loTDT = []string{"trans_time", "lo_create_time", "update_time"}

// has NonString 空串判定:Jackson 对 Integer/LocalDate(Time)/Boolean 字段把 "" 绑为 null
// (8765 实测:carriagePlatform:"" 不拼条件);String 字段空串仍参与 like '%%'
func dtof(m map[string]any, k string) (any, bool) {
	v, ok := m[k]
	if !ok || v == nil {
		return nil, false
	}
	if s, isStr := v.(string); isStr && strings.TrimSpace(s) == "" {
		return nil, false
	}
	return v, true
}

// loListByPageSearch Java listByPage 手工拼的 searchSql(值原样内插,注入语义保留)
func loListByPageSearch(m map[string]any) string {
	sqlStr := "where 1=1 "
	if _, ok := dtof(m, "loadingAddress"); ok {
		sqlStr += " and loading_address like '%" + sqlEscape(str(m, "loadingAddress")) + "%'"
	}
	if _, ok := dtof(m, "unloadingPoint"); ok {
		sqlStr += " and unloading_point like '%" + sqlEscape(str(m, "unloadingPoint")) + "%'"
	}
	if _, ok := dtof(m, "billNumber"); ok {
		sqlStr += " and bill_number like '%" + sqlEscape(str(m, "billNumber")) + "%'"
	}
	if _, ok := dtof(m, "transTime"); ok {
		t, err := parseLocalDateTime(str(m, "transTime"))
		if err == nil {
			sqlStr += " and DATE_FORMAT(trans_time, '%Y-%m-%d') = '" + t.Format("2006-01-02") + "'"
		} else {
			sqlStr += "\nBAD_BIND" // Jackson 绑定失败 → 999
		}
	}
	if _, ok := dtof(m, "startTime"); ok {
		t, ok2 := parseLocalDate(str(m, "startTime"))
		if !ok2 {
			return "\nBAD_BIND"
		}
		sqlStr += " and DATE_FORMAT(trans_time, '%Y-%m-%d') >= '" + t.Format("2006-01-02") + "'"
	}
	if _, ok := dtof(m, "endTime"); ok {
		t, ok2 := parseLocalDate(str(m, "endTime"))
		if !ok2 {
			return "\nBAD_BIND"
		}
		sqlStr += " and DATE_FORMAT(trans_time, '%Y-%m-%d') <= '" + t.Format("2006-01-02") + "'"
	}
	if _, ok := dtof(m, "carriagePlatform"); ok {
		sqlStr += " and carriage_platform = " + sqlString(str(m, "carriagePlatform"))
	}
	if _, ok := dtof(m, "notBindOrder"); ok {
		sqlStr += " and (order_id is null or order_id='' or order_id=" + sqlString(str(m, "orderId")) + ")"
	} else if _, ok := dtof(m, "orderId"); ok {
		sqlStr += " and order_id =" + sqlString(str(m, "orderId"))
	}
	return sqlStr
}

// loListByPage POST /chainLogistics/listByPage。mapper:select * from chain_logistics
// ${searchSql} and CONCAT(...) like '%${searchStr}%' order by create_time desc
func (s *Service) loListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["page"].(map[string]any); !ok {
		platform.Unknown(w) // Java mapper requires a Page object.
		return
	}
	d := s.pick(r)
	searchSql := loListByPageSearch(m)
	if strings.Contains(searchSql, "BAD_BIND") {
		platform.Unknown(w)
		return
	}
	current, size := pageOf(m)
	where := searchSql + " and CONCAT(IFNULL(bill_number, ''), IFNULL(trans_time,''), IFNULL(car_number, ''), IFNULL(loading_address, ''), IFNULL(unloading_point, ''), IFNULL(subject_matter, ''), IFNULL(freight_unit, ''), IFNULL(freight_amount, 0), IFNULL(original_tonnage, ''), IFNULL(tonnage_paid, ''), IFNULL(pay_money, ''), IFNULL(carriage_platform, '')) like '%" + searchStrOf(m) + "%'"
	q := "SELECT * FROM chain_logistics " + where + " ORDER BY create_time DESC"
	countQ := "SELECT count(1) FROM chain_logistics " + where
	ip, err := paged(d, q, countQ, nil, current, size)
	if err != nil {
		platform.Unknown(w)
		return
	}
	recs, _ := ip["records"].([]map[string]any)
	ip["records"] = entLo.rows(recs)
	platform.Success(w, ip)
}

// loList POST /chainLogistics/list = allEq(maps,false)
func (s *Service) loList(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	cond, ok := allEqCond(m, entLo.cols)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows, err := platform.ListWhere(s.pick(r), "chain_logistics", cond)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if rows == nil {
		rows = []map[string]any{}
	}
	platform.Success(w, entLo.rows(rows))
}

// loGetById getById:MP selectById(查不到 → success(null))
func (s *Service) loGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	var row map[string]any
	if id, ok := int64Of(m, "id"); ok {
		var err error
		row, err = platform.GetByID(s.pick(r), "chain_logistics", id)
		if err != nil {
			platform.Unknown(w)
			return
		}
	}
	if row == nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, entLo.rows([]map[string]any{row})[0])
}

// loAdd save:fill 实体(Insert 自动 create/update 时间),响应=内存实体回显
func (s *Service) loAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, loSpaceDT, loTDT) {
		platform.Unknown(w) // Jackson 绑定失败
		return
	}
	d := s.pick(r)
	vals := pick(m, entLo.cols...)
	delete(vals, "id") // @TableId(IdType.AUTO):body id 忽略走自增
	id, err := platform.Insert(d, "chain_logistics", vals)
	if err != nil {
		platform.Unknown(w)
		return
	}
	shapeFilelist(m)
	platform.Success(w, platform.EntityEcho(m, entLo.echoKeys(), map[string]any{
		"id":         json.Number(strconv.FormatInt(id, 10)),
		"createTime": nowSpace(),
		"updateTime": nowTNanos(),
	}))
}

// loDeleteById removeById + success("删除成功，id为：X")
func (s *Service) loDeleteById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	id, ok := int64Of(m, "id")
	if ok {
		if err := platform.DeleteByID(d, "chain_logistics", id); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, "删除成功，id为："+strconv.FormatInt(id, 10))
		return
	}
	// id 缺失:MP deleteById(null) 按 id=null 条件删除 0 行,消息拼接字面 null
	if _, err := d.Exec("DELETE FROM chain_logistics WHERE id = NULL"); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功，id为：null")
}

// loUpdateById updateChainLogistics:updateById(fill)+filelist(chain_common_file 批更),
// controller 恒回显内存实体(无 try-catch,MP 异常 → 999)
func (s *Service) loUpdateById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !validateDT(m, loSpaceDT, loTDT) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	id, ok := int64Of(m, "id")
	if !ok {
		platform.Unknown(w) // MP updateById 无 id → MybatisPlusException
		return
	}
	vals := pick(m, entLo.cols...)
	delete(vals, "id")
	if _, err := platform.UpdateByID(d, "chain_logistics", id, vals); err != nil {
		platform.Unknown(w)
		return
	}
	s.updateFilelist(d, m, id, u)
	shapeFilelistSvc(m, id, u.UserID)
	platform.Success(w, platform.EntityEcho(m, entLo.echoKeys(), map[string]any{"updateTime": nowTNanos()}))
}

// updateFilelist 对照 feign updBthById 前置赋值:fileType/ownerType/ownerId/relationId/
// relationType/effective 固定后按 id 批更 chain_common_file(无 fill 实体)
func (s *Service) updateFilelist(d *sql.DB, m map[string]any, id int64, u *auth.CurrentUser) {
	files, ok := m["filelist"].([]any)
	if !ok {
		return
	}
	for _, it := range files {
		f, ok := it.(map[string]any)
		if !ok || f["id"] == nil {
			continue
		}
		_, _ = platform.UpdateByIDNoFill(d, "chain_common_file", f["id"], map[string]any{
			"file_type": "1", "owner_type": "1", "owner_id": strconv.FormatInt(u.UserID, 10),
			"relation_id": json.Number(strconv.FormatInt(id, 10)), "relation_type": 1, "effective": true,
		})
	}
}

// loDelBatch removeByIds:null/空集合 → MP 短路 0 行,照 8765 实测 200 "删除成功"
func (s *Service) loDelBatch(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	ids, ok := m["ids"].([]any)
	if !ok || len(ids) == 0 {
		platform.Success(w, "删除成功")
		return
	}
	numIds := make([]any, 0, len(ids))
	for _, it := range ids {
		n, ok := it.(json.Number)
		if !ok {
			platform.Unknown(w)
			return
		}
		numIds = append(numIds, n.String())
	}
	if err := platform.DeleteByIDs(s.pick(r), "chain_logistics", numIds); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, "删除成功")
}

// loSaveBatch saveBatchLogistics:系统编号 WLGD 顺序赋值后 saveBatch,响应固定 "新增成功"
func (s *Service) loSaveBatch(w http.ResponseWriter, r *http.Request) {
	s.saveBatchWithCode(w, r, "chain_logistics", entLo, "WLYD.")
}

// saveBatchWithCode 对照 ChainLogistics/Weighing saveIrt(List, request):
// getSystemCode(jwt) 起头,逐条 setSystemCode + 递增,saveBatch
func (s *Service) saveBatchWithCode(w http.ResponseWriter, r *http.Request, table string, ent entShape, idt string) {
	items, ok := bodyList(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	u := currentUser(r)
	curcode, ok2 := nextSystemCode(d, table, idt+u.EnterpriseCode+time.Now().Format("20060102"), "1")
	if !ok2 {
		platform.Unknown(w)
		return
	}
	for _, it := range items {
		e, ok := it.(map[string]any)
		if !ok {
			continue
		}
		if table == "chain_logistics" && !validateDT(e, loSpaceDT, loTDT) {
			platform.Unknown(w)
			return
		}
		if table == "chain_logistics_weighing_list" && !validateDT(e, wlSpaceDT, wlTDT) {
			platform.Unknown(w)
			return
		}
		e["systemCode"] = curcode
		next, ok3 := bumpSystemCode(curcode)
		if !ok3 {
			platform.Unknown(w)
			return
		}
		curcode = next
		vals := pick(e, ent.cols...)
		delete(vals, "id")
		if ent.fillInsert {
			if _, err := platform.Insert(d, table, vals); err != nil {
				platform.Unknown(w)
				return
			}
		} else {
			if _, err := platform.InsertNoFill(d, table, vals); err != nil {
				platform.Unknown(w)
				return
			}
		}
	}
	platform.Success(w, "新增成功")
}

// loUpdTransLocation 自动匹配收发货地定位(腾讯地图 place/search,Key 与 Java 同源硬编码)。
// 任一环节失败 → catch → "更新运单数据0条"。外部依赖:apis.map.qq.com 可达性。
func (s *Service) loUpdTransLocation(w http.ResponseWriter, r *http.Request) {
	d := s.pick(r)
	count, err := s.updTransLocation(d)
	if err != nil {
		platform.Success(w, "更新运单数据0条")
		return
	}
	platform.Success(w, "更新运单数据"+strconv.Itoa(count)+"条")
}

const tencentMapKey = "CM7BZ-7FEWP-4H2DF-VYNXK-3EQJJ-2YBP6"

func (s *Service) updTransLocation(d *sql.DB) (int, error) {
	rows, err := d.Query("SELECT DISTINCT loading_address, unloading_point FROM chain_logistics WHERE (loading_area IS NULL OR unloading_area IS NULL) AND (loading_address IS NOT NULL AND unloading_point IS NOT NULL)")
	if err != nil {
		return 0, err
	}
	ms, err := platform.ScanMaps(rows)
	if err != nil {
		return 0, err
	}
	maps := map[string]bool{}
	count := 0
	for _, lo := range ms {
		loadingAddress := btsStr(lo["loading_address"])
		if parts := strings.Split(loadingAddress, "-"); len(parts) > 1 {
			loadingAddress = parts[1]
		}
		if i := strings.Index(loadingAddress, "#"); i >= 0 {
			loadingAddress = loadingAddress[:i]
		}
		if maps[loadingAddress] {
			continue
		}
		loadingArea, loadingLnglog, err := tencentSearch(loadingAddress)
		if err != nil {
			return 0, err
		}
		time.Sleep(500 * time.Millisecond)
		unloadingPoint := btsStr(lo["unloading_point"])
		if parts := strings.Split(unloadingPoint, "-"); len(parts) > 1 {
			unloadingPoint = parts[1]
		}
		if i := strings.Index(unloadingPoint, "#"); i >= 0 {
			unloadingPoint = unloadingPoint[:i]
		}
		if maps[unloadingPoint] {
			continue
		}
		unloadingArea, unloadingLnglog, err := tencentSearch(unloadingPoint)
		if err != nil {
			return 0, err
		}
		res, err := d.Exec("UPDATE chain_logistics SET loading_area=?, loading_lnglog=?, unloading_area=?, unloading_lnglog=? WHERE loading_address=? AND unloading_point=?",
			loadingArea, loadingLnglog, unloadingArea, unloadingLnglog, btsStr(lo["loading_address"]), btsStr(lo["unloading_point"]))
		if err != nil {
			return 0, err
		}
		n, _ := res.RowsAffected()
		count += int(n)
	}
	return count, nil
}

// tencentSearch 对照 updTransList:keyword 搜索 → ad_info.district → 再搜 district 取 location。
// count==0 时 Java location.get("lng") NPE → 此处以 error 表达(同收 0 条结果)
func tencentSearch(address string) (area, location string, err error) {
	res, err := mapSearch(address)
	if err != nil {
		return "", "", err
	}
	// bodyMap 走 UseNumber,count 是 json.Number:原先按 float64 断言恒得 0,功能永远"更新 0 条"
	if num(res["count"]) <= 0 { // 错误响应无 count 键 → Java getCount() 拆箱 NPE → 同 0 条
		return "", "", fmt.Errorf("无结果")
	}
	data, _ := res["data"].([]any) // 第三方返回异常结构时不 panic
	if len(data) == 0 {
		return "", "", fmt.Errorf("无结果")
	}
	first, _ := data[0].(map[string]any)
	if ad, ok := first["ad_info"].(map[string]any); ok && ad["district"] != nil {
		area, _ = ad["district"].(string)
		res2, err := mapSearch(area)
		if err != nil {
			return "", "", err
		}
		if c, _ := res2["count"].(float64); c > 0 {
			d2 := res2["data"].([]any)
			f2, _ := d2[0].(map[string]any)
			if loc, ok := f2["location"].(map[string]any); ok {
				location = fmt.Sprintf("%v,%v", loc["lng"], loc["lat"])
			}
		}
		return area, location, nil
	}
	area, _ = first["title"].(string)
	if loc, ok := first["location"].(map[string]any); ok {
		location = fmt.Sprintf("%v,%v", loc["lng"], loc["lat"])
	}
	return area, location, nil
}

func mapSearch(keyword string) (map[string]any, error) {
	url := "https://apis.map.qq.com/ws/place/v1/search?boundary=region(%E4%B8%AD%E5%9B%BD%EF%BC%8C0)&page_size=10&page_index=1&keyword=" + urlQueryEscape(keyword) + "&key=" + tencentMapKey
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
