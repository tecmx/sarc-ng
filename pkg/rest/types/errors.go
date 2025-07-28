package types

import "fmt"

// APIError represents a standard API error response
type APIError struct {
	Message string `json:"error"`
	Code    int    `json:"code,omitempty"`
}

// Error implements the error interface
func (e APIError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("API Error (%d): %s", e.Code, e.Message)
	}
	return fmt.Sprintf("API Error: %s", e.Message)
}

// ValidationError represents field validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e ValidationError) Error() string {
	return fmt.Sprintf("Validation error for field '%s': %s", e.Field, e.Message)
}

// CLIError represents CLI-specific errors
type CLIError struct {
	Operation string
	Message   string
	Cause     error
}

// Error implements the error interface
func (e CLIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("CLI Error in %s: %s (caused by: %v)", e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("CLI Error in %s: %s", e.Operation, e.Message)
}

// NewCLIError creates a new CLI error
func NewCLIError(operation, message string, cause error) *CLIError {
	return &CLIError{
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}
