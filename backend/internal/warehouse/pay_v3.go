package warehouse

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) v3MembershipPay(r *http.Request, m map[string]any, amount, order string) (map[string]any, error) {
	cfg := platform.LoadWxConfig()
	app := "wx0ff59ec3274e0a53"
	path := "/v3/pay/transactions/app"
	fee, ok := new(big.Rat).SetString(amount)
	if !ok {
		return nil, fmt.Errorf("金额格式错误")
	}
	fee.Mul(fee, big.NewRat(100, 1))
	total := int32(new(big.Int).Quo(fee.Num(), fee.Denom()).Int64())
	body := map[string]any{"appid": app, "mchid": cfg.MchID, "description": "会员充值", "out_trade_no": order, "notify_url": platform.Jnull(cfg.NotifyURL) + "/saas-warehouse/busMyBill/appPayNotify", "amount": map[string]any{"currency": "CNY", "total": total}}
	js := m["isJsPay"] == true
	if js {
		// 服务器 Java 的 1 分钱、毫秒订单号和字面 openId 均原样保留。
		app, path = "wx2c78c1ed3e6c6167", "/v3/pay/transactions/jsapi"
		body["appid"], body["mchid"], body["out_trade_no"], body["notify_url"], body["amount"], body["payer"] = app, "1643733121", fmt.Sprint(time.Now().UnixMilli()), "https://scm.coalcloudtech.com/saas/saas-warehouse/busMyBill/payNotify", map[string]any{"currency": "CNY", "total": 1}, map[string]any{"openid": "openId"}
	}
	res, err := platform.WechatV3JSON(r.Context(), "POST", path, body)
	if err != nil {
		return nil, err
	}
	prepay := str(res, "prepay_id")
	if prepay == "" {
		return nil, fmt.Errorf("微信预支付单号缺失")
	}
	stamp, nonce := fmt.Sprint(time.Now().Unix()), uuid32()
	packageValue := prepay
	if js {
		packageValue = "prepay_id=" + prepay
	}
	signature, err := platform.WechatRSASign(app + "\n" + stamp + "\n" + nonce + "\n" + packageValue + "\n")
	if err != nil {
		return nil, err
	}
	if js {
		return map[string]any{"appId": app, "timeStamp": stamp, "nonceStr": nonce, "packageVal": packageValue, "signType": "RSA", "paySign": signature}, nil
	}
	return map[string]any{"appid": app, "partnerId": cfg.MchID, "prepayId": prepay, "packageVal": "Sign=WXPay", "nonceStr": nonce, "timestamp": stamp, "sign": signature}, nil
}
