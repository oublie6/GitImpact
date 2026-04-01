// Package response 封装统一的 HTTP 响应体格式。
package response

import "github.com/gin-gonic/gin"

// Body 是全部 API 默认使用的响应结构。
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK 返回业务成功响应。
func OK(c *gin.Context, data interface{}) { c.JSON(200, Body{Code: 0, Message: "ok", Data: data}) }

// Err 返回业务失败响应，HTTP 状态码与 body.code 保持一致。
func Err(c *gin.Context, status int, msg string) { c.JSON(status, Body{Code: status, Message: msg}) }
