package main

import (
	"strings"
	"testing"
)

// 日志里的 DSN 不能带密码:旧实现按"首个冒号下标<12"猜位置,用户名一长密码就原样输出。
func TestMaskDSNHidesPassword(t *testing.T) {
	for _, dsn := range []string{
		"root:p@ss:w0rd@tcp(127.0.0.1:3306)/chain",
		"saas_admin_user:s3cret@tcp(127.0.0.1:3307)/chain",
	} {
		masked := maskDSN(dsn)
		if strings.Contains(masked, "s3cret") || strings.Contains(masked, "w0rd") || !strings.Contains(masked, ":***@") {
			t.Fatalf("密码未打码: %q", masked)
		}
	}
	if got := maskDSN("root@tcp(127.0.0.1:3306)/chain"); got != "root:***@tcp(127.0.0.1:3306)/chain" {
		t.Fatalf("无密码 DSN 的用户名/地址应保留: %q", got)
	}
}
