package collector

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// CPUCollector collects CPU usage metrics.
type CPUCollector struct{}

// NewCPUCollector returns a new CPUCollector instance.
func NewCPUCollector() *CPUCollector {
	return &CPUCollector{}
}

// Collect retrieves CPU usage percentages.
// It waits for one second to sample CPU usage.
func (c *CPUCollector) Collect() (interface{}, error) {
	// Get overall CPU usage percentage over 1 second.
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// You can expand this to include per-core data if needed.
	return map[string]interface{}{
		"usage": percentages, // slice of percentages; overall if percpu=false
	}, nil
}
