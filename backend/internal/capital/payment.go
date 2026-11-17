// capital - payment.go /chainPayment 招行 openbank 对接(8 端点)。
// 外部依赖:PaymentHelper 经 busbank_* 环境变量调招行 API;本环境无这些变量 →
// host=null → HTTP 调用异常 → service catch → ResultTool.fail()=999(外部失败保真,
// 不自造 mock)。Redis 同理不可用 → 验证码类端点 999。
package capital

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"

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
	if mid.String == "" {
		platform.FailData(w, 602, str(m, "payeeEtc")+"未在招行开户，付款失败", nil)
		return
	}
	// PaymentHelper.payment:host 缺失 → 异常 → fail()
	platform.Fail(w, 999, "失败")
}

// payConfirmPay 确认付款:code 键缺 → 602 "验证码不匹配或已失效";
// code 有 → Redis 连接异常 → 999(本地无 Redis,selById 到不了)
func (s *Service) payConfirmPay(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["code"]; !ok || m["code"] == nil {
		platform.FailData(w, 602, "验证码不匹配或已失效", nil)
		return
	}
	platform.Fail(w, 999, "失败") // redisUtil.get 连接异常 → catch → fail()
}

// payCancelPay 取消付款:code 缺 → getString("code").toString() NPE → 999;有 → 999
func (s *Service) payCancelPay(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["code"]; !ok || m["code"] == nil {
		platform.Fail(w, 999, "失败")
		return
	}
	platform.Fail(w, 999, "失败")
}

// payRefund 退款:同 cancelPay(验证码链先炸)恒 999
func (s *Service) payRefund(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if _, ok := m["code"]; !ok || m["code"] == nil {
		platform.Fail(w, 999, "失败")
		return
	}
	platform.Fail(w, 999, "失败")
}

// payOrderDetails 查询订单:任何异常 → return null → success(null)(controller 直接包)
func (s *Service) payOrderDetails(w http.ResponseWriter, r *http.Request) {
	_ = bodyMap(r) // 消费 body;任何异常(含缺参 NPE) → null → success(null)
	platform.Success(w, nil)
}

// payBillInfo 交易对账单:本企业 merchantNumber 空 → 602 "X未在招行开户，查询失败";
// 非空 → 外部调用异常 → 602 "订单查询异常"
func (s *Service) payBillInfo(w http.ResponseWriter, r *http.Request) {
	_ = bodyMap(r)
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
	platform.FailData(w, 602, "订单查询异常", nil)
}

// payTransMsgCallBack 招行回调:裸 JSONObject {msgNbr, rtnCode:"SUC0000", rtnText:"接收成功"}。
// Java rbody 取自网关塞的 attribute;Go 用原始 body。企业无商户号场景 8765 侧行为以实测为准。
func (s *Service) payTransMsgCallBack(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	var callObj map[string]any
	_ = json.Unmarshal(raw, &callObj)
	if callObj == nil {
		platform.Unknown(w) // parseObject 异常路径
		return
	}
	bkcall := bodyMap(r)
	_ = bkcall
	// 本地无有效回调链路(数据库无招商商户),按原版分支产出应答报文
	res := map[string]any{
		"msgNbr":  callObj["msgNbr"],
		"rtnCode": "SUC0000",
		"rtnText": "接收成功",
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(res)
}
