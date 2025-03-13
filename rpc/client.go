package rpc

import (
	"blockchain-client/models"
	"blockchain-client/pkg/errors"
	"blockchain-client/pkg/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// EnhancedClient implements JSON-RPC over HTTP for blockchain interactions
// with improved error handling and logging
type EnhancedClient struct {
	rpcURL  string
	httpClient  *http.Client
	timeout time.Duration
}

// NewEnhancedClient creates a new RPC client with enhanced error handling
func NewEnhancedClient(rpcURL string, timeout time.Duration) *EnhancedClient {
	if timeout <= 0 {
		timeout = 10 * time.Second // Default timeout
	}

	logger.Debug("Initializing enhanced RPC client", 
		zap.String("rpc_url", rpcURL), 
		zap.Duration("timeout", timeout))

	return &EnhancedClient{
		rpcURL: rpcURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// GetLatestBlockNumber gets the latest block number from the blockchain
func (c *EnhancedClient) GetLatestBlockNumber() (string, error) {
	// Create JSON-RPC request
	requestBody := models.RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      1,
	}
	
	var response models.BlockNumberResponse
	err := c.doRequest(requestBody, &response)
	if err != nil {
		logger.Error("Failed to get latest block number", zap.Error(err))
		return "", errors.NewBlockchainError("Failed to get latest block number", err)
	}
	
	logger.Debug("Received latest block number", zap.String("block_number", response.Result))
	return response.Result, nil
}

// GetBlockByNumber retrieves a block by its number
// To maintain backward compatibility, we default includeTransactions to true
func (c *EnhancedClient) GetBlockByNumber(blockNumber string) (*models.Block, error) {
	return c.getBlockByNumber(blockNumber, true)
}

// getBlockByNumber is the internal implementation that allows control over the includeTransactions parameter
func (c *EnhancedClient) getBlockByNumber(blockNumber string, includeTransactions bool) (*models.Block, error) {
	// Create JSON-RPC request
	requestBody := models.RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{blockNumber, includeTransactions},
		ID:      1,
	}
	
	var response models.BlockResponse
	err := c.doRequest(requestBody, &response)
	if err != nil {
		logger.Error("Failed to get block by number", 
			zap.String("block_number", blockNumber), 
			zap.Error(err))
		return nil, errors.NewBlockchainError(fmt.Sprintf("Failed to get block data for block %s", blockNumber), err)
	}
	
	if response.Result == nil {
		logger.Warn("Block not found", zap.String("block_number", blockNumber))
		errData := make(map[string]interface{})
		errData["block_number"] = blockNumber
		return nil, errors.NewNotFoundError("Block not found", nil).WithData(errData)
	}
	
	return response.Result, nil
}

// doRequest performs an HTTP request to the RPC endpoint
func (c *EnhancedClient) doRequest(request models.RPCRequest, response interface{}) error {
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return errors.NewInternalError("Failed to marshal JSON request", err)
	}
	
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	reqStartTime := time.Now()
	logger.Debug("Sending RPC request", 
		zap.String("method", request.Method), 
		zap.String("url", c.rpcURL))
	
	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rpcURL, bytes.NewReader(requestJSON))
	if err != nil {
		return errors.NewInternalError("Failed to create HTTP request", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			logger.Warn("RPC request timed out",
				zap.String("method", request.Method),
				zap.Duration("elapsed", time.Since(reqStartTime)))
			return errors.NewTimeoutError("RPC request timed out", err)
		}
		
		logger.Error("RPC request failed", 
			zap.String("method", request.Method), 
			zap.Error(err))
		return errors.NewInternalError("Failed to execute HTTP request", err)
	}
	defer resp.Body.Close()
	
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.NewInternalError("Failed to read response body", err)
	}
	
	// Log response status and time
	logger.Debug("Received RPC response", 
		zap.String("method", request.Method),
		zap.Int("status", resp.StatusCode),
		zap.Duration("elapsed", time.Since(reqStartTime)))
	
	if resp.StatusCode != http.StatusOK {
		logger.Warn("Non-200 response from RPC",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(bodyBytes)))
		
		errData := make(map[string]interface{})
		errData["status_code"] = resp.StatusCode
		errData["response"] = string(bodyBytes)
		return errors.NewBlockchainError(
			fmt.Sprintf("RPC server returned non-200 response: %d", resp.StatusCode), nil).WithData(errData)
	}
	
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		logger.Error("Failed to unmarshal response",
			zap.Error(err),
			zap.String("response", string(bodyBytes)))
		return errors.NewInternalError("Failed to unmarshal JSON response", err)
	}
	
	// Check for RPC error response
	var rpcError models.RPCErrorResponse
	if err := json.Unmarshal(bodyBytes, &rpcError); err == nil && rpcError.Error.Code != 0 {
		logger.Error("RPC returned error",
			zap.Int("error_code", rpcError.Error.Code),
			zap.String("error_message", rpcError.Error.Message))
		
		errData := make(map[string]interface{})
		errData["error_code"] = rpcError.Error.Code
		errData["error_message"] = rpcError.Error.Message
		return errors.NewBlockchainError(
			fmt.Sprintf("RPC error: %s (code: %d)", rpcError.Error.Message, rpcError.Error.Code), nil).WithData(errData)
	}
	
	return nil
}
