package util

import "github.com/gin-gonic/gin"

// Response 统一 API 响应格式
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Data:    data,
		Message: "ok",
	})
}

// Error 返回错误响应
func Error(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Data:    nil,
		Message: message,
	})
}
