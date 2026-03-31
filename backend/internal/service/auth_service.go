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
	ErrRegisterDisabled = errors.New("register disabled")
	ErrInvalidLogin     = errors.New("invalid username or password")
	ErrUserDisabled     = errors.New("user disabled")
	ErrUserExists       = errors.New("username already exists")
)

// AuthService 处理注册/登录与用户来源校验。
type AuthService struct {
	cfg      *config.AppConfig
	userRepo *repository.UserRepository
}

func NewAuthService(cfg *config.AppConfig, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{cfg: cfg, userRepo: userRepo}
}

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

func (s *AuthService) validateConfigUser(username, password string) (*model.User, bool) {
	for _, cu := range s.cfg.Auth.ConfigUsers {
		if cu.Username != username {
			continue
		}
		if cu.Status == model.UserStatusDisabled {
			return nil, false
		}
		if cu.PasswordHash != "" && bcrypt.CompareHashAndPassword([]byte(cu.PasswordHash), []byte(password)) == nil {
			return &model.User{ID: 0, Username: cu.Username, DisplayName: cu.DisplayName, Email: cu.Email, Role: normalizeRole(cu.Role), Status: cu.Status, Source: model.UserSourceConfig}, true
		}
		if cu.AllowDevPlain && cu.PasswordPlain != "" && cu.PasswordPlain == password {
			return &model.User{ID: 0, Username: cu.Username, DisplayName: cu.DisplayName, Email: cu.Email, Role: normalizeRole(cu.Role), Status: cu.Status, Source: model.UserSourceConfig}, true
		}
	}
	return nil, false
}

func normalizeRole(role string) string {
	if strings.TrimSpace(role) == "" {
		return model.RoleUser
	}
	return role
}

func (s *AuthService) GenerateToken(u *model.User) (string, error) {
	return auth.GenerateToken(s.cfg.Auth.JWTSecret, s.cfg.TokenTTL(), u.ID, u.Username, u.Role)
}

func (s *AuthService) GetDBUserByUsername(username string) (*model.User, error) {
	u, err := s.userRepo.GetByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidLogin
	}
	return u, err
}
