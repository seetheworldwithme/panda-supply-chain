package logistics

import (
	"strings"
	"testing"
)

func TestWlSpecOrder(t *testing.T) {
	sp := newSpec(wlFields)
	for _, c := range []struct{ cell, wantCol string }{
		{"企业id", "company_id"},
		{"订单id", "id"}, // cia 序 id 在前(劫持,与 8765 一致)
		{"票号", "bill_number"},
		{"时间间隔（分钟）", "time_interval_minute"},
	} {
		if f, ok := sp.matchHeader(c.cell); !ok || f.col != c.wantCol {
			t.Errorf("%s → %v(%v) 期望 %v", c.cell, f.col, ok, c.wantCol)
		}
	}
	lo := newSpec(loFields)
	if f, ok := lo.matchHeader("企业id"); !ok || f.col != "company_id" {
		t.Errorf("lo 企业id → %v 期望 company_id", f.col)
	}
	tw := newSpec(twFields)
	if f, ok := tw.matchHeader("运单编号"); !ok || f.col != "bill_number" {
		t.Errorf("tw 运单编号 → %v 期望 bill_number", f.col)
	}
}

func TestWlOrderDump(t *testing.T) {
	sp := newSpec(wlFields)
	t.Log("GO|", strings.Join(sp.order, "|"))
}
