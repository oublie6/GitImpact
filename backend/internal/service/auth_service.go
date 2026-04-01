// auth_service.go 实现注册、登录、配置用户优先级和默认管理员初始化。
//
// 该服务同时连接配置文件中的静态用户与数据库用户，是认证链路的核心。
package service

import (
	"errors"
	"strings"
	"time"

	"gitimpact/backend/internal/config"
	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/repository"
	"gitimpact/backend/pkg/auth"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	// ErrRegisterDisabled 表示当前认证模式不允许开放注册。
	ErrRegisterDisabled = errors.New("register disabled")
	// ErrInvalidLogin 表示用户名不存在或密码不匹配。
	ErrInvalidLogin = errors.New("invalid username or password")
	// ErrUserDisabled 表示账号已被停用。
	ErrUserDisabled = errors.New("user disabled")
	// ErrUserExists 表示注册时用户名已存在。
	ErrUserExists = errors.New("username already exists")
)

// AuthService 处理注册/登录与用户来源校验。
type AuthService struct {
	cfg      *config.AppConfig
	userRepo *repository.UserRepository
}

func NewAuthService(cfg *config.AppConfig, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{cfg: cfg, userRepo: userRepo}
}

// InitDefaultAdmin 在允许的模式下确保数据库内至少存在 admin 用户。
// 这里不会影响 config_users，只针对数据库用户初始化。
func (s *AuthService) InitDefaultAdmin() error {
	if !s.cfg.Auth.InitAdminEnabled || (s.cfg.Auth.Mode == "config") {
		return nil
	}
	_, err := s.userRepo.GetByUsername("admin")
	if err == nil {
		return nil
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("Admin@123456"), bcrypt.DefaultCost)
	return s.userRepo.Create(&model.User{Username: "admin", PasswordHash: string(hash), DisplayName: "系统管理员", Role: model.RoleAdmin, Status: model.UserStatusActive, Source: model.UserSourceDB})
}

// Register 创建数据库用户。
// 配置模式下禁止注册，mixed 模式允许注册但只写入数据库表。
func (s *AuthService) Register(username, password, email string) error {
	if !s.cfg.Auth.EnableRegister || s.cfg.Auth.Mode == "config" {
		return ErrRegisterDisabled
	}
	if _, err := s.userRepo.GetByUsername(username); err == nil {
		return ErrUserExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userRepo.Create(&model.User{Username: username, PasswordHash: string(hash), DisplayName: username, Email: email, Role: model.RoleUser, Status: model.UserStatusActive, Source: model.UserSourceDB})
}

// ValidateLogin 登录校验顺序：mixed 模式先校验 config_users，再校验 DB。
func (s *AuthService) ValidateLogin(username, password string) (*model.User, error) {
	// config 与 mixed 模式先尝试命中静态配置用户，保持代码行为与文档一致。
	if s.cfg.Auth.Mode != "db" {
		if user, ok := s.validateConfigUser(username, password); ok {
			return user, nil
		}
		if s.cfg.Auth.Mode == "config" {
			return nil, ErrInvalidLogin
		}
	}
	if s.cfg.Auth.Mode != "config" {
		u, err := s.userRepo.GetByUsername(username)
		if err != nil {
			return nil, ErrInvalidLogin
		}
		if u.Status == model.UserStatusDisabled {
			return nil, ErrUserDisabled
		}
		// 数据库用户统一使用 bcrypt 校验，避免与 config_users 的明文密码策略混淆。
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
			return nil, ErrInvalidLogin
		}
		now := time.Now()
		u.LastLoginAt = &now
		_ = s.userRepo.Update(u)
		return u, nil
	}
	return nil, ErrInvalidLogin
}

// validateConfigUser 在配置用户列表中顺序查找匹配账号。
// 返回值中的 ID 固定为 0，因为该用户并未持久化到数据库。
func (s *AuthService) validateConfigUser(username, password string) (*model.User, bool) {
	for _, cu := range s.cfg.Auth.ConfigUsers {
		if cu.Username != username {
			continue
		}
		if cu.Status == model.UserStatusDisabled {
			return nil, false
		}
		if cu.Password != "" && cu.Password == password {
			return &model.User{ID: 0, Username: cu.Username, DisplayName: cu.DisplayName, Email: cu.Email, Role: normalizeRole(cu.Role), Status: cu.Status, Source: model.UserSourceConfig}, true
		}
	}
	return nil, false
}

// normalizeRole 为未填写角色的配置用户补默认值。
func normalizeRole(role string) string {
	if strings.TrimSpace(role) == "" {
		return model.RoleUser
	}
	return role
}

// GenerateToken 根据当前配置生成 JWT。
func (s *AuthService) GenerateToken(u *model.User) (string, error) {
	return auth.GenerateToken(s.cfg.Auth.JWTSecret, s.cfg.TokenTTL(), u.ID, u.Username, u.Role)
}

// GetDBUserByUsername 仅从数据库读取用户，供需要绕开 config_users 的场景复用。
func (s *AuthService) GetDBUserByUsername(username string) (*model.User, error) {
	u, err := s.userRepo.GetByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidLogin
	}
	return u, err
}
