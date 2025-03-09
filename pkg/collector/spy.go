package collector

import (
	"time"

	"github.com/shirou/gopsutil/process"
)

// Spy collects information about running processes.
type Spy struct{}

// NewSpy returns a new Spy instance.
func NewSpy() *Spy {
	return &Spy{}
}

// Collect retrieves process details.
func (s *Spy) Collect() (interface{}, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	// Collect basic info for each process.
	results := make([]map[string]interface{}, 0, len(procs))
	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			continue
		}
		results = append(results, map[string]interface{}{
			"pid":  p.Pid,
			"name": name,
			"time": time.Now(),
		})
	}
	return results, nil
}
