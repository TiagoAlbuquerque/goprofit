package utils

import (
	"fmt"
	"sync"
	"time"
)

type Metric struct {
	TotalDuration time.Duration
	Count         int64
}

var (
	metrics      = make(map[string]*Metric)
	metricsMutex sync.RWMutex
)

// StartTimer begins a timer for the given sector name and returns a stop function
// It accumulates the duration to the total time for that sector
func StartTimer(name string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		metricsMutex.Lock()
		if _, ok := metrics[name]; !ok {
			metrics[name] = &Metric{}
		}
		metrics[name].TotalDuration += duration
		metrics[name].Count++
		metricsMutex.Unlock()
	}
}

// GetMetrics returns a map of all recorded metrics with friendly string values
func GetMetrics() map[string]interface{} {
	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	result := make(map[string]interface{})
	for k, v := range metrics {
		avg := time.Duration(0)
		if v.Count > 0 {
			avg = v.TotalDuration / time.Duration(v.Count)
		}
		result[k] = map[string]string{
			"total_time": v.TotalDuration.String(),
			"count":      fmt.Sprintf("%d", v.Count),
			"avg_time":   avg.String(),
		}
	}
	return result
}

// GetRawMetrics returns the raw metrics data
func GetRawMetrics() map[string]*Metric {
	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	// Return a copy to be safe
	copy := make(map[string]*Metric)
	for k, v := range metrics {
		copy[k] = &Metric{
			TotalDuration: v.TotalDuration,
			Count:         v.Count,
		}
	}
	return copy
}

// MeasureFunc is a helper to wrap a function call with a timer
func MeasureFunc(name string, f func()) {
	defer StartTimer(name)()
	f()
}
