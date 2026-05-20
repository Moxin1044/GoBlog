package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// GetAdminList 管理员列表
func GetAdminList(c *gin.Context) {
	var admins []model.Admin
	if err := database.DB.Order("id ASC").Find(&admins).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取管理员列表失败")
		return
	}
	responseSuccess(c, admins)
}

// CreateAdmin 创建管理员
func CreateAdmin(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"omitempty,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	if req.Role == "" {
		req.Role = "admin"
	}

	// 检查用户名唯一性
	var count int64
	database.DB.Model(&model.Admin{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		responseError(c, "用户名已存在")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	admin := model.Admin{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		Status:   "active",
	}
	if err := database.DB.Create(&admin).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "创建管理员失败")
		return
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "创建管理员", req.Username, "成功", utils.GetClientIP(c))

	responseSuccess(c, gin.H{"id": admin.ID})
}

// UpdateAdmin 更新管理员
func UpdateAdmin(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var admin model.Admin
	if err := database.DB.First(&admin, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "管理员不存在")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			responseErrorWithCode(c, http.StatusInternalServerError, "密码加密失败")
			return
		}
		updates["password"] = hashedPassword
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}

	database.DB.Model(&admin).Updates(updates)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新管理员", admin.Username, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

// UpdateAdminStatus 启用/禁用管理员
func UpdateAdminStatus(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	// 不能禁用自己
	if fmt.Sprintf("%d", adminID) == id && req.Status == "disabled" {
		responseError(c, "不能禁用自己")
		return
	}

	result := database.DB.Model(&model.Admin{}).Where("id = ?", id).Update("status", req.Status)
	if result.RowsAffected == 0 {
		responseErrorWithCode(c, http.StatusNotFound, "管理员不存在")
		return
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "更新管理员状态", id, req.Status, utils.GetClientIP(c))

	responseSuccess(c, nil)
}
