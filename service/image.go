package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/moxin/GoBlog/config"
)

func SaveUploadedFile(fileHeader io.Reader, originalName string) (string, error) {
	ext := strings.ToLower(filepath.Ext(originalName))

	// 检查文件类型
	allowed := false
	for _, t := range config.AppConfig.Upload.AllowTypes {
		if ext == t {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 生成文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dateDir := time.Now().Format("20060102")
	dir := filepath.Join(config.AppConfig.Upload.Path, dateDir)
	os.MkdirAll(dir, 0755)

	fullPath := filepath.Join(dir, filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, fileHeader); err != nil {
		return "", err
	}

	// 返回相对URL路径
	return fmt.Sprintf("/uploads/%s/%s", dateDir, filename), nil
}

// 下载网络图片并转存到本地
func DownloadAndSaveImage(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: %s", resp.Status)
	}

	// 从URL或Content-Type推断扩展名
	ext := filepath.Ext(imageURL)
	if ext == "" {
		ext = ".jpg"
	}

	return SaveUploadedFile(resp.Body, "downloaded"+ext)
}
