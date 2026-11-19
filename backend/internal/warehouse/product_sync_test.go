package warehouse

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/db"
)

func TestProductSyncPersistsPriceHistoryRewardsAndRelatedState(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("本机独立数据库验收需显式启用")
	}
	admin, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	name := fmt.Sprintf("audit_sync_%d", time.Now().UnixNano())
	if _, err = admin.Exec("CREATE DATABASE " + name); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + name)
	d, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+name)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	for _, table := range []string{"bus_product", "bus_price_log", "bus_pro_collection", "bus_msg", "bus_user_msg", "bus_videos", "chain_basic_enterprise", "chain_common_file", "bus_pro_log", "bus_red_money_clog", "bus_red_envelope_range", "chain_auth_user", "bus_member_benefits", "bus_config"} {
		if _, err = d.Exec("CREATE TABLE " + table + " LIKE chain." + table); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{
		"INSERT INTO chain_basic_enterprise(id,enterprise_name,is_authentication,group_id) VALUES(1,'fixture',1,3)",
		"INSERT INTO chain_auth_user(id,phone,member_type) VALUES(1,'00000001122',1)",
		"INSERT INTO bus_member_benefits(id,member_type,is_show_etc_name,is_show_etc_phone) VALUES(1,1,1,1)",
		"INSERT INTO bus_product(id,company_id,status,ex_factory_price,plate_price,create_user_id,create_phone,sec_name,warehouse_nb,enterprise_name,is_real_check) VALUES(1,1,1,110,210,1,'00000001122','fixture','fixture','fixture',1)",
		"INSERT INTO bus_price_log(user_id,pro_id,price,price_type) VALUES('2',1,100,1),('2',1,200,2)",
		"INSERT INTO bus_pro_collection(user_id,pro_id) VALUES(1,1)",
		"INSERT INTO bus_videos(id,pro_id,status) VALUES(1,1,0)",
		"INSERT INTO bus_red_envelope_range(receive_type,min_amount,max_amount) VALUES(8,5,5)",
	} {
		if _, err = d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PANDA_WAREHOUSE_RED_ENABLED", "0")
	t.Setenv("PANDA_WAREHOUSE_SEND_MSG_NOTICE", "false")
	t.Setenv("PANDA_WAREHOUSE_SEND_OFFICE_NOTICE", "false")
	s := NewService(&db.DB{Chain: d, Zjjg: d})
	r := httptest.NewRequest("POST", "/fixture", nil)
	if err = s.bpSyncOtherInfo(r, "1", "2", false, "1"); err != nil {
		t.Fatal(err)
	}
	var count int
	for _, q := range []string{"SELECT count(*) FROM bus_price_log WHERE pro_id=1", "SELECT count(*) FROM bus_red_money_clog WHERE receive_type=8", "SELECT count(*) FROM bus_msg WHERE msg_type=1"} {
		if err = d.QueryRow(q).Scan(&count); err != nil {
			t.Fatal(err)
		}
		want := 2
		if q == "SELECT count(*) FROM bus_price_log WHERE pro_id=1" {
			want = 4
		}
		if q == "SELECT count(*) FROM bus_msg WHERE msg_type=1" {
			want = 1
		}
		if count != want {
			t.Fatalf("%s got %d want %d", q, count, want)
		}
	}
	var oldFactory, oldPlate, lastFactory string
	var cert, group, video int
	if err = d.QueryRow("SELECT pre_ex_factory_price,pre_plate_price,last_check_fac_price,is_certified_etc,group_id FROM bus_product WHERE id=1").Scan(&oldFactory, &oldPlate, &lastFactory, &cert, &group); err != nil {
		t.Fatal(err)
	}
	if oldFactory != "100.00" || oldPlate != "200.00" || lastFactory != "110.00" || cert != 1 || group != 3 {
		t.Fatalf("产品同步字段错误: %s %s %s %d %d", oldFactory, oldPlate, lastFactory, cert, group)
	}
	d.QueryRow("SELECT status FROM bus_videos WHERE id=1").Scan(&video)
	if video != 1 {
		t.Fatal("审核未同步视频状态")
	}
	if got := s.bpOriginalUpdPriceScore(httptest.NewRecorder(), d, map[string]any{"status": "1"}, map[string]any{"id": "1", "status": "1", "exFactoryPrice": "120"}); got != 0 {
		t.Fatalf("相同状态的价格变动不应获得积分：%d", got)
	}
}
