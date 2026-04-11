package main

import (
	"context"
	"fmt"
	"log"

	"github.com/plexusone/omnillm-core"
	"github.com/plexusone/omnillm-core/provider"
)

// Example of a 3rd party provider implementation following omnillm's internal pattern
// This would be in an external package like github.com/someone/omnillm-custom

// Step 1: HTTP Client (like providers/ollama/ollama.go)
type httpClient struct {
	name   string
	apiKey string
}

func newHTTPClient(name, apiKey string) *httpClient {
	return &httpClient{name: name, apiKey: apiKey}
}

func (c *httpClient) Name() string {
	return c.name
}

// Step 2: Provider Adapter (like the adapters in providers.go)
type customProvider struct {
	client *httpClient
}

// NewCustomProvider creates a new custom provider following omnillm's architecture pattern
// Note: Now uses the public provider.Provider interface that external packages can import
func NewCustomProvider(name, apiKey string) provider.Provider {
	client := newHTTPClient(name, apiKey)
	return &customProvider{client: client}
}

func (p *customProvider) Name() string {
	return p.client.Name()
}

func (p *customProvider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
	// This would normally call p.client.CreateCompletion() and convert the response
	// Mock implementation for demonstration
	return &provider.ChatCompletionResponse{
		ID:      "custom-123",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   req.Model,
		Choices: []provider.ChatCompletionChoice{
			{
				Index: 0,
				Message: provider.Message{
					Role:    provider.RoleAssistant,
					Content: fmt.Sprintf("Hello from %s! You asked: %s", p.client.Name(), req.Messages[len(req.Messages)-1].Content),
				},
				FinishReason: &[]string{"stop"}[0],
			},
		},
		Usage: provider.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}, nil
}

func (p *customProvider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
	return nil, fmt.Errorf("streaming not implemented in custom provider demo")
}

func (p *customProvider) Close() error {
	return nil
}

func main() {
	fmt.Println("=== 3rd Party Custom Provider Example ===")

	// Create a custom provider (this could be from an external package)
	customProv := NewCustomProvider("MyCustomLLM", "custom-api-key")

	// Inject the custom provider directly into omnillm
	client, err := omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{
			{CustomProvider: customProv}, // Direct provider injection!
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Printf("Using custom provider: %s\n", client.Provider().Name())

	// Use the same omnillm API with the custom provider
	response, err := client.CreateChatCompletion(context.Background(), &omnillm.ChatCompletionRequest{
		Model: "custom-model-v1",
		Messages: []omnillm.Message{
			{
				Role:    omnillm.RoleUser,
				Content: "Hello from a 3rd party provider!",
			},
		},
		MaxTokens:   &[]int{50}[0],
		Temperature: &[]float64{0.7}[0],
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)

	fmt.Println()
	fmt.Println("🎉 3rd party providers can now extend omnillm without modifying core!")
	fmt.Println("   - No registry needed")
	fmt.Println("   - No global state")
	fmt.Println("   - Compile-time type safety")
	fmt.Println("   - Clean dependency injection")
}
