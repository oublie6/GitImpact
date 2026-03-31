package model

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

const (
	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
)

const (
	UserSourceConfig = "config"
	UserSourceDB     = "db"
)

const (
	TaskStatusPending = "pending"
	TaskStatusRunning = "running"
	TaskStatusSuccess = "success"
	TaskStatusFailed  = "failed"
)

// User 用户信息。
type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	DisplayName  string     `gorm:"size:128" json:"display_name"`
	Email        string     `gorm:"size:128" json:"email"`
	Role         string     `gorm:"size:32;not null" json:"role"`
	Status       string     `gorm:"size:32;not null" json:"status"`
	Source       string     `gorm:"size:32;not null" json:"source"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Repository 代码仓库配置。
type Repository struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:128;uniqueIndex;not null" json:"name"`
	RepoURL       string    `gorm:"size:500;not null" json:"repo_url"`
	DefaultBranch string    `gorm:"size:128;not null" json:"default_branch"`
	LocalCacheDir string    `gorm:"size:500;not null" json:"local_cache_dir"`
	AuthNote      string    `gorm:"size:500" json:"auth_note"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AnalysisTask 分析任务。
type AnalysisTask struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	TaskName           string    `gorm:"size:200;not null" json:"task_name"`
	Mode               string    `gorm:"size:64;not null" json:"mode"`
	OldRepoID          uint      `json:"old_repo_id"`
	OldRef             string    `gorm:"size:128;not null" json:"old_ref"`
	NewRepoID          uint      `json:"new_repo_id"`
	NewRef             string    `gorm:"size:128;not null" json:"new_ref"`
	GenerateMarkdown   bool      `json:"generate_markdown"`
	GenerateStructured bool      `json:"generate_structured"`
	CustomFocus        string    `gorm:"type:text" json:"custom_focus"`
	Remark             string    `gorm:"size:500" json:"remark"`
	Status             string    `gorm:"size:32;not null" json:"status"`
	ErrorMessage       string    `gorm:"type:text" json:"error_message"`
	CreatedBy          uint      `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// AnalysisReport 任务报告。
type AnalysisReport struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TaskID           uint      `gorm:"index;not null" json:"task_id"`
	MarkdownReport   string    `gorm:"type:longtext" json:"markdown_report"`
	StructuredReport string    `gorm:"type:longtext" json:"structured_report"`
	RawStdout        string    `gorm:"type:longtext" json:"raw_stdout"`
	RawStderr        string    `gorm:"type:longtext" json:"raw_stderr"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TaskLog 任务执行日志。
type TaskLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    uint      `gorm:"index;not null" json:"task_id"`
	Level     string    `gorm:"size:16;not null" json:"level"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskArtifact 任务产物路径。
type TaskArtifact struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TaskID      uint      `gorm:"index;not null" json:"task_id"`
	ArtifactKey string    `gorm:"size:128;not null" json:"artifact_key"`
	FilePath    string    `gorm:"size:500;not null" json:"file_path"`
	CreatedAt   time.Time `json:"created_at"`
}

// SystemSetting 系统设置。
type SystemSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"size:128;uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
