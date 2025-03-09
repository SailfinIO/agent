package storage

import "time"

// Snapshot represents a single set of collected metrics.
type Snapshot struct {
	Timestamp time.Time              `json:"timestamp"`
	Metrics   map[string]interface{} `json:"metrics"`
}

// Storage defines the interface for storing snapshots.
type Storage interface {
	Save(snapshot Snapshot) error
	GetAll() ([]Snapshot, error)
	Query(from, to time.Time) ([]Snapshot, error)
}

// InMemoryStorage is a simple storage backend that keeps snapshots in memory.
type InMemoryStorage struct {
	snapshots []Snapshot
}

// NewInMemoryStorage returns a new InMemoryStorage instance.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		snapshots: []Snapshot{},
	}
}

// Save appends a new snapshot.
func (s *InMemoryStorage) Save(snapshot Snapshot) error {
	s.snapshots = append(s.snapshots, snapshot)
	return nil
}

// GetAll returns all snapshots.
func (s *InMemoryStorage) GetAll() ([]Snapshot, error) {
	return s.snapshots, nil
}

// Query returns snapshots between two timestamps.
func (s *InMemoryStorage) Query(from, to time.Time) ([]Snapshot, error) {
	var result []Snapshot
	for _, snap := range s.snapshots {
		if snap.Timestamp.After(from) && snap.Timestamp.Before(to) {
			result = append(result, snap)
		}
	}
	return result, nil
}
