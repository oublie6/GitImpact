// router_test.go 验证前端静态托管和 history 路由回退行为。
package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"gitimpact/backend/internal/config"
	"gitimpact/backend/internal/handler"
	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/repository"
	"gitimpact/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouterForStaticTest(t *testing.T, distDir string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Repository{}, &model.AnalysisTask{}, &model.AnalysisReport{}, &model.TaskLog{}, &model.SystemSetting{}, &model.TaskArtifact{}); err != nil {
		t.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	repoRepo := repository.NewRepoRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	reportRepo := repository.NewReportRepository(db)
	settingRepo := repository.NewSettingRepository(db)

	authSvc := service.NewAuthService(&config.AppConfig{Auth: config.AuthConfig{Mode: "db", JWTSecret: "secret", EnableRegister: true, TokenExpireMinutes: 60}}, userRepo)
	repoSvc := service.NewRepositoryService(repoRepo)
	taskSvc := service.NewTaskService(repoRepo, taskRepo, reportRepo)
	settingSvc := service.NewSettingService(settingRepo)

	r := gin.New()
	Register(
		r,
		handler.NewAuthHandler(authSvc),
		handler.NewRepoHandler(repoSvc),
		handler.NewTaskHandler(taskSvc, nil),
		handler.NewSettingHandler(settingSvc),
		"secret",
		config.FrontendConfig{Enabled: true, DistDir: distDir},
	)
	return r
}

func TestRegisterServesFrontendIndexAndAssets(t *testing.T) {
	distDir := t.TempDir()
	assetsDir := filepath.Join(distDir, "assets")
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(distDir, "index.html"), []byte("<html><body>frontend</body></html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assetsDir, "app.js"), []byte("console.log('ok')"), 0o644); err != nil {
		t.Fatal(err)
	}

	r := setupRouterForStaticTest(t, distDir)

	for _, target := range []string{"/", "/tasks/1", "/repositories"} {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, target, nil)
		req.Header.Set("Accept", "text/html")
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("%s should return index.html, got %d", target, w.Code)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/assets/app.js", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("asset should be served, got %d", w.Code)
	}
}

func TestRegisterDoesNotHijackMissingAPI(t *testing.T) {
	distDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(distDir, "index.html"), []byte("<html><body>frontend</body></html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	r := setupRouterForStaticTest(t, distDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/not-found", nil)
	req.Header.Set("Accept", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("missing api should return 404, got %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["code"] != float64(http.StatusNotFound) {
		t.Fatalf("unexpected response body: %v", body)
	}
}

func TestJWTProtectedRoutesRemainProtected(t *testing.T) {
	distDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(distDir, "index.html"), []byte("<html><body>frontend</body></html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := setupRouterForStaticTest(t, distDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/tasks", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected protected api to remain unauthorized, got %d", w.Code)
	}
}
