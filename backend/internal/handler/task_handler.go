// task_handler.go 提供任务创建、查询、日志和报告下载接口。
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

// NewTaskHandler 创建任务处理器。
func NewTaskHandler(s *service.TaskService, w *worker.TaskWorker) *TaskHandler {
	return &TaskHandler{svc: s, worker: w}
}

// Create 创建任务并立即异步入队。
// 请求返回时任务通常仍处于 pending/running，实际执行在后台 goroutine 中完成。
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

// List 返回任务列表。
func (h *TaskHandler) List(c *gin.Context) {
	list, err := h.svc.ListTasks()
	if err != nil {
		response.Err(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

// Detail 返回单个任务详情。
func (h *TaskHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := h.svc.GetTask(uint(id))
	if err != nil {
		response.Err(c, 404, err.Error())
		return
	}
	response.OK(c, t)
}

// Logs 返回任务执行日志。
func (h *TaskHandler) Logs(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	logs, err := h.svc.ListLogs(uint(id))
	if err != nil {
		response.Err(c, 500, err.Error())
		return
	}
	response.OK(c, logs)
}

// Report 返回任务报告详情。
func (h *TaskHandler) Report(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	rp, err := h.svc.GetReport(uint(id))
	if err != nil {
		response.Err(c, 404, err.Error())
		return
	}
	response.OK(c, rp)
}

// DownloadReport 以附件形式下载 Markdown 或 JSON 报告。
// format 默认为 md，传 json 时直接返回 structured_report 的原始字符串。
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
