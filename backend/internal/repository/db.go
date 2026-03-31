package repository

import (
	"fmt"

	"gitimpact/backend/internal/config"
	"gitimpact/backend/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDB 创建数据库连接。目前 MySQL 与 达梦都通过兼容 DSN 使用 gorm mysql driver。
func NewDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("database.dsn is required")
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&model.User{}, &model.Repository{}, &model.AnalysisTask{}, &model.AnalysisReport{},
		&model.TaskLog{}, &model.SystemSetting{}, &model.TaskArtifact{},
	); err != nil {
		return nil, err
	}
	return db, nil
}
