package service

import (
	"testing"

	"gitimpact/backend/internal/config"
	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newAuthSvc(t *testing.T, mode string) (*AuthService, *repository.UserRepository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	_ = db.AutoMigrate(&model.User{})
	repo := repository.NewUserRepository(db)
	cfg := &config.AppConfig{Auth: config.AuthConfig{Mode: mode, EnableRegister: true, JWTSecret: "secret", TokenExpireMinutes: 60}}
	return NewAuthService(cfg, repo), repo
}

func TestRegisterSuccessAndDuplicate(t *testing.T) {
	svc, _ := newAuthSvc(t, "db")
	if err := svc.Register("u1", "p1", "a@b.com"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Register("u1", "p1", "a@b.com"); err == nil {
		t.Fatal("expected duplicate")
	}
}

func TestRegisterDisabled(t *testing.T) {
	svc, _ := newAuthSvc(t, "config")
	if err := svc.Register("u", "p", ""); err == nil {
		t.Fatal("expected disabled")
	}
}

func TestLoginSuccessAndFailure(t *testing.T) {
	svc, repo := newAuthSvc(t, "db")
	hash, _ := bcrypt.GenerateFromPassword([]byte("p1"), bcrypt.DefaultCost)
	_ = repo.Create(&model.User{Username: "u1", PasswordHash: string(hash), Status: model.UserStatusActive, Role: model.RoleUser, Source: model.UserSourceDB})
	if _, err := svc.ValidateLogin("u1", "p1"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ValidateLogin("u1", "wrong"); err == nil {
		t.Fatal("expected failure")
	}
}

func TestDisabledUserLoginFail(t *testing.T) {
	svc, repo := newAuthSvc(t, "db")
	hash, _ := bcrypt.GenerateFromPassword([]byte("p1"), bcrypt.DefaultCost)
	_ = repo.Create(&model.User{Username: "u2", PasswordHash: string(hash), Status: model.UserStatusDisabled, Role: model.RoleUser, Source: model.UserSourceDB})
	if _, err := svc.ValidateLogin("u2", "p1"); err == nil {
		t.Fatal("expected disabled")
	}
}

func TestConfigModeLogin(t *testing.T) {
	svc, _ := newAuthSvc(t, "config")
	svc.cfg.Auth.ConfigUsers = []config.ConfigUser{{Username: "cfg", PasswordPlain: "p", AllowDevPlain: true, Status: model.UserStatusActive, Role: model.RoleAdmin}}
	u, err := svc.ValidateLogin("cfg", "p")
	if err != nil || u.Source != model.UserSourceConfig {
		t.Fatal("expect config user")
	}
}

func TestMixedModeLoginOrder(t *testing.T) {
	svc, repo := newAuthSvc(t, "mixed")
	svc.cfg.Auth.ConfigUsers = []config.ConfigUser{{Username: "same", PasswordPlain: "cfgp", AllowDevPlain: true, Status: model.UserStatusActive, Role: model.RoleAdmin}}
	hash, _ := bcrypt.GenerateFromPassword([]byte("dbp"), bcrypt.DefaultCost)
	_ = repo.Create(&model.User{Username: "same", PasswordHash: string(hash), Status: model.UserStatusActive, Role: model.RoleUser, Source: model.UserSourceDB})
	u, err := svc.ValidateLogin("same", "cfgp")
	if err != nil || u.Source != model.UserSourceConfig {
		t.Fatal("mixed should prioritize config")
	}
	if _, err := svc.ValidateLogin("same", "dbp"); err != nil {
		t.Fatal("db fallback expected")
	}
}
