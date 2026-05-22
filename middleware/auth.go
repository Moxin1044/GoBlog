package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/utils"
)

// CORS中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// JWT用户认证中间件（同时支持user和admin token）
func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			utils.ResponseError(c, http.StatusUnauthorized, "未登录或token无效")
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.ResponseError(c, http.StatusUnauthorized, "token无效或已过期")
			c.Abort()
			return
		}
		// 同时接受 user 和 admin token
		if claims.Type == "user" {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
		} else if claims.Type == "admin" {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", "admin")
			c.Set("admin_id", claims.UserID)
			c.Set("admin_name", claims.Username)
			c.Set("admin_role", claims.Role)
		} else {
			utils.ResponseError(c, http.StatusUnauthorized, "无效的token类型")
			c.Abort()
			return
		}
		c.Next()
	}
}

// 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			utils.ResponseError(c, http.StatusUnauthorized, "未登录或token无效")
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.ResponseError(c, http.StatusUnauthorized, "token无效或已过期")
			c.Abort()
			return
		}
		if claims.Type != "admin" {
			utils.ResponseError(c, http.StatusUnauthorized, "无效的token类型")
			c.Abort()
			return
		}
		c.Set("admin_id", claims.UserID)
		c.Set("admin_name", claims.Username)
		c.Set("admin_role", claims.Role)
		c.Next()
	}
}

// 超级管理员权限中间件
func SuperAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("admin_role")
		if role != "superadmin" {
			utils.ResponseError(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}

// 限流中间件（简单实现）
func RateLimit(maxRequests int, duration time.Duration) gin.HandlerFunc {
	limiter := make(map[string][]time.Time)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		requests := limiter[ip]
		// 清理过期记录
		validRequests := []time.Time{}
		for _, t := range requests {
			if now.Sub(t) < duration {
				validRequests = append(validRequests, t)
			}
		}
		if len(validRequests) >= maxRequests {
			utils.ResponseError(c, http.StatusTooManyRequests, "请求过于频繁")
			c.Abort()
			return
		}
		limiter[ip] = append(validRequests, now)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return c.Query("token")
}
