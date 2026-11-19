package commonrest

import (
	"encoding/json"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

func legacyOCRResult(response map[string]any, typ int) platform.JsonResult {
	if typ == 2 && strings.TrimSpace(str(response, "Address")) == "" {
		return platform.JsonResult{Success: false, Code: 701, Msg: "识别失败", Data: "FailedOperation.OcrFailed"}
	}
	var data any = response
	if typ == 2 {
		data = ocrBeanJSON(response)
	}
	if typ == 3 {
		data = ocrBeanJSON(response["TextDetections"])
	}
	encoded, _ := json.Marshal(data)
	return platform.JsonResult{Success: true, Code: 200, Msg: "成功", Data: string(encoded)} // 原 Java data 是 JSON 字符串。
}

func ocrBeanJSON(value any) any {
	switch v := value.(type) {
	case map[string]any:
		out := map[string]any{}
		for key, value := range v {
			if value != nil && key != "" {
				out[strings.ToLower(key[:1])+key[1:]] = ocrBeanJSON(value)
			}
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			out[i] = ocrBeanJSON(item)
		}
		return out
	default:
		return value
	}
}
