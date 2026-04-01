// Package config 定义 GitImpact 的配置结构与加载逻辑。
//
// 该文件对应 backend/config.yaml 的 YAML 结构，是启动流程、数据库连接、
// OpenCode 调用方式以及运行时目录规划的统一来源。
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// AppConfig 是应用总配置结构。
type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Auth     AuthConfig     `yaml:"auth"`
	Database DatabaseConfig `yaml:"database"`
	OpenCode OpenCodeConfig `yaml:"opencode"`
	Workdir  WorkdirConfig  `yaml:"workdir"`
}

// ServerConfig 描述 HTTP 服务监听配置。
type ServerConfig struct {
	Port string `yaml:"port"`
}

// AuthConfig 描述认证模式、JWT 参数与配置用户列表。
type AuthConfig struct {
	Mode               string       `yaml:"mode"`
	JWTSecret          string       `yaml:"jwt_secret"`
	TokenExpireMinutes int          `yaml:"token_expire_minutes"`
	EnableRegister     bool         `yaml:"enable_register"`
	InitAdminEnabled   bool         `yaml:"init_admin_enabled"`
	ConfigUsers        []ConfigUser `yaml:"config_users"`
}

// ConfigUser 表示写在配置文件中的静态用户。
// 这类用户不会落库，适合开发环境或应急管理账号。
type ConfigUser struct {
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	DisplayName string `yaml:"display_name"`
	Email       string `yaml:"email"`
	Role        string `yaml:"role"`
	Status      string `yaml:"status"`
}

// DatabaseConfig 描述数据库类型、直连 DSN 与分数据库类型的补充字段。
type DatabaseConfig struct {
	Type   string            `yaml:"type"`
	DSN    string            `yaml:"dsn"`
	MySQL  map[string]string `yaml:"mysql"`
	Dameng map[string]string `yaml:"dameng"`
}

// OpenCodeConfig 描述 OpenCode CLI 的调用方式。
type OpenCodeConfig struct {
	BinaryPath     string `yaml:"binary_path"`
	AttachURL      string `yaml:"attach_url"`
	DefaultModel   string `yaml:"default_model"`
	DefaultAgent   string `yaml:"default_agent"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

// WorkdirConfig 定义运行期目录分工。
// Root 为根目录，其余目录分别存放仓库缓存、任务材料和最终报告。
type WorkdirConfig struct {
	Root      string `yaml:"root"`
	RepoCache string `yaml:"repo_cache"`
	Artifacts string `yaml:"artifacts"`
	Reports   string `yaml:"reports"`
}

// TokenTTL 将分钟级配置转换为 time.Duration，供 JWT 生成复用。
func (a AppConfig) TokenTTL() time.Duration {
	return time.Duration(a.Auth.TokenExpireMinutes) * time.Minute
}

// Load 从 yaml 文件加载配置。
func Load(path string) (*AppConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}
	cfg := &AppConfig{}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	// 这里补齐最基本的默认值，避免配置文件只写关键项时服务无法启动。
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Auth.TokenExpireMinutes == 0 {
		cfg.Auth.TokenExpireMinutes = 120
	}
	if cfg.OpenCode.TimeoutSeconds == 0 {
		cfg.OpenCode.TimeoutSeconds = 600
	}
	return cfg, nil
}
