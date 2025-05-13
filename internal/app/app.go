// Package app provides the core application management for the metric server.
package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/sanchey92/metric-server/internal/config"
	"github.com/sanchey92/metric-server/internal/flusher"
	"github.com/sanchey92/metric-server/internal/http-server/server"
	"github.com/sanchey92/metric-server/internal/storage"
)

// App is the main application struct that orchestrates the HTTP server,
// metrics flusher, and database storage components.
// It manages their lifecycle and handles graceful shutdown.
type App struct {
	server  *server.Server
	flusher *flusher.Flusher
	db      *storage.PostgresStorage
	errCh   chan error
}

// New creates and initializes a new App instance.
// It sets up the memory storage, database connection, HTTP server, and metrics flusher.
func New(ctx context.Context, cfg *config.Config) (*App, error) {
	memStorage := storage.NewMemStorage()

	db, err := storage.NewPostgresStorage(ctx, cfg.PgDSN)
	if err != nil {
		return nil, err
	}

	s, err := server.New(cfg, memStorage)
	if err != nil {
		return nil, err
	}

	interval := 1 * time.Minute

	f := flusher.New(interval, memStorage, db)

	return &App{
		server:  s,
		flusher: f,
		db:      db,
		errCh:   make(chan error, 2),
	}, nil
}

// Run starts the application components and manages their lifecycle.
// It handles graceful shutdown on receiving termination signals.
func (a *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Println("Starting HTTPServer on port :8080")
		if err := a.server.Run(); err != nil {
			a.errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	go func() {
		fmt.Println("Starting metrics flusher")
		if err := a.flusher.Run(ctx); err != nil {
			a.errCh <- fmt.Errorf("flusher error: %w", err)
		}
	}()

	select {
	case err := <-a.errCh:
		fmt.Printf("application error: %v", err)
		return err
	case <-ctx.Done():
		fmt.Println("application shutdown initiated")
	}

	return a.shutdown()
}

// shutdown performs the orderly shutdown of application components.
// It attempts to stop the HTTP server and close the database connection
// with a timeout to prevent hanging.
func (a *App) shutdown() error {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	if err := a.db.Close(); err != nil {
		return err
	}

	fmt.Printf("Application shutdown complete")
	return nil
}
