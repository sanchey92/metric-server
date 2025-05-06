// Package models defines the core data structures used throughout the metric-server application.
// These structures represent the domain objects and their JSON representations for API communication.
package models

// Metric represents a single measurement or data point collected by the system.
// It is used for both storage and API payloads, with JSON tags defining the serialization format.
type Metric struct {
	Name  string  `json:"name"`
	MType string  `json:"type"`
	Value float64 `json:"value"`
}
