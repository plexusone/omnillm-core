package omnillm

import (
	"fmt"

	"github.com/plexusone/omnillm-core/provider"
)

// TokenEstimator estimates token counts for messages before sending to the API.
// This is useful for validating requests won't exceed model limits.
type TokenEstimator interface {
	// EstimateTokens estimates the token count for a set of messages.
	// The estimate may not be exact but should be reasonably close.
	EstimateTokens(model string, messages []provider.Message) (int, error)

	// GetContextWindow returns the maximum context window size for a model.
	// Returns 0 if the model is unknown.
	GetContextWindow(model string) int
}

// TokenEstimatorConfig configures token estimation behavior
type TokenEstimatorConfig struct {
	// CharactersPerToken is the average number of characters per token.
	// Default: 4.0 (reasonable for English text)
	// Lower values (e.g., 3.0) give more conservative estimates.
	CharactersPerToken float64

	// CustomContextWindows allows overriding context window sizes for specific models.
	// Keys should be model IDs (e.g., "gpt-4o", "claude-3-opus").
	CustomContextWindows map[string]int

	// TokenOverheadPerMessage is extra tokens added per message for formatting.
	// Default: 4 (accounts for role, separators, etc.)
	TokenOverheadPerMessage int
}

// DefaultTokenEstimatorConfig returns a TokenEstimatorConfig with sensible defaults
func DefaultTokenEstimatorConfig() TokenEstimatorConfig {
	return TokenEstimatorConfig{
		CharactersPerToken:      4.0,
		TokenOverheadPerMessage: 4,
	}
}

// defaultTokenEstimator implements TokenEstimator using character-based estimation
type defaultTokenEstimator struct {
	config TokenEstimatorConfig
}

// NewTokenEstimator creates a new token estimator with the given configuration.
// If config has zero values, defaults are used for those fields.
func NewTokenEstimator(config TokenEstimatorConfig) TokenEstimator {
	if config.CharactersPerToken == 0 {
		config.CharactersPerToken = 4.0
	}
	if config.TokenOverheadPerMessage == 0 {
		config.TokenOverheadPerMessage = 4
	}

	return &defaultTokenEstimator{config: config}
}

// EstimateTokens estimates the token count using character-based approximation.
// This provides a reasonable estimate for most use cases but is not exact.
func (e *defaultTokenEstimator) EstimateTokens(model string, messages []provider.Message) (int, error) {
	if len(messages) == 0 {
		return 0, nil
	}

	var totalChars int

	for _, msg := range messages {
		// Count role characters
		totalChars += len(msg.Role)

		// Count content characters
		totalChars += len(msg.Content)

		// Count tool calls if present
		for _, tc := range msg.ToolCalls {
			totalChars += len(tc.ID)
			totalChars += len(tc.Type)
			totalChars += len(tc.Function.Name)
			totalChars += len(tc.Function.Arguments)
		}

		// Count tool call ID if present
		if msg.ToolCallID != nil {
			totalChars += len(*msg.ToolCallID)
		}

		// Count name if present
		if msg.Name != nil {
			totalChars += len(*msg.Name)
		}
	}

	// Calculate base tokens from characters
	tokens := int(float64(totalChars) / e.config.CharactersPerToken)

	// Add per-message overhead
	tokens += len(messages) * e.config.TokenOverheadPerMessage

	// Add a small buffer for system formatting
	tokens += 3 // priming tokens

	return tokens, nil
}

// GetContextWindow returns the context window size for a model.
// Checks custom overrides first, then falls back to built-in knowledge.
func (e *defaultTokenEstimator) GetContextWindow(model string) int {
	// Check custom overrides first
	if e.config.CustomContextWindows != nil {
		if window, ok := e.config.CustomContextWindows[model]; ok {
			return window
		}
	}

	// Check built-in ModelInfo
	if info := GetModelInfo(model); info != nil && info.MaxTokens > 0 {
		return info.MaxTokens
	}

	// Fall back to extended lookup
	return getExtendedContextWindow(model)
}

// getExtendedContextWindow provides context window sizes for common models
func getExtendedContextWindow(model string) int {
	windows := map[string]int{
		// OpenAI Models
		"gpt-4o":              128000,
		"gpt-4o-mini":         128000,
		"gpt-4o-2024-05-13":   128000,
		"gpt-4o-2024-08-06":   128000,
		"gpt-4-turbo":         128000,
		"gpt-4-turbo-preview": 128000,
		"gpt-4-1106-preview":  128000,
		"gpt-4":               8192,
		"gpt-4-32k":           32768,
		"gpt-3.5-turbo":       16385,
		"gpt-3.5-turbo-16k":   16385,
		"gpt-3.5-turbo-1106":  16385,
		"o1":                  200000,
		"o1-preview":          128000,
		"o1-mini":             128000,

		// Anthropic Models
		"claude-opus-4":            200000,
		"claude-sonnet-4":          200000,
		"claude-3-opus":            200000,
		"claude-3-opus-20240229":   200000,
		"claude-3-sonnet":          200000,
		"claude-3-sonnet-20240229": 200000,
		"claude-3-haiku":           200000,
		"claude-3-haiku-20240307":  200000,
		"claude-3.5-sonnet":        200000,
		"claude-3.5-haiku":         200000,
		"claude-2.1":               200000,
		"claude-2":                 100000,
		"claude-instant-1.2":       100000,

		// Google Gemini Models
		"gemini-2.5-pro":          1000000,
		"gemini-2.5-flash":        1000000,
		"gemini-1.5-pro":          2000000,
		"gemini-1.5-pro-latest":   2000000,
		"gemini-1.5-flash":        1000000,
		"gemini-1.5-flash-latest": 1000000,
		"gemini-1.0-pro":          32768,
		"gemini-pro":              32768,

		// X.AI Grok Models
		"grok-4":      128000,
		"grok-4-fast": 128000,
		"grok-3":      128000,
		"grok-3-fast": 128000,
		"grok-2":      128000,
		"grok-beta":   128000,

		// Ollama Local Models (common defaults)
		"llama3":         8192,
		"llama3:8b":      8192,
		"llama3:70b":     8192,
		"llama2":         4096,
		"llama2:7b":      4096,
		"llama2:13b":     4096,
		"llama2:70b":     4096,
		"mistral":        32768,
		"mistral:7b":     32768,
		"mixtral":        32768,
		"mixtral:8x7b":   32768,
		"codellama":      16384,
		"codellama:7b":   16384,
		"codellama:13b":  16384,
		"codellama:34b":  16384,
		"gemma":          8192,
		"gemma:2b":       8192,
		"gemma:7b":       8192,
		"qwen":           32768,
		"qwen:7b":        32768,
		"qwen:14b":       32768,
		"deepseek-coder": 16384,
		"phi":            2048,
		"phi:2.7b":       2048,
	}

	if window, ok := windows[model]; ok {
		return window
	}

	// Default fallback
	return 4096
}

// TokenValidation contains the result of token validation
type TokenValidation struct {
	// EstimatedTokens is the estimated prompt token count
	EstimatedTokens int

	// ContextWindow is the model's maximum context window
	ContextWindow int

	// MaxCompletionTokens is the requested max completion tokens
	MaxCompletionTokens int

	// AvailableTokens is how many tokens are available for completion
	// (ContextWindow - EstimatedTokens)
	AvailableTokens int

	// ExceedsLimit is true if the prompt exceeds the context window
	ExceedsLimit bool

	// ExceedsWithCompletion is true if prompt + max_tokens exceeds context
	ExceedsWithCompletion bool
}

// ValidateTokens checks if the request fits within model limits.
// Returns validation details including whether limits are exceeded.
func ValidateTokens(
	estimator TokenEstimator,
	model string,
	messages []provider.Message,
	maxCompletionTokens int,
) (*TokenValidation, error) {
	estimated, err := estimator.EstimateTokens(model, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate tokens: %w", err)
	}

	contextWindow := estimator.GetContextWindow(model)
	available := contextWindow - estimated

	validation := &TokenValidation{
		EstimatedTokens:     estimated,
		ContextWindow:       contextWindow,
		MaxCompletionTokens: maxCompletionTokens,
		AvailableTokens:     available,
		ExceedsLimit:        estimated > contextWindow,
	}

	// Check if prompt + completion would exceed limit
	if maxCompletionTokens > 0 {
		validation.ExceedsWithCompletion = (estimated + maxCompletionTokens) > contextWindow
	}

	return validation, nil
}

// TokenLimitError is returned when a request exceeds token limits
type TokenLimitError struct {
	// EstimatedTokens is the estimated prompt token count
	EstimatedTokens int

	// ContextWindow is the model's maximum context window
	ContextWindow int

	// AvailableTokens is how many tokens are available (may be negative)
	AvailableTokens int

	// Model is the model ID
	Model string
}

func (e *TokenLimitError) Error() string {
	return fmt.Sprintf("estimated tokens (%d) exceed context window (%d) for model %s",
		e.EstimatedTokens, e.ContextWindow, e.Model)
}

// EstimatePromptTokens is a convenience function that creates a default estimator
// and estimates tokens for a set of messages.
func EstimatePromptTokens(model string, messages []provider.Message) (int, error) {
	estimator := NewTokenEstimator(DefaultTokenEstimatorConfig())
	return estimator.EstimateTokens(model, messages)
}

// GetModelContextWindow is a convenience function that returns the context window
// for a model using the default estimator.
func GetModelContextWindow(model string) int {
	estimator := NewTokenEstimator(DefaultTokenEstimatorConfig())
	return estimator.GetContextWindow(model)
}
