// product_list.go POST /busProduct/listByPage(对照 BusProductServiceImpl.listByPage 159-423 行 +
// BusProductMapper 的 listByPage/listNoPage/listByFenCi/listNoPFenCi/relListByPage/selTags、
// BusUserClickProMapper.selLastClickPro;SQL 原文逐字照抄,含 Java 端错列 bug)。
//
// HanLP 分词由同版本 portable-1.8.4 模型原生移植，离线黄金用原服务器 JAR 对照。
// 数字词的"硫/挥发/水"分支仍照抄 Java(可解析数字不含汉字，属于原实现的不可达分支)。
//
// 其他登记差异:
//  1. 推荐分支 page 缺失:Java 先跑 relListNoPage(结果丢弃)再 relListByPage(null) 连环异常 → 999;
//     Go 在分支入口直接 999(只读等价,响应逐字节一致)。
//  2. 匿名验签请求(UserID=0)视同 Java jwtInfo=null:maps.userId=0、推荐分支跳过、
//     warehouseType=1 → NPE 对应 999(authGuard 已拦未登录,正常登录流不触达)。
//  3. JSON 数字字面小数位:Java BigDecimal 保 "550.00",Go float64 得 "550"(仓库既有近似,
//     前端以字符串传数可保字面,bindVal 已支持)。
//  4. MP count 优化(optimizeCountSql)以 SELECT COUNT(*) FROM (原SQL去 order by) cnt 近似
//     (仓库既有约定,见 template.go stripOrderBy)。
package warehouse

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		mux.HandleFunc("POST /saas-warehouse/busProduct/listByPage", authGuard(s.bpListByPage))
	})
}

// ── mapper SQL 原文(逐字照抄 BusProductMapper.java / BusUserClickProMapper.java)──

const bpUnstr25 = "fstProId,bigProId,threeProId,minCalorificValue,maxCalorificValue,minWaterContent,maxWaterContent,minAshContent,maxAshContent,minStd,maxStd,minVolatiles,maxVolatiles,minCsr,maxCsr,minBondIndex,maxBondIndex,minY,maxY,isPlatformCertification,isRealCheck,provinceCode,cityCode,areaCode,groupId"

const bpSqlListByPage = `select false as isDivider, false as loading, (SELECT count(1) from chain_common_file where association_id=p.association_id and f_type=4 and is_compressed_file=0) mzkfyq, (SELECT count(1) from bus_pro_log where operate_obj=4 and operate_type=7 and pro_id=p.id and target_user_id=p.create_user_id) visitedTotal, (SELECT count(DISTINCT user_id) from bus_pro_log where operate_obj=4 and operate_type=7 and pro_id=p.id and target_user_id=p.create_user_id) vtPeopCount, (SELECT count(1) from bus_pro_log where operate_obj=4 and operate_type=4 and pro_id=p.id and target_user_id=p.create_user_id) shareTotal, (SELECT count(1) from bus_pro_log where operate_obj=4 and operate_type=5 and pro_id=p.id and target_user_id=p.create_user_id) fwdTotal, (SELECT count(1) from bus_pro_log where operate_obj=4 and operate_type=6 and pro_id=p.id and target_user_id=p.create_user_id) interestTotal, (select warehouse_video from bus_warehouse where id=p.warehouse_id) warehouseVideo, (select warehouse_img from bus_warehouse where id=p.warehouse_id) warehouseImg, (SELECT IFNULL(count(1), 0) FROM bus_pro_collection WHERE pro_id = p.id and user_id=#{maps.userId}) collCount, (SELECT qua_cert_score+bas_info_score+coal_verify_score+upd_share_score+report_score+qualified_score+upd_price_score from chain_basic_enterprise where id=p.company_id) sumScore, (select is_show_etc_name  from bus_member_benefits b, chain_auth_user uu where b.member_type=uu.member_type and uu.id=p.create_user_id) isShowEtcName, (select is_show_etc_phone  from bus_member_benefits b, chain_auth_user uu where b.member_type=uu.member_type and uu.id=p.create_user_id) isShowEtcPhone, (select ifnull(count(1), 0) from bus_verify_info where pro_id=p.id) isVerify, p.* from bus_product p ${searchSql} and CONCAT(IFNULL(enterprise_name, ''),IFNULL(warehouse_nb, ''),IFNULL(region, ''),IFNULL(product_name, ''),IFNULL(resource_number, ''),IFNULL(send_addres, ''),IFNULL(fst_name, ''),IFNULL(sec_name, ''),IFNULL(th_name, ''),IFNULL(warse_code, ''),IFNULL(coal_code, ''),IFNULL(create_time, ''),IFNULL(update_time, ''),IFNULL(create_time, ''),IFNULL(update_time, ''),IFNULL(business_name, ''),IFNULL(area_buser_name, ''),IFNULL(coal_rel_name, ''),IFNULL(create_phone, '')) like '%${maps.searchStr}%' order by status asc, update_time desc`

const bpSqlListNoPage = `select (select warehouse_video from bus_warehouse where id=p.warehouse_id) warehouseVideo, (select warehouse_img from bus_warehouse where id=p.warehouse_id) warehouseImg, (SELECT IFNULL(count(1), 0) FROM bus_pro_collection WHERE pro_id = p.id and user_id=#{maps.userId}) collCount,(select is_show_etc_name  from bus_member_benefits b, chain_auth_user uu where b.member_type=uu.member_type and uu.id=p.create_user_id) isShowEtcName, (select is_show_etc_phone  from bus_member_benefits b, chain_auth_user uu where b.member_type=uu.member_type and uu.id=p.create_user_id) isShowEtcPhone,p.* from bus_product p ${searchSql} and CONCAT(IFNULL(enterprise_name, ''),IFNULL(warehouse_nb, ''),IFNULL(region, ''),IFNULL(product_name, ''),IFNULL(resource_number, ''),IFNULL(send_addres, ''),IFNULL(fst_name, ''),IFNULL(sec_name, ''),IFNULL(th_name, ''),IFNULL(warse_code, ''),IFNULL(coal_code, ''),IFNULL(create_time, ''),IFNULL(update_time, ''),IFNULL(create_time, ''),IFNULL(update_time, ''),IFNULL(business_name, ''),IFNULL(area_buser_name, ''),IFNULL(coal_rel_name, ''),IFNULL(create_phone, '')) like '%${maps.searchStr}%' order by status asc, update_time desc`

const bpSqlListByFenCi = `select (SELECT qua_cert_score+bas_info_score+coal_verify_score+upd_share_score+report_score+qualified_score+upd_price_score from chain_basic_enterprise where id=p.company_id) sumScore, (select warehouse_video from bus_warehouse where id=p.warehouse_id) warehouseVideo, (select warehouse_img from bus_warehouse where id=p.warehouse_id) warehouseImg, (SELECT IFNULL(count(1), 0) FROM bus_pro_collection WHERE pro_id = p.id and user_id=#{maps.userId}) collCount,p.* from bus_product p ${searchSql} order by status asc, update_time desc`

const bpSqlListNoPFenCi = `select (select warehouse_video from bus_warehouse where id=p.warehouse_id) warehouseVideo, (select warehouse_img from bus_warehouse where id=p.warehouse_id) warehouseImg, (SELECT IFNULL(count(1), 0) FROM bus_pro_collection WHERE pro_id = p.id and user_id=#{maps.userId}) collCount,p.* from bus_product p ${searchSql} order by status asc, update_time desc`

const bpSqlRelListByPage = `select p.*, (SELECT qua_cert_score+bas_info_score+coal_verify_score+upd_share_score+report_score+qualified_score+upd_price_score from chain_basic_enterprise where id=p.company_id) sumScore ,(select is_show_etc_name  from bus_member_benefits b, chain_auth_user uu where b.member_type=uu.member_type and uu.id=p.create_user_id) isShowEtcName, (select is_show_etc_phone  from bus_member_benefits b, chain_auth_user uu where b.member_type=uu.member_type and uu.id=p.create_user_id) isShowEtcPhone from bus_user_click_pro cp RIGHT JOIN bus_product p on cp.three_pro_id=p.three_pro_id and cp.user_id=${maps.userId} ${searchSql} and CONCAT(IFNULL(warehouse_nb, ''),IFNULL(region, ''),IFNULL(product_name, ''),IFNULL(resource_number, ''),IFNULL(send_addres, '')) like #{maps.searchStr} order by cp.is_last_click desc, cp.three_pro_id DESC, p.update_time desc`

const bpSqlSelLastClickPro = `select *, (SELECT qua_cert_score+bas_info_score+coal_verify_score+upd_share_score+report_score+qualified_score+upd_price_score from chain_basic_enterprise where id=p.company_id) sumScore,(select is_show_etc_name from bus_member_benefits b, chain_auth_user uu where b.member_type=uu.member_type and uu.id=p.create_user_id) isShowEtcName from bus_product p, bus_user_click_pro cp where p.id=cp.last_pro_id and status=1 and publish_type=1 and is_last_click=1 and user_id=${userId}`

// ── 小工具 ──

// bpB 校验已过的 body 标量取 Jackson 绑定值;validateBody 通过后不会失败。
// 缺失/JSON null → nil;i/d/D 类型的 "" 也绑 null(照 Jackson);s 类型 "" 保持空串
func bpB(kind string, v any) any {
	b, _ := bindVal(kind, v)
	return b
}

// bpBoolTrue 对照 `null != getX() && getX()`(仅字面 true 生效)
func bpBoolTrue(m map[string]any, k string) bool {
	t, _ := bpB("b", m[k]).(bool)
	return t
}

// bpCount 对照 count(new QueryWrapper<BusProduct>().like(col, v)):
// SELECT COUNT(*) FROM bus_product WHERE (col LIKE '%v%'),v null → 字面 '%null%'
func bpCount(d *sql.DB, col, v string) (int64, error) {
	var n int64
	err := d.QueryRow("SELECT COUNT(*) FROM bus_product WHERE ("+col+" LIKE ?)", "%"+v+"%").Scan(&n)
	return n, err
}

func bpCachedPhone(r *http.Request, uid string) (string, error) {
	hash, err := platform.RedisHash(r.Context(), "userId"+uid)
	if err != nil {
		return "", err
	}
	phone, ok := hash["phone"]
	if !ok {
		return "null", nil
	}
	return phone, nil
}

// bpParseNum 对照 Double.parseDouble(word) 成功(仅数字/小数词;字母词 NaN/Inf 不视为数字)
func bpParseNum(w string) bool {
	hasDigit := false
	for _, c := range w {
		if c >= '0' && c <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false
	}
	_, err := strconv.ParseFloat(w, 64)
	return err == nil
}

// bpKeyOnly 仅保留实体键(selLastClickPro 的 select * 带出 bus_user_click_pro 列:同名列
// MyBatis 自动映射按列序后写覆盖,shapeRows 行为一致;非实体列(userId/lastProId/isLastClick
// 等)Jackson 不序列化,这里删除对齐)
func bpKeyOnly(e *ent, rows []map[string]any) []map[string]any {
	keep := make(map[string]bool, len(e.keys))
	for _, k := range e.keys {
		keep[k] = true
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		for k := range row {
			if !keep[k] {
				delete(row, k)
			}
		}
	}
	return rows
}

// bpIPage 分页执行(MP 3.4.2 语义):count(原SQL去 order by 包一层)→ 窗口;size<0 不分页;
// 失败写 999 返 nil
func (s *Service) bpIPage(w http.ResponseWriter, r *http.Request, q string, args []any, current, size int64) ([]map[string]any, int64) {
	var total int64
	if err := s.pick(r).QueryRow("SELECT COUNT(*) FROM ("+stripOrderBy(q)+") cnt", args...).Scan(&total); err != nil {
		platform.Unknown(w)
		return nil, 0
	}
	if size < 0 { // MP:size<0 不分页
		rows := s.qEntEnt2(w, r, entBusProduct, q, args...)
		if rows == nil {
			return nil, 0
		}
		return rows, total
	}
	rows := s.qEntEnt2(w, r, entBusProduct, q+" LIMIT ? OFFSET ?", append(append([]any{}, args...), size, (current-1)*size)...)
	if rows == nil {
		return nil, 0
	}
	return rows, total
}

// ── 端点 ──

// bpListByPage 对照 listByPage(jwtInfo 三参版):条件按 Java 源顺序拼 searchSql(错误列
// bug 照抄);!isDirect 走近似分词分支;threeProId 空+isClientSearch 走推荐分支;
// page 空 → listNoPage(List),否则 listByPage(IPage);带 @Results 的查询补 tags。
func (s *Service) bpListByPage(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	e := entBusProduct
	if !e.validateBody(m) {
		platform.Unknown(w) // Jackson 绑定失败
		return
	}
	d := s.pick(r)
	uid := currentUser(r).UserID
	uidStr := "0" // jwtInfo null → maps.userId=0
	if uid != 0 {
		uidStr = strconv.FormatInt(uid, 10)
	}
	// searchStr 绑定值与"是否缺失"(Java null 语义)
	searchRaw, searchNil := m["searchStr"], m["searchStr"] == nil
	ss := "null" // like(col, null) → '%null%' 字面
	var searchArg any
	if !searchNil {
		ss = bpB("s", searchRaw).(string)
		searchArg = "%" + ss + "%" // getMaps %v% 包裹
	}
	hasPage := false
	current, size := int64(1), int64(10)
	if p, ok := m["page"].(map[string]any); ok {
		hasPage = true
		if v, ok2 := p["current"]; ok2 && v != nil {
			current = int64(num(v))
		}
		if v, ok2 := p["size"]; ok2 && v != nil {
			size = int64(num(v))
		}
	}

	// 1. getSql(busProduct, 25 字段 unstr)
	sqlStr, args, ok := e.getSqlUnstr(m, bpUnstr25)
	if !ok {
		platform.Unknown(w)
		return
	}
	// 2. 指标区间(Java 错列 bug 照抄:minWater→max_water_content、minAsh→max_ash_content、
	//    minStd→max_std、minVolatiles→max_volatiles、maxCsr→min_csr、maxBondIndex→min_bond_index、
	//    maxY→min_y,其余正列;数值裸拼)
	for _, rg := range []struct{ key, col, op string }{
		{"minCalorificValue", "min_calorific_value", ">="},
		{"maxCalorificValue", "max_calorific_value", "<="},
		{"minWaterContent", "max_water_content", ">="},
		{"maxWaterContent", "max_water_content", "<="},
		{"minAshContent", "max_ash_content", ">="},
		{"maxAshContent", "max_ash_content", "<="},
		{"minStd", "max_std", ">="},
		{"maxStd", "max_std", "<="},
		{"minVolatiles", "max_volatiles", ">="},
		{"maxVolatiles", "max_volatiles", "<="},
		{"minCsr", "min_csr", ">="},
		{"maxCsr", "min_csr", "<="},
		{"minBondIndex", "min_bond_index", ">="},
		{"maxBondIndex", "min_bond_index", "<="},
		{"minY", "min_y", ">="},
		{"maxY", "min_y", "<="},
	} {
		if v := bpB("d", m[rg.key]); v != nil {
			sqlStr += " and " + rg.col + " is not null and " + rg.col + "!=0 and " + rg.col + " " + rg.op + " " + v.(string)
		}
	}
	// 3. p.fst/big/three_pro_id 裸拼
	for _, pr := range []struct{ key, col string }{
		{"fstProId", "p.fst_pro_id"}, {"bigProId", "p.big_pro_id"}, {"threeProId", "p.three_pro_id"},
	} {
		if v := bpB("i", m[pr.key]); v != nil {
			sqlStr += " and " + pr.col + " = " + v.(string)
		}
	}
	// 4. isRealCheck 1/0;省市县码裸拼(非数字 → SQL 语法错误 → 999,照抄)
	if v := bpB("b", m["isRealCheck"]); v != nil {
		bit := "0"
		if v.(bool) {
			bit = "1"
		}
		sqlStr += " and is_real_check = " + bit
	}
	for _, cd := range []struct{ key, col string }{
		{"provinceCode", "province_code"}, {"cityCode", "city_code"}, {"areaCode", "area_code"},
	} {
		if v := bpB("s", m[cd.key]); v != nil {
			sqlStr += " and " + cd.col + " = " + v.(string)
		}
	}
	// 5. 6 个日期区间(#{maps.X} 参数化,LocalDate 原值绑定)
	for _, dt := range []struct{ key, col, op string }{
		{"createStartTimeDate", "create_time", ">="},
		{"createEndTimeDate", "create_time", "<="},
		{"updateStartTimeDate", "update_time", ">="},
		{"updateEndTimeDate", "update_time", "<="},
		{"facPriceStartTimeDate", "fprice_upd_time", ">="},
		{"facPriceEndTimeDate", "fprice_upd_time", "<="},
	} {
		if v := bpB("D", m[dt.key]); v != nil {
			sqlStr += " and DATE_FORMAT(" + dt.col + ", '%Y-%m-%d') " + dt.op + " ?"
			args = append(args, v.(string))
		}
	}
	// 6. isRelationEtc
	if v := bpB("b", m["isRelationEtc"]); v != nil {
		if v.(bool) {
			sqlStr += " and p.company_id is not null"
		} else {
			sqlStr += " and p.company_id is null"
		}
	}
	// 7. warehouseType 0
	if v := bpB("i", m["warehouseType"]); v != nil && v.(string) == "0" {
		sqlStr += " and status = 1 and (is_certified_etc=1 or ware_identity=4)"
	}
	// 8. 原 Java 直接拼接 Redis phone，未命中保留字面 null；不以数据库补缓存。
	if v := bpB("i", m["warehouseType"]); v != nil && v.(string) == "1" {
		if uid == 0 {
			platform.Unknown(w)
			return
		}
		phone, err := bpCachedPhone(r, uidStr)
		if err != nil {
			platform.Unknown(w)
			return
		}
		sqlStr += " and (create_user_id = " + uidStr + " or business_id= " + uidStr + " or area_buser_id= " + uidStr + " or etc_rel_phone=" + phone + ")"
	}
	// 9. warehouseType 2:otherUserId(null → 字面 null);redis 无记录 → etc_rel_phone=null
	if v := bpB("i", m["warehouseType"]); v != nil && v.(string) == "2" {
		ouid := "null"
		if ov := bpB("i", m["otherUserId"]); ov != nil {
			ouid = ov.(string)
		}
		phone, err := bpCachedPhone(r, ouid)
		if err != nil {
			platform.Unknown(w)
			return
		}
		sqlStr += " and status != 3 and (create_user_id = " + ouid + " or business_id= " + ouid + " or area_buser_id= " + ouid + " or etc_rel_phone=" + phone + ")"
	}
	// 10/11. isUpDown / isPlatformCertification
	if bpBoolTrue(m, "isUpDown") {
		sqlStr += " and (status = 1 or status = 3)"
	}
	if bpBoolTrue(m, "isPlatformCertification") {
		sqlStr += " and is_platform_certification=1"
	}
	// 12. isDirect 判定(5 个 like count 实查;searchStr null 且首 count=0 → contains NPE → 999)
	isDirect := false
	etcN, err := bpCount(d, "enterprise_name", ss)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if etcN > 0 {
		isDirect = true
	} else if searchNil {
		platform.Unknown(w) // searchStr.contains("公司") NPE
		return
	} else if strings.Contains(ss, "公司") {
		isDirect = true
	}
	for _, col := range []string{"th_name", "warehouse_nb", "warse_code", "create_phone"} {
		n, err := bpCount(d, col, ss)
		if err != nil {
			platform.Unknown(w)
			return
		}
		if n > 0 {
			isDirect = true
		}
	}
	// 13. isGroupSearch(groupId null → 字面 null 拼接)
	if v := bpB("b", m["isGroupSearch"]); v != nil {
		if v.(bool) {
			g := "null"
			if gv := bpB("i", m["groupId"]); gv != nil {
				g = gv.(string)
			}
			sqlStr += " and (p.company_id = " + g + " or group_id = " + g + ")"
		} else {
			sqlStr += " and ((p.company_id is null or p.company_id !=420) and (group_id is null or group_id != 420))"
		}
	}
	// 14. isShared/isforwared/isCared/isVisited exists 子查询
	if bpBoolTrue(m, "isShared") {
		sqlStr += " and exists (select 1 from bus_pro_log where operate_obj=4 and operate_type=4 and pro_id=p.id)"
	}
	if bpBoolTrue(m, "isforwared") {
		sqlStr += " and exists (select 1 from bus_pro_log where operate_obj=4 and operate_type=5 and pro_id=p.id)"
	}
	if bpBoolTrue(m, "isCared") {
		sqlStr += " and exists (select 1 from bus_pro_collection where pro_id=p.id)"
	}
	if bpBoolTrue(m, "isVisited") {
		sqlStr += " and exists (select 1 from bus_pro_log where operate_obj=4 and operate_type=7 and pro_id=p.id)"
	}
	// 15. 分词分支(!isDirect;使用服务器同版本 HanLP portable 模型)
	if !isDirect {
		words := hanlpNative(ss)
		hasEtcName, hasCoalName := false, false
		for i, word := range words {
			if !hasEtcName {
				n, err := bpCount(d, "enterprise_name", word)
				if err != nil {
					platform.Unknown(w)
					return
				}
				if n > 0 {
					sqlStr += " and enterprise_name like ?"
					args = append(args, "%"+word+"%")
					hasEtcName = true
					continue
				}
			}
			if !hasCoalName {
				n, err := bpCount(d, "th_name", word)
				if err != nil {
					platform.Unknown(w)
					return
				}
				if n > 0 {
					sqlStr += " and th_name like ?"
					args = append(args, "%"+word+"%")
					hasCoalName = true
					continue
				}
			}
			if bpParseNum(word) { // Java 死分支照抄:parseable 词不含汉字,"硫/挥发/水"不可达
				if strings.Contains(word, "硫") {
					if i+1 >= len(words) {
						platform.Unknown(w) // termList.get(i+1) 越界
						return
					}
					sqlStr += " and max_std =" + words[i+1]
				} else if strings.Contains(word, "挥发") {
					if i+1 >= len(words) {
						platform.Unknown(w)
						return
					}
					sqlStr += " and max_volatiles =" + words[i+1]
				} else if strings.Contains(word, "水") {
					if i+1 >= len(words) {
						platform.Unknown(w)
						return
					}
					sqlStr += " and max_water_content =" + words[i+1]
				} else {
					sqlStr += " and (max_ash_content=" + word + " or max_volatiles=" + word + " or max_std=" + word + " or max_water_content=" + word + " or min_bond_index=" + word + " or min_csr=" + word + " or min_y=" + word + " or lithofacies=" + word + " or clean_coal_recovery_rate=" + word + " or gangue=" + word + ")"
				}
				continue
			}
		}
		if len(words) > 1 {
			if !hasPage {
				rows := s.qEntEnt2(w, r, e, strings.Replace(strings.Replace(bpSqlListNoPFenCi, "${searchSql}", sqlStr, 1), "#{maps.userId}", "?", 1), append([]any{uidStr}, args...)...)
				if rows == nil {
					return
				}
				if err := bpAttachTags(d, rows); err != nil {
					platform.Unknown(w)
					return
				}
				platform.Success(w, rows)
				return
			}
			q := strings.Replace(strings.Replace(bpSqlListByFenCi, "${searchSql}", sqlStr, 1), "#{maps.userId}", "?", 1)
			rows, total := s.bpIPage(w, r, q, append([]any{uidStr}, args...), current, size)
			if rows == nil {
				return
			}
			if err := bpAttachTags(d, rows); err != nil {
				platform.Unknown(w)
				return
			}
			platform.Success(w, ipageOf(rows, total, current, size))
			return
		}
	}
	// 16. 推荐煤种分支(threeProId 空 + jwtInfo 非空 + isClientSearch)
	if m["threeProId"] == nil && uid != 0 && bpBoolTrue(m, "isClientSearch") {
		if !hasPage {
			// Java:relListNoPage(结果丢弃)+ relListByPage(null) 连环异常 → 999(见文件头差异 2)
			platform.Unknown(w)
			return
		}
		lastPros := s.qEntEnt2(w, r, e, strings.Replace(bpSqlSelLastClickPro, "${userId}", uidStr, 1))
		if lastPros == nil {
			return
		}
		lastPros = bpKeyOnly(e, lastPros)
		if len(lastPros) > 0 {
			pid := "null"
			if idv := lastPros[0]["id"]; idv != nil {
				pid = decStr(idv)
			}
			sqlStr += " and p.id != " + pid
		} else {
			size = 10 // page.setSize(10)
		}
		q := strings.Replace(bpSqlRelListByPage, "${maps.userId}", "?", 1)
		q = strings.Replace(q, "${searchSql}", sqlStr, 1)
		q = strings.Replace(q, "#{maps.searchStr}", "?", 1)
		rows, total := s.bpIPage(w, r, q, append(append([]any{uidStr}, args...), searchArg), current, size)
		if rows == nil {
			return
		}
		// relListByPage 无 @Results → 无 tags(照抄 Java)
		platform.Success(w, map[string]any{"lastPros": lastPros, "prolist": ipageOf(rows, total, current, size)})
		return
	}
	// 17. 常规分支
	if !hasPage {
		q := strings.Replace(strings.Replace(bpSqlListNoPage, "${searchSql}", sqlStr, 1), "#{maps.userId}", "?", 1)
		inner := "%null%" // getMaps 已包 %v%,模板再包 '%' → '%v%' 语义;searchStr 缺失为字面 null
		if !searchNil {
			inner = searchArg.(string)
		}
		q = strings.Replace(q, "${maps.searchStr}", inner, 1)
		rows := s.qEntEnt2(w, r, e, q, append([]any{uidStr}, args...)...)
		if rows == nil {
			return
		}
		if err := bpAttachTags(d, rows); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, rows)
		return
	}
	q := strings.Replace(strings.Replace(bpSqlListByPage, "${searchSql}", sqlStr, 1), "#{maps.userId}", "?", 1)
	inner := "%null%"
	if !searchNil {
		inner = searchArg.(string)
	}
	q = strings.Replace(q, "${maps.searchStr}", inner, 1)
	rows, total := s.bpIPage(w, r, q, append([]any{uidStr}, args...), current, size)
	if rows == nil {
		return
	}
	if err := bpAttachTags(d, rows); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, ipageOf(rows, total, current, size))
}
