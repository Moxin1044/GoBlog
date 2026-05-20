package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// GetSystemConfig 获取系统配置（从SiteConfig KV表读取）
func GetSystemConfig(c *gin.Context) {
	var configs []model.SiteConfig
	if err := database.DB.Find(&configs).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取系统配置失败")
		return
	}

	// 转为KV map
	result := make(map[string]string)
	for _, cfg := range configs {
		result[cfg.Key] = cfg.Value
	}

	responseSuccess(c, result)
}

// UpdateSystemConfig 更新系统配置
func UpdateSystemConfig(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	for key, value := range req {
		setConfigValue(key, value)
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新系统配置", "", "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// GetSiteConfig 前台获取站点配置（网站名称、Logo、版权、ICP、注册开关等）
func GetSiteConfig(c *gin.Context) {
	keys := []string{
		"site_name", "site_logo", "site_copyright", "site_icp",
		"site_description", "site_keywords", "register_enabled",
		"site_footer", "site_notice",
	}

	result := make(map[string]string)
	for _, key := range keys {
		result[key] = getConfigValue(key)
	}

	// 默认值
	if result["site_name"] == "" {
		result["site_name"] = "GoBlog"
	}
	if result["register_enabled"] == "" {
		result["register_enabled"] = "true"
	}

	responseSuccess(c, result)
}
