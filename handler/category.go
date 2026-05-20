package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
)

// GetCategories 获取所有分类
func GetCategories(c *gin.Context) {
	var categories []model.Category
	if err := database.DB.Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取分类失败")
		return
	}
	responseSuccess(c, categories)
}

// GetTags 获取所有标签
func GetTags(c *gin.Context) {
	var tags []model.Tag
	if err := database.DB.Order("id ASC").Find(&tags).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取标签失败")
		return
	}
	responseSuccess(c, tags)
}
