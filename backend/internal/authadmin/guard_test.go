package authadmin

import (
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

// 账号边界:非平台管理员不能把自己或别人变成平台管理员,也不能动别人的账号。
func TestUserWriteBoundary(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	defer admin.Close()
	schema := fmt.Sprintf("audit_guard_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_auth_user", "chain_auth_role", "chain_auth_use_role", "chain_auth_menu_role", "chain_auth_menu", "chain_auth_menu_user"} {
		if _, err := d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	keep, _ := bcrypt.GenerateFromPassword([]byte("keep-me"), bcrypt.MinCost)
	if _, err := d.Exec(`INSERT INTO chain_auth_user(id,account,account_type,level,handle_result,company_id,password) VALUES
		(1,'ptadmin',2,2,1,NULL,?),(2,'normal',1,1,1,10,NULL),(3,'entadmin',3,1,1,10,NULL),
		(4,'colleague',1,1,1,10,NULL),(5,'other',1,1,1,20,NULL)`, string(keep)); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Exec("INSERT INTO chain_auth_role(id,name,role_type) VALUES(1,'pt',1),(2,'ent',2)"); err != nil {
		t.Fatal(err)
	}
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	kp := &auth.KeyPair{Private: key, Public: &key.PublicKey}
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	mux := http.NewServeMux()
	Register(mux, s)
	call := func(user, path, body string) *httptest.ResponseRecorder {
		uid, cid := map[string]string{"ptadmin": "1", "normal": "2", "entadmin": "3"}[user], "0"
		if user != "ptadmin" {
			cid = "10"
		}
		token, _ := kp.Generate(auth.JWTClaims{Account: user, UserID: uid, CompanyID: cid})
		r := httptest.NewRequest("POST", path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		auth.Middleware(kp, d)(mux).ServeHTTP(w, r)
		return w
	}
	var code int
	expect := func(name string, w *httptest.ResponseRecorder, want int) {
		t.Helper()
		if _, err := fmt.Sscanf(strings.SplitN(w.Body.String(), `"code":`, 2)[1], "%d", &code); err != nil || code != want {
			t.Fatalf("%s: want code %d, got %d %s", name, want, code, w.Body.String())
		}
	}
	// 普通用户:重置管理员密码、新建管理员账号、改同事 → 一律拒绝
	expect("resetPwd", call("normal", "/saas-auth/chainAuthUser/resetPwd", `{"id":1}`), 602)
	expect("add", call("normal", "/saas-auth/chainAuthUser/add", `{"account":"x","password":"x","accountType":5,"level":2,"handleResult":1}`), 602)
	expect("update other", call("normal", "/saas-auth/chainAuthUser/updateById", `{"id":4,"nickname":"h"}`), 602)
	// 普通用户改自己:资料生效,平台属性 / 状态 / 会员字段被丢弃
	expect("update self", call("normal", "/saas-auth/chainAuthUser/updateById", `{"id":2,"nickname":"n","accountType":"5","level":2,"companyId":20,"handleResult":1}`), 200)
	expect("update self member", call("normal", "/saas-auth/chainAuthUser/updateById", `{"id":2,"memberType":1,"expirationTime":"2099-12-31 00:00:00","activateType":1}`), 200)
	var accountType, level, company, handle, member, nickname sql.NullString
	if err := d.QueryRow("SELECT account_type,level,company_id,handle_result,member_type,nickname FROM chain_auth_user WHERE id=2").
		Scan(&accountType, &level, &company, &handle, &member, &nickname); err != nil {
		t.Fatal(err)
	}
	if accountType.String != "1" || level.String != "1" || company.String != "10" || member.String != "0" || nickname.String != "n" { // member_type 列默认 0 = 非会员
		t.Fatalf("普通用户自改写入了受限字段: type=%s level=%s company=%s member=%s nick=%s", accountType.String, level.String, company.String, member.String, nickname.String)
	}
	// 企业管理员:能停用本企业账号,动不了平台管理员,也分不到平台角色
	expect("ent update peer", call("entadmin", "/saas-auth/chainAuthUser/updateById", `{"id":4,"handleResult":2}`), 200)
	expect("ent update ptadmin", call("entadmin", "/saas-auth/chainAuthUser/updateById", `{"id":1,"handleResult":2}`), 602)
	expect("ent assign pt role", call("entadmin", "/saas-auth/chainAuthUserList/assignRole", `{"userId":4,"userRoles":[{"role_id":1}]}`), 602)
	expect("ent assign ent role", call("entadmin", "/saas-auth/chainAuthUserList/assignRole", `{"userId":4,"userRoles":[{"role_id":2}]}`), 200)
	// 平台管理员:原有能力不变
	expect("pt resetPwd", call("ptadmin", "/saas-auth/chainAuthUser/resetPwd", `{"id":2}`), 200)
	var hash string
	if err := d.QueryRow("SELECT password FROM chain_auth_user WHERE id=2").Scan(&hash); err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte("123456")) != nil {
		t.Fatalf("平台管理员重置密码未生效: %v", err)
	}
	if err := d.QueryRow("SELECT password FROM chain_auth_user WHERE id=1").Scan(&hash); err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte("keep-me")) != nil {
		t.Fatalf("平台管理员密码被改动: %v", err)
	}
}

// 忘记密码:未登录时 updPwd 只能改"刚在 valCode 通过短信码校验的手机号"下的账号,不能凭 id 改任何人的密码。
func TestSetPwdRequiresVerifiedPhoneForAnonymous(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	defer admin.Close()
	schema := fmt.Sprintf("audit_setpwd_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_auth_user", "chain_auth_menu_role", "chain_auth_menu", "chain_auth_menu_user", "chain_auth_use_role"} {
		if _, err := d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	keep, _ := bcrypt.GenerateFromPassword([]byte("keep-me"), bcrypt.MinCost)
	if _, err := d.Exec("INSERT INTO chain_auth_user(id,account,account_type,level,handle_result,phone,password) VALUES(1,'ptadmin',2,2,1,'00000001199',?),(2,'normal',1,1,1,'00000001188',?)", string(keep), string(keep)); err != nil {
		t.Fatal(err)
	}
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	s.Auth = auth.NewService(d, nil)
	w := httptest.NewRecorder()
	s.userSetPwd(w, httptest.NewRequest("POST", "/saas-auth/chainAuthUser/updPwd", strings.NewReader(`{"id":1,"newPassword":"Hacked123"}`)))
	if !strings.Contains(w.Body.String(), `"code":602`) {
		t.Fatalf("未验证身份就能改密码: %s", w.Body.String())
	}
	// 普通用户通过了 1388 的短信码校验:只能改 1388 的账号,改不了管理员的
	s.Auth.MarkPhoneVerified("00000001188")
	w = httptest.NewRecorder()
	s.userSetPwd(w, httptest.NewRequest("POST", "/saas-auth/chainAuthUser/updPwd", strings.NewReader(`{"id":1,"newPassword":"Hacked123"}`)))
	if !strings.Contains(w.Body.String(), `"code":602`) {
		t.Fatalf("验证的是别的手机号却改了管理员密码: %s", w.Body.String())
	}
	w = httptest.NewRecorder()
	s.userSetPwd(w, httptest.NewRequest("POST", "/saas-auth/chainAuthUser/updPwd", strings.NewReader(`{"id":2,"newPassword":"Hacked123"}`)))
	var hash string
	if err := d.QueryRow("SELECT password FROM chain_auth_user WHERE id=2").Scan(&hash); err != nil || !strings.Contains(w.Body.String(), `"code":200`) || bcrypt.CompareHashAndPassword([]byte(hash), []byte("Hacked123")) != nil {
		t.Fatalf("本人通过短信码校验后应能改自己的密码: %s %v", w.Body.String(), err)
	}
	if err := d.QueryRow("SELECT password FROM chain_auth_user WHERE id=1").Scan(&hash); err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte("keep-me")) != nil {
		t.Fatalf("管理员密码被改动: %v", err)
	}
}

// 菜单增删改、全局配置、e签宝配置、平台角色与权限分配:只有平台管理员能调。
func TestPlatformOnlyWriteEndpoints(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	defer admin.Close()
	schema := fmt.Sprintf("audit_ptonly_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_auth_user", "chain_auth_menu", "chain_auth_menu_authy", "chain_auth_role", "chain_auth_menu_role", "chain_auth_use_role", "chain_auth_menu_user",
		"chain_auth_authority", "chain_auth_authority_role", "chain_auth_authority_user", "chain_global_configuration", "chain_eqb_config", "bus_red_envelope_range"} {
		if _, err := d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	keep, _ := bcrypt.GenerateFromPassword([]byte("keep-me"), bcrypt.MinCost)
	if _, err := d.Exec("INSERT INTO chain_auth_user(id,account,account_type,level,handle_result,company_id,password) VALUES(1,'ptadmin',2,2,1,NULL,?),(2,'entadmin',3,1,1,10,?),(3,'normal',1,1,1,10,?)", string(keep), string(keep), string(keep)); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Exec("INSERT INTO chain_auth_menu(id,name,route_name) VALUES(1,'fixture-menu','UserList')"); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Exec("INSERT INTO chain_global_configuration(id,red_uplimit_days) VALUES(1,0)"); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Exec("INSERT INTO bus_red_envelope_range(id,min_amount,max_amount,receive_type) VALUES(1,100,200,3)"); err != nil {
		t.Fatal(err)
	}
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	kp := &auth.KeyPair{Private: key, Public: &key.PublicKey}
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	mux := http.NewServeMux()
	Register(mux, s)
	guarded := auth.Middleware(kp, d)(mux)
	call := func(user string, path, body string) int {
		token, _ := kp.Generate(auth.JWTClaims{Account: user, UserID: map[string]string{"ptadmin": "1", "entadmin": "2", "normal": "3"}[user], CompanyID: map[string]string{"entadmin": "10", "normal": "10"}[user]})
		r := httptest.NewRequest("POST", path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		guarded.ServeHTTP(w, r)
		var code int
		fmt.Sscanf(strings.SplitN(w.Body.String(), `"code":`, 2)[1], "%d", &code)
		return code
	}
	// 菜单名、红包金额下限、提现上限天数:被拒后必须原样,平台管理员调用后必须改掉
	read := func() (name string, min, days int) {
		d.QueryRow("SELECT name FROM chain_auth_menu WHERE id=1").Scan(&name)
		d.QueryRow("SELECT min_amount FROM bus_red_envelope_range WHERE id=1").Scan(&min)
		d.QueryRow("SELECT red_uplimit_days FROM chain_global_configuration WHERE id=1").Scan(&days)
		return
	}
	rangeBody := `[{"id":1,"minAmount":101,"maxAmount":201,"receiveType":3,"redUplimitDays":7}]`
	for _, path := range []string{"/saas-auth/chainAuthMenu/add", "/saas-auth/chainAuthMenu/updateById", "/saas-auth/chainAuthMenu/deleteById", "/saas-auth/chainAuthMenu/delBatch",
		"/saas-auth/chainAuthRole/addMenu", "/saas-auth/chainAuthRole/addPtRole", "/saas-auth/chainAuthRole/savePtAuths", "/saas-auth/chainAuthRole/updRoleForPt",
		"/saas-common/chainGlobalConfiguration/updateById", "/saas-common/chainEqbConfig/add", "/saas-common/chainEqbConfig/updateById", "/saas-common/chainEqbConfig/delBatch",
		"/saas-warehouse/busRedEnvelopeRange/updBatch"} {
		for _, user := range []string{"entadmin", "normal"} {
			body := `{"id":1,"name":"x","roleType":1,"redUplimitDays":7}`
			if strings.HasSuffix(path, "busRedEnvelopeRange/updBatch") {
				body = rangeBody
			}
			if code := call(user, path, body); code != 602 {
				t.Errorf("%s 调 %s 应被拒绝,得到 %d", user, path, code)
			}
		}
	}
	if name, min, days := read(); name != "fixture-menu" || min != 100 || days != 0 {
		t.Fatalf("被拒的请求改动了数据: %s %d %d", name, min, days)
	}
	// 用切片固定顺序:全局配置与 updBatch 改的是同一列,后调用的生效
	for _, c := range []struct{ path, body string }{
		{"/saas-auth/chainAuthMenu/updateById", `{"id":1,"name":"renamed"}`},
		{"/saas-common/chainGlobalConfiguration/updateById", `{"id":1,"redUplimitDays":9}`},
		{"/saas-warehouse/busRedEnvelopeRange/updBatch", rangeBody},
	} {
		if code := call("ptadmin", c.path, c.body); code != 200 {
			t.Errorf("平台管理员调 %s 应成功,得到 %d", c.path, code)
		}
	}
	// updBatch 排在全局配置之后调用:它按首行的 redUplimitDays 覆盖同一列,最终值是 7
	if name, min, days := read(); name != "renamed" || min != 101 || days != 7 {
		t.Errorf("平台管理员的修改未生效: %s %d %d", name, min, days)
	}
}
