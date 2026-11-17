package logistics

import "testing"

func TestJavaMapOrder(t *testing.T) {
	// 8765 实测的三个判定:企业id→companyId、装货地经纬度→装货地、创建时间→时间
	lo := newSpec(loFields)
	if f, ok := lo.matchHeader("企业id"); !ok || f.col != "company_id" {
		t.Fatalf("企业id → %v %v", f.col, ok)
	}
	if f, ok := lo.matchHeader("装货地经纬度"); !ok || f.col != "loading_address" {
		t.Fatalf("装货地经纬度 → %v %v", f.col, ok)
	}
	if f, ok := lo.matchHeader("创建时间"); !ok || f.col != "trans_time" {
		t.Fatalf("创建时间 → %v %v", f.col, ok)
	}
	if f, ok := lo.matchHeader("卸货地经纬度"); !ok || f.col != "unloading_point" {
		t.Fatalf("卸货地经纬度 → %v %v", f.col, ok)
	}
	if f, ok := lo.matchHeader("运单编号"); !ok || f.col != "bill_number" {
		t.Fatalf("运单编号 → %v %v", f.col, ok)
	}
	t.Log("matchHeader 模拟与 8765 实测一致")
}
