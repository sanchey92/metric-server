// Package server provides the core HTTP server implementation for the metric-server application.
// It handles server initialization, configuration, and lifecycle management.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sanchey92/metric-server/internal/config"
	"github.com/sanchey92/metric-server/internal/http-server/handler"
	"github.com/sanchey92/metric-server/internal/http-server/router"
	"github.com/sanchey92/metric-server/internal/storage"
)

// Server represents the HTTP server for the metric service.
// It encapsulates the http.Server along with its configuration and dependencies.
type Server struct {
	srv *http.Server
}

// New creates and configures a new Server instance with all required dependencies.
// It initializes the storage, handlers, and router based on the provided configuration.
func New(cfg *config.Config, memStorage *storage.MemStorage) (*Server, error) {
	h := handler.New(memStorage)
	r := router.New(h)

	address := fmt.Sprintf("%s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port)

	srv := &http.Server{
		Addr:         address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &Server{srv: srv}, nil
}

// Run starts the HTTP server and begins accepting connections.
// It blocks until the server is shut down and returns any error encountered.
// The method gracefully handles http.ErrServerClosed as a normal shutdown case.
func (s *Server) Run() error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server, allowing in-flight requests to complete.
// It uses the provided context to control the shutdown timeout duration.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
