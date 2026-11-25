package warehouse

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type redFixture struct {
	d      *sql.DB
	kp     *auth.KeyPair
	mux    http.Handler
	amount []float64 // 每次打款记两项:[金额, 收款人(1=owner 2=other 3=ptadmin)]
}

const (
	receivePath  = "/saas-warehouse/busRedEnvelopeRecord/receiveAmount"
	generatePath = "/saas-warehouse/busRedEnvelopeRecord/generateAmount"
)

// 红包提现 / 生成的边界:金额与收款人以流水为准、只能提自己的、服务端发放的类型不能自己造;
// 意见 / 视频审核只允许平台管理员。打款走本地假的微信商户转账接口。
func newRedFixture(t *testing.T) *redFixture {
	t.Helper()
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机隔离写入验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	t.Cleanup(func() { admin.Close() })
	schema := fmt.Sprintf("audit_redb_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { admin.Exec("DROP DATABASE " + schema) })
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	t.Cleanup(func() { d.Close() })
	for _, table := range []string{"chain_auth_user", "chain_global_configuration", "bus_red_money_clog", "bus_red_envelope_record",
		"bus_red_envelope_batch", "bus_red_envelope_detail", "bus_red_envelope_range", "bus_opinion", "bus_videos"} {
		if _, err := d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{
		"INSERT INTO chain_auth_user(id,account,handle_result,level,account_type,phone,openid) VALUES(1,'owner',1,1,1,'00000001122','openid-owner'),(2,'other',1,1,1,'00000001123','openid-other'),(3,'ptadmin',1,2,2,'00000001124','openid-ptadmin')",
		"INSERT INTO chain_global_configuration(id,red_uplimit_days) VALUES(1,1000)",
		"INSERT INTO bus_red_money_clog(id,user_id,receive_type,create_amount,is_write_off) VALUES(1,1,7,99,0),(2,1,3,50,0),(3,1,3,50,0),(4,1,3,50,0),(5,1,3,50,0)",
		"INSERT INTO bus_red_envelope_range(receive_type,min_amount,max_amount) VALUES(3,50,50),(7,40,40),(10,60,60)",
		"INSERT INTO bus_opinion(id,user_id,phone,content,status) VALUES(1,1,'00000001122','fixture',0)",
		"INSERT INTO bus_videos(id,user_id,video_name,status) VALUES(1,1,'fixture-video',0)",
	} {
		if _, err := d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	der, _ := x509.MarshalPKCS8PrivateKey(key)
	path := filepath.Join(t.TempDir(), "key.pem")
	os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0600)
	cert := &x509.Certificate{SerialNumber: big.NewInt(123), NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	encoded, _ := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	old := platform.LoadWxConfig
	t.Cleanup(func() { platform.LoadWxConfig = old })
	platform.LoadWxConfig = func() platform.WxLocalConfig {
		return platform.WxLocalConfig{MchID: "fixture", MerchantSerial: "fixture", MerchantPrivateKeyPath: path, PlatformCert: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: encoded}))}
	}
	fx := &redFixture{d: d, kp: &auth.KeyPair{Private: key, Public: &key.PublicKey}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		fx.amount = append(fx.amount, body["total_amount"].(float64))
		openid := body["transfer_detail_list"].([]any)[0].(map[string]any)["openid"].(string)
		fx.amount = append(fx.amount, map[string]float64{"openid-owner": 1, "openid-other": 2, "openid-ptadmin": 3}[openid])
		data := `{"batch_id":"fixture"}`
		stamp, nonce := fmt.Sprint(time.Now().Unix()), "fixture"
		hash := sha256.Sum256([]byte(stamp + "\n" + nonce + "\n" + data + "\n"))
		signature, _ := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
		w.Header().Set("Wechatpay-Timestamp", stamp)
		w.Header().Set("Wechatpay-Nonce", nonce)
		w.Header().Set("Wechatpay-Serial", "7b")
		w.Header().Set("Wechatpay-Signature", base64.StdEncoding.EncodeToString(signature))
		io.WriteString(w, data)
	}))
	t.Cleanup(server.Close)
	t.Setenv("WECHAT_PAY_ENDPOINT", server.URL)
	t.Setenv("PANDA_WAREHOUSE_RED_ENABLED", "1")
	t.Setenv("PANDA_WAREHOUSE_RED_START_TIME", "00:00:00")
	t.Setenv("PANDA_WAREHOUSE_RED_END_TIME", "23:59:59")
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+receivePath, s.rrReceiveAmount)
	mux.HandleFunc("POST "+generatePath, s.rrGenerateAmount)
	mux.HandleFunc("POST /saas-warehouse/busOpinion/updateById", s.updateOpinion)
	mux.HandleFunc("POST /saas-warehouse/busVideos/updateById", s.updateVideo)
	fx.mux = auth.Middleware(fx.kp, d)(platform.ValidateJavaBodies(mux))
	return fx
}

// call 以指定账号发起请求,返回响应体里的 code。
func (fx *redFixture) call(t *testing.T, user, path, body string) (int, string) {
	t.Helper()
	token, err := fx.kp.Generate(auth.JWTClaims{Account: user, UserID: map[string]string{"owner": "1", "other": "2", "ptadmin": "3"}[user]})
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("POST", path, strings.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	fx.mux.ServeHTTP(w, r)
	var code int
	if _, err := fmt.Sscanf(strings.SplitN(w.Body.String(), `"code":`, 2)[1], "%d", &code); err != nil {
		t.Fatalf("响应里没有 code: %s", w.Body.String())
	}
	return code, w.Body.String()
}

func TestRedEnvelopeWithdrawalBoundary(t *testing.T) {
	fx := newRedFixture(t)
	// 金额、类型、收款人都以流水为准:请求体写 7 元、收款人写别人,实际按流水生成时的 99 打给流水所属人
	if code, body := fx.call(t, "owner", receivePath, `{"userId":2,"receiveType":3,"receiveAmount":7,"clogId":1}`); code != 200 ||
		len(fx.amount) != 2 || fx.amount[0] != 99 || fx.amount[1] != 1 {
		t.Fatalf("提现未按流水出款: code=%d paid=%v %s", code, fx.amount, body)
	}
	// 提别人的流水:拒绝且不打款
	fx.amount = nil
	if code, _ := fx.call(t, "other", receivePath, `{"userId":2,"receiveType":3,"receiveAmount":50,"clogId":2}`); code != 603 || len(fx.amount) != 0 {
		t.Fatalf("提现了别人的流水: code=%d paid=%v", code, fx.amount)
	}
	// 原先没有次数上限的类型:每人每类每天最多 3 次
	for i := 0; i < 3; i++ {
		if code, body := fx.call(t, "owner", receivePath, fmt.Sprintf(`{"userId":1,"receiveType":3,"receiveAmount":50,"clogId":%d}`, 2+i)); code != 200 {
			t.Fatalf("第 %d 次领取应成功: code=%d %s", i+1, code, body)
		}
	}
	if code, _ := fx.call(t, "owner", receivePath, `{"userId":1,"receiveType":3,"receiveAmount":50,"clogId":5}`); code != 603 || len(fx.amount) != 6 {
		t.Fatalf("第 4 次领取应被拒绝: code=%d paid=%v", code, fx.amount)
	}
	// 服务端在业务动作里发放的类型,客户端不能自己造流水
	if code, _ := fx.call(t, "owner", generatePath, `{"receiveType":7}`); code != 603 {
		t.Fatalf("意见奖励流水不应能自造: code=%d", code)
	}
	if code, _ := fx.call(t, "owner", generatePath, `{"receiveType":3}`); code != 200 {
		t.Fatalf("正常生成流水失败: code=%d", code)
	}
}

func TestAuditEndpointsRequirePlatformAdmin(t *testing.T) {
	fx := newRedFixture(t)
	// 普通用户"采纳"自己的意见、审核自己的视频:拒绝且状态不变
	if code, _ := fx.call(t, "owner", "/saas-warehouse/busOpinion/updateById", `{"id":1,"status":1,"userId":1}`); code != 602 {
		t.Fatalf("普通用户不应能审核意见: code=%d", code)
	}
	if code, _ := fx.call(t, "owner", "/saas-warehouse/busVideos/updateById", `{"id":1,"status":1}`); code != 602 {
		t.Fatalf("普通用户不应能审核视频: code=%d", code)
	}
	for _, table := range []string{"bus_opinion", "bus_videos"} {
		var status int
		if err := fx.d.QueryRow("SELECT status FROM " + table + " WHERE id=1").Scan(&status); err != nil || status != 0 || len(fx.amount) != 0 {
			t.Fatalf("%s 被改动或触发了打款: status=%d err=%v paid=%v", table, status, err, fx.amount)
		}
	}
	// 视频作者仍可改自己视频的非审核字段
	if code, _ := fx.call(t, "owner", "/saas-warehouse/busVideos/updateById", `{"id":1,"videoName":"renamed"}`); code != 200 {
		t.Fatalf("作者修改视频资料失败: code=%d", code)
	}
	// 平台管理员采纳意见:请求体把收款人写成自己也不生效,奖励打给意见提交人
	if code, body := fx.call(t, "ptadmin", "/saas-warehouse/busOpinion/updateById", `{"id":1,"status":1,"userId":3,"phone":"00000001124"}`); code != 200 ||
		len(fx.amount) != 2 || fx.amount[0] != 40 || fx.amount[1] != 1 {
		t.Fatalf("意见奖励未打给提交人: code=%d paid=%v %s", code, fx.amount, body)
	}
	// 平台管理员审核视频通过:奖励打给视频作者
	if code, body := fx.call(t, "ptadmin", "/saas-warehouse/busVideos/updateById", `{"id":1,"status":1}`); code != 200 ||
		len(fx.amount) != 4 || fx.amount[2] != 60 || fx.amount[3] != 1 {
		t.Fatalf("视频奖励未打给作者: code=%d paid=%v %s", code, fx.amount, body)
	}
}
