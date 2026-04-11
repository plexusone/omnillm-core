package omnillm

import "github.com/plexusone/omnillm-core/provider"

// Type aliases for backward compatibility and convenience.
// These allow thick providers to import from omnillm-core root package.
// Note: Provider and ChatCompletionStream are defined in provider.go
type (
	// Message types
	Role         = provider.Role
	Message      = provider.Message
	ToolCall     = provider.ToolCall
	ToolFunction = provider.ToolFunction

	// Request/Response types
	ChatCompletionRequest  = provider.ChatCompletionRequest
	ChatCompletionResponse = provider.ChatCompletionResponse
	ChatCompletionChoice   = provider.ChatCompletionChoice
	ChatCompletionChunk    = provider.ChatCompletionChunk

	// Tool types
	Tool     = provider.Tool
	ToolSpec = provider.ToolSpec

	// Usage
	Usage = provider.Usage

	// Response format
	ResponseFormat = provider.ResponseFormat
)

// Capabilities describes the features supported by a provider.
// Thick providers can implement a Capabilities() method returning this struct.
// Note: This is not part of the Provider interface but useful for feature detection.
type Capabilities struct {
	// Tools indicates support for tool/function calling.
	Tools bool

	// Streaming indicates support for streaming responses.
	Streaming bool

	// Vision indicates support for image inputs in messages.
	Vision bool

	// JSON indicates support for JSON response format.
	JSON bool

	// SystemRole indicates support for system messages.
	SystemRole bool

	// MaxContextWindow is the maximum context window size in tokens.
	MaxContextWindow int

	// SupportsMaxTokens indicates if the provider supports the max_tokens parameter.
	SupportsMaxTokens bool
}

// Role constants for convenience
const (
	RoleSystem    = provider.RoleSystem
	RoleUser      = provider.RoleUser
	RoleAssistant = provider.RoleAssistant
	RoleTool      = provider.RoleTool
)

// ModelInfo represents information about a model
type ModelInfo struct {
	ID        string       `json:"id"`
	Provider  ProviderName `json:"provider"`
	Name      string       `json:"name"`
	MaxTokens int          `json:"max_tokens"`
}

// GetModelInfo returns model information
func GetModelInfo(modelID string) *ModelInfo {
	modelMap := map[string]ModelInfo{
		ModelGPT4o: {
			ID:        ModelGPT4o,
			Provider:  ProviderNameOpenAI,
			Name:      "GPT-4o",
			MaxTokens: 128000,
		},
		ModelClaude3Opus: {
			ID:        ModelClaude3Opus,
			Provider:  ProviderNameAnthropic,
			Name:      "Claude 3 Opus",
			MaxTokens: 200000,
		},
		ModelBedrockClaude3Sonnet: {
			ID:        ModelBedrockClaude3Sonnet,
			Provider:  ProviderNameBedrock,
			Name:      "Claude 3 Sonnet (Bedrock)",
			MaxTokens: 200000,
		},
		ModelOllamaLlama3_8B: {
			ID:        ModelOllamaLlama3_8B,
			Provider:  ProviderNameOllama,
			Name:      "Llama 3 8B",
			MaxTokens: 8192,
		},
		ModelOllamaMistral7B: {
			ID:        ModelOllamaMistral7B,
			Provider:  ProviderNameOllama,
			Name:      "Mistral 7B",
			MaxTokens: 32768,
		},
		ModelOllamaCodeLlama: {
			ID:        ModelOllamaCodeLlama,
			Provider:  ProviderNameOllama,
			Name:      "CodeLlama 13B",
			MaxTokens: 16384,
		},
	}

	if info, exists := modelMap[modelID]; exists {
		return &info
	}
	return nil
}
