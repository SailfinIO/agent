package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SailfinIO/agent/pkg/collector"
	"github.com/SailfinIO/agent/pkg/config"
	"github.com/SailfinIO/agent/pkg/server"
	"github.com/SailfinIO/agent/pkg/storage"
	"github.com/SailfinIO/agent/pkg/utils"
)

// Collector defines the interface for metric collectors.
type Collector interface {
	Collect() (interface{}, error)
}

// Agent ties together configuration, collectors, the HTTP server, storage, and logging.
type Agent struct {
	cfg        *config.Config
	httpServer *server.HTTPServer
	collectors []Collector
	storage    storage.Storage
	logger     utils.Logger
}

// NewAgent creates a new Agent instance.
func NewAgent(cfg *config.Config) (*Agent, error) {
	// Initialize collectors, including memory, process spy, and CPU.
	collectors := []Collector{
		collector.NewMemoryCollector(),
		collector.NewSpy(),
		collector.NewCPUCollector(),
	}

	// Create an HTTP mux that will serve the /metrics endpoint.
	mux := http.NewServeMux()
	a := &Agent{
		cfg:        cfg,
		collectors: collectors,
		storage:    storage.NewInMemoryStorage(),
		logger:     utils.New().WithContext("agent"),
	}
	mux.HandleFunc("/metrics", a.handleMetrics)

	srv := server.NewHTTPServer(cfg.ServerAddress, mux)
	a.httpServer = srv
	return a, nil
}

// Start runs the agent, including periodic metric collection and the HTTP server.
func (a *Agent) Start() error {
	// Launch a goroutine that collects and stores snapshots every 30 seconds.
	go func() {
		for {
			metrics, err := aggregateMetrics(a.collectors)
			if err != nil {
				a.logger.Error(fmt.Sprintf("Error collecting metrics: %v", err))
			} else {
				snapshot := storage.Snapshot{
					Timestamp: time.Now(),
					Metrics:   metrics,
				}
				if err := a.storage.Save(snapshot); err != nil {
					a.logger.Error(fmt.Sprintf("Error saving snapshot: %v", err))
				} else {
					a.logger.Info(fmt.Sprintf("Saved snapshot at %v", snapshot.Timestamp))
				}
			}
			time.Sleep(30 * time.Second)
		}
	}()
	return a.httpServer.Start()
}

// CollectMetrics returns a fresh aggregated metric snapshot (without storage).
func (a *Agent) CollectMetrics() ([]byte, error) {
	metrics, err := aggregateMetrics(a.collectors)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(metrics, "", "  ")
}

// aggregateMetrics collects metrics from all collectors and returns the data as a map.
func aggregateMetrics(collectors []Collector) (map[string]interface{}, error) {
	aggregated := make(map[string]interface{})
	for _, c := range collectors {
		switch v := c.(type) {
		case *collector.MemoryCollector:
			data, err := v.Collect()
			if err != nil {
				return nil, err
			}
			aggregated["memory"] = data
		case *collector.Spy:
			data, err := v.Collect()
			if err != nil {
				return nil, err
			}
			aggregated["processes"] = data
		case *collector.CPUCollector:
			data, err := v.Collect()
			if err != nil {
				return nil, err
			}
			aggregated["cpu"] = data
		default:
			aggregated["unknown"] = "collector type not recognized"
		}
	}
	return aggregated, nil
}

// handleMetrics serves HTTP requests to /metrics.
// It supports query parameters:
//   - limit: number of latest snapshots to return.
//   - from and to: Unix timestamps to define a time window.
func (a *Agent) handleMetrics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// If both "from" and "to" are provided, use a time range query.
	if fromStr, toStr := query.Get("from"), query.Get("to"); fromStr != "" && toStr != "" {
		fromUnix, err1 := parseUnix(fromStr)
		toUnix, err2 := parseUnix(toStr)
		if err1 != nil || err2 != nil {
			http.Error(w, "Invalid from/to parameters", http.StatusBadRequest)
			return
		}
		from := time.Unix(fromUnix, 0)
		to := time.Unix(toUnix, 0)
		snaps, err := a.GetSnapshotsByTime(from, to)
		if err != nil {
			http.Error(w, "Error querying snapshots", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(snaps)
		return
	}

	// If "limit" is provided, return the latest N snapshots.
	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconvAtoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		snaps, err := a.GetSnapshotsByLimit(limit)
		if err != nil {
			http.Error(w, "Error retrieving snapshots", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(snaps)
		return
	}

	// Default: return the latest snapshot.
	snap, err := a.GetLatestSnapshot()
	if err != nil {
		http.Error(w, "Error retrieving snapshots", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(snap)
}

// GetSnapshotsByTime returns snapshots collected between the given times.
func (a *Agent) GetSnapshotsByTime(from, to time.Time) ([]storage.Snapshot, error) {
	return a.storage.Query(from, to)
}

// GetSnapshotsByLimit returns the latest 'limit' snapshots.
func (a *Agent) GetSnapshotsByLimit(limit int) ([]storage.Snapshot, error) {
	all, err := a.storage.GetAll()
	if err != nil {
		return nil, err
	}
	if limit < len(all) {
		return all[len(all)-limit:], nil
	}
	return all, nil
}

// GetLatestSnapshot returns the most recent snapshot.
func (a *Agent) GetLatestSnapshot() (storage.Snapshot, error) {
	all, err := a.storage.GetAll()
	if err != nil {
		return storage.Snapshot{}, err
	}
	if len(all) == 0 {
		return storage.Snapshot{}, fmt.Errorf("no snapshots available")
	}
	return all[len(all)-1], nil
}

// parseUnix converts a string to an int64 Unix timestamp.
func parseUnix(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// strconvAtoi is a helper to convert string to int.
func strconvAtoi(s string) (int, error) {
	return strconv.Atoi(s)
}
