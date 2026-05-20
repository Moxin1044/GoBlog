package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// GetDashboard 博客数据统计（总访问量、总文章数、总字数）
func GetDashboard(c *gin.Context) {
	var articleCount int64
	var totalViews int64
	var totalWords int64

	database.DB.Model(&model.Article{}).Where("status = ?", "published").Count(&articleCount)

	type Result struct {
		Total int64
	}
	var viewResult Result
	database.DB.Model(&model.Article{}).Select("COALESCE(SUM(view_count), 0) as total").Where("status = ?", "published").Scan(&viewResult)
	totalViews = viewResult.Total

	var wordResult Result
	database.DB.Model(&model.Article{}).Select("COALESCE(SUM(CHAR_LENGTH(content)), 0) as total").Where("status = ?", "published").Scan(&wordResult)
	totalWords = wordResult.Total

	var commentCount int64
	database.DB.Model(&model.Comment{}).Count(&commentCount)

	var userCount int64
	database.DB.Model(&model.User{}).Count(&userCount)

	responseSuccess(c, gin.H{
		"article_count": articleCount,
		"total_views":   totalViews,
		"total_words":   totalWords,
		"comment_count": commentCount,
		"user_count":    userCount,
	})
}

// GetServerMonitor 服务器监控数据（CPU、内存、磁盘、网络）
func GetServerMonitor(c *gin.Context) {
	// CPU信息
	cpuPercent, _ := cpu.Percent(0, false)
	cpuCounts, _ := cpu.Counts(true)

	// 内存信息
	memInfo, _ := mem.VirtualMemory()

	// 磁盘信息
	diskInfo, _ := disk.Usage("/")

	// 系统负载
	loadInfo, _ := load.Avg()

	// 主机信息
	hostInfo, _ := host.Info()

	// 网络IO
	netInfo, _ := net.IOCounters(false)

	// 运行时间
	uptime := uint64(0)
	if hostInfo != nil {
		uptime = hostInfo.Uptime
	}

	var netStat map[string]interface{}
	if len(netInfo) > 0 {
		netStat = map[string]interface{}{
			"bytes_sent": netInfo[0].BytesSent,
			"bytes_recv": netInfo[0].BytesRecv,
		}
	}

	responseSuccess(c, gin.H{
		"cpu": map[string]interface{}{
			"percent":    cpuPercent,
			"core_count": cpuCounts,
		},
		"memory": map[string]interface{}{
			"total":        memInfo.Total,
			"used":         memInfo.Used,
			"available":    memInfo.Available,
			"used_percent": memInfo.UsedPercent,
		},
		"disk": map[string]interface{}{
			"total":        diskInfo.Total,
			"used":         diskInfo.Used,
			"free":         diskInfo.Free,
			"used_percent": diskInfo.UsedPercent,
		},
		"load": map[string]interface{}{
			"load1":  loadInfo.Load1,
			"load5":  loadInfo.Load5,
			"load15": loadInfo.Load15,
		},
		"host": map[string]interface{}{
			"hostname": hostInfo.Hostname,
			"os":       hostInfo.OS,
			"platform": hostInfo.Platform,
			"uptime":   uptime,
		},
		"network": netStat,
	})
}

// GetVisitMapData 访客地域数据
func GetVisitMapData(c *gin.Context) {
	var visitLogs []model.VisitLog
	database.DB.Select("location, COUNT(*) as count").Where("location != '' AND location IS NOT NULL").
		Group("location").Order("count DESC").Limit(50).Find(&visitLogs)

	type LocationData struct {
		Location string `json:"location"`
		Count    int64  `json:"count"`
	}

	var result []LocationData
	database.DB.Model(&model.VisitLog{}).
		Select("location, COUNT(*) as count").
		Where("location != '' AND location IS NOT NULL").
		Group("location").
		Order("count DESC").
		Limit(50).
		Scan(&result)

	responseSuccess(c, result)
}
