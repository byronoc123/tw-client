package main

import (
	"os"
	"strconv"
	"time"

	"blockchain-client/pkg/logger"
	"blockchain-client/rpc"
	"blockchain-client/server"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger with rotation for production use
	isProduction := os.Getenv("GIN_MODE") == "release"
	rotationConfig := logger.DefaultRotationConfig()

	// Set the appropriate log level based on the environment
	logLevel := "info"
	if !isProduction {
		logLevel = "debug" // More verbose logging in development
	}

	logger.InitWithRotation(logLevel, rotationConfig)
	defer logger.Sync()

	logger.Info("Starting blockchain client application")

	// Get configuration from environment variables
	rpcURL := getEnv("RPC_URL", "https://polygon-rpc.com/")
	timeoutStr := getEnv("TIMEOUT_SECONDS", "10")
	port := getEnv("PORT", "8080")

	// Parse timeout
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		logger.Fatal("Invalid timeout value", zap.String("timeout", timeoutStr), zap.Error(err))
	}

	// Create enhanced RPC client
	logger.Info("Initializing blockchain RPC client", zap.String("url", rpcURL))
	client := rpc.NewEnhancedClient(rpcURL, time.Duration(timeout)*time.Second)

	// Create and start server with rate limiting and metrics
	logger.Info("Initializing enhanced HTTP server", zap.String("port", port))
	srv := server.NewEnhanced(client, port)

	// Log startup message
	logger.Info("Server initialized with rate limiting, metrics, and enhanced logging",
		zap.String("port", port),
		zap.String("metrics_endpoint", "/metrics"),
		zap.String("log_file", rotationConfig.Filename))

	// Start the server
	if err := srv.Start(); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
