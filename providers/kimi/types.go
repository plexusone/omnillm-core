// Package kimi provides Kimi (Moonshot AI) API types (OpenAI-compatible format)
package kimi

// Request represents a Kimi API request (OpenAI-compatible format)
type Request struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	MaxTokens        *int      `json:"max_tokens,omitempty"`
	Temperature      *float64  `json:"temperature,omitempty"`
	TopP             *float64  `json:"top_p,omitempty"`
	Stream           *bool     `json:"stream,omitempty"`
	Stop             []string  `json:"stop,omitempty"`
	PresencePenalty  *float64  `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64  `json:"frequency_penalty,omitempty"`
	Seed             *int      `json:"seed,omitempty"`
}

// Message represents a message in Kimi format (OpenAI-compatible)
type Message struct {
	Role    string  `json:"role"`
	Content string  `json:"content"`
	Name    *string `json:"name,omitempty"`
}

// Response represents a Kimi API response (OpenAI-compatible)
type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice in Kimi response
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason *string `json:"finish_reason"`
}

// Usage represents token usage in Kimi response
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamChunk represents a chunk in Kimi streaming response (OpenAI-compatible)
type StreamChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []StreamDelta `json:"choices"`
	Usage   *Usage        `json:"usage,omitempty"`
}

// StreamDelta represents delta content in a streaming chunk
type StreamDelta struct {
	Index        int          `json:"index"`
	Delta        *DeltaChange `json:"delta,omitempty"`
	FinishReason *string      `json:"finish_reason"`
}

// DeltaChange represents the actual content change in a stream
type DeltaChange struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}
