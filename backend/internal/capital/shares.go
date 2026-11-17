// capital - shares.go 写链路共享小件:filelist 批更、变更日志(文本细节债务,同四期仅前缀)。
package capital

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/auth"
	"github.com/seetheworldwithme/panda-supply-chain/backend/internal/platform"
)

// updateFilelist 对照 feign updBthById 前置赋值:fileType/ownerType/ownerId/relationId/
// relationType/effective 固定后按 id 批更 chain_common_file(无 fill 实体)
func updateFilelist(d *sql.DB, m map[string]any, id int64, u *auth.CurrentUser) {
	files, ok := m["filelist"].([]any)
	if !ok {
		return
	}
	for _, it := range files {
		f, ok := it.(map[string]any)
		if !ok || f["id"] == nil {
			continue
		}
		_, _ = platform.UpdateByIDNoFill(d, "chain_common_file", f["id"], map[string]any{
			"file_type": "1", "owner_type": "1", "owner_id": strconv.FormatInt(u.UserID, 10),
			"relation_id": strconv.FormatInt(id, 10), "relation_type": 1, "effective": true,
		})
	}
}

// isEmptyStr 数值/布尔字段空串判定(Jackson 把 "" 绑为 null,不拼条件)
func isEmptyStr(v any) bool {
	s, ok := v.(string)
	return ok && strings.TrimSpace(s) == ""
}

// changeLog 对照 EntityUtils.getChangeLog 文本(细节为债务,仅拼固定头部;operation 只入日志表)
func changeLog(before, after map[string]any, keys []string) string {
	if before == nil {
		return "对比结果为：\n"
	}
	return "对比结果为：\n"
}

// echoOf add/updateById 内存实体回显:EntityEcho + 原生 boolean 载体键缺省 false
func echoOf(ent *entShape, m map[string]any, extra map[string]any) map[string]any {
	out := platform.EntityEcho(m, ent.echoKeys(), extra)
	for _, k := range ent.primFalse {
		if out[k] == nil {
			out[k] = false
		}
	}
	return out
}
