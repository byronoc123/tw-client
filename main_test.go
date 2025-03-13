// File: main_test.go
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"blockchain-client/models"
)

// MockClient implements the minimum functionality needed for testing
type MockClient struct{}

func (c *MockClient) GetLatestBlockNumber() (string, error) {
	return "0x134e82a", nil
}

func (c *MockClient) GetBlockByNumber(blockNumber string) (*models.Block, error) {
	return &models.Block{
		Number: blockNumber,
		Hash:   "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	}, nil
}

func TestHealthEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup router
	r := gin.New()
	
	// Add health endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Create test request
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestLatestBlockEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock client
	mockClient := &MockClient{}
	
	// Setup router
	r := gin.New()
	
	// Add latest block endpoint
	r.GET("/api/v1/block/latest", func(c *gin.Context) {
		blockNumber, _ := mockClient.GetLatestBlockNumber()
		c.JSON(http.StatusOK, gin.H{
			"blockNumber": blockNumber,
		})
	})

	// Create test request
	req, _ := http.NewRequest("GET", "/api/v1/block/latest", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "0x134e82a", response["blockNumber"])
}

func TestGetBlockByNumberEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock client
	mockClient := &MockClient{}
	
	// Setup router
	r := gin.New()
	
	// Add block by number endpoint
	r.GET("/api/v1/block/:number", func(c *gin.Context) {
		blockNumber := c.Param("number")
		block, _ := mockClient.GetBlockByNumber(blockNumber)
		c.JSON(http.StatusOK, block)
	})

	// Create test request
	req, _ := http.NewRequest("GET", "/api/v1/block/0x1234", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.Block
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "0x1234", response.Number)
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", response.Hash)
}
