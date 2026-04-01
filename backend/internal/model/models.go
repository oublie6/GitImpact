// Package model 定义数据库实体与跨层共享的领域常量。
//
// 这些结构体既用于 GORM 建模，也会直接作为 API 响应数据返回给前端。
package model

import "time"

const (
	// RoleAdmin 表示管理员角色。
	RoleAdmin = "admin"
	// RoleUser 表示普通用户角色。
	RoleUser = "user"
)

const (
	// UserStatusActive 表示账号可正常登录。
	UserStatusActive = "active"
	// UserStatusDisabled 表示账号已停用。
	UserStatusDisabled = "disabled"
)

const (
	// UserSourceConfig 表示用户来自配置文件。
	UserSourceConfig = "config"
	// UserSourceDB 表示用户来自数据库。
	UserSourceDB = "db"
)

const (
	// TaskStatusPending 表示任务已创建、尚未执行。
	TaskStatusPending = "pending"
	// TaskStatusRunning 表示任务已进入 worker 执行阶段。
	TaskStatusRunning = "running"
	// TaskStatusSuccess 表示任务执行完成且报告已保存。
	TaskStatusSuccess = "success"
	// TaskStatusFailed 表示任务执行过程中出现错误。
	TaskStatusFailed = "failed"
)

// User 用户信息。
type User struct {
	// ID 是数据库主键。
	ID uint `gorm:"primaryKey" json:"id"`
	// Username 是登录名，全局唯一。
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	// PasswordHash 仅对数据库用户生效，配置用户不落库。
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	// DisplayName 是界面展示名称。
	DisplayName string `gorm:"size:128" json:"display_name"`
	// Email 主要用于展示和联系，目前未接入邮件流程。
	Email string `gorm:"size:128" json:"email"`
	// Role 决定用户权限级别。
	Role string `gorm:"size:32;not null" json:"role"`
	// Status 标识账号是否可用。
	Status string `gorm:"size:32;not null" json:"status"`
	// Source 记录用户来自配置还是数据库。
	Source string `gorm:"size:32;not null" json:"source"`
	// LastLoginAt 仅对数据库用户在成功登录时回写。
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Repository 代码仓库配置。
type Repository struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// Name 是仓库别名，供界面与任务引用。
	Name string `gorm:"size:128;uniqueIndex;not null" json:"name"`
	// RepoURL 是 git clone/fetch 使用的远端地址。
	RepoURL string `gorm:"size:500;not null" json:"repo_url"`
	// DefaultBranch 记录默认分支，当前主要用于展示和默认值。
	DefaultBranch string `gorm:"size:128;not null" json:"default_branch"`
	// LocalCacheDir 是仓库在运行机上的本地缓存目录。
	LocalCacheDir string `gorm:"size:500;not null" json:"local_cache_dir"`
	// AuthNote 用于记录鉴权或访问方式说明，不参与程序逻辑。
	AuthNote  string    `gorm:"size:500" json:"auth_note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AnalysisTask 分析任务。
type AnalysisTask struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// TaskName 是任务显示名称。
	TaskName string `gorm:"size:200;not null" json:"task_name"`
	// Mode 描述比较模式，例如 same_repo_commits。
	Mode string `gorm:"size:64;not null" json:"mode"`
	// OldRepoID 与 NewRepoID 分别表示比较两端仓库。
	OldRepoID uint   `json:"old_repo_id"`
	OldRef    string `gorm:"size:128;not null" json:"old_ref"`
	NewRepoID uint   `json:"new_repo_id"`
	NewRef    string `gorm:"size:128;not null" json:"new_ref"`
	// GenerateMarkdown / GenerateStructured 控制是否产出两类报告。
	GenerateMarkdown   bool `json:"generate_markdown"`
	GenerateStructured bool `json:"generate_structured"`
	// CustomFocus 会直接拼到分析提示词中，作为关注点补充。
	CustomFocus string `gorm:"type:text" json:"custom_focus"`
	// Remark 仅用于记录业务备注，不参与分析逻辑。
	Remark string `gorm:"size:500" json:"remark"`
	// Status 与 ErrorMessage 用于展示任务状态与失败原因。
	Status       string `gorm:"size:32;not null" json:"status"`
	ErrorMessage string `gorm:"type:text" json:"error_message"`
	// CreatedBy 来源于 JWT 中的 user_id；配置用户目前会写入 0。
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AnalysisReport 任务报告。
type AnalysisReport struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// TaskID 是与任务的一对一关联键。
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// MarkdownReport 是面向人阅读的分析结果。
	MarkdownReport string `gorm:"type:longtext" json:"markdown_report"`
	// StructuredReport 是结构化 JSON 字符串。
	StructuredReport string `gorm:"type:longtext" json:"structured_report"`
	// RawStdout / RawStderr 保存 OpenCode 原始输出，便于排障。
	RawStdout string    `gorm:"type:longtext" json:"raw_stdout"`
	RawStderr string    `gorm:"type:longtext" json:"raw_stderr"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskLog 任务执行日志。
type TaskLog struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// TaskID 指向所属任务。
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// Level 当前只使用 INFO / ERROR。
	Level string `gorm:"size:16;not null" json:"level"`
	// Message 保存执行过程中的关键信息。
	Message   string    `gorm:"type:text;not null" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskArtifact 任务产物路径。
type TaskArtifact struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	TaskID uint `gorm:"index;not null" json:"task_id"`
	// ArtifactKey 是产物名称，例如 diff.patch、commit_log.txt。
	ArtifactKey string `gorm:"size:128;not null" json:"artifact_key"`
	// FilePath 是产物在磁盘上的绝对或相对路径。
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}

// SystemSetting 系统设置。
type SystemSetting struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// Key 是唯一设置名。
	Key string `gorm:"size:128;uniqueIndex;not null" json:"key"`
	// Value 为字符串化后的设置值，当前未做类型系统。
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
