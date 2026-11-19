package warehouse

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html"
	"image"
	_ "image/jpeg"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func (s *Service) bpWXToken(r *http.Request) any {
	token, err := platform.WechatAccessToken(r.Context(), s.pick(r), s.HTTPClient, "cloudAccount")
	if err != nil {
		return ""
	} // 原 Feign getAccessToken 捕获异常后成功返回空串。
	return token
}

func (s *Service) bpTemplate(r *http.Request, body map[string]any) error {
	token := s.bpWXToken(r)
	if token == nil {
		return fmt.Errorf("Java accessToken 为 null")
	}
	_, _ = platform.WechatJSON(r.Context(), s.HTTPClient, "POST", "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token="+url.QueryEscape(platform.JnullAny(token)), body)
	return nil // 原 pushTemplateMsg 吞掉发送异常。
}

func (s *Service) bpPushOfficialActual(w http.ResponseWriter, r *http.Request) {
	if err := s.bpTemplate(r, map[string]any{}); err != nil {
		platform.Unknown(w)
		return
	}
	platform.Success(w, nil)
}

func (s *Service) bpGetCareActual(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	if m["openid"] == nil {
		platform.Unknown(w)
		return
	}
	token := s.bpWXToken(r)
	if token == nil {
		platform.Unknown(w)
		return
	}
	res, err := platform.WechatJSON(r.Context(), s.HTTPClient, "GET", "https://api.weixin.qq.com/cgi-bin/user/info?"+url.Values{"access_token": {platform.JnullAny(token)}, "openid": {platform.JnullAny(m["openid"])}, "lang": {"zh_CN"}}.Encode(), nil)
	if err != nil || res["subscribe"] == nil {
		platform.FailData(w, 999, "失败", "判断是否关注公众号异常：")
		return
	}
	text := "未关注"
	if num(res["subscribe"]) == 1 {
		text = "已关注"
	}
	platform.Success(w, text)
}

func warehouseNoticeEnabled(name string) bool {
	v := os.Getenv("PANDA_WAREHOUSE_" + name)
	return v == "1" || strings.EqualFold(v, "true")
}

func (s *Service) bpBenefits(r *http.Request, p map[string]any) (map[string]any, error) {
	u, err := platform.GetByID(s.pick(r), "chain_auth_user", p["createUserId"])
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("Java createUser 为 null")
	}
	rows, err := queryCam(s.pick(r), "SELECT * FROM bus_member_benefits WHERE member_type=?", u["member_type"])
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, fmt.Errorf("Java 会员权益不存在或不唯一")
	}
	return rows[0], nil
}

func (s *Service) bpNotifyReview(r *http.Request, p map[string]any, orgStatus string) error {
	if orgStatus == "0" {
		u, err := platform.GetByID(s.pick(r), "chain_auth_user", p["createUserId"])
		if err != nil {
			return err
		}
		if u == nil {
			return fmt.Errorf("Java 审核用户为 null")
		}
		if warehouseNoticeEnabled("SEND_MSG_NOTICE") {
			check := ""
			if num(p["status"]) == 1 {
				check = "通过"
			} else if num(p["status"]) == 2 {
				check = "拒绝"
			}
			platform.SendLegacySMS(r.Context(), platform.JnullAny(u["phone"]), []string{platform.JnullAny(p["enterpriseName"]) + " " + platform.JnullAny(p["thName"]) + "，" + platform.JnullAny(p["exFactoryPrice"]) + "出厂价", check}, "1767390", "山西海云启辰科技")
		}
	} else if orgStatus == "1" && num(p["status"]) == 3 && warehouseNoticeEnabled("SEND_MSG_NOTICE") {
		rows, err := queryCam(s.pick(r), "SELECT * FROM bus_config WHERE msg_type=1")
		if err != nil {
			return err
		}
		for _, row := range rows {
			platform.SendLegacySMS(r.Context(), str(row, "adminPhone"), []string{platform.JnullAny(p["enterpriseName"]) + " " + platform.JnullAny(p["thName"]) + " " + platform.JnullAny(p["exFactoryPrice"]) + "元（出厂价）"}, "2119567", "山西海云启辰科技")
		}
	}
	return nil
}

func (s *Service) bpNotifyPending(r *http.Request, p map[string]any, oldPrice any, userID string) error {
	diff, ok := decimalSubtractExact(p["exFactoryPrice"], oldPrice)
	if !ok {
		return fmt.Errorf("Java 价格为 null")
	}
	priceContent := ""
	if num(diff) > 0 {
		priceContent = "，" + platform.JnullAny(oldPrice) + "元更新到" + platform.JnullAny(p["exFactoryPrice"]) + "元，上涨" + diff + "元"
	} else if num(diff) < 0 {
		priceContent = "，" + platform.JnullAny(oldPrice) + "元更新到" + platform.JnullAny(p["exFactoryPrice"]) + "元，下降" + strings.TrimPrefix(diff, "-") + "元"
	}
	if !warehouseNoticeEnabled("SEND_MSG_NOTICE") {
		return nil
	}
	u, err := platform.GetByID(s.pick(r), "chain_auth_user", userID)
	if err != nil {
		return err
	}
	if u == nil {
		return fmt.Errorf("Java 操作用户为 null")
	}
	rows, err := queryCam(s.pick(r), "SELECT * FROM bus_config WHERE msg_type=1")
	if err != nil {
		return err
	}
	for _, row := range rows {
		platform.SendLegacySMS(r.Context(), platform.JnullAny(row["adminPhone"]), []string{platform.JnullAny(u["phone"]), platform.JnullAny(p["enterpriseName"]) + " " + platform.JnullAny(p["thName"]) + priceContent}, "2068879", "")
	}
	return nil
}

func (s *Service) bpNotifyPrice(r *http.Request, p map[string]any, oldPrice, newPrice any) error {
	if warehouseNoticeEnabled("SEND_MSG_NOTICE") {
		rows, err := queryCam(s.pick(r), "SELECT u.id,u.phone FROM bus_pro_collection c JOIN chain_auth_user u ON u.id=c.user_id WHERE c.pro_id=?", p["id"])
		if err != nil {
			return err
		}
		now := time.Now()
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, time.Local)
		phones := []string{}
		for _, row := range rows {
			if _, fresh := s.redCachePut("smsPhone"+str(row, "id"), time.Until(end), false); fresh {
				phones = append(phones, platform.JnullAny(row["phone"]))
			}
		}
		params := []string{platform.JnullAny(p["enterpriseName"]) + " " + platform.JnullAny(p["thName"]) + "，" + platform.JnullAny(oldPrice) + "更新到" + platform.JnullAny(p["exFactoryPrice"]) + "，"}
		for len(phones) > 0 {
			n := len(phones)
			if n > 200 {
				n = 200
			}
			platform.SendLegacySMSMany(r.Context(), phones[:n], params, "2165435", "山西海云启辰科技")
			phones = phones[n:]
		}
	}
	if err := s.bpAppNotice(r, p, 2, oldPrice, newPrice); err != nil {
		return err
	}
	return s.bpOfficialNotice(r, p, oldPrice, newPrice, false)
}

func (s *Service) bpNotifyAdded(r *http.Request, p map[string]any) error {
	if err := s.bpAppNotice(r, p, 1, nil, nil); err != nil {
		return err
	}
	return s.bpOfficialNotice(r, p, nil, nil, true)
}

func (s *Service) bpAppNotice(r *http.Request, p map[string]any, kind int, oldPrice, newPrice any) error {
	benefits, err := s.bpBenefits(r, p)
	if err != nil {
		return err
	}
	if !warehouseNoticeEnabled("SEND_OFFICE_NOTICE") {
		return nil
	}
	content := platform.JnullAny(p["warehouseNb"]) + " " + platform.JnullAny(p["exFactoryPrice"]) + "元（出厂价）"
	title := "今日新增煤种"
	if num(benefits["isShowEtcName"]) == 1 {
		content = platform.JnullAny(p["enterpriseName"]) + " " + content
	}
	if kind == 2 {
		diff, ok := decimalSubtractExact(newPrice, oldPrice)
		if !ok {
			return nil
		} // 原 try 内异常吞掉。
		direction := "上升"
		if num(diff) < 0 {
			direction = "下降"
		}
		content = platform.JnullAny(p["warehouseNb"]) + "今日" + direction + strings.TrimPrefix(diff, "-") + "元"
		title = "今日价格变动"
	}
	body := map[string]any{"appId": os.Getenv("PANDA_WAREHOUSE_UNIPUSH_APP_ID"), "title": title, "content": content, "payload": map[string]any{"url": "/pages_bus/specifics/specifics?speid=null&id=" + platform.JnullAny(p["id"]) + "&pname=" + platform.JnullAny(p["warehouseNb"]), "pathType": 0}}
	_, _ = platform.WechatJSON(r.Context(), s.HTTPClient, "POST", os.Getenv("PANDA_WAREHOUSE_UNIPUSH_ADR"), body)
	return nil
}

func (s *Service) bpOfficialNotice(r *http.Request, p map[string]any, oldPrice, newPrice any, added bool) error {
	if !warehouseNoticeEnabled("SEND_OFFICE_NOTICE") {
		return nil
	}
	key := "updProNotice" + str(p, "id") + platform.JnullAny(oldPrice) + platform.JnullAny(newPrice)
	if added {
		key = "addProNotice" + str(p, "id")
	}
	if _, fresh := s.redCachePut(key, time.Hour, false); !fresh {
		return nil
	}
	benefits, err := s.bpBenefits(r, p)
	if err != nil {
		return err
	}
	name := p["warehouseNb"]
	if num(benefits["isShowEtcName"]) == 1 {
		name = p["enterpriseName"]
	}
	remark := ""
	if added {
		price := ""
		if int(num(p["exFactoryPrice"])) != 0 {
			price = "价格" + strconv.Itoa(int(num(p["exFactoryPrice"]))) + "元/吨"
		} else if int(num(p["platePrice"])) != 0 {
			price = "价格" + strconv.Itoa(int(num(p["platePrice"]))) + "元/吨（车板价）"
		}
		remark = "企业新增煤种  " + price
	} else {
		diff, ok := decimalSubtractExact(newPrice, oldPrice)
		if !ok {
			return fmt.Errorf("Java 通知价格为 null")
		}
		direction := "上涨"
		if num(diff) <= 0 {
			direction = "下跌"
		}
		remark = "价格" + direction + strconv.Itoa(int(num(strings.TrimPrefix(diff, "-")))) + "元"
	}
	users, err := queryCam(s.pick(r), "SELECT official_openid FROM chain_auth_user WHERE official_openid IS NOT NULL")
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, u := range users {
		openid := str(u, "officialOpenid")
		if seen[openid] {
			continue
		}
		seen[openid] = true
		body := map[string]any{"touser": openid, "template_id": "-HEsls2xJP20Xt5WGwQ6oyWwcA3t03KunyCb8Zg2w6Q", "url": "#小程序://熊猫煤仓/mj5LeVKo8hpukgm", "miniprogram": map[string]any{"appid": os.Getenv("PANDA_WAREHOUSE_PUSH_APPID"), "pagepath": "pages_bus/specifics/specifics?id=" + str(p, "id")}, "data": map[string]any{"thing11": map[string]any{"value": name}, "thing4": map[string]any{"value": "硫≤" + platform.JnullAny(p["maxStd"]) + "灰≤" + platform.JnullAny(p["maxAshContent"]) + "..."}, "thing9": map[string]any{"value": remark}}}
		if err = s.bpTemplate(r, body); err != nil {
			return err
		}
	}
	return nil
}

// 分享图复用已安装的 LibreOffice/Poppler，保留 Java 缺字段或模板时返回 null 的行为。
func (s *Service) bpGetDetailImgActual(w http.ResponseWriter, r *http.Request) {
	m := bodyMap(r)
	name, err := s.bpSharePicture(r, m["id"])
	if err != nil {
		platform.Success(w, nil)
		return
	}
	platform.Success(w, name)
}

func (s *Service) bpSharePicture(r *http.Request, id any) (string, error) {
	raw, err := platform.GetByID(s.pick(r), "bus_product", id)
	if err != nil {
		return "", err
	}
	if raw == nil {
		return "", fmt.Errorf("Java 煤种为 null")
	}
	p := platform.Camelize(raw)
	diff, ok := decimalSubtractExact(p["exFactoryPrice"], p["preExFactoryPrice"])
	if !ok {
		return "", fmt.Errorf("Java 分享图价格为 null")
	}
	benefits, err := s.bpBenefits(r, p)
	if err != nil {
		return "", err
	}
	if benefits["isShowEtcName"] == nil || benefits["isShowEtcPhone"] == nil {
		return "", fmt.Errorf("Java 分享图权益为 null")
	}
	root := os.Getenv("PANDA_UPLOAD_DIR")
	if root == "" {
		root = "data/uploads"
	}
	background, err := os.ReadFile(filepath.Join(root, "frame73.jpg"))
	if err != nil {
		return "", err
	}
	im, _, err := image.DecodeConfig(bytes.NewReader(background))
	if err != nil {
		return "", err
	}
	var svg strings.Builder
	fmt.Fprintf(&svg, `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="%dpx" height="%dpx" viewBox="0 0 %d %d"><image width="100%%" height="100%%" xlink:href="data:image/jpeg;base64,%s"/>`, im.Width, im.Height, im.Width, im.Height, base64.StdEncoding.EncodeToString(background))
	text := func(value string, size, x, y int, color string) {
		fmt.Fprintf(&svg, `<text x="%d" y="%d" font-family="Microsoft YaHei,Noto Sans CJK SC,sans-serif" font-size="%d" fill="%s">%s</text>`, x, y+size, size, color, html.EscapeString(value))
	}
	text(platform.JnullAny(p["warehouseNb"]), 40, 210, 130, "black")
	text(time.Now().Format("2006/01/02"), 30, 440, 50, "white")
	if num(benefits["isShowEtcName"]) == 1 {
		text(platform.JnullAny(p["enterpriseName"]), 25, 90, 230, "black")
	}
	contact, phone := "***", "***********"
	if num(benefits["isShowEtcPhone"]) == 1 {
		contact, phone = platform.JnullAny(p["coalRelName"]), platform.JnullAny(p["coalRelPhone"])
	}
	text(contact, 25, 280, 751, "black")
	text(phone, 20, 280, 793, "black")
	text("出厂现汇价："+platform.JnullAny(p["exFactoryPrice"])+"元", 25, 90, 275, "black")
	if num(diff) != 0 {
		arrow := "▲"
		if num(diff) < 0 {
			arrow = "▼"
		}
		text("较上次："+arrow+strings.TrimPrefix(diff, "-")+"元", 25, 380, 275, "black")
	}
	text(platform.JnullAny(p["businessName"]), 25, 340, 1016, "white")
	text(platform.JnullAny(p["businessPhone"]), 20, 280, 1060, "white")
	index := 0
	for _, field := range [][2]string{{"maxAshContent", "灰分(Ad)%≤"}, {"maxWaterContent", "水分(Mt)%≤"}, {"maxStd", "硫(St,d)%≤"}, {"maxVolatiles", "挥发分(Vdaf)%≤"}, {"minBondIndex", "粘结指数GR.I≥"}, {"minY", "Ymm≥"}, {"minCsr", "热态CSR≥"}, {"lithofacies", "岩相S"}} {
		if p[field[0]] == nil {
			continue
		}
		x := 90
		if index%2 == 1 {
			x = 370
		}
		text(field[1]+platform.JnullAny(p[field[0]]), 28, x, 330+(index/2)*55, "black")
		index++
	}
	svg.WriteString("</svg>")
	name := strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
	if err = renderProductSVG(r.Context(), svg.String(), filepath.Join(root, name), im.Width, im.Height); err != nil {
		return "", err
	}
	_, err = platform.InsertNoFill(s.pick(r), "chain_common_file", map[string]any{"file_real_name": name, "file_name": name, "file_path": root, "file_type": "5", "owner_type": "2"})
	return name, err
}

func renderProductSVG(ctx context.Context, svg, target string, width, height int) error {
	tmp, err := os.MkdirTemp("", "panda-share-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	src := filepath.Join(tmp, "share.svg")
	if err = os.WriteFile(src, []byte(svg), 0600); err != nil {
		return err
	}
	bin := os.Getenv("PANDA_SOFFICE_BIN")
	if bin == "" {
		bin = "soffice"
	}
	raster := os.Getenv("PANDA_PDFTOPPM_BIN")
	if raster == "" {
		raster = "pdftoppm"
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	profile := (&url.URL{Scheme: "file", Path: filepath.Join(tmp, "profile")}).String()
	if out, err := exec.CommandContext(ctx, bin, "-env:UserInstallation="+profile, "--headless", "--convert-to", "pdf:draw_pdf_Export", "--outdir", tmp, src).CombinedOutput(); err != nil {
		return fmt.Errorf("分享图转换失败: %w: %s", err, out)
	}
	base := filepath.Join(tmp, "share-image")
	if out, err := exec.CommandContext(ctx, raster, "-jpeg", "-singlefile", "-scale-to-x", strconv.Itoa(width), "-scale-to-y", strconv.Itoa(height), filepath.Join(tmp, "share.pdf"), base).CombinedOutput(); err != nil {
		return fmt.Errorf("分享图栅格化失败: %w: %s", err, out)
	}
	data, err := os.ReadFile(base + ".jpg")
	if err != nil {
		return err
	}
	return os.WriteFile(target, data, 0644)
}
