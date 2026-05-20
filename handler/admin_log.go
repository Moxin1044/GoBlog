package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// GetOperationLogs 操作日志列表
func GetOperationLogs(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	query := database.DB.Model(&model.OperationLog{})

	if action := c.Query("action"); action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}
	if adminName := c.Query("admin_name"); adminName != "" {
		query = query.Where("admin_name LIKE ?", "%"+adminName+"%")
	}

	var total int64
	query.Count(&total)

	var logs []model.OperationLog
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取操作日志失败")
		return
	}

	utils.ResponsePage(c, logs, total, page, pageSize)
}
