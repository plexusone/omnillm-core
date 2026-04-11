# Custom Providers

External packages can create providers without modifying the core library. This is the recommended approach for adding new LLM backends.

## Creating a Provider

### Step 1: Implement the Provider Interface

```go
// In your external package (e.g., github.com/yourname/omnillm-myprovider)
package myprovider

import (
    "context"
    "github.com/plexusone/omnillm-core/provider"
)

// HTTP Client
type Client struct {
    apiKey string
    // your HTTP client implementation
}

func New(apiKey string) *Client {
    return &Client{apiKey: apiKey}
}

// Provider Adapter
type Provider struct {
    client *Client
}

func NewProvider(apiKey string) provider.Provider {
    return &Provider{client: New(apiKey)}
}

func (p *Provider) Name() string { return "myprovider" }
func (p *Provider) Close() error { return p.client.Close() }

func (p *Provider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
    // Convert provider.ChatCompletionRequest to your API format
    // Make HTTP call via p.client
    // Convert response back to provider.ChatCompletionResponse
}

func (p *Provider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
    // Your streaming implementation
}
```

### Step 2: Use Your Provider

```go
import (
    "github.com/plexusone/omnillm-core"
    "github.com/yourname/omnillm-myprovider"
)

func main() {
    customProvider := myprovider.NewProvider("your-api-key")

    client, err := omnillm.NewClient(omnillm.ClientConfig{
        Providers: []omnillm.ProviderConfig{
            {CustomProvider: customProvider},
        },
    })

    // Use the same omnillm API
    response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
        Model: "my-model",
        Messages: []omnillm.Message{{Role: omnillm.RoleUser, Content: "Hello!"}},
    })
}
```

## Provider Interface

```go
type Provider interface {
    // Name returns the provider name (e.g., "openai", "anthropic")
    Name() string

    // CreateChatCompletion performs a synchronous chat completion
    CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)

    // CreateChatCompletionStream performs a streaming chat completion
    CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error)

    // Close releases any resources
    Close() error
}
```

## Reference Implementations

Look at any built-in provider as a reference:

- `providers/openai/` - OpenAI implementation
- `providers/anthropic/` - Anthropic implementation
- [omnillm-bedrock](https://github.com/plexusone/omnillm-bedrock) - External provider example

## Benefits

- **No Core Changes**: External providers don't require modifying the core library
- **Clean Injection**: Use `ProviderConfig.CustomProvider` to inject your provider
- **Same Interface**: Both internal and external providers use the same `provider.Provider` interface
- **Easy Testing**: Mock the provider interface for unit tests
