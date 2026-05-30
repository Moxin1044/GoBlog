package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/config"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// GetUserInfo 获取用户信息（支持普通用户和管理员）
func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 先尝试从 users 表查找
	var user model.User
	if err := database.DB.First(&user, userID).Error; err == nil {
		responseSuccess(c, user)
		return
	}

	// 如果是管理员，从 admins 表查找
	adminID, exists := c.Get("admin_id")
	if exists {
		var admin model.Admin
		if err := database.DB.First(&admin, adminID).Error; err == nil {
			responseSuccess(c, gin.H{
				"id":         admin.ID,
				"username":   admin.Username,
				"email":      admin.Email,
				"nickname":   admin.Username,
				"avatar":     "",
				"role":       admin.Role,
				"is_admin":   true,
			})
			return
		}
	}

	responseErrorWithCode(c, http.StatusNotFound, "用户不存在")
}

// UpdateUserInfo 更新用户信息（昵称、邮箱）
func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email" binding:"omitempty,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		var count int64
		database.DB.Model(&model.User{}).Where("email = ? AND id != ?", req.Email, userID).Count(&count)
		if count > 0 {
			responseError(c, "邮箱已被其他用户使用")
			return
		}
		updates["email"] = req.Email
	}

	if len(updates) == 0 {
		responseError(c, "没有需要更新的字段")
		return
	}

	if err := database.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "更新失败")
		return
	}

	responseSuccess(c, nil)
}

// ChangePassword 修改密码（支持普通用户和管理员）
func ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	// 先尝试从 users 表查找
	var user model.User
	if err := database.DB.First(&user, userID).Error; err == nil {
		if !utils.CheckPassword(req.OldPassword, user.Password) {
			responseError(c, "原密码错误")
			return
		}

		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			responseErrorWithCode(c, http.StatusInternalServerError, "密码加密失败")
			return
		}

		database.DB.Model(&user).Update("password", hashedPassword)
		responseSuccess(c, nil)
		return
	}

	// 如果是管理员，从 admins 表查找
	adminID, exists := c.Get("admin_id")
	if exists {
		var admin model.Admin
		if err := database.DB.First(&admin, adminID).Error; err == nil {
			if !utils.CheckPassword(req.OldPassword, admin.Password) {
				responseError(c, "原密码错误")
				return
			}

			hashedPassword, err := utils.HashPassword(req.NewPassword)
			if err != nil {
				responseErrorWithCode(c, http.StatusInternalServerError, "密码加密失败")
				return
			}

			database.DB.Model(&admin).Update("password", hashedPassword)
			responseSuccess(c, nil)
			return
		}
	}

	responseErrorWithCode(c, http.StatusNotFound, "用户不存在")
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("avatar")
	if err != nil {
		responseError(c, "请选择头像文件")
		return
	}

	// 校验文件大小（最大2MB）
	if file.Size > 2*1024*1024 {
		responseError(c, "头像文件大小不能超过2MB")
		return
	}

	// 校验文件格式
	contentType := file.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		responseError(c, "仅支持 JPG、PNG、GIF、WebP 格式")
		return
	}

	// 保存文件
	filename := fmt.Sprintf("avatar_%d_%s", userID, file.Filename)
	savePath := filepath.Join(config.AppConfig.Upload.Path, "avatars", filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "头像保存失败")
		return
	}

	avatarURL := "/uploads/avatars/" + filename
	database.DB.Model(&model.User{}).Where("id = ?", userID).Update("avatar", avatarURL)

	responseSuccess(c, gin.H{"avatar": avatarURL})
}
