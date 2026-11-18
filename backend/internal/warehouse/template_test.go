package warehouse

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExcelMetadataPreservesTagColumns(t *testing.T) {
	if len(entBusTag.excel) != 8 {
		t.Fatalf("BusTag export columns=%d, want 8 Java @ExcelColumn headers", len(entBusTag.excel))
	}
	if entBusTag.excel[0].header != "标签名称" {
		t.Fatalf("first header = %q", entBusTag.excel[0].header)
	}
}

func TestOrderQualityColumnsUseMybatisPropertyMatching(t *testing.T) {
	row := map[string]any{"p_ash_content": "9.50", "p_bond_index": "60.00", "p_water_content": "10.00"}
	got := entBusOrder.shapeRows([]map[string]any{row}, []string{"p_ash_content", "p_bond_index", "p_water_content"})[0]
	for _, key := range []string{"pashcontent", "pbondindex", "pwatercontent"} {
		if got[key] == nil {
			t.Errorf("%s lost value from underscore column", key)
		}
	}
}

func TestRequestDecimalRetainsScaleAndLargeInteger(t *testing.T) {
	m := bodyMap(httptest.NewRequest("POST", "/", strings.NewReader(`{"price":123.4500,"id":9007199254740993}`)))
	v, ok := m["price"].(json.Number)
	if !ok || v.String() != "123.4500" {
		t.Fatalf("decimal lost scale: %#v", m["price"])
	}
	for _, key := range []string{"price", "id"} {
		kind := "d"
		if key == "id" {
			kind = "l"
		}
		if _, ok := bindVal(kind, m[key]); !ok {
			t.Fatalf("number rejected for %s", key)
		}
	}
	if javaStr(m["id"]) != "9007199254740993" {
		t.Fatal("integer rounded")
	}
}

func TestSQLLiteralPreservesJSONNumber(t *testing.T) {
	if got := javaNumStr(json.Number("3.50")); got != "3.50" {
		t.Fatalf("SQL literal=%q", got)
	}
}
