package service

import (
	"testing"

	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRepositoryServiceBasic(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&model.Repository{})
	svc := NewRepositoryService(repository.NewRepoRepository(db))
	r := &model.Repository{Name: "r1", RepoURL: "https://example.com/a.git", DefaultBranch: "main", LocalCacheDir: "/tmp/a"}
	if err := svc.Create(r); err != nil {
		t.Fatal(err)
	}
	list, err := svc.List()
	if err != nil || len(list) != 1 {
		t.Fatal("list failed")
	}
}
