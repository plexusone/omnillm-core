package omnillm

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	// Common errors
	ErrUnsupportedProvider  = errors.New("unsupported provider")
	ErrBedrockExternal      = errors.New("bedrock provider moved to github.com/plexusone/omnillm-bedrock; use CustomProvider to inject it")
	ErrInvalidConfiguration = errors.New("invalid configuration")
	ErrNoProviders          = errors.New("at least one provider must be configured")
	ErrEmptyAPIKey          = errors.New("API key cannot be empty")
	ErrEmptyModel           = errors.New("model cannot be empty")
	ErrEmptyMessages        = errors.New("messages cannot be empty")
	ErrStreamClosed         = errors.New("stream is closed")
	ErrInvalidResponse      = errors.New("invalid response format")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrQuotaExceeded        = errors.New("quota exceeded")
	ErrInvalidRequest       = errors.New("invalid request")
	ErrModelNotFound        = errors.New("model not found")
	ErrServerError          = errors.New("server error")
	ErrNetworkError         = errors.New("network error")

	// Aliases for thick provider compatibility
	ErrInvalidAPIKey = ErrEmptyAPIKey
)

// APIError represents an error response from the API
type APIError struct {
	StatusCode int          `json:"status_code"`
	Message    string       `json:"message"`
	Type       string       `json:"type"`
	Code       string       `json:"code"`
	Provider   ProviderName `json:"provider"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s (status: %d, type: %s, code: %s)",
		e.Provider, e.Message, e.StatusCode, e.Type, e.Code)
}

// NewAPIError creates a new API error.
// This signature is compatible with thick providers that pass (provider, statusCode, errorType, message).
func NewAPIError(provider string, statusCode int, errorType, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Type:       errorType,
		Provider:   ProviderName(provider),
	}
}

// NewAPIErrorFull creates a new API error with all fields.
func NewAPIErrorFull(provider ProviderName, statusCode int, message, errorType, code string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Type:       errorType,
		Code:       code,
		Provider:   provider,
	}
}

// ErrorCategory classifies errors for retry/fallback logic
type ErrorCategory int

const (
	// ErrorCategoryUnknown indicates the error type could not be determined
	ErrorCategoryUnknown ErrorCategory = iota
	// ErrorCategoryRetryable indicates the error is transient and the request can be retried
	// Examples: rate limits (429), server errors (5xx), network errors
	ErrorCategoryRetryable
	// ErrorCategoryNonRetryable indicates the error is permanent and retrying won't help
	// Examples: auth errors (401/403), invalid requests (400), not found (404)
	ErrorCategoryNonRetryable
)

// String returns the string representation of the error category
func (c ErrorCategory) String() string {
	switch c {
	case ErrorCategoryRetryable:
		return "retryable"
	case ErrorCategoryNonRetryable:
		return "non-retryable"
	default:
		return "unknown"
	}
}

// ClassifyError determines the category of an error for retry/fallback decisions
func ClassifyError(err error) ErrorCategory {
	if err == nil {
		return ErrorCategoryUnknown
	}

	// Check for APIError with status code
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return classifyStatusCode(apiErr.StatusCode)
	}

	// Check for network errors
	if isNetworkError(err) {
		return ErrorCategoryRetryable
	}

	// Check for known error types
	if errors.Is(err, ErrRateLimitExceeded) || errors.Is(err, ErrServerError) || errors.Is(err, ErrNetworkError) {
		return ErrorCategoryRetryable
	}

	if errors.Is(err, ErrInvalidRequest) || errors.Is(err, ErrModelNotFound) ||
		errors.Is(err, ErrEmptyAPIKey) || errors.Is(err, ErrEmptyModel) ||
		errors.Is(err, ErrEmptyMessages) || errors.Is(err, ErrInvalidConfiguration) {
		return ErrorCategoryNonRetryable
	}

	// Check error message for common patterns
	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "rate limit") || strings.Contains(errMsg, "too many requests") {
		return ErrorCategoryRetryable
	}
	if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "connection refused") {
		return ErrorCategoryRetryable
	}
	if strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "forbidden") ||
		strings.Contains(errMsg, "invalid api key") || strings.Contains(errMsg, "authentication") {
		return ErrorCategoryNonRetryable
	}

	return ErrorCategoryUnknown
}

// classifyStatusCode maps HTTP status codes to error categories
func classifyStatusCode(statusCode int) ErrorCategory {
	switch {
	case statusCode == 429: // Rate limit
		return ErrorCategoryRetryable
	case statusCode >= 500 && statusCode < 600: // Server errors
		return ErrorCategoryRetryable
	case statusCode == 408: // Request timeout
		return ErrorCategoryRetryable
	case statusCode == 401 || statusCode == 403: // Auth errors
		return ErrorCategoryNonRetryable
	case statusCode == 400 || statusCode == 404: // Client errors
		return ErrorCategoryNonRetryable
	case statusCode == 422: // Validation error
		return ErrorCategoryNonRetryable
	case statusCode >= 400 && statusCode < 500: // Other client errors
		return ErrorCategoryNonRetryable
	default:
		return ErrorCategoryUnknown
	}
}

// isNetworkError checks if the error is a network-related error
func isNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Check for net.Error (timeout, temporary)
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	// Check for net.OpError
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}

	// Check error message for network-related patterns
	errMsg := strings.ToLower(err.Error())
	networkPatterns := []string{
		"connection refused",
		"connection reset",
		"no such host",
		"network is unreachable",
		"i/o timeout",
		"tls handshake",
		"eof",
	}
	for _, pattern := range networkPatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

// IsRetryableError returns true if the error is transient and the request can be retried.
// This is useful for fallback provider logic - only retry on retryable errors.
func IsRetryableError(err error) bool {
	category := ClassifyError(err)
	// Treat unknown errors as retryable to be safe (fail-open)
	return category == ErrorCategoryRetryable || category == ErrorCategoryUnknown
}

// IsNonRetryableError returns true if the error is permanent and retrying won't help.
func IsNonRetryableError(err error) bool {
	return ClassifyError(err) == ErrorCategoryNonRetryable
}
