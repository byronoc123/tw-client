package rpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"blockchain-client/models"
	"blockchain-client/pkg/logger"

	"go.uber.org/zap"
)

// HealthCheck performs a health check on the RPC endpoint
func (c *EnhancedClient) HealthCheck(ctx context.Context) (bool, string, error) {
	logger.Debug("Performing RPC health check")
	
	// Create a context with timeout for health check
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	// Try to get net version as a lightweight check
	healthy, details, err := c.checkNetVersion(checkCtx)
	if err != nil {
		logger.Warn("RPC health check failed", zap.Error(err))
		return false, "Failed to connect to RPC endpoint", err
	}
	
	// Format description
	var description string
	if healthy {
		if chainName, ok := details["chainName"].(string); ok && chainName != "" {
			description = fmt.Sprintf("Connected to %s (Network ID: %s)", 
				chainName, details["networkId"])
		} else {
			description = fmt.Sprintf("Connected to RPC endpoint (Network ID: %s)", 
				details["networkId"])
		}
	} else {
		description = "Unhealthy RPC connection"
	}
	
	return healthy, description, nil
}

// checkNetVersion checks the RPC connection by getting the network version
func (c *EnhancedClient) checkNetVersion(ctx context.Context) (bool, map[string]interface{}, error) {
	// Create request for net_version
	requestBody := models.RPCRequest{
		JSONRPC: "2.0",
		Method:  "net_version",
		Params:  []interface{}{},
		ID:      1,
	}
	
	// Use a map to receive the response
	var response struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}
	
	// Send the request with context
	err := c.doRequestWithContext(ctx, requestBody, &response)
	if err != nil {
		return false, nil, err
	}
	
	// Check for RPC error
	if response.Error != nil {
		return false, nil, fmt.Errorf("RPC error: %s (code: %d)", 
			response.Error.Message, response.Error.Code)
	}
	
	// Basic response validation
	if response.Result == "" {
		return false, nil, fmt.Errorf("empty network ID")
	}
	
	// Get chain name based on network ID
	details := map[string]interface{}{
		"networkId": response.Result,
		"chainName": getChainNameFromNetworkID(response.Result),
	}
	
	return true, details, nil
}

// doRequestWithContext performs an RPC request with the provided context
func (c *EnhancedClient) doRequestWithContext(ctx context.Context, requestBody models.RPCRequest, responseObj interface{}) error {
	// Implementation similar to doRequest but with context
	// [Simplified for brevity - full implementation would be similar to doRequest with context support]
	
	// For simplicity in this example, we'll just wrap the existing method
	// In a real implementation, the entire HTTP request process would use the provided context
	resultChan := make(chan error, 1)
	
	go func() {
		resultChan <- c.doRequest(requestBody, responseObj)
	}()
	
	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// getChainNameFromNetworkID returns a human-readable chain name from network ID
func getChainNameFromNetworkID(networkID string) string {
	switch networkID {
	case "1":
		return "Ethereum Mainnet"
	case "3":
		return "Ropsten Testnet"
	case "4":
		return "Rinkeby Testnet"
	case "5":
		return "Goerli Testnet"
	case "42":
		return "Kovan Testnet"
	case "56":
		return "Binance Smart Chain"
	case "137":
		return "Polygon Mainnet"
	case "42161":
		return "Arbitrum One"
	case "10":
		return "Optimism"
	default:
		if strings.HasPrefix(networkID, "2018") {
			return "Ethereum Classic"
		}
		return ""
	}
}
