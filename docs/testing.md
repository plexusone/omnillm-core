# Testing

OmniLLM includes a comprehensive test suite with both unit tests and integration tests.

## Running Tests

```bash
# Run all unit tests (no API keys required)
go test ./... -short

# Run with coverage
go test ./... -short -cover

# Run integration tests (requires API keys)
ANTHROPIC_API_KEY=your-key go test ./providers/anthropic -v
OPENAI_API_KEY=your-key go test ./providers/openai -v
XAI_API_KEY=your-key go test ./providers/xai -v

# Run all tests including integration
ANTHROPIC_API_KEY=your-key OPENAI_API_KEY=your-key XAI_API_KEY=your-key go test ./... -v
```

## Test Coverage

- **Unit Tests**: Mock-based tests that run without external dependencies
- **Integration Tests**: Real API tests that skip gracefully when API keys are not set
- **Memory Tests**: Comprehensive conversation memory management tests
- **Provider Tests**: Adapter logic, message conversion, and streaming tests

## Writing Tests

The clean interface design makes testing straightforward:

```go
// Mock the Provider interface for testing
type mockProvider struct{}

func (m *mockProvider) CreateChatCompletion(ctx context.Context, req *omnillm.ChatCompletionRequest) (*omnillm.ChatCompletionResponse, error) {
    return &omnillm.ChatCompletionResponse{
        Choices: []omnillm.ChatCompletionChoice{
            {
                Message: omnillm.Message{
                    Role:    omnillm.RoleAssistant,
                    Content: "Mock response",
                },
            },
        },
    }, nil
}

func (m *mockProvider) CreateChatCompletionStream(ctx context.Context, req *omnillm.ChatCompletionRequest) (omnillm.ChatCompletionStream, error) {
    return nil, nil
}

func (m *mockProvider) Close() error { return nil }
func (m *mockProvider) Name() string { return "mock" }
```

## Conditional Integration Tests

Integration tests automatically skip when API keys are not available:

```go
func TestAnthropicIntegration_Streaming(t *testing.T) {
    apiKey := os.Getenv("ANTHROPIC_API_KEY")
    if apiKey == "" {
        t.Skip("Skipping integration test: ANTHROPIC_API_KEY not set")
    }
    // Test code here...
}
```

## Mock KVS for Memory Testing

OmniLLM provides a mock KVS implementation for testing memory functionality:

```go
import omnillmtest "github.com/plexusone/omnillm-core/testing"

// Create mock KVS for testing
mockKVS := omnillmtest.NewMockKVS()

client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "test-key"},
    },
    Memory: mockKVS,
})
```
