// repo_service.go 实现仓库元数据管理与 Git 仓库同步。
//
// 这里既负责数据库中的仓库记录，也负责 git clone/fetch 的执行入口。
package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/repository"
)

// RepositoryService 封装仓库 CRUD 与同步操作。
type RepositoryService struct{ repo *repository.RepoRepository }

// NewRepositoryService 创建仓库服务。
func NewRepositoryService(repo *repository.RepoRepository) *RepositoryService {
	return &RepositoryService{repo: repo}
}

// Create 新增仓库记录。
func (s *RepositoryService) Create(r *model.Repository) error { return s.repo.Create(r) }

// Update 更新仓库记录。
func (s *RepositoryService) Update(r *model.Repository) error { return s.repo.Update(r) }

// List 返回全部仓库记录。
func (s *RepositoryService) List() ([]model.Repository, error) { return s.repo.List() }

// GetByID 按 ID 读取仓库。
func (s *RepositoryService) GetByID(id uint) (*model.Repository, error) { return s.repo.GetByID(id) }

// Fetch 确保仓库本地缓存存在，并执行 git fetch --all --prune 更新远端引用。
// 如果缓存目录尚未初始化为 Git 仓库，会先执行 git clone。
func (s *RepositoryService) Fetch(repo *model.Repository) error {
	dir := repo.LocalCacheDir
	if _, err := os.Stat(filepath.Join(dir, ".git")); err != nil {
		cmd := exec.Command("git", "clone", repo.RepoURL, dir)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("clone failed: %s %w", string(out), err)
		}
	}
	cmd := exec.Command("git", "-C", dir, "fetch", "--all", "--prune")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("fetch failed: %s %w", string(out), err)
	}
	return nil
}
