// Package router 负责把所有 HTTP 接口按功能分组注册到 Gin 引擎。
package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitimpact/backend/internal/config"
	"gitimpact/backend/internal/handler"
	"gitimpact/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Register 注册公开接口与受 JWT 保护的业务接口。
// 这里是后端 API 总表，文档中的路由应始终与此文件保持一致。
func Register(r *gin.Engine, authH *handler.AuthHandler, repoH *handler.RepoHandler, taskH *handler.TaskHandler, settingH *handler.SettingHandler, jwtSecret string, frontend config.FrontendConfig) {
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

	registerFrontendFallback(r, frontend)
}

// registerFrontendFallback 在单服务部署模式下托管前端静态资源。
// 它会优先返回真实静态文件，并在命中前端 history 路由时回退到 index.html。
func registerFrontendFallback(r *gin.Engine, frontend config.FrontendConfig) {
	if !frontend.Enabled || strings.TrimSpace(frontend.DistDir) == "" {
		return
	}

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "not found"})
			return
		}
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusNotFound)
			return
		}

		target := frontendFilePath(frontend.DistDir, path)
		if target != "" {
			c.File(target)
			return
		}

		if shouldServeFrontendIndex(c.Request.URL.Path, c.GetHeader("Accept")) {
			indexPath := filepath.Join(frontend.DistDir, "index.html")
			if fileExists(indexPath) {
				c.File(indexPath)
				return
			}
		}

		c.Status(http.StatusNotFound)
	})
}

// frontendFilePath 把 URL 路径映射为 dist 下的实际文件路径。
func frontendFilePath(distDir, requestPath string) string {
	cleaned := strings.TrimPrefix(filepath.Clean(filepath.FromSlash(requestPath)), string(filepath.Separator))
	if cleaned == "." || cleaned == "" {
		candidate := filepath.Join(distDir, "index.html")
		if fileExists(candidate) {
			return candidate
		}
		return ""
	}

	candidate := filepath.Join(distDir, cleaned)
	if fileExists(candidate) {
		return candidate
	}

	return ""
}

// shouldServeFrontendIndex 决定当前请求是否应该回退到前端 index.html。
func shouldServeFrontendIndex(path, accept string) bool {
	if path == "" || path == "/" {
		return true
	}
	if strings.Contains(filepath.Base(path), ".") {
		return false
	}
	if accept == "" {
		return true
	}
	return strings.Contains(accept, "text/html") || strings.Contains(accept, "*/*")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
