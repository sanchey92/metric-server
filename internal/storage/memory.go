// Package storage provides metric storage implementations.
package storage

import "sync"

// MemStorage implements an in-memory thread-safe key-value store for metric data.
// It uses a read-write mutex to allow multiple concurrent readers or a single writer.
type MemStorage struct {
	mu   sync.RWMutex
	data map[string]float64
}

// NewMemStorage creates and returns a new initialized MemStorage instance.
// The returned storage is ready to use with an empty data map.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]float64),
	}
}

// Set stores a metric value with the given name in the storage.
// The operation is thread-safe and will overwrite any existing value.
func (s *MemStorage) Set(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[name] = value
}

// Snapshot creates and returns a thread-safe copy of all current metric values.
// The snapshot is a new map containing all key-value pairs at the time of calling.
func (s *MemStorage) Snapshot() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := make(map[string]float64, len(s.data))
	for key, value := range s.data {
		snapshot[key] = value
	}

	return snapshot
}
