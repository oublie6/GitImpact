package handler

import (
	"gitimpact/backend/internal/service"
	"gitimpact/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct{ svc *service.SettingService }

func NewSettingHandler(s *service.SettingService) *SettingHandler { return &SettingHandler{svc: s} }
func (h *SettingHandler) List(c *gin.Context) {
	items, err := h.svc.List()
	if err != nil {
		response.Err(c, 500, err.Error())
		return
	}
	response.OK(c, items)
}
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
