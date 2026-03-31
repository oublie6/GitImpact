package handler

import (
	"strconv"

	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/service"
	"gitimpact/backend/internal/worker"
	"gitimpact/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	svc    *service.TaskService
	worker *worker.TaskWorker
}

func NewTaskHandler(s *service.TaskService, w *worker.TaskWorker) *TaskHandler {
	return &TaskHandler{svc: s, worker: w}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req model.AnalysisTask
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 400, err.Error())
		return
	}
	req.CreatedBy = c.GetUint("user_id")
	if err := h.svc.CreateTask(&req); err != nil {
		response.Err(c, 400, err.Error())
		return
	}
	h.worker.Enqueue(req.ID)
	response.OK(c, req)
}
func (h *TaskHandler) List(c *gin.Context) {
	list, err := h.svc.ListTasks()
	if err != nil {
		response.Err(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}
func (h *TaskHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := h.svc.GetTask(uint(id))
	if err != nil {
		response.Err(c, 404, err.Error())
		return
	}
	response.OK(c, t)
}
func (h *TaskHandler) Logs(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	logs, err := h.svc.ListLogs(uint(id))
	if err != nil {
		response.Err(c, 500, err.Error())
		return
	}
	response.OK(c, logs)
}
func (h *TaskHandler) Report(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	rp, err := h.svc.GetReport(uint(id))
	if err != nil {
		response.Err(c, 404, err.Error())
		return
	}
	response.OK(c, rp)
}

func (h *TaskHandler) DownloadReport(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	rp, err := h.svc.GetReport(uint(id))
	if err != nil {
		response.Err(c, 404, err.Error())
		return
	}
	format := c.DefaultQuery("format", "md")
	if format == "json" {
		c.Header("Content-Disposition", "attachment; filename=report.json")
		c.Data(200, "application/json", []byte(rp.StructuredReport))
		return
	}
	c.Header("Content-Disposition", "attachment; filename=report.md")
	c.Data(200, "text/markdown; charset=utf-8", []byte(rp.MarkdownReport))
}
