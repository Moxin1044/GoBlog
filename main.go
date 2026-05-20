package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/config"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/middleware"
	"github.com/moxin/GoBlog/router"
	"github.com/moxin/GoBlog/service"
)

func main() {
	// 加载配置
	if err := config.InitConfig("config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(&config.AppConfig.Database); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 启动定时任务
	go service.StartCronJobs()

	// 设置Gin模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())

	// 静态文件
	r.Static("/uploads", config.AppConfig.Upload.Path)

	// 注册API路由
	router.SetupRouter(r)

	// 前端静态文件服务
	setupFrontend(r)

	// 启动服务
	log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
	if err := r.Run(":" + config.AppConfig.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupFrontend 配置前端静态文件服务和SPA路由回退
func setupFrontend(r *gin.Engine) {
	distDir := filepath.Join("web", "dist")
	if _, err := os.Stat(distDir); os.IsNotExist(err) {
		log.Println("Frontend dist directory not found, skipping static file serving")
		return
	}

	// 静态资源（JS/CSS/图片等）
	r.Static("/assets", filepath.Join(distDir, "assets"))

	// 其他静态文件（favicon等）
	r.StaticFile("/vite.svg", filepath.Join(distDir, "vite.svg"))
	r.StaticFile("/favicon.ico", filepath.Join(distDir, "favicon.ico"))

	// SPA回退：所有非API、非静态文件的请求返回index.html
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// API请求不回退
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
			return
		}
		// uploads请求不回退
		if len(path) >= 8 && path[:8] == "/uploads" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
			return
		}
		c.File(filepath.Join(distDir, "index.html"))
	})
}
