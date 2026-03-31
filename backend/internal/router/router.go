package router

import (
	"gitimpact/backend/internal/handler"
	"gitimpact/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, authH *handler.AuthHandler, repoH *handler.RepoHandler, taskH *handler.TaskHandler, settingH *handler.SettingHandler, jwtSecret string) {
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/logout", middleware.JWT(jwtSecret), authH.Logout)
		auth.GET("/me", middleware.JWT(jwtSecret), authH.Me)

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
