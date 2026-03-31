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

func NewTaskService(repoRepo *repository.RepoRepository, taskRepo *repository.TaskRepository, reportRepo *repository.ReportRepository) *TaskService {
	return &TaskService{repoRepo: repoRepo, taskRepo: taskRepo, reportRepo: reportRepo}
}

func (s *TaskService) CreateTask(task *model.AnalysisTask) error {
	task.Status = model.TaskStatusPending
	return s.taskRepo.Create(task)
}
func (s *TaskService) ListTasks() ([]model.AnalysisTask, error) { return s.taskRepo.List() }
func (s *TaskService) GetTask(id uint) (*model.AnalysisTask, error) {
	return s.taskRepo.GetByID(id)
}
func (s *TaskService) ListLogs(id uint) ([]model.TaskLog, error) { return s.taskRepo.ListLogs(id) }
func (s *TaskService) GetReport(id uint) (*model.AnalysisReport, error) {
	return s.reportRepo.GetByTaskID(id)
}
func (s *TaskService) TaskRepo() *repository.TaskRepository     { return s.taskRepo }
func (s *TaskService) RepoRepo() *repository.RepoRepository     { return s.repoRepo }
func (s *TaskService) ReportRepo() *repository.ReportRepository { return s.reportRepo }
