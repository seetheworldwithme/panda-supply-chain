// member2.go 批2 会员/积分/红包 bespoke 端点(selUserMemberInfo/exchangeRecords/
// totoalMoney/searchSignDate/signIn/generateAmount/receiveAmount)。
// 红包 generateAmount/receiveAmount:本地 nacos 无 openRedAmount 配置 → Java NPE 恒 999
// (8765 实测),Go 照抄恒 999;真实配置上线时的行为登记差异交徐先生。
package warehouse

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		reg := func(path string, h http.HandlerFunc) {
			mux.HandleFunc("POST /saas-warehouse"+path, authGuard(h))
		}
		reg("/busUserMember/selUserMemberInfo", s.umSelInfo)
		reg("/busScoreLog/exchangeRecords", s.slExchangeRecords)
		reg("/busRedEnvelopeRecord/totoalMoney", s.rrTotoalMoney)
		reg("/busRedEnvelopeRecord/generateAmount", s.rrGenerateAmount)
		reg("/busRedEnvelopeRecord/receiveAmount", s.rrReceiveAmount)
		reg("/busUserSign/searchSignDate", s.usSearchSignDate)
		reg("/busUserSign/signIn", s.usSignIn)
	})
}

// allEqCondTrue MP allEq(maps) 默认变体:null 值 → IS NULL 条件;非列键 → SQL 错 999
func allEqCondTrue(m map[string]any, cols []string) (string, []any, bool) {
	allow := make(map[string]bool, len(cols))
	for _, c := range cols {
		allow[c] = true
	}
	var conds []string
	var args []any
	for k, v := range m {
		if !allow[k] {
			return "", nil, false
		}
		if v == nil {
			conds = append(conds, k+" IS NULL")
			continue
		}
		conds = append(conds, k+" = ?")
		args = append(args, v)
	}
	if len(conds) == 0 {
		return "", nil, true
	}
	return " WHERE " + strings.Join(conds, " AND "), args, true
}

// getOneRow MP getOne:0 行 nil、>1 行 TooManyResults → 999
func (s *Service) getOneRow(w http.ResponseWriter, r *http.Request, e *ent, m map[string]any) map[string]any {
	where, args, ok := allEqCondTrue(m, e.dbcols)
	if !ok {
		platform.Unknown(w)
		return nil
	}
	rows, cols, err := s.queryEnt(r, "SELECT * FROM "+e.table+where, args...)
	if err != nil {
		platform.Unknown(w)
		return nil
	}
	if len(rows) == 0 {
		return nil
	}
	if len(rows) > 1 {
		platform.Unknown(w) // TooManyResultsException
		return nil
	}
	return e.shapeRows(rows, cols)[0]
}

// umSelInfo = selUserMemberInfo:getOne(allEq(maps) 默认变体);空表/无行 → null
func (s *Service) umSelInfo(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	row := s.getOneRow(w, r, entBusUserMember, m)
	if row == nil && w.Header().Get("Content-Type") != "" {
		return
	}
	platform.Success(w, row)
}

// slExchangeRecords = exchangeRecords:totalScore(黑金卷)+logs(标准分页)
func (s *Service) slExchangeRecords(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusScoreLog.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	uid := strconvI64(currentUser(r).UserID)
	var totalScore int64
	var urow map[string]any
	urows, ucols, err := s.queryEnt(r, "SELECT * FROM bus_user_score WHERE user_id=?", uid)
	if err == nil && len(urows) > 0 && urows[0] != nil {
		urow = entBusUserScore.shapeRows(urows, ucols)[0]
	}
	if urow != nil {
		totalScore = int64(num(urow["hjjNum"]))
	}
	sqlStr, args, ok := entBusScoreLog.getSql(m)
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
	where := "SELECT * FROM bus_score_log " + sqlStr + " and CONCAT(id) like ? order by create_time desc"
	countQ := "SELECT COUNT(*) FROM bus_score_log " + sqlStr + " and CONCAT(id) like ?"
	_, hasPage := m["page"].(map[string]any)
	var logs any
	if !hasPage {
		lr, _, err := s.queryEnt(r, where, append(args, search)...)
		if err != nil {
			platform.Unknown(w)
			return
		}
		logs = entBusScoreLog.shape(lr)
	} else {
		current, size := pageOf(m)
		var total int64
		if err := d.QueryRow(countQ, append(args, search)...).Scan(&total); err != nil {
			platform.Unknown(w)
			return
		}
		off := (current - 1) * size
		lr, _, err := s.queryEnt(r, where+" LIMIT ? OFFSET ?", append(append(args, search), size, off)...)
		if err != nil {
			platform.Unknown(w)
			return
		}
		logs = ipageOf(entBusScoreLog.shape(lr), total, current, size)
	}
	platform.Success(w, map[string]any{"totalScore": jsonNumberInt64(totalScore), "logs": logs})
}

// rrTotoalMoney = totoalMoney:sum(receive_amount) + DATE_FORMAT 字面注入
func (s *Service) rrTotoalMoney(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusRedEnvelopeRecord.validateBody(m) {
		platform.Unknown(w)
		return
	}
	e := entBusRedEnvelopeRecord
	sqlStr, args, ok := e.getSql(m)
	if !ok {
		platform.Unknown(w)
		return
	}
	if v, has := m["startTime"]; has && v != nil {
		d, ok2 := bindVal("D", v)
		if !ok2 || d == nil {
			platform.Unknown(w)
			return
		}
		sqlStr += " and DATE_FORMAT(create_time, '%Y-%m-%d')>= '" + d.(string) + "'"
	}
	if v, has := m["endTime"]; has && v != nil {
		dd, ok2 := bindVal("D", v)
		if !ok2 || dd == nil {
			platform.Unknown(w)
			return
		}
		sqlStr += " and DATE_FORMAT(create_time, '%Y-%m-%d')<= '" + dd.(string) + "'"
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
	var total sql.NullString
	if err := s.pick(r).QueryRow("select sum(receive_amount) from bus_red_envelope_record "+sqlStr+" and CONCAT(id, user_phone) like ?", append(args, search)...).Scan(&total); err != nil {
		platform.Unknown(w)
		return
	}
	v := int64(0)
	if total.Valid && total.String != "" {
		v = int64(num(total.String))
	}
	platform.Success(w, jsonNumberInt64(v))
}

// rrGenerateAmount:本地 nacos 无 openRedAmount → NPE 恒 999(8765 实测;真实配置行为登记差异)
func (s *Service) rrGenerateAmount(w http.ResponseWriter, r *http.Request) {
	platform.Unknown(w)
}

// rrReceiveAmount:同上恒 999
func (s *Service) rrReceiveAmount(w http.ResponseWriter, r *http.Request) {
	platform.Unknown(w)
}

// usSearchSignDate = searchSignDate:当月逐日签到图(已签 "+N"/过去未签 0/未来与今日 "签")
func (s *Service) usSearchSignDate(w http.ResponseWriter, r *http.Request) {
	d := s.pick(r)
	uid := strconvI64(currentUser(r).UserID)
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, -1)
	start := monthStart.Format("2006-01-02")
	end := monthEnd.Format("2006-01-02")
	// 签到映射:sign_time → hjj_num
	srows, _, err := s.queryEnt(r, "SELECT * FROM bus_user_sign WHERE user_id=?", uid)
	if err != nil {
		platform.Unknown(w)
		return
	}
	signMap := map[string]string{}
	for _, sr := range srows {
		st := str(sr, "sign_time")
		if len(st) >= 10 {
			signMap[st[:10]] = str(sr, "hjjNum")
		}
	}
	today := now.Format("2006-01-02")
	signMaps := []any{}
	cur := start
	for {
		day := map[string]any{"date": cur[5:]}
		if v, ok := signMap[cur]; ok {
			day["hjjNum"] = "+" + v
		} else if cur < today {
			day["hjjNum"] = jsonNumberInt64(0)
		} else {
			day["hjjNum"] = "签"
		}
		signMaps = append(signMaps, day)
		if cur == end {
			break
		}
		t, _ := time.Parse("2006-01-02", cur)
		cur = t.AddDate(0, 0, 1).Format("2006-01-02")
	}
	// 累计签到行(实体)/黑金卷/今日是否已签(按 create_time——原版按 create_time 而非 sign_time)
	var signs any
	sr2, sc2, err := s.queryEnt(r, "SELECT * FROM bus_user_sign WHERE user_id=?", uid)
	if err != nil {
		platform.Unknown(w)
		return
	}
	signs = entBusUserSign.shapeRows(sr2, sc2)
	var sumScore int64
	var urow map[string]any
	urows, ucols, err := s.queryEnt(r, "SELECT * FROM bus_user_score WHERE user_id=?", uid)
	if err == nil && len(urows) > 0 && urows[0] != nil {
		urow = entBusUserScore.shapeRows(urows, ucols)[0]
	}
	if urow != nil {
		sumScore = int64(num(urow["hjjNum"]))
	}
	var curCnt int64
	_ = d.QueryRow("SELECT count(*) FROM bus_user_sign WHERE user_id=? AND create_time BETWEEN CONCAT(CURDATE(),' 00:00:00') AND CONCAT(CURDATE(),' 23:59:59')", uid).Scan(&curCnt)
	platform.Success(w, map[string]any{
		"signsDays": signs, "signData": signMaps,
		"sumScore": jsonNumberInt64(sumScore), "isSign": curCnt > 0,
	})
}

// usSignIn = signIn:当日已签静默;昨日连续 +1 否则 1;黑金卷同步加
func (s *Service) usSignIn(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusUserSign.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	uid := strconvI64(currentUser(r).UserID)
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var todayCnt int64
	_ = d.QueryRow("SELECT count(*) FROM bus_user_sign WHERE user_id=? AND sign_time=?", uid, today).Scan(&todayCnt)
	if todayCnt > 0 {
		platform.Success(w, nil)
		return
	}
	var yCnt int64
	var yHjj int64
	_ = d.QueryRow("SELECT count(*) FROM bus_user_sign WHERE user_id=? AND sign_time=?", uid, yesterday).Scan(&yCnt)
	addScore := int64(1)
	if yCnt > 0 {
		_ = d.QueryRow("SELECT hjj_num FROM bus_user_sign WHERE user_id=? AND sign_time=? ORDER BY id DESC LIMIT 1", uid, yesterday).Scan(&yHjj)
		addScore = yHjj + 1
	}
	if _, err := d.Exec("INSERT INTO bus_user_sign (user_id, sign_time, hjj_num) VALUES (?,?,?)", uid, today, strconvI64(addScore)); err != nil {
		platform.Unknown(w)
		return
	}
	var uCount int64
	_ = d.QueryRow("SELECT count(*) FROM bus_user_score WHERE user_id=?", uid).Scan(&uCount)
	if uCount == 0 {
		if _, err := d.Exec("INSERT INTO bus_user_score (user_id, hjj_num) VALUES (?,1)", uid); err != nil {
			platform.Unknown(w)
			return
		}
	} else {
		if _, err := d.Exec("UPDATE bus_user_score SET hjj_num=IFNULL(hjj_num,0)+? WHERE user_id=?", strconvI64(addScore), uid); err != nil {
			platform.Unknown(w)
			return
		}
	}
	platform.Success(w, nil)
}
