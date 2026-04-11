# Architecture

OmniLLM uses a clean, modular architecture that separates concerns and enables easy extensibility.

## Directory Structure

```
omnillm/
├── client.go            # Main ChatClient wrapper
├── providers.go         # Factory functions for built-in providers
├── types.go             # Type aliases for backward compatibility
├── memory.go            # Conversation memory management
├── observability.go     # ObservabilityHook interface
├── errors.go            # Unified error handling
├── *_test.go            # Comprehensive unit tests
├── provider/            # Public interface package for external providers
│   ├── interface.go     # Provider interface
│   └── types.go         # Unified request/response types
├── providers/           # Individual provider packages
│   ├── openai/
│   │   ├── openai.go    # HTTP client
│   │   ├── types.go     # OpenAI-specific types
│   │   ├── adapter.go   # provider.Provider implementation
│   │   └── *_test.go    # Provider tests
│   ├── anthropic/
│   ├── gemini/
│   ├── xai/
│   └── ollama/
└── testing/             # Test utilities
    └── mock_kvs.go      # Mock KVS for memory testing
```

## Key Architecture Benefits

- **Public Interface**: The `provider` package exports the `Provider` interface that external packages can implement
- **Reference Implementation**: Internal providers follow the exact same structure that external providers should use
- **Direct Injection**: External providers are injected via `ClientConfig.CustomProvider` without modifying core code
- **Modular Design**: Each provider is self-contained with its own HTTP client, types, and adapter
- **Testable**: Clean interfaces that can be easily mocked and tested
- **Extensible**: New providers can be added without touching existing code
- **Native Implementation**: Uses standard `net/http` for direct API communication (no official SDK dependencies)

## Provider Interface

```go
type Provider interface {
    // Name returns the provider name
    Name() string

    // CreateChatCompletion performs a synchronous chat completion
    CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)

    // CreateChatCompletionStream performs a streaming chat completion
    CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error)

    // Close releases any resources
    Close() error
}
```

## Request/Response Flow

```
┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐
│   Client   │───►│  Provider  │───►│   HTTP     │───►│  LLM API   │
│   Code     │    │  Adapter   │    │   Client   │    │  Backend   │
└────────────┘    └────────────┘    └────────────┘    └────────────┘
      │                 │                 │                 │
      │  Unified API    │  Convert to     │  Native HTTP    │
      │  Request        │  Provider       │  Request        │
      │                 │  Format         │                 │
      ▼                 ▼                 ▼                 ▼
```

## Adding New Providers

### External (Recommended)

1. Create a new Go module
2. Import `github.com/plexusone/omnillm-core/provider`
3. Implement the `Provider` interface
4. Users inject via `CustomProvider`

### Built-in (Core Contributors)

1. Create provider package: `providers/newprovider/`
2. Implement HTTP client, types, and adapter
3. Add factory function in `providers.go`
4. Add provider constant
