// Package handler provides HTTP handlers for processing and managing metrics.
// It defines a Handler type that receives metrics in JSON format and pushes them
// to a buffer service for further processing.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sanchey92/metric-server/internal/models"
)

// MemStorage defines an interface for storing metrics in memory.
// It provides methods for setting and retrieving metric values.
type MemStorage interface {
	Set(name string, value float64)
}

// Handler provides HTTP handlers for metric processing operations.
type Handler struct {
	storage MemStorage
}

// New creates and returns a new Handler instance with the provided BufferService.
func New(storage MemStorage) *Handler {
	return &Handler{
		storage: storage,
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

	for _, value := range metrics {
		h.storage.Set(value.Name, value.Value)
	}

	w.WriteHeader(http.StatusOK)
}
