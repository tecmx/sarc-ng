package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Define metrics as package-level variables to ensure they are registered with Prometheus
var (
	// RequestsTotal tracks the total number of HTTP requests
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// RequestDuration tracks the duration of HTTP requests
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// RequestsInProgress tracks the number of in-progress HTTP requests
	requestsInProgress = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_progress",
			Help: "Number of in-progress HTTP requests",
		},
		[]string{"method", "path"},
	)
)

// Metrics returns a middleware that collects Prometheus metrics
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics endpoint itself to avoid recursion
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// Get path and method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
			if path == "" {
				path = "unknown"
			}
		}
		method := c.Request.Method

		// Track in-progress requests
		requestsInProgress.WithLabelValues(method, path).Inc()
		defer requestsInProgress.WithLabelValues(method, path).Dec()

		// Track request duration
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime).Seconds()

		// Get status code
		status := strconv.Itoa(c.Writer.Status())

		// Record metrics
		requestsTotal.WithLabelValues(method, path, status).Inc()
		requestDuration.WithLabelValues(method, path, status).Observe(duration)
	}
}
 