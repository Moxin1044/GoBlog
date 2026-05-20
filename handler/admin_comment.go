package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// AdminGetComments 评论列表（支持状态筛选）
func AdminGetComments(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	query := database.DB.Model(&model.Comment{})

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if articleID := c.Query("article_id"); articleID != "" {
		query = query.Where("article_id = ?", articleID)
	}

	var total int64
	query.Count(&total)

	var comments []model.Comment
	offset := (page - 1) * pageSize
	if err := query.Preload("Article").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取评论列表失败")
		return
	}

	utils.ResponsePage(c, comments, total, page, pageSize)
}

// ReviewComment 审核评论
func ReviewComment(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var comment model.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "评论不存在")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	now := time.Now()
	database.DB.Model(&comment).Updates(map[string]interface{}{
		"status":      req.Status,
		"reviewed_at": &now,
		"reviewer_id": adminID,
	})

	// 如果审核通过，更新文章评论数
	if req.Status == "approved" {
		database.DB.Model(&model.Article{}).Where("id = ?", comment.ArticleID).
			UpdateColumn("comment_count", database.DB.Model(&model.Comment{}).
				Where("article_id = ? AND status = ?", comment.ArticleID, "approved").Select("COUNT(*)"))
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "审核评论", id, req.Status, utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var comment model.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "评论不存在")
		return
	}

	database.DB.Delete(&comment)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "删除评论", id, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// BatchReviewComments 批量审核
func BatchReviewComments(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var req struct {
		IDs    []uint `json:"ids" binding:"required"`
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	now := time.Now()
	database.DB.Model(&model.Comment{}).Where("id IN ?", req.IDs).Updates(map[string]interface{}{
		"status":      req.Status,
		"reviewed_at": &now,
		"reviewer_id": adminID,
	})

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "批量审核评论", fmt.Sprintf("共%d条", len(req.IDs)), req.Status, utils.GetClientIP(c))

	responseSuccess(c, nil)
}
