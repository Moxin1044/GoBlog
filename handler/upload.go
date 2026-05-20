package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/config"
)

// UploadImage 通用图片上传（校验格式、大小、保存到uploads目录）
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		responseError(c, "请选择图片文件")
		return
	}

	// 校验文件大小（默认最大5MB）
	maxSize := config.AppConfig.Upload.MaxSize * 1024 * 1024
	if maxSize == 0 {
		maxSize = 5 * 1024 * 1024
	}
	if file.Size > maxSize {
		responseError(c, fmt.Sprintf("文件大小不能超过%dMB", config.AppConfig.Upload.MaxSize))
		return
	}

	// 校验文件格式
	contentType := file.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/gif":     true,
		"image/webp":    true,
		"image/svg+xml": true,
	}
	if !allowedTypes[contentType] {
		responseError(c, "仅支持 JPG、PNG、GIF、WebP、SVG 格式")
		return
	}

	// 生成文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	uploadDir := filepath.Join(config.AppConfig.Upload.Path, "images")

	// 确保目录存在
	os.MkdirAll(uploadDir, 0755)

	savePath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "图片保存失败")
		return
	}

	imageURL := "/uploads/images/" + filename
	responseSuccess(c, gin.H{"url": imageURL})
}
