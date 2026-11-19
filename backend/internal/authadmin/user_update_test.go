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

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

func TestUserUpdateRouteUsesServerBusinessImplementation(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	defer admin.Close()
	schema := fmt.Sprintf("audit_user_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	if _, err := d.Exec("CREATE TABLE chain_auth_user LIKE chain.chain_auth_user"); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Exec("INSERT INTO chain_auth_user(id,account,handle_result,level,account_type,nickname) VALUES(1,'fixture',1,1,1,'old')"); err != nil {
		t.Fatal(err)
	}
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	kp := &auth.KeyPair{Private: key, Public: &key.PublicKey}
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	mux := http.NewServeMux()
	Register(mux, s)
	token, _ := kp.Generate(auth.JWTClaims{Account: "fixture", UserID: "1"})
	r := httptest.NewRequest("POST", "/saas-auth/chainAuthUser/updateById", strings.NewReader(`{"id":1,"nickname":"fixture-updated"}`))
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	auth.Middleware(kp, d)(mux).ServeHTTP(w, r)
	var nickname string
	if err := d.QueryRow("SELECT nickname FROM chain_auth_user WHERE id=1").Scan(&nickname); err != nil || nickname != "fixture-updated" || w.Code != 200 {
		t.Fatalf("用户修改入口仍被本机 JAR 缺陷覆盖: http=%d nickname=%s err=%v", w.Code, nickname, err)
	}
}
