package service

import (
	"fmt"
	"log"

	"github.com/moxin/GoBlog/config"
	"github.com/robfig/cron/v3"
)

func StartCronJobs() {
	c := cron.New()

	// 服务器监控检查（每5分钟）
	if config.AppConfig.Monitor.Enabled {
		c.AddFunc("*/5 * * * *", func() {
			checkServerHealth()
		})
	}

	// 自动备份（每天凌晨3点）
	c.AddFunc("0 3 * * *", func() {
		if _, err := CreateDatabaseBackup("auto"); err != nil {
			log.Printf("Auto backup failed: %v", err)
		} else {
			log.Println("Auto backup completed")
		}
	})

	c.Start()
	log.Println("Cron jobs started")
}

func checkServerHealth() {
	data, err := GetMonitorData()
	if err != nil {
		return
	}

	// CPU超过90%告警
	if data.CPU.Usage > 90 {
		AlertServerAbnormal("CPU", fmt.Sprintf("CPU占用率 %.1f%%", data.CPU.Usage))
	}

	// 内存超过90%告警
	if data.Memory.UsagePercent > 90 {
		AlertServerAbnormal("内存", fmt.Sprintf("内存占用率 %.1f%%", data.Memory.UsagePercent))
	}
}
