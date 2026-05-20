package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// getPaginationParams 从query获取page, page_size
func getPaginationParams(c *gin.Context) (page, pageSize int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return
}

// recordOperationLog 记录操作日志
func recordOperationLog(adminID uint, adminName, action, target, result, ip string) {
	log := model.OperationLog{
		AdminID:   adminID,
		AdminName: adminName,
		Action:    action,
		Target:    target,
		Result:    result,
		IP:        ip,
	}
	database.DB.Create(&log)
}

// notifyAdmin 通知管理员（飞书+邮件）
func notifyAdmin(subject, content string) {
	// 飞书通知
	if config := getConfigValue("feishu_enabled"); config == "true" {
		token := getConfigValue("feishu_token")
		if token != "" {
			sendFeishuNotification(token, subject, content)
		}
	}

	// 邮件通知
	adminEmail := getConfigValue("admin_email")
	if adminEmail != "" {
		_ = utils.SendEmail(adminEmail, subject, content)
	}
}

// sendFeishuNotification 发送飞书通知
func sendFeishuNotification(token, title, content string) {
	// 飞书Webhook通知的简化实现
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", token)
	// 实际项目中应使用http.Post发送JSON
	_ = url
}

// getConfigValue 从SiteConfig KV表读取配置
func getConfigValue(key string) string {
	var sc model.SiteConfig
	if err := database.DB.Where("`key` = ?", key).First(&sc).Error; err != nil {
		return ""
	}
	return sc.Value
}

// setConfigValue 写入SiteConfig KV表
func setConfigValue(key, value string) {
	var sc model.SiteConfig
	result := database.DB.Where("`key` = ?", key).First(&sc)
	if result.Error != nil {
		database.DB.Create(&model.SiteConfig{Key: key, Value: value})
	} else {
		database.DB.Model(&sc).Update("value", value)
	}
}

// responseSuccess 统一成功响应
func responseSuccess(c *gin.Context, data interface{}) {
	utils.ResponseSuccess(c, data)
}

// responseError 统一错误响应
func responseError(c *gin.Context, msg string) {
	utils.ResponseError(c, http.StatusBadRequest, msg)
}

// responseErrorWithCode 带HTTP状态码的错误响应
func responseErrorWithCode(c *gin.Context, code int, msg string) {
	utils.ResponseError(c, code, msg)
}
