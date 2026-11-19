package commonrest

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

// 以服务器运行包为基准，通过注册后的 HTTP 路由验收，防止内部函数存在但入口仍固定报错。

type wxFakeTransport struct{ host string }

func (t wxFakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = t.host
	return http.DefaultTransport.RoundTrip(req)
}

func wxEncryptPhone(t *testing.T, key []byte, phone string) (string, string) {
	t.Helper()
	iv := make([]byte, 16)
	rand.Read(iv)
	plain := []byte(fmt.Sprintf(`{"phoneNumber":%q}`, phone))
	pad := 16 - len(plain)%16
	for i := 0; i < pad; i++ {
		plain = append(plain, byte(pad))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	out := make([]byte, len(plain))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, plain)
	return base64.StdEncoding.EncodeToString(out), base64.StdEncoding.EncodeToString(iv)
}

func TestLoginChatWithFakeWeixin(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机隔离写入验收需显式启用")
	}
	d, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/chain?parseTime=true&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	const phone = "00000001120"
	const openid = "audit1118-wx-openid"
	var count int
	if err = d.QueryRow("SELECT count(*) FROM chain_auth_user WHERE phone=? OR openid=?", phone, openid).Scan(&count); err != nil || count != 0 {
		t.Fatalf("fixture 冲突: count=%d err=%v", count, err)
	}
	var pool int
	if err = d.QueryRow("SELECT mo_get_sum_score FROM bus_score_rule WHERE id=1").Scan(&pool); err != nil {
		t.Fatal(err)
	}
	sessionKey := make([]byte, 16)
	rand.Read(sessionKey)
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/sns/jscode2session" {
			fmt.Fprintf(w, `{"session_key":%q,"openid":%q}`, base64.StdEncoding.EncodeToString(sessionKey), openid)
			return
		}
		w.WriteHeader(500)
	}))
	defer fake.Close()
	defer func() {
		var uid string
		if d.QueryRow("SELECT id FROM chain_auth_user WHERE phone=? OR openid=?", phone, openid).Scan(&uid) == nil {
			d.Exec("DELETE FROM bus_score_log WHERE user_id=?", uid)
			d.Exec("DELETE FROM bus_user_score WHERE user_id=?", uid)
			d.Exec("DELETE FROM chain_auth_user WHERE id=?", uid)
		}
		d.Exec("UPDATE bus_score_rule SET mo_get_sum_score=? WHERE id=1", pool)
	}()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	s := NewService(&db.DB{Chain: d})
	s.Auth = auth.NewService(d, &auth.KeyPair{Private: key, Public: &key.PublicKey})
	s.HTTPClient = &http.Client{Transport: wxFakeTransport{host: fake.Listener.Addr().String()}}
	s.BatchMux = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"code":200,"data":"true"}`)) })
	mux := http.NewServeMux()
	s.registerChatWeixin(mux)

	call := func(body map[string]any) *httptest.ResponseRecorder {
		b, _ := json.Marshal(body)
		r := httptest.NewRequest("POST", "/saas-auth/chainAuthChat/loginChat", strings.NewReader(string(b)))
		r.Header.Set("User-Agent", "microMessenger probe")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w
	}

	// A：新手机号完整成功链（建用户+注册积分+名片查询+权益默认值）
	enc, iv := wxEncryptPhone(t, sessionKey, phone)
	w := call(map[string]any{"appIdentification": "chain", "code": "x", "encryptedData": enc, "iv": iv})
	if w.Code != 200 {
		t.Fatalf("成功链失败: %d %s", w.Code, w.Body.String())
	}
	var resp struct {
		Code int            `json:"code"`
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil || resp.Code != 200 {
		t.Fatalf("响应形状: %v %s", err, w.Body.String())
	}
	if resp.Data["newUser"] != true || resp.Data["openid"] != openid {
		t.Fatalf("newUser/openid 不符: %v", resp.Data)
	}
	if resp.Data["registScore"] == nil || resp.Data["MemberBenefits"] == nil || resp.Data["isCareOfficial"] != false {
		t.Fatalf("积分/权益/关注字段缺失: %v", resp.Data)
	}
	var nickname, mini string
	if err := d.QueryRow("SELECT nickname,is_mini_pro_user FROM chain_auth_user WHERE phone=?", phone).Scan(&nickname, &mini); err != nil {
		t.Fatal(err)
	}
	if nickname != "匿名用户"+phone[7:] || mini != "1" {
		t.Fatalf("落库昵称/小程序标记不符: %q %q", nickname, mini)
	}
	if err := d.QueryRow("SELECT count(*) FROM bus_score_log WHERE user_id=(SELECT id FROM chain_auth_user WHERE phone=?) AND score_get_type=1", phone).Scan(&count); err != nil || count != 1 {
		t.Fatalf("注册积分日志缺失: %d %v", count, err)
	}

	// B：searchUser + 已有 openid 无手机号 → 603 带用户实体
	d.Exec("UPDATE chain_auth_user SET phone=NULL,last_login_time='2026-01-01 00:00:00' WHERE phone=?", phone)
	w = call(map[string]any{"appIdentification": "chain", "code": "x", "searchUser": "1"})
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"code":603`) {
		t.Fatalf("603 分支不符: %d %s", w.Code, w.Body.String())
	}

	// C：searchUser + 未知 openid → 605 且新用户已提交（Java 自动提交语义）
	d.Exec("DELETE FROM bus_score_log WHERE user_id=(SELECT id FROM chain_auth_user WHERE openid=?)", openid)
	d.Exec("DELETE FROM bus_user_score WHERE user_id=(SELECT id FROM chain_auth_user WHERE openid=?)", openid)
	d.Exec("DELETE FROM chain_auth_user WHERE openid=?", openid)
	w = call(map[string]any{"appIdentification": "chain", "code": "x", "searchUser": "1"})
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"code":605`) {
		t.Fatalf("605 分支不符: %d %s", w.Code, w.Body.String())
	}
	if err := d.QueryRow("SELECT count(*) FROM chain_auth_user WHERE openid=? AND phone IS NULL", openid).Scan(&count); err != nil || count != 1 {
		t.Fatalf("605 新用户应已落库: %d %v", count, err)
	}
}

func TestPhoneLoginRoutesUseServerJavaBusinessHandlers(t *testing.T) {
	s := NewService(&db.DB{})
	s.Auth = auth.NewService(nil, nil)
	mux := http.NewServeMux()
	s.registerChatLocal(mux)
	for path, code := range map[string]string{"loginPhone": "601", "loginChatPhone": "602"} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/saas-auth/chainAuthChat/"+path, strings.NewReader(`{"phone":"fixture","code":"invalid"}`))
		mux.ServeHTTP(w, r)
		if w.Code != 200 || !strings.Contains(w.Body.String(), `"code":`+code) {
			t.Errorf("%s 未接业务处理函数: %d %s", path, w.Code, w.Body.String())
		}
	}
}
