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

// Register 用户注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	// 检查注册开关
	if getConfigValue("register_enabled") == "false" {
		responseErrorWithCode(c, http.StatusForbidden, "注册功能已关闭")
		return
	}

	// 验证码校验
	var vc model.VerificationCode
	if err := database.DB.Where("email = ? AND code = ? AND type = ? AND used = ? AND expired_at > ?",
		req.Email, req.Code, "register", false, time.Now()).First(&vc).Error; err != nil {
		responseError(c, "验证码无效或已过期")
		return
	}

	// 标记验证码已使用
	database.DB.Model(&vc).Update("used", true)

	// 用户名唯一性检查
	var count int64
	database.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		responseError(c, "用户名已存在")
		return
	}

	// 邮箱唯一性检查
	database.DB.Model(&model.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		responseError(c, "邮箱已被注册")
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Username,
		Role:     "user",
		Status:   "active",
	}
	if err := database.DB.Create(&user).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "注册失败")
		return
	}

	// 通知管理员
	go notifyAdmin("新用户注册", fmt.Sprintf("用户 %s (%s) 完成注册", req.Username, req.Email))

	responseSuccess(c, gin.H{"id": user.ID})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	// 支持用户名/邮箱登录
	var user model.User
	if err := database.DB.Where("username = ? OR email = ?", req.Account, req.Account).First(&user).Error; err != nil {
		responseError(c, "用户名或密码错误")
		return
	}

	// 检查用户状态
	if user.Status == "disabled" {
		responseErrorWithCode(c, http.StatusForbidden, "账号已被禁用")
		return
	}

	// 密码校验
	if !utils.CheckPassword(req.Password, user.Password) {
		responseError(c, "用户名或密码错误")
		return
	}

	// 生成JWT
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, "user")
	if err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "生成token失败")
		return
	}

	// 记录登录日志
	ip := utils.GetClientIP(c)
	database.DB.Create(&model.VisitLog{
		IP:        ip,
		Path:      "/api/auth/login",
		UserAgent: utils.GetUserAgent(c),
	})

	responseSuccess(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	})
}

// SendVerifyCode 发送验证码
func SendVerifyCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Type  string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	// 生成6位随机码
	code := generateVerifyCode()

	// 存入VerificationCode表
	vc := model.VerificationCode{
		Email:     req.Email,
		Code:      code,
		Type:      req.Type,
		ExpiredAt: time.Now().Add(5 * time.Minute),
		Used:      false,
	}
	if err := database.DB.Create(&vc).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "验证码创建失败")
		return
	}

	// 发送邮件
	if err := utils.SendVerificationCode(req.Email, code); err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "邮件发送失败: "+err.Error())
		return
	}

	responseSuccess(c, nil)
}

// generateVerifyCode 生成6位随机验证码
func generateVerifyCode() string {
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}
