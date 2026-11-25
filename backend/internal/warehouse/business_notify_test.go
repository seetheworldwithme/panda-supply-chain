package warehouse

import (
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

func TestBusinessSendsOriginalSMSAfterDatabaseWrites(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立库验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	defer admin.Close()
	schema := fmt.Sprintf("audit_notify_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_auth_user", "bus_product", "bus_config", "bus_order", "bus_videos", "chain_common_file", "bus_user_care", "bus_pro_log", "bus_red_envelope_range", "bus_red_money_clog"} {
		if _, err := d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{
		"INSERT INTO chain_auth_user(id,account,phone,channel,handle_result,level,account_type) VALUES(1,'fixture','00000001121','2',1,1,1),(2,'fixture-admin','00000001129','2',1,2,2)",
		"INSERT INTO bus_product(id,warehouse_nb,enterprise_name,notice_time,ex_factory_price) VALUES(1,'fixture-warehouse','fixture-enterprise',0,100)",
		"INSERT INTO bus_config(id,admin_phone,msg_type) VALUES(1,'00000001122',1)",
		"INSERT INTO chain_common_file(id,file_name) VALUES(1,'fixture.mp4')",
		"INSERT INTO bus_red_envelope_range(receive_type,min_amount,max_amount) VALUES(11,5,5)",
		"UPDATE chain_auth_user SET create_time='2000-01-01' WHERE id=1",
	} {
		if _, err := d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	var received []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		received = append(received, body)
		io.WriteString(w, `{"Response":{"RequestId":"fixture"}}`)
	}))
	defer server.Close()
	t.Setenv("TENCENT_SMS_ENDPOINT", server.URL)
	t.Setenv("TENCENT_SECRET_ID", "fixture")
	t.Setenv("TENCENT_SECRET_KEY", "fixture")
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	kp := &auth.KeyPair{Private: key, Public: &key.PublicKey}
	token, _ := kp.Generate(auth.JWTClaims{Account: "fixture", UserID: "1"})
	adminToken, _ := kp.Generate(auth.JWTClaims{Account: "fixture-admin", UserID: "2"})
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	mux := http.NewServeMux()
	Register(mux, s)
	for _, tc := range []struct{ path, body, template, as string }{
		{"/busCustManage/noticeCust", `{"proId":1,"phones":["00000001123","00000001124"]}`, "2115763", "author"},
		{"/busOrder/placeOrder", `{"proId":1}`, "1990680", "author"},
		{"/busVideos/add", `{"videoId":1}`, "1792052", "author"},
		{"/busVideos/updateById", `{"id":1,"status":0}`, "1792052", "author"}, // 作者提交审核
		{"/busVideos/updateById", `{"id":1,"status":2}`, "1794840", "admin"},  // 审核拒绝:平台管理员
	} {
		t.Run(tc.template, func(t *testing.T) {
			before := len(received)
			r := httptest.NewRequest("POST", "/saas-warehouse"+tc.path, strings.NewReader(tc.body))
			r.Header.Set("Authorization", "Bearer "+map[string]string{"author": token, "admin": adminToken}[tc.as])
			w := httptest.NewRecorder()
			auth.Middleware(kp, d)(mux).ServeHTTP(w, r)
			if !strings.Contains(w.Body.String(), `"code":200`) {
				t.Fatal(w.Body.String())
			}
			if len(received) != before+1 || received[before]["TemplateID"] != tc.template {
				t.Fatalf("原业务短信未发送: %v", received[before:])
			}
		})
	}
	r := httptest.NewRequest("POST", "/saas-warehouse/busVideos/getById", strings.NewReader(`{"id":1,"shareUserId":2,"busProLog":{"userId":1,"videoId":1,"operateObj":3,"operateType":7,"shareUserId":2}}`))
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	auth.Middleware(kp, d)(mux).ServeHTTP(w, r)
	if !strings.Contains(w.Body.String(), `"code":200`) {
		t.Fatal(w.Body.String())
	}
	var generated, visits int
	d.QueryRow("SELECT count(*) FROM bus_red_money_clog WHERE receive_type=11").Scan(&generated)
	d.QueryRow("SELECT count(*) FROM bus_pro_log WHERE video_id=1").Scan(&visits)
	if generated != 2 || visits != 1 {
		t.Fatalf("未保留视频分享分支 count为0仍奖励的原问题: rewards=%d visits=%d", generated, visits)
	}
}
