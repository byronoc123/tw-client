package metrics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestsTotal counts the total number of requests
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "blockchain_client_requests_total",
			Help: "The total number of API requests",
		},
		[]string{"endpoint", "method", "status"},
	)

	// RequestDuration tracks the duration of requests
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "blockchain_client_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint", "method"},
	)

	// RPCRequestsTotal counts RPC requests to the blockchain
	RPCRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "blockchain_client_rpc_requests_total",
			Help: "The total number of RPC requests to the blockchain",
		},
		[]string{"method", "status"},
	)

	// RPCRequestDuration tracks the duration of RPC requests
	RPCRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "blockchain_client_rpc_request_duration_seconds",
			Help:    "RPC request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// BlockProcessingTime tracks the time to process a block
	BlockProcessingTime = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "blockchain_client_block_processing_seconds",
			Help:    "Time to process a block in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	// BlockchainHeight tracks the current height of the blockchain
	BlockchainHeight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "blockchain_client_blockchain_height",
			Help: "Current height of the blockchain",
		},
	)
)

// RecordAPIRequest records metrics for an API request
func RecordAPIRequest(endpoint, method, status string, duration time.Duration) {
	RequestsTotal.WithLabelValues(endpoint, method, status).Inc()
	RequestDuration.WithLabelValues(endpoint, method).Observe(duration.Seconds())
}

// RecordRPCRequest records metrics for an RPC request
func RecordRPCRequest(method, status string, duration time.Duration) {
	RPCRequestsTotal.WithLabelValues(method, status).Inc()
	RPCRequestDuration.WithLabelValues(method).Observe(duration.Seconds())
}

// RecordBlockProcessing records the time taken to process a block
func RecordBlockProcessing(duration time.Duration) {
	BlockProcessingTime.Observe(duration.Seconds())
}

// UpdateBlockchainHeight updates the gauge for blockchain height
func UpdateBlockchainHeight(height float64) {
	BlockchainHeight.Set(height)
}

// MetricsMiddleware returns a Gin middleware for collecting metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Record metrics
		status := http.StatusText(c.Writer.Status())
		duration := time.Since(start)
		RecordAPIRequest(path, method, status, duration)
	}
}

// RegisterMetricsEndpoint registers the Prometheus metrics endpoint
func RegisterMetricsEndpoint(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
