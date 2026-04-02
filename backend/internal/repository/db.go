// Package repository 提供数据库初始化与数据访问封装。
//
// db.go 负责根据配置选择具体数据库驱动，并在启动时校验核心表是否已初始化。
package repository

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	dameng "github.com/godoes/gorm-dameng"
	"gitimpact/backend/internal/config"
	"gitimpact/backend/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDB 创建数据库连接，根据 database.type 选择对应驱动。
func NewDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch strings.ToLower(strings.TrimSpace(cfg.Type)) {
	case "", "mysql":
		if cfg.DSN == "" {
			return nil, fmt.Errorf("database.dsn is required for mysql")
		}
		dialector = mysql.Open(cfg.DSN)
	case "dameng":
		dsn, err := resolveDamengDSN(cfg)
		if err != nil {
			return nil, err
		}
		dialector = dameng.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database.type: %s", cfg.Type)
	}
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if cfg.AutoMigrate {
		log.Println("[db] database.auto_migrate=true, applying GORM AutoMigrate (not recommended for dameng)")
		if err := db.AutoMigrate(
			&model.User{}, &model.Repository{}, &model.AnalysisTask{}, &model.AnalysisReport{},
			&model.TaskLog{}, &model.SystemSetting{}, &model.TaskArtifact{},
		); err != nil {
			return nil, fmt.Errorf("auto migrate failed: %w", err)
		}
	} else {
		log.Println("[db] database.auto_migrate=false, skip GORM AutoMigrate and require manual SQL initialization")
	}

	if err := verifySchemaInitialized(db, cfg.Type); err != nil {
		return nil, err
	}
	return db, nil
}

// resolveDamengDSN 在达梦场景下统一解析最终连接串。
// 当 database.dsn 已配置时直接使用；否则根据 dameng.* 字段拼接。
func resolveDamengDSN(cfg config.DatabaseConfig) (string, error) {
	if strings.TrimSpace(cfg.DSN) != "" {
		return cfg.DSN, nil
	}
	host := strings.TrimSpace(cfg.Dameng["host"])
	portText := strings.TrimSpace(cfg.Dameng["port"])
	user := strings.TrimSpace(cfg.Dameng["user"])
	password := cfg.Dameng["password"]
	schema := strings.TrimSpace(cfg.Dameng["dbname"])
	if host == "" || portText == "" || user == "" || password == "" {
		return "", fmt.Errorf("dameng config requires host/port/user/password when database.dsn is empty")
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		return "", fmt.Errorf("invalid dameng port: %w", err)
	}
	options := map[string]string{}
	if schema != "" {
		options["schema"] = schema
	}
	return dameng.BuildUrl(user, password, host, port, options), nil
}

func verifySchemaInitialized(db *gorm.DB, dbType string) error {
	missingTables := make([]string, 0, 4)
	coreTables := []string{"users", "repositories", "analysis_tasks", "analysis_reports"}
	for _, table := range coreTables {
		if !db.Migrator().HasTable(table) {
			missingTables = append(missingTables, table)
		}
	}
	if len(missingTables) == 0 {
		return nil
	}

	sqlPath := "sql/mysql/init.sql"
	if strings.EqualFold(strings.TrimSpace(dbType), "dameng") {
		sqlPath = "sql/dameng/init.sql"
	}
	return fmt.Errorf(
		"database schema is not initialized, missing core table(s): %s; please run %s before starting GitImpact",
		strings.Join(missingTables, ", "),
		sqlPath,
	)
}
