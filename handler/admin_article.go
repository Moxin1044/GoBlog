package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// AdminGetArticleList 文章列表（含草稿、支持筛选）
func AdminGetArticleList(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	query := database.DB.Model(&model.Article{})

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 关键词搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var articles []model.Article
	offset := (page - 1) * pageSize
	if err := query.Preload("Category").Preload("Tags").Preload("Author").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取文章列表失败")
		return
	}

	utils.ResponsePage(c, articles, total, page, pageSize)
}

// CreateArticle 创建文章（含分类、标签关联）
func CreateArticle(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var req struct {
		Title      string `json:"title" binding:"required"`
		Cover      string `json:"cover"`
		Content    string `json:"content" binding:"required"`
		Summary    string `json:"summary"`
		CategoryID uint   `json:"category_id"`
		TagIDs     []uint `json:"tag_ids"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	if req.Status == "" {
		req.Status = "draft"
	}

	article := model.Article{
		Title:      req.Title,
		Cover:      req.Cover,
		Content:    req.Content,
		Summary:    req.Summary,
		CategoryID: req.CategoryID,
		AuthorID:   adminID,
		Status:     req.Status,
	}

	if req.Status == "published" {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := database.DB.Create(&article).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "创建文章失败")
		return
	}

	// 关联标签
	if len(req.TagIDs) > 0 {
		var tags []model.Tag
		database.DB.Where("id IN ?", req.TagIDs).Find(&tags)
		database.DB.Model(&article).Association("Tags").Replace(tags)
	}

	// 记录操作日志
	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "创建文章", article.Title, "成功", utils.GetClientIP(c))

	responseSuccess(c, gin.H{"id": article.ID})
}

// AdminGetArticle 获取文章详情（管理用）
func AdminGetArticle(c *gin.Context) {
	id := c.Param("id")
	var article model.Article
	if err := database.DB.Preload("Category").Preload("Tags").Preload("Author").First(&article, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "文章不存在")
		return
	}
	responseSuccess(c, article)
}

// UpdateArticle 更新文章
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var article model.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "文章不存在")
		return
	}

	var req struct {
		Title      string `json:"title"`
		Cover      string `json:"cover"`
		Content    string `json:"content"`
		Summary    string `json:"summary"`
		CategoryID uint   `json:"category_id"`
		TagIDs     []uint `json:"tag_ids"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Status != "" {
		updates["status"] = req.Status
		if req.Status == "published" && article.PublishedAt == nil {
			now := time.Now()
			updates["published_at"] = &now
		}
	}

	database.DB.Model(&article).Updates(updates)

	// 更新标签关联
	if req.TagIDs != nil {
		var tags []model.Tag
		database.DB.Where("id IN ?", req.TagIDs).Find(&tags)
		database.DB.Model(&article).Association("Tags").Replace(tags)
	}

	// 记录操作日志
	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新文章", article.Title, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var article model.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "文章不存在")
		return
	}

	database.DB.Select("Tags").Delete(&article)

	// 记录操作日志
	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "删除文章", article.Title, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// UpdateArticleStatus 更新文章状态（发布/下架）
func UpdateArticleStatus(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		responseError(c, "无效的文章ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{
		"status": req.Status,
	}
	if req.Status == "published" {
		now := time.Now()
		updates["published_at"] = &now
	}

	database.DB.Model(&model.Article{}).Where("id = ?", articleID).Updates(updates)

	// 记录操作日志
	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新文章状态", id, req.Status, utils.GetClientIP(c))

	responseSuccess(c, nil)
}
