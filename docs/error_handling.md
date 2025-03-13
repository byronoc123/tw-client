# Error Handling and Logging

The blockchain client now implements a robust error handling and logging system.

## Error Handling System

### Error Classification

Errors are categorized by type:
- `CLIENT`: Client-side errors (400 Bad Request)
- `VALIDATION`: Input validation errors (400 Bad Request)
- `BLOCKCHAIN`: Blockchain RPC-related errors (503 Service Unavailable)
- `INTERNAL`: Server-side errors (500 Internal Server Error)
- `NOT_FOUND`: Resource not found errors (404 Not Found)
- `UNAUTHORIZED`: Authentication errors (401 Unauthorized)

### Key Features

- **Contextual Information**: Errors include additional context data to aid in debugging
- **HTTP Status Codes**: Each error type maps to an appropriate HTTP status code
- **Middleware Integration**: Centralized error handling middleware provides consistent error responses
- **Error Wrapping**: Original errors are preserved while adding context

### Usage Example

```go
// Creating a new error
err := errors.New(errors.ErrorTypeValidation, "Invalid block number format", nil)

// Wrapping an existing error
wrappedErr := errors.Wrap(err, errors.ErrorTypeBlockchain, "Failed to process block request")

// Adding additional context
contextErr := wrappedErr.WithData("block_number", "0x12345")
```

## Structured Logging

The application uses [zap](https://github.com/uber-go/zap) for structured logging:

### Key Features

- **JSON-formatted logs** in production for easy parsing
- **Development-friendly logs** with colors in development mode
- **Request logging** with method, path, status code, and latency information
- **Error context** including error type, message, and additional data
- **Performance** through zap's optimized logging approach

### Log Levels

- `DEBUG`: Detailed information for debugging purposes
- `INFO`: General operational information
- `WARN`: Warning conditions that should be addressed
- `ERROR`: Error conditions that prevent normal operation
- `FATAL`: Critical errors that cause the application to terminate

### Example Log Output

Development mode:
```
2025-03-12T16:15:00.000Z INFO blockchain-client/server.go:42 Server starting {"address": ":8080"}
2025-03-12T16:15:01.123Z INFO blockchain-client/middleware/logger.go:35 HTTP Request {"method": "GET", "path": "/api/v1/block/latest", "status": 200, "latency": "15ms"}
2025-03-12T16:15:02.456Z WARN blockchain-client/server.go:98 Block not found {"block_number": "0x1234567890"}
```

Production mode (JSON format):
```json
{"level":"info","timestamp":"2025-03-12T16:15:00.000Z","caller":"blockchain-client/server.go:42","msg":"Server starting","service":"blockchain-client","address":":8080"}
{"level":"info","timestamp":"2025-03-12T16:15:01.123Z","caller":"blockchain-client/middleware/logger.go:35","msg":"HTTP Request","service":"blockchain-client","method":"GET","path":"/api/v1/block/latest","status":200,"latency":"15ms"}
{"level":"warn","timestamp":"2025-03-12T16:15:02.456Z","caller":"blockchain-client/server.go:98","msg":"Block not found","service":"blockchain-client","block_number":"0x1234567890"}
```

## Configuration

Logging verbosity is controlled through the application's mode:
- Development mode: Verbose, colorized output
- Production mode: JSON-formatted, optimized for machine processing

Set the environment variable `GIN_MODE=release` to enable production mode logging.
