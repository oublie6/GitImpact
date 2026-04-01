// setting_handler.go 提供系统设置查询与保存接口。
package handler

import (
	"gitimpact/backend/internal/service"
	"gitimpact/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// SettingHandler 负责 settings 资源的 HTTP 入口。
type SettingHandler struct{ svc *service.SettingService }

// NewSettingHandler 创建设置处理器。
func NewSettingHandler(s *service.SettingService) *SettingHandler { return &SettingHandler{svc: s} }

// List 返回全部设置项。
func (h *SettingHandler) List(c *gin.Context) {
	items, err := h.svc.List()
	if err != nil {
		response.Err(c, 500, err.Error())
		return
	}
	response.OK(c, items)
}

// Save 保存单个 key/value 设置。
func (h *SettingHandler) Save(c *gin.Context) {
	var req struct{ Key, Value string }
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 400, err.Error())
		return
	}
	if err := h.svc.Save(req.Key, req.Value); err != nil {
		response.Err(c, 400, err.Error())
		return
	}
	response.OK(c, req)
}
