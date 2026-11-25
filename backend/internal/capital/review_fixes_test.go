package capital

import "testing"

// 月末起点:AddDate 会把"1月31日+1个月"进位成 3 月,原循环永远等不到 2 月直到 OOM。
func TestCountSeriesMonthEndStart(t *testing.T) {
	series, _, ok := countSeries(map[string]any{"startTime": "2026-01-31", "endTime": "2026-02-10"}, false)
	if !ok || len(series) != 2 || series[0] != "2026-01" || series[1] != "2026-02" {
		t.Fatalf("月份序列错误: %v ok=%v", series, ok)
	}
	if series, _, _ := countSeries(map[string]any{"startTime": "2026-03-01", "endTime": "2026-01-01"}, false); len(series) != 0 {
		t.Fatalf("起晚于止应为空序列: %v", series)
	}
}
