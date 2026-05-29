package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// GetNavigationList 获取导航列表（前台，只返回启用的）
func GetNavigationList(c *gin.Context) {
	var navigations []model.Navigation
	if err := database.DB.Where("enabled = ?", true).Order("sort ASC, id ASC").Find(&navigations).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取导航列表失败")
		return
	}

	tree := buildNavigationTree(navigations, nil)
	responseSuccess(c, tree)
}

// AdminGetNavigationList 管理员获取导航列表（含禁用项）
func AdminGetNavigationList(c *gin.Context) {
	var navigations []model.Navigation
	if err := database.DB.Order("sort ASC, id ASC").Find(&navigations).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取导航列表失败")
		return
	}

	tree := buildNavigationTree(navigations, nil)
	responseSuccess(c, tree)
}

// GetNavigation 获取单个导航
func GetNavigation(c *gin.Context) {
	id := c.Param("id")
	var navigation model.Navigation
	if err := database.DB.First(&navigation, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "导航不存在")
		return
	}
	responseSuccess(c, navigation)
}

// CreateNavigation 创建导航
func CreateNavigation(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var req struct {
		Name       string `json:"name" binding:"required"`
		NameEn     string `json:"name_en"`
		ParentID   *uint  `json:"parent_id"`
		Type       string `json:"type"`
		Link       string `json:"link"`
		CategoryID uint   `json:"category_id"`
		NewTab     bool   `json:"new_tab"`
		Sort       int    `json:"sort"`
		Enabled    bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	navigation := model.Navigation{
		Name:       req.Name,
		NameEn:     req.NameEn,
		ParentID:   req.ParentID,
		Type:       req.Type,
		Link:       req.Link,
		CategoryID: req.CategoryID,
		NewTab:     req.NewTab,
		Sort:       req.Sort,
		Enabled:    req.Enabled,
	}
	if err := database.DB.Create(&navigation).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "创建导航失败")
		return
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "创建导航", req.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, navigation)
}

// UpdateNavigation 更新导航
func UpdateNavigation(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var navigation model.Navigation
	if err := database.DB.First(&navigation, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "导航不存在")
		return
	}

	var req struct {
		Name       string `json:"name"`
		NameEn     string `json:"name_en"`
		ParentID   *uint  `json:"parent_id"`
		Type       string `json:"type"`
		Link       string `json:"link"`
		CategoryID uint   `json:"category_id"`
		NewTab     bool   `json:"new_tab"`
		Sort       int    `json:"sort"`
		Enabled    bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"name_en":     req.NameEn,
		"parent_id":   req.ParentID,
		"type":        req.Type,
		"link":        req.Link,
		"category_id": req.CategoryID,
		"new_tab":     req.NewTab,
		"sort":        req.Sort,
		"enabled":     req.Enabled,
	}
	database.DB.Model(&navigation).Updates(updates)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新导航", navigation.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// DeleteNavigation 删除导航
func DeleteNavigation(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var navigation model.Navigation
	if err := database.DB.First(&navigation, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "导航不存在")
		return
	}

	// 检查是否有子导航
	var count int64
	database.DB.Model(&model.Navigation{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		responseError(c, "请先删除子导航")
		return
	}

	database.DB.Delete(&navigation)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "删除导航", navigation.Name, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// UpdateNavigationSort 更新导航排序
func UpdateNavigationSort(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var sortData []struct {
		ID   uint `json:"id"`
		Sort int  `json:"sort"`
	}
	if err := c.ShouldBindJSON(&sortData); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	for _, item := range sortData {
		database.DB.Model(&model.Navigation{}).Where("id = ?", item.ID).Update("sort", item.Sort)
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新导航排序", strconv.Itoa(len(sortData))+"项", "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// buildNavigationTree 构建导航树形结构
func buildNavigationTree(navigations []model.Navigation, parentID *uint) []model.Navigation {
	var tree []model.Navigation
	for _, nav := range navigations {
		if equalParentID(nav.ParentID, parentID) {
			children := buildNavigationTree(navigations, &nav.ID)
			nav.Children = children
			tree = append(tree, nav)
		}
	}
	return tree
}

func equalParentID(a, b *uint) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
