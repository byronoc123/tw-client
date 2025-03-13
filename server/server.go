package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"blockchain-client/models"
	"blockchain-client/pkg/errors"
	"blockchain-client/pkg/logger"
	"blockchain-client/pkg/metrics"
	"blockchain-client/pkg/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BlockchainClient interface for blockchain operations
type BlockchainClient interface {
	GetLatestBlockNumber() (string, error)
	GetBlockByNumber(blockNumber string) (*models.Block, error)
}

// EnhancedBlockchainClient interface for blockchain operations with metrics support
type EnhancedBlockchainClient interface {
	BlockchainClient
	// Additional methods can be added as needed
}

// EnhancedServer represents the HTTP server with enhanced features
type EnhancedServer struct {
	router  *gin.Engine
	client  EnhancedBlockchainClient
	address string
}

// NewEnhanced creates and configures a new enhanced server
func NewEnhanced(client EnhancedBlockchainClient, port string) *EnhancedServer {
	// Configure router
	router := gin.New()
	
	// Use our custom middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(metrics.MetricsMiddleware())

	// Configure rate limiters
	middleware.ConfigureRateLimiters(router)
	
	// Register metrics endpoint
	metrics.RegisterMetricsEndpoint(router)

	server := &EnhancedServer{
		router:  router,
		client:  client,
		address: fmt.Sprintf(":%s", port),
	}

	// Set up routes
	server.setupRoutes()

	return server
}

// Start starts the HTTP server
func (s *EnhancedServer) Start() error {
	logger.Info("Enhanced server starting", zap.String("address", s.address))
	return s.router.Run(s.address)
}

// setupRoutes configures the API routes
func (s *EnhancedServer) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Get latest block number
		api.GET("/block/latest", s.getLatestBlockNumber)
		
		// Get block by number
		api.GET("/block/:number", s.getBlockByNumber)
	}
}

// getLatestBlockNumber handles requests for the latest block number
func (s *EnhancedServer) getLatestBlockNumber(c *gin.Context) {
	// Start metrics timer
	start := time.Now()
	
	blockNumber, err := s.client.GetLatestBlockNumber()
	
	// Record RPC metrics
	duration := time.Since(start).Seconds()
	if err != nil {
		metrics.RPCRequestsTotal.WithLabelValues("eth_blockNumber", "error").Inc()
		logger.Error("Failed to get latest block number", zap.Error(err))
		c.Error(errors.Wrap(err, errors.ErrorTypeBlockchain, "Failed to get latest block number"))
		return
	}
	
	// Record successful RPC metrics
	metrics.RPCRequestsTotal.WithLabelValues("eth_blockNumber", "success").Inc()
	metrics.RPCRequestDuration.WithLabelValues("eth_blockNumber").Observe(duration)
	
	// Update blockchain height metric - convert hex string to float64
	// Remove "0x" prefix and parse as hexadecimal
	if len(blockNumber) > 2 && blockNumber[:2] == "0x" {
		if blockVal, err := strconv.ParseUint(blockNumber[2:], 16, 64); err == nil {
			metrics.UpdateBlockchainHeight(float64(blockVal))
		}
	}
	
	logger.Debug("Retrieved latest block number", zap.String("block_number", blockNumber))
	c.JSON(http.StatusOK, gin.H{
		"blockNumber": blockNumber,
	})
}

// getBlockByNumber handles requests for a specific block by number
func (s *EnhancedServer) getBlockByNumber(c *gin.Context) {
	blockNumberParam := c.Param("number")
	
	// Log the incoming request
	logger.Debug("Block details requested", zap.String("block_number", blockNumberParam))
	
	// Validate and format block number
	formattedBlockNumber, err := validateAndFormatBlockNumber(blockNumberParam)
	if err != nil {
		logger.Warn("Invalid block number format", 
			zap.String("input", blockNumberParam), 
			zap.Error(err))
		c.Error(errors.Wrap(err, errors.ErrorTypeValidation, "Invalid block number format"))
		return
	}
	
	// Start metrics timer
	start := time.Now()
	
	// Get block details
	block, err := s.client.GetBlockByNumber(formattedBlockNumber)
	
	// Record RPC metrics
	duration := time.Since(start).Seconds()
	if err != nil {
		metrics.RPCRequestsTotal.WithLabelValues("eth_getBlockByNumber", "error").Inc()
		
		if errors.IsType(err, errors.ErrorTypeNotFound) {
			logger.Warn("Block not found", 
				zap.String("block_number", formattedBlockNumber))
			c.Error(err)
		} else {
			logger.Error("Failed to get block details", 
				zap.String("block_number", formattedBlockNumber), 
				zap.Error(err))
			
			// Create a data map for the error
			errData := map[string]interface{}{
				"block_number": formattedBlockNumber,
			}
			
			c.Error(errors.Wrap(err, errors.ErrorTypeBlockchain, 
				"Failed to get block data").WithData(errData))
		}
		return
	}
	
	// Record successful RPC metrics
	metrics.RPCRequestsTotal.WithLabelValues("eth_getBlockByNumber", "success").Inc()
	metrics.RPCRequestDuration.WithLabelValues("eth_getBlockByNumber").Observe(duration)
	
	logger.Debug("Successfully retrieved block",
		zap.String("block_number", block.Number),
		zap.String("block_hash", block.Hash))
	
	c.JSON(http.StatusOK, block)
}

// validateAndFormatBlockNumber validates and formats block number string
func validateAndFormatBlockNumber(blockNumber string) (string, error) {
	// Handle special case for "latest"
	if blockNumber == "latest" {
		return "latest", nil
	}

	// Handle numeric inputs
	if len(blockNumber) > 0 {
		// Add "0x" prefix if not already there and not empty
		if blockNumber[0] != '0' || (len(blockNumber) > 1 && blockNumber[1] != 'x') {
			// Strip "0x" if it exists
			if len(blockNumber) > 2 && blockNumber[0] == '0' && blockNumber[1] == 'x' {
				blockNumber = blockNumber[2:]
			}
			
			// For numeric inputs without 0x prefix, validate and add the prefix
			blockNumber = "0x" + blockNumber
		}
		return blockNumber, nil
	}

	return "", errors.New(errors.ErrorTypeValidation, "Invalid block number format")
}
