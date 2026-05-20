package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
)

// GetSubscription 获取用户订阅配置
func GetSubscription(c *gin.Context) {
	userID := c.GetUint("user_id")

	var sub model.Subscription
	if err := database.DB.Preload("Categories").Where("user_id = ?", userID).First(&sub).Error; err != nil {
		// 没有订阅配置，返回空
		responseSuccess(c, gin.H{
			"user_id":       userID,
			"email":         "",
			"feishu_token":  "",
			"notify_email":  false,
			"notify_feishu": false,
			"categories":    []model.Category{},
		})
		return
	}

	responseSuccess(c, sub)
}

// UpdateSubscription 更新订阅配置（邮箱、飞书Token、通知渠道、订阅分类）
func UpdateSubscription(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Email        string `json:"email"`
		FeishuToken  string `json:"feishu_token"`
		NotifyEmail  bool   `json:"notify_email"`
		NotifyFeishu bool   `json:"notify_feishu"`
		CategoryIDs  []uint `json:"category_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	var sub model.Subscription
	result := database.DB.Where("user_id = ?", userID).First(&sub)
	if result.Error != nil {
		// 创建新订阅
		sub = model.Subscription{
			UserID:       userID,
			Email:        req.Email,
			FeishuToken:  req.FeishuToken,
			NotifyEmail:  req.NotifyEmail,
			NotifyFeishu: req.NotifyFeishu,
		}
		if err := database.DB.Create(&sub).Error; err != nil {
			responseErrorWithCode(c, http.StatusInternalServerError, "创建订阅失败")
			return
		}
	} else {
		// 更新订阅
		database.DB.Model(&sub).Updates(map[string]interface{}{
			"email":         req.Email,
			"feishu_token":  req.FeishuToken,
			"notify_email":  req.NotifyEmail,
			"notify_feishu": req.NotifyFeishu,
		})
	}

	// 更新订阅分类
	if req.CategoryIDs != nil {
		database.DB.Model(&sub).Association("Categories").Replace(req.CategoryIDs)
	}

	responseSuccess(c, nil)
}
