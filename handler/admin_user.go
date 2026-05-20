package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// AdminGetUsers 用户列表
func AdminGetUsers(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	query := database.DB.Model(&model.User{})

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var users []model.User
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	utils.ResponsePage(c, users, total, page, pageSize)
}

// AdminGetUser 用户详情
func AdminGetUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "用户不存在")
		return
	}
	responseSuccess(c, user)
}

// UpdateUserStatus 启用/禁用用户
func UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	result := database.DB.Model(&model.User{}).Where("id = ?", id).Update("status", req.Status)
	if result.RowsAffected == 0 {
		responseErrorWithCode(c, http.StatusNotFound, "用户不存在")
		return
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新用户状态", id, req.Status, utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// ResetUserPassword 重置用户密码
func ResetUserPassword(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	result := database.DB.Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword)
	if result.RowsAffected == 0 {
		responseErrorWithCode(c, http.StatusNotFound, "用户不存在")
		return
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "重置用户密码", id, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}
