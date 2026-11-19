package warehouse

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type redMemory struct {
	count   int
	expires time.Time
}

func (s *Service) redCachePut(key string, ttl time.Duration, increment bool) (int, bool) {
	s.redMu.Lock()
	defer s.redMu.Unlock()
	if s.redCache == nil {
		s.redCache = map[string]redMemory{}
	}
	old, ok := s.redCache[key]
	if !time.Now().Before(old.expires) {
		old, ok = redMemory{}, false
	}
	if !increment && ok {
		return old.count, false
	}
	if increment && old.count == 3 {
		return old.count, false
	}
	old.count++
	old.expires = time.Now().Add(ttl)
	s.redCache[key] = old
	return old.count, true
}

func (s *Service) redReceive(r *http.Request, m map[string]any) (int, string, error) {
	if enabled, configured := redEnvelopeEnabled(); !configured || !enabled {
		return 603, "", nil
	}
	now := time.Now()
	today := now.Format("2006-01-02")
	start, err := time.ParseInLocation("2006-01-02T15:04:05", today+"T"+os.Getenv("PANDA_WAREHOUSE_RED_START_TIME"), time.Local)
	if err != nil {
		return 0, "", err
	}
	end, err := time.ParseInLocation("2006-01-02T15:04:05", today+"T"+os.Getenv("PANDA_WAREHOUSE_RED_END_TIME"), time.Local)
	if err != nil {
		return 0, "", err
	}
	if today < "2024-01-11" || now.Before(start) || now.After(end) {
		return 603, "", nil
	}
	d := s.pick(r)
	var maximum, total sql.NullInt64
	if err = d.QueryRow("SELECT red_uplimit_days FROM chain_global_configuration LIMIT 1").Scan(&maximum); err != nil {
		return 0, "", err
	}
	if err = d.QueryRow("SELECT sum(receive_amount) FROM bus_red_envelope_record WHERE TO_DAYS(create_time)=TO_DAYS(now())").Scan(&total); err != nil {
		return 0, "", err
	}
	if total.Valid {
		if !maximum.Valid {
			return 0, "", fmt.Errorf("Java 红包上限为 null")
		}
		if total.Int64 >= maximum.Int64 {
			return 603, "", nil
		}
	}
	clogID := m["clogId"]
	if clogID == nil {
		if clog, ok := m["redMoneyClog"].(map[string]any); ok {
			clogID = clog["id"]
		}
	}
	if clogID == nil || m["receiveAmount"] == nil {
		return 0, "", fmt.Errorf("Java 红包金额或流水为 null")
	}
	if _, ok := s.redCachePut("clog:"+javaStr(clogID), time.Hour, false); !ok {
		return 603, "红包金额重复领取", nil
	}
	clog, err := platform.GetByID(d, "bus_red_money_clog", clogID)
	if err != nil {
		return 0, "", err
	}
	if clog == nil || clog["is_write_off"] == nil {
		return 0, "", fmt.Errorf("Java 红包流水为 null")
	}
	if num(clog["is_write_off"]) != 0 {
		return 603, "红包已提现", nil
	}
	if m["receiveType"] == nil || clog["receive_type"] == nil {
		return 0, "", fmt.Errorf("Java receiveType 为 null")
	}
	if num(m["receiveType"]) != num(clog["receive_type"]) {
		return 603, "红包提现类型不一致", nil
	}
	m["clogId"] = clog["id"]
	kind := int(num(m["receiveType"]))
	uid := m["userId"]
	if uid == nil {
		return 0, "", fmt.Errorf("Java userId 为 null")
	}
	if kind == 1 {
		if _, ok := s.redCachePut("login:"+javaStr(uid), time.Hour, false); !ok {
			return 603, "", nil
		}
	}
	if kind == 6 {
		if _, ok := s.redCachePut("care:"+javaStr(uid), 24*time.Hour, true); !ok {
			return 603, "今日关注奖励已达领取上限", nil
		}
	}
	if kind == 1 || kind == 6 || kind == 11 || kind == 12 {
		q := "SELECT count(*) FROM bus_red_envelope_record WHERE user_id=? AND receive_type=?"
		if kind != 1 {
			q += " AND TO_DAYS(create_time)=TO_DAYS(now())"
		}
		var count int
		if err = d.QueryRow(q, uid, kind).Scan(&count); err != nil {
			return 0, "", err
		}
		if kind == 1 && count > 0 {
			return 603, "", nil
		}
		if kind == 6 && count == 3 {
			return 603, "今日关注奖励已达领取上限", nil
		}
		if (kind == 11 || kind == 12) && count >= 3 {
			return 603, "", nil
		}
	}
	user, err := platform.GetByID(d, "chain_auth_user", uid)
	if err != nil {
		return 0, "", err
	}
	if user == nil {
		return 0, "", fmt.Errorf("Java 红包用户为 null")
	}
	if m["userPhone"] == nil {
		m["userPhone"] = user["phone"]
	}
	environment := "测试环境"
	switch os.Getenv("PANDA_WAREHOUSE_ENVIRONMENT_TYPE") {
	case "2":
		environment = "熊猫煤仓-"
	case "3":
		environment = "准生产环境"
	}
	remark := "首次登录"
	remarks := map[int]string{2: "发布煤种并审核通过", 3: "完善名片资料", 4: "完善企业资料", 5: "分享转发", 6: "关注他人煤种", 7: "意见建议", 8: "价格更新", 9: "完善个人资料", 10: "发布视频并审核通过"}
	if text, ok := remarks[kind]; ok {
		remark = text
	}
	batch := now.Format("20060102150405") + javaStr(uid)
	detail := "D" + batch
	amount := m["receiveAmount"]
	_, err = platform.WechatV3JSON(r.Context(), "POST", "/v3/transfer/batches", map[string]any{"appid": "wx2c78c1ed3e6c6167", "out_batch_no": batch, "batch_name": environment + "营销红包", "batch_remark": environment + remark, "total_amount": amount, "total_num": 1, "transfer_detail_list": []any{map[string]any{"out_detail_no": detail, "transfer_amount": amount, "transfer_remark": remark, "openid": user["openid"]}}})
	if err != nil {
		return 0, "", err
	}
	// 原 Java 不核对客户端金额与生成金额，外部转账成功后再依次写表，保留该顺序与问题。
	if _, err = platform.Insert(d, "bus_red_envelope_record", insertNonNil(pick(m, entBusRedEnvelopeRecord.dbcols))); err != nil {
		return 0, "", err
	}
	if _, err = platform.Insert(d, "bus_red_envelope_batch", map[string]any{"out_batch_no": batch, "batch_name": environment + "营销红包", "batch_remark": remark, "total_amount": amount, "total_num": 1}); err != nil {
		return 0, "", err
	}
	if _, err = platform.Insert(d, "bus_red_envelope_detail", map[string]any{"out_batch_no": batch, "out_detail_no": detail, "transfer_amount": amount, "transfer_remark": remark, "openid": user["openid"]}); err != nil {
		return 0, "", err
	}
	if _, err = platform.UpdateByID(d, "bus_red_money_clog", clog["id"], map[string]any{"is_write_off": true, "write_off_time": now.Format("2006-01-02 15:04:05")}); err != nil {
		return 0, "", err
	}
	return 200, "", nil
}

func insertNonNil(m map[string]any) map[string]any {
	for k, v := range m {
		if v == nil {
			delete(m, k)
		}
	}
	return m
}
