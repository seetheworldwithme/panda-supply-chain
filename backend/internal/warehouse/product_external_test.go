package warehouse

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

type productWXTransport struct{ host string }

func (t productWXTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "http"
	r.URL.Host = t.host
	return http.DefaultTransport.RoundTrip(r)
}

func TestProductExternalUsesWechatAndCreatesSharePicture(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立库和原生渲染验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	defer admin.Close()
	schema := fmt.Sprintf("audit_product_external_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_auth_chat_authorization", "chain_auth_user", "bus_member_benefits", "bus_product", "chain_common_file"} {
		if _, err := d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{
		"INSERT INTO chain_auth_chat_authorization(app_identification,appid,secret) VALUES('cloudAccount','fixture-app','fixture-secret')",
		"INSERT INTO chain_auth_user(id,member_type,official_openid) VALUES(1,1,'fixture-openid')",
		"INSERT INTO bus_member_benefits(member_type,is_show_etc_name,is_show_etc_phone) VALUES(1,1,1)",
		"INSERT INTO bus_product(id,create_user_id,warehouse_nb,enterprise_name,ex_factory_price,pre_ex_factory_price,coal_rel_name,coal_rel_phone,business_name,business_phone,max_std) VALUES(1,1,'MC001','测试煤矿',100,95,'联系人','00000000001','业务员','00000000002',0.8)",
	} {
		if _, err := d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	sends := 0
	appSends := 0
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/stable_token":
			fmt.Fprint(w, `{"access_token":"fixture-token"}`)
		case "/cgi-bin/user/info":
			if r.URL.Query().Get("openid") != "fixture-openid" {
				t.Error("关注查询未透传openid")
			}
			fmt.Fprint(w, `{"subscribe":1}`)
		case "/cgi-bin/message/template/send":
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			sends++
			if sends == 1 && len(body) != 0 {
				t.Error("测试推送未保留旧控制器空消息")
			}
			if sends > 1 {
				data, _ := body["data"].(map[string]any)
				if body["touser"] != "fixture-openid" || data["thing9"] == nil {
					t.Error("业务通知缺收件人或价格内容")
				}
			}
			if r.URL.Query().Get("access_token") != "fixture-token" {
				t.Error("消息未使用应用token")
			}
			fmt.Fprint(w, `{"errcode":0}`)
		case "/push":
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["title"] != "今日新增煤种" {
				t.Error("App推送标题不一致")
			}
			appSends++
			fmt.Fprint(w, `{}`)
		default:
			t.Errorf("意外外部请求 %s", r.URL.Path)
			w.WriteHeader(500)
		}
	}))
	defer fake.Close()
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	s.HTTPClient = &http.Client{Transport: productWXTransport{fake.Listener.Addr().String()}}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{}`))
	s.bpPushOfficialActual(w, r)
	if !strings.Contains(w.Body.String(), `"code":200`) || sends != 1 {
		t.Fatalf("推送未发送 %s", w.Body.String())
	}
	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/", strings.NewReader(`{"openid":"fixture-openid"}`))
	s.bpGetCareActual(w, r)
	if !strings.Contains(w.Body.String(), `"data":"已关注"`) {
		t.Fatalf("关注结果 %s", w.Body.String())
	}
	t.Setenv("PANDA_WAREHOUSE_SEND_OFFICE_NOTICE", "true")
	t.Setenv("PANDA_WAREHOUSE_UNIPUSH_ADR", fake.URL+"/push")
	p := map[string]any{"id": 1, "createUserId": 1, "warehouseNb": "MC001", "enterpriseName": "测试煤矿", "exFactoryPrice": "100.00", "maxStd": "0.80"}
	if err := s.bpNotifyAdded(r, p); err != nil || sends != 2 || appSends != 1 {
		t.Fatalf("新增通知未发送或不保留字段: sends=%d app=%d err=%v", sends, appSends, err)
	}
	if err := s.bpNotifyAdded(r, p); err != nil || sends != 2 {
		t.Fatal("官方消息未按原缓存去重")
	}
	root, err := os.MkdirTemp("/tmp", "audit-share-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)
	t.Setenv("PANDA_UPLOAD_DIR", root)
	f, err := os.Create(filepath.Join(root, "frame73.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	jpeg.Encode(f, image.NewRGBA(image.Rect(0, 0, 750, 1200)), &jpeg.Options{Quality: 80})
	f.Close()
	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/", strings.NewReader(`{"id":1}`))
	name, renderErr := s.bpSharePicture(r, 1)
	if renderErr != nil {
		t.Fatal(renderErr)
	}
	platform.Success(w, name)
	var res struct{ Data string }
	if err = json.Unmarshal(w.Body.Bytes(), &res); err != nil || res.Data == "" {
		t.Fatalf("未生成分享图 %s", w.Body.String())
	}
	f, err = os.Open(filepath.Join(root, res.Data))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	im, err := jpeg.Decode(f)
	if err != nil || im.Bounds().Dx() != 750 || im.Bounds().Dy() != 1200 {
		t.Fatalf("图片尺寸或格式错误 %v", err)
	}
	var count int
	d.QueryRow("SELECT count(*) FROM chain_common_file WHERE file_name=?", res.Data).Scan(&count)
	if count != 1 {
		t.Error("分享图未登记附件")
	}
	if err = os.Remove(filepath.Join(root, "frame73.jpg")); err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/", strings.NewReader(`{"id":1}`))
	s.bpGetDetailImgActual(w, r)
	if !strings.Contains(w.Body.String(), `"data":null`) {
		t.Error("模板缺失未保留Java空结果")
	}
}
