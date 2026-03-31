package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitimpact/backend/internal/config"
	"gitimpact/backend/internal/middleware"
	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/repository"
	"gitimpact/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthRouter(t *testing.T, mode string) *gin.Engine {
	t.Helper()
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&model.User{})
	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(&config.AppConfig{Auth: config.AuthConfig{Mode: mode, JWTSecret: "secret", EnableRegister: true, TokenExpireMinutes: 60}}, repo)
	r := gin.New()
	h := NewAuthHandler(svc)
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/me", middleware.JWT("secret"), h.Me)
	return r
}

func TestRegisterLoginAndMe(t *testing.T) {
	r := setupAuthRouter(t, "db")
	payload := []byte(`{"username":"u","password":"p","email":"a@b.com"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("register fail %d", w.Code)
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("login fail %d", w.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	token := resp["data"].(map[string]any)["token"].(string)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("me fail %d", w.Code)
	}
}

func TestLoginFailedAndJWTRequired(t *testing.T) {
	r := setupAuthRouter(t, "db")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username":"x","password":"y"}`)))
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Fatal("expected login failure")
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/me", nil)
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Fatal("expected 401")
	}
}
