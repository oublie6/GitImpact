package service

import (
	"testing"

	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateTask(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&model.AnalysisTask{}, &model.Repository{}, &model.AnalysisReport{})
	svc := NewTaskService(repository.NewRepoRepository(db), repository.NewTaskRepository(db), repository.NewReportRepository(db))
	task := &model.AnalysisTask{TaskName: "t1", Mode: "same_repo_commits", OldRepoID: 1, OldRef: "a", NewRepoID: 1, NewRef: "b"}
	if err := svc.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	if task.Status != model.TaskStatusPending {
		t.Fatal("status should pending")
	}
}
