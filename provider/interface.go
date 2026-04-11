// Package provider defines the core interfaces that external LLM providers must implement.
// External provider packages should import this package to implement the Provider interface.
package provider

import "context"

// Provider defines the interface that all LLM providers must implement.
// External packages can implement this interface and inject via omnillm.ClientConfig.CustomProvider.
//
// Example usage in external package:
//
//	import "github.com/plexusone/omnillm-core/provider"
//
//	func NewMyProvider(apiKey string) provider.Provider {
//	    return &myProvider{apiKey: apiKey}
//	}
type Provider interface {
	// CreateChatCompletion creates a new chat completion
	CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)

	// CreateChatCompletionStream creates a streaming chat completion
	CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error)

	// Close closes the provider and cleans up resources
	Close() error

	// Name returns the provider name
	Name() string
}

// ChatCompletionStream represents a streaming chat completion response
type ChatCompletionStream interface {
	// Recv receives the next chunk from the stream
	Recv() (*ChatCompletionChunk, error)

	// Close closes the stream
	Close() error
}
