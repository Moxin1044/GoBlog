package handler

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/config"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// GetBackupList 获取备份列表
func GetBackupList(c *gin.Context) {
	var backups []model.Backup
	if err := database.DB.Order("created_at DESC").Find(&backups).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取备份列表失败")
		return
	}
	responseSuccess(c, backups)
}

// CreateBackup 创建备份
func CreateBackup(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	cfg := config.AppConfig.Database
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("goblog_backup_%s.sql", timestamp)
	backupDir := filepath.Join(config.AppConfig.Upload.Path, "backups")

	// 确保备份目录存在
	os.MkdirAll(backupDir, 0755)

	backupPath := filepath.Join(backupDir, filename)

	// 使用mysqldump命令导出
	cmd := exec.Command("mysqldump",
		"-h", cfg.Host,
		"-P", fmt.Sprintf("%d", cfg.Port),
		"-u", cfg.User,
		fmt.Sprintf("-p%s", cfg.Password),
		cfg.DBName,
		"-r", backupPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "备份失败: "+string(output))
		return
	}

	// 获取文件大小
	var fileInfo os.FileInfo
	fileInfo, _ = os.Stat(backupPath)
	var fileSize int64
	if fileInfo != nil {
		fileSize = fileInfo.Size()
	}

	backup := model.Backup{
		Filename: filename,
		Size:     fileSize,
		Type:     "manual",
	}
	database.DB.Create(&backup)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "创建备份", filename, "成功", utils.GetClientIP(c))

	responseSuccess(c, backup)
}

// DownloadBackup 下载备份
func DownloadBackup(c *gin.Context) {
	id := c.Param("id")

	var backup model.Backup
	if err := database.DB.First(&backup, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "备份不存在")
		return
	}

	backupPath := filepath.Join(config.AppConfig.Upload.Path, "backups", backup.Filename)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		responseErrorWithCode(c, http.StatusNotFound, "备份文件不存在")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", backup.Filename))
	c.Header("Content-Type", "application/octet-stream")
	c.File(backupPath)
}

// DeleteBackup 删除备份
func DeleteBackup(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var backup model.Backup
	if err := database.DB.First(&backup, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "备份不存在")
		return
	}

	// 删除文件
	backupPath := filepath.Join(config.AppConfig.Upload.Path, "backups", backup.Filename)
	os.Remove(backupPath)

	database.DB.Delete(&backup)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "删除备份", backup.Filename, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}
