// Package platform 统一响应体与通用工具。JSON 形状对照 Java ResultTool/JsonResult:
// 成功 {"success":true,"code":200,"msg":"成功","data":...}
package platform

import (
	"encoding/json"
	"net/http"
)

type JsonResult struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func Success(w http.ResponseWriter, data interface{}) {
	write(w, JsonResult{Success: true, Code: 200, Msg: "成功", Data: data})
}

// Fail 业务失败。code 语义照 Java ResultCode: 999 未知错误 2003 账号密码错误 3002 验证码过期等
func Fail(w http.ResponseWriter, code int, msg string) {
	write(w, JsonResult{Success: false, Code: code, Msg: msg})
}

// FailMsg 对照 Java ResultTool.fail(T data): code 999 固定,消息放 data(前端按 code==999 读 data 弹窗)
func FailMsg(w http.ResponseWriter, msg string) {
	write(w, JsonResult{Success: false, Code: 999, Msg: "失败", Data: msg})
}

// FailData 业务失败且带 data(照 Java ResultTool.fail(code, data) 两参形态,前端有读取 data 的场景)
func FailData(w http.ResponseWriter, code int, msg string, data interface{}) {
	write(w, JsonResult{Success: false, Code: code, Msg: msg, Data: data})
}

// NotAuth 403 网关级未授权(前端 http.js 走 HTTP 错误分支弹 toast)。
// Java 版该响应 HTTP 状态仍是 200、靠 body code 判断,这里保持一致
func NotAuth(w http.ResponseWriter) {
	writeRawCode(w, JsonResult{Success: false, Code: 403, Msg: "请求未授权，访问失败"})
}

// Unknown 999 与 Java 一致:HTTP 200 + code 999
func Unknown(w http.ResponseWriter) {
	writeRawCode(w, JsonResult{Success: false, Code: 999, Msg: "失败", Data: "系统发生未知错误，请稍后重试。"})
}

func write(w http.ResponseWriter, r JsonResult) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(r)
}

func writeRawCode(w http.ResponseWriter, r JsonResult) {
	// 与 Java 一致:HTTP 状态 200,错误码在 body
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(r)
}
