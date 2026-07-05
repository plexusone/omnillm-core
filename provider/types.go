package provider

// Role represents the role of a message sender
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message represents a chat message
type Message struct {
	Role       Role       `json:"role"`
	Content    string     `json:"content"`
	Name       *string    `json:"name,omitempty"`
	ToolCallID *string    `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall represents a tool function call
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents the function being called
type ToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatCompletionRequest represents a request for chat completion
type ChatCompletionRequest struct {
	Model            string          `json:"model"`
	Messages         []Message       `json:"messages"`
	MaxTokens        *int            `json:"max_tokens,omitempty"`
	Temperature      *float64        `json:"temperature,omitempty"`
	TopP             *float64        `json:"top_p,omitempty"`
	TopK             *int            `json:"top_k,omitempty"` // Anthropic, Gemini, Ollama
	Stream           *bool           `json:"stream,omitempty"`
	Stop             []string        `json:"stop,omitempty"`
	PresencePenalty  *float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int  `json:"logit_bias,omitempty"`
	User             *string         `json:"user,omitempty"`
	Tools            []Tool          `json:"tools,omitempty"`
	ToolChoice       any             `json:"tool_choice,omitempty"`
	Seed             *int            `json:"seed,omitempty"`             // OpenAI, X.AI - for reproducible outputs
	N                *int            `json:"n,omitempty"`                // OpenAI - number of completions
	ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`  // OpenAI, Gemini - JSON mode
	Logprobs         *bool           `json:"logprobs,omitempty"`         // OpenAI - return log probabilities
	TopLogprobs      *int            `json:"top_logprobs,omitempty"`     // OpenAI - number of top logprobs
	ReasoningEffort  *string         `json:"reasoning_effort,omitempty"` // OpenAI, X.AI - "none", "low", "medium", "high"
	Thinking         *ThinkingConfig `json:"thinking,omitempty"`         // Anthropic - extended thinking configuration
}

// ResponseFormat specifies the format of the response
type ResponseFormat struct {
	Type string `json:"type"` // "text" or "json_object"
}

// ReasoningEffort values for OpenAI-compatible providers.
// Use these with the ReasoningEffort field in ChatCompletionRequest.
const (
	ReasoningEffortNone   = "none"   // Disables reasoning entirely
	ReasoningEffortLow    = "low"    // Light reasoning (default for reasoning models)
	ReasoningEffortMedium = "medium" // Moderate reasoning for complex tasks
	ReasoningEffortHigh   = "high"   // Deep reasoning for demanding problems
)

// ThinkingType values for Anthropic-style extended thinking.
// Use these with ThinkingConfig.Type.
const (
	ThinkingTypeEnabled  = "enabled"  // Enable thinking with explicit BudgetTokens
	ThinkingTypeDisabled = "disabled" // Disable thinking entirely
	ThinkingTypeAdaptive = "adaptive" // Let the model decide thinking depth
)

// ThinkingConfig configures Anthropic-style extended thinking.
// For Anthropic models that support extended thinking (Claude with thinking capability).
type ThinkingConfig struct {
	// Type specifies the thinking mode: "enabled", "disabled", or "adaptive".
	Type string `json:"type"`

	// BudgetTokens is the token budget for thinking.
	// Required when Type is "enabled".
	// Minimum 1024, must be less than max_tokens.
	BudgetTokens *int64 `json:"budget_tokens,omitempty"`
}

// Tool represents a tool that can be called
type Tool struct {
	Type     string   `json:"type"`
	Function ToolSpec `json:"function"`
}

// ToolSpec defines a tool specification
type ToolSpec struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

// ChatCompletionResponse represents a response from chat completion
type ChatCompletionResponse struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	SystemFingerprint *string                `json:"system_fingerprint,omitempty"`
	Choices           []ChatCompletionChoice `json:"choices"`
	Usage             Usage                  `json:"usage"`
	ProviderMetadata  map[string]any         `json:"provider_metadata,omitempty"` // Provider-specific metadata
}

// ChatCompletionChoice represents a single choice in the response
type ChatCompletionChoice struct {
	Index        int      `json:"index"`
	Message      Message  `json:"message"`
	Delta        *Message `json:"delta,omitempty"`
	FinishReason *string  `json:"finish_reason"`
	Logprobs     any      `json:"logprobs,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionChunk represents a chunk in streaming response
type ChatCompletionChunk struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	SystemFingerprint *string                `json:"system_fingerprint,omitempty"`
	Choices           []ChatCompletionChoice `json:"choices"`
	Usage             *Usage                 `json:"usage,omitempty"`
	ProviderMetadata  map[string]any         `json:"provider_metadata,omitempty"` // Provider-specific metadata
}

// EmbeddingEncodingFormat specifies the format of the embedding vectors
type EmbeddingEncodingFormat string

const (
	// EmbeddingEncodingFormatFloat returns embeddings as float arrays (default)
	EmbeddingEncodingFormatFloat EmbeddingEncodingFormat = "float"
	// EmbeddingEncodingFormatBase64 returns embeddings as base64-encoded bytes
	EmbeddingEncodingFormatBase64 EmbeddingEncodingFormat = "base64"
)

// EmbeddingRequest represents a request for text embeddings
type EmbeddingRequest struct {
	// Model is the embedding model to use (e.g., "text-embedding-3-small")
	Model string `json:"model"`

	// Input is the text(s) to embed. Can be a single string or array of strings.
	Input []string `json:"input"`

	// EncodingFormat specifies the format of the returned embeddings (optional)
	// Defaults to "float". Some providers support "base64" for efficiency.
	EncodingFormat EmbeddingEncodingFormat `json:"encoding_format,omitempty"`

	// Dimensions specifies the number of dimensions for the output vectors (optional)
	// Only supported by some models (e.g., text-embedding-3-small/large)
	Dimensions *int `json:"dimensions,omitempty"`

	// User is an optional unique identifier for the end-user
	User *string `json:"user,omitempty"`
}

// EmbeddingResponse represents a response from embedding creation
type EmbeddingResponse struct {
	// Object is the object type (always "list" for embeddings)
	Object string `json:"object"`

	// Data contains the embedding vectors
	Data []EmbeddingData `json:"data"`

	// Model is the model used to generate the embeddings
	Model string `json:"model"`

	// Usage contains token usage information
	Usage EmbeddingUsage `json:"usage"`

	// ProviderMetadata contains provider-specific metadata
	ProviderMetadata map[string]any `json:"provider_metadata,omitempty"`
}

// EmbeddingData represents a single embedding vector
type EmbeddingData struct {
	// Object is the object type (always "embedding")
	Object string `json:"object"`

	// Index is the index of this embedding in the input array
	Index int `json:"index"`

	// Embedding is the vector representation as float64 values
	Embedding []float64 `json:"embedding"`
}

// EmbeddingUsage represents token usage for embedding requests
type EmbeddingUsage struct {
	// PromptTokens is the number of tokens in the input
	PromptTokens int `json:"prompt_tokens"`

	// TotalTokens is the total number of tokens used
	TotalTokens int `json:"total_tokens"`
}
