// Package router provides HTTP route configuration and middleware setup
// for the metric server application. It uses the chi router to define
// API endpoints and apply middleware.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sanchey92/metric-server/internal/http-server/middleware"
)

// MetricHandler defines the interface for handlers that process metric-related
// HTTP requests. This abstraction allows for flexible handler implementations
// while maintaining a consistent router interface.
type MetricHandler interface {
	HandleMetrics(w http.ResponseWriter, r *http.Request)
}

// New creates and configures a new chi router instance with:
// - Gzip middleware for request/response compression
// - POST /update route for metric submissions
func New(handler MetricHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)
	r.Post("/update", handler.HandleMetrics)

	return r
}
