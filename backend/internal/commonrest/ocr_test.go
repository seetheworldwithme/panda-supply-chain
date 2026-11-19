package commonrest

import (
	"strings"
	"testing"
)

func TestOCRPreservesJavaStringAndFieldShapes(t *testing.T) {
	response := map[string]any{"Name": "fixture", "Address": "fixture-address", "RequestId": "fixture-request", "TextDetections": []any{map[string]any{"DetectedText": "fixture-text", "Confidence": 99}}}
	for typ, wanted := range map[int]string{1: `"Name":"fixture"`, 2: `"address":"fixture-address"`, 3: `"detectedText":"fixture-text"`} {
		result := legacyOCRResult(response, typ)
		data, ok := result.Data.(string)
		if !ok || !strings.Contains(data, wanted) {
			t.Errorf("type=%d 返回形状不符: %#v", typ, result)
		}
	}
	if result := legacyOCRResult(map[string]any{"Address": ""}, 2); result.Code != 701 || result.Data != "FailedOperation.OcrFailed" {
		t.Fatal(result)
	}
}
