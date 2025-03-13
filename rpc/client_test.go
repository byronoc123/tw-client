package rpc

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestBlockNumber(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		assert.Equal(t, "POST", r.Method)

		// Check content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x134e82a"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewEnhancedClient(server.URL, 10*time.Second)

	// Call the method
	blockNumber, err := client.GetLatestBlockNumber()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "0x134e82a", blockNumber)
}

func TestGetBlockByNumber(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		assert.Equal(t, "POST", r.Method)

		// Check content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
			"jsonrpc":"2.0",
			"id":1,
			"result":{
				"number":"0x134e82a",
				"hash":"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				"parentHash":"0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				"nonce":"0x0000000000000000",
				"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"stateRoot":"0xd7f8974fb5ac78d9ac099b9ad5018bedc2ce0a72dad1827a1709da30580f0544",
				"receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"miner":"0x0000000000000000000000000000000000000000",
				"difficulty":"0x0",
				"totalDifficulty":"0x0",
				"extraData":"0x",
				"size":"0x1000",
				"gasLimit":"0x1000000",
				"gasUsed":"0x500000",
				"timestamp":"0x60123456",
				"transactions":[],
				"uncles":[]
			}
		}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewEnhancedClient(server.URL, 10*time.Second)

	// Call the method
	block, err := client.GetBlockByNumber("0x134e82a")

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, block)
	assert.Equal(t, "0x134e82a", block.Number)
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", block.Hash)
}

func TestErrorHandling(t *testing.T) {
	// Create a server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Internal server error"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewEnhancedClient(server.URL, 10*time.Second)

	// Call the method and expect an error
	_, err := client.GetLatestBlockNumber()
	assert.Error(t, err)
}
