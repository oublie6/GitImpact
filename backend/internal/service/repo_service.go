package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/repository"
)

type RepositoryService struct{ repo *repository.RepoRepository }

func NewRepositoryService(repo *repository.RepoRepository) *RepositoryService {
	return &RepositoryService{repo: repo}
}
func (s *RepositoryService) Create(r *model.Repository) error           { return s.repo.Create(r) }
func (s *RepositoryService) Update(r *model.Repository) error           { return s.repo.Update(r) }
func (s *RepositoryService) List() ([]model.Repository, error)          { return s.repo.List() }
func (s *RepositoryService) GetByID(id uint) (*model.Repository, error) { return s.repo.GetByID(id) }

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
