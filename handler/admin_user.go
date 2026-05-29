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

// AdminCreateUser 管理员创建用户
func AdminCreateUser(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	var count int64
	database.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		responseError(c, "用户名已存在")
		return
	}

	database.DB.Model(&model.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		responseError(c, "邮箱已被注册")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: nickname,
		Phone:    req.Phone,
		Role:     "user",
		Status:   "active",
	}
	if err := database.DB.Create(&user).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "创建用户失败")
		return
	}

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "创建用户", req.Username, "成功", utils.GetClientIP(c))

	responseSuccess(c, gin.H{"id": user.ID})
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
