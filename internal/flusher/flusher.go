// Package flusher provides periodic metric data synchronization between
// in-memory storage and persistent database storage. It implements a
// ticker-based flushing mechanism with configurable intervals.
package flusher

import (
	"context"
	"fmt"
	"time"
)

// MemStorage defines the interface for in-memory metric storage that
// can provide snapshots of current metrics.
type MemStorage interface {
	Snapshot() map[string]float64
}

// PostgresStorage defines the interface for persistent metric storage
// that can save batches of metrics.
type PostgresStorage interface {
	Save(ctx context.Context, data map[string]float64) error
}

// Flusher implements periodic synchronization of metrics from memory to database.
// It runs at configured intervals until the context is canceled.
type Flusher struct {
	interval   time.Duration
	memStorage MemStorage
	db         PostgresStorage
}

// New creates a new Flusher instance with the specified configuration.
func New(interval time.Duration, storage MemStorage, db PostgresStorage) *Flusher {
	return &Flusher{
		interval:   interval,
		memStorage: storage,
		db:         db,
	}
}

// Run starts the periodic flushing process and blocks until the context is canceled.
// It performs flushes both on interval ticks and during graceful shutdown.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: Any error encountered during final flush on shutdown
//
// Operation:
// 1. Creates a ticker that triggers at the configured interval
// 2. On each tick:
//   - Takes snapshot of in-memory metrics
//   - Persists to database
//
// 3. On context cancellation:
//   - Performs one final flush
//   - Returns any flush error
func (f *Flusher) Run(ctx context.Context) error {
	ticker := time.NewTicker(f.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return f.flush(ctx)
		case <-ticker.C:
			if err := f.flush(ctx); err != nil {
				fmt.Println("failed to save metrics: %w\n", err)
			}
		}
	}
}

// flush performs the actual synchronization of metrics from memory to database.
// It's called both periodically and during shutdown.
func (f *Flusher) flush(ctx context.Context) error {
	snapshot := f.memStorage.Snapshot()

	if len(snapshot) == 0 {
		return nil
	}

	if err := f.db.Save(ctx, snapshot); err != nil {
		return fmt.Errorf("failed to save metrics: %w", err)
	}

	fmt.Printf("Successfully flushed %d metrics\n", len(snapshot))
	return nil
}
