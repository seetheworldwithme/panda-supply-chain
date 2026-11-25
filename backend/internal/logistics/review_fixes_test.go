package logistics

import (
	"strings"
	"testing"
)

// 手工拼接的查询条件必须转义单引号,否则 searchStr 就是 SQL 注入面。
func TestListSearchEscapesQuotes(t *testing.T) {
	sql := loListByPageSearch(map[string]any{"billNumber": "x' OR 1=1 -- ", "orderId": "0 OR 1=1"})
	if strings.Contains(sql, "x' OR") || !strings.Contains(sql, "x'' OR 1=1") || !strings.Contains(sql, "order_id ='0 OR 1=1'") {
		t.Fatalf("查询条件未转义: %s", sql)
	}
	if searchStrOf(map[string]any{"searchStr": "a'b"}) != "a''b" {
		t.Fatal("searchStrOf 未转义")
	}
}
