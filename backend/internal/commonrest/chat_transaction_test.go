package commonrest

import (
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

// 仅显式启用时操作本机 fixture，不读取部署环境数据库配置。
func TestPhoneLoginRollsBackUserWhenCardServiceFails(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机隔离写入验收需显式启用")
	}
	d, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/chain?parseTime=true&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	const phone = "00000001119"
	var count int
	if err = d.QueryRow("SELECT count(*) FROM chain_auth_user WHERE phone=?", phone).Scan(&count); err != nil || count != 0 {
		t.Fatalf("fixture 发生冲突: count=%d err=%v", count, err)
	}
	defer d.Exec("DELETE FROM chain_auth_user WHERE phone=?", phone)
	s := NewService(&db.DB{Chain: d})
	s.Auth = auth.NewService(d, nil)
	s.BatchMux = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) })
	code := s.Auth.IssuePhoneCode(phone)
	r := httptest.NewRequest("POST", "/saas-auth/chainAuthChat/loginPhone", strings.NewReader(fmt.Sprintf(`{"phone":%q,"code":%q}`, phone, code)))
	w := httptest.NewRecorder()
	s.chatPhoneLogin(w, r, false)
	if w.Code != 500 {
		t.Fatalf("模拟依赖故障未触发: %d", w.Code)
	}
	if err = d.QueryRow("SELECT count(*) FROM chain_auth_user WHERE phone=?", phone).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("名片依赖失败后，新用户仍已提交")
	}
}

type failingTransport struct{}

func (failingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("模拟微信接口故障")
}

// 老用户 openid/app_client_id 已写入事务后，公众号关注依赖失败必须整体回滚。
func TestPhoneLoginRollsBackExistingUserWhenOfficialCareFails(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机隔离写入验收需显式启用")
	}
	d, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/chain?parseTime=true&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	const phone = "00000001119"
	var count int
	if err = d.QueryRow("SELECT count(*) FROM chain_auth_user WHERE phone=?", phone).Scan(&count); err != nil || count != 0 {
		t.Fatalf("fixture 发生冲突: count=%d err=%v", count, err)
	}
	// cloudAccount 应用配置必须存在，保证失败点在微信 HTTP 依赖而不是更早的配置查询。
	if err = d.QueryRow("SELECT count(*) FROM chain_auth_chat_authorization WHERE app_identification='cloudAccount'").Scan(&count); err != nil || count != 1 {
		t.Fatalf("cloudAccount 配置不唯一，无法定位失败点: count=%d err=%v", count, err)
	}
	res, err := d.Exec("INSERT INTO chain_auth_user(phone,account,level,channel,nickname,is_agree_call,openid,app_client_id,official_openid) VALUES (?,?,'1','2','回滚fixture',1,'fixture-old-openid','fixture-old-client','fixture-official-openid')", phone, "fixture-rollback-account")
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()
	defer d.Exec("DELETE FROM chain_auth_user WHERE id=?", id)
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	s := NewService(&db.DB{Chain: d})
	s.Auth = auth.NewService(d, &auth.KeyPair{Private: key, Public: &key.PublicKey})
	s.HTTPClient = &http.Client{Transport: failingTransport{}}
	s.BatchMux = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"code":200,"data":"true"}`)) // 已有默认名片，跳过 add
	})
	code := s.Auth.IssuePhoneCode(phone)
	r := httptest.NewRequest("POST", "/saas-auth/chainAuthChat/loginChatPhone", strings.NewReader(fmt.Sprintf(`{"phone":%q,"code":%q,"openid":"fixture-new-openid","pushClientid":"fixture-new-client"}`, phone, code)))
	w := httptest.NewRecorder()
	s.chatPhoneLogin(w, r, true)
	if w.Code != 500 {
		t.Fatalf("模拟公众号依赖故障未触发: %d %s", w.Code, w.Body.String())
	}
	var openid, client string
	if err = d.QueryRow("SELECT openid,app_client_id FROM chain_auth_user WHERE id=?", id).Scan(&openid, &client); err != nil {
		t.Fatal(err)
	}
	if openid != "fixture-old-openid" || client != "fixture-old-client" {
		t.Fatalf("公众号依赖失败后，老用户 openid/app_client_id 仍被提交: %q %q", openid, client)
	}
}
