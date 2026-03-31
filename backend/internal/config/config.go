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

type ServerConfig struct {
	Port string `yaml:"port"`
}

type AuthConfig struct {
	Mode               string       `yaml:"mode"`
	JWTSecret          string       `yaml:"jwt_secret"`
	TokenExpireMinutes int          `yaml:"token_expire_minutes"`
	EnableRegister     bool         `yaml:"enable_register"`
	InitAdminEnabled   bool         `yaml:"init_admin_enabled"`
	ConfigUsers        []ConfigUser `yaml:"config_users"`
}

type ConfigUser struct {
	Username      string `yaml:"username"`
	PasswordHash  string `yaml:"password_hash"`
	PasswordPlain string `yaml:"password_plain"`
	DisplayName   string `yaml:"display_name"`
	Email         string `yaml:"email"`
	Role          string `yaml:"role"`
	Status        string `yaml:"status"`
	AllowDevPlain bool   `yaml:"allow_dev_plain"`
}

type DatabaseConfig struct {
	Type   string            `yaml:"type"`
	DSN    string            `yaml:"dsn"`
	MySQL  map[string]string `yaml:"mysql"`
	Dameng map[string]string `yaml:"dameng"`
}

type OpenCodeConfig struct {
	BinaryPath     string `yaml:"binary_path"`
	AttachURL      string `yaml:"attach_url"`
	DefaultModel   string `yaml:"default_model"`
	DefaultAgent   string `yaml:"default_agent"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

type WorkdirConfig struct {
	Root      string `yaml:"root"`
	RepoCache string `yaml:"repo_cache"`
	Artifacts string `yaml:"artifacts"`
	Reports   string `yaml:"reports"`
}

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
