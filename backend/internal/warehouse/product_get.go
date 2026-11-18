// product_get.go POST /saas-warehouse/busProduct/getById(规格 3.3 = getBusProduct+updByVisitor)。
// feign 全部换本地同库直查(实体 JSON 驼峰);微信/unipush 按本地失败形态。
// 已登记差异(徐先生决策项):
//  1. 分享红包 generateAmount 金额=区间内随机整数(Java Random 同分布,响应不含金额)。
//  2. receiveAmount 核销链不实现:本地形态=微信商家转账 HTTP 失败被自身 catch →
//     fail(603,"") 零落库(bus_red_envelope_record/batch/detail 与 clog 核销均不发生,
//     响应不受影响),故只落 generateAmount+shareRedEnvelopes 各一条 bus_red_money_clog
//     (Java 同一动作确实落 2 条,第二条 id 为 receiveAmount 的 clogId)。
//  3. 短信/unipush 不发送;getDetailImg 不在本文件(规格 3.17,/home/ubuntu/uploadfiles 缺失)。
package warehouse

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		mux.HandleFunc("POST /saas-warehouse/busProduct/getById", authGuard(s.bpGetById))
	})
}

// ── feign 实体 JSON 等价整形(chain 侧实体无 Go ent 元数据,手工补载体键/Boolean/ISO 时间) ──

// bpUserShape chain_auth_user 行 → ChainAuthUser 实体 JSON:载体键 null,
// Boolean 列 → true/false,lastLoginTime/updateTime 无 @JsonFormat → ISO T
// (expirationTime/createTime 有 @JsonFormat,驱动原生空格格式即正确;brokerStartDate LocalDate 原生)。
func bpUserShape(row map[string]any) map[string]any {
	if row == nil {
		return nil
	}
	row = entChainAuthUser.decorateRow(row)
	for _, k := range []string{"searchStr", "page", "ids", "newPassword", "code", "effSeconds", "userRoles", "roleNames", "startTime", "endTime"} {
		row[k] = nil
	}
	for _, k := range []string{"singleDatabase", "msgAlarmFlag", "isAgreeCall", "isPerfect", "isMiniProUser", "isAppUser", "isWebUser", "isPromptCare"} {
		if v, ok := row[k]; ok && v != nil {
			if b, isBool := v.(bool); isBool {
				row[k] = b // decorate 已布尔化,勿二次转换
			} else {
				row[k] = num(v) != 0
			}
		}
	}
	for _, k := range []string{"lastLoginTime", "updateTime"} {
		if s, ok := row[k].(string); ok && len(s) == 19 {
			row[k] = s[:10] + "T" + s[11:]
		}
	}
	if s, ok := row["brokerStartDate"].(string); ok && len(s) == 19 {
		row["brokerStartDate"] = s[:10] // DATE 列 parseTime 形态 → LocalDate "yyyy-MM-dd"
	}
	// expirationTimeStr 实体计算 getter;ids 请求回环初始化为 []
	if et, ok := row["expirationTime"].(string); ok && len(et) >= 10 {
		row["expirationTimeStr"] = et[:10]
	} else {
		row["expirationTimeStr"] = nil
	}
	if row["ids"] == nil {
		row["ids"] = []any{}
	}
	return row
}

// bpChainFiles chain_common_file 行 → ChainCommonFile 实体 JSON:载体键 null,
// updateTime 无 @JsonFormat → ISO T,effective/isCompressedFile Boolean。
func bpChainFiles(rows []map[string]any) []map[string]any {
	for _, r := range rows {
		if r == nil {
			continue
		}
		for _, k := range []string{"searchStr", "page", "ids", "configFilePath"} {
			r[k] = nil
		}
		if s, ok := r["updateTime"].(string); ok && len(s) == 19 {
			r["updateTime"] = s[:10] + "T" + s[11:]
		}
		for _, k := range []string{"effective", "isCompressedFile"} {
			if v, ok := r[k]; ok && v != nil {
				r[k] = num(v) != 0
			}
		}
	}
	return rows
}

// bpBrokerUser businessId/areaBuserId 经纪人行:borkerYears=Period(brokerStartDate,今天)
// 整年数+isBorker;brokerStartDate 缺失/解析失败 → false(→999)。
func bpBrokerUser(row map[string]any) (map[string]any, bool) {
	u := bpUserShape(row)
	sd := str(u, "brokerStartDate")
	if len(sd) < 10 {
		return nil, false // null.toString NPE
	}
	t, err := time.Parse("2006-01-02", sd[:10])
	if err != nil {
		return nil, false // LocalDate.parse DateTimeParseException
	}
	now := time.Now()
	years := now.Year() - t.Year() // Period.between(...).getYears()
	if now.Month() < t.Month() || (now.Month() == t.Month() && now.Day() < t.Day()) {
		years--
	}
	u["borkerYears"] = years
	u["isBorker"] = true
	return u, true
}

// ── 主流程 ──

// bpGetById = getBusProduct(Java BusProductServiceImpl 1520-1657)
func (s *Service) bpGetById(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusProduct.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	// 1. getById(id):qEntEnt2 有序整形(实体 camel 键);查询失败已写 999
	rows := s.qEntEnt2(w, r, entBusProduct, "SELECT * FROM bus_product WHERE id=?", m["id"])
	if rows == nil {
		return
	}
	if len(rows) == 0 || rows[0] == nil {
		platform.Fail(w, 602, "产品已下架")
		return
	}
	p := rows[0]
	id := str(p, "id")
	// 2. association_id null → 32 位 uuid 写回(updateById 全实体回写=值不变,仅 update_time 刷新)
	if p["associationId"] == nil {
		uid := uuid32()
		p["associationId"] = uid
		if _, err := d.Exec("UPDATE bus_product SET association_id=?, update_time=now() WHERE id=?", uid, id); err != nil {
			platform.Unknown(w)
			return
		}
	}
	// 3. status==0 且 warse_code null → createNewCode 仅内存(status 拆箱 null → 999)
	if p["status"] == nil {
		platform.Unknown(w) // 0 == getStatus() 拆箱 NPE
		return
	}
	println("DBGSTATUS raw=", fmt.Sprintf("%v", rows[0]["status"]), "shaped=", fmt.Sprintf("%v", p["status"]), "id=", fmt.Sprintf("%v", p["id"]), "warseCodeNil=", p["warseCode"] == nil)
	if num(p["status"]) == 0 && p["warseCode"] == nil {
		where := " where warse_code like '0%'" // companyId null 或 ∉{44,420} 走此分支,不炸
		switch int(num(p["companyId"])) {
		case 44:
			where = " where warse_code like 'HY%'" // parseInt 带字母必炸 → 999
		case 420:
			where = " where warse_code like 'KJ%'"
		}
		var wc sql.NullString
		if err := d.QueryRow("SELECT warse_code from bus_product" + where + " order by id desc limit 1").Scan(&wc); err != nil || !wc.Valid {
			platform.Unknown(w) // 空结果 parseInt(null) NPE
			return
		}
		n, err := strconv.Atoi(wc.String)
		if err != nil {
			platform.Unknown(w) // HY/KJ/脏数据 NumberFormatException
			return
		}
		p["warseCode"] = fmt.Sprintf("%05d", n+1)
	}
	// 4. userMap:feign /chainAuthUser/getById → 直查(实体 JSON;data null → 后续 memberType NPE 999)
	if p["createUserId"] != nil {
		u, err := queryCam(d, "SELECT * FROM chain_auth_user WHERE id=?", str(p, "createUserId"))
		if err != nil {
			platform.Unknown(w)
			return
		}
		if len(u) > 0 {
			if u[0]["nickname"] == nil || u[0]["nickname"] == "" {
				ph := str(u[0], "phone")
				suffix := ""
				if len(ph) >= 4 {
					suffix = ph[len(ph)-4:]
				}
				nn := strconvI64(time.Now().UnixNano()/1e6) + suffix
				_, _ = d.Exec("UPDATE chain_auth_user SET nickname=? WHERE id=?", nn, str(u[0], "id"))
				u[0]["nickname"] = nn
			}
			p["userMap"] = bpUserShape(u[0])
		}
	}
	// 5. etcInfo/groupInfo(selEtcById Map 保留 snake 列标签)
	if p["companyId"] != nil {
		etc, err := queryRows(d, "SELECT * from chain_basic_enterprise where id=?", str(p, "companyId"))
		if err != nil {
			platform.Unknown(w)
			return
		}
		if len(etc) > 0 {
			p["etcInfo"] = mapShape(etc[0], loadMapTypes(d, "chain_basic_enterprise"))
		} else {
			p["etcInfo"] = nil
		}
	}
	if p["groupId"] != nil {
		g, err := queryRows(d, "SELECT * from chain_basic_enterprise where id=?", str(p, "groupId"))
		if err != nil {
			platform.Unknown(w)
			return
		}
		if len(g) > 0 {
			p["groupInfo"] = mapShape(g[0], loadMapTypes(d, "chain_basic_enterprise"))
		} else {
			p["groupInfo"] = nil
		}
	} else if p["companyId"] != nil {
		gl, err := queryRows(d, "SELECT * from chain_basic_enterprise where group_id=?", str(p, "companyId"))
		if err != nil {
			platform.Unknown(w)
			return
		}
		if len(gl) > 0 {
			g, err := queryRows(d, "SELECT * from chain_basic_enterprise where id=?", str(p, "companyId"))
			if err != nil {
				platform.Unknown(w)
				return
			}
			if len(g) > 0 {
				p["groupInfo"] = mapShape(g[0], loadMapTypes(d, "chain_basic_enterprise"))
			}
		}
	}
	// 6. 三类附件(chainCommonFile findList=实体 JSON)+短视频表;顺序照 Java
	assoc := str(p, "associationId")
	rp, err := queryCam(d, "SELECT * FROM chain_common_file WHERE f_type=? AND is_compressed_file=? AND association_id=?", 2, 0, assoc)
	if err != nil {
		platform.Unknown(w)
		return
	}
	p["reportListMap"] = bpChainFiles(rp)
	im, err := queryCam(d, "SELECT * FROM chain_common_file WHERE f_type=? AND is_compressed_file=? AND association_id=?", 4, 1, assoc)
	if err != nil {
		platform.Unknown(w)
		return
	}
	p["imgListMap"] = bpChainFiles(im)
	vv, err := queryCam(d, "SELECT * FROM bus_videos WHERE pro_id=?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	p["videoList"] = entBusVideos.shape(vv)
	vm, err := queryCam(d, "SELECT * FROM chain_common_file WHERE f_type=? AND is_compressed_file=? AND association_id=?", 1, 0, assoc)
	if err != nil {
		platform.Unknown(w)
		return
	}
	p["videoListMap"] = bpChainFiles(vm)
	// 7. 同地区煤种(云仓号首段 like;mapper 无 thName 子查询,照原文)
	if p["warehouseNb"] != nil {
		first := strings.Split(str(p, "warehouseNb"), " ")[0]
		ap := s.qEntEnt2(w, r, entBusProduct, "select * from bus_product where warehouse_nb like ? and id != ? and status=1 ORDER BY create_time", first+"%", id)
		if ap == nil {
			return
		}
		p["areaProList"] = ap
	}
	// 8. 客服:configs.get(0) 无论配置是否为空都再取一次 → 空配置/查无此人 → IndexOOB 999
	var adminPhone sql.NullString
	_ = d.QueryRow("SELECT admin_phone FROM bus_config WHERE msg_type=3 ORDER BY id LIMIT 1").Scan(&adminPhone)
	if adminPhone.Valid {
		p["clientPhone"] = adminPhone.String
		p["clientName"] = "李志强" // 常量照抄
	}
	cu, err := queryCam(d, "SELECT * FROM chain_auth_user WHERE phone=? AND channel=2", adminPhone.String)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(cu) == 0 {
		platform.Unknown(w) // clientUsers.get(0) 越界(含 configs 空)
		return
	}
	cu0 := bpUserShape(cu[0])
	p["businessHead"] = ""
	if cu0["wxFaceImg"] != nil {
		p["businessHead"] = str(cu0, "wxFaceImg")
	}
	// 9. brokerList(LinkedHashMap 键序:username/phone/isBorker)
	brokerList := []any{}
	if p["coalRelName"] != nil {
		brokerList = append(brokerList, map[string]any{"username": p["coalRelName"], "phone": p["coalRelPhone"], "isBorker": false})
	}
	if p["warehouseNb"] != nil && strings.Contains(str(p, "warehouseNb"), "凯嘉") {
		brokerList = append(brokerList, map[string]any{"username": "张慧庆", "phone": "18634880999", "isBorker": false})
	}
	if p["businessId"] != nil {
		bu, err := queryCam(d, "SELECT * FROM chain_auth_user WHERE id=?", str(p, "businessId"))
		if err != nil || len(bu) == 0 {
			platform.Unknown(w) // feign data null → toString NPE
			return
		}
		u, ok := bpBrokerUser(bu[0])
		if !ok {
			platform.Unknown(w) // brokerStartDate 缺失/解析失败
			return
		}
		brokerList = append(brokerList, u)
	}
	if p["areaBuserId"] != nil {
		if p["businessId"] == nil {
			platform.Unknown(w) // businessId.intValue() NPE(先于 feign)
			return
		}
		if num(p["businessId"]) != num(p["areaBuserId"]) {
			au, err := queryCam(d, "SELECT * FROM chain_auth_user WHERE id=?", str(p, "areaBuserId"))
			if err != nil || len(au) == 0 {
				platform.Unknown(w)
				return
			}
			u, ok := bpBrokerUser(au[0])
			if !ok {
				platform.Unknown(w)
				return
			}
			brokerList = append(brokerList, u)
		}
	}
	p["brokerList"] = brokerList
	// 10. tags(selTags;无 → [])
	tags, err := queryCam(d, "SELECT t.* from bus_tag t, bus_pro_tag pt where t.id = pt.sign_id and pt.pro_id = ?", id)
	if err != nil {
		platform.Unknown(w)
		return
	}
	p["tags"] = entBusTag.shape(tags)
	// 11. 访问人数
	var vt int64
	if err := d.QueryRow("SELECT count(*) FROM bus_visitor_record WHERE pro_id=?", id).Scan(&vt); err != nil {
		platform.Unknown(w)
		return
	}
	p["vtPeopCount"] = jsonNumberInt64(vt)
	// 12. 三级分类名(分类缺失 → getById null → getProductName NPE 999)
	if p["threeProId"] == nil {
		platform.Unknown(w)
		return
	}
	var thName sql.NullString
	if err := d.QueryRow("SELECT product_name FROM bus_big_prodect WHERE id=?", str(p, "threeProId")).Scan(&thName); err != nil || !thName.Valid {
		platform.Unknown(w)
		return
	}
	p["thName"] = thName.String
	// 13. 创建者权益(getOne member_type;userMap 缺 → NPE 999;多行 TooManyResults 999)
	um, _ := p["userMap"].(map[string]any)
	if um == nil {
		platform.Unknown(w)
		return
	}
	bf, err := queryCam(d, "SELECT * FROM bus_member_benefits WHERE member_type=?", str(um, "memberType"))
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(bf) > 1 {
		platform.Unknown(w) // getOne 多行
		return
	}
	if len(bf) == 1 {
		p["busMemberBenefits"] = entBusMemberBenefits.shapeOne(bf[0])
	} else {
		p["busMemberBenefits"] = nil
	}
	// 15. updByVisitor(authGuard 包裹,jwtInfo 恒非 null)
	if !s.bpUpdByVisitor(w, r, d, m, p) {
		return
	}
	// 16. prags=[]
	p["prags"] = []any{}
	platform.Success(w, p)
}

// ── updByVisitor(Java 1700-1760) ──

// bpUpdByVisitor 收藏计数/最后点击/关注奖励/最后分享/分享红包/访客消息;false=已写响应
func (s *Service) bpUpdByVisitor(w http.ResponseWriter, r *http.Request, d *sql.DB, m, p map[string]any) bool {
	id := str(p, "id")
	uid := strconvI64(currentUser(r).UserID)
	now := nowSpace()
	// collCount(selectOne pro_id+user_id;多行 TooManyResults → 999)
	p["collCount"] = jsonNumberInt64(0)
	var cc int64
	if err := d.QueryRow("SELECT count(*) FROM bus_pro_collection WHERE pro_id=? AND user_id=?", id, uid).Scan(&cc); err != nil {
		platform.Unknown(w)
		return false
	}
	if cc > 1 {
		platform.Unknown(w)
		return false
	}
	if cc == 1 {
		p["collCount"] = jsonNumberInt64(1)
	}
	// 最后一次点击 upsert(getOne user_id 多行 → 999;构造器无 lastTypeId,插入不落该列)
	var ucCnt, ucid int64
	if err := d.QueryRow("SELECT count(*), COALESCE(MIN(id),0) FROM bus_user_click_pro WHERE user_id=?", uid).Scan(&ucCnt, &ucid); err != nil {
		platform.Unknown(w)
		return false
	}
	if ucCnt > 1 {
		platform.Unknown(w)
		return false
	}
	if ucCnt == 0 {
		if _, err := d.Exec("INSERT INTO bus_user_click_pro (user_id, last_pro_id, three_pro_id, click_time, is_last_click, create_time, update_time) VALUES (?,?,?,?,1,?,?)",
			uid, id, p["threeProId"], now, now, now); err != nil {
			platform.Unknown(w)
			return false
		}
	} else {
		if _, err := d.Exec("UPDATE bus_user_click_pro SET last_type_id=?, last_pro_id=?, click_time=?, is_last_click=1, three_pro_id=?, update_time=? WHERE id=?",
			p["threeProId"], id, now, p["threeProId"], now, strconv.FormatInt(ucid, 10)); err != nil {
			platform.Unknown(w)
			return false
		}
	}
	// isCanCareReward:当日 type6 obj4 ≥3 → false;该用户该 pro 的 type6 记录 >0 → false
	isCan := true
	var careCount int64
	if err := d.QueryRow("select count(1) from bus_pro_log where user_id=? and pro_id=? and operate_type=6 and operate_obj=4 and to_days(create_time) = to_days(now())", uid, id).Scan(&careCount); err != nil {
		platform.Unknown(w)
		return false
	}
	if careCount >= 3 {
		isCan = false
	}
	var careProCount int64
	if err := d.QueryRow("SELECT count(*) FROM bus_pro_log WHERE user_id=? AND pro_id=? AND operate_obj=4 AND operate_type=6", uid, id).Scan(&careProCount); err != nil {
		platform.Unknown(w)
		return false
	}
	if careProCount > 0 {
		isCan = false
	}
	p["isCanCareReward"] = isCan
	// lastShareTime:当日 type4/5 obj4 最新 create_time(LocalDateTime 无 @JsonFormat → ISO T)
	var ls sql.NullString
	if err := d.QueryRow("select create_time from bus_pro_log where user_id=? and pro_id=? and (operate_type=4 or operate_type=5) and operate_obj=4 and to_days(create_time) = to_days(now()) order by create_time desc limit 1", uid, id).Scan(&ls); err != nil && err != sql.ErrNoRows {
		platform.Unknown(w)
		return false
	}
	if ls.Valid && len(ls.String) == 19 {
		p["lastShareTime"] = ls.String[:10] + "T" + ls.String[11:]
	} else {
		p["lastShareTime"] = nil
	}
	// 分享红包:selTodayUser 当日新访客 × body.shareUserId 非空 × ≠访客 × 一生未提供过收益
	var todayNew int64
	if err := d.QueryRow("select count(1) from chain_auth_user where id=? and to_days(create_time) = to_days(now())", uid).Scan(&todayNew); err != nil {
		platform.Unknown(w)
		return false
	}
	if suv := m["shareUserId"]; todayNew > 0 && suv != nil && num(suv) != num(uid) {
		var lifeCnt int64
		if err := d.QueryRow("select count(*) from bus_pro_log where user_id=? and share_user_id is not null and operate_type=7 and (operate_obj=4 or operate_obj=3)", uid).Scan(&lifeCnt); err != nil {
			platform.Unknown(w)
			return false
		}
		if lifeCnt == 0 && !s.bpShareRed(w, d, uid, now) {
			return false
		}
	}
	// 访客记录:visitUser 查无 → NPE 999;createUserId≠访客 且 channel=2 → saveVisitInfo
	var vch, vphone, vusername, vnick, vware sql.NullString
	if err := d.QueryRow("SELECT channel, phone, username, nickname, ware_identity FROM chain_auth_user WHERE id=?", uid).Scan(&vch, &vphone, &vusername, &vnick, &vware); err != nil {
		platform.Unknown(w) // selectById null → getChannel NPE
		return false
	}
	if num(p["createUserId"]) != num(uid) && vch.String == "2" {
		if !s.bpSaveVisitInfo(w, d, m, p, uid, vphone.String, vusername.String, vnick.String, vware, now) {
			return false
		}
	}
	return true
}

// bpShareRed = shareRedEnvelopes 的 generateAmount 部分(已核实的落库面):
// bus_red_envelope_range getOne(receive_type=5;0 行 NPE/多行 TooManyResults/边界坏值 → 999),
// 金额=[min,max] 闭区间随机整数;同一动作 Java 落 2 条 bus_red_money_clog
// (generateAmount 内一条+shareRedEnvelopes 一条);receiveAmount 核销链=登记差异(文件头)。
func (s *Service) bpShareRed(w http.ResponseWriter, d *sql.DB, uid, now string) bool {
	rg, err := queryCam(d, "SELECT * FROM bus_red_envelope_range WHERE receive_type=5")
	if err != nil {
		platform.Unknown(w)
		return false
	}
	if len(rg) != 1 || rg[0] == nil || rg[0]["minAmount"] == nil || rg[0]["maxAmount"] == nil {
		platform.Unknown(w) // getOne 0 行 NPE / 多行异常 / min·max null intValue NPE
		return false
	}
	lo, hi := int(num(rg[0]["minAmount"])), int(num(rg[0]["maxAmount"]))
	if hi < lo {
		platform.Unknown(w) // nextInt(负数) IllegalArgumentException
		return false
	}
	amt := strconv.Itoa(lo + rand.Intn(hi-lo+1))
	for i := 0; i < 2; i++ {
		if _, err := d.Exec("INSERT INTO bus_red_money_clog (user_id, receive_type, create_amount, is_write_off, create_time, update_time) VALUES (?,?,?,0,?,?)",
			uid, 5, amt, now, now); err != nil {
			platform.Unknown(w)
			return false
		}
	}
	return true
}

// bpSaveVisitInfo = saveVisitInfo(Java 1762-1800):bus_msg(type4)+bus_user_msg+
// bus_visitor_record upsert(wx_head_img 用 nickname 的 bug 照抄)+bus_pro_log(type7,obj4)。
func (s *Service) bpSaveVisitInfo(w http.ResponseWriter, d *sql.DB, m, p map[string]any, uid, phone, username, nickname string, ware sql.NullString, now string) bool {
	id := str(p, "id")
	res, err := d.Exec("INSERT INTO bus_msg (msg_type, content, remark, jump_id, create_time, update_time) VALUES (4, ?, ?, ?, ?, ?)",
		"有1位客户浏览过您发布的煤种，快去查看吧!", id, id, now, now)
	if err != nil {
		platform.Unknown(w)
		return false
	}
	msgId, _ := res.LastInsertId()
	if _, err := d.Exec("INSERT INTO bus_user_msg (user_id, msg_id, create_time, update_time) VALUES (?,?,?,?)",
		str(p, "createUserId"), strconv.FormatInt(msgId, 10), now, now); err != nil {
		platform.Unknown(w)
		return false
	}
	var vcnt int64
	if err := d.QueryRow("SELECT count(*) FROM bus_visitor_record WHERE visitor_phone=? AND pro_id=?", phone, id).Scan(&vcnt); err != nil {
		platform.Unknown(w)
		return false
	}
	if vcnt > 1 {
		platform.Unknown(w) // getOne 多行
		return false
	}
	if vcnt == 0 {
		if _, err := d.Exec("INSERT INTO bus_visitor_record (pro_id, visitor_id, visitor_name, nickname, visitor_phone, visitor_headimg, visitor_time, visitor_resource, create_time, update_time) VALUES (?,?,?,?,?,?,?,?,?,?)",
			id, uid, username, nickname, phone, nickname, now, bpVisitorRes(m), now, now); err != nil {
			platform.Unknown(w)
			return false
		}
	} else {
		var vid int64
		if err := d.QueryRow("SELECT id FROM bus_visitor_record WHERE visitor_phone=? AND pro_id=? ORDER BY id LIMIT 1", phone, id).Scan(&vid); err != nil {
			platform.Unknown(w)
			return false
		}
		if _, err := d.Exec("UPDATE bus_visitor_record SET visitor_id=?, update_time=? WHERE id=?", uid, now, strconv.FormatInt(vid, 10)); err != nil {
			platform.Unknown(w)
			return false
		}
	}
	var shareUID, wareArg any
	if v := m["shareUserId"]; v != nil {
		shareUID = str(m, "shareUserId")
	}
	if ware.Valid {
		wareArg = ware.String
	}
	if _, err := d.Exec("INSERT INTO bus_pro_log (user_id, user_phone, pro_id, target_user_id, operate_obj, operate_type, share_user_id, ware_identity, create_time, update_time) VALUES (?,?,?,?,4,7,?,?,?,?)",
		uid, phone, id, str(p, "createUserId"), shareUID, wareArg, now, now); err != nil {
		platform.Unknown(w)
		return false
	}
	return true
}

// bpVisitorRes body.visitorResource 原样(null → NULL)
func bpVisitorRes(m map[string]any) any {
	if v, ok := m["visitorResource"]; ok && v != nil {
		return str(m, "visitorResource")
	}
	return nil
}
