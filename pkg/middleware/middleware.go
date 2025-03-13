package middleware

import (
	"blockchain-client/pkg/logger"
	"blockchain-client/pkg/metrics"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"go.uber.org/zap"
)

// RateLimiterConfig defines configuration for rate limiting middleware
type RateLimiterConfig struct {
	Limit          int
	Period         time.Duration
	BlockDuration  time.Duration
	ClientIPHeader string
}

// DefaultRateLimiterConfig returns a default rate limiter configuration
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Limit:          100,
		Period:         time.Minute,
		BlockDuration:  time.Minute * 5,
		ClientIPHeader: "X-Forwarded-For",
	}
}

// Logger returns a middleware that logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Log after request is processed
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		logger.Info("HTTP Request",
			zap.String("path", path),
			zap.String("method", method),
			zap.Int("status", status),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}

// Metrics returns a middleware that records Prometheus metrics
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start)
		status := http.StatusText(c.Writer.Status())

		metrics.RecordAPIRequest(path, method, status, duration)
	}
}

// Recovery returns a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Request panicked",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}

// RateLimiter returns a middleware that limits request rates
func RateLimiter(config RateLimiterConfig) gin.HandlerFunc {
	// Create rate limiter store
	store := memory.NewStore()

	// Create rate limiter instance
	rate := limiter.Rate{
		Limit:  int64(config.Limit),
		Period: config.Period,
	}

	rateLimiter := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get client IP from header or fallback to RemoteAddr
		clientIP := c.ClientIP()
		if header := c.GetHeader(config.ClientIPHeader); header != "" {
			clientIP = header
		}

		// Get limiter context for this request
		limiterCtx, err := rateLimiter.Get(c, clientIP)
		if err != nil {
			logger.Error("Rate limiter error", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		// Check if request is limited
		if limiterCtx.Reached {
			logger.Warn("Rate limit exceeded",
				zap.String("client_ip", clientIP),
				zap.Int("limit", config.Limit),
				zap.Duration("period", config.Period))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

// ErrorHandler returns a middleware that handles errors from handlers
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there were any errors during processing
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error
		err := c.Errors.Last()

		// Log the error
		logger.Error("Request error",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Error(err.Err))

		// Determine the error type and appropriate status code
		statusCode := http.StatusInternalServerError
		errorMessage := "Internal server error"

		// Check for known error types
		if err.IsType(gin.ErrorTypePublic) {
			// Public errors can be shown to the client
			errorMessage = err.Error()
		} else if err.IsType(gin.ErrorTypeBind) {
			// Binding errors (invalid request parameters)
			statusCode = http.StatusBadRequest
			errorMessage = "Invalid request parameters"
		}

		// Record metrics for errors
		metrics.RPCRequestsTotal.WithLabelValues(c.Request.Method, "error").Inc()

		// Send error response if one hasn't been sent already
		if !c.Writer.Written() {
			c.JSON(statusCode, gin.H{
				"error": errorMessage,
			})
		}
	}
}

// ConfigureRateLimiters sets up rate limiting for various API endpoints
func ConfigureRateLimiters(router *gin.Engine) {
	// API endpoints - allow more frequent access
	apiConfig := DefaultRateLimiterConfig()
	apiConfig.Limit = 200 // Higher limit for API calls

	// Block height endpoint - very frequent access allowed
	blockHeightConfig := DefaultRateLimiterConfig()
	blockHeightConfig.Limit = 500 // Even higher limit for block height queries

	// Setup rate limiting for specific API groups
	router.Group("/api").
		Use(RateLimiter(apiConfig))

	// Higher limits for block queries
	router.Group("/api/blocks").
		Use(RateLimiter(blockHeightConfig))

	// Default rate limiting for all other endpoints
	defaultConfig := DefaultRateLimiterConfig()
	router.Use(RateLimiter(defaultConfig))
}
