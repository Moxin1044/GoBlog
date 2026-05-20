package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// AdminGetCategories 管理员获取分类列表
func AdminGetCategories(c *gin.Context) {
	var categories []model.Category
	if err := database.DB.Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取分类失败")
		return
	}
	responseSuccess(c, categories)
}

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var req struct {
		Name   string `json:"name" binding:"required"`
		NameEn string `json:"name_en"`
		Sort   int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	category := model.Category{
		Name:   req.Name,
		NameEn: req.NameEn,
		Sort:   req.Sort,
	}
	if err := database.DB.Create(&category).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "创建分类失败")
		return
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "创建分类", req.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, category)
}

// UpdateCategory 更新分类
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var category model.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "分类不存在")
		return
	}

	var req struct {
		Name   string `json:"name"`
		NameEn string `json:"name_en"`
		Sort   int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	updates["name_en"] = req.NameEn
	updates["sort"] = req.Sort

	database.DB.Model(&category).Updates(updates)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新分类", category.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var category model.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "分类不存在")
		return
	}

	database.DB.Delete(&category)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "删除分类", category.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// AdminGetTags 管理员获取标签列表
func AdminGetTags(c *gin.Context) {
	var tags []model.Tag
	if err := database.DB.Order("id ASC").Find(&tags).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取标签失败")
		return
	}
	responseSuccess(c, tags)
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var req struct {
		Name   string `json:"name" binding:"required"`
		NameEn string `json:"name_en"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	tag := model.Tag{
		Name:   req.Name,
		NameEn: req.NameEn,
	}
	if err := database.DB.Create(&tag).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "创建标签失败")
		return
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "创建标签", req.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, tag)
}

// UpdateTag 更新标签
func UpdateTag(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var tag model.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "标签不存在")
		return
	}

	var req struct {
		Name   string `json:"name"`
		NameEn string `json:"name_en"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	updates["name_en"] = req.NameEn

	database.DB.Model(&tag).Updates(updates)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新标签", tag.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var tag model.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "标签不存在")
		return
	}

	database.DB.Delete(&tag)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "删除标签", tag.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}
