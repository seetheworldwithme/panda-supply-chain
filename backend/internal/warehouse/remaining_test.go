package warehouse

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInformationAndPeripheralRoutesRegistered(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, &Service{})
	for _, path := range []string{"/xmmcCoalInfo/importExcel", "/xmmcCokingInfo/importExcel", "/xmmcRailwayInfo/importExcel", "/xmmcWashingInfo/importExcel", "/xmmcCoalInfo/generateAndDownloadExcel", "/xmmcCokingInfo/generateAndDownloadExcel", "/xmmcRailwayInfo/generateAndDownloadExcel", "/xmmcWashingInfo/generateAndDownloadExcel", "/wechat/transfer", "/busVr/getVr"} {
		method := "POST"
		if path == "/busVr/getVr" || strings.HasSuffix(path, "generateAndDownloadExcel") {
			method = "GET"
		}
		req := httptest.NewRequest(method, "/saas-warehouse"+path, nil)
		_, pattern := mux.Handler(req)
		if pattern == "" {
			t.Errorf("unregistered: %s %s", method, path)
		}
	}
}

func TestEnterpriseLeasingWriteRoutesRegistered(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, &Service{})
	for _, path := range []string{"/busEnterpriseLeasing/add", "/busEnterpriseLeasing/updateById"} {
		req := httptest.NewRequest("POST", "/saas-warehouse"+path, nil)
		if _, pattern := mux.Handler(req); pattern == "" {
			t.Errorf("unregistered: POST %s", path)
		}
	}
}

func TestInformationFiltersUseQueryWrapperRanges(t *testing.T) {
	where, args := infoWhere(entXmmcCoalInfo, map[string]any{"city": "晋中", "capacity": float64(10), "sulfur": float64(1), "sulfurMin": float64(0.5), "factoryTaxPriceMax": float64(500), "searchStr": "ignored", "priceUpdTimeStart": "2025-01-01T00:00:00"})
	want := " WHERE city LIKE ? AND sulfur >= ? AND factory_tax_price <= ? AND price_upd_time >= ?"
	if where != want {
		t.Fatalf("where=%q, want %q", where, want)
	}
	if len(args) != 4 || args[0] != "%晋中%" || args[1] != "0.5" {
		t.Fatalf("args=%#v", args)
	}
}

func TestInformationChangeLogKeepsDecimalScaleAndSkipsEqualValues(t *testing.T) {
	original := map[string]any{"factory_tax_price": "500.00", "city": "晋中"}
	vals := map[string]any{"factory_tax_price": "501.50", "city": "晋中"}
	got := infoChangeLog(entXmmcCoalInfo, original, vals)
	if !strings.Contains(got, "500.00——》501.50") || strings.Contains(got, "晋中") {
		t.Fatalf("log=%q", got)
	}
	if !decimalEqual("500.00", "500") {
		t.Fatal("scale must not affect Java BigDecimal.compareTo")
	}
}

func TestInformationCustomDateRanges(t *testing.T) {
	m := map[string]any{"priceUpdTimeStart": "2025-01-01", "priceUpdTimeEnd": "2025-01-02", "updateTimeStart": "bad", "createTimeEnd": "2025-01-03T01:02:03.123"}
	normalizeInfoRanges(m)
	if m["priceUpdTimeStart"] != "2025-01-01T00:00:00" || m["priceUpdTimeEnd"] != "2025-01-02T23:59:59.999999999" || m["updateTimeStart"] != nil || m["createTimeEnd"] != "2025-01-03T01:02:03.123" {
		t.Fatalf("ranges=%#v", m)
	}
	if !entXmmcCoalInfo.validateBody(m) {
		t.Fatal("Java custom date deserializer values must bind")
	}
}
