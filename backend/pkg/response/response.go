package response

import "github.com/gin-gonic/gin"

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{})        { c.JSON(200, Body{Code: 0, Message: "ok", Data: data}) }
func Err(c *gin.Context, status int, msg string) { c.JSON(status, Body{Code: status, Message: msg}) }
