package warehouse

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

func TestLeasingWritesKeepAttachmentsOwnerAndOriginalPartialCommits(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("租售写入隔离数据库验收需显式启用")
	}
	admin, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := fmt.Sprintf("audit_leasing_%d", time.Now().UnixNano())
	if _, err = admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	for _, table := range []string{"chain_auth_user", "chain_common_file", "bus_enterprise_leasing"} {
		if _, err = d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{
		"INSERT INTO chain_auth_user(id,account,phone,channel,handle_result,level,account_type) VALUES(41,'fixture-leasing','00000001122','2',1,1,1)",
		"INSERT INTO chain_common_file(id,file_name,relation_id,effective) VALUES(1,'fixture-1.jpg',77,0),(2,'fixture-2.jpg',77,0),(3,'fixture-3.jpg',77,0),(4,'fixture-4.jpg',77,0),(5,'fixture-5.jpg',77,0)",
	} {
		if _, err = d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	kp := &auth.KeyPair{Private: key, Public: &key.PublicKey}
	token, _ := kp.Generate(auth.JWTClaims{Account: "fixture-leasing", UserID: "41"})
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	mux := http.NewServeMux()
	mux.HandleFunc("POST /add", authGuard(s.leasingAdd))
	mux.HandleFunc("POST /update", authGuard(s.leasingUpdate))
	request := func(path, body string) string {
		r := httptest.NewRequest("POST", path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		auth.Middleware(kp, d)(mux).ServeHTTP(w, r)
		return w.Body.String()
	}
	if response := request("/add", `{"enterpriseName":"fixture","status":0,"createUserId":999,"createPhone":"forged","cerList":[{"id":1}],"photoList":[{"id":2}]}`); !strings.Contains(response, `"code":200`) {
		t.Fatal(response)
	}
	var id, owner int
	var cert, photo, phone string
	if err = d.QueryRow("SELECT id,cert_association_id,photo_association_id,create_user_id,create_phone FROM bus_enterprise_leasing LIMIT 1").Scan(&id, &cert, &photo, &owner, &phone); err != nil {
		t.Fatal(err)
	}
	if len(cert) != 32 || len(photo) != 32 || cert == photo || owner != 41 || phone != "00000001122" {
		t.Fatalf("租售关联/发布人错误: %q %q %d %q", cert, photo, owner, phone)
	}
	checkFile := func(fileID, relation int, association string) {
		t.Helper()
		var ft, ow, rel, rt, active int
		var ot, assoc string
		if err := d.QueryRow("SELECT f_type,owner_type,owner_id,relation_id,relation_type,effective,association_id FROM chain_common_file WHERE id=?", fileID).Scan(&ft, &ot, &ow, &rel, &rt, &active, &assoc); err != nil {
			t.Fatal(err)
		}
		if ft != 4 || ot != "1" || ow != 41 || rel != relation || rt != 1 || active != 1 || assoc != association {
			t.Fatalf("附件%d未关联: %d %q %d %d %d %d %q", fileID, ft, ot, ow, rel, rt, active, assoc)
		}
	}
	// Java 保存附件时主表尚未自增，relationId=null 不覆盖已有关系值。
	checkFile(1, 77, cert)
	checkFile(2, 77, photo)
	if response := request("/update", fmt.Sprintf(`{"id":%d,"status":1,"cerList":[{"id":3}],"photoList":[{"id":4}]}`, id)); !strings.Contains(response, `"code":200`) {
		t.Fatal(response)
	}
	var updatedCert, updatedPhoto string
	if err = d.QueryRow("SELECT cert_association_id,photo_association_id FROM bus_enterprise_leasing WHERE id=?", id).Scan(&updatedCert, &updatedPhoto); err != nil {
		t.Fatal(err)
	}
	if updatedCert == cert || updatedPhoto == photo {
		t.Fatal("修改附件未重新生成关联编号")
	}
	checkFile(3, id, updatedCert)
	checkFile(4, id, updatedPhoto)
	// 主表写入失败时，之前独立提交的证书附件仍保留；不擅自增加总事务。
	response := request("/add", `{"enterpriseName":"`+strings.Repeat("x", 1024)+`","cerList":[{"id":5}]}`)
	if !strings.Contains(response, `"code":999`) {
		t.Fatal(response)
	}
	var associated string
	if err = d.QueryRow("SELECT association_id FROM chain_common_file WHERE id=5").Scan(&associated); err != nil || len(associated) != 32 {
		t.Fatalf("主表失败错误回滚附件: %q %v", associated, err)
	}
	// 照片 null 元素由服务抛错；之前证书批次已提交，主表尚未更新。
	response = request("/update", fmt.Sprintf(`{"id":%d,"status":0,"cerList":[{"id":5}],"photoList":[null]}`, id))
	if !strings.Contains(response, `"code":999`) {
		t.Fatal(response)
	}
	var changedAssociation string
	if err = d.QueryRow("SELECT association_id FROM chain_common_file WHERE id=5").Scan(&changedAssociation); err != nil || changedAssociation == associated {
		t.Fatalf("后批次失败错误回滚前批次: %q %v", changedAssociation, err)
	}
	var status int
	if err = d.QueryRow("SELECT status FROM bus_enterprise_leasing WHERE id=?", id).Scan(&status); err != nil || status != 1 {
		t.Fatalf("照片失败却更新了主记录: %d %v", status, err)
	}
	// 原旧行检查在更新附件之后，缺行不会阻止已经发生的附件更新。
	response = request("/update", `{"id":999,"status":1,"cerList":[{"id":5}]}`)
	if !strings.Contains(response, `"code":999`) {
		t.Fatal(response)
	}
	var relation int
	if err = d.QueryRow("SELECT relation_id FROM chain_common_file WHERE id=5").Scan(&relation); err != nil || relation != 999 {
		t.Fatalf("缺旧行未保留前置附件更新: %d %v", relation, err)
	}
}
