package service

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type MonitorData struct {
	CPU     CPUData     `json:"cpu"`
	Memory  MemoryData  `json:"memory"`
	Disk    DiskData    `json:"disk"`
	Network NetworkData `json:"network"`
}

type CPUData struct {
	Usage float64 `json:"usage"` // 百分比
	Cores int     `json:"cores"`
}

type MemoryData struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Available    uint64  `json:"available"`
	UsagePercent float64 `json:"usage_percent"`
}

type DiskData struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	UsagePercent float64 `json:"usage_percent"`
}

type NetworkData struct {
	BytesSent    uint64  `json:"bytes_sent"`
	BytesRecv    uint64  `json:"bytes_recv"`
	UploadRate   float64 `json:"upload_rate"`   // bytes/s
	DownloadRate float64 `json:"download_rate"` // bytes/s
}

var lastNetStats *net.IOCountersStat
var lastNetTime time.Time

func GetMonitorData() (*MonitorData, error) {
	data := &MonitorData{}

	// CPU
	cpuPercents, _ := cpu.Percent(time.Second, false)
	cpuCores, _ := cpu.Counts(true)
	if len(cpuPercents) > 0 {
		data.CPU = CPUData{Usage: cpuPercents[0], Cores: cpuCores}
	}

	// Memory
	memInfo, _ := mem.VirtualMemory()
	if memInfo != nil {
		data.Memory = MemoryData{
			Total: memInfo.Total, Used: memInfo.Used,
			Available: memInfo.Available, UsagePercent: memInfo.UsedPercent,
		}
	}

	// Disk
	diskInfo, _ := disk.Usage("/")
	if diskInfo != nil {
		data.Disk = DiskData{
			Total: diskInfo.Total, Used: diskInfo.Used,
			Free: diskInfo.Free, UsagePercent: diskInfo.UsedPercent,
		}
	}

	// Network
	netInfos, _ := net.IOCounters(false)
	if len(netInfos) > 0 {
		currentStats := &netInfos[0]
		now := time.Now()
		if lastNetStats != nil {
			duration := now.Sub(lastNetTime).Seconds()
			if duration > 0 {
				data.Network = NetworkData{
					BytesSent: currentStats.BytesSent, BytesRecv: currentStats.BytesRecv,
					UploadRate:   float64(currentStats.BytesSent-lastNetStats.BytesSent) / duration,
					DownloadRate: float64(currentStats.BytesRecv-lastNetStats.BytesRecv) / duration,
				}
			}
		} else {
			data.Network = NetworkData{
				BytesSent: currentStats.BytesSent, BytesRecv: currentStats.BytesRecv,
			}
		}
		lastNetStats = currentStats
		lastNetTime = now
	}

	return data, nil
}
