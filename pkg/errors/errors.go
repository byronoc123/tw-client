package errors

import (
	"errors"
	"fmt"
)

// Common error types for the application
const (
	ErrTypeInternal       = "internal_error"
	ErrTypeRPC            = "rpc_error"
	ErrTypeValidation     = "validation_error"
	ErrTypeTimeout        = "timeout_error"
	ErrTypeAuthentication = "auth_error"
	ErrTypeAuthorization  = "authorization_error"
	ErrTypeNotFound       = "not_found_error"
	ErrorTypeBlockchain   = "blockchain_error"
	ErrorTypeNotFound     = "not_found_error"  // Duplicate with different name for backward compatibility
	ErrorTypeValidation   = "validation_error" // For backward compatibility
	ErrTypePermission     = "permission_error" // For permission-related errors
)

// Standard errors
var (
	ErrNotFound = errors.New("not found")
)

// AppError represents a structured application error
type AppError struct {
	Type    string
	Message string
	Err     error
	Data    map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithData adds contextual data to the error
func (e *AppError) WithData(data map[string]interface{}) *AppError {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}

	for k, v := range data {
		e.Data[k] = v
	}

	return e
}

// New creates a new application error with the given type and message
func New(errType, message string) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
	}
}

// NewAppError creates a new application error
func NewAppError(errType, message string, err error) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// Wrap wraps an existing error with a type and message
func Wrap(err error, errType, message string) *AppError {
	return NewAppError(errType, message, err)
}

// NewInternalError creates a new internal error
func NewInternalError(message string, err error) *AppError {
	return NewAppError(ErrTypeInternal, message, err)
}

// NewRPCError creates a new RPC error
func NewRPCError(message string, err error) *AppError {
	return NewAppError(ErrTypeRPC, message, err)
}

// NewBlockchainError creates a new blockchain error
func NewBlockchainError(message string, err error) *AppError {
	return NewAppError(ErrorTypeBlockchain, message, err)
}

// NewValidationError creates a new validation error
func NewValidationError(message string, err error) *AppError {
	return NewAppError(ErrTypeValidation, message, err)
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(message string, err error) *AppError {
	return NewAppError(ErrTypeTimeout, message, err)
}

// NewPermissionError creates a new permission error
func NewPermissionError(message string, err error) *AppError {
	return NewAppError(ErrTypePermission, message, err)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(ErrTypeNotFound, message, err)
}

// IsAppError checks if an error is an AppError and returns it
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// IsType checks if an error is of a specific type
func IsType(err error, errType string) bool {
	appErr, ok := IsAppError(err)
	if !ok {
		return false
	}
	return appErr.Type == errType
}
