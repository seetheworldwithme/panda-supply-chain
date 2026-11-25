// 招行 openbank 七个入口，按原 Controller/Service/PaymentHelper 保留参数、返回与落库顺序。
// 外部请求使用 busbank_* 私有运行配置；旧 Java 的错误分支与计数问题原样保留。
package capital

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) handlePayment(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	switch {
	case strings.HasSuffix(p, "/payment"):
		s.payPayment(w, r)
	case strings.HasSuffix(p, "/confirmPay"):
		s.payConfirmPay(w, r)
	case strings.HasSuffix(p, "/cancelPay"):
		s.payCancelPay(w, r)
	case strings.HasSuffix(p, "/refund"):
		s.payRefund(w, r)
	case strings.HasSuffix(p, "/orderDetails"):
		s.payOrderDetails(w, r)
	case strings.HasSuffix(p, "/billInfo"):
		s.payBillInfo(w, r)
	case strings.HasSuffix(p, "/transMsgCallBack"):
		s.payTransMsgCallBack(w, r)
	default:
		platform.Unknown(w)
	}
}

// payPayment 付款:payeeEtc 缺键/查无 → 999;merchantNumber 空 → 602;
// 后续外部调用异常 → 999(本地恒达)
func (s *Service) payPayment(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	d := s.pick(r)
	u := currentUser(r)
	systemCode, ok := s.transSysCode(d, u)
	if !ok {
		platform.Fail(w, 999, "失败")
		return
	}
	if _, ok := m["payeeEtc"]; !ok || m["payeeEtc"] == nil {
		platform.Fail(w, 999, "失败") // maps.get("payeeEtc").toString() NPE → fail()
		return
	}
	var mid sql.NullString
	err := d.QueryRow("SELECT merchant_number FROM chain_basic_enterprise WHERE enterprise_name = ?", str(m, "payeeEtc")).Scan(&mid)
	if err != nil {
		platform.Fail(w, 999, "失败") // payeeEtc 查无 → NPE → fail()
		return
	}
	if strings.TrimSpace(mid.String) == "" {
		platform.FailData(w, 602, str(m, "payeeEtc")+"未在招行开户，付款失败", nil)
		return
	}
	order := "DDXX" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	serial := "LSXX" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	dto := bankMerchant(mid.String)
	bankCopy(dto, m, "orderSubtype", "orderAmt")
	for k, v := range map[string]any{"mchOrderNbr": order, "mchTransNbr": serial, "orderType": "PAY", "orderCcyNbr": "10", "orderPostTime": time.Now().Format("20060102150405"), "urlTimeLimit": 30, "payerIssrId": "", "pltRemark": "", "orderRemark": ""} {
		dto[k] = v
	}
	for _, k := range []string{"frontNtfUrl", "asyncNtfUrl"} {
		if v, ok := os.LookupEnv(k); ok {
			dto[k] = v
		}
	}
	res, err := bankPost(r, "/order-pay", map[string]any{"orderInfoPayDto": dto, "orderExtendPayDto": bankOrderExtension()})
	if err != nil {
		platform.Fail(w, 999, "失败")
		return
	}
	head, ok := res["head"].(map[string]any)
	if !ok {
		platform.Fail(w, 999, "失败")
		return
	}
	if head["rtnCode"] != "SUC0000" {
		platform.FailData(w, 602, str(head, "rtnText"), nil)
		return
	}
	// Java 外部下单在参数拆箱与保存之前；后续失败不会撤销银行订单。
	for _, k := range []string{"contractId", "contractCode", "orderAmt", "orderSubtype"} {
		if m[k] == nil {
			platform.Fail(w, 999, "失败")
			return
		}
	}
	contractID, err := strconv.ParseInt(str(m, "contractId"), 10, 32)
	if err != nil {
		platform.Fail(w, 999, "失败")
		return
	}
	buyer, err := platform.GetByID(d, "chain_basic_enterprise", u.CompanyID)
	if err != nil || buyer == nil {
		platform.Fail(w, 999, "失败")
		return
	}
	var sellerID int64
	if err = d.QueryRow("SELECT id FROM chain_basic_enterprise WHERE enterprise_name=?", m["payeeEtc"]).Scan(&sellerID); err != nil {
		platform.Fail(w, 999, "失败")
		return
	}
	if _, err = platform.Insert(d, "chain_capital_transinfo", map[string]any{"system_code": systemCode, "mch_order_nbr": order, "mch_trans_nbr": serial, "contract_id": contractID, "contract_code": str(m, "contractCode"), "pay_money": m["orderAmt"], "status": 0, "pay_unit_id": u.CompanyID, "payee_unit_id": sellerID, "orderSubtype": str(m, "orderSubtype"), "pay_unit_name": buyer["enterprise_name"], "payee_unit_name": str(m, "payeeEtc"), "mch_nbr": mid.String}); err != nil {
		platform.Fail(w, 999, "失败")
		return
	}
	if _, err = d.Exec("UPDATE chain_contract SET is_pay=1 WHERE id=?", str(m, "contractId")); err != nil {
		platform.Fail(w, 999, "失败")
		return
	}
	body, ok := res["body"].(map[string]any)
	if !ok {
		platform.Fail(w, 999, "失败")
		return
	}
	platform.Success(w, body["payUrl"])
}

func bankOrderExtension() map[string]any {
	out := map[string]any{}
	for _, k := range strings.Fields("subMchNbr subMchType subMchName transTerminalNbr accInputType transIp transMac transImei transImsi transIccid transWifiMac deviceId deviceModel deviceLanguage deviceGps deviceSimNbr deviceSimCnt riskScore riskReason mchUserId mchUserRegTime mchUserRegEmail rcvCtry rcvProvince rcvCity rcvAddress goodsType goodsDesc goodsPrice goodsQuantity") {
		out[k] = ""
	}
	out["transTerminalType"], out["transChannelType"] = "07", "01"
	return out
}

// bankCode 保留 Java Redis.get 后 String.equals 的空值与错误码分支。
func (s *Service) bankCode(w http.ResponseWriter, r *http.Request, m map[string]any, confirm bool) bool {
	if m["code"] == nil {
		if confirm {
			platform.FailData(w, 602, "验证码不匹配或已失效", nil)
		} else {
			platform.Fail(w, 999, "失败")
		}
		return false
	}
	if s.Auth == nil {
		platform.Fail(w, 999, "失败")
		return false
	}
	message := s.Auth.PhoneCodeError(strconv.FormatInt(currentUser(r).UserID, 10), str(m, "code"))
	if message == "" {
		return true
	}
	if message == "未发送或验证码已过期" {
		platform.Fail(w, 999, "失败")
		return false
	}
	text := "验证码不匹配"
	if confirm {
		text = "验证码不匹配或已失效"
	}
	platform.FailData(w, 602, text, nil)
	return false
}

func bankConfirmBody(m map[string]any) map[string]any {
	dto := bankMerchant(m["mchNbr"])
	bankCopy(dto, m, "mchOrderNbr", "mchTransNbr", "transType", "transAmt")
	dto["transTime"], dto["pltRemark"], dto["transRemark"] = time.Now().Format("20060102150405"), "", ""
	return map[string]any{"orderCfmDto": dto}
}

func (s *Service) payConfirmPay(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.bankCode(w, r, m, true) {
		return
	}
	if m["id"] == nil || m["mchNbr"] == nil {
		platform.Fail(w, 999, "失败")
		return
	}
	row, err := platform.GetByID(s.pick(r), "chain_capital_transinfo", m["id"])
	if err != nil || row == nil {
		platform.Fail(w, 999, "失败")
		return
	}
	if row["bus_status"] != "WCF" {
		platform.FailData(w, 602, "订单处理中，请稍后重试", nil)
		return
	}
	res, err := bankPost(r, "/order-cfm", bankConfirmBody(m))
	if err != nil {
		platform.Fail(w, 999, "失败")
		return
	}
	head, ok := res["head"].(map[string]any)
	if !ok {
		platform.Fail(w, 999, "失败")
		return
	}
	if head["rtnCode"] != "SUC0000" {
		if head["rtnText"] == nil {
			platform.Fail(w, 999, "失败")
			return
		}
		text := str(head, "rtnText")
		if strings.Contains(text, "BHY2634") {
			text = "银行清算处理中，请24小时后再尝试！"
		}
		platform.FailData(w, 602, text, nil)
		return
	}
	platform.Success(w, "操作成功")
}

func bankRefundBody(m map[string]any) map[string]any {
	dto := bankMerchant(m["mchNbr"])
	bankCopy(dto, m, "oriMchOrderNbr", "oriMchTransNbr", "oriOrderAmt")
	dto["mchOrderNbr"] = "BK" + platform.JnullAny(m["oriMchOrderNbr"])
	dto["mchTransNbr"] = "BK" + platform.JnullAny(m["oriMchTransNbr"])
	dto["orderType"], dto["orderSubtype"], dto["orderCcyNbr"] = "BCK", "BCK", "10"
	if m["oriOrderAmt"] != nil {
		dto["orderAmt"] = m["oriOrderAmt"]
	}
	dto["orderPostTime"], dto["pltRemark"], dto["orderRemark"] = time.Now().Format("20060102150405"), "", ""
	if url, ok := os.LookupEnv("asyncNtfUrl"); ok {
		dto["asyncNtfUrl"] = url
	}
	return map[string]any{"orderInfoBckDto": dto, "orderExtendAddDto": bankOrderExtension()}
}

func (s *Service) payCancelPay(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.bankCode(w, r, m, false) {
		return
	}
	list, ok := m["refundList"].([]any)
	if !ok {
		platform.Fail(w, 999, "失败")
		return
	}
	errors := []string{}
	count := 0
	for _, item := range list {
		row, ok := item.(map[string]any)
		if !ok || row["id"] == nil || row["mchNbr"] == nil || row["payerBankFlag"] == nil {
			platform.Fail(w, 999, "失败")
			return
		}
		stored, err := platform.GetByID(s.pick(r), "chain_capital_transinfo", row["id"])
		if err != nil || stored == nil {
			platform.Fail(w, 999, "失败")
			return
		}
		if stored["bus_status"] != "WCF" {
			errors = append(errors, "订单号为"+platform.JnullAny(row["mchOrderNbr"])+"的记录为非 待确认 状态")
		}
		path := "/order-cfm"
		body := bankConfirmBody(row)
		switch row["payerBankFlag"] {
		case "Y":
		case "N":
			path = "/order-bck"
			body = bankRefundBody(map[string]any{"mchNbr": row["mchNbr"], "oriMchOrderNbr": row["mchOrderNbr"], "oriMchTransNbr": row["mchTransNbr"], "oriOrderAmt": row["transAmt"]})
		default:
			platform.FailData(w, 602, "退款失败，行内外标志为空", nil)
			return
		}
		res, err := bankPost(r, path, body)
		if err != nil {
			platform.Fail(w, 999, "失败")
			return
		}
		head, ok := res["head"].(map[string]any)
		if !ok {
			platform.Fail(w, 999, "失败")
			return
		}
		if head["rtnCode"] != "SUC0000" {
			errors = append(errors, "订单号为"+platform.JnullAny(row["mchOrderNbr"])+"的记录异常，信息为"+platform.JnullAny(head["rtnText"]))
		}
		count++ // Java 银行失败也增加成功数，原样保留。
	}
	platform.Success(w, map[string]any{"resMsg": fmt.Sprintf("退款订单%d笔，成功%d笔，失败%d笔", len(list), count, len(list)-count), "errInfos": errors})
}

func (s *Service) payRefund(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if !s.bankCode(w, r, m, false) {
		return
	}
	list, ok := m["refundList"].([]any)
	if !ok {
		platform.Fail(w, 999, "失败")
		return
	}
	for _, item := range list {
		row, ok := item.(map[string]any)
		if !ok || row["mchNbr"] == nil {
			platform.Fail(w, 999, "失败")
			return
		}
		res, err := bankPost(r, "/order-bck", bankRefundBody(row))
		if err != nil {
			platform.Fail(w, 999, "失败")
			return
		}
		head, ok := res["head"].(map[string]any)
		if !ok {
			platform.Fail(w, 999, "失败")
			return
		}
		if head["rtnCode"] != "SUC0000" {
			platform.FailData(w, 602, str(head, "rtnText"), nil)
			return
		}
		if _, err = s.pick(r).Exec("UPDATE chain_capital_transinfo SET bck_order_nbr=?,bck_trans_nbr=? WHERE mch_order_nbr=?", "BK"+platform.JnullAny(row["oriMchOrderNbr"]), "BK"+platform.JnullAny(row["oriMchTransNbr"]), row["oriMchOrderNbr"]); err != nil {
			platform.Fail(w, 999, "失败")
			return
		}
	}
	platform.Success(w, "操作成功")
}

// payOrderDetails 查询订单:任何异常 → return null → success(null)(controller 直接包)
func (s *Service) payOrderDetails(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	platform.Success(w, bankOrderDetails(r, m))
}

func bankOrderDetails(r *http.Request, m map[string]any) map[string]any {
	if m["mchNbr"] == nil {
		return nil
	}
	dto := bankMerchant(m["mchNbr"])
	bankCopy(dto, m, "mchOrdNbr", "mchTransNbr")
	result, err := bankPost(r, "/order-qry", map[string]any{"orderInfoQryDto": dto})
	if err != nil {
		return nil
	} // Java catch → null，Controller 仍 success(null)。
	return result
}

// payBillInfo 交易对账单:本企业 merchantNumber 空 → 602 "X未在招行开户，查询失败";
// 非空 → 外部调用异常 → 602 "订单查询异常"
func (s *Service) payBillInfo(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	u := currentUser(r)
	d := s.pick(r)
	var name, entName sql.NullString
	if err := d.QueryRow("SELECT merchant_number, enterprise_name FROM chain_basic_enterprise WHERE id = ?", u.CompanyID).Scan(&name, &entName); err != nil {
		platform.FailData(w, 602, "订单查询异常", nil)
		return
	}
	if name.String == "" {
		platform.FailData(w, 602, entName.String+"未在招行开户，查询失败", nil)
		return
	}
	if m["pageNum"] == nil || m["pageSize"] == nil {
		platform.FailData(w, 602, "订单查询异常", nil)
		return
	}
	dto := bankMerchant(name.String)
	bankCopy(dto, m, "orderDate")
	dto["pageNum"], dto["pageSize"] = str(m, "pageNum"), str(m, "pageSize")
	res, err := bankPost(r, "/bill-info", map[string]any{"billQryDto": dto})
	if err != nil {
		platform.FailData(w, 602, "订单查询异常", nil)
		return
	}
	platform.Success(w, res)
}

// payTransMsgCallBack 招行回调:rbody 取自网关塞的 request attribute,本地直连(无网关
// filter)时 Java 侧 attr null → .toString() NPE → 恒 999(8765 实测),照实复刻。
func (s *Service) payTransMsgCallBack(w http.ResponseWriter, r *http.Request) {
	// 必须是 RSA 验签通过的回调:原先只看请求头里有没有 CMB-BodySign,
	// 带有效 JWT 的请求不验签也能把订单改成"已支付"
	if !auth.From(r.Context()).BankSigned {
		platform.Unknown(w)
		return
	}
	m := bodyMap(r)
	// 原 Java transMsgCallBack 由 REQUIRED 事务包裹，后段失败要撤销回调日志以允许重试。
	if err := platform.WithTx(s.pick(r), func(tx *sql.Tx) error { return bankCallback(r, m, tx) }); err != nil {
		platform.Unknown(w)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{"msgNbr": m["msgNbr"], "rtnCode": "SUC0000", "rtnText": "接收成功"})
}

func bankCallback(r *http.Request, m map[string]any, tx *sql.Tx) error {
	status, orderType, flag := str(m, "status"), str(m, "orderType"), str(m, "rtnFlag")
	change, value := false, 0
	d := tx
	if status == "WCF" {
		value, change = 1, true
		_ = bankOrderDetails(r, map[string]any{"mchNbr": m["mchNbr"], "mchOrdNbr": m["mchOrdNbr"], "mchTransNbr": m["mchTransNbr"]})
		if m["payerBankName"] == nil {
			return fmt.Errorf("招商回调处理异常")
		}
		bankFlag := "N"
		if strings.Contains(str(m, "payerBankName"), "招商银行") {
			bankFlag = "Y"
		}
		if _, err := d.Exec("UPDATE chain_capital_transinfo SET payer_bank_flag=? WHERE mch_order_nbr=?", bankFlag, m["mchOrderNbr"]); err != nil {
			return fmt.Errorf("招商回调处理异常")
		}
	}
	if status == "FIN" {
		if orderType == "PAY" {
			switch flag {
			case "S":
				value, change = 2, true
			case "C":
				value, change = 3, true
			case "F", "B":
				value, change = 4, true
			}
		} else if orderType == "BCK" && flag == "S" {
			value, change = 3, true
		}
	}
	if m["mchOrderNbr"] == nil {
		return fmt.Errorf("招商回调处理异常")
	}
	order := str(m, "mchOrderNbr")
	if orderType == "BCK" {
		if len(order) < 2 {
			return fmt.Errorf("招商回调处理异常")
		}
		order = order[2:]
	}
	var maxMessage sql.NullString
	if err := d.QueryRow("SELECT max(msg_nbr) FROM chain_capital_bkcall WHERE mch_order_nbr=?", order).Scan(&maxMessage); err != nil {
		return fmt.Errorf("招商回调处理异常")
	}
	newer := !maxMessage.Valid
	if maxMessage.Valid {
		incoming, err := strconv.ParseInt(str(m, "msgNbr"), 10, 32)
		previous, oldErr := strconv.ParseInt(maxMessage.String, 10, 32)
		if err != nil || oldErr != nil {
			return fmt.Errorf("招商回调处理异常")
		}
		newer = incoming > previous
	}
	if change && newer {
		vals := pick(m, entBkcall.cols...)
		if m["msgtext"] != nil {
			vals["msgText"] = m["msgtext"]
		}
		if _, err := platform.TxInsert(tx, "chain_capital_bkcall", vals); err != nil {
			return fmt.Errorf("招商回调处理异常")
		}
		order = str(m, "mchOrderNbr")
		if orderType != "PAY" {
			order = strings.ReplaceAll(order, "BK", "")
		}
		if _, err := d.Exec("UPDATE chain_capital_transinfo SET bus_status=?,status=? WHERE mch_order_nbr=?", status, value, order); err != nil {
			return fmt.Errorf("招商回调处理异常")
		}
	}
	return nil
}
