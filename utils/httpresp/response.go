package httpresp

import (
	"encoding/json"
	"net/http"
	"tk-common/utils/codes"
)

// Envelope 统一响应结构。
type Envelope struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Write 按传入 HTTP 状态码输出统一响应结构。
func Write(w http.ResponseWriter, statusCode int, code int, msg string, data interface{}) {
	// 明确 JSON 响应头，避免浏览器/代理误判。
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(Envelope{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// OK 输出成功响应（HTTP 200）。
func OK(w http.ResponseWriter, data interface{}) {
	Write(w, http.StatusOK, codes.OK, "ok", data)
}

// Fail 输出失败响应（支持自定义 HTTP 状态码）。
func Fail(w http.ResponseWriter, statusCode int, code int, msg string) {
	Write(w, statusCode, code, msg, nil)
}

// BizFail 输出业务失败响应（HTTP 固定 200）。
func BizFail(w http.ResponseWriter, code int, msg string) {
	Write(w, http.StatusOK, code, msg, nil)
}
