package warehouse

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) videoVisit(r *http.Request, m map[string]any) error {
	d := s.pick(r)
	uid := currentUser(r).UserID
	// 原代码只判断 count 非 null，0 也进入，保留旧访客同样触发的情况。
	var todayCount int
	if err := d.QueryRow("SELECT count(*) FROM chain_auth_user WHERE id=? AND TO_DAYS(create_time)=TO_DAYS(now())", uid).Scan(&todayCount); err != nil {
		return err
	}
	if m["shareUserId"] != nil && num(m["shareUserId"]) != float64(uid) {
		var count int
		if err := d.QueryRow("SELECT count(*) FROM bus_pro_log WHERE user_id=? AND share_user_id IS NOT NULL AND operate_type=7 AND operate_obj IN(3,4)", uid).Scan(&count); err != nil {
			return err
		}
		if count == 0 {
			amount, err := randomRedAmount(d, 11)
			if err != nil {
				return err
			}
			var clogID int64
			for i := 0; i < 2; i++ {
				clogID, err = platform.Insert(d, "bus_red_money_clog", map[string]any{"user_id": uid, "receive_type": 11, "create_amount": amount, "is_write_off": false})
				if err != nil {
					return err
				}
			}
			_, _, _ = s.redReceive(r, map[string]any{"userId": m["shareUserId"], "receiveType": 11, "receiveAmount": amount, "clogId": clogID})
		}
	}
	if raw, exists := m["busProLog"]; exists && raw != nil {
		log, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("Java 视频访问日志类型无效")
		}
		_, err := platform.Insert(d, "bus_pro_log", pick(log, entBusProLog.dbcols))
		return err
	}
	return nil
}

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		mux.HandleFunc("POST /saas-warehouse/busOpinion/add", authGuard(s.addOpinion))
		mux.HandleFunc("POST /saas-warehouse/busOpinion/updateById", authGuard(s.updateOpinion))
		mux.HandleFunc("POST /saas-warehouse/busVideos/add", authGuard(s.addVideo))
	})
}

// The Go service resolves the current phone from the user row instead of the
// Java login cache. The caller cannot supply another user's phone or ID.
func (s *Service) addOpinion(w http.ResponseWriter, r *http.Request) {
	body, ok := bodyMapStrict(w, r)
	if !ok {
		return
	}
	if !entBusOpinion.validateBody(body) {
		platform.Unknown(w)
		return
	}
	phone, ok := s.userPhone(r)
	if !ok {
		platform.Unknown(w)
		return
	}
	vals := pick(body, entBusOpinion.dbcols)
	delete(vals, "id")
	vals["user_id"] = currentUser(r).UserID
	vals["phone"] = phone
	if _, err := platform.Insert(s.pick(r), "bus_opinion", vals); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, platform.EntityEcho(nil, entBusRedMoneyClog.keys, nil))
}

func (s *Service) addVideo(w http.ResponseWriter, r *http.Request) {
	body, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	videoID, ok := bindVal("i", body["videoId"])
	if !ok || videoID == nil {
		platform.Fail(w, 602, "视频上传失败，请重试")
		return
	}
	d := s.pick(r)
	uid := currentUser(r).UserID
	user, err := platform.GetByID(d, "chain_auth_user", uid)
	if err != nil || user == nil {
		platform.Unknown(w)
		return
	}
	vals := pick(body, entBusVideos.dbcols)
	delete(vals, "id")
	vals["user_id"] = uid
	if nickname := user["nickname"]; nickname != nil {
		vals["nickname"] = nickname
	}
	tx, err := d.BeginTx(r.Context(), nil)
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer tx.Rollback()
	if _, err = platform.TxInsert(tx, "bus_videos", vals); err != nil {
		platform.Unknown(w)
		return
	}
	_, err = tx.ExecContext(r.Context(), "UPDATE chain_common_file SET owner_id=?, effective=1, association_id=?, relation_id=COALESCE(?, relation_id) WHERE id=?", uid, uuid32(), vals["pro_id"], videoID)
	if err != nil {
		platform.Unknown(w)
		return
	}
	var selfCare int
	if err = tx.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM bus_user_care WHERE user_id=? AND care_user_id=?", uid, uid).Scan(&selfCare); err != nil {
		platform.Unknown(w)
		return
	}
	if selfCare == 0 {
		if _, err = platform.TxInsert(tx, "bus_user_care", map[string]any{"user_id": uid, "care_user_id": uid}); err != nil {
			platform.Unknown(w)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		platform.Unknown(w)
		return
	}
	s.videoAdminSMS(r)
	platform.Success(w, nil)
}

func (s *Service) videoAdminSMS(r *http.Request) {
	rows, err := queryRows(s.pick(r), "SELECT admin_phone FROM bus_config WHERE msg_type=1")
	if err != nil {
		return
	}
	for _, row := range rows {
		platform.SendLegacySMS(r.Context(), javaStr(row["admin_phone"]), nil, "1792052", "熊猫煤仓")
	}
}

func (s *Service) updateVideo(w http.ResponseWriter, r *http.Request) {
	m, ok := vidEntBody(w, r)
	if !ok {
		return
	}
	d := s.pick(r)
	values := pick(m, entBusVideos.dbcols)
	delete(values, "id")
	if _, err := platform.UpdateByID(d, "bus_videos", m["id"], values); err != nil {
		platform.Unknown(w)
		return
	}
	video, err := platform.GetByID(d, "bus_videos", m["id"])
	if err != nil || video == nil {
		platform.Unknown(w)
		return
	}
	if video["status"] == nil {
		platform.Success(w, nil)
		return
	}
	if num(video["status"]) == 0 {
		s.videoAdminSMS(r)
		platform.Success(w, nil)
		return
	}
	user, err := platform.GetByID(d, "chain_auth_user", video["user_id"])
	if err != nil || user == nil {
		platform.Unknown(w)
		return
	}
	text := ""
	if num(video["status"]) == 1 {
		text = "通过"
		if _, fresh := s.redCachePut("video"+javaStr(video["id"]), time.Hour, false); fresh {
			uid := currentUser(r).UserID
			amount, err := randomRedAmount(d, 10)
			if err != nil {
				platform.Unknown(w)
				return
			}
			reward := map[string]any{"user_id": uid, "receive_type": 10, "create_amount": amount, "is_write_off": false}
			if _, err := platform.Insert(d, "bus_red_money_clog", reward); err != nil {
				platform.Unknown(w)
				return
			}
			var count int
			if err := d.QueryRow("SELECT count(*) FROM bus_red_envelope_record WHERE user_id=? AND receive_type IN(2,10) AND TO_DAYS(create_time)=TO_DAYS(now())", uid).Scan(&count); err != nil {
				platform.Unknown(w)
				return
			}
			if count < 3 {
				clogID, err := platform.Insert(d, "bus_red_money_clog", reward)
				if err != nil {
					platform.Unknown(w)
					return
				}
				_, _, _ = s.redReceive(r, map[string]any{"userId": video["user_id"], "userPhone": user["phone"], "receiveType": 10, "receiveAmount": amount, "clogId": clogID})
			}
		}
	} else if num(video["status"]) == 2 {
		text = "拒绝"
	}
	platform.SendLegacySMS(r.Context(), javaStr(user["phone"]), []string{text}, "1794840", "熊猫煤仓")
	platform.Success(w, nil)
}

// 意见审核按 Java 生成两条红包流水，第二条用于提现。
func (s *Service) updateOpinion(w http.ResponseWriter, r *http.Request) {
	body, ok := bodyMapStrict(w, r)
	if !ok {
		return
	}
	if !entBusOpinion.validateBody(body) {
		platform.Unknown(w)
		return
	}
	tx, err := s.pick(r).BeginTx(r.Context(), nil)
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer tx.Rollback()
	var withdrawal map[string]any
	var before sql.NullInt64
	id := body["id"]
	if err := tx.QueryRowContext(r.Context(), "SELECT status FROM bus_opinion WHERE id=? FOR UPDATE", id).Scan(&before); err != nil || !before.Valid {
		platform.Unknown(w)
		return
	}
	vals := pick(body, entBusOpinion.dbcols)
	delete(vals, "id")
	if err := platform.TxUpdateByID(tx, "bus_opinion", id, vals); err != nil {
		platform.Unknown(w)
		return
	}
	if before.Int64 == 0 {
		status, valid := bindVal("i", body["status"])
		if !valid || status == nil {
			platform.Unknown(w)
			return
		}
		if num(status) == 1 {
			const opinionReward = 7
			amount, err := randomRedAmount(tx, opinionReward)
			if err != nil {
				platform.Unknown(w)
				return
			}
			reward := map[string]any{"user_id": currentUser(r).UserID, "receive_type": opinionReward, "create_amount": amount, "is_write_off": false}
			for i := 0; i < 2; i++ {
				clogID, err := platform.TxInsert(tx, "bus_red_money_clog", reward)
				if err != nil {
					platform.Unknown(w)
					return
				}
				// 原 Java 从请求取收款人，生成流水却记在审核人名下。
				withdrawal = map[string]any{"clogId": clogID, "receiveType": opinionReward, "receiveAmount": amount, "userId": body["userId"], "phone": body["phone"]}
			}
		}
	}
	if err := tx.Commit(); err != nil {
		platform.Unknown(w)
		return
	}
	if withdrawal != nil {
		_, _, _ = s.redReceive(r, withdrawal)
	}
	platform.Success(w, nil)
}
