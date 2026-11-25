// videos.go BusVideos 16 个 bespoke 端点(toLike/cancelLike/collect/cancelCollect/care/
// cancelCare/forward/play/myLikeVideos/allStatisticsCount/updByUserId/cancelRelation/
// delByVideoId/allLikeCount/allCcollectCount/allCareCount)。
// 对照 BusVideosController/BusVideosServiceImpl/BusVideosMapper,及 BusUserVideoLc/
// BusUserCare/BusVisitorLog 关系表与 saas-common ChainCommonFileServiceImpl.removeFile
// (delByVideoId 的 feign 链)。Java 网关 8765 活体实测样本摘录在各端点注释。要点:
//   - @RequestBody 必填体:零字节/坏 JSON → 999(HTTP 500,body code 999,8765 实测);
//     {} 空对象正常绑定;body "null" 按 {} 处理
//   - 本地运行版实测点赞/收藏链为事务:行缺失/计数列 NULL → 999 且关系行回滚。
//   - 简报称 updByUserId/cancelRelation/cancelCare 空体 999,8765 实测均 200 data:null
//     (MP fill update_time 兜住空 SET;removeByMap/allEq 默认变体 null → IS NULL),照实测实现
package warehouse

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		reg := func(path string, h http.HandlerFunc) {
			mux.HandleFunc("POST /saas-warehouse"+path, authGuard(h))
		}
		reg("/busVideos/toLike", s.vToLike)
		reg("/busVideos/cancelLike", s.vCancelLike)
		reg("/busVideos/collect", s.vCollect)
		reg("/busVideos/cancelCollect", s.vCancelCollect)
		reg("/busVideos/care", s.vCare)
		reg("/busVideos/cancelCare", s.vCancelCare)
		reg("/busVideos/forward", s.vForward)
		reg("/busVideos/play", s.vPlay)
		reg("/busVideos/myLikeVideos", s.vMyLikeVideos)
		reg("/busVideos/allStatisticsCount", s.vAllStatisticsCount)
		reg("/busVideos/updByUserId", s.vUpdByUserId)
		reg("/busVideos/cancelRelation", s.vCancelRelation)
		reg("/busVideos/delByVideoId", s.vDelByVideoId)
		reg("/busVideos/allLikeCount", s.vAllLikeCount)
		reg("/busVideos/allCcollectCount", s.vAllCcollectCount)
		reg("/busVideos/allCareCount", s.vAllCareCount)
	})
}

// ── 共享小件 ──

// vidBody 对照 @RequestBody 必填体:零字节/纯空白/坏 JSON → 999(8765 实测 toLike/
// cancelCare/allLikeCount 零字节均 999;{bad 同 999);body "null" 绑 null 按 {} 处理
func vidBody(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	b, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil || len(bytes.TrimSpace(b)) == 0 {
		platform.Unknown(w)
		return nil, false
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		platform.Unknown(w)
		return nil, false
	}
	if m == nil {
		m = map[string]any{}
	}
	return m, true
}

// vidEntBody vidBody + Jackson 实体绑定校验(8765 实测 {"id":"abc"} → 999)
func vidEntBody(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	m, ok := vidBody(w, r)
	if !ok {
		return nil, false
	}
	if !entBusVideos.validateBody(m) {
		platform.Unknown(w)
		return nil, false
	}
	return m, true
}

// vidLit body 整数字段 → ${} 注入字面(Jackson 已校验整数面;缺失为 "")
func vidLit(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	lit, _ := bindVal("i", v)
	s, _ := lit.(string)
	return s
}

// vidDelEq MP removeByMap/allEq 默认变体:null 值 → IS NULL 条件
func vidDelEq(w http.ResponseWriter, d *sql.DB, table string, cond map[string]any) bool {
	var conds []string
	var args []any
	for _, k := range platform.SortedKeys(cond) {
		if cond[k] == nil {
			conds = append(conds, k+" IS NULL")
			continue
		}
		conds = append(conds, k+" = ?")
		args = append(args, cond[k])
	}
	if _, err := d.Exec("DELETE FROM "+table+" WHERE "+strings.Join(conds, " AND "), args...); err != nil {
		platform.Unknown(w)
		return false
	}
	return true
}

// vidExecUpdate MP update(entity, UpdateWrapper.allEq({key:val})):SET 非空列
// (fill update_time 已带,恒非空 SET —— {} 空实体也 200 的原因,8765 实测),
// WHERE key=val / key IS NULL
func vidExecUpdate(w http.ResponseWriter, d *sql.DB, table string, vals map[string]any, keyCol string, keyVal any) bool {
	cols := platform.SortedKeys(vals)
	var sets []string
	var args []any
	for _, k := range cols {
		sets = append(sets, k+" = ?")
		args = append(args, vals[k])
	}
	cond := keyCol + " IS NULL"
	if keyVal != nil {
		cond = keyCol + " = ?"
		args = append(args, keyVal)
	}
	if _, err := d.Exec("UPDATE "+table+" SET "+strings.Join(sets, ",")+" WHERE "+cond, args...); err != nil {
		platform.Unknown(w)
		return false
	}
	return true
}

// vidBump MP getById 后计数列 ±1 再 updateById:全实体非空列重写 + fill update_time。
// 行缺失/计数列 NULL → Java NPE → 999;减到 0 以下不执行(Java if(likeNum>0) 守卫)
func vidBump(w http.ResponseWriter, d *sql.DB, col string, id any, delta int64) bool {
	row, err := platform.GetByID(d, "bus_videos", id)
	if err != nil || row == nil || row[col] == nil {
		platform.Unknown(w) // getById 空行或计数 NULL → busVideos.getLikeNum() NPE
		return false
	}
	cur := int64(num(row[col]))
	if delta < 0 && cur <= 0 {
		return true
	}
	vals := map[string]any{}
	for _, c := range entBusVideos.dbcols {
		if c != "id" && row[c] != nil {
			vals[c] = row[c]
		}
	}
	vals[col] = cur + delta
	if _, err := platform.UpdateByID(d, "bus_videos", id, vals); err != nil {
		platform.Unknown(w)
		return false
	}
	return true
}

// ── 点赞 / 收藏(toLike/cancelLike/collect/cancelCollect 同构)──

// vLcToggle delta>0 即加(先插 lc 行),delta<0 即减(先删 lc 行,默认变体 null→IS NULL)
func (s *Service) vLcToggle(w http.ResponseWriter, r *http.Request, typ int64, col string, delta int64) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	uid := strconvI64(currentUser(r).UserID)
	// Java service 方法整体事务:save/删 lc 后 getById(video) 缺行 NPE → 回滚 → 行不落库
	// (99977703 实测 java 侧 0 行;999 → 无副作用)
	tx, err := d.Begin()
	if err != nil {
		platform.Unknown(w)
		return
	}
	rollback999 := func() {
		tx.Rollback()
		platform.Unknown(w)
	}
	if delta > 0 {
		// BusUserVideoLc(uid, id, type) save + fill;video_id NULL 照插(事务内,失败回滚)
		if _, err := tx.Exec("INSERT INTO bus_user_video_lc (user_id, video_id, type, create_time, update_time) VALUES (?,?,?,NOW(),NOW())",
			uid, m["id"], typ); err != nil {
			rollback999()
			return
		}
	} else {
		if _, err := tx.Exec("DELETE FROM bus_user_video_lc WHERE user_id=? AND video_id=? AND type=?",
			uid, m["id"], typ); err != nil {
			rollback999()
			return
		}
	}
	var counter sql.NullInt64
	if err := tx.QueryRow("SELECT "+col+" FROM bus_videos WHERE id=?", m["id"]).Scan(&counter); err != nil || !counter.Valid {
		rollback999() // getById 缺行 → NPE → 回滚
		return
	}
	if delta > 0 {
		if _, err := tx.Exec("UPDATE bus_videos SET "+col+"="+col+"+1, update_time=NOW() WHERE id=?", m["id"]); err != nil {
			rollback999()
			return
		}
	} else {
		if _, err := tx.Exec("UPDATE bus_videos SET "+col+"="+col+"-1, update_time=NOW() WHERE id=? AND "+col+">0", m["id"]); err != nil {
			rollback999()
			return
		}
	}
	if err := tx.Commit(); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// vToLike = toLike:插 bus_user_video_lc(type=1)后 like_num+1。
// 8765 实测 {} → 999，插入关系行随事务回滚。
func (s *Service) vToLike(w http.ResponseWriter, r *http.Request) {
	s.vLcToggle(w, r, 1, "like_num", 1)
}

// vCancelLike = cancelLike:删 lc(user_id,video_id,type=1)后 like_num>0 才 -1。
// 8765 实测 {} → 999(delete video_id IS NULL 0 行后 getById 空行 NPE)
func (s *Service) vCancelLike(w http.ResponseWriter, r *http.Request) {
	s.vLcToggle(w, r, 1, "like_num", -1)
}

// vCollect = collect:type=2、collect_num+1。{} 实测 999(同 toLike 链)
func (s *Service) vCollect(w http.ResponseWriter, r *http.Request) {
	s.vLcToggle(w, r, 2, "collect_num", 1)
}

// vCancelCollect = cancelCollect:type=2、collect_num>0 才 -1。{} 实测 999
func (s *Service) vCancelCollect(w http.ResponseWriter, r *http.Request) {
	s.vLcToggle(w, r, 2, "collect_num", -1)
}

// ── 关注 / 取消关注 ──

// vCare = care:feignChainAuthService.getById(careUserId) 取对方 wxFaceImg/nickname
// (Go 直查 chain_auth_user),任一 NULL → userMap.get(...).toString() NPE → 999;
// 否则插 bus_user_care(user_id, care_user_id, care_head_img, care_nickname)+fill。
// 8765 实测 {} → 999(feign data null → NPE)
func (s *Service) vCare(w http.ResponseWriter, r *http.Request) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	urow, err := platform.GetByID(d, "chain_auth_user", m["careUserId"])
	if err != nil || urow == nil {
		platform.Unknown(w)
		return
	}
	face, nick := urow["wx_face_img"], urow["nickname"]
	if face == nil || nick == nil {
		platform.Unknown(w) // 对方头像/昵称 NULL → NPE
		return
	}
	if _, err := platform.Insert(d, "bus_user_care", map[string]any{
		"user_id":       strconvI64(currentUser(r).UserID),
		"care_user_id":  m["careUserId"],
		"care_head_img": face,
		"care_nickname": nick,
	}); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// vCancelCare = cancelCare:jwtInfo 非空即删 bus_user_care(user_id + care_user_id,
// 默认变体 null → IS NULL)。8765 实测 {} → 200 data:null(删 0 行也成功)
func (s *Service) vCancelCare(w http.ResponseWriter, r *http.Request) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	if !vidDelEq(w, s.pick(r), "bus_user_care", map[string]any{
		"user_id":      strconvI64(currentUser(r).UserID),
		"care_user_id": m["careUserId"],
	}) {
		return
	}
	platform.Success(w, nil)
}

// ── 转发 / 播放 ──

// vForward = forward:getById 后 forward_num+1 全行重写。8765 实测 {} → 999(空行 NPE)
func (s *Service) vForward(w http.ResponseWriter, r *http.Request) {
	s.vBumpOnly(w, r, "forward_num")
}

// vPlay = play:getById 后 play_num+1。{} 实测 999
func (s *Service) vPlay(w http.ResponseWriter, r *http.Request) {
	s.vBumpOnly(w, r, "play_num")
}

func (s *Service) vBumpOnly(w http.ResponseWriter, r *http.Request, col string) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	if !vidBump(w, s.pick(r), col, m["id"], 1) {
		return
	}
	platform.Success(w, nil)
}

// ── 计数读端点 ──

// vAllLikeCount = allLikeCount:sum(like_num) by jwt userId,映射 Integer(null → data:null)。
// 8765 实测 hyqcbjhd 无视频 → {"code":200,"data":null};{"id":"5"} 字符串数字照常 200
func (s *Service) vAllLikeCount(w http.ResponseWriter, r *http.Request) {
	s.vSumCount(w, r, "like_num")
}

// vAllCcollectCount = allCcollectCount:sum(collect_num),同上 {} → 200 data:null
func (s *Service) vAllCcollectCount(w http.ResponseWriter, r *http.Request) {
	s.vSumCount(w, r, "collect_num")
}

func (s *Service) vSumCount(w http.ResponseWriter, r *http.Request, col string) {
	var su sql.NullString
	if err := s.pick(r).QueryRow("SELECT sum("+col+") FROM bus_videos WHERE user_id=?",
		strconvI64(currentUser(r).UserID)).Scan(&su); err != nil {
		platform.Unknown(w)
		return
	}
	if !su.Valid || su.String == "" {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, jsonNumberInt64(int64(num(su.String))))
}

// vAllCareCount = allCareCount:service 恒 return null → 200 data:null(8765 实测)
func (s *Service) vAllCareCount(w http.ResponseWriter, r *http.Request) {
	platform.Success(w, nil)
}

// vAllStatisticsCount = allStatisticsCount:userId 取自 body(非 jwt),7 键计数图。
// 8765 实测 userId=2 → {"likeCount":0,"collectCount":0,"fanceCount":0,"headImg":null,
// "hzCount":0,"careCount":0,"nickName":null};{} 或不存在用户(id=1 等) → 999
// (chainAuthUserMapper.selectById 空 → user.getWxFaceImg() NPE)。
// hzMap=sum 查询:全 NULL 单行被 MyBatis 丢弃(returnInstanceForEmptyRow=false)→ hzCount 0;
// likeCount NULL 而 collectCount 非 NULL → parseInt(null) NPE → 999
func (s *Service) vAllStatisticsCount(w http.ResponseWriter, r *http.Request) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	uid := vidLit(m, "userId")
	urow, err := platform.GetByID(d, "chain_auth_user", m["userId"])
	if err != nil || urow == nil {
		platform.Unknown(w)
		return
	}
	// lc 分组计数(type=1 赞 / 2 藏;cout 为 count(1))
	likeCount, collectCount := int64(0), int64(0)
	lrows, err := queryRows(d, "SELECT type, count(1) cout FROM bus_user_video_lc WHERE user_id=? GROUP BY type", uid)
	if err != nil {
		platform.Unknown(w)
		return
	}
	for _, lr := range lrows {
		switch int64(num(lr["type"])) {
		case 1:
			likeCount = int64(num(lr["cout"]))
		case 2:
			collectCount = int64(num(lr["cout"]))
		}
	}
	// 视频表点赞/收藏合计(对应 busVideosMapper.lcCount)
	hzCount := int64(0)
	var likeSum, collSum sql.NullString
	if err := d.QueryRow("SELECT sum(like_num) likeCount, sum(collect_num) collectCount FROM bus_videos WHERE user_id=?", uid).Scan(&likeSum, &collSum); err != nil {
		platform.Unknown(w)
		return
	}
	if likeSum.Valid && likeSum.String != "" {
		hzCount = int64(num(likeSum.String))
	} else if collSum.Valid && collSum.String != "" {
		platform.Unknown(w) // hzMap.get("likeCount").toString() NPE
		return
	}
	var careCount, fanceCount int64
	if err := d.QueryRow("SELECT count(*) FROM bus_user_care WHERE user_id=?", uid).Scan(&careCount); err != nil {
		platform.Unknown(w)
		return
	}
	if err := d.QueryRow("SELECT count(*) FROM bus_user_care WHERE care_user_id=?", uid).Scan(&fanceCount); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, map[string]any{
		"likeCount": jsonNumberInt64(likeCount), "collectCount": jsonNumberInt64(collectCount),
		"careCount": jsonNumberInt64(careCount), "fanceCount": jsonNumberInt64(fanceCount),
		"hzCount": jsonNumberInt64(hzCount), "headImg": urow["wx_face_img"], "nickName": urow["nickname"],
	})
}

// ── 我喜欢/收藏的作品列表 ──

// videosMyLikeSQL BusVideosMapper.myLikeVideos 原文(8765 实测 uid=297 searchStr=""
// page{1,1} → records 93 个驼峰键含 null,datetime→"2024-03-20T09:33:02"、
// decimal(12,2)→Double 最短表示(1540.0)、tinyint(1)→bool、count 子查询→int;
// size=3 → total=40 pages=14)。无 page → 999;无 userId → ${maps.userId} 空串致
// SQL 语法错 → 999;searchStr 缺失绑 NULL → LIKE NULL 空集 200
const videosMyLikeSQL = `SELECT v.user_id userId, v.id as id, v.video_id as videoId,video_name as videoName,like_num as likeNum,collect_num as collectNum,u.phone as phone,v.status as status,v.check_remark as checkRemark,pro_id as proId,v.create_time as createTime,v.update_time as updateTime,pro_describe as proDescribe,u.nickname as nickname,u.wx_face_img as userHeadImg,forward_num as forwardNum,play_num as playNum, warehouse_id as warehouseId,etc_pro_id as etcProId,warehouse_nb as warehouseNb,enterprise_name as enterpriseName,fst_pro_id as fstProId,big_pro_id as bigProId,three_pro_id as threeProId,region as region,product_name as productName,ble_tonnage as bleTonnage,available_tonnage as availableTonnage,min_granularity as minGranularity,max_granularity as maxGranularity,min_ash_content as minAshContent,max_ash_content as maxAshContent,is_show_ash_content as isShowAshContent,min_coal_ash as minCoalAsh,max_coal_ash as maxCoalAsh,min_volatiles as minVolatiles,max_volatiles as maxVolatiles,is_show_volatiles as isShowVolatiles,min_std as minStd,max_std as maxStd,is_show_std as isShowStd,min_internal_sulfur as minInternalSulfur,max_internal_sulfur as maxInternalSulfur,min_water_content as minWaterContent,max_water_content as maxWaterContent,is_show_water_content as isShowWaterContent,min_bond_index as minBondIndex,max_bond_index as maxBondIndex,is_show_bond_index as isShowBondIndex,min_calorific_value as minCalorificValue,max_calorific_value as maxCalorificValue,min_csr as minCsr,max_csr as maxCsr,is_show_csr as isShowCsr,min_y as minY,max_y as maxY,is_show_y as isShowY,min_hgi as minHgi,max_hgi as maxHgi,is_show_hgi as isShowHgi,recovery_analysis as recoveryAnalysis,is_show_recovery_analysis as isShowRecoveryAnalysis,lithofacies as lithofacies,is_show_lithofacies as isShowLithofacies,calorific_value as calorificValue,is_show_calorific_value as isShowCalorificValue,remark as remark,resource_number as resourceNumber,laboratory as laboratory,sampling as sampling,test_number as testNumber,type_of_shipping as typeOfShipping,send_addres as sendAddres,pre_ex_factory_price as preExFactoryPrice,ex_factory_price as exFactoryPrice,plate_price as platePrice,last_test_date as lastTestDate,last_transaction_time as lastTransactionTime,last_tonnage_traded as lastTonnageTraded,last_transaction_price as lastTransactionPrice,ckhy_file_name as ckhyFileName,data_provider as dataProvider,resource_type as resourceType,business_name as businessName,business_phone as businessPhone,notice_time as noticeTime,is_recommende as isRecommende,create_user_id as createUserId,create_phone as createPhone,sharing_times as sharingTimes,view_times as viewTimes,association_id as associationId,img_ids as imgIds,msg_ids as msgIds,is_platform_certification as isPlatformCertification, (SELECT count(1) from bus_user_video_lc where video_id=v.id and type=1 and user_id=${maps.userId}) ifLike,  (SELECT count(1) from bus_user_video_lc where video_id=v.id and type=2 and user_id=${maps.userId}) ifCollect, (SELECT count(1) from bus_user_care where care_user_id=v.user_id and user_id=${maps.userId}) ifCare from bus_videos v join bus_user_video_lc lc on v.id=lc.video_id left join bus_product p on v.pro_id=p.id left join chain_auth_user u on v.user_id=u.id ${searchSql} and CONCAT(IFNULL(v.pro_describe, '')) like #{maps.searchStr} order by v.create_time desc`

// vidAliasCol SELECT 别名 → 源列名(count 子查询别名 ifLike/ifCollect/ifCare 不在表内,
// 走整型兜底)。类型按列名在 bus_videos/bus_product/chain_auth_user 三表内首个命中
// (同名冲突列 status/id/create_time 等两表同型,实测无碍;不同名列 Java 侧本就歧义报错)
var vidAliasCol = map[string]string{
	"id": "id", "userId": "user_id", "videoId": "video_id", "videoName": "video_name", "videoTitle": "video_title",
	"likeNum": "like_num", "collectNum": "collect_num", "phone": "phone", "status": "status",
	"checkRemark": "check_remark", "proId": "pro_id", "createTime": "create_time",
	"updateTime": "update_time", "proDescribe": "pro_describe", "nickname": "nickname",
	"userHeadImg": "wx_face_img", "forwardNum": "forward_num", "playNum": "play_num",
	"warehouseId": "warehouse_id", "etcProId": "etc_pro_id", "warehouseNb": "warehouse_nb",
	"enterpriseName": "enterprise_name", "fstProId": "fst_pro_id", "bigProId": "big_pro_id",
	"threeProId": "three_pro_id", "region": "region", "productName": "product_name",
	"bleTonnage": "ble_tonnage", "availableTonnage": "available_tonnage",
	"minGranularity": "min_granularity", "maxGranularity": "max_granularity",
	"minAshContent": "min_ash_content", "maxAshContent": "max_ash_content",
	"isShowAshContent": "is_show_ash_content", "minCoalAsh": "min_coal_ash",
	"maxCoalAsh": "max_coal_ash", "minVolatiles": "min_volatiles", "maxVolatiles": "max_volatiles",
	"isShowVolatiles": "is_show_volatiles", "minStd": "min_std", "maxStd": "max_std",
	"isShowStd": "is_show_std", "minInternalSulfur": "min_internal_sulfur",
	"maxInternalSulfur": "max_internal_sulfur", "minWaterContent": "min_water_content",
	"maxWaterContent": "max_water_content", "isShowWaterContent": "is_show_water_content",
	"minBondIndex": "min_bond_index", "maxBondIndex": "max_bond_index",
	"isShowBondIndex": "is_show_bond_index", "minCalorificValue": "min_calorific_value",
	"maxCalorificValue": "max_calorific_value", "minCsr": "min_csr", "maxCsr": "max_csr",
	"isShowCsr": "is_show_csr", "minY": "min_y", "maxY": "max_y", "isShowY": "is_show_y",
	"minHgi": "min_hgi", "maxHgi": "max_hgi", "isShowHgi": "is_show_hgi",
	"recoveryAnalysis": "recovery_analysis", "isShowRecoveryAnalysis": "is_show_recovery_analysis",
	"lithofacies": "lithofacies", "isShowLithofacies": "is_show_lithofacies",
	"calorificValue": "calorific_value", "isShowCalorificValue": "is_show_calorific_value",
	"remark": "remark", "resourceNumber": "resource_number", "laboratory": "laboratory",
	"sampling": "sampling", "testNumber": "test_number", "typeOfShipping": "type_of_shipping",
	"sendAddres": "send_addres", "preExFactoryPrice": "pre_ex_factory_price",
	"exFactoryPrice": "ex_factory_price", "platePrice": "plate_price",
	"lastTestDate": "last_test_date", "lastTransactionTime": "last_transaction_time",
	"lastTonnageTraded": "last_tonnage_traded", "lastTransactionPrice": "last_transaction_price",
	"ckhyFileName": "ckhy_file_name", "dataProvider": "data_provider",
	"resourceType": "resource_type", "businessName": "business_name",
	"businessPhone": "business_phone", "noticeTime": "notice_time",
	"isRecommende": "is_recommende", "createUserId": "create_user_id",
	"createPhone": "create_phone", "sharingTimes": "sharing_times", "viewTimes": "view_times",
	"associationId": "association_id", "imgIds": "img_ids", "msgIds": "msg_ids",
	"isPlatformCertification": "is_platform_certification",
	"coalRelName":             "coal_rel_name", "coalRelPhone": "coal_rel_phone",
	"thName": "product_name", "isShowEtcName": "is_show_etc_name", "isShowEtcPhone": "is_show_etc_phone",
	"ifLike": "", "ifCollect": "", "ifCare": "", // count(1) 子查询别名,列名留空走整型兜底
}

// vidDoubleJson DECIMAL 列在 MyBatis Map 行按 Java Double 序列化:最短表示,整数补 .0
// (8765 实测 plate_price decimal(12,2) 值 1540 → "1540.0"、0 → "0.0")
// # ponytail: |x|≥1e7 时 Java Double.toString 走科学计数,煤价/吨位到不了,遇真值再补
func vidDoubleJson(f float64) any {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return json.Number(s)
}

// vidShapeRows 按 information_schema 实型整形 myLikeVideos 的 Map 行(键=SELECT 别名,
// null 键保留 —— callSettersOnNulls=true 实证)
func vidShapeRows(d *sql.DB, rows []map[string]any) []map[string]any {
	types := map[string]colType{}
	for _, t := range []string{"bus_videos", "bus_product", "chain_auth_user", "bus_member_benefits"} {
		for _, ct := range loadMapTypes(d, t) {
			if _, dup := types[ct.name]; !dup {
				types[ct.name] = ct
			}
		}
	}
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		rec := make(map[string]any, len(vidAliasCol))
		for alias, col := range vidAliasCol {
			v, has := row[alias]
			if !has {
				continue // 行内不存在的别名不补(如 listByPage2 无 if* 列,实测)
			}
			if v == nil {
				rec[alias] = nil
				continue
			}
			// MariaDB JDBC exposes these scalar-subquery expressions as numeric
			// values even when their underlying columns are tinyint(1).
			if alias == "isShowEtcName" || alias == "isShowEtcPhone" {
				rec[alias] = jsonNumberInt64(int64(num(v)))
				continue
			}
			ct, ok := types[col]
			if !ok { // count 子查询别名与跨表布尔列(bus_member_benefits)
				if col == "is_show_etc_name" || col == "is_show_etc_phone" {
					rec[alias] = num(v) != 0 // tinyint(1)→Boolean
				} else {
					rec[alias] = jsonNumberInt64(int64(num(v)))
				}
				continue
			}
			switch {
			case ct.dataType == "tinyint" && ct.columnType == "tinyint(1)":
				rec[alias] = num(v) != 0
			case ct.dataType == "decimal":
				rec[alias] = json.Number(strOf(v)) // JDBC DECIMAL is BigDecimal, preserving scale.
			case ct.dataType == "double" || ct.dataType == "float":
				rec[alias] = vidDoubleJson(num(v))
			case ct.dataType == "datetime" || ct.dataType == "timestamp":
				if s, isS := v.(string); isS && len(s) == 19 {
					rec[alias] = s[:10] + "T" + s[11:]
				} else {
					rec[alias] = v
				}
			case ct.dataType == "date":
				if s, isS := v.(string); isS && len(s) == 19 {
					rec[alias] = s[:10]
				} else {
					rec[alias] = v
				}
			case ct.dataType == "int" || ct.dataType == "bigint" || ct.dataType == "smallint" || ct.dataType == "mediumint" || ct.dataType == "tinyint":
				rec[alias] = jsonNumberInt64(int64(num(v)))
			default:
				rec[alias] = v
			}
		}
		out = append(out, rec)
	}
	return out
}

// vidPgNum MP Page 的 long 字段数值绑定面(字符串整数/浮点整数收,非整数/布尔/对象拒)
func vidPgNum(v any) (int64, bool) {
	switch t := v.(type) {
	case string:
		n, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	case float64:
		if t != float64(int64(t)) {
			return 0, false
		}
		return int64(t), true
	}
	return 0, false
}

// vidPgOpt MP Page 的 long 字段三态:缺失/null/"" → not-found(保留默认);可解析 →
// (值,true,true)—— 显式 0 照回显(实测 Java current=0 回显 0、偏移取首页);失败 → 999
func vidPgOpt(pg map[string]any, key string) (int64, bool, bool) {
	v, has := pg[key]
	if !has || v == nil {
		return 0, false, true
	}
	if s, isS := v.(string); isS && s == "" {
		return 0, false, true
	}
	n, ok := vidPgNum(v)
	if !ok {
		return 0, false, false
	}
	return n, true, true
}

// vMyLikeVideos = myLikeVideos:lc 内联 join 视频/产品/用户,Map 行 + MP 分页。
// searchSql=getSql(busVideos,"userId,status")(含 contains("id") 连 id 一起跳过的
// Java 原样怪癖)+ lc.user_id/v.status/lc.type 字面拼接;无 page/无 userId → 999(实测)
func (s *Service) vMyLikeVideos(w http.ResponseWriter, r *http.Request) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	pg, hasPage := m["page"].(map[string]any)
	if !hasPage {
		platform.Unknown(w) // mapper 收 null Page → 999(实测)
		return
	}
	current, size := int64(1), int64(10)
	if n, has, ok2 := vidPgOpt(pg, "current"); !ok2 {
		platform.Unknown(w)
		return
	} else if has {
		current = n
	}
	if n, has, ok2 := vidPgOpt(pg, "size"); !ok2 {
		platform.Unknown(w)
		return
	} else if has {
		size = n
	}
	uidLit := vidLit(m, "userId")
	if uidLit == "" {
		platform.Unknown(w) // ${maps.userId} 缺失 → 空串 → SQL 语法错(实测 999)
		return
	}
	sqlStr, args, ok2 := entBusVideos.getSqlUnstr(m, "userId,status")
	if !ok2 {
		platform.Unknown(w)
		return
	}
	sqlStr += " and lc.user_id=" + uidLit
	if lit := vidLit(m, "status"); lit != "" {
		sqlStr += " and v.status=" + lit
	}
	if lit := vidLit(m, "careType"); lit != "" {
		sqlStr += " and lc.type=" + lit
	}
	var search any
	if v, has := m["searchStr"]; has && v != nil {
		ss, ok3 := bindVal("s", v)
		if !ok3 {
			platform.Unknown(w)
			return
		}
		search = "%" + ss.(string) + "%"
	}
	q := strings.Replace(videosMyLikeSQL, "${searchSql}", sqlStr, 1)
	q = strings.Replace(q, "#{maps.searchStr}", "?", 1)
	q = strings.ReplaceAll(q, "${maps.userId}", uidLit)
	countQ := "SELECT COUNT(*) FROM (" + stripOrderBy(q) + ") cnt"
	baseArgs := append(args, search)
	d := s.pick(r)
	off := (current - 1) * size
	if off < 0 {
		off = 0 // MP Page.getOffset:current<=0 → 0(实测 current=0 返回首页且回显 current=0)
	}
	var total int64
	if size < 0 { // MP:size<0 不分页且不 count(实测 total=0、records 全量)
		if rows, err := queryRows(d, q, baseArgs...); err != nil {
			platform.Unknown(w)
			return
		} else {
			platform.Success(w, ipageOf(vidShapeRows(d, rows), 0, current, size))
			return
		}
	}
	if err := d.QueryRow(countQ, baseArgs...).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	if rows, err := queryRows(d, q+" LIMIT ? OFFSET ?", append(append([]any{}, baseArgs...), size, off)...); err != nil {
		platform.Unknown(w)
		return
	} else {
		platform.Success(w, ipageOf(vidShapeRows(d, rows), total, current, size))
	}
}

// ── 用户信息同步 / 取消关联 / 按文件删 ──

// vUpdByUserId = updByUserId:三连 MP update(entity, UpdateWrapper.allEq)——
// ①bus_videos SET body 非空实体列(id 除外)+fill update_time WHERE user_id=body userId
// ②bus_user_care SET care_head_img/care_nickname(body 有值才 SET)WHERE care_user_id
// ③bus_visitor_log SET target_head_img/target_nickname WHERE visitor_id
// 默认变体 null → IS NULL;{} 实测 200(SET 仅 update_time 匹配 IS NULL 删 0 行)
func (s *Service) vUpdByUserId(w http.ResponseWriter, r *http.Request) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	uv := pick(m, entBusVideos.dbcols)
	vals := map[string]any{}
	for c, v := range uv {
		if c != "id" && v != nil {
			vals[c] = v
		}
	}
	vals["update_time"] = nowSpace()
	if !vidExecUpdate(w, d, "bus_videos", vals, "user_id", m["userId"]) {
		return
	}
	head, nick := m["userHeadImg"], m["nickname"]
	cvals := map[string]any{"update_time": nowSpace()}
	if head != nil {
		cvals["care_head_img"] = head
	}
	if nick != nil {
		cvals["care_nickname"] = nick
	}
	if !vidExecUpdate(w, d, "bus_user_care", cvals, "care_user_id", m["userId"]) {
		return
	}
	gvals := map[string]any{"update_time": nowSpace()}
	if head != nil {
		gvals["target_head_img"] = head
	}
	if nick != nil {
		gvals["target_nickname"] = nick
	}
	if !vidExecUpdate(w, d, "bus_visitor_log", gvals, "visitor_id", m["userId"]) {
		return
	}
	platform.Success(w, nil)
}

// vCancelRelation = cancelRelation:@Update 字面 update bus_videos set pro_id=null
// where id=#{id}(参数化绑 NULL 匹配 0 行)。8765 实测 {} → 200 data:null
func (s *Service) vCancelRelation(w http.ResponseWriter, r *http.Request) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	if _, err := s.pick(r).Exec("UPDATE bus_videos SET pro_id=NULL WHERE id=?", m["id"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// vDelByVideoId = delByVideoId:①删 bus_videos video_id=?(默认变体 null→IS NULL)
// ②feign saas-common removeFile:删 chain_common_file id → 清企业宣传片/官网图引用
// (delEtcVideo/delEtcImg)→ 删磁盘文件与 _img 目录(getById 命中才删)→ 删转换后图片
// 文件(selByFileId 仅删文件不删行)→ 删 source_id 压缩件。videoId 缺失时 removeFile 的
// delEtcVideo ${id} 空串 → SQL 语法错 → 999(8765 实测 {} → 999,①先删 video_id IS NULL
// 行的副作用同 Java)。成功 data "删除成功，id为：<body id,缺省字面 null>"(控制器拼 getId)
func (s *Service) vDelByVideoId(w http.ResponseWriter, r *http.Request) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	vid := m["videoId"]
	if !vidDelEq(w, d, "bus_videos", map[string]any{"video_id": vid}) {
		return
	}
	if vid == nil {
		platform.Unknown(w) // removeFile delEtcVideo ${id}→"" → SQL 语法错
		return
	}
	if _, err := d.Exec("DELETE FROM chain_common_file WHERE id=?", vid); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("UPDATE chain_basic_enterprise SET advertising_video_id=NULL WHERE advertising_video_id=?", vid); err != nil {
		platform.Unknown(w)
		return
	}
	if _, err := d.Exec("UPDATE chain_basic_enterprise SET website_img_id=NULL, website_ced_img_id=NULL WHERE website_img_id=? OR website_ced_img_id=?", vid, vid); err != nil {
		platform.Unknown(w)
		return
	}
	frow, err := platform.GetByID(d, "chain_common_file", vid)
	if err != nil {
		platform.Unknown(w)
		return
	}
	if frow != nil {
		fp, fn := str(frow, "file_path"), str(frow, "file_name")
		if fp != "" && fn != "" {
			dot := strings.Index(fn, ".")
			if dot < 0 {
				platform.Unknown(w) // Java substring 越界
				return
			}
			_ = os.Remove(fp + "/" + fn) // Java file.delete() 忽略结果
			_ = os.RemoveAll(fp + "/" + fn[:dot] + "_img")
		}
	}
	irows, err := queryRows(d, "SELECT file_path, file_name FROM chain_common_img WHERE file_id=?", vid)
	if err != nil {
		platform.Unknown(w)
		return
	}
	for _, ir := range irows {
		_ = os.Remove(str(ir, "file_path") + "/" + str(ir, "file_name"))
	}
	if _, err := d.Exec("DELETE FROM chain_common_file WHERE source_id=?", vid); err != nil {
		platform.Unknown(w)
		return
	}
	idStr := "null"
	if m["id"] != nil {
		idStr = str(m, "id")
	}
	platform.Success(w, "删除成功，id为："+idStr)
}
