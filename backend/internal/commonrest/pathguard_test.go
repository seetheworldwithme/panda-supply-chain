package commonrest

import "testing"

// 文件行的物理路径会被预览/下载接口直接当磁盘路径:客户端写入必须被约束在上传根目录内。
func TestConfineFilePath(t *testing.T) {
	vals := map[string]any{"file_path": "/etc", "file_name": "passwd"}
	if msg := confineFilePath(vals, "chain_common_file", true); msg == "" {
		t.Fatal("上传根目录之外的 file_path 未被拒绝")
	}
	delete(vals, "file_path")
	vals["file_path"] = externalRoot() + "/2026/a.jpg"
	if msg := confineFilePath(vals, "chain_common_file", true); msg != "" {
		t.Fatal(msg)
	}
	vals["file_name"] = "../keys/jwt_private.pem"
	if msg := confineFilePath(vals, "chain_common_img", true); msg == "" {
		t.Fatal("带路径分隔符的 file_name 未被拒绝")
	}
	vals = map[string]any{"file_path": "/etc", "file_name": "x.jpg"}
	if msg := confineFilePath(vals, "chain_common_file", false); msg != "" || vals["file_path"] != nil {
		t.Fatalf("更新时应忽略 file_path 而不是报错: msg=%q vals=%v", msg, vals)
	}
	if msg := confineFilePath(map[string]any{"file_path": "/etc"}, "chain_bill", true); msg != "" {
		t.Fatal("非文件表不受影响")
	}
	if extRowUsable(map[string]any{"file_path": "", "file_name": ".env"}, "") {
		t.Fatal("目录为空的文件行可被用来读进程工作目录下的文件")
	}
	if !extRowUsable(map[string]any{"file_path": "/x", "file_name": "a.jpg"}, "") {
		t.Fatal("正常的文件行被误拒")
	}
}
