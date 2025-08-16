package internal

import (
	"sync"
	"time"
)

// Metrics holds basic operational metrics
type Metrics struct {
	mu               sync.RWMutex
	RequestCount     int64
	ErrorCount       int64
	LastRequestTime  time.Time
	AverageRequestMs float64
	requestDurations []time.Duration
	maxDurations     int
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		maxDurations: 100, // Keep last 100 request durations for average calculation
	}
}

// RecordRequest records a successful request with its duration
func (m *Metrics) RecordRequest(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.RequestCount++
	m.LastRequestTime = time.Now()

	// Add duration to sliding window
	m.requestDurations = append(m.requestDurations, duration)
	if len(m.requestDurations) > m.maxDurations {
		m.requestDurations = m.requestDurations[1:]
	}

	// Calculate new average
	var total time.Duration
	for _, d := range m.requestDurations {
		total += d
	}
	m.AverageRequestMs = float64(total.Nanoseconds()) / float64(len(m.requestDurations)) / 1e6
}

// RecordError records an error occurrence
func (m *Metrics) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorCount++
}

// GetStats returns a snapshot of current metrics
func (m *Metrics) GetStats() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return MetricsSnapshot{
		RequestCount:     m.RequestCount,
		ErrorCount:       m.ErrorCount,
		LastRequestTime:  m.LastRequestTime,
		AverageRequestMs: m.AverageRequestMs,
	}
}

// MetricsSnapshot represents a point-in-time view of metrics
type MetricsSnapshot struct {
	RequestCount     int64     `json:"request_count"`
	ErrorCount       int64     `json:"error_count"`
	LastRequestTime  time.Time `json:"last_request_time"`
	AverageRequestMs float64   `json:"average_request_ms"`
}
