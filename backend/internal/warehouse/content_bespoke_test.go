package warehouse

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// 路由从模板换到定制 handler:Skip 漏一个 Register 就会因重复注册 panic。
func TestBespokeWriteRoutesRegistered(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, &Service{})
	for _, path := range []string{"/busBusinessCard/add", "/busContacts/add", "/busEnterpriseIntroduction/add",
		"/busScoreExchangeLog/add", "/busVisitorLog/add", "/busVerifyInfo/add", "/busVerifyInfo/updateById"} {
		if _, pattern := mux.Handler(httptest.NewRequest("POST", "/saas-warehouse"+path, nil)); pattern == "" {
			t.Errorf("unregistered: POST %s", path)
		}
	}
}

// filelist 是 @Many 载体,模板校验不覆盖:绑定规则对照 Jackson。
func TestVerifyFilelistFollowsJacksonBinding(t *testing.T) {
	for body, want := range map[string]bool{
		`{}`: true, `{"filelist":null}`: true, `{"filelist":[{"id":5}]}`: true,
		`{"filelist":[null]}`: true, // null 元素收得下,服务里才空指针
		`{"filelist":"x"}`:    false, `{"filelist":[1]}`: false, `{"filelist":[{"id":"abc"}]}`: false,
	} {
		if _, ok := filelistOf(bodyMap(httptest.NewRequest("POST", "/", strings.NewReader(body)))); ok != want {
			t.Errorf("%s 绑定结果应为 %v", body, want)
		}
	}
}

// bespokeFixture 七个定制写端点的隔离库夹具:登录用户 41 / 42、没有积分行的 43、当月额度另设的 44;
// 短信打到本地假端点。
type bespokeFixture struct {
	d    *sql.DB
	kp   *auth.KeyPair
	mux  http.Handler // 网关中间件 + 业务路由(外部请求)
	raw  http.Handler // 业务路由本身(commonrest.callJSON 进程内调用直接分发到这里)
	sms  []map[string]any
	user map[string]string // 账号名 → 用户 id
}

func newBespokeFixture(t *testing.T) *bespokeFixture {
	t.Helper()
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机隔离写入验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	t.Cleanup(func() { admin.Close() })
	schema := fmt.Sprintf("audit_bespoke_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { admin.Exec("DROP DATABASE " + schema) })
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	t.Cleanup(func() { d.Close() })
	for _, table := range []string{"chain_auth_user", "chain_auth_secret", "chain_common_file", "bus_business_card", "bus_contacts",
		"bus_enterprise_introduction", "bus_visitor_log", "bus_verify_info", "bus_product", "bus_score_exchange_log", "bus_score_pro",
		"bus_user_score", "bus_score_rule", "bus_score_log", "bus_msg", "bus_user_msg", "bus_config"} {
		if _, err := d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{
		"INSERT INTO chain_auth_user(id,account,phone,channel,handle_result,level,account_type,nickname,wx_face_img) VALUES" +
			"(41,'bs-owner','00000001141','2',1,1,1,'访客甲','face-41.png'),(42,'bs-other','00000001142','2',1,1,1,'路人乙','face-42.png')," +
			"(43,'bs-noscore','00000001143','2',1,1,1,'无分',''),(44,'bs-capped','00000001144','2',1,1,1,'封顶','')",
		"INSERT INTO chain_auth_secret(id,authorization_id,secret) VALUES(1,'fixture-partner','partner-secret')",
		"INSERT INTO bus_business_card(id,user_id,name,is_default) VALUES(1,41,'old-41',1),(2,42,'old-42',1)",
		"INSERT INTO chain_common_file(id,file_name,file_path,effective) VALUES(1,'a.jpg','/keep/a.jpg',0),(2,'','',0),(3,'','',0),(4,'','',0),(5,'','',0),(6,'','',0),(7,'','',0),(8,'','',0),(9,'','',0)",
		"INSERT INTO bus_product(id,is_real_check) VALUES(1,0),(2,0)",
		"INSERT INTO bus_user_score(id,user_id,hjj_num,dj_hjj_num) VALUES(1,41,100,7),(2,42,100,0),(3,44,100,0)",
		"INSERT INTO bus_score_pro(id,pro_name,pro_score) VALUES(1,'茶杯',60),(2,'等额',100),(3,'豪车',500)",
		"INSERT INTO bus_score_rule(id,mo_get_sum_score) VALUES(1,100)",
		"INSERT INTO bus_config(id,admin_phone,msg_type) VALUES(1,'00000001129',1)",
	} {
		if _, err := d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	fx := &bespokeFixture{d: d, user: map[string]string{"owner": "41", "other": "42", "noscore": "43", "capped": "44"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		fx.sms = append(fx.sms, body)
		io.WriteString(w, `{"Response":{"RequestId":"fixture"}}`)
	}))
	t.Cleanup(server.Close)
	t.Setenv("TENCENT_SMS_ENDPOINT", server.URL)
	t.Setenv("TENCENT_SECRET_ID", "fixture")
	t.Setenv("TENCENT_SECRET_KEY", "fixture")
	t.Setenv("PANDA_WAREHOUSE_SEND_MSG_NOTICE", "true")
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	fx.kp = &auth.KeyPair{Private: key, Public: &key.PublicKey}
	mux := http.NewServeMux()
	Register(mux, NewService(&db.DB{Chain: d, Zjjg: d}))
	fx.raw = mux
	fx.mux = auth.Middleware(fx.kp, d)(platform.ValidateJavaBodies(mux))
	return fx
}

// call 以指定账号(user 为空 = 不带 token)调 warehouse 路由,返回响应体里的 code
func (fx *bespokeFixture) call(t *testing.T, user, path, body string, headers map[string]string) (int, string) {
	t.Helper()
	r := httptest.NewRequest("POST", path, strings.NewReader(body))
	if user != "" {
		token, err := fx.kp.Generate(auth.JWTClaims{Account: "bs-" + user, UserID: fx.user[user]})
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Authorization", "Bearer "+token)
	}
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	fx.mux.ServeHTTP(w, r)
	var code int
	if _, err := fmt.Sscanf(strings.SplitN(w.Body.String(), `"code":`, 2)[1], "%d", &code); err != nil {
		t.Fatalf("响应里没有 code: %s", w.Body.String())
	}
	return code, w.Body.String()
}

func (fx *bespokeFixture) expect(t *testing.T, want int, user, path, body string) {
	t.Helper()
	fx.expectWith(t, want, user, path, body, nil)
}

func (fx *bespokeFixture) expectWith(t *testing.T, want int, user, path, body string, headers map[string]string) {
	t.Helper()
	if code, resp := fx.call(t, user, path, body, headers); code != want {
		t.Fatalf("%s %s 想要 code=%d 实得 %d: %s", user, path, want, code, resp)
	}
}

// row 取一行各列的文本值(NULL → "<nil>");expectRow 按逗号拼接比对
func (fx *bespokeFixture) row(t *testing.T, q string, args ...any) []string {
	t.Helper()
	rows, err := fx.d.Query(q, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	if !rows.Next() {
		t.Fatalf("无行: %s", q)
	}
	vals := make([]sql.NullString, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		t.Fatal(err)
	}
	out := make([]string, len(cols))
	for i, v := range vals {
		out[i] = "<nil>"
		if v.Valid {
			out[i] = v.String
		}
	}
	return out
}

func (fx *bespokeFixture) expectRow(t *testing.T, want, q string, args ...any) {
	t.Helper()
	if got := strings.Join(fx.row(t, q, args...), ","); got != want {
		t.Fatalf("%s 得 %q,想要 %q", q, got, want)
	}
}

// 名片与联系人:身份取登录人;新名片成为默认并清掉本人旧默认;联系人查重。
func TestCardAndContactsWritesUseLoginUser(t *testing.T) {
	fx := newBespokeFixture(t)
	card, contacts := "/saas-warehouse/busBusinessCard/add", "/saas-warehouse/busContacts/add"
	fx.expect(t, 200, "owner", card, `{"userId":42,"name":"new-41","isDefault":false}`)
	fx.expectRow(t, "41,1", "SELECT user_id,is_default FROM bus_business_card WHERE id=3")
	fx.expectRow(t, "2", "SELECT concat(is_default) FROM bus_business_card WHERE id=1") // 本人旧默认名片
	fx.expectRow(t, "1", "SELECT concat(is_default) FROM bus_business_card WHERE id=2") // 别人的名片不动
	// 同一用户对同一张名片只落一条联系人,重复提交照常回成功;card_id 为空按 IS NULL 查重
	fx.expect(t, 200, "owner", contacts, `{"userId":42,"cardId":7,"companyId":3}`)
	fx.expect(t, 200, "owner", contacts, `{"cardId":7}`)
	fx.expectRow(t, "1,41", "SELECT count(*),max(user_id) FROM bus_contacts WHERE card_id=7")
	fx.expect(t, 200, "owner", contacts, `{"cardId":null}`)
	fx.expect(t, 200, "owner", contacts, `{"cardId":null}`)
	fx.expectRow(t, "1", "SELECT count(*) FROM bus_contacts WHERE card_id IS NULL")
	// 绑定失败(Integer 字段给了非数字):999 且不写
	fx.expect(t, 999, "owner", card, `{"headImgName":"abc"}`)
	fx.expectRow(t, "3", "SELECT count(*) FROM bus_business_card")
}

// 企业介绍:logo / 图片文件置为生效;文件号不是整数回 999,该批不生效,企业介绍已落库(同 Java)。
func TestEnterpriseIntroductionActivatesFiles(t *testing.T) {
	fx := newBespokeFixture(t)
	add := "/saas-warehouse/busEnterpriseIntroduction/add"
	fx.expect(t, 200, "owner", add, `{"userId":42,"cname":"甲公司","logoId":"5","imgIds":"6,7"}`)
	fx.expectRow(t, "41", "SELECT user_id FROM bus_enterprise_introduction")
	fx.expectRow(t, "1,1,1,0", "SELECT group_concat(effective ORDER BY id) FROM chain_common_file WHERE id IN (5,6,7,8)")
	// logo 那一批先走完(9 生效),imgIds 里有非整数 → 999,8 不生效
	fx.expect(t, 999, "owner", add, `{"cname":"乙","logoId":"9","imgIds":"8,bad"}`)
	fx.expectRow(t, "0,1", "SELECT group_concat(effective ORDER BY id) FROM chain_common_file WHERE id IN (8,9)")
	fx.expectRow(t, "2", "SELECT count(*) FROM bus_enterprise_introduction")
}

// 积分兑换:余额不足不兑换;足够时按商品价扣登录人的分,并写流水、站内消息、管理员短信。
func TestScoreExchangeDeductsAndNotifies(t *testing.T) {
	fx := newBespokeFixture(t)
	exchange := "/saas-warehouse/busScoreExchangeLog/add"
	score := "SELECT concat(hjj_num) FROM bus_user_score WHERE user_id="
	// 余额不足:请求体把 userId 写成别人、consumeScore 写成 1 分也换不到——不落记录、不扣分、不发消息和短信
	fx.expect(t, 602, "owner", exchange, `{"userId":42,"goodsId":3,"consumeScore":1,"proName":"x"}`)
	// 余额恰好等于商品价(Java 的判断是严格大于)同样不兑换
	fx.expect(t, 602, "owner", exchange, `{"goodsId":2,"consumeScore":1}`)
	for _, q := range []string{"SELECT count(*) FROM bus_score_exchange_log", "SELECT count(*) FROM bus_msg", "SELECT count(*) FROM bus_score_log"} {
		fx.expectRow(t, "0", q)
	}
	fx.expectRow(t, "100", score+"41")
	fx.expectRow(t, "100", score+"42")
	if len(fx.sms) != 0 {
		t.Fatalf("兑换被拒却发了短信: %v", fx.sms)
	}
	// 足够:商品名、扣分数取积分商品表,扣的一定是登录人
	fx.expect(t, 200, "owner", exchange, `{"userId":42,"goodsId":1,"consumeScore":1,"proName":"伪造","consigneeName":"张三"}`)
	fx.expectRow(t, "41,茶杯,60", "SELECT user_id,pro_name,consume_score FROM bus_score_exchange_log")
	fx.expectRow(t, "40", score+"41")
	fx.expectRow(t, "100", score+"42")
	// 积分流水照抄 Java 类型 7 分支错位的列:变动后可用积分写成 2、规则类型写成变动前余额
	fx.expectRow(t, "41,兑换商品-茶杯,0,100,2,7,7,100,7",
		"SELECT user_id,content,score,orgin_available_score,final_available_score,orgin_freeze_score,final_freeze_score,rule_type,score_get_type FROM bus_score_log")
	fx.expectRow(t, "5,您已成功兑换茶杯，3到7日内发货，请留意物流信息", "SELECT msg_type,content FROM bus_msg")
	fx.expectRow(t, "41", "SELECT user_id FROM bus_user_msg WHERE msg_id=(SELECT id FROM bus_msg)")
	if len(fx.sms) != 1 || fmt.Sprint(fx.sms[0]["TemplateID"]) != "2246413" ||
		fmt.Sprint(fx.sms[0]["TemplateParamSet"]) != "[茶杯]" || fmt.Sprint(fx.sms[0]["PhoneNumberSet"]) != "[8600000001129]" {
		t.Fatalf("管理员短信不对: %v", fx.sms)
	}
	// 再换余额不够;没有积分行 / 商品不存在 → 999
	fx.expect(t, 602, "owner", exchange, `{"goodsId":1}`)
	fx.expect(t, 999, "noscore", exchange, `{"goodsId":1}`)
	fx.expect(t, 999, "owner", exchange, `{"goodsId":999}`)
	fx.expectRow(t, "1", "SELECT count(*) FROM bus_score_exchange_log")
	// 当月积分额度为 0:Java 到这里不扣分却已落兑换记录;这里照常扣分
	if _, err := fx.d.Exec("UPDATE bus_score_rule SET mo_get_sum_score=0"); err != nil {
		t.Fatal(err)
	}
	fx.expect(t, 200, "capped", exchange, `{"goodsId":1}`)
	fx.expectRow(t, "40", score+"44")
	fx.expectRow(t, "2", "SELECT count(*) FROM bus_score_exchange_log")
}

// 访客记录:首次用登录人自己的头像昵称新增;同一目标再来不会重复新增;带 id 时按请求体更新且归属不能改。
func TestVisitorLogRecordsLoginVisitorOnce(t *testing.T) {
	fx := newBespokeFixture(t)
	add := "/saas-warehouse/busVisitorLog/add"
	fx.expect(t, 200, "owner", add, `{"visitorId":42,"targetUserId":77,"companyId":5}`)
	fx.expectRow(t, "41,77,face-41.png,访客甲,<nil>", "SELECT visitor_id,target_user_id,target_head_img,target_nickname,company_id FROM bus_visitor_log")
	// 已有记录且请求体没有 id:Java 的 updateById(请求体实体) 等于空更新
	fx.expect(t, 200, "owner", add, `{"visitorId":42,"targetUserId":77}`)
	fx.expectRow(t, "1", "SELECT count(*) FROM bus_visitor_log")
	fx.expect(t, 200, "owner", add, `{"id":1,"visitorId":42,"targetUserId":77,"companyId":9}`)
	fx.expectRow(t, "41,9", "SELECT visitor_id,company_id FROM bus_visitor_log WHERE id=1")
	// 另一个登录人访问同一目标:各记各的
	fx.expect(t, 200, "other", add, `{"targetUserId":77}`)
	fx.expectRow(t, "42,face-42.png", "SELECT visitor_id,target_head_img FROM bus_visitor_log WHERE id=2")
}

// 核实信息:关联号由服务端生成、附件挂到登录人名下、磁盘路径列不被改写、煤种 is_real_check 置真。
func TestVerifyInfoLinksFilesAndMarksProduct(t *testing.T) {
	fx := newBespokeFixture(t)
	add, upd := "/saas-warehouse/busVerifyInfo/add", "/saas-warehouse/busVerifyInfo/updateById"
	file := func(id int) string {
		return strings.Join(fx.row(t, "SELECT owner_type,owner_id,effective,association_id,file_path,file_name FROM chain_common_file WHERE id="+strconv.Itoa(id)), ",")
	}
	fx.expect(t, 200, "owner", add, `{"proId":1,"evaluate":"e1","associationId":"forged","filelist":[{"id":1,"filePath":"/etc/passwd","fileName":"../../x.jpg","fileRealName":"real.jpg"},{"id":2}]}`)
	assoc := fx.row(t, "SELECT association_id FROM bus_verify_info WHERE id=1")[0]
	if len(assoc) != 32 || assoc == "forged" {
		t.Fatalf("关联号应由服务端生成: %q", assoc)
	}
	if got := file(1); got != "2,41,1,"+assoc+",/keep/a.jpg,a.jpg" { // 路径穿越的 fileName / filePath 都不生效
		t.Fatalf("附件1: %q", got)
	}
	if got := file(2); got != "2,41,1,"+assoc+",," {
		t.Fatalf("附件2: %q", got)
	}
	fx.expectRow(t, "1", "SELECT concat(is_real_check) FROM bus_product WHERE id=1")
	fx.expectRow(t, "0", "SELECT concat(is_real_check) FROM bus_product WHERE id=2")
	// 修改:只改请求体字段;没带关联号就不动附件的关联号
	fx.expect(t, 200, "owner", upd, `{"id":1,"evaluate":"e2","filelist":[{"id":3}]}`)
	fx.expectRow(t, "e2", "SELECT evaluate FROM bus_verify_info WHERE id=1")
	if got := file(3); got != "2,41,1,<nil>,," {
		t.Fatalf("附件3: %q", got)
	}
	// 煤种不存在:999,但附件已先挂上(同 Java);filelist 不是数组 → 绑定失败不写主记录;null 元素 → 主记录先写入再 999
	fx.expect(t, 999, "owner", add, `{"proId":999,"filelist":[{"id":4}]}`)
	fx.expectRow(t, "1", "SELECT concat(effective) FROM chain_common_file WHERE id=4")
	fx.expect(t, 999, "owner", add, `{"proId":1,"filelist":"x"}`)
	fx.expectRow(t, "2", "SELECT count(*) FROM bus_verify_info")
	fx.expect(t, 999, "owner", add, `{"proId":1,"filelist":[null]}`) // 附件预检先于主记录,不再留下没有附件的核实记录(Java 是主记录先写入再空指针)
	fx.expectRow(t, "2", "SELECT count(*) FROM bus_verify_info")
	// 别人的附件不能挂到自己名下:602,不留下核实记录,附件原样;附件本人提交照常
	if _, err := fx.d.Exec("UPDATE chain_common_file SET owner_id=42 WHERE id=8"); err != nil {
		t.Fatal(err)
	}
	fx.expect(t, 602, "owner", add, `{"proId":1,"filelist":[{"id":8}]}`)
	fx.expect(t, 602, "owner", upd, `{"id":1,"filelist":[{"id":8}]}`)
	fx.expectRow(t, "2", "SELECT count(*) FROM bus_verify_info")
	if got := file(8); got != "<nil>,42,0,<nil>,," {
		t.Fatalf("他人附件被改动: %q", got)
	}
	fx.expect(t, 200, "other", add, `{"proId":1,"filelist":[{"id":8}]}`)
	fx.expectRow(t, "3", "SELECT count(*) FROM bus_verify_info")
	if got := file(8); got != "2,42,1,"+strings.Join(fx.row(t, "SELECT association_id FROM bus_verify_info WHERE id=3"), ",")+",," {
		t.Fatalf("本人附件未挂上: %q", got)
	}
}

// 仅带签名的匿名请求(没有登录身份)回 999;登录链在进程内调用名片 add(此时还没有 token)
// 走 Java 的 jwt==null 分支,按 body.userId 建默认名片。
func TestBespokeWritesRejectAnonymousButAllowLoginChainCall(t *testing.T) {
	fx := newBespokeFixture(t)
	mac := hmac.New(sha256.New, []byte("partner-secret"))
	mac.Write([]byte("fixture-partner" + "hyqc" + "partner-secret"))
	signed := map[string]string{"AuthorizationId": "fixture-partner", "Signature": hex.EncodeToString(mac.Sum(nil))}
	for _, path := range []string{"/saas-warehouse/busBusinessCard/add", "/saas-warehouse/busContacts/add", "/saas-warehouse/busScoreExchangeLog/add"} {
		fx.expectWith(t, 999, "", path, `{"userId":42,"cardId":7,"goodsId":1}`, signed)
	}
	fx.expectRow(t, "2", "SELECT count(*) FROM bus_business_card")
	fx.expectRow(t, "0", "SELECT count(*) FROM bus_contacts")
	fx.expectRow(t, "0", "SELECT count(*) FROM bus_score_exchange_log")

	// 取一个"验签通过的匿名请求"的上下文,再按 callJSON 的方式带上进程内标记分发到业务路由
	var ctx context.Context
	probe := auth.Middleware(fx.kp, fx.d)(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { ctx = r.Context() }))
	pr := httptest.NewRequest("POST", "/saas-warehouse/busBusinessCard/add", nil)
	pr.Header.Set("AuthorizationId", signed["AuthorizationId"])
	pr.Header.Set("Signature", signed["Signature"])
	probe.ServeHTTP(httptest.NewRecorder(), pr)
	inner := func(body string) string {
		r := httptest.NewRequest("POST", "/saas-warehouse/busBusinessCard/add", strings.NewReader(body)).WithContext(auth.WithInternalCall(ctx))
		w := httptest.NewRecorder()
		fx.raw.ServeHTTP(w, r)
		return w.Body.String()
	}
	if resp := inner(`{"userId":42,"phone":"00000001142"}`); !strings.Contains(resp, `"code":200`) {
		t.Fatal(resp)
	}
	fx.expectRow(t, "42,1", "SELECT user_id,is_default FROM bus_business_card WHERE id=3")
	fx.expectRow(t, "2", "SELECT concat(is_default) FROM bus_business_card WHERE id=2")
	// 进程内调用但请求体没有 userId:Java 的 getUserId().toString() 空指针
	if resp := inner(`{"phone":"00000001142"}`); !strings.Contains(resp, `"code":999`) {
		t.Fatal(resp)
	}
	fx.expectRow(t, "3", "SELECT count(*) FROM bus_business_card")
}
