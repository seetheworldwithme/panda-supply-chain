package commonrest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) tlCacheAdd(key string) bool {
	s.tlTasksMu.Lock()
	defer s.tlTasksMu.Unlock()
	if s.tlCache == nil {
		s.tlCache = map[string]time.Time{}
	}
	if expiry := s.tlCache[key]; time.Now().Before(expiry) {
		return false
	}
	s.tlCache[key] = time.Now().Add(time.Hour)
	return true
}
func (s *Service) tlCacheDelete(key string) {
	s.tlTasksMu.Lock()
	delete(s.tlCache, key)
	s.tlTasksMu.Unlock()
}

func (s *Service) tlConfirm(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entChainTlzfOrder.validateBody(m) {
		platform.Unknown(w)
		return
	}
	if m["code"] == nil {
		platform.Fail(w, 603, "请求失败")
		return
	}
	// 原 Java 固定从该手机号的 Redis 键取验证码，保留这一业务问题。
	if s.Auth == nil || !s.Auth.VerifyPhoneCode("13994243151", str(m, "code")) {
		platform.Fail(w, 602, "验证码错误")
		return
	}
	if m["reqTimeMinners"] == nil {
		platform.Fail(w, 603, "请求失败")
		return
	}
	key := str(m, "reqTimeMinners")
	if !s.tlCacheAdd(key) {
		platform.Fail(w, 602, "业务处理中，请稍后")
		return
	}
	code, text, err := s.tlConfirmOrders(r, m, key)
	if err != nil {
		platform.Fail(w, 603, "请求失败")
		return
	}
	if code != 200 {
		platform.Fail(w, code, text)
		return
	}
	platform.Success(w, text)
}

func (s *Service) tlConfirmOrders(r *http.Request, m map[string]any, key string) (int, string, error) {
	if m["confirmType"] == nil || m["contractId"] == nil {
		return 0, "", fmt.Errorf("Java 参数拆箱 null")
	}
	d := s.pick(r)
	var available sql.NullString
	if err := d.QueryRow("SELECT sum(rel_operate_amt) FROM chain_tlzf_order WHERE trans_type=1 AND trans_status=1 AND contract_id=?", m["contractId"]).Scan(&available); err != nil {
		return 0, "", err
	}
	if !available.Valid {
		return 0, "", fmt.Errorf("Java 可操作金额为 null")
	}
	can, err := tlAmount(available.String)
	if err != nil {
		return 0, "", err
	}
	need, err := tlAmount(m["amtJf"])
	if err != nil {
		return 0, "", err
	}
	if can.Cmp(need) < 0 {
		return 602, "操作失败，剩余可操作金额" + available.String + "小于解冻/退款金额" + str(m, "amtJf"), nil
	}
	if m["transType"] == nil && m["actTraceFee"] != nil {
		return 0, "", fmt.Errorf("Java transType 为 null")
	}
	fee := big.NewRat(0, 1)
	if m["actTraceFee"] != nil && num(m["transType"]) != 3 {
		fee, err = tlAmount(m["actTraceFee"])
		if err != nil {
			return 0, "", err
		}
	}
	rows, err := queryRows(d, "SELECT * FROM chain_tlzf_order WHERE contract_id=? AND trans_status=1 AND trans_type=1 ORDER BY rel_operate_amt ASC", m["contractId"])
	if err != nil {
		return 0, "", err
	}
	row := pick(m, entChainTlzfOrder.dbcols)
	row["act_trace_fee"], row["batch_no"] = tlMoney(fee), strings.ReplaceAll(uuidV4(), "-", "")
	relFee := new(big.Rat).Set(fee)
	operated := false
	for _, original := range rows {
		order := strings.ReplaceAll(uuidV4(), "-", "")
		for _, k := range strings.Fields("mcht_no sign_num cur_type buy_ect_name acct_num buy_bnk_id buy_bnk_name sale_ect_name sale_acct_num sale_bnk_id amt_tran") {
			row[k] = original[k]
		}
		row["sale_bnk_name"] = original["buy_bnk_name"] // 原 Java 误用买方银行名。
		row["trans_date"], row["mch_order_nbr"], row["mch_trans_nbr"] = time.Now().Format("20060102"), order, order
		row["org_mch_order_nbr"], row["org_mch_trans_nbr"], row["org_trans_date"] = original["mch_order_nbr"], original["mch_trans_nbr"], original["trans_date"]
		remainder, err := tlAmount(original["rel_operate_amt"])
		if err != nil {
			return 0, "", err
		}
		if remainder.Sign() == 0 {
			continue
		}
		operated = true
		oldFee, err := tlAmount(original["rel_trace_fee"])
		if err != nil {
			return 0, "", err
		}
		if m["transType"] == nil {
			return 0, "", fmt.Errorf("Java transType 为 null")
		}
		if remainder.Cmp(oldFee) >= 0 && num(m["transType"]) == 2 {
			target, err := tlAmount(row["amt_jf"])
			if err != nil {
				return 0, "", err
			}
			if oldFee.Cmp(relFee) < 0 || remainder.Cmp(target) < 0 {
				s.tlCacheDelete(key)
				return 602, "操作失败，当前订单可操作金额为：" + str(original, "rel_operate_amt") + "，可分润金额为：" + str(original, "rel_trace_fee"), nil
			}
			oldFee.Sub(oldFee, fee)
			trace, err := tlAmount(original["trace_fee"])
			if err != nil {
				return 0, "", err
			}
			original["rel_trace_fee"], original["act_trace_fee"] = tlMoney(oldFee), tlMoney(new(big.Rat).Sub(trace, oldFee))
		}
		row["act_trace_fee"] = tlMoney(fee)
		relFee.Sub(relFee, fee)
		if remainder.Cmp(need) < 0 {
			need.Sub(need, remainder)
			row["amt_jf"] = tlMoney(remainder)
			if err = s.tlSaveOrder(d, row); err != nil {
				return 0, "", err
			}
			original["rel_operate_amt"] = 0
			if _, err = platform.UpdateByID(d, "chain_tlzf_order", original["id"], insertVals(pick(original, entChainTlzfOrder.dbcols))); err != nil {
				return 0, "", err
			}
			if err = s.tlConfirmRequest(r, d, row, remainder, fee, m["confirmType"]); err != nil {
				return 0, "", err
			}
			continue
		}
		original["rel_operate_amt"] = tlMoney(new(big.Rat).Sub(remainder, need))
		if _, err = platform.UpdateByID(d, "chain_tlzf_order", original["id"], insertVals(pick(original, entChainTlzfOrder.dbcols))); err != nil {
			return 0, "", err
		}
		if num(m["transType"]) == 2 {
			need.Sub(need, fee)
		}
		row["amt_jf"] = tlMoney(need)
		if err = s.tlConfirmRequest(r, d, row, need, fee, m["confirmType"]); err != nil {
			return 0, "", err
		}
		if err = s.tlSaveOrder(d, row); err != nil {
			return 0, "", err
		}
		break
	}
	if !operated {
		return 602, "操作失败，可操作金额为0", nil
	}
	action, kind := "退款", 11
	if num(m["confirmType"]) == 1 {
		action, kind = "解冻", 10
	}
	if err = s.tlOperationLog(r, d, action, str(row, "amt_jf"), kind); err != nil {
		return 0, "", err
	}
	s.startTLTask()
	return 200, "请求成功", nil
}

func (s *Service) tlSaveOrder(d *sql.DB, row map[string]any) error {
	id, err := platform.Insert(d, "chain_tlzf_order", insertVals(row))
	if err == nil {
		row["id"] = id
	}
	return err // 保留多订单拆分时 Java save 同一实体重复主键的问题。
}

func (s *Service) tlConfirmRequest(r *http.Request, d *sql.DB, row map[string]any, amount, fee *big.Rat, confirm any) error {
	req := map[string]any{"_trans_date": row["trans_date"], "req_trace_num": row["mch_order_nbr"], "sign_num": row["sign_num"], "org_req_trace_num": row["org_mch_trans_nbr"], "org_order_num": row["org_mch_order_nbr"], "org_processing_code": "2089", "org_trans_date": row["org_trans_date"], "cfm_processing_code": confirm, "amt_tran": json.Number(tlCents(new(big.Rat).Add(amount, fee))), "resp_url": os.Getenv("TL_ASYNC_ADDRESS"), "addtnl_data1": row["remark1"]}
	if num(confirm) == 3 {
		req["back_mode"] = 0
	}
	if num(confirm) == 1 {
		req["sep_detail"] = tlShare(row, amount, fee, true)
	}
	_, _, err := s.tlCall(r, d, "2090", req, true)
	return err
}
