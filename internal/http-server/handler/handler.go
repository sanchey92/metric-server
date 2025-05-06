// Package handler provides HTTP handlers for processing and managing metrics.
// It defines a Handler type that receives metrics in JSON format and pushes them
// to a buffer service for further processing.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sanchey92/metric-server/internal/models"
)

// BufferService defines the interface for services that can accept and process
// batches of metrics. This abstraction allows the handler to work with different
// metric processing implementations.
type BufferService interface {
	Push([]models.Metric)
}

// Handler provides HTTP handlers for metric processing operations.
type Handler struct {
	service BufferService
}

// New creates and returns a new Handler instance with the provided BufferService.
func New(service BufferService) *Handler {
	return &Handler{
		service: service,
	}
}

// HandleMetrics processes incoming HTTP requests containing metric data.
// It expects a JSON array of metrics in the request body and pushes them
// to the configured BufferService.
func (h *Handler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	var metrics []models.Metric

	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	h.service.Push(metrics)
	w.WriteHeader(http.StatusOK)
}
