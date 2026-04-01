// repositories.go 定义各业务实体的 Repository 封装。
//
// 这些仓储是 service 层唯一直接依赖的持久化入口，用来隔离 GORM 细节。
package repository

import (
	"gitimpact/backend/internal/model"

	"gorm.io/gorm"
)

// UserRepository 负责用户数据的增删改查。
type UserRepository struct{ db *gorm.DB }

// RepoRepository 负责代码仓库元数据的持久化。
type RepoRepository struct{ db *gorm.DB }

// TaskRepository 负责分析任务、任务日志与任务产物记录。
type TaskRepository struct{ db *gorm.DB }

// ReportRepository 负责分析报告读写。
type ReportRepository struct{ db *gorm.DB }

// SettingRepository 负责系统设置键值对存储。
type SettingRepository struct{ db *gorm.DB }

// NewUserRepository 创建用户仓储。
func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

// NewRepoRepository 创建仓库仓储。
func NewRepoRepository(db *gorm.DB) *RepoRepository { return &RepoRepository{db: db} }

// NewTaskRepository 创建任务仓储。
func NewTaskRepository(db *gorm.DB) *TaskRepository { return &TaskRepository{db: db} }

// NewReportRepository 创建报告仓储。
func NewReportRepository(db *gorm.DB) *ReportRepository { return &ReportRepository{db: db} }

// NewSettingRepository 创建系统设置仓储。
func NewSettingRepository(db *gorm.DB) *SettingRepository { return &SettingRepository{db: db} }

// Create 新增数据库用户。
func (r *UserRepository) Create(u *model.User) error { return r.db.Create(u).Error }

// GetByUsername 按用户名读取用户。
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Update 保存用户最新状态，例如最近登录时间。
func (r *UserRepository) Update(u *model.User) error { return r.db.Save(u).Error }

// Create 新增仓库记录。
func (r *RepoRepository) Create(repo *model.Repository) error { return r.db.Create(repo).Error }

// Update 更新仓库记录。
func (r *RepoRepository) Update(repo *model.Repository) error { return r.db.Save(repo).Error }

// List 返回仓库列表，按 ID 倒序展示最新记录。
func (r *RepoRepository) List() ([]model.Repository, error) {
	var list []model.Repository
	return list, r.db.Order("id desc").Find(&list).Error
}

// GetByID 按主键读取仓库详情。
func (r *RepoRepository) GetByID(id uint) (*model.Repository, error) {
	var repo model.Repository
	if err := r.db.First(&repo, id).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

// Create 新增任务记录。
func (r *TaskRepository) Create(t *model.AnalysisTask) error { return r.db.Create(t).Error }

// Update 保存任务状态变化。
func (r *TaskRepository) Update(t *model.AnalysisTask) error { return r.db.Save(t).Error }

// List 返回任务列表，按 ID 倒序，便于前端优先看到最新任务。
func (r *TaskRepository) List() ([]model.AnalysisTask, error) {
	var list []model.AnalysisTask
	return list, r.db.Order("id desc").Find(&list).Error
}

// GetByID 按主键读取任务详情。
func (r *TaskRepository) GetByID(id uint) (*model.AnalysisTask, error) {
	var t model.AnalysisTask
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// AddLog 为任务追加一条执行日志。
func (r *TaskRepository) AddLog(log *model.TaskLog) error { return r.db.Create(log).Error }

// ListLogs 读取单个任务的历史日志。
func (r *TaskRepository) ListLogs(taskID uint) ([]model.TaskLog, error) {
	var logs []model.TaskLog
	return logs, r.db.Where("task_id = ?", taskID).Order("id asc").Find(&logs).Error
}

// AddArtifact 记录任务产生的中间材料文件路径。
func (r *TaskRepository) AddArtifact(a *model.TaskArtifact) error { return r.db.Create(a).Error }

// DB 暴露底层 GORM 连接，留给需要事务或复杂查询的调用方。
func (r *TaskRepository) DB() *gorm.DB { return r.db }

// Upsert 以 task_id 为唯一键保存报告，保证同一任务只有一份最新报告。
func (r *ReportRepository) Upsert(taskID uint, report *model.AnalysisReport) error {
	var exist model.AnalysisReport
	err := r.db.Where("task_id = ?", taskID).First(&exist).Error
	if err == nil {
		report.ID = exist.ID
		return r.db.Save(report).Error
	}
	return r.db.Create(report).Error
}

// GetByTaskID 按任务 ID 读取报告。
func (r *ReportRepository) GetByTaskID(taskID uint) (*model.AnalysisReport, error) {
	var rp model.AnalysisReport
	if err := r.db.Where("task_id = ?", taskID).First(&rp).Error; err != nil {
		return nil, err
	}
	return &rp, nil
}

// List 返回全部系统设置。
func (r *SettingRepository) List() ([]model.SystemSetting, error) {
	var items []model.SystemSetting
	return items, r.db.Order("id asc").Find(&items).Error
}

// Save 以 key 为唯一键执行保存或覆盖。
func (r *SettingRepository) Save(key, value string) error {
	var s model.SystemSetting
	err := r.db.Where("`key` = ?", key).First(&s).Error
	if err == nil {
		s.Value = value
		return r.db.Save(&s).Error
	}
	s.Key = key
	s.Value = value
	return r.db.Create(&s).Error
}
