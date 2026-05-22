package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// GetAIModelList 获取AI模型提供商列表
func GetAIModelList(c *gin.Context) {
	var models []model.AIModel
	if err := database.DB.Order("id ASC").Find(&models).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取AI模型列表失败")
		return
	}
	responseSuccess(c, models)
}

// CreateAIModel 创建AI模型提供商
func CreateAIModel(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var req struct {
		Name     string `json:"name" binding:"required"`
		Provider string `json:"provider"`
		APIUrl   string `json:"api_url" binding:"required"`
		Models   string `json:"models"`
		Enabled  bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	aiModel := model.AIModel{
		Name:     req.Name,
		Provider: req.Provider,
		APIUrl:   req.APIUrl,
		Models:   req.Models,
		Enabled:  req.Enabled,
	}
	if err := database.DB.Create(&aiModel).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "创建AI模型提供商失败")
		return
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "创建AI模型提供商", req.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, aiModel)
}

// UpdateAIModel 更新AI模型提供商
func UpdateAIModel(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var aiModel model.AIModel
	if err := database.DB.First(&aiModel, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "AI模型提供商不存在")
		return
	}

	var req struct {
		Name     string `json:"name"`
		Provider string `json:"provider"`
		APIUrl   string `json:"api_url"`
		Models   string `json:"models"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Provider != "" {
		updates["provider"] = req.Provider
	}
	if req.APIUrl != "" {
		updates["api_url"] = req.APIUrl
	}
	updates["models"] = req.Models

	database.DB.Model(&aiModel).Updates(updates)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新AI模型提供商", aiModel.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// UpdateAIModelStatus 更新AI模型提供商状态
func UpdateAIModelStatus(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	result := database.DB.Model(&model.AIModel{}).Where("id = ?", id).Update("enabled", req.Enabled)
	if result.RowsAffected == 0 {
		responseErrorWithCode(c, http.StatusNotFound, "AI模型提供商不存在")
		return
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新AI模型提供商状态", id, fmt.Sprintf("enabled=%v", req.Enabled), utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// DeleteAIModel 删除AI模型提供商
func DeleteAIModel(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var aiModel model.AIModel
	if err := database.DB.First(&aiModel, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "AI模型提供商不存在")
		return
	}

	database.DB.Delete(&aiModel)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "删除AI模型提供商", aiModel.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}
