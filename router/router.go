package router

import (
	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/handler"
	"github.com/moxin/GoBlog/middleware"
)

func SetupRouter(r *gin.Engine) {
	// 公开API
	api := r.Group("/api")
	{
		// 站点配置
		api.GET("/site/config", handler.GetSiteConfig)

		// 用户认证
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
			auth.POST("/verify-code", handler.SendVerifyCode)
		}

		// 文章（公开）
		article := api.Group("/article")
		{
			article.GET("/list", handler.GetArticleList)
			article.GET("/:id", handler.GetArticleDetail)
			article.POST("/:id/like", handler.LikeArticle)
			article.GET("/:id/comments", handler.GetComments)
			article.POST("/:id/comment", handler.SubmitComment)
		}

		// 分类标签
		api.GET("/categories", handler.GetCategories)
		api.GET("/tags", handler.GetTags)

		// RSS
		api.GET("/rss", handler.GetRSS)

		// 用户需要登录的接口
		user := api.Group("/user")
		user.Use(middleware.UserAuth())
		{
			user.GET("/info", handler.GetUserInfo)
			user.PUT("/info", handler.UpdateUserInfo)
			user.PUT("/password", handler.ChangePassword)
			user.POST("/avatar", handler.UploadAvatar)

			// 订阅
			user.GET("/subscription", handler.GetSubscription)
			user.PUT("/subscription", handler.UpdateSubscription)

			// AI配置
			user.GET("/ai-config", handler.GetAIConfig)
			user.PUT("/ai-config", handler.UpdateAIConfig)

			// AI对话
			user.POST("/chat", handler.Chat)
			user.GET("/chat/history", handler.GetChatHistory)
		}

		// 管理员接口
		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuth())
		{
			// 仪表盘
			admin.GET("/dashboard", handler.GetDashboard)
			admin.GET("/monitor", handler.GetServerMonitor)
			admin.GET("/visit-map", handler.GetVisitMapData)

			// 文章管理
			articleAdmin := admin.Group("/article")
			{
				articleAdmin.GET("/list", handler.AdminGetArticleList)
				articleAdmin.POST("", handler.CreateArticle)
				articleAdmin.GET("/:id", handler.AdminGetArticle)
				articleAdmin.PUT("/:id", handler.UpdateArticle)
				articleAdmin.DELETE("/:id", handler.DeleteArticle)
				articleAdmin.PUT("/:id/status", handler.UpdateArticleStatus)
			}

			// 分类管理
			categoryAdmin := admin.Group("/category")
			{
				categoryAdmin.GET("/list", handler.AdminGetCategories)
				categoryAdmin.POST("", handler.CreateCategory)
				categoryAdmin.PUT("/:id", handler.UpdateCategory)
				categoryAdmin.DELETE("/:id", handler.DeleteCategory)
			}

			// 标签管理
			tagAdmin := admin.Group("/tag")
			{
				tagAdmin.GET("/list", handler.AdminGetTags)
				tagAdmin.POST("", handler.CreateTag)
				tagAdmin.PUT("/:id", handler.UpdateTag)
				tagAdmin.DELETE("/:id", handler.DeleteTag)
			}

			// 评论管理
			commentAdmin := admin.Group("/comment")
			{
				commentAdmin.GET("/list", handler.AdminGetComments)
				commentAdmin.PUT("/:id/review", handler.ReviewComment)
				commentAdmin.DELETE("/:id", handler.DeleteComment)
				commentAdmin.POST("/batch-review", handler.BatchReviewComments)
			}

			// 用户管理
			userAdmin := admin.Group("/user")
			{
				userAdmin.GET("/list", handler.AdminGetUsers)
				userAdmin.GET("/:id", handler.AdminGetUser)
				userAdmin.PUT("/:id/status", handler.UpdateUserStatus)
				userAdmin.PUT("/:id/reset-password", handler.ResetUserPassword)
			}

			// 管理员管理（仅超级管理员）
			adminMgmt := admin.Group("/admin-mgmt")
			adminMgmt.Use(middleware.SuperAdminAuth())
			{
				adminMgmt.GET("/list", handler.GetAdminList)
				adminMgmt.POST("", handler.CreateAdmin)
				adminMgmt.PUT("/:id", handler.UpdateAdmin)
				adminMgmt.PUT("/:id/status", handler.UpdateAdminStatus)
			}

			// 系统配置
			configAdmin := admin.Group("/config")
			configAdmin.Use(middleware.SuperAdminAuth())
			{
				configAdmin.GET("", handler.GetSystemConfig)
				configAdmin.PUT("", handler.UpdateSystemConfig)
			}

			// AI模型管理
			aiAdmin := admin.Group("/ai-model")
			aiAdmin.Use(middleware.SuperAdminAuth())
			{
				aiAdmin.GET("/list", handler.GetAIModelList)
				aiAdmin.POST("", handler.CreateAIModel)
				aiAdmin.PUT("/:id", handler.UpdateAIModel)
				aiAdmin.PUT("/:id/status", handler.UpdateAIModelStatus)
				aiAdmin.DELETE("/:id", handler.DeleteAIModel)
			}

			// 数据备份
			backupAdmin := admin.Group("/backup")
			backupAdmin.Use(middleware.SuperAdminAuth())
			{
				backupAdmin.GET("/list", handler.GetBackupList)
				backupAdmin.POST("", handler.CreateBackup)
				backupAdmin.GET("/:id/download", handler.DownloadBackup)
				backupAdmin.DELETE("/:id", handler.DeleteBackup)
			}

			// 操作日志
			admin.GET("/logs", handler.GetOperationLogs)

			// 图片上传
			admin.POST("/upload", handler.UploadImage)
		}
	}
}
