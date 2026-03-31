package service

import "gitimpact/backend/internal/repository"

type SettingService struct{ repo *repository.SettingRepository }

func NewSettingService(repo *repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}
func (s *SettingService) List() (interface{}, error)   { return s.repo.List() }
func (s *SettingService) Save(key, value string) error { return s.repo.Save(key, value) }
