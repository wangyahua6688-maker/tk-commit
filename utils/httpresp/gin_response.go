package httpresp

import "github.com/gin-gonic/gin"

// GinOK 输出 Gin 成功响应（HTTP 200）。
func GinOK(c *gin.Context, data interface{}) {
	c.JSON(200, Envelope{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

// GinError 输出 Gin 失败响应（HTTP 200，业务码透传）。
func GinError(c *gin.Context, code int, msg string) {
	c.JSON(200, Envelope{
		Code: code,
		Msg:  msg,
	})
}

// GinFailWithStatus 输出 Gin 失败响应（可指定 HTTP 状态码）。
func GinFailWithStatus(c *gin.Context, statusCode int, code int, msg string) {
	c.JSON(statusCode, Envelope{
		Code: code,
		Msg:  msg,
	})
}
