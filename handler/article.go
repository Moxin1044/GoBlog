package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// GetArticleList 文章列表（分页、分类筛选、标签筛选、关键词搜索、只返回已发布文章）
func GetArticleList(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	query := database.DB.Model(&model.Article{}).Where("status = ?", "published")

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 标签筛选
	if tagID := c.Query("tag_id"); tagID != "" {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", tagID)
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
		Order("published_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取文章列表失败")
		return
	}

	utils.ResponsePage(c, articles, total, page, pageSize)
}

// GetArticleDetail 文章详情（浏览量+1、返回文章完整信息含分类标签、评论数、点赞数）
func GetArticleDetail(c *gin.Context) {
	id := c.Param("id")
	var article model.Article
	if err := database.DB.Preload("Category").Preload("Tags").Preload("Author").
		First(&article, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "文章不存在")
		return
	}

	// 浏览量+1
	database.DB.Model(&article).UpdateColumn("view_count", article.ViewCount+1)
	article.ViewCount++

	// 记录访问日志
	ip := utils.GetClientIP(c)
	database.DB.Create(&model.VisitLog{
		IP:        ip,
		Path:      "/api/article/" + id,
		UserAgent: utils.GetUserAgent(c),
	})

	responseSuccess(c, article)
}

// LikeArticle 点赞（IP/用户去重、更新点赞数）
func LikeArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		responseError(c, "无效的文章ID")
		return
	}

	// 检查文章是否存在
	var article model.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "文章不存在")
		return
	}

	ip := utils.GetClientIP(c)
	userID, _ := c.Get("user_id")

	// 去重检查
	var like model.Like
	likeQuery := database.DB.Where("article_id = ?", id)
	if uid, ok := userID.(uint); ok && uid > 0 {
		likeQuery = likeQuery.Where("user_id = ? OR ip = ?", uid, ip)
	} else {
		likeQuery = likeQuery.Where("ip = ?", ip)
	}
	if err := likeQuery.First(&like).Error; err == nil {
		responseError(c, "已经点过赞了")
		return
	}

	// 创建点赞记录
	newLike := model.Like{
		ArticleID: uint(id),
		IP:        ip,
	}
	if uid, ok := userID.(uint); ok && uid > 0 {
		newLike.UserID = uid
	}
	database.DB.Create(&newLike)

	// 更新点赞数
	database.DB.Model(&article).UpdateColumn("like_count", article.LikeCount+1)

	responseSuccess(c, gin.H{"like_count": article.LikeCount + 1})
}
