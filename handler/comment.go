package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// GetComments 获取文章评论（只返回已审核通过的，支持嵌套回复）
func GetComments(c *gin.Context) {
	articleID := c.Param("id")
	page, pageSize := getPaginationParams(c)

	var total int64
	query := database.DB.Model(&model.Comment{}).Where("article_id = ? AND status = ?", articleID, "approved")
	query.Count(&total)

	var comments []model.Comment
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取评论失败")
		return
	}

	utils.ResponsePage(c, comments, total, page, pageSize)
}

// SubmitComment 提交评论（支持匿名/登录用户、默认pending状态、通知管理员）
func SubmitComment(c *gin.Context) {
	articleID := c.Param("id")

	var req struct {
		Nickname string `json:"nickname"`
		Content  string `json:"content" binding:"required"`
		ParentID uint   `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	// 检查文章是否存在
	var article model.Article
	if err := database.DB.First(&article, articleID).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "文章不存在")
		return
	}

	ip := utils.GetClientIP(c)
	userID, _ := c.Get("user_id")

	comment := model.Comment{
		ArticleID: article.ID,
		Content:   req.Content,
		ParentID:  req.ParentID,
		Status:    "pending",
		IP:        ip,
	}

	// 登录用户
	if uid, ok := userID.(uint); ok && uid > 0 {
		comment.UserID = uid
		var user model.User
		if err := database.DB.First(&user, uid).Error; err == nil {
			comment.Nickname = user.Nickname
		}
	} else {
		// 匿名用户
		if req.Nickname == "" {
			req.Nickname = "匿名用户"
		}
		comment.Nickname = req.Nickname
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "评论提交失败")
		return
	}

	// 更新文章评论数
	database.DB.Model(&article).UpdateColumn("comment_count", article.CommentCount+1)

	// 通知管理员
	go notifyAdmin("新评论通知",
		fmt.Sprintf("文章《%s》收到新评论：%s", article.Title, req.Content))

	responseSuccess(c, gin.H{"id": comment.ID, "status": comment.Status})
}
