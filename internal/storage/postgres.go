// Package storage provides metric storage implementations.
package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStorage implements metric storage using PostgreSQL as the backend.
// It maintains a connection pool for efficient database access and provides
// transactional guarantees for metric updates.
type PostgresStorage struct {
	pool *pgxpool.Pool
}

// NewPostgresStorage creates and initializes a new PostgreSQL-backed storage.
// It establishes a connection pool with the database and verifies connectivity.
func NewPostgresStorage(ctx context.Context, dsn string) (*PostgresStorage, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context canceled before connecting to postgres")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config")
	}

	config.MaxConns = 10
	config.MinConns = 1
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool")
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres db")
	}

	return &PostgresStorage{
		pool: pool,
	}, nil
}

// Close gracefully shuts down the storage by closing all database connections.
// It should be called when the storage is no longer needed to release resources.
func (s *PostgresStorage) Close() error {
	if s.pool != nil {
		s.pool.Close()
		fmt.Println("close connection to postgres")
	}

	return nil
}

// Save persists a batch of metrics to PostgreSQL using a transaction.
// It performs atomic upsert operations (insert new or update existing metrics).
func (s *PostgresStorage) Save(ctx context.Context, data map[string]float64) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to init transaction")
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			fmt.Println("rollback error")
		}
	}()

	query := `INSERT INTO metrics (name, value)
			 VALUES ($1, $2)
			 ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value`

	for name, value := range data {
		_, err = tx.Exec(ctx, query, name, value)
		if err != nil {
			return fmt.Errorf("exec tx error")
		}
	}

	return tx.Commit(ctx)
}
