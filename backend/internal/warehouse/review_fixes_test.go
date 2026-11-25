package warehouse

import "testing"

// 请求键名不能变成 SQL 片段:Java getSql 忽略实体里没有的键。
func TestKnownKeyRejectsUnknownKeys(t *testing.T) {
	if !entBusProduct.knownKey("id") {
		t.Fatal("实体字段应被识别")
	}
	for _, k := range []string{"id>0 or 1=1 or id", "a b", "id like '%' or id"} {
		if entBusProduct.knownKey(k) {
			t.Fatalf("未知键 %q 不应被当成列名", k)
		}
	}
	if got := sqlEscape("a'b\\c"); got != "a''b\\\\c" {
		t.Fatalf("sqlEscape = %q", got)
	}
}
