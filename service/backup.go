package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/moxin/GoBlog/config"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
)

func CreateDatabaseBackup(backupType string) (*model.Backup, error) {
	cfg := config.AppConfig.Database
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("goblog_backup_%s.sql", timestamp)
	backupDir := "./backups"
	os.MkdirAll(backupDir, 0755)
	filePath := filepath.Join(backupDir, filename)

	// mysqldump命令
	cmd := exec.Command("mysqldump",
		"-h", cfg.Host,
		"-P", fmt.Sprintf("%d", cfg.Port),
		"-u", cfg.User,
		fmt.Sprintf("-p%s", cfg.Password),
		cfg.DBName,
		"-r", filePath,
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("backup failed: %v", err)
	}

	// 获取文件大小
	fileInfo, _ := os.Stat(filePath)
	size := int64(0)
	if fileInfo != nil {
		size = fileInfo.Size()
	}

	backup := &model.Backup{
		Filename: filename,
		Size:     size,
		Type:     backupType,
	}

	if err := database.DB.Create(backup).Error; err != nil {
		return nil, err
	}

	return backup, nil
}
