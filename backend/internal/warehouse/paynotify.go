package warehouse

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// 微信支付回调。payNotify 走 v2 XML+MD5 真实验签；appPayNotify 走 v3 平台证书验签+AES-GCM 解密。
// 原 Java 不校验金额且无幂等防护（重复通知会再次延长会员），属原实现问题，差异报告单列。
// 本地 Java appPayNotify 因商户私钥文件缺失类初始化失败恒 500，是环境缺陷不是业务行为；
// Go 按生产语义实现（验签失败→200 空体），凭据缺失同样走失败分支，不伪造成功。

func init() {
	addExtra(func(mux *http.ServeMux, s *Service) {
		// 网关 JwtAuthorizationTokenFilter 对 paynotify 自建身份，不走 authGuard。
		mux.HandleFunc("POST /saas-warehouse/busMyBill/payNotify", s.billPayNotify)
		mux.HandleFunc("GET /saas-warehouse/busMyBill/payNotify", s.billPayNotify)
		mux.HandleFunc("POST /saas-warehouse/busMyBill/appPayNotify", s.billAppPayNotify)
		mux.HandleFunc("GET /saas-warehouse/busMyBill/appPayNotify", s.billAppPayNotify)
	})
}

// wxXMLToMap 对照 WxPayKit.xmlToMap：空串/非法 XML 抛异常（Java 500），合法空文档返回空 map。
func wxXMLToMap(text string) (map[string]string, error) {
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("empty xml")
	}
	dec := xml.NewDecoder(strings.NewReader(text))
	out := make(map[string]string)
	var key string
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			key = t.Name.Local
		case xml.CharData:
			if key != "" && key != "xml" {
				out[key] = string(t)
			}
		case xml.EndElement:
			key = ""
		}
	}
	return out, nil
}

// wxVerifyNotify 对照 WxPayKit.verifyNotify(MD5)：空值与 sign 不参与，键升序拼接后 &key=密钥 取 MD5 大写。
func wxVerifyNotify(params map[string]string, key string) bool {
	return params["sign"] != "" && wxLinkSign(params, key) == params["sign"]
}

func (s *Service) billPayNotify(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 公网回调:微信通知体很小,封顶防超大包打 OOM
	if err != nil {
		platform.Unknown(w)
		return
	}
	params, err := wxXMLToMap(string(body))
	if err != nil {
		platform.Unknown(w)
		return
	}
	key := platform.LoadWxConfig().MchSecret
	if key == "" {
		// 没配商户密钥就无法验签:不能退化成字面 "null" 去比对,否则伪造的通知会被当成已支付
		w.WriteHeader(http.StatusOK) // 与验签失败一致:Java return null → 200 空体
		return
	}
	if !wxVerifyNotify(params, key) || params["return_code"] != "SUCCESS" {
		w.WriteHeader(http.StatusOK) // Java return null → 200 空体
		return
	}
	s.billPaySuccess(w, r, params["out_trade_no"])
}

func (s *Service) billAppPayNotify(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 公网回调:微信通知体很小,封顶防超大包打 OOM
	if err != nil {
		platform.Unknown(w)
		return
	}
	outTradeNo, err := wxV3NotifyOutTradeNo(r, body, platform.LoadWxConfig())
	if err != nil {
		w.WriteHeader(http.StatusOK) // getTransAction catch → null → 200 空体
		return
	}
	s.billPaySuccess(w, r, outTradeNo)
}

// wxV3NotifyOutTradeNo 对照 NotificationParser.parse：平台证书验签 + AEAD_AES_256_GCM 解密 resource。
func wxV3NotifyOutTradeNo(r *http.Request, body []byte, cfg platform.WxLocalConfig) (string, error) {
	if cfg.APIV3Key == "" {
		return "", errors.New("微信支付 v3 凭据未配置")
	}
	cert, err := platform.WechatPlatformCertificate(r.Context(), cfg, r.Header.Get("wechatpay-serial"))
	if err != nil {
		return "", err
	}
	message := r.Header.Get("wechatpay-timestamp") + "\n" + r.Header.Get("wechatpay-nonce") + "\n" + string(body) + "\n"
	sig, err := base64.StdEncoding.DecodeString(r.Header.Get("wechatpay-signature"))
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256([]byte(message))
	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok || rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sig) != nil {
		return "", errors.New("验签失败")
	}
	var notify struct {
		Resource struct {
			Algorithm      string `json:"algorithm"`
			Ciphertext     string `json:"ciphertext"`
			AssociatedData string `json:"associated_data"`
			Nonce          string `json:"nonce"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(body, &notify); err != nil {
		return "", err
	}
	if notify.Resource.Algorithm != "AEAD_AES_256_GCM" {
		return "", errors.New("不支持的解密算法")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(notify.Resource.Ciphertext)
	if err != nil {
		return "", err
	}
	blockCipher, err := aes.NewCipher([]byte(cfg.APIV3Key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}
	if len(notify.Resource.Nonce) != gcm.NonceSize() {
		return "", errors.New("微信通知nonce长度无效")
	}
	plain, err := gcm.Open(nil, []byte(notify.Resource.Nonce), ciphertext, []byte(notify.Resource.AssociatedData))
	if err != nil {
		return "", err
	}
	var trans struct {
		OutTradeNo string `json:"out_trade_no"`
	}
	if err := json.Unmarshal(plain, &trans); err != nil {
		return "", err
	}
	return trans.OutTradeNo, nil
}

// billPaySuccess 验签通过后的落库链：订单 status=2 → 会员副作用 → SUCCESS XML。
// 原服务方法经 TxAdviceInterceptor("*"→REQUIRED) 包裹，后段 NPE 时前段更新整体回滚。
func (s *Service) billPaySuccess(w http.ResponseWriter, r *http.Request, outTradeNo string) {
	d := s.pick(r)
	q, args := "SELECT * FROM bus_my_bill WHERE order_nb=?", []any{outTradeNo}
	if outTradeNo == "" {
		q, args = "SELECT * FROM bus_my_bill WHERE order_nb IS NULL", nil
	}
	bills, err := queryRows(d, q, args...)
	if err != nil || len(bills) != 1 {
		platform.Unknown(w) // getOne 0 行 NPE / 多行 TooManyResults
		return
	}
	bill := bills[0]
	err = platform.WithTx(d, func(tx *sql.Tx) error {
		if err := platform.TxUpdateByID(tx, "bus_my_bill", bill["id"], map[string]any{"status": 2}); err != nil {
			return err
		}
		return s.billUpdUserInfo(r, tx, bill)
	})
	if err != nil {
		platform.Unknown(w)
		return
	}
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	_, _ = w.Write([]byte("<xml><return_code>SUCCESS</return_code><return_msg>OK</return_msg></xml>"))
}

// billUpdUserInfo 对照同名方法：删优惠券；transType=1 置 member_type=2，
// 其余置 member_type=3 并按原到期时间续期（memberType 刚置 3，Java 条件恒真）。
// 开通短信异常被 Java 吞掉、Redis member 缓存 Go 不依赖，二者无契约影响（见差异报告）。
func (s *Service) billUpdUserInfo(r *http.Request, tx *sql.Tx, bill map[string]any) error {
	if _, err := tx.Exec("DELETE FROM bus_user_coupon WHERE id=?", bill["user_coupon_id"]); err != nil {
		return err
	}
	transType := num(bill["trans_type"])
	if transType == 1 {
		_, err := tx.Exec("UPDATE chain_auth_user SET member_type=2 WHERE id=?", bill["user_id"])
		return err
	}
	user, err := platform.TxGetByID(tx, "chain_auth_user", bill["user_id"])
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user missing") // user.setMemberType NPE
	}
	var newExp any
	if transType >= 2 && transType <= 5 {
		expText, ok := user["expiration_time"].(string)
		if !ok || expText == "" {
			return errors.New("expirationTime null") // plusYears/plusMonths NPE，更新尚未发生
		}
		exp, err := time.ParseInLocation("2006-01-02 15:04:05", expText, time.Local)
		if err != nil {
			return err
		}
		months := map[float64]int{2: 12, 3: 6, 4: 3, 5: 1}[transType]
		newExp = plusMonthsClamped(exp, months).Format("2006-01-02 15:04:05")
	}
	vals := map[string]any{"member_type": 3}
	if newExp != nil {
		vals["expiration_time"] = newExp
	}
	if err := platform.TxUpdateByIDNoFill(tx, "chain_auth_user", user["id"], vals); err != nil {
		return err
	}
	if newExp == nil {
		return errors.New("expirationTime null") // transType 越界：更新后 Date.from(null) NPE
	}
	if err := auth.ScheduleMemberExpiryText(r.Context(), javaStr(user["id"]), fmt.Sprint(newExp)); err != nil {
		return err
	}
	platform.SendLegacySMS(r.Context(), str(user, "phone"), []string{"", fmt.Sprint(newExp)}, "2193101", "")
	return nil
}

// plusMonthsClamped LocalDateTime.plusMonths：月末溢出钳制到目标月最后一天（Go AddDate 会进位）。
func plusMonthsClamped(t time.Time, months int) time.Time {
	y, m, d := t.Date()
	first := time.Date(y, m+time.Month(months), 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	if last := first.AddDate(0, 1, -1).Day(); d > last {
		d = last
	}
	return time.Date(first.Year(), first.Month(), d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// miniAppPay 对照 BusMyBillServiceImpl.miniAppPay：v2 统一下单 + 调起支付签名。
// 任何异常/失败 → null（Java catch 同值），调用方 success(null) 且事务提交。
// 商户凭据走私有运行配置；本地缺失时微信侧拒绝下单，与 Java 本地行为一致。
func (s *Service) miniAppPay(r *http.Request, openID, amount, orderNb, billRemark string, appid any) map[string]any {
	cfg := platform.LoadWxConfig()
	host, err := os.Hostname()
	if err != nil {
		return nil
	}
	hosts, err := net.LookupHost(host)
	if err != nil || len(hosts) == 0 {
		return nil // InetAddress.getLocalHost 异常 → catch → null
	}
	ip := hosts[0]
	tradeType := "APP"
	appID := ""
	if appid == nil {
		rows, err := queryRows(s.pick(r), "SELECT * FROM chain_auth_chat_authorization WHERE app_identification='cloudWare'")
		if err != nil || len(rows) != 1 || rows[0]["appid"] == nil {
			return nil
		}
		appID = str(rows[0], "appid")
		tradeType = "JSAPI"
	} else {
		appID = fmt.Sprint(appid)
	}
	fee, ok := new(big.Rat).SetString(amount)
	if !ok {
		return nil
	}
	fee.Mul(fee, big.NewRat(100, 1))
	totalFee := new(big.Int).Quo(fee.Num(), fee.Denom()).String()
	params := map[string]string{
		"nonce_str":        uuid32(),
		"mch_id":           cfg.MchID,
		"body":             "会员充值",
		"attach":           billRemark,
		"notify_url":       platform.Jnull(cfg.NotifyURL) + "/saas-warehouse/busMyBill/payNotify",
		"spbill_create_ip": ip,
		"appid":            appID,
		"out_trade_no":     orderNb,
		"total_fee":        totalFee,
		"trade_type":       tradeType,
	}
	if appid == nil {
		params["openid"] = openID
	}
	sign := wxLinkSign(params, platform.Jnull(cfg.MchSecret))
	params["sign"] = sign
	var xmlBody strings.Builder
	xmlBody.WriteString("<xml>")
	for _, k := range sortedKeys(params) {
		if params[k] != "" {
			xmlBody.WriteString("<" + k + "><![CDATA[" + params[k] + "]]></" + k + ">")
		}
	}
	xmlBody.WriteString("</xml>")
	req, err := http.NewRequestWithContext(r.Context(), "POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", strings.NewReader(xmlBody.String()))
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "text/xml")
	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req) // DefaultClient 无超时,微信接口挂起会一直占着调用方的事务连接
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	result, err := wxXMLToMap(string(raw))
	if err != nil || result["return_code"] != "SUCCESS" || result["result_code"] != "SUCCESS" {
		return nil
	}
	prepayID := result["prepay_id"]
	nonce := uuid32()
	ts := fmt.Sprintf("%d", time.Now().Unix())
	pkg := map[string]string{"appId": appID, "timeStamp": ts, "nonceStr": nonce, "package": "prepay_id=" + prepayID, "signType": "MD5"}
	paySign := wxLinkSign(pkg, cfg.MiniAppPayKey)
	return map[string]any{
		"appId": appID, "timeStamp": ts, "nonceStr": nonce, "package": "prepay_id=" + prepayID, "signType": "MD5", "paySign": paySign,
		"partnerid": cfg.MchID, "prepayid": orderNb, "resprepayid": prepayID, "sign1": sign, "sign2": result["sign"],
	}
}

// wxLinkSign PayKit.createLinkString + MD5 大写：空值与 sign 不参与，键升序。
func wxLinkSign(params map[string]string, key string) string {
	keys := sortedKeys(params)
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		if k != "sign" && params[k] != "" {
			pairs = append(pairs, k+"="+params[k])
		}
	}
	sum := md5.Sum([]byte(strings.Join(pairs, "&") + "&key=" + key))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
