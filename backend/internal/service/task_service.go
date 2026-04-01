// task_service.go 负责任务元数据、日志与报告的读取封装。
//
// 任务真正的执行发生在 worker 中，这里的职责是入库和查询。
package service

import (
	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/repository"
)

type TaskService struct {
	repoRepo   *repository.RepoRepository
	taskRepo   *repository.TaskRepository
	reportRepo *repository.ReportRepository
}

// NewTaskService 创建任务服务。
func NewTaskService(repoRepo *repository.RepoRepository, taskRepo *repository.TaskRepository, reportRepo *repository.ReportRepository) *TaskService {
	return &TaskService{repoRepo: repoRepo, taskRepo: taskRepo, reportRepo: reportRepo}
}

// CreateTask 创建任务并设置初始状态为 pending。
func (s *TaskService) CreateTask(task *model.AnalysisTask) error {
	task.Status = model.TaskStatusPending
	return s.taskRepo.Create(task)
}

// ListTasks 返回全部任务。
func (s *TaskService) ListTasks() ([]model.AnalysisTask, error) { return s.taskRepo.List() }

// GetTask 按主键查询任务详情。
func (s *TaskService) GetTask(id uint) (*model.AnalysisTask, error) {
	return s.taskRepo.GetByID(id)
}

// ListLogs 返回任务执行日志。
func (s *TaskService) ListLogs(id uint) ([]model.TaskLog, error) { return s.taskRepo.ListLogs(id) }

// GetReport 返回任务最终报告。
func (s *TaskService) GetReport(id uint) (*model.AnalysisReport, error) {
	return s.reportRepo.GetByTaskID(id)
}

// TaskRepo 暴露任务仓储，方便 handler/worker 复用现有装配结果。
func (s *TaskService) TaskRepo() *repository.TaskRepository { return s.taskRepo }

// RepoRepo 暴露仓库仓储。
func (s *TaskService) RepoRepo() *repository.RepoRepository { return s.repoRepo }

// ReportRepo 暴露报告仓储。
func (s *TaskService) ReportRepo() *repository.ReportRepository { return s.reportRepo }
