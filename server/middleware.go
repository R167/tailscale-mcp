package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/R167/tailscale-mcp/internal"
)

// RequestIDKey is the context key for request IDs
type RequestIDKey struct{}

// LoggerKey is the context key for structured loggers
type LoggerKey struct{}

var globalMetrics = internal.NewMetrics()

// RequestMiddleware adds request context, correlation IDs, structured logging, and metrics collection
func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID
		requestID := generateRequestID()

		// Create structured logger with request context
		logger := slog.With(
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		// Add to request context
		ctx := context.WithValue(r.Context(), RequestIDKey{}, requestID)
		ctx = context.WithValue(ctx, LoggerKey{}, logger)
		r = r.WithContext(ctx)

		// Log request start
		logger.Info("Request started")

		// Wrap response writer to capture status
		wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log request completion and record metrics
		duration := time.Since(start)
		logger.Info("Request completed",
			"status_code", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
		)

		// Record metrics
		if wrapped.statusCode >= 400 {
			globalMetrics.RecordError()
		} else {
			globalMetrics.RecordRequest(duration)
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// generateRequestID creates a random request ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random fails
		return hex.EncodeToString([]byte(time.Now().String()))[:32]
	}
	return hex.EncodeToString(bytes)
}

// GetRequestID extracts the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey{}).(string); ok {
		return requestID
	}
	return ""
}

// GetLogger extracts the structured logger from context
func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(LoggerKey{}).(*slog.Logger); ok {
		return logger
	}
	// Fallback to default logger
	return slog.Default()
}

// GetMetrics returns the current metrics snapshot
func GetMetrics() internal.MetricsSnapshot {
	return globalMetrics.GetStats()
}
