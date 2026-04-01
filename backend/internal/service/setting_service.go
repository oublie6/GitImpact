// setting_service.go 实现系统设置的简单键值对读写。
package service

import "gitimpact/backend/internal/repository"

// SettingService 封装系统设置仓储。
type SettingService struct{ repo *repository.SettingRepository }

// NewSettingService 创建设置服务。
func NewSettingService(repo *repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

// List 返回当前全部设置项。
func (s *SettingService) List() (interface{}, error) { return s.repo.List() }

// Save 按 key 保存设置值。
func (s *SettingService) Save(key, value string) error { return s.repo.Save(key, value) }
