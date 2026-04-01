// Package router 负责把所有 HTTP 接口按功能分组注册到 Gin 引擎。
package router

import (
	"gitimpact/backend/internal/handler"
	"gitimpact/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Register 注册公开接口与受 JWT 保护的业务接口。
// 这里是后端 API 总表，文档中的路由应始终与此文件保持一致。
func Register(r *gin.Engine, authH *handler.AuthHandler, repoH *handler.RepoHandler, taskH *handler.TaskHandler, settingH *handler.SettingHandler, jwtSecret string) {
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/logout", middleware.JWT(jwtSecret), authH.Logout)
		auth.GET("/me", middleware.JWT(jwtSecret), authH.Me)

		// 核心业务接口全部依赖 JWT 鉴权，中间件会把用户信息写入上下文。
		core := api.Group("")
		core.Use(middleware.JWT(jwtSecret))
		core.POST("/repositories", repoH.Create)
		core.PUT("/repositories/:id", repoH.Update)
		core.GET("/repositories", repoH.List)
		core.GET("/repositories/:id", repoH.Detail)
		core.POST("/repositories/:id/fetch", repoH.Fetch)

		core.POST("/tasks", taskH.Create)
		core.GET("/tasks", taskH.List)
		core.GET("/tasks/:id", taskH.Detail)
		core.GET("/tasks/:id/logs", taskH.Logs)
		core.GET("/tasks/:id/report", taskH.Report)
		core.GET("/tasks/:id/report/download", taskH.DownloadReport)

		core.GET("/settings", settingH.List)
		core.POST("/settings", settingH.Save)
	}
}
