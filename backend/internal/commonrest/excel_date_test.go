package commonrest

import "testing"

func TestExportLocalDateHasNoTime(t *testing.T) {
	if got := exportCell("2026-11-18 00:00:00", "D"); got != "2026-11-18" {
		t.Fatalf("LocalDate 导出多出了时间: %s", got)
	}
}
