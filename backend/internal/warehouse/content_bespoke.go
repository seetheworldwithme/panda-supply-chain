// content_bespoke.go 七个被路由生成器误当模板直通的写端点,按 Java 字节码补回定制逻辑
// (BusXxxServiceImpl.saveBusXxx / updbusVerifyInfo、BusProductServiceImpl.updUserScore 类型 7 分支、
// BusBusinessCardMapper.clearDefault)。身份一律取登录人:请求体里的 userId / visitorId 不生效;
// 控制器都回不带 data 的 success()。除积分兑换(要求原子)外沿用 Java 的逐步落库,不另加总事务。
package warehouse

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		reg := func(path string, h http.HandlerFunc) {
			mux.HandleFunc("POST /saas-warehouse"+path, authGuard(h))
		}
		reg("/busBusinessCard/add", s.cardAdd)
		reg("/busContacts/add", s.contactsAdd)
		reg("/busEnterpriseIntroduction/add", s.entIntroAdd)
		reg("/busScoreExchangeLog/add", s.scoreExchangeAdd)
		reg("/busVisitorLog/add", s.visitorLogAdd)
		reg("/busVerifyInfo/add", s.verifyInfoAdd)
		reg("/busVerifyInfo/updateById", s.verifyInfoUpdate)
	})
}

// loginUID 定制写端点的登录人:未登录(含仅带签名的匿名请求)对照 Java 里 jwt 为 null 的空指针 → 999
func loginUID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	uid := currentUser(r).UserID
	if uid == 0 {
		platform.Unknown(w)
		return 0, false
	}
	return uid, true
}

// cardAdd = saveBusBusinessCard:clearDefault(该用户全部名片 is_default=2)→ 新名片 is_default=true、
// user_id=登录人 → save。
// Java 在 jwt 为 null 时改用名片自带的 userId:那是登录链(还没有 token)经 Feign 调本接口的分支,
// 这里只对进程内调用(commonrest.callJSON)保留;外部匿名请求回 999。
func (s *Service) cardAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusBusinessCard.validateBody(m) {
		platform.Unknown(w)
		return
	}
	owner := strconvI64(currentUser(r).UserID)
	if owner == "0" {
		v, _ := bindVal("i", m["userId"])
		if v == nil || !auth.IsInternalCall(r.Context()) {
			platform.Unknown(w) // 外部匿名请求,或 card.getUserId().toString() 空指针
			return
		}
		owner = v.(string)
	}
	d := s.pick(r)
	if _, err := d.Exec("UPDATE bus_business_card SET is_default=2 WHERE user_id=?", owner); err != nil {
		platform.Unknown(w)
		return
	}
	vals := pick(m, entBusBusinessCard.dbcols)
	delete(vals, "id")
	vals["is_default"], vals["user_id"] = true, owner
	if _, err := platform.Insert(d, "bus_business_card", vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// contactsAdd = saveBusContacts:登录人对该名片(card_id)还没有联系人才写入;已有则什么都不做、照常回成功。
// card_id 为空按 MP allEq(map) 默认变体匹配 IS NULL。
func (s *Service) contactsAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusContacts.validateBody(m) {
		platform.Unknown(w)
		return
	}
	uid, ok := loginUID(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	var n int
	if err := d.QueryRow("SELECT count(1) FROM bus_contacts WHERE user_id=? AND card_id<=>?", uid, m["cardId"]).Scan(&n); err != nil || n > 0 {
		if err != nil {
			platform.Unknown(w)
		} else {
			platform.Success(w, nil)
		}
		return
	}
	vals := pick(m, entBusContacts.dbcols)
	delete(vals, "id")
	vals["user_id"] = uid
	if _, err := platform.Insert(d, "bus_contacts", vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// entIntroAdd = saveBusEnterpriseIntroduction:user_id=登录人 → save → logoId、imgIds(逗号分隔)对应的
// chain_common_file 置为生效(feign updBthById)。文件号须为整数,否则 999(此时企业介绍已落库,同 Java)。
func (s *Service) entIntroAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusEnterpriseIntroduction.validateBody(m) {
		platform.Unknown(w)
		return
	}
	uid, ok := loginUID(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	vals := pick(m, entBusEnterpriseIntroduction.dbcols)
	delete(vals, "id")
	vals["user_id"] = uid
	if _, err := platform.Insert(d, "bus_enterprise_introduction", vals); err != nil {
		platform.Unknown(w)
		return
	}
	// Java 分两批:logo 一批、imgIds 一批;每批先把编号全部 parseInt 再调 feign
	for _, key := range []string{"logoId", "imgIds"} {
		if m[key] == nil {
			continue
		}
		parts := []string{str(m, key)}
		if key == "imgIds" {
			parts = strings.Split(parts[0], ",")
		}
		var ids []any
		for _, p := range parts {
			n, err := strconv.ParseInt(p, 10, 32) // Integer.parseInt
			if err != nil {
				platform.Unknown(w)
				return
			}
			ids = append(ids, n)
		}
		for _, id := range ids {
			if _, err := d.Exec("UPDATE chain_common_file SET effective=1 WHERE id=?", id); err != nil {
				platform.Unknown(w)
				return
			}
		}
	}
	platform.Success(w, nil)
}

var errInsufficientScore = errors.New("积分不足")

// scoreExchangeAdd = saveBusScoreExchangeLog + updUserScore(类型 7):写兑换记录 → 扣分并写积分流水 →
// 站内消息("您已成功兑换X，3到7日内发货，请留意物流信息")→ 给 bus_config(msg_type=1) 的管理员发普通短信(模板 2246413)。
// 全部落库在一个事务里(Java 无事务,逐步提交)。对 Java 的刻意收紧:
//   - Java 先落兑换记录再判余额,余额不够(hjj_num <= pro_score)照样回成功、发消息、通知发货;这里回 602,不落任何记录。
//   - Java 扣的是请求体 consume_score(与余额校验用的商品 pro_score 不是一个数),并按请求体 user_id 扣分;
//     这里商品名、扣分数取积分商品表,扣的一定是登录人的分。
//   - Java 在当月积分额度 mo_get_sum_score<=0 时直接返回:不扣分,兑换记录却已落库;这里照常扣分。
func (s *Service) scoreExchangeAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusScoreExchangeLog.validateBody(m) {
		platform.Unknown(w)
		return
	}
	uid, ok := loginUID(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	vals := pick(m, entBusScoreExchangeLog.dbcols)
	delete(vals, "id")
	vals["user_id"] = uid
	var proName string
	err := platform.WithTx(d, func(tx *sql.Tx) error {
		var name sql.NullString
		var proScore sql.NullInt64
		if err := tx.QueryRow("SELECT pro_name, pro_score FROM bus_score_pro WHERE id=?", m["goodsId"]).Scan(&name, &proScore); err != nil || !proScore.Valid {
			return errors.New("积分商品不存在或未定价") // getById 无行 / pro_score 为空 → Java 空指针
		}
		proName = name.String
		vals["pro_name"], vals["consume_score"] = proName, proScore.Int64
		if _, err := platform.TxInsert(tx, "bus_score_exchange_log", vals); err != nil {
			return err
		}
		var ruleID, moScore sql.NullInt64
		if err := tx.QueryRow("SELECT id, mo_get_sum_score FROM bus_score_rule").Scan(&ruleID, &moScore); err != nil || !moScore.Valid {
			return errors.New("积分规则缺失") // list().get(0) 越界 / getMoGetSumScore 空指针
		}
		// 行锁:并发兑换串行校验余额
		rows, err := tx.Query("SELECT id, hjj_num, dj_hjj_num FROM bus_user_score WHERE user_id=? FOR UPDATE", uid)
		if err != nil {
			return err
		}
		var scoreID, balance, frozen sql.NullInt64
		n := 0
		for rows.Next() {
			if err := rows.Scan(&scoreID, &balance, &frozen); err != nil {
				return err
			}
			n++
		}
		if err := rows.Err(); err != nil {
			return err
		}
		if n != 1 || !balance.Valid { // 无行空指针 / 多行 MP getOne 抛异常
			return errors.New("用户积分不存在")
		}
		if balance.Int64 <= proScore.Int64 {
			return errInsufficientScore
		}
		if err := platform.TxUpdateByID(tx, "bus_user_score", scoreID.Int64, map[string]any{"hjj_num": balance.Int64 - proScore.Int64}); err != nil {
			return err
		}
		// 积分流水:Java 在类型 7 分支里复用了局部变量,"变动后可用积分"写成了规则类型 2、
		// "规则类型"写成了变动前余额,这里照抄
		flow := map[string]any{"user_id": uid, "user_phone": vals["user_phone"], "content": "兑换商品-" + proName, "score": 0,
			"orgin_available_score": balance.Int64, "final_available_score": 2, "rule_type": balance.Int64, "score_get_type": 7}
		if frozen.Valid {
			flow["orgin_freeze_score"], flow["final_freeze_score"] = frozen.Int64, frozen.Int64
		}
		if _, err := platform.TxInsert(tx, "bus_score_log", flow); err != nil {
			return err
		}
		// Java 扣 0 分后仍 updateById(规则):值不变,只刷新 update_time
		if err := platform.TxUpdateByID(tx, "bus_score_rule", ruleID.Int64, map[string]any{"mo_get_sum_score": moScore.Int64}); err != nil {
			return err
		}
		msgID, err := platform.TxInsert(tx, "bus_msg", map[string]any{"msg_type": 5, "content": "您已成功兑换" + proName + "，3到7日内发货，请留意物流信息", "remark": ""})
		if err != nil {
			return err
		}
		_, err = platform.TxInsert(tx, "bus_user_msg", map[string]any{"user_id": uid, "msg_id": msgID})
		return err
	})
	if err != nil {
		if errors.Is(err, errInsufficientScore) {
			platform.Fail(w, 602, "积分不足，无法兑换")
			return
		}
		platform.Unknown(w)
		return
	}
	s.exchangeAdminSMS(r, proName)
	platform.Success(w, nil)
}

// exchangeAdminSMS 对照 Java 末尾:按 bus_config(msg_type=1) 给管理员手机号发普通短信;没配手机号只记错误日志
func (s *Service) exchangeAdminSMS(r *http.Request, proName string) {
	rows, err := queryRows(s.pick(r), "SELECT admin_phone FROM bus_config WHERE msg_type=1")
	if err != nil {
		return
	}
	if len(rows) == 0 {
		log.Printf("通知平台短信推送失败，请配置相关平台管理员短信通知手机号")
		return
	}
	for _, row := range rows {
		if warehouseNoticeEnabled("SEND_MSG_NOTICE") {
			platform.SendLegacySMS(r.Context(), javaStr(row["admin_phone"]), []string{proName}, "2246413", "")
		}
	}
}

// visitorLogAdd = saveBusVisitorLog:登录人(visitor)对 target_user_id 已有记录 → updateById(请求体实体:
// 字节码传的是入参而不是查到的那条,体里没有 id 时等于空更新);没有 → 用登录人的头像、昵称新建一条。
func (s *Service) visitorLogAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusVisitorLog.validateBody(m) {
		platform.Unknown(w)
		return
	}
	uid, ok := loginUID(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	var n int
	if err := d.QueryRow("SELECT count(1) FROM bus_visitor_log WHERE visitor_id=? AND target_user_id<=>?", uid, m["targetUserId"]).Scan(&n); err != nil || n > 1 {
		platform.Unknown(w) // 多行:MP getOne 抛异常
		return
	}
	vals := pick(m, entBusVisitorLog.dbcols)
	delete(vals, "id")
	delete(vals, "visitor_id") // 归属不能经本接口改写
	var err error
	if n == 0 {
		var face, nick sql.NullString
		if err = d.QueryRow("SELECT wx_face_img, nickname FROM chain_auth_user WHERE id=?", uid).Scan(&face, &nick); err != nil {
			platform.Unknown(w) // selectById 无行 → 空指针
			return
		}
		vals = map[string]any{"visitor_id": uid, "target_user_id": m["targetUserId"]}
		if face.Valid {
			vals["target_head_img"] = face.String
		}
		if nick.Valid {
			vals["target_nickname"] = nick.String
		}
		_, err = platform.Insert(d, "bus_visitor_log", vals)
	} else {
		_, err = platform.UpdateByID(d, "bus_visitor_log", m["id"], vals)
	}
	if err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// filelistOf 核实信息 filelist(@Many 载体,模板校验不覆盖):不是数组、或元素不是附件对象 → Jackson 绑定失败;
// null 元素 Jackson 收得下,留到挂附件时空指针(同 Java:主记录此时已写入)
func filelistOf(m map[string]any) ([]any, bool) {
	if m["filelist"] == nil {
		return nil, true
	}
	arr, ok := m["filelist"].([]any)
	if !ok {
		return nil, false
	}
	for _, it := range arr {
		f, isObj := it.(map[string]any)
		if it != nil && (!isObj || !fileEnt.validateBody(f)) {
			return nil, false
		}
	}
	return arr, true
}

var errNotFileOwner = errors.New("附件不属于当前用户")

// checkFileOwners 核实附件只能是自己上传的(或尚未认领的)文件:原先按附件 id 直接挂到登录人名下,
// 等于可以把别人的附件改挂给自己。在写主记录之前校验,避免留下没有附件的核实记录。
func checkFileOwners(d *sql.DB, files []any, uid int64) error {
	for _, it := range files {
		f, _ := it.(map[string]any)
		if f == nil {
			return errors.New("核实附件为空") // Java 对 null 元素 setOwnerType 空指针
		}
		var owner sql.NullInt64
		if err := d.QueryRow("SELECT owner_id FROM chain_common_file WHERE id=?", f["id"]).Scan(&owner); err != nil {
			return errors.New("核实附件不存在")
		}
		if owner.Valid && owner.Int64 != uid {
			return errNotFileOwner
		}
	}
	return nil
}

// linkVerifyFiles 对照 feign updBthById 的前置赋值 + updateBatchById(非空子集):每项 owner_type=2、
// owner_id=登录人、effective=true、association_id=核实信息关联号(没给就不改这一列)。
// 沿用 commonrest 的约束:file_path 不接受客户端改写,file_name 只接受纯文件名。
func (s *Service) linkVerifyFiles(d *sql.DB, files []any, uid int64, association any) error {
	for _, it := range files {
		f, _ := it.(map[string]any)
		if f == nil {
			return errors.New("核实附件为空") // Java 对 null 元素 setOwnerType 空指针
		}
		vals := pick(f, fileEnt.dbcols)
		delete(vals, "id")
		delete(vals, "file_path")
		if v, ok := vals["file_name"].(string); ok && (v == "" || v != filepath.Base(v) || v == ".." || v == ".") {
			delete(vals, "file_name")
		}
		vals["owner_type"], vals["owner_id"], vals["effective"] = "2", uid, true
		if association != nil {
			vals["association_id"] = association
		}
		if _, err := platform.UpdateByIDNoFill(d, "chain_common_file", f["id"], vals); err != nil {
			return err
		}
	}
	return nil
}

// verifyInfoAdd = saveBusVerifyInfo:association_id=UUID → save → filelist 挂到登录人名下并生效 →
// 煤种 is_real_check=true(Java 是把查出的整行回写,这里只改这一列:效果相同,且不会盖掉并发修改)。
func (s *Service) verifyInfoAdd(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	files, ok := filelistOf(m)
	if !ok || !entBusVerifyInfo.validateBody(m) {
		platform.Unknown(w)
		return
	}
	uid, ok := loginUID(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	if err := checkFileOwners(d, files, uid); err != nil {
		bespokeFileError(w, err)
		return
	}
	association := uuid32()
	vals := pick(m, entBusVerifyInfo.dbcols)
	delete(vals, "id")
	vals["association_id"] = association
	if _, err := platform.Insert(d, "bus_verify_info", vals); err != nil {
		platform.Unknown(w)
		return
	}
	if err := s.linkVerifyFiles(d, files, uid, association); err != nil {
		platform.Unknown(w)
		return
	}
	var one int
	if err := d.QueryRow("SELECT 1 FROM bus_product WHERE id=?", m["proId"]).Scan(&one); err != nil {
		platform.Unknown(w) // getById 无行 → setIsRealCheck 空指针(附件此前已更新,同 Java)
		return
	}
	if _, err := platform.UpdateByID(d, "bus_product", m["proId"], map[string]any{"is_real_check": true}); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// verifyInfoUpdate = updbusVerifyInfo:updateById(请求体实体)→ filelist 同上,association_id 取请求体
func (s *Service) verifyInfoUpdate(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	files, ok := filelistOf(m)
	if !ok || !entBusVerifyInfo.validateBody(m) {
		platform.Unknown(w)
		return
	}
	uid, ok := loginUID(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	if err := checkFileOwners(d, files, uid); err != nil {
		bespokeFileError(w, err)
		return
	}
	vals := pick(m, entBusVerifyInfo.dbcols)
	delete(vals, "id")
	if _, err := platform.UpdateByID(d, "bus_verify_info", m["id"], vals); err != nil {
		platform.Unknown(w)
		return
	}
	if err := s.linkVerifyFiles(d, files, uid, m["associationId"]); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

// bespokeFileError 附件归属不合法回 602(业务失败),其余按 Java 空指针 / SQL 异常回 999
func bespokeFileError(w http.ResponseWriter, err error) {
	if errors.Is(err, errNotFileOwner) {
		platform.Fail(w, 602, "无权限操作")
		return
	}
	platform.Unknown(w)
}
