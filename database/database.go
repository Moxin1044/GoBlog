package database

import (
	"fmt"
	"log"

	"github.com/moxin/GoBlog/config"
	"github.com/moxin/GoBlog/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	if err := model.AutoMigrate(DB); err != nil {
		log.Printf("AutoMigrate error: %v", err)
		return err
	}

	// 初始化默认超级管理员
	initDefaultAdmin()

	return nil
}

func initDefaultAdmin() {
	var count int64
	DB.Model(&model.Admin{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := model.Admin{
			Username: "admin",
			Email:    "admin@goblog.com",
			Password: string(hashedPassword),
			Role:     "superadmin",
			Status:   "active",
		}
		DB.Create(&admin)
		log.Println("Default super admin created: admin / admin123")
	}
}
