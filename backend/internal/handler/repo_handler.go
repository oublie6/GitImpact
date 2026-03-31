package handler

import (
	"strconv"

	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/service"
	"gitimpact/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type RepoHandler struct{ svc *service.RepositoryService }

func NewRepoHandler(s *service.RepositoryService) *RepoHandler { return &RepoHandler{svc: s} }

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
func (h *RepoHandler) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		response.Err(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}
func (h *RepoHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	r, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Err(c, 404, err.Error())
		return
	}
	response.OK(c, r)
}
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
