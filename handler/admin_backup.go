package handler

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/config"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
	"gorm.io/gorm"
)

func GetBackupList(c *gin.Context) {
	var backups []model.Backup
	if err := database.DB.Order("created_at DESC").Find(&backups).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取备份列表失败")
		return
	}
	responseSuccess(c, backups)
}

func CreateBackup(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("goblog_backup_%s.sql", timestamp)
	backupDir := filepath.Join(config.AppConfig.Upload.Path, "backups")
	os.MkdirAll(backupDir, 0755)
	backupPath := filepath.Join(backupDir, filename)

	f, err := os.Create(backupPath)
	if err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "备份失败: "+err.Error())
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(fmt.Sprintf("-- GoBlog Backup\n-- Date: %s\n-- Database: %s\n\n",
		time.Now().Format("2006-01-02 15:04:05"), config.AppConfig.Database.DBName))
	w.WriteString("SET NAMES utf8mb4;\nSET FOREIGN_KEY_CHECKS = 0;\n\n")

	tables := []interface{}{
		&model.User{}, &model.Article{}, &model.Category{}, &model.Tag{}, &model.ArticleTag{},
		&model.Comment{}, &model.Like{}, &model.Subscription{}, &model.SubscriptionCategory{},
		&model.AIModel{}, &model.UserAIConfig{}, &model.ChatMessage{}, &model.SiteConfig{},
		&model.Admin{}, &model.OperationLog{}, &model.VerificationCode{}, &model.VisitLog{},
		&model.Navigation{},
	}

	for _, table := range tables {
		stmt := &gorm.Statement{DB: database.DB}
		stmt.Parse(table)
		tableName := stmt.Schema.Table

		w.WriteString(fmt.Sprintf("-- ----------------------------\n-- Table structure for `%s`\n-- ----------------------------\n", tableName))
		w.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", tableName))

		var createSQL string
		var tblName string
		database.DB.Raw(fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)).Row().Scan(&tblName, &createSQL)
		w.WriteString(createSQL + ";\n\n")

		w.WriteString(fmt.Sprintf("-- ----------------------------\n-- Records of `%s`\n-- ----------------------------\n", tableName))

		rows, err := database.DB.Table(tableName).Rows()
		if err != nil {
			rows.Close()
			continue
		}

		columns, _ := rows.Columns()
		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range columns {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)

			colNames := ""
			colValues := ""
			for i, col := range columns {
				if i > 0 {
					colNames += ", "
					colValues += ", "
				}
				colNames += fmt.Sprintf("`%s`", col)
				val := values[i]
				if val == nil {
					colValues += "NULL"
				} else {
					switch v := val.(type) {
					case []byte:
						colValues += fmt.Sprintf("'%s'", escapeString(string(v)))
					case string:
						colValues += fmt.Sprintf("'%s'", escapeString(v))
					case time.Time:
						colValues += fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
					default:
						colValues += fmt.Sprintf("'%v'", v)
					}
				}
			}
			w.WriteString(fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s);\n", tableName, colNames, colValues))
		}
		rows.Close()
		w.WriteString("\n")
	}

	w.WriteString("SET FOREIGN_KEY_CHECKS = 1;\n")
	w.Flush()

	fileInfo, _ := os.Stat(backupPath)
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

func escapeString(s string) string {
	result := ""
	for _, c := range s {
		switch c {
		case '\'':
			result += "\\'"
		case '\\':
			result += "\\\\"
		case '\n':
			result += "\\n"
		case '\r':
			result += "\\r"
		case '\t':
			result += "\\t"
		case 0:
			result += "\\0"
		default:
			result += string(c)
		}
	}
	return result
}

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

func DeleteBackup(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")

	var backup model.Backup
	if err := database.DB.First(&backup, id).Error; err != nil {
		responseErrorWithCode(c, http.StatusNotFound, "备份不存在")
		return
	}

	backupPath := filepath.Join(config.AppConfig.Upload.Path, "backups", backup.Filename)
	os.Remove(backupPath)

	database.DB.Delete(&backup)

	adminName, _ := c.Get("admin_name")
	recordOperationLog(adminID, adminName.(string), "删除备份", backup.Filename, "成功", utils.GetClientIP(c))

	responseSuccess(c, nil)
}

func UpdateAutoBackupConfig(c *gin.Context) {
	var req struct {
		Enabled   bool   `json:"enabled"`
		Frequency string `json:"frequency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	setConfigValue("auto_backup_enabled", fmt.Sprintf("%v", req.Enabled))
	setConfigValue("auto_backup_frequency", req.Frequency)

	responseSuccess(c, nil)
}
