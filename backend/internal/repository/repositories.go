package repository

import (
	"gitimpact/backend/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct{ db *gorm.DB }
type RepoRepository struct{ db *gorm.DB }
type TaskRepository struct{ db *gorm.DB }
type ReportRepository struct{ db *gorm.DB }
type SettingRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository       { return &UserRepository{db: db} }
func NewRepoRepository(db *gorm.DB) *RepoRepository       { return &RepoRepository{db: db} }
func NewTaskRepository(db *gorm.DB) *TaskRepository       { return &TaskRepository{db: db} }
func NewReportRepository(db *gorm.DB) *ReportRepository   { return &ReportRepository{db: db} }
func NewSettingRepository(db *gorm.DB) *SettingRepository { return &SettingRepository{db: db} }

func (r *UserRepository) Create(u *model.User) error { return r.db.Create(u).Error }
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *UserRepository) Update(u *model.User) error { return r.db.Save(u).Error }

func (r *RepoRepository) Create(repo *model.Repository) error { return r.db.Create(repo).Error }
func (r *RepoRepository) Update(repo *model.Repository) error { return r.db.Save(repo).Error }
func (r *RepoRepository) List() ([]model.Repository, error) {
	var list []model.Repository
	return list, r.db.Order("id desc").Find(&list).Error
}
func (r *RepoRepository) GetByID(id uint) (*model.Repository, error) {
	var repo model.Repository
	if err := r.db.First(&repo, id).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *TaskRepository) Create(t *model.AnalysisTask) error { return r.db.Create(t).Error }
func (r *TaskRepository) Update(t *model.AnalysisTask) error { return r.db.Save(t).Error }
func (r *TaskRepository) List() ([]model.AnalysisTask, error) {
	var list []model.AnalysisTask
	return list, r.db.Order("id desc").Find(&list).Error
}
func (r *TaskRepository) GetByID(id uint) (*model.AnalysisTask, error) {
	var t model.AnalysisTask
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
func (r *TaskRepository) AddLog(log *model.TaskLog) error { return r.db.Create(log).Error }
func (r *TaskRepository) ListLogs(taskID uint) ([]model.TaskLog, error) {
	var logs []model.TaskLog
	return logs, r.db.Where("task_id = ?", taskID).Order("id asc").Find(&logs).Error
}
func (r *TaskRepository) AddArtifact(a *model.TaskArtifact) error { return r.db.Create(a).Error }
func (r *TaskRepository) DB() *gorm.DB                            { return r.db }

func (r *ReportRepository) Upsert(taskID uint, report *model.AnalysisReport) error {
	var exist model.AnalysisReport
	err := r.db.Where("task_id = ?", taskID).First(&exist).Error
	if err == nil {
		report.ID = exist.ID
		return r.db.Save(report).Error
	}
	return r.db.Create(report).Error
}
func (r *ReportRepository) GetByTaskID(taskID uint) (*model.AnalysisReport, error) {
	var rp model.AnalysisReport
	if err := r.db.Where("task_id = ?", taskID).First(&rp).Error; err != nil {
		return nil, err
	}
	return &rp, nil
}

func (r *SettingRepository) List() ([]model.SystemSetting, error) {
	var items []model.SystemSetting
	return items, r.db.Order("id asc").Find(&items).Error
}
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
