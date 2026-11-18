package warehouse

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBiddingFetchParsesRowsAndPreservesDecimal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" {
			t.Error("vendor request contract")
		}
		_, _ = w.Write([]byte(`{"data":[{"preorderBeginTime":"2025-01-01 12:00:00","price":123.4500}]}`))
	}))
	defer server.Close()
	rows, err := fetchBiddingRows(t.Context(), server.Client(), server.URL)
	if err != nil || len(rows) != 1 || javaStr(rows[0]["price"]) != "123.4500" {
		t.Fatalf("rows=%#v err=%v", rows, err)
	}
}
func TestBiddingSelectionRequiresStrictlyNewerTimestamp(t *testing.T) {
	cutoff, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 12:00:00")
	rows := []map[string]any{{"preorderBeginTime": "2025-01-01 12:00:00"}, {"preorderBeginTime": "2025-01-02 12:00:00"}}
	selected, err := selectBiddingRows(rows, &cutoff)
	if err != nil || len(selected) != 1 {
		t.Fatalf("rows=%#v err=%v", selected, err)
	}
	if _, err := selectBiddingRows([]map[string]any{{}}, nil); err == nil {
		t.Fatal("missing date must fail like Java parse/get.toString")
	}
}

func TestMembershipMoneySubtractionIsExact(t *testing.T) {
	got, ok := decimalSubtractExact("100.0500", "0.01")
	if !ok || got != "100.0400" {
		t.Fatalf("amount=%q ok=%v", got, ok)
	}
	if _, ok := decimalSubtractExact(nil, "0"); ok {
		t.Fatal("Java BigDecimal null must fail")
	}
}
