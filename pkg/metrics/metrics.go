package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// DatabaseOperations tracks database operations
	DatabaseOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "entity"},
	)

	// DatabaseErrors tracks database errors
	DatabaseErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_errors_total",
			Help: "Total number of database errors",
		},
		[]string{"operation", "entity"},
	)

	// DatabaseOperationDuration tracks the duration of database operations
	DatabaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_operation_duration_seconds",
			Help:    "Database operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "entity"},
	)

	// ApplicationInfo provides application version information
	ApplicationInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "application_info",
			Help: "Application information including version",
		},
		[]string{"version", "build_date", "commit"},
	)

	// CacheHits tracks cache hits
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache"},
	)

	// CacheMisses tracks cache misses
	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache"},
	)
)

// Initialize sets up initial metrics values
func Initialize(version, buildDate, commit string) {
	// Set application info
	ApplicationInfo.WithLabelValues(version, buildDate, commit).Set(1)
}
