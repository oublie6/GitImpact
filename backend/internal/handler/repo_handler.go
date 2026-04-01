// repo_handler.go 提供仓库增改查与手动拉取接口。
package handler

import (
	"strconv"

	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/service"
	"gitimpact/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// RepoHandler 负责仓库资源的 HTTP 适配。
type RepoHandler struct{ svc *service.RepositoryService }

// NewRepoHandler 创建仓库处理器。
func NewRepoHandler(s *service.RepositoryService) *RepoHandler { return &RepoHandler{svc: s} }

// Create 创建仓库记录。
func (h *RepoHandler) Create(c *gin.Context) {
	var req model.Repository
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 400, err.Error())
		return
	}
	if err := h.svc.Create(&req); err != nil {
		response.Err(c, 400, err.Error())
		return
	}
	response.OK(c, req)
}

// Update 更新指定仓库记录。
func (h *RepoHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.Repository
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 400, err.Error())
		return
	}
	req.ID = uint(id)
	if err := h.svc.Update(&req); err != nil {
		response.Err(c, 400, err.Error())
		return
	}
	response.OK(c, req)
}

// List 返回仓库列表。
func (h *RepoHandler) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		response.Err(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

// Detail 返回单个仓库详情。
func (h *RepoHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	r, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Err(c, 404, err.Error())
		return
	}
	response.OK(c, r)
}

// Fetch 触发仓库缓存同步。
// 这里不会创建任务，只负责 git clone/fetch。
func (h *RepoHandler) Fetch(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	r, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Err(c, 404, err.Error())
		return
	}
	if err := h.svc.Fetch(r); err != nil {
		response.Err(c, 500, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "fetched"})
}
