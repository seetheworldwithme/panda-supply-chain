package contract

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

// delOnce:无条件、条件全为空值时都不能清库;按批次号删除照常。
func TestDelOnceRefusesEmptyAndAllNullConditions(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	defer admin.Close()
	schema := fmt.Sprintf("audit_delonce_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_contract", "chain_capital_settlement", "chain_capital_detailed", "chain_bill", "chain_bill_rate", "chain_logistics", "chain_logistics_weighing_list"} {
		if _, err := d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := d.Exec("INSERT INTO chain_contract(id, im_batch_number) VALUES (1, NULL), (2, 'B1')"); err != nil {
		t.Fatal(err)
	}
	s := &Service{DB: &db.DB{Chain: d, Zjjg: d}}
	call := func(body string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		s.delOnce(w, httptest.NewRequest("POST", "/saas-contract/chainContract/delOnce", strings.NewReader(body)))
		return w
	}
	count := func() int {
		var n int
		if err := d.QueryRow("SELECT COUNT(*) FROM chain_contract").Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	for _, body := range []string{`{}`, `{"im_batch_number":null}`} {
		if w := call(body); strings.Contains(w.Body.String(), `"code":200`) || count() != 2 {
			t.Fatalf("条件为空/全为 null 时不得删除: body=%s resp=%s rows=%d", body, w.Body.String(), count())
		}
	}
	if w := call(`{"im_batch_number":"B1"}`); !strings.Contains(w.Body.String(), `"code":200`) || count() != 1 {
		t.Fatalf("按批次号删除应照常生效: resp=%s rows=%d", w.Body.String(), count())
	}
}
