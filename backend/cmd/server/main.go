// GitImpact 后端服务启动入口。
//
// 这个文件负责串联配置加载、数据库初始化、仓储与服务装配、任务 worker 创建
// 以及 Gin 路由注册，是理解整个后端启动链路的最佳入口。
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
	// 优先读取环境变量指定的配置文件，便于不同部署环境复用同一套二进制。
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

	// 运行时目录用于存放仓库缓存、任务材料和报告文件；目录不存在时自动创建。
	_ = os.MkdirAll(cfg.Workdir.Root, 0o755)
	_ = os.MkdirAll(cfg.Workdir.RepoCache, 0o755)
	_ = os.MkdirAll(cfg.Workdir.Artifacts, 0o755)
	_ = os.MkdirAll(cfg.Workdir.Reports, 0o755)

	// 先构造 repository，再向上组装 service 和 worker，保持依赖方向单向清晰。
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

	// 当前分析器实现为 CLIAnalyzer，ServerAnalyzer 仍然是后续扩展预留。
	workerInst := worker.NewTaskWorker(cfg, analyzer.NewCLIAnalyzer(cfg.OpenCode), taskRepo, repoRepo, reportRepo)

	r := gin.Default()
	router.Register(r, handler.NewAuthHandler(authSvc), handler.NewRepoHandler(repoSvc), handler.NewTaskHandler(taskSvc, workerInst), handler.NewSettingHandler(settingSvc), cfg.Auth.JWTSecret, cfg.Frontend)

	// Gin 默认监听在配置指定端口，启动失败直接终止进程，避免服务处于半初始化状态。
	log.Fatal(r.Run(":" + cfg.Server.Port))
}
