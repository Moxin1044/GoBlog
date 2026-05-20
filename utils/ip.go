package utils

import "github.com/gin-gonic/gin"

func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Real-IP")
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}

func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

func IsAjaxRequest(c *gin.Context) bool {
	return c.GetHeader("X-Requested-With") == "XMLHttpRequest" ||
		c.GetHeader("Accept") != "" &&
			c.GetHeader("Accept") != "text/html"
}
