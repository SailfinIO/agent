// pkg/collector/system_stats.go

package collector

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// SystemStatsCollector collects detailed system stats.
type SystemStatsCollector struct{}

// NewSystemStatsCollector returns a new instance.
func NewSystemStatsCollector() *SystemStatsCollector {
	return &SystemStatsCollector{}
}

// Collect gathers system statistics and returns a map.
func (s *SystemStatsCollector) Collect() (interface{}, error) {
	// Get host information.
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	// Get load averages.
	loadAvg, err := load.Avg()
	if err != nil {
		return nil, err
	}

	// Get memory usage.
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// Get swap usage.
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	// Get disk usage for root (you might iterate over all partitions as needed).
	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	// Get CPU info.
	cpuInfos, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	// Get network interfaces.
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// Build a slice of CPU details.
	cpus := make([]map[string]interface{}, 0, len(cpuInfos))
	for _, info := range cpuInfos {
		cpus = append(cpus, map[string]interface{}{
			"model": info.ModelName,
			"cores": info.Cores,
			"speed": info.Mhz / 1000.0, // Convert MHz to GHz.
		})
	}

	// Construct the system stats map.
	stats := map[string]interface{}{
		"uptime":   hostInfo.Uptime,
		"hostname": hostInfo.Hostname,
		"platform": hostInfo.Platform,
		"arch":     runtime.GOARCH,
		"loadAvg":  []float64{loadAvg.Load1, loadAvg.Load5, loadAvg.Load15},
		"memoryUsage": map[string]interface{}{
			"totalMem":   vmStat.Total,
			"usedMem":    vmStat.Used,
			"freeMem":    vmStat.Available,
			"percentage": vmStat.UsedPercent,
		},
		"swapInfo": map[string]interface{}{
			"totalSwap":    swapStat.Total,
			"freeSwap":     swapStat.Free,
			"usedSwap":     swapStat.Used,
			"usagePercent": swapStat.UsedPercent,
		},
		"diskUsage": map[string]interface{}{
			"filesystem":  diskStat.Fstype,
			"total":       diskStat.Total,
			"used":        diskStat.Used,
			"free":        diskStat.Free,
			"usedPercent": diskStat.UsedPercent,
			"mount":       diskStat.Path,
		},
		"cpus":              cpus,
		"networkInterfaces": interfaces, // This will include a lot of data; you can filter fields as needed.
		"timestamp":         time.Now().Unix(),
	}

	return stats, nil
}
