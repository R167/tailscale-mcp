package internal

import (
	"testing"
	"time"
)

func TestMetrics_RecordRequest(t *testing.T) {
	metrics := NewMetrics()

	// Record a request
	duration := 100 * time.Millisecond
	metrics.RecordRequest(duration)

	stats := metrics.GetStats()
	if stats.RequestCount != 1 {
		t.Errorf("Expected request count 1, got %d", stats.RequestCount)
	}

	if stats.ErrorCount != 0 {
		t.Errorf("Expected error count 0, got %d", stats.ErrorCount)
	}

	if stats.AverageRequestMs <= 0 {
		t.Errorf("Expected positive average request time, got %f", stats.AverageRequestMs)
	}

	// Approximate check for average (should be around 100ms)
	expectedMs := 100.0
	tolerance := 5.0 // 5ms tolerance
	if stats.AverageRequestMs < expectedMs-tolerance || stats.AverageRequestMs > expectedMs+tolerance {
		t.Errorf("Expected average around %fms, got %fms", expectedMs, stats.AverageRequestMs)
	}
}

func TestMetrics_RecordError(t *testing.T) {
	metrics := NewMetrics()

	// Record an error
	metrics.RecordError()

	stats := metrics.GetStats()
	if stats.RequestCount != 0 {
		t.Errorf("Expected request count 0, got %d", stats.RequestCount)
	}

	if stats.ErrorCount != 1 {
		t.Errorf("Expected error count 1, got %d", stats.ErrorCount)
	}
}

func TestMetrics_MultipleRequests(t *testing.T) {
	metrics := NewMetrics()

	// Record multiple requests with different durations
	durations := []time.Duration{
		50 * time.Millisecond,
		100 * time.Millisecond,
		150 * time.Millisecond,
	}

	for _, duration := range durations {
		metrics.RecordRequest(duration)
	}

	stats := metrics.GetStats()
	if stats.RequestCount != 3 {
		t.Errorf("Expected request count 3, got %d", stats.RequestCount)
	}

	// Average should be around 100ms
	expectedMs := 100.0
	tolerance := 10.0 // 10ms tolerance
	if stats.AverageRequestMs < expectedMs-tolerance || stats.AverageRequestMs > expectedMs+tolerance {
		t.Errorf("Expected average around %fms, got %fms", expectedMs, stats.AverageRequestMs)
	}
}

func TestMetrics_SlidingWindow(t *testing.T) {
	metrics := NewMetrics()
	metrics.maxDurations = 3 // Small window for testing

	// Record more requests than window size
	durations := []time.Duration{
		10 * time.Millisecond,  // This should be dropped
		20 * time.Millisecond,  // This should be dropped
		100 * time.Millisecond, // Keep
		200 * time.Millisecond, // Keep
		300 * time.Millisecond, // Keep
	}

	for _, duration := range durations {
		metrics.RecordRequest(duration)
	}

	stats := metrics.GetStats()
	if stats.RequestCount != 5 {
		t.Errorf("Expected request count 5, got %d", stats.RequestCount)
	}

	// Average should be around 200ms (100+200+300)/3
	expectedMs := 200.0
	tolerance := 10.0
	if stats.AverageRequestMs < expectedMs-tolerance || stats.AverageRequestMs > expectedMs+tolerance {
		t.Errorf("Expected average around %fms, got %fms", expectedMs, stats.AverageRequestMs)
	}
}

func TestMetrics_Concurrent(t *testing.T) {
	metrics := NewMetrics()

	// Test concurrent access
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			metrics.RecordRequest(time.Millisecond)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 50; i++ {
			metrics.RecordError()
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	stats := metrics.GetStats()
	if stats.RequestCount != 100 {
		t.Errorf("Expected request count 100, got %d", stats.RequestCount)
	}

	if stats.ErrorCount != 50 {
		t.Errorf("Expected error count 50, got %d", stats.ErrorCount)
	}
}
