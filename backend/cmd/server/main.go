package main

import (
	"log"
	"os"

	"gitimpact/backend/internal/analyzer"
	"gitimpact/backend/internal/config"
	"gitimpact/backend/internal/handler"
	"gitimpact/backend/internal/repository"
	"gitimpact/backend/internal/router"
	"gitimpact/backend/internal/service"
	"gitimpact/backend/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	cfgPath := os.Getenv("GITIMPACT_CONFIG")
	if cfgPath == "" {
		cfgPath = "./config.yaml"
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}
	db, err := repository.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	_ = os.MkdirAll(cfg.Workdir.Root, 0o755)
	_ = os.MkdirAll(cfg.Workdir.RepoCache, 0o755)
	_ = os.MkdirAll(cfg.Workdir.Artifacts, 0o755)
	_ = os.MkdirAll(cfg.Workdir.Reports, 0o755)

	userRepo := repository.NewUserRepository(db)
	repoRepo := repository.NewRepoRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	reportRepo := repository.NewReportRepository(db)
	settingRepo := repository.NewSettingRepository(db)

	authSvc := service.NewAuthService(cfg, userRepo)
	if err := authSvc.InitDefaultAdmin(); err != nil {
		log.Printf("init admin failed: %v", err)
	}
	repoSvc := service.NewRepositoryService(repoRepo)
	taskSvc := service.NewTaskService(repoRepo, taskRepo, reportRepo)
	settingSvc := service.NewSettingService(settingRepo)

	workerInst := worker.NewTaskWorker(cfg, analyzer.NewCLIAnalyzer(cfg.OpenCode), taskRepo, repoRepo, reportRepo)

	r := gin.Default()
	router.Register(r, handler.NewAuthHandler(authSvc), handler.NewRepoHandler(repoSvc), handler.NewTaskHandler(taskSvc, workerInst), handler.NewSettingHandler(settingSvc), cfg.Auth.JWTSecret)
	log.Fatal(r.Run(":" + cfg.Server.Port))
}
