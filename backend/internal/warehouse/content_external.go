package warehouse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func fetchBiddingRows(ctx context.Context, client *http.Client, endpoint string) ([]map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBufferString(`{"page":1,"pageSize":50,"status":2,"requireType":0,"tradeSources":"quote"}`))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("bidding vendor HTTP %d", response.StatusCode)
	}
	var envelope struct {
		Data []map[string]any `json:"data"`
	}
	decoder := json.NewDecoder(io.LimitReader(response.Body, 16<<20))
	decoder.UseNumber()
	if err := decoder.Decode(&envelope); err != nil {
		return nil, err
	}
	if envelope.Data == nil {
		return nil, fmt.Errorf("bidding vendor response has null/missing data")
	}
	return envelope.Data, nil
}

func selectBiddingRows(rows []map[string]any, cutoff *time.Time) ([]map[string]any, error) {
	selected := []map[string]any{}
	for _, row := range rows {
		begin, err := time.Parse("2006-01-02 15:04:05", javaStr(row["preorderBeginTime"]))
		if err != nil {
			return nil, err
		}
		if cutoff == nil || begin.After(*cutoff) {
			selected = append(selected, row)
		}
	}
	return selected, nil
}

func (s *Service) syncBiddingInfo(w http.ResponseWriter, r *http.Request) {
	d := s.pick(r)
	latest, err := queryRows(d, "SELECT preorder_begin_time FROM bus_bidding_info ORDER BY preorder_begin_time DESC LIMIT 1")
	if err != nil {
		platform.Unknown(w)
		return
	}
	var cutoff *time.Time
	if len(latest) > 0 {
		parsed, err := time.Parse("2006-01-02 15:04:05", btsStr(latest[0]["preorder_begin_time"]))
		if err != nil {
			platform.Unknown(w)
			return
		}
		cutoff = &parsed
	}
	endpoint := os.Getenv("BIDDING_SOURCE_URL")
	if endpoint == "" {
		endpoint = "https://www.taomeiwang.com/trade/order/bidList"
	}
	rows, err := fetchBiddingRows(r.Context(), &http.Client{Timeout: 25 * time.Second}, endpoint)
	if err != nil {
		platform.Unknown(w)
		return
	}
	selected, err := selectBiddingRows(rows, cutoff)
	if err != nil {
		platform.Unknown(w)
		return
	}
	tx, err := d.Begin()
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer tx.Rollback()
	for _, row := range selected {
		vals := pick(row, entBusBiddingInfo.dbcols)
		delete(vals, "id")
		if _, err := platform.TxInsert(tx, "bus_bidding_info", vals); err != nil {
			platform.Unknown(w)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

func decimalSubtractExact(a, b any) (string, bool) {
	av, bv := btsStr(a), btsStr(b)
	ar, ok := new(big.Rat).SetString(av)
	if !ok {
		return "", false
	}
	br, ok := new(big.Rat).SetString(bv)
	if !ok {
		return "", false
	}
	scale := 0
	for _, v := range []string{av, bv} {
		if n := strings.Index(v, "."); n >= 0 && len(v)-n-1 > scale {
			scale = len(v) - n - 1
		}
	}
	return new(big.Rat).Sub(ar, br).FloatString(scale), true
}

// The local Java transaction advisor covers openPay even though the method has
// no @Transactional annotation. A missing openid/vendor failure rolls the bill
// INSERT back. Keeping the prepayment work in the transaction also reproduces
// database errors and coupon/config null checks before reaching the adapter.
func (s *Service) openMembershipPay(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entBusMyBill.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	coupon, err := platform.GetByID(d, "bus_user_coupon", m["userCouponId"])
	if err != nil {
		platform.Unknown(w)
		return
	}
	discount := any("0")
	if coupon != nil {
		discount = coupon["coupon_money"]
	}
	// MapUtils includes member_type:null; MP allEq default generates IS NULL.
	cfg := s.getOneRow(w, r, entBusMemberConfig, map[string]any{"member_type": m["memberType"]})
	if cfg == nil && w.Header().Get("Content-Type") != "" {
		return
	}
	tt, ok := bindVal("i", m["transType"])
	if !ok || tt == nil {
		platform.Unknown(w)
		return
	}
	trans := int64(num(tt))
	fee := any("0")
	remark := "预付会员定金"
	if trans >= 1 && trans <= 5 {
		if cfg == nil {
			platform.Unknown(w)
			return
		}
		keys := []string{"advanceDeposit", "finalYearPrice", "finalHalfYearPrice", "finalQuarterYearPrice", "finalMonthPrice"}
		fee = cfg[keys[trans-1]]
		if trans >= 2 {
			renew, ok := bindVal("b", m["isRenew"])
			if !ok || renew == nil {
				platform.Unknown(w)
				return
			}
			term := []string{"年度", "半年度", "季度", "月度"}[trans-2]
			remark = "开通" + term + "会员"
			if renew.(bool) {
				remark = "续费" + term + "会员"
			}
		}
	}
	amount, ok := decimalSubtractExact(fee, discount)
	if !ok {
		platform.Unknown(w)
		return
	}
	tx, err := d.Begin()
	if err != nil {
		platform.Unknown(w)
		return
	}
	defer tx.Rollback()
	orderNb := uuid32()
	vals := map[string]any{"order_nb": orderNb, "user_id": currentUser(r).UserID, "amount": amount, "amount_unit": "元", "status": 1, "bill_remark": remark, "trans_type": tt, "is_renew": m["isRenew"], "user_coupon_id": m["userCouponId"]}
	if _, err := platform.TxInsert(tx, "bus_my_bill", vals); err != nil {
		platform.Unknown(w)
		return
	}
	// 原方法经 TxAdviceInterceptor 全量事务：openid NPE 回滚账单；
	// miniAppPay 内部 catch → null 是正常返回，账单提交且 success(null)。
	u, err := platform.GetByID(d, "chain_auth_user", currentUser(r).UserID)
	if err != nil || u == nil || u["openid"] == nil {
		platform.Unknown(w) // feign getById 后 openid.toString() NPE
		return
	}
	if m["isJsPay"] == true || m["appid"] != nil {
		pay, err := s.v3MembershipPay(r, m, amount, orderNb)
		if err != nil {
			platform.Unknown(w)
			return
		}
		if err = tx.Commit(); err != nil {
			platform.Unknown(w)
			return
		}
		platform.Success(w, pay)
		return
	}
	payObj := s.miniAppPay(r, str(u, "openid"), amount, orderNb, remark, nil)
	if err := tx.Commit(); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, payObj)
}
