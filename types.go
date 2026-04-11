package omnillm

import "github.com/plexusone/omnillm-core/provider"

// Type aliases for backward compatibility and convenience
type Role = provider.Role
type Message = provider.Message
type ToolCall = provider.ToolCall
type ToolFunction = provider.ToolFunction
type ChatCompletionRequest = provider.ChatCompletionRequest
type Tool = provider.Tool
type ToolSpec = provider.ToolSpec
type ChatCompletionResponse = provider.ChatCompletionResponse
type ChatCompletionChoice = provider.ChatCompletionChoice
type Usage = provider.Usage
type ChatCompletionChunk = provider.ChatCompletionChunk

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
