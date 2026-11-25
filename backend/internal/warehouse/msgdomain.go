// msgdomain.go 批8:busMsg/busMyBill/busBusinessCard/busCustManage/busBiddingInfo/
// busCommonLog/busVisitorRecord/busVisitorLog 的 16+2 个 bespoke 端点。
// 全部行为 2026-11-17 对活体 Java 网关(localhost:8765,账号 hyqcbjhd)逐端点 curl 实证,
// 实测样本摘要写在各 handler 注释;999 形状 = {"success":false,"code":999,...} HTTP 500。
package warehouse

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		reg := func(path string, h http.HandlerFunc) {
			mux.HandleFunc("POST /saas-warehouse"+path, authGuard(h))
		}
		reg("/busMsg/selPriceMsg", s.msgSelPriceMsg)
		reg("/busMsg/selWdCount", s.msgSelWdCount)
		reg("/busMsg/selMyMsg", s.msgSelMyMsg)
		reg("/busMsg/selNoReadMsg", s.msgSelNoReadMsg)
		reg("/busMsg/updBthById", s.msgUpdBthById)
		reg("/busMsg/oncClickRead", s.msgOncClickRead)
		reg("/busMyBill/getPayInfo", s.billGetPayInfo)
		reg("/busBusinessCard/searchDefault", s.cardSearchDefault)
		reg("/busBusinessCard/searchHasDefault", s.cardSearchHasDefault)
		reg("/busBusinessCard/listBsCard", s.cardListBsCard)
		reg("/busBusinessCard/searchByUser", s.cardSearchByUser)
		reg("/busCustManage/noticeCust", s.custNoticeCust)
		reg("/busBiddingInfo/getBiddingInfo", s.bidGetBiddingInfo)
		reg("/busCommonLog/saveShare", s.clSaveShare)
		reg("/busCommonLog/listByPage", s.clListByPage)
		reg("/busVisitorRecord/selVisited", s.vrSelVisited)
		reg("/busVisitorRecord/listByPage", s.vrListByPage)
		reg("/busVisitorLog/list", s.vlAllLog)
	})
}

// bodyMapStrict @RequestBody Map 版:空体按 {};非对象 JSON(如 [])Jackson 绑定失败 → 999
func bodyMapStrict(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, 32<<20))
	if err != nil {
		platform.Unknown(w)
		return nil, false
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		return map[string]any{}, true
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil || m == nil {
		platform.Unknown(w)
		return nil, false
	}
	return m, true
}

// entFieldSet 实体 @ExcelColumn 字段集合(StringHelper.getSql(Map,cls) 的 checkFields 只认
// 带 @ExcelColumn 注解的字段,searchStr/page 等载体键不生成条件——8765 实测 searchStr/nickname 体均 200)
func entFieldSet(e *ent) map[string]bool {
	out := make(map[string]bool, len(e.q))
	for _, f := range e.q {
		out[f.Name] = true
	}
	return out
}

// ── busMsg ──

// msgSelPriceMsg = selPriceMsg:getSql(searchSql) + 固定 join,当前用户维度的价格变动消息。
// Java 实测(8765):{} → 200 data:[](hyqcbjhd 无消息行);
// SQL:select * from bus_msg m, bus_user_msg um ${searchSql} and m.id=um.msg_id and um.user_id=#{maps.userId}
func (s *Service) msgSelPriceMsg(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusMsg.validateBody(m) {
		platform.Unknown(w)
		return
	}
	sqlStr, args, ok := entBusMsg.getSql(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	q := "select * from bus_msg m, bus_user_msg um " + sqlStr + " and m.id=um.msg_id and um.user_id=?"
	rows := s.qEntEnt2(w, r, entBusMsg, q, append(args, strconvI64(currentUser(r).UserID))...)
	if rows == nil {
		return
	}
	platform.Success(w, rows)
}

// msgSelWdCount = selWdCount:未读消息分组计数(controller 无 @RequestBody,体内容不参与绑定)。
// Java 实测:{} → 200 data={"coalAddCount":0,"sumCount":0,"checkCount":0,"priceMsgCount":0,"custCount":0}
// (hyqcbjhd 实测全 0);sumCount = priceMsgCount + checkCount(类型 1+3)动态累加
func (s *Service) msgSelWdCount(w http.ResponseWriter, r *http.Request) {
	rows, err := queryRows(s.pick(r),
		"select msg_type msgType, count(1) msgCount from bus_msg m, bus_user_msg um where um.read_status=0 and m.id=um.msg_id and um.user_id=? and m.is_effect=1 GROUP BY m.msg_type",
		strconvI64(currentUser(r).UserID))
	if err != nil {
		platform.Unknown(w)
		return
	}
	zero := jsonNumberInt64(0)
	res := map[string]any{"priceMsgCount": zero, "coalAddCount": zero, "checkCount": zero, "custCount": zero}
	var sum int64
	for _, row := range rows {
		n := jsonNumberInt64(int64(num(row["msgCount"])))
		switch int64(num(row["msgType"])) {
		case 1:
			res["priceMsgCount"] = n
			sum += int64(num(row["msgCount"]))
		case 2:
			res["coalAddCount"] = n
		case 3:
			res["checkCount"] = n
			sum += int64(num(row["msgCount"]))
		case 4:
			res["custCount"] = n
		}
	}
	res["sumCount"] = jsonNumberInt64(sum)
	platform.Success(w, res)
}

// msgSelMyMsgPage selMyMsg/selNoReadMsg 共用:read_status 字面条件(Java " and read_status="+Integer)
// + 分页(mapper 带 Page 参数 → MP 拦截器分页,返回类型 List → data 恒为当前页数组;无 page → 全量)
func (s *Service) msgSelMyMsgPage(w http.ResponseWriter, r *http.Request, e *ent, q string, baseArgs []any, m map[string]any) {
	_, hasPage := m["page"].(map[string]any)
	if !hasPage {
		rows := s.qEntEnt2(w, r, e, q, baseArgs...)
		if rows == nil {
			return
		}
		platform.Success(w, rows)
		return
	}
	current, size := pageOf(m)
	countQ := "SELECT COUNT(*) FROM (" + stripOrderBy(q) + ") cnt"
	var total int64
	if err := s.pick(r).QueryRow(countQ, baseArgs...).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	off := (current - 1) * size
	rows := s.qEntEnt2(w, r, e, q+" LIMIT ? OFFSET ?", append(append([]any{}, baseArgs...), size, off)...)
	if rows == nil {
		return
	}
	platform.Success(w, rows)
}

func msgReadStatusCond(m map[string]any) (string, bool) {
	if v, has := m["readStatus"]; has && v != nil {
		bv, ok := bindVal("i", v)
		if !ok || bv == nil {
			return "", false
		}
		return " and read_status=" + bv.(string), true
	}
	return "", true
}

// msgSelMyMsg = selMyMsg:我的消息(join bus_user_msg 限定当前用户 + is_effect=1)。
// Java 实测:{} → 200 data:[](hyqcbjhd 名下无消息行);实体键含 userMsgId/readStatus
func (s *Service) msgSelMyMsg(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusMsg.validateBody(m) {
		platform.Unknown(w)
		return
	}
	sqlStr, args, ok := entBusMsg.getSql(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	cond, ok2 := msgReadStatusCond(m)
	if !ok2 {
		platform.Unknown(w)
		return
	}
	q := "select DISTINCT m.*, um.id as userMsgId, um.read_status as readStatus from bus_msg m, bus_user_msg um " +
		sqlStr + cond + " and m.id=um.msg_id and um.user_id=? and m.is_effect=1 order by m.create_time desc"
	s.msgSelMyMsgPage(w, r, entBusMsg, q, append(args, strconvI64(currentUser(r).UserID)), m)
}

// selNoReadMsg matches the running Java warehouse JAR: it does not update
// read_status. Marking messages read belongs to the explicit write routes.
func (s *Service) msgSelNoReadMsg(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusMsg.validateBody(m) {
		platform.Unknown(w)
		return
	}

	sqlStr, args, ok := entBusMsg.getSql(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	cond, ok2 := msgReadStatusCond(m)
	if !ok2 {
		platform.Unknown(w)
		return
	}
	q := "select DISTINCT m.*, um.id as userMsgId, um.read_status as readStatus from bus_msg m, bus_user_msg um " +
		sqlStr + cond + " and m.id=um.msg_id and m.is_effect=1 order by m.create_time desc"
	s.msgSelMyMsgPage(w, r, entBusMsg, q, args, m)
}

// msgUpdBthById = updBthById:@RequestBody List<BusMsg> 逐行 updateBatchById(非空子集+fill)。
// Java 实测:零字节体 → 999(Jackson 缺体),{} → 999(对象绑 List 失败),[] → 200 data:"批量修改成功"
func (s *Service) msgUpdBthById(w http.ResponseWriter, r *http.Request) {
	var items []any
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		platform.Unknown(w)
		return
	}
	for _, it := range items {
		em, _ := it.(map[string]any)
		if em == nil || !entBusMsg.validateBody(em) {
			platform.Unknown(w)
			return
		}
		vals := pick(em, entBusMsg.dbcols)
		delete(vals, "id")
		if _, err := platform.UpdateByID(s.pick(r), entBusMsg.table, em["id"], vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, "批量修改成功")
}

// msgOncClickRead = oncClickRead:按消息类型一键已读(${msgType} 字面注入,缺 → MyBatis 异常)。
// Java 实测:{} → 999;{"msgType":99999999} → 200 data:null(update 0 行后 success())
func (s *Service) msgOncClickRead(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusMsg.validateBody(m) {
		platform.Unknown(w)
		return
	}
	mt, ok := bindVal("i", m["msgType"])
	if !ok || mt == nil {
		platform.Unknown(w) // ${msgType} null → 异常(实测 {} 999)
		return
	}
	if _, err := s.pick(r).Exec("update bus_user_msg um, bus_msg m set read_status=1 where um.msg_id=m.id and msg_type=" + mt.(string)); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// ── busMyBill ──

// billGetPayInfo follows the Java transaction, fee/coupon calculation and bill
// creation; the local unavailable merchant boundary rolls the transaction back.
func (s *Service) billGetPayInfo(w http.ResponseWriter, r *http.Request) {
	s.openMembershipPay(w, r)
}

// ── busBusinessCard ──

// cardUidOf searchDefault/searchHasDefault 共用:body userId 非 null 且非 0 → 覆盖,否则当前用户
func cardUidOf(m map[string]any, r *http.Request) (string, bool) {
	if v, has := m["userId"]; has && v != nil {
		bv, ok := bindVal("i", v)
		if !ok {
			return "", false
		}
		if bv != nil && num(bv.(string)) != 0 {
			return bv.(string), true
		}
	}
	return strconvI64(currentUser(r).UserID), true
}

// msgSearchDefaultQ searchDefault mapper SQL 原文:${userId} 字面
func msgSearchDefaultQ(uid string) string {
	return "select * from bus_business_card where user_id=" + uid + " and is_default = 1"
}

// cardSearchDefault = searchDefault:默认名片列表。
// Java 实测:{} → 200 data:[](hyqcbjhd 无 is_default=1 名片)
func (s *Service) cardSearchDefault(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusBusinessCard.validateBody(m) {
		platform.Unknown(w)
		return
	}
	uid, ok := cardUidOf(m, r)
	if !ok {
		platform.Unknown(w)
		return
	}
	rows := s.qEntEnt2(w, r, entBusBusinessCard, msgSearchDefaultQ(uid))
	if rows == nil {
		return
	}
	platform.Success(w, rows)
}

// cardSearchHasDefault = searchHasDefault:body userId null → NPE 999(实测 {});
// 有 userId → searchDefault 计数,0 → data:"false",否则 data:"true"
func (s *Service) cardSearchHasDefault(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusBusinessCard.validateBody(m) {
		platform.Unknown(w)
		return
	}
	v, has := m["userId"]
	if !has || v == nil {
		platform.Unknown(w) // getUserId().toString() NPE(实测 {} 999)
		return
	}
	bv, ok := bindVal("i", v)
	if !ok || bv == nil {
		platform.Unknown(w)
		return
	}
	rows, err := queryRows(s.pick(r), "select id from bus_business_card where user_id="+bv.(string)+" and is_default = 1")
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(rows) == 0 {
		platform.Success(w, "false")
		return
	}
	platform.Success(w, "true")
}

// cardListBsCard = listBsCard:我的名片夹模糊查询(bus_contacts × bus_business_card)。
// searchStr null 先置 ""(getMaps 包裹 → "%%");searchSql + " and c.user_id=<jwt uid>"。
// Java 实测:{} → 200 data:[](hyqcbjhd 无联系人)
func (s *Service) cardListBsCard(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusBusinessCard.validateBody(m) {
		platform.Unknown(w)
		return
	}
	sqlStr, args, ok := entBusBusinessCard.getSql(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	sqlStr += " and c.user_id=" + strconvI64(currentUser(r).UserID)
	q := "SELECT bc.*, c.create_time as contactsTime, c.id as cid from bus_contacts c, bus_business_card bc " +
		sqlStr + " and c.card_id=bc.id and CONCAT(IFNULL(bc.name, ''), IFNULL(bc.duty, ''), IFNULL(bc.company, ''), IFNULL(bc.phone, ''), IFNULL(bc.region, ''), IFNULL(bc.address, '')) like ? order by create_time desc"
	search := "%" + str(m, "searchStr") + "%" // setSearchStr(null→"") 后 getMaps 包裹
	rows := s.qEntEnt2(w, r, entBusBusinessCard, q, append(args, search)...)
	if rows == nil {
		return
	}
	platform.Success(w, rows)
}

// cardSearchByUser = searchByUser:cardInfos(该用户名片)+userInfo(chain_auth_user 首行)+
// etcInfo(其企业)。Java 实测:{} → 200 data:{cardInfos:[], userInfo:<chain_auth_user id=2 行>,
// etcInfo:<chain_basic_enterprise id=1 全键>}——feign findList({id:null}) 实测返回全表首行(最小 id),
// users.get(0) 空表会 IndexOutOfBounds → 999。
// 差异登记:saas-common getById 关联填充的 pubFlows/perFlows(签署流程,Go 库无该表)恒 null,
// saleEtcs/etcs/urcs 实测(企业1)为空数组,Go 恒 []。
func (s *Service) cardSearchByUser(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusBusinessCard.validateBody(m) {
		platform.Unknown(w)
		return
	}
	var uid any // mapper searchByUser(#{userId}):缺参 → user_id=NULL 空行(实测 {} cardInfos:[])
	if v, has := m["userId"]; has && v != nil {
		bv, ok := bindVal("i", v)
		if !ok {
			platform.Unknown(w)
			return
		}
		uid = bv
	}
	cards := s.qEntEnt2(w, r, entBusBusinessCard, "SELECT * from bus_business_card where user_id=?", uid)
	if cards == nil {
		return
	}
	urows, ucols, err := s.queryEnt(r, "SELECT * FROM chain_auth_user ORDER BY id LIMIT 1")
	if err != nil {
		platform.Unknown(w)
		return
	}
	if len(urows) == 0 {
		platform.Unknown(w) // users.get(0) IndexOutOfBounds
		return
	}
	userInfo := entChainAuthUserLive.shapeRows(urows, ucols)[0]
	userInfo["password"] = nil // 任意登录用户可调,不回传口令散列
	var etcInfo any
	if companyId := userInfo["companyId"]; companyId != nil {
		erow, err := platform.GetByID(s.pick(r), "chain_basic_enterprise", companyId)
		if err != nil {
			platform.Unknown(w)
			return
		}
		if erow != nil {
			ecols := loadMapTypes(s.pick(r), "chain_basic_enterprise")
			etcs := entChainBasicEnt.shapeRows([]map[string]any{erow}, colNames(ecols))
			if len(etcs) > 0 {
				etcInfo = etcs[0]
				if etcInfo != nil {
					em := etcInfo.(map[string]any)
					// 签署流程/关联清单(对照 ChainBasicEnterpriseServiceImpl.getEtc 381-409)
					pubRows, pubCols, err := s.queryEnt(r, "SELECT * FROM chain_eqb_config WHERE flow_type=1 ORDER BY id")
					if err != nil {
						platform.Unknown(w)
						return
					}
					em["pubFlows"] = entChainEqbConfig.shapeRows(pubRows, pubCols)
					perRows, perCols, err := s.queryEnt(r, "SELECT * FROM chain_eqb_config WHERE flow_type=2 ORDER BY id")
					if err != nil {
						platform.Unknown(w)
						return
					}
					em["perFlows"] = entChainEqbConfig.shapeRows(perRows, perCols)
					em["saleEtcs"] = []any{} // chain_tlzf_sale_etc buy_company_id(企业1 实测 0 行)
					em["etcs"] = []any{}     // chain_basic_enterprise group_id(实测 0 行)
					em["urcs"] = []any{}     // chain_auth_user company_id+channel=2(实测 0 行)
				}
			}
		}
	}
	platform.Success(w, map[string]any{"cardInfos": cards, "userInfo": userInfo, "etcInfo": etcInfo})
}

// entChainEqbConfig chain_eqb_config(对照 saas-common esign/ChainEqbConfig 模型:
// 字段全小写尾缀,DB 列混大小写,filepassword/autoarchive/...;键名=Jackson 原样)
var entChainEqbConfig = &ent{
	table: "chain_eqb_config",
	keys: []string{"id", "planName", "encryption", "filepassword", "autoarchive", "autoinitiate", "contractremind",
		"contractvalidity", "batchdropseal", "billisolationcode", "countdown", "mobileshieldway", "noticedeveloperurl",
		"noticetype", "orgavailableauthtypes", "personauthadvancedenabled", "personavailableauthtypes",
		"redirectdelaytime", "redirecturl", "signplatform", "willtypes", "remark", "signvalidity", "platformsign",
		"actorindentitytype", "assignedposbean", "autoexecute", "pospage", "posx", "posy", "sealtype", "fontsize",
		"format", "qpospage", "qposx", "qposy", "signdatebeantype", "signtype", "width", "signorder",
		"thirdorderno", "flowType", "forcereadtime", "createTime", "updateTime", "remark1", "remark2", "remark3", "remark4", "remark5", "searchStr", "page", "ids"},
	dbcols: []string{"id", "plan_name", "encryption", "filePassword", "autoArchive", "autoInitiate", "contractRemind",
		"contractvalidity", "batchDropSeal", "billIsolationCode", "countdown", "mobileShieldWay", "noticeDeveloperUrl",
		"noticeType", "orgAvailableAuthTypes", "personAuthAdvancedEnabled", "personAvailableAuthTypes",
		"redirectDelayTime", "redirectUrl", "signPlatform", "willTypes", "remark", "signValidity", "platformSign",
		"actorIndentityType", "assignedPosbean", "autoExecute", "posPage", "posX", "posY", "sealType", "fontSize",
		"format", "qposPage", "qposX", "qposY", "signDateBeanType", "signType", "width", "signOrder",
		"thirdOrderNo", "flow_type", "forceReadTime", "create_time", "update_time"},
	types: []jf{{"id", "i"}, {"planName", "s"}, {"encryption", "i"}, {"filepassword", "s"}, {"autoarchive", "b"},
		{"autoinitiate", "b"}, {"contractremind", "i"}, {"contractvalidity", "i"}, {"batchdropseal", "b"},
		{"billisolationcode", "s"}, {"countdown", "i"}, {"mobileshieldway", "i"}, {"noticedeveloperurl", "s"},
		{"noticetype", "s"}, {"orgavailableauthtypes", "s"}, {"personauthadvancedenabled", "s"},
		{"personavailableauthtypes", "s"}, {"redirectdelaytime", "i"}, {"redirecturl", "s"}, {"signplatform", "s"},
		{"willtypes", "s"}, {"remark", "s"}, {"signvalidity", "i"}, {"platformsign", "b"},
		{"actorindentitytype", "i"}, {"assignedposbean", "b"}, {"autoexecute", "b"}, {"pospage", "s"},
		{"posx", "s"}, {"posy", "s"}, {"sealtype", "s"}, {"fontsize", "i"}, {"format", "s"}, {"qpospage", "s"},
		{"qposx", "s"}, {"qposy", "s"}, {"signdatebeantype", "i"}, {"signtype", "i"}, {"width", "i"},
		{"signorder", "i"}, {"thirdorderno", "s"}, {"flowType", "i"}, {"forcereadtime", "i"},
		{"createTime", "t"}, {"updateTime", "t"}, {"remark1", "s"}, {"remark2", "s"}, {"remark3", "s"}, {"remark4", "s"}, {"remark5", "s"}, {"searchStr", "s"}, {"page", "page"}, {"ids", "list"}},
}

// colNames 列类型清单转列名
func colNames(cols []colType) []string {
	out := make([]string, 0, len(cols))
	for _, c := range cols {
		out = append(out, c.name)
	}
	return out
}

// ── busCustManage ──

// custNoticeCust = noticeCust:一键通知我的客户。Java 顺序:getById(proId)→null NPE 999 →
// noticeTime+1 updateById → phones null NPE 999 → 真实短信(SendMessageUtil)。
// Java 实测:{} → 999(proId null → getById(null) null → NPE,无写库)。
// phones 就绪后通过原腾讯云短信协议发送。
func (s *Service) custNoticeCust(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusCustManage.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	bp, err := platform.GetByID(d, "bus_product", m["proId"])
	if err != nil || bp == nil {
		platform.Unknown(w) // getById(null/无效) → bp null → NPE(实测 {} 999)
		return
	}
	if bp["notice_time"] == nil {
		platform.Unknown(w)
		return
	} // Java Integer unboxing precedes any write.
	if _, err := d.Exec("UPDATE bus_product SET notice_time=?, update_time=? WHERE id=?",
		jsonNumberInt64(int64(num(bp["notice_time"]))+1), nowSpace(), strOf(bp["id"])); err != nil {
		platform.Unknown(w)
		return
	}
	ph, has := m["phones"]
	if !has || ph == nil {
		platform.Unknown(w) // for (String phone : getPhones()) NPE(update 已执行,Java 同序)
		return
	}
	values, ok := ph.([]any)
	if !ok {
		platform.Unknown(w) // phones 绑 List<String> 失败
		return
	}
	phones := make([]string, 0, len(values))
	for _, value := range values {
		phones = append(phones, javaStr(value))
	}
	price := ""
	if num(bp["ex_factory_price"]) != 0 {
		price = javaStr(bp["ex_factory_price"]) + "元（出厂价）"
	} else if num(bp["plate_price"]) != 0 {
		price = javaStr(bp["plate_price"]) + "元（车板价）"
	}
	platform.SendLegacySMSMany(r.Context(), phones, []string{javaStr(bp["enterprise_name"]) + " " + javaStr(bp["warehouse_nb"]) + " " + price}, "2115763", "")
	platform.Success(w, nil)
}

// ── busBiddingInfo ──

// bidGetBiddingInfo performs the Java vendor fetch, strict timestamp selection,
// and transactional saveBatch. BIDDING_SOURCE_URL allows a local fixture server.
func (s *Service) bidGetBiddingInfo(w http.ResponseWriter, r *http.Request) {
	s.syncBiddingInfo(w, r)
}

// ── busCommonLog ──

// clSaveShare = saveShare:每次调用固定写一条操作日志再 success()(data null)。
// Java 实测:{} → 200 data:null(save(new BusCommonLog(uid, username, "", 0, "小程序分享\n", 4)))
func (s *Service) clSaveShare(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	vals := map[string]any{
		"user_id":        jsonNumberInt64(u.UserID),
		"user_name":      u.Username,
		"user_phone":     "",
		"pro_id":         jsonNumberInt64(0),
		"operation":      "小程序分享\n",
		"operation_type": jsonNumberInt64(4),
	}
	if _, err := platform.Insert(s.pick(r), "bus_common_log", vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// entListByPage 自定义 count 的 listByPage(带 page → MP JsqlParser 优化 count:去 order by、
// select 列表换 COUNT(*),不包派生表——JOIN/* 派生表在 MySQL 报 Duplicate column,
// Go 通用模板 `SELECT COUNT(*) FROM (sql) cnt` 对这两个 mapper 恒 999 的根因,8799 实测+库内复现);
// 无 page → mapper 传 null Page → "can't found IPage for args" 999(实测 {} 双双 999)
func (s *Service) entListByPage(w http.ResponseWriter, r *http.Request, e *ent, selectSQL, countSQL string) {
	m := bodyMap(r)
	if _, hasPage := m["page"].(map[string]any); !hasPage {
		platform.Unknown(w) // null Page → MP Assert 异常(实测 {} 999)
		return
	}
	if !e.validateBody(m) {
		platform.Unknown(w)
		return
	}
	sqlStr, args, ok := e.getSql(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	var search any
	if v, has := m["searchStr"]; has && v != nil {
		ss, ok2 := bindVal("s", v)
		if !ok2 {
			platform.Unknown(w)
			return
		}
		search = "%" + ss.(string) + "%"
	}
	q := strings.Replace(selectSQL, "${searchSql}", sqlStr, 1)
	q = strings.Replace(q, "#{maps.searchStr}", "?", 1)
	countQ := strings.Replace(countSQL, "${searchSql}", sqlStr, 1)
	countQ = strings.Replace(countQ, "#{maps.searchStr}", "?", 1)
	baseArgs := append(append([]any{}, args...), search)
	current, size := int64(1), int64(10)
	if pageObj, hasPage := m["page"].(map[string]any); hasPage {
		if v, ok := pageObj["current"]; ok && v != nil {
			current = int64(num(v))
		}
		if v, ok := pageObj["size"]; ok && v != nil {
			size = int64(num(v))
		}
	}
	var total int64
	if err := s.pick(r).QueryRow(countQ, baseArgs...).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	off := (current - 1) * size
	rows := s.qEntEnt2(w, r, e, q+" LIMIT ? OFFSET ?", append(baseArgs, size, off)...)
	if rows == nil {
		return
	}
	platform.Success(w, ipageOf(rows, total, current, size))
}

// clListByPage = busCommonLog/listByPage:
// Java 实测:{} → 999;{"page":{}} → 200 IPage(records:[],total:0,pages:0)
func (s *Service) clListByPage(w http.ResponseWriter, r *http.Request) {
	s.entListByPage(w, r, entBusCommonLog,
		"select l.*, u.phone from bus_common_log l, chain_auth_user u ${searchSql} and l.user_id=u.id and CONCAT(u.phone) like #{maps.searchStr} order by l.create_time desc",
		"SELECT COUNT(*) from bus_common_log l, chain_auth_user u ${searchSql} and l.user_id=u.id and CONCAT(u.phone) like #{maps.searchStr}")
}

// vrListByPage = busVisitorRecord/listByPage:同上(visitor_headimg 子查询 + bus_visitor_record.*
// 派生表列重名,通用 count 包装恒 999 的另一根因)
func (s *Service) vrListByPage(w http.ResponseWriter, r *http.Request) {
	s.entListByPage(w, r, entBusVisitorRecord,
		"select bus_visitor_record.*, (select wx_face_img from chain_auth_user where id=visitor_id) visitor_headimg from bus_visitor_record ${searchSql} and CONCAT(id) like #{maps.searchStr} order by create_time desc",
		"SELECT COUNT(*) from bus_visitor_record ${searchSql} and CONCAT(id) like #{maps.searchStr}")
}

// ── busVisitorRecord ──

// isoTime 原 "yyyy-MM-dd HH:mm:ss" → LocalDateTime.toString 的 ISO T(MyBatis Map 行经 Jackson)
func isoTime(v any) any {
	if s, ok := v.(string); ok && len(s) == 19 {
		return s[:10] + "T" + s[11:]
	}
	return v
}

// vrSelVisited = selVisited:产品维度访客记录(mapper select 只用 ${maps.proId},searchSql 是死参数;
// ${maps.proId} 缺键 → MyBatis 异常)。Java 实测:{} → 999;{"proId":1} → 200 data:[];
// {"proId":1,"page":{}} → 200 IPage(records:[],total:0,pages:0)
func (s *Service) vrSelVisited(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusVisitorRecord.validateBody(m) {
		platform.Unknown(w)
		return
	}
	pv, ok := bindVal("i", m["proId"])
	if !ok || pv == nil {
		platform.Unknown(w) // ${maps.proId} 缺键(实测 {} 999)
		return
	}
	q := "SELECT (SELECT wx_face_img from chain_auth_user where phone=visitor_phone and channel=2) wxFaceImg,visitor_phone visitorPhone, update_time visitedTime from bus_visitor_record where pro_id=" + pv.(string)
	_, hasPage := m["page"].(map[string]any)
	if !hasPage {
		rows, err := queryRows(s.pick(r), q)
		if err != nil {
			platform.Unknown(w)
			return
		}
		for _, row := range rows {
			row["visitedTime"] = isoTime(row["visitedTime"])
		}
		platform.Success(w, rows)
		return
	}
	current, size := pageOf(m)
	countQ := "SELECT COUNT(*) FROM (" + q + ") cnt"
	var total int64
	if err := s.pick(r).QueryRow(countQ).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	off := (current - 1) * size
	rows, err := queryRows(s.pick(r), q+" LIMIT ? OFFSET ?", size, off)
	if err != nil {
		platform.Unknown(w)
		return
	}
	for _, row := range rows {
		row["visitedTime"] = isoTime(row["visitedTime"])
	}
	platform.Success(w, platform.IPage(rows, total, current, size))
}

// ── busVisitorLog ──

// vlAllLog = busVisitorLog/list(controller /list 是 bespoke allLog,非模板 allEq——已从模板 skip):
// @RequestBody Map;getSql(Map,cls) 对 body 中命中实体字段的键生成 ` and snake like #{maps.key}`
// (参数化原始值);target_user_id 键(非实体字段)字面拼接;mapper ${maps.searchStr} 字面注入,
// 缺键 → MyBatis ${} 空串 → like '%%'(任务给定语义)。返回 List<BusVisitorLog>(ifCare 子查询计数)。
// Java 实测:{} → 200 data:[{id:158,...,"ifCare":1,...}](create_time desc,createTime 空格格式)
func (s *Service) vlAllLog(w http.ResponseWriter, r *http.Request) {
	m, ok := bodyMapStrict(w, r)
	if !ok {
		return
	}
	sqlStr := "where 1=1"
	var args []any
	// Generic StringHelper.getSql(Map) contributes no predicates in Java.
	if v, has := m["target_user_id"]; has && v != nil {
		// Java 是 maps.get → 字面拼接;这里按整数绑定(非整数照旧 999),防 SQL 注入
		tid, okT := bindVal("i", v)
		if !okT || tid == nil {
			platform.Unknown(w)
			return
		}
		sqlStr += " and target_user_id=?"
		args = append(args, tid)
	}
	inner := "%%" // 缺键时 Java 得空串 → like '%%'
	if v, has := m["searchStr"]; has && v != nil {
		inner = "%" + strOf(v) + "%"
	}
	q := "select *, (select count(1) from bus_user_care where care_user_id=bus_visitor_log.target_user_id) ifCare from bus_visitor_log " +
		sqlStr + " and CONCAT(id) like ? order by create_time desc"
	args = append(args, inner)
	rows := s.qEntEnt2(w, r, entBusVisitorLog, q, args...)
	if rows == nil {
		return
	}
	platform.Success(w, rows)
}

// ── 跨服务模型实体(searchByUser 响应整形用;字段序=Java 模型声明序=Jackson 键序)──

// entChainAuthUserLive = meta_gen entChainAuthUser(模型 68 键)+ 部署版 saas-auth 实测多出的
// expirationTimeStr 键(8765 实测 userInfo 共 69 键;值 null)。浅拷贝追加,不改 meta_gen 已有实体
var entChainAuthUserLive = func() *ent {
	e := *entChainAuthUser
	e.keys = append(append([]string{}, entChainAuthUser.keys...), "expirationTimeStr")
	e.types = append(append([]jf{}, entChainAuthUser.types...), jf{"expirationTimeStr", "s"})
	return &e
}()

// entChainBasicEnt saas-common ChainBasicEnterprise(saas-supply-chain/saas-common/business/entity,
// 113 字段,8765 实测 etcInfo 全键与类型逐键核对)
var entChainBasicEnt = &ent{
	table: "chain_basic_enterprise",
	keys:  []string{"id", "enterpriseCode", "adminName", "enterpriseName", "simpleName", "account", "creditCode", "address", "legalPerson", "controllerPerson", "contactsPerson", "enterpriseType", "provinceCode", "cityCode", "areaCode", "countyCode", "type", "regTime", "rateNumber", "shTime", "phone", "legIdentity", "regMoney", "regAdress", "businesser", "countyName", "boss", "statusMark", "designPro", "relOutput", "proName", "cityName", "areaName", "countryName", "isNetwork", "rotationFlag", "pubFlowId", "perFlowId", "eqbOrgId", "merchantNumber", "merchantAccount", "signFreeCount", "signSurplusCount", "signChargeCount", "serviceExpirationTime", "advertisingVideoId", "vrLink", "websiteImgId", "websiteCedImgId", "isChainEtc", "isBusinessEtc", "isSignEtc", "isMemberEtc", "wareIdentity", "logo", "imageUrl", "recCap", "econKind", "updatedDate", "scope", "isAuthentication", "groupId", "tlSignNum", "tlAcctSubNo", "tlMchtNo", "tlOpenTime", "busLicenseId", "accOpeningPermitId", "legalFtCardId", "legalBkCardId", "quaCertScore", "basInfoScore", "coalVerifyScore", "updShareScore", "reportScore", "qualifiedScore", "updPriceScore", "ycoalSilo", "jcoalSilo", "ycoalSiloCapacity", "jcoalSiloCapacity", "testEquipmentType", "isDisplayed", "createTime", "updateTime", "remark1", "remark2", "remark3", "remark4", "remark5", "searchStr", "page", "ids", "userId", "filelist", "imagelist", "password", "yyzz", "lbts", "pubFlows", "perFlows", "isAutoAdd", "startTime", "endTime", "isTlEtc", "tlzfConfig", "saleEtcs", "etcs", "isOtherEtc", "sumScore", "urcs", "selEtcType", "proCount"},
	types: []jf{{"id", "i"}, {"enterpriseCode", "s"}, {"adminName", "s"}, {"enterpriseName", "s"}, {"simpleName", "s"}, {"account", "s"}, {"creditCode", "s"}, {"address", "s"}, {"legalPerson", "s"}, {"controllerPerson", "s"}, {"contactsPerson", "s"}, {"enterpriseType", "s"}, {"provinceCode", "s"}, {"cityCode", "s"}, {"areaCode", "s"}, {"countyCode", "s"}, {"type", "s"}, {"regTime", "D"}, {"rateNumber", "s"}, {"shTime", "s"}, {"phone", "s"}, {"legIdentity", "s"}, {"regMoney", "s"}, {"regAdress", "s"}, {"businesser", "s"}, {"countyName", "s"}, {"boss", "s"}, {"statusMark", "i"}, {"designPro", "s"}, {"relOutput", "s"}, {"proName", "s"}, {"cityName", "s"}, {"areaName", "s"}, {"countryName", "s"}, {"isNetwork", "b"}, {"rotationFlag", "b"}, {"pubFlowId", "i"}, {"perFlowId", "i"}, {"eqbOrgId", "s"}, {"merchantNumber", "s"}, {"merchantAccount", "s"}, {"signFreeCount", "i"}, {"signSurplusCount", "i"}, {"signChargeCount", "i"}, {"serviceExpirationTime", "D"}, {"advertisingVideoId", "i"}, {"vrLink", "s"}, {"websiteImgId", "i"}, {"websiteCedImgId", "i"}, {"isChainEtc", "b"}, {"isBusinessEtc", "b"}, {"isSignEtc", "b"}, {"isMemberEtc", "b"}, {"wareIdentity", "i"}, {"logo", "i"}, {"imageUrl", "s"}, {"recCap", "s"}, {"econKind", "s"}, {"updatedDate", "s"}, {"scope", "s"}, {"isAuthentication", "b"}, {"groupId", "i"}, {"tlSignNum", "s"}, {"tlAcctSubNo", "s"}, {"tlMchtNo", "s"}, {"tlOpenTime", "T"}, {"busLicenseId", "s"}, {"accOpeningPermitId", "s"}, {"legalFtCardId", "s"}, {"legalBkCardId", "s"}, {"quaCertScore", "f"}, {"basInfoScore", "f"}, {"coalVerifyScore", "f"}, {"updShareScore", "f"}, {"reportScore", "f"}, {"qualifiedScore", "f"}, {"updPriceScore", "f"}, {"ycoalSilo", "i"}, {"jcoalSilo", "i"}, {"ycoalSiloCapacity", "d"}, {"jcoalSiloCapacity", "d"}, {"testEquipmentType", "i"}, {"isDisplayed", "b"}, {"createTime", "t"}, {"updateTime", "T"}, {"remark1", "s"}, {"remark2", "s"}, {"remark3", "s"}, {"remark4", "s"}, {"remark5", "s"}, {"searchStr", "s"}, {"page", "page"}, {"ids", "list"}, {"userId", "i"}, {"filelist", "list"}, {"imagelist", "list"}, {"password", "s"}, {"yyzz", "s"}, {"lbts", "list"}, {"pubFlows", "list"}, {"perFlows", "list"}, {"isAutoAdd", "b"}, {"startTime", "D"}, {"endTime", "D"}, {"isTlEtc", "b"}, {"tlzfConfig", "list"}, {"saleEtcs", "list"}, {"etcs", "list"}, {"isOtherEtc", "b"}, {"sumScore", "d"}, {"urcs", "list"}, {"selEtcType", "i"}, {"proCount", "i"}},
}
