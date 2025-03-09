package collector

import (
	"github.com/shirou/gopsutil/mem"
)

// MemoryCollector gathers memory usage metrics.
type MemoryCollector struct{}

// NewMemoryCollector creates a new MemoryCollector.
func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{}
}

// Collect retrieves virtual memory statistics.
func (m *MemoryCollector) Collect() (interface{}, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"total":       vmStat.Total,
		"available":   vmStat.Available,
		"used":        vmStat.Used,
		"usedPercent": vmStat.UsedPercent,
	}, nil
}
