package commonrest

import (
	"context"
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

func tlAmount(value any) (*big.Rat, error) {
	if value == nil {
		return nil, fmt.Errorf("Java BigDecimal 为 null")
	}
	r, ok := new(big.Rat).SetString(btsStr(value))
	if !ok {
		return nil, fmt.Errorf("金额格式错误")
	}
	return r, nil
}
func tlCents(value *big.Rat) string {
	r := new(big.Rat).Mul(value, big.NewRat(100, 1))
	return new(big.Int).Quo(r.Num(), r.Denom()).String()
}
func tlMoney(value *big.Rat) json.Number { return json.Number(value.FloatString(2)) }

func (s *Service) tlCall(r *http.Request, d *sql.DB, code string, request map[string]any, verify bool) (map[string]any, string, error) {
	endpoint := os.Getenv("TL_REQUEST_URL")
	if endpoint == "" {
		endpoint = os.Getenv("TL_REQUEST_URL_LEGACY")
	}
	if code == "3001" || code == "3002" {
		endpoint = os.Getenv("TL_QUERY_REQUEST_URL")
	}
	response, req, raw, err := platform.AllinpayPost(r.Context(), endpoint, os.Getenv("TL_INST_ID"), code, request, verify)
	if raw != "" {
		_, saveErr := platform.Insert(d, "chain_tlzf_log", map[string]any{"req_content": req, "res_content": raw})
		if saveErr != nil {
			return nil, raw, saveErr
		}
	}
	return response, raw, err
}

func tlDecoded(value any) any {
	if text, ok := value.(string); ok {
		var out any
		d := json.NewDecoder(strings.NewReader(text))
		d.UseNumber()
		if d.Decode(&out) == nil {
			return out
		}
	}
	return value
}

func (s *Service) tlBalance(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	id := m["companyId"]
	if id == nil {
		id = currentUser(r).CompanyID
	}
	d := s.pick(r)
	row, err := platform.GetByID(d, "chain_basic_enterprise", id)
	if err != nil || row == nil {
		platform.Fail(w, 603, "请求失败")
		return
	}
	res, _, err := s.tlCall(r, d, "3002", map[string]any{"req_trace_num": strings.ReplaceAll(uuidV4(), "-", ""), "sign_num": row["tl_sign_num"], "acct_sub_no": row["tl_acct_sub_no"]}, false)
	if err != nil {
		platform.Fail(w, 603, "请求失败")
		return
	}
	platform.Success(w, tlDecoded(res["qur_rst"]))
}

func (s *Service) tlQueryOrder(r *http.Request, d *sql.DB, row map[string]any) (map[string]any, error) {
	date := str(row, "create_time")
	if len(date) < 10 {
		return nil, fmt.Errorf("Java createTime 为 null")
	}
	res, _, err := s.tlCall(r, d, "3001", map[string]any{"req_trace_num": strings.ReplaceAll(uuidV4(), "-", ""), "org_req_trace_num": row["mch_trans_nbr"], "org_trans_date": strings.ReplaceAll(date[:10], "-", ""), "org_order_num": row["mch_order_nbr"]}, false)
	return res, err
}

func (s *Service) tlTransResult(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entChainTlzfOrder.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	row, err := platform.GetByID(d, "chain_tlzf_order", m["id"])
	if err != nil || row == nil {
		platform.Fail(w, 603, "请求失败")
		return
	}
	res, err := s.tlQueryOrder(r, d, row)
	if err != nil {
		platform.Fail(w, 603, "请求失败")
		return
	}
	for k, v := range res {
		res[k] = tlDecoded(v)
	}
	platform.Success(w, res)
}

func tlShare(row map[string]any, amount, fee *big.Rat, confirm bool) []any {
	merchant, account, bank := "84600000", "171712165", "03050000"
	// 原 Java 平台分账账户是业务常量，允许运行配置覆盖以保持部署数据归属。
	if v := os.Getenv("TL_SHARE_MERCHANT"); v != "" {
		merchant = v
	}
	if v := os.Getenv("TL_SHARE_ACCOUNT"); v != "" {
		account = v
	}
	if v := os.Getenv("TL_SHARE_BANK"); v != "" {
		bank = v
	}
	main := new(big.Rat).Set(amount)
	if !confirm {
		main.Sub(main, fee)
	}
	list := []any{map[string]any{"sep_no": "000", "sep_acct_type": "6", "sep_merch_code": row["mcht_no"], "sep_acct_num": row["sale_acct_num"], "sep_trans_amount": json.Number(tlCents(main)), "sep_bank_id": row["sale_bnk_id"]}}
	cent := tlCents(fee)
	if !confirm || fee.Sign() > 0 {
		list = append(list, map[string]any{"sep_no": "001", "sep_acct_type": "1", "sep_merch_code": merchant, "sep_acct_num": account, "sep_trans_amount": json.Number(cent), "sep_bank_id": bank})
	}
	return list
}

func (s *Service) tlPayMoney(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entChainTlzfOrder.validateBody(m) {
		platform.Unknown(w)
		return
	}
	raw, err := s.tlPay(r, m)
	if err != nil {
		platform.Fail(w, 603, "请求失败")
		return
	}
	platform.Success(w, raw)
}

func (s *Service) tlPay(r *http.Request, m map[string]any) (string, error) {
	d := s.pick(r)
	buyer, err := platform.GetByID(d, "chain_basic_enterprise", m["buyCompanyId"])
	if err != nil {
		return "", err
	}
	if buyer == nil {
		return "", fmt.Errorf("Java buyEtc null")
	}
	seller, err := platform.GetByID(d, "chain_basic_enterprise", m["saleCompanyId"])
	if err != nil {
		return "", err
	}
	if seller == nil {
		return "", fmt.Errorf("Java saleEtc null")
	}
	amount, err := tlAmount(m["amtTran"])
	if err != nil {
		return "", err
	}
	// Java new BigDecimal(0.01) 保留 double 的二进制误差后 HALF_UP 到两位。
	feeText := new(big.Rat).Mul(amount, new(big.Rat).SetFloat64(0.01)).FloatString(2)
	fee, _ := new(big.Rat).SetString(feeText)
	order := strings.ReplaceAll(uuidV4(), "-", "")
	row := pick(m, entChainTlzfOrder.dbcols)
	for k, v := range map[string]any{"mch_order_nbr": order, "mch_trans_nbr": order, "mcht_no": seller["tl_mcht_no"], "sign_num": buyer["tl_sign_num"], "trans_type": 1, "cur_type": "156", "buy_ect_name": buyer["enterprise_name"], "rel_operate_amt": m["amtTran"], "trans_date": time.Now().Format("20060102"), "batch_no": strings.ReplaceAll(uuidV4(), "-", ""), "trace_fee": json.Number(feeText), "act_trace_fee": 0, "rel_trace_fee": json.Number(feeText)} {
		row[k] = v
	}
	req := map[string]any{"_trans_date": row["trans_date"], "req_trace_num": order, "reserve_purchase_type": "0", "pay_mode": "8", "frozen_acct_flag": "9", "bnk_id": row["buy_bnk_id"], "acct_type": "6", "acct_num": row["acct_num"], "cur_type": "156", "amt_tran": json.Number(tlCents(amount)), "prod_import_flag": "0", "order_num": order, "resp_url": os.Getenv("TL_ASYNC_ADDRESS"), "mcht_no": row["mcht_no"], "sign_num": row["sign_num"], "addtnl_data1": row["remark1"], "sep_detail": tlShare(row, amount, fee, false)}
	_, raw, err := s.tlCall(r, d, "2089", req, false)
	if err != nil {
		return "", err
	}
	_, err = platform.Insert(d, "chain_tlzf_order", insertVals(row))
	if err != nil {
		return "", err
	}
	if err = s.tlOperationLog(r, d, "付款", btsStr(m["amtTran"]), 9); err != nil {
		return "", err
	}
	s.startTLTask()
	return raw, nil
}

func (s *Service) tlOperationLog(r *http.Request, d *sql.DB, action, amount string, kind int) error {
	u := currentUser(r)
	_, err := platform.Insert(d, "chain_common_log", map[string]any{"user_id": u.UserID, "user_name": u.Username, "operation": fmt.Sprintf("%s操作了%s，用户id为：%d，金额为%s", u.Username, action, u.UserID, amount), "operation_type": kind})
	return err
}

func (s *Service) tlConfig(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !entChainTlzfConfig.validateBody(m) {
		platform.Unknown(w)
		return
	}
	d := s.pick(r)
	etc, err := platform.GetByID(d, "chain_basic_enterprise", currentUser(r).CompanyID)
	if err != nil || etc == nil {
		platform.Fail(w, 603, "请求失败")
		return
	}
	change := strings.HasSuffix(r.URL.Path, "/tlChangeAcc")
	if !change && strings.TrimSpace(str(etc, "tl_sign_num")) == "" {
		res, _, err := s.tlCall(r, d, "1087", map[string]any{"req_trace_num": strings.ReplaceAll(uuidV4(), "-", ""), "cor_name": etc["legal_person"], "sign_type": "0", "real_name": etc["enterprise_name"], "cer_type": "14", "cer_num": etc["credit_code"]}, true)
		if err != nil {
			platform.Fail(w, 603, "请求失败")
			return
		}
		if res["resp_code"] == "0000" {
			if res["sign_num"] == nil || res["acct_sub_no"] == nil {
				platform.Fail(w, 603, "请求失败")
				return
			}
			etc["tl_sign_num"], etc["tl_acct_sub_no"] = res["sign_num"], res["acct_sub_no"]
			if _, err = platform.UpdateByIDNoFill(d, "chain_basic_enterprise", etc["id"], map[string]any{"tl_sign_num": res["sign_num"], "tl_acct_sub_no": res["acct_sub_no"]}); err != nil {
				platform.Fail(w, 603, "请求失败")
				return
			}
		}
	}
	code := "1182"
	if change {
		code = "1088"
	}
	req := map[string]any{"req_trace_num": strings.ReplaceAll(uuidV4(), "-", ""), "sign_type": "3", "sign_num": etc["tl_sign_num"], "acct_sub_no": etc["tl_acct_sub_no"], "bnk_id": m["bnkId"], "acct_type": "6", "acct_num": m["acctNum"], "cnaps_code": m["cnapsCode"], "hld_name": etc["enterprise_name"], "real_name": etc["enterprise_name"], "cer_type": "14", "cer_num": etc["credit_code"]}
	if change {
		delete(req, "real_name")
		req["tel_num"], req["org_bnk_id"], req["org_acct_type"], req["org_acct_num"] = etc["phone"], m["orgBnkId"], "6", m["orgAcctNum"]
	}
	_, _, err = s.tlCall(r, d, code, req, true)
	if err == nil {
		vals := insertVals(pick(m, entChainTlzfConfig.dbcols))
		if change {
			_, err = platform.UpdateByID(d, "chain_tlzf_config", m["id"], vals)
		} else {
			_, err = platform.Insert(d, "chain_tlzf_config", vals)
		}
	}
	if err != nil {
		platform.Fail(w, 603, "请求失败")
		return
	}
	platform.Success(w, "请求成功")
}

// tlRequest 用于原 MQ 消费者与定时任务：选库由调用者固定，不复用用户请求体。
func tlRequest(ctx context.Context, tenant string) *http.Request {
	r, _ := http.NewRequestWithContext(ctx, "POST", "http://localhost/saas-common/chainTlzfOrder/synchronousStatus", nil)
	r.Header.Set("applytype", tenant)
	return r
}
