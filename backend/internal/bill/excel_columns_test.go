package bill

import "testing"

func TestExcelCellReferencesUseZeroBasedColumns(t *testing.T) {
	for ref, want := range map[string]int{"A1": 0, "D5": 3, "AA2": 26} {
		got, err := colIndex(ref)
		if err != nil || got != want {
			t.Errorf("%s 的列号=%d，应为 %d", ref, got, want)
		}
	}
}
