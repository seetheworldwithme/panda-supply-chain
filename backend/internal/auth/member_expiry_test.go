package auth

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type expiryRedisFixture struct {
	sync.Mutex
	values  map[string]string
	expires map[string]int64
}

func startExpiryRedisFixture(t *testing.T) *expiryRedisFixture {
	t.Helper()
	f := &expiryRedisFixture{values: map[string]string{}, expires: map[string]int64{}}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { listener.Close() })
	t.Setenv("PANDA_REDIS_ADDR", listener.Addr().String())
	t.Setenv("PANDA_REDIS_PASSWORD", "")
	t.Setenv("PANDA_REDIS_DB", "0")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				reader := bufio.NewReader(conn)
				for {
					line, err := reader.ReadString('\n')
					if err != nil {
						return
					}
					n, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "*")))
					args := make([]string, n)
					for i := range args {
						line, err = reader.ReadString('\n')
						if err != nil {
							return
						}
						size, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "$")))
						data := make([]byte, size+2)
						if _, err = io.ReadFull(reader, data); err != nil {
							return
						}
						args[i] = string(data[:size])
					}
					f.Lock()
					switch args[0] {
					case "SET":
						f.values[args[1]] = args[2]
						fmt.Fprint(conn, "+OK\r\n")
					case "SETEX":
						f.values[args[1]] = args[3]
						seconds, _ := strconv.ParseInt(args[2], 10, 64)
						f.expires[args[1]] = seconds
						fmt.Fprint(conn, "+OK\r\n")
					case "PEXPIREAT":
						milliseconds, _ := strconv.ParseInt(args[2], 10, 64)
						f.expires[args[1]] = milliseconds
						fmt.Fprint(conn, ":1\r\n")
					case "GET":
						value, ok := f.values[args[1]]
						if !ok {
							fmt.Fprint(conn, "$-1\r\n")
						} else {
							fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(value), value)
						}
					default:
						fmt.Fprint(conn, "+OK\r\n")
					}
					f.Unlock()
				}
			}()
		}
	}()
	return f
}

func TestMemberExpiryKeepsDowngradeAndOriginalDelayedSMSCondition(t *testing.T) {
	if os.Getenv("PANDA_LOCAL_AUDIT") != "1" {
		t.Skip("会员降级隔离数据库验收需显式启用")
	}
	admin, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/")
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	name := fmt.Sprintf("audit_member_%d", time.Now().UnixNano())
	if _, err = admin.Exec("CREATE DATABASE " + name); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP DATABASE " + name)
	d, err := sql.Open("mysql", "root@unix(/tmp/mysql.sock)/"+name)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	if _, err = d.Exec("CREATE TABLE chain_auth_user LIKE chain.chain_auth_user"); err != nil {
		t.Fatal(err)
	}
	if _, err = d.Exec("INSERT INTO chain_auth_user(id,phone,member_type) VALUES(42,'00000001122',3)"); err != nil {
		t.Fatal(err)
	}
	redis := startExpiryRedisFixture(t)
	var sent atomic.Int32
	fakeSMS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body["TemplateID"] != "2165510" {
			t.Errorf("到期模板不符: %v", body["TemplateID"])
		}
		sent.Add(1)
		fmt.Fprint(w, `{"Response":{"RequestId":"fixture"}}`)
	}))
	defer fakeSMS.Close()
	t.Setenv("TENCENT_SMS_ENDPOINT", fakeSMS.URL)
	t.Setenv("TENCENT_SECRET_ID", "fixture")
	t.Setenv("TENCENT_SECRET_KEY", "fixture")
	s := NewService(d, nil)
	noon := time.Date(2026, 11, 19, 12, 0, 0, 0, time.Local)
	if err = s.handleMemberExpired(context.Background(), "member42", noon); err != nil {
		t.Fatal(err)
	}
	var member int
	if err = d.QueryRow("SELECT member_type FROM chain_auth_user WHERE id=42").Scan(&member); err != nil || member != 0 {
		t.Fatalf("会员未降级: %d %v", member, err)
	}
	if sent.Load() != 1 {
		t.Fatal("白天到期未发短信")
	}
	night := time.Date(2026, 11, 19, 23, 15, 0, 0, time.Local)
	if err = s.handleMemberExpired(context.Background(), "member42", night); err != nil {
		t.Fatal(err)
	}
	redis.Lock()
	expiry := redis.expires["delayed-00000001122"]
	redis.Unlock()
	if want := time.Date(2026, 11, 20, 9, 30, 0, 0, time.Local).UnixMilli(); expiry != want {
		t.Fatalf("夜间提醒到期错误 %d want %d", expiry, want)
	}
	if sent.Load() != 1 {
		t.Fatal("夜间不得立即发送")
	}
	if err = s.handleMemberExpired(context.Background(), "delayed-00000001122", noon); err != nil {
		t.Fatal(err)
	}
	if sent.Load() != 1 {
		t.Fatal("原Java冷sendMsg缓存不发送，不能修复此问题")
	}
	redis.Lock()
	ttl := redis.expires["sendMsg00000001122"]
	redis.Unlock()
	if ttl != 3600 {
		t.Fatalf("发送标记TTL错误 %d", ttl)
	}
	if err = s.handleMemberExpired(context.Background(), "delayed-00000001122", noon); err != nil {
		t.Fatal(err)
	}
	if sent.Load() != 2 {
		t.Fatal("原Java热sendMsg缓存应发送")
	}
}

func TestMemberReminderTimePreservesExactJavaTimeBoundaries(t *testing.T) {
	for _, tc := range []struct{ h, m, s, ns, days int }{{7, 59, 0, 0, 0}, {7, 59, 0, 1, 1}, {23, 0, 0, 0, 1}, {0, 0, 0, 0, 0}} {
		now := time.Date(2026, 11, 19, tc.h, tc.m, tc.s, tc.ns, time.Local)
		want := time.Date(2026, 11, 19+tc.days, 9, 30, 0, 0, time.Local)
		if got := memberReminderTime(now); !got.Equal(want) {
			t.Fatalf("%v -> %v want %v", now, got, want)
		}
	}
}

func TestScheduleMemberExpiryUsesOriginalKeyAndMilliseconds(t *testing.T) {
	redis := startExpiryRedisFixture(t)
	expires := time.Date(2026, 12, 1, 12, 34, 56, 789000000, time.Local)
	if err := ScheduleMemberExpiry(context.Background(), "42", expires); err != nil {
		t.Fatal(err)
	}
	redis.Lock()
	defer redis.Unlock()
	if redis.values["member42"] != "11" || redis.expires["member42"] != expires.UnixMilli() {
		t.Fatalf("原会员缓存键或毫秒到期不符: %v %v", redis.values, redis.expires)
	}
}
