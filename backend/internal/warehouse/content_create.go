package warehouse

import (
	"net/http"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		mux.HandleFunc("POST /saas-warehouse/busOpinion/add", authGuard(s.addOpinion))
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
	// Review SMS belongs to the external notification adapter.
	platform.Success(w, nil)
}
