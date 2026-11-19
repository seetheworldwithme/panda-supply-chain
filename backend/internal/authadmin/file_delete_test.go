package authadmin

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

func TestFileDeletePreservesOriginalBugButDeletesConvertedPictures(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立库验收需显式启用")
	}
	admin, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	defer admin.Close()
	schema := fmt.Sprintf("audit_file_delete_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + schema)
	d, _ := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+schema)
	defer d.Close()
	for _, table := range []string{"chain_common_file", "chain_common_img", "chain_basic_enterprise"} {
		if _, err := d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	root, err := os.MkdirTemp("/tmp", "audit-file-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)
	for _, name := range []string{"original.pdf", "converted.png"} {
		if err = os.WriteFile(filepath.Join(root, name), []byte("fixture"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = d.Exec("INSERT INTO chain_common_file(id,file_path,file_name) VALUES(1,?,'original.pdf')", root); err != nil {
		t.Fatal(err)
	}
	if _, err = d.Exec("INSERT INTO chain_common_img(file_id,file_path,file_name) VALUES(1,?,'converted.png')", root); err != nil {
		t.Fatal(err)
	}
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	s.removeFileRow("1")
	if _, err = os.Stat(filepath.Join(root, "original.pdf")); err != nil {
		t.Fatal("原Java先删行再查原文件，原件应保留")
	}
	if _, err = os.Stat(filepath.Join(root, "converted.png")); !os.IsNotExist(err) {
		t.Fatal("原Java会删除转换图片文件")
	}
}
