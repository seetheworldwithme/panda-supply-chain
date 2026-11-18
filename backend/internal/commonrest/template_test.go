package commonrest

import (
	"encoding/json"
	"testing"
)

func TestCommonEntityBindingAndShape(t *testing.T) {
	e := ents["ChainCommonAddress"]
	m := map[string]any{"id": "5", "addressName": "测试"}
	if !e.validateBody(m) {
		t.Fatal("numeric strings must bind")
	}
	if e.validateBody(map[string]any{"id": "abc"}) {
		t.Fatal("invalid Integer must reject")
	}
	row := e.shapeOne(map[string]any{"id": int64(5), "address_name": "测试"})
	if row["addressName"] != "测试" || row["page"] != nil {
		t.Fatal(row)
	}
	if _, ok := row["ids"]; !ok {
		t.Fatal("entity carrier key omitted")
	}
}
func TestDecimalScale(t *testing.T) {
	e := ents["ChainTlzfOrder"]
	row := e.shapeOne(map[string]any{"amt_tran": "1.2300"})
	if row["amtTran"] != json.Number("1.2300") {
		t.Fatal(row["amtTran"])
	}
}

func TestJSONNumberBindingPreservesScale(t *testing.T) {
	v, ok := bindVal("d", json.Number("1.2300"))
	if !ok || v != "1.2300" {
		t.Fatalf("%v %v", v, ok)
	}
}

func TestCommonFieldAnnotationsDoNotLeak(t *testing.T) {
	if entChainMsg.typeOf("updateTime") != "T" {
		t.Fatal("updateTime must not inherit createTime JsonFormat")
	}
	for _, f := range entChainMsg.excel {
		if f.col == "msg_content" && f.header != "消息内容" {
			t.Fatal(f)
		}
	}
	if entChainCommonAddress.fill {
		t.Fatal("address has no fill annotations")
	}
	if !entChainTlzfOrder.fill {
		t.Fatal("Allinpay order must fill timestamps")
	}
}
