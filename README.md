# OmniLLM - Unified Go SDK for Large Language Models

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/omnillm/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/omnillm/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/omnillm/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/omnillm/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/omnillm/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/omnillm/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/omnillm
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/omnillm
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/omnillm
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/omnillm
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fomnillm
 [loc-svg]: https://tokei.rs/b1/github/plexusone/omnillm
 [repo-url]: https://github.com/plexusone/omnillm
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/omnillm/blob/master/LICENSE

OmniLLM is a unified Go SDK that provides a consistent interface for interacting with multiple Large Language Model (LLM) providers including OpenAI, Anthropic (Claude), Google Gemini, X.AI (Grok), GLM (Zhipu AI), Kimi (Moonshot AI), Qwen (Alibaba Cloud), and Ollama. It implements the Chat Completions API pattern and offers both synchronous and streaming capabilities. Additional providers like AWS Bedrock are available as [external modules](#external-providers).

## ✨ Features

- **🔌 Multi-Provider Support**: OpenAI, Anthropic (Claude), Google Gemini, X.AI (Grok), GLM (Zhipu AI), Kimi (Moonshot AI), Qwen (Alibaba Cloud), Ollama, plus [external providers](#external-providers) (AWS Bedrock, etc.)
- **🎯 Unified API**: Same interface across all providers
- **📡 Streaming Support**: Real-time response streaming for all providers
- **🧠 Conversation Memory**: Persistent conversation history using Key-Value Stores
- **🔀 Fallback Providers**: Automatic failover to backup providers when primary fails
- **⚡ Circuit Breaker**: Prevent cascading failures by temporarily skipping unhealthy providers
- **🔢 Token Estimation**: Pre-flight token counting to validate requests before sending
- **💾 Response Caching**: Cache identical requests with configurable TTL to reduce costs
- **📊 Observability Hooks**: Extensible hooks for tracing, logging, and metrics without modifying core library
- **🔄 Retry with Backoff**: Automatic retries for transient failures (rate limits, 5xx errors)
- **🧪 Comprehensive Testing**: Unit tests, integration tests, and mock implementations included
- **🔧 Extensible**: Easy to add new LLM providers
- **📦 Modular**: Provider-specific implementations in separate packages
- **🏗️ Reference Architecture**: Internal providers serve as reference implementations for external providers
- **🔌 3rd Party Friendly**: External providers can be injected without modifying core library
- **⚡ Type Safe**: Full Go type safety with comprehensive error handling

## 🏗️ Architecture

OmniLLM uses a clean, modular architecture that separates concerns and enables easy extensibility:

```
omnillm/
├── client.go            # Main ChatClient wrapper
├── providers.go         # Factory functions for built-in providers
├── types.go             # Type aliases for backward compatibility
├── memory.go            # Conversation memory management
├── observability.go     # ObservabilityHook interface for tracing/logging/metrics
├── errors.go            # Unified error handling
├── *_test.go            # Comprehensive unit tests
├── provider/            # 🎯 Public interface package for external providers
│   ├── interface.go     # Provider interface that all providers must implement
│   └── types.go         # Unified request/response types
├── providers/           # 📦 Individual provider packages (reference implementations)
│   ├── openai/          # OpenAI implementation
│   │   ├── openai.go    # HTTP client
│   │   ├── types.go     # OpenAI-specific types
│   │   ├── adapter.go   # provider.Provider implementation
│   │   └── *_test.go    # Provider tests
│   ├── anthropic/       # Anthropic implementation
│   │   ├── anthropic.go # HTTP client (SSE streaming)
│   │   ├── types.go     # Anthropic-specific types
│   │   ├── adapter.go   # provider.Provider implementation
│   │   └── *_test.go    # Provider and integration tests
│   ├── gemini/          # Google Gemini implementation
│   ├── xai/             # X.AI Grok implementation
│   ├── glm/             # Zhipu AI GLM implementation
│   ├── kimi/            # Moonshot AI Kimi implementation
│   ├── qwen/            # Alibaba Cloud Qwen implementation
│   └── ollama/          # Ollama implementation
└── testing/             # 🧪 Test utilities
    └── mock_kvs.go      # Mock KVS for memory testing
```

### Key Architecture Benefits

- **🎯 Public Interface**: The `provider` package exports the `Provider` interface that external packages can implement
- **🏗️ Reference Implementation**: Internal providers follow the exact same structure that external providers should use
- **🔌 Direct Injection**: External providers are injected via `ClientConfig.CustomProvider` without modifying core code
- **📦 Modular Design**: Each provider is self-contained with its own HTTP client, types, and adapter
- **🧪 Testable**: Clean interfaces that can be easily mocked and tested
- **🔧 Extensible**: New providers can be added without touching existing code
- **⚡ Native Implementation**: Uses standard `net/http` for direct API communication (no official SDK dependencies)

## 🚀 Quick Start

### Installation

```bash
go get github.com/plexusone/omnillm
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/plexusone/omnillm"
)

func main() {
    // Create a client for OpenAI
    client, err := omnillm.NewClient(omnillm.ClientConfig{
        Providers: []omnillm.ProviderConfig{
            {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-openai-api-key"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create a chat completion request
    response, err := client.CreateChatCompletion(context.Background(), &omnillm.ChatCompletionRequest{
        Model: omnillm.ModelGPT4o,
        Messages: []omnillm.Message{
            {
                Role:    omnillm.RoleUser,
                Content: "Hello! How can you help me today?",
            },
        },
        MaxTokens:   &[]int{150}[0],
        Temperature: &[]float64{0.7}[0],
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
    fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)
}
```

## 🔧 Supported Providers

### OpenAI

- **Models**: GPT-5, GPT-4.1, GPT-4o, GPT-4o-mini, GPT-4-turbo, GPT-3.5-turbo
- **Features**: Chat completions, streaming, function calling

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-openai-api-key"},
    },
})
```

### Anthropic (Claude)

- **Models**: Claude-Opus-4.1, Claude-Opus-4, Claude-Sonnet-4, Claude-3.7-Sonnet, Claude-3.5-Haiku, Claude-3-Opus, Claude-3-Sonnet, Claude-3-Haiku
- **Features**: Chat completions, streaming, system message support

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameAnthropic, APIKey: "your-anthropic-api-key"},
    },
})
```

### Google Gemini

- **Models**: Gemini-2.5-Pro, Gemini-2.5-Flash, Gemini-1.5-Pro, Gemini-1.5-Flash
- **Features**: Chat completions, streaming

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameGemini, APIKey: "your-gemini-api-key"},
    },
})
```

### AWS Bedrock (External Provider)

AWS Bedrock is available as an external module to avoid pulling AWS SDK dependencies for users who don't need it.

```bash
go get github.com/plexusone/omnillm-bedrock
```

```go
import (
    "github.com/plexusone/omnillm"
    "github.com/plexusone/omnillm-bedrock"
)

// Create the Bedrock provider
bedrockProvider, err := bedrock.NewProvider("us-east-1")
if err != nil {
    log.Fatal(err)
}

// Use it with omnillm via CustomProvider
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {CustomProvider: bedrockProvider},
    },
})
```

See [External Providers](#external-providers) for more details.

### X.AI (Grok)

- **Models**: Grok-4.1-Fast (Reasoning/Non-Reasoning), Grok-4 (0709), Grok-4-Fast (Reasoning/Non-Reasoning), Grok-Code-Fast, Grok-3, Grok-3-Mini, Grok-2, Grok-2-Vision
- **Features**: Chat completions, streaming, OpenAI-compatible API, 2M context window (4.1/4-Fast models)

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameXAI, APIKey: "your-xai-api-key"},
    },
})
```

### GLM (Zhipu AI)

- **Models**: GLM-5, GLM-4.7, GLM-4.6, GLM-4.5 series
- **Features**: Chat completions, streaming, OpenAI-compatible API, thinking modes, up to 200K context

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameGLM, APIKey: "your-glm-api-key"},
    },
})
```

### Kimi (Moonshot AI)

- **Models**: Kimi K2.5, Kimi K2 series, Moonshot V1 (8k/32k/128k)
- **Features**: Chat completions, streaming, OpenAI-compatible API, up to 256K context, vision support

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameKimi, APIKey: "your-kimi-api-key"},
    },
})
```

### Qwen (Alibaba Cloud)

- **Models**: Qwen3 Max, QwQ, Qwen3.5, Qwen3, Qwen2.5 series
- **Features**: Chat completions, streaming, OpenAI-compatible API, thinking modes, wide global availability

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameQwen, APIKey: "your-qwen-api-key"},
    },
})
```

### Ollama (Local Models)

- **Models**: Llama 3, Mistral, CodeLlama, Gemma, Qwen2.5, DeepSeek-Coder
- **Features**: Local inference, no API keys required, optimized for Apple Silicon

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOllama, BaseURL: "http://localhost:11434"},
    },
})
```

## 🔌 External Providers

Some providers with heavy SDK dependencies are available as separate modules to keep the core library lightweight. These are injected via `ClientConfig.CustomProvider`.

| Provider | Module | Why External |
|----------|--------|--------------|
| AWS Bedrock | [github.com/plexusone/omnillm-bedrock](https://github.com/plexusone/omnillm-bedrock) | AWS SDK v2 adds 17+ transitive dependencies |

### Using External Providers

```go
import (
    "github.com/plexusone/omnillm"
    "github.com/plexusone/omnillm-bedrock"  // or your custom provider
)

// Create the external provider
bedrockProv, err := bedrock.NewProvider("us-east-1")
if err != nil {
    log.Fatal(err)
}

// Inject via CustomProvider in Providers slice
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {CustomProvider: bedrockProv},
    },
})
```

### Creating Your Own External Provider

External providers implement the `provider.Provider` interface:

```go
import "github.com/plexusone/omnillm/provider"

type MyProvider struct{}

func (p *MyProvider) Name() string { return "myprovider" }
func (p *MyProvider) Close() error { return nil }

func (p *MyProvider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
    // Your implementation
}

func (p *MyProvider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
    // Your streaming implementation
}
```

See the [omnillm-bedrock](https://github.com/plexusone/omnillm-bedrock) source code as a reference implementation.

## 📡 Streaming Example

```go
stream, err := client.CreateChatCompletionStream(context.Background(), &omnillm.ChatCompletionRequest{
    Model: omnillm.ModelGPT4o,
    Messages: []omnillm.Message{
        {
            Role:    omnillm.RoleUser,
            Content: "Tell me a short story about AI.",
        },
    },
    MaxTokens:   &[]int{200}[0],
    Temperature: &[]float64{0.8}[0],
})
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

fmt.Print("AI Response: ")
for {
    chunk, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    
    if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
        fmt.Print(chunk.Choices[0].Delta.Content)
    }
}
fmt.Println()
```

## 🧠 Conversation Memory

OmniLLM supports persistent conversation memory using any Key-Value Store that implements the [Sogo KVS interface](https://github.com/grokify/sogo/blob/master/database/kvs/definitions.go). This enables multi-turn conversations that persist across application restarts.

### Memory Configuration

```go
// Configure memory settings
memoryConfig := omnillm.MemoryConfig{
    MaxMessages: 50,                    // Keep last 50 messages per session
    TTL:         24 * time.Hour,       // Messages expire after 24 hours
    KeyPrefix:   "myapp:conversations", // Custom key prefix
}

// Create client with memory (using Redis, DynamoDB, etc.)
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-api-key"},
    },
    Memory:       kvsClient,          // Your KVS implementation
    MemoryConfig: &memoryConfig,
})
```

### Memory-Aware Completions

```go
// Create a session with system message
err = client.CreateConversationWithSystemMessage(ctx, "user-123", 
    "You are a helpful assistant that remembers our conversation history.")

// Use memory-aware completion - automatically loads conversation history
response, err := client.CreateChatCompletionWithMemory(ctx, "user-123", &omnillm.ChatCompletionRequest{
    Model: omnillm.ModelGPT4o,
    Messages: []omnillm.Message{
        {Role: omnillm.RoleUser, Content: "What did we discuss last time?"},
    },
    MaxTokens: &[]int{200}[0],
})

// The response will include context from previous conversations in this session
```

### Memory Management

```go
// Load conversation history
conversation, err := client.LoadConversation(ctx, "user-123")

// Get just the messages
messages, err := client.GetConversationMessages(ctx, "user-123")

// Manually append messages
err = client.AppendMessage(ctx, "user-123", omnillm.Message{
    Role:    omnillm.RoleUser,
    Content: "Remember this important fact: I prefer JSON responses.",
})

// Delete conversation
err = client.DeleteConversation(ctx, "user-123")
```

### KVS Backend Support

Memory works with any KVS implementation:
- **Redis**: For high-performance, distributed memory
- **DynamoDB**: For AWS-native storage
- **In-Memory**: For testing and development
- **Custom**: Any implementation of the Sogo KVS interface

```go
// Example with Redis (using a hypothetical Redis KVS implementation)
redisKVS := redis.NewKVSClient("localhost:6379")
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-key"},
    },
    Memory: redisKVS,
})
```

## 📊 Observability Hooks

OmniLLM supports observability hooks that allow you to add tracing, logging, and metrics to LLM calls without modifying the core library. This is useful for integrating with observability platforms like OpenTelemetry, Datadog, or custom monitoring solutions.

### ObservabilityHook Interface

```go
// LLMCallInfo provides metadata about the LLM call
type LLMCallInfo struct {
    CallID       string    // Unique identifier for correlating BeforeRequest/AfterResponse
    ProviderName string    // e.g., "openai", "anthropic"
    StartTime    time.Time // When the call started
}

// ObservabilityHook allows external packages to observe LLM calls
type ObservabilityHook interface {
    // BeforeRequest is called before each LLM call.
    // Returns a new context for trace/span propagation.
    BeforeRequest(ctx context.Context, info LLMCallInfo, req *provider.ChatCompletionRequest) context.Context

    // AfterResponse is called after each LLM call completes (success or failure).
    AfterResponse(ctx context.Context, info LLMCallInfo, req *provider.ChatCompletionRequest, resp *provider.ChatCompletionResponse, err error)

    // WrapStream wraps a stream for observability of streaming responses.
    // Note: AfterResponse is only called if stream creation fails. For streaming
    // completion timing, handle Close() or EOF detection in your wrapper.
    WrapStream(ctx context.Context, info LLMCallInfo, req *provider.ChatCompletionRequest, stream provider.ChatCompletionStream) provider.ChatCompletionStream
}
```

### Basic Usage

```go
// Create a simple logging hook
type LoggingHook struct{}

func (h *LoggingHook) BeforeRequest(ctx context.Context, info omnillm.LLMCallInfo, req *omnillm.ChatCompletionRequest) context.Context {
    log.Printf("[%s] LLM call started: provider=%s model=%s", info.CallID, info.ProviderName, req.Model)
    return ctx
}

func (h *LoggingHook) AfterResponse(ctx context.Context, info omnillm.LLMCallInfo, req *omnillm.ChatCompletionRequest, resp *omnillm.ChatCompletionResponse, err error) {
    duration := time.Since(info.StartTime)
    if err != nil {
        log.Printf("[%s] LLM call failed: provider=%s duration=%v error=%v", info.CallID, info.ProviderName, duration, err)
    } else {
        log.Printf("[%s] LLM call completed: provider=%s duration=%v tokens=%d", info.CallID, info.ProviderName, duration, resp.Usage.TotalTokens)
    }
}

func (h *LoggingHook) WrapStream(ctx context.Context, info omnillm.LLMCallInfo, req *omnillm.ChatCompletionRequest, stream omnillm.ChatCompletionStream) omnillm.ChatCompletionStream {
    return stream // Return unwrapped for simple logging, or wrap for streaming metrics
}

// Use the hook when creating a client
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-api-key"},
    },
    ObservabilityHook: &LoggingHook{},
})
```

### OpenTelemetry Integration Example

```go
type OTelHook struct {
    tracer trace.Tracer
}

func (h *OTelHook) BeforeRequest(ctx context.Context, info omnillm.LLMCallInfo, req *omnillm.ChatCompletionRequest) context.Context {
    ctx, span := h.tracer.Start(ctx, "llm.chat_completion",
        trace.WithAttributes(
            attribute.String("llm.provider", info.ProviderName),
            attribute.String("llm.model", req.Model),
        ),
    )
    return ctx
}

func (h *OTelHook) AfterResponse(ctx context.Context, info omnillm.LLMCallInfo, req *omnillm.ChatCompletionRequest, resp *omnillm.ChatCompletionResponse, err error) {
    span := trace.SpanFromContext(ctx)
    defer span.End()

    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
    } else if resp != nil {
        span.SetAttributes(
            attribute.Int("llm.tokens.total", resp.Usage.TotalTokens),
            attribute.Int("llm.tokens.prompt", resp.Usage.PromptTokens),
            attribute.Int("llm.tokens.completion", resp.Usage.CompletionTokens),
        )
    }
}

func (h *OTelHook) WrapStream(ctx context.Context, info omnillm.LLMCallInfo, req *omnillm.ChatCompletionRequest, stream omnillm.ChatCompletionStream) omnillm.ChatCompletionStream {
    return &observableStream{stream: stream, ctx: ctx, info: info}
}
```

### Key Benefits

- **Non-Invasive**: Add observability without modifying core library code
- **Provider Agnostic**: Works with all LLM providers (OpenAI, Anthropic, Gemini, etc.)
- **Streaming Support**: Wrap streams to observe streaming responses
- **Context Propagation**: Pass trace context through the entire call chain
- **Flexible**: Implement only the methods you need; all are called if the hook is set

## 🔀 Fallback Providers

OmniLLM supports automatic failover to backup providers when the primary provider fails. Fallback only triggers on retryable errors (rate limits, server errors, network issues) - authentication errors and invalid requests do not trigger fallback.

### Basic Usage

```go
// Providers[0] is primary, Providers[1+] are fallbacks
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "openai-key"},       // Primary
        {Provider: omnillm.ProviderNameAnthropic, APIKey: "anthropic-key"}, // Fallback 1
        {Provider: omnillm.ProviderNameGemini, APIKey: "gemini-key"},       // Fallback 2
    },
})

// If OpenAI fails with a retryable error, automatically tries Anthropic, then Gemini
response, err := client.CreateChatCompletion(ctx, request)
```

### With Circuit Breaker

Enable circuit breaker to temporarily skip providers that are failing repeatedly:

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "openai-key"},
        {Provider: omnillm.ProviderNameAnthropic, APIKey: "anthropic-key"},
    },
    CircuitBreakerConfig: &omnillm.CircuitBreakerConfig{
        FailureThreshold: 5,               // Open after 5 consecutive failures
        SuccessThreshold: 2,               // Close after 2 successes in half-open
        Timeout:          30 * time.Second, // Wait before trying again
    },
})
```

### Error Classification

Fallback uses intelligent error classification:

| Error Type | Triggers Fallback |
|------------|-------------------|
| Rate limits (429) | ✅ Yes |
| Server errors (5xx) | ✅ Yes |
| Network errors | ✅ Yes |
| Auth errors (401/403) | ❌ No |
| Invalid requests (400) | ❌ No |

## ⚡ Circuit Breaker

The circuit breaker pattern prevents cascading failures by temporarily skipping providers that are unhealthy.

### States

- **Closed**: Normal operation, requests flow through
- **Open**: Provider is failing, requests skip it immediately
- **Half-Open**: Testing if provider has recovered

### Configuration

```go
cbConfig := &omnillm.CircuitBreakerConfig{
    FailureThreshold:     5,               // Failures before opening
    SuccessThreshold:     2,               // Successes to close from half-open
    Timeout:              30 * time.Second, // Wait before half-open
    FailureRateThreshold: 0.5,             // 50% failure rate opens circuit
    MinimumRequests:      10,              // Minimum requests for rate calculation
}
```

## 🔢 Token Estimation

OmniLLM provides pre-flight token estimation to validate requests before sending them to the API. This helps avoid hitting context window limits.

### Basic Usage

```go
// Create estimator with default config
estimator := omnillm.NewTokenEstimator(omnillm.DefaultTokenEstimatorConfig())

// Estimate tokens for messages
tokens, err := estimator.EstimateTokens("gpt-4o", messages)

// Get model's context window
window := estimator.GetContextWindow("gpt-4o") // Returns 128000
```

### Automatic Validation

Enable automatic token validation in client:

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-key"},
    },
    TokenEstimator: omnillm.NewTokenEstimator(omnillm.DefaultTokenEstimatorConfig()),
    ValidateTokens: true, // Rejects requests that exceed context window
})

// Returns TokenLimitError if request exceeds model limits
response, err := client.CreateChatCompletion(ctx, request)
if tlErr, ok := err.(*omnillm.TokenLimitError); ok {
    fmt.Printf("Request has %d tokens, but model only supports %d\n",
        tlErr.EstimatedTokens, tlErr.ContextWindow)
}
```

### Built-in Context Windows

Token estimator includes context windows for 40+ models:

| Provider | Models | Context Window |
|----------|--------|----------------|
| OpenAI | GPT-4o, GPT-4o-mini | 128,000 |
| OpenAI | o1 | 200,000 |
| Anthropic | Claude 3/3.5/4 | 200,000 |
| Google | Gemini 2.5 | 1,000,000 |
| Google | Gemini 1.5 Pro | 2,000,000 |
| X.AI | Grok 3/4 | 128,000 |

### Custom Configuration

```go
config := omnillm.TokenEstimatorConfig{
    CharactersPerToken: 3.5, // More conservative estimate
    CustomContextWindows: map[string]int{
        "my-custom-model": 500000,
        "gpt-4o":          200000, // Override built-in
    },
}
estimator := omnillm.NewTokenEstimator(config)
```

## 💾 Response Caching

OmniLLM supports response caching to reduce API costs for identical requests. Caching uses the same KVS backend as conversation memory.

### Basic Usage

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-key"},
    },
    Cache: kvsClient, // Your KVS implementation (Redis, DynamoDB, etc.)
    CacheConfig: &omnillm.CacheConfig{
        TTL:       1 * time.Hour,        // Cache duration
        KeyPrefix: "myapp:llm-cache",    // Key prefix in KVS
    },
})

// First call hits the API
response1, _ := client.CreateChatCompletion(ctx, request)

// Second identical call returns cached response
response2, _ := client.CreateChatCompletion(ctx, request)

// Check if response was from cache
if response2.ProviderMetadata["cache_hit"] == true {
    fmt.Println("Response was cached!")
}
```

### Cache Configuration

```go
cacheConfig := &omnillm.CacheConfig{
    TTL:                1 * time.Hour,       // Time-to-live
    KeyPrefix:          "omnillm:cache",     // Key prefix
    SkipStreaming:      true,                // Don't cache streaming (default)
    CacheableModels:    []string{"gpt-4o"},  // Only cache specific models (nil = all)
    IncludeTemperature: true,                // Temperature affects cache key
    IncludeSeed:        true,                // Seed affects cache key
}
```

### Cache Key Generation

Cache keys are generated from a SHA-256 hash of:

- Model name
- Messages (role, content, name, tool_call_id)
- MaxTokens, Temperature, TopP, TopK, Seed, Stop sequences

Different parameter values = different cache keys.

## 🔄 Provider Switching

The unified interface makes it easy to switch between providers:

```go
// Same request works with any provider
request := &omnillm.ChatCompletionRequest{
    Model: omnillm.ModelGPT4o, // or omnillm.ModelClaude3Sonnet, etc.
    Messages: []omnillm.Message{
        {Role: omnillm.RoleUser, Content: "Hello, world!"},
    },
    MaxTokens: &[]int{100}[0],
}

// OpenAI
openaiClient, _ := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "openai-key"},
    },
})

// Anthropic
anthropicClient, _ := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameAnthropic, APIKey: "anthropic-key"},
    },
})

// Gemini
geminiClient, _ := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameGemini, APIKey: "gemini-key"},
    },
})

// Same API call for all providers
response1, _ := openaiClient.CreateChatCompletion(ctx, request)
response2, _ := anthropicClient.CreateChatCompletion(ctx, request)
response3, _ := geminiClient.CreateChatCompletion(ctx, request)
```

## 🧪 Testing

OmniLLM includes a comprehensive test suite with both unit tests and integration tests.

### Running Tests

```bash
# Run all unit tests (no API keys required)
go test ./... -short

# Run with coverage
go test ./... -short -cover

# Run integration tests (requires API keys)
ANTHROPIC_API_KEY=your-key go test ./providers/anthropic -v
OPENAI_API_KEY=your-key go test ./providers/openai -v
XAI_API_KEY=your-key go test ./providers/xai -v
GLM_API_KEY=your-key go test ./providers/glm -v
KIMI_API_KEY=your-key go test ./providers/kimi -v
QWEN_API_KEY=your-key go test ./providers/qwen -v

# Run all tests including integration
ANTHROPIC_API_KEY=your-key OPENAI_API_KEY=your-key XAI_API_KEY=your-key GLM_API_KEY=your-key KIMI_API_KEY=your-key QWEN_API_KEY=your-key go test ./... -v
```

### Test Coverage

- **Unit Tests**: Mock-based tests that run without external dependencies
- **Integration Tests**: Real API tests that skip gracefully when API keys are not set
- **Memory Tests**: Comprehensive conversation memory management tests
- **Provider Tests**: Adapter logic, message conversion, and streaming tests

### Writing Tests

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

### Conditional Integration Tests

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

### Mock KVS for Memory Testing

OmniLLM provides a mock KVS implementation for testing memory functionality:

```go
import omnillmtest "github.com/plexusone/omnillm/testing"

// Create mock KVS for testing
mockKVS := omnillmtest.NewMockKVS()

client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "test-key"},
    },
    Memory: mockKVS,
})
```

## 📚 Examples

The repository includes comprehensive examples:

- **Basic Usage**: Simple chat completions with each provider
- **Streaming**: Real-time response handling
- **Conversation**: Multi-turn conversations with context
- **Memory Demo**: Persistent conversation memory with KVS backend
- **Architecture Demo**: Overview of the provider architecture
- **Custom Provider**: How to create and use 3rd party providers

Run examples:
```bash
go run examples/basic/main.go
go run examples/streaming/main.go
go run examples/anthropic_streaming/main.go
go run examples/conversation/main.go
go run examples/memory_demo/main.go
go run examples/providers_demo/main.go
go run examples/xai/main.go
go run examples/ollama/main.go
go run examples/ollama_streaming/main.go
go run examples/gemini/main.go
go run examples/custom_provider/main.go
```

## 🔧 Configuration

### Environment Variables

- `OPENAI_API_KEY`: Your OpenAI API key
- `ANTHROPIC_API_KEY`: Your Anthropic API key
- `GEMINI_API_KEY`: Your Google Gemini API key
- `XAI_API_KEY`: Your X.AI API key
- `GLM_API_KEY`: Your Zhipu AI GLM API key
- `KIMI_API_KEY`: Your Moonshot AI Kimi API key
- `QWEN_API_KEY`: Your Alibaba Cloud Qwen API key

### Advanced Configuration

```go
config := omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {
            Provider: omnillm.ProviderNameOpenAI,
            APIKey:   "your-api-key",
            BaseURL:  "https://custom-endpoint.com/v1",
            Extra: map[string]any{
                "timeout": 60, // Custom provider-specific settings
            },
        },
    },
}
```

### Request Parameters

`ChatCompletionRequest` supports the following parameters with provider-specific availability:

| Parameter | Type | Providers | Description |
|-----------|------|-----------|-------------|
| `Model` | `string` | All | Model identifier (required) |
| `Messages` | `[]Message` | All | Conversation messages (required) |
| `MaxTokens` | `*int` | All | Maximum tokens to generate |
| `Temperature` | `*float64` | All | Randomness (0.0-2.0) |
| `TopP` | `*float64` | All | Nucleus sampling threshold |
| `TopK` | `*int` | Anthropic, Gemini, Ollama | Top K token selection |
| `Stop` | `[]string` | All | Stop sequences |
| `PresencePenalty` | `*float64` | OpenAI, X.AI | Penalize tokens by presence |
| `FrequencyPenalty` | `*float64` | OpenAI, X.AI | Penalize tokens by frequency |
| `Seed` | `*int` | OpenAI, X.AI, Ollama | Reproducible outputs |
| `N` | `*int` | OpenAI | Number of completions |
| `ResponseFormat` | `*ResponseFormat` | OpenAI, Gemini | JSON mode (`{"type": "json_object"}`) |
| `Logprobs` | `*bool` | OpenAI | Return log probabilities |
| `TopLogprobs` | `*int` | OpenAI | Top logprobs count (0-20) |
| `User` | `*string` | OpenAI | End-user identifier |
| `LogitBias` | `map[string]int` | OpenAI | Token bias adjustments |

```go
// Helper for pointer values
func ptr[T any](v T) *T { return &v }

// Example: Reproducible outputs with seed
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model:    omnillm.ModelGPT4o,
    Messages: messages,
    Seed:     ptr(42), // Same seed = same output
})

// Example: JSON mode response
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model:    omnillm.ModelGPT4o,
    Messages: messages,
    ResponseFormat: &omnillm.ResponseFormat{Type: "json_object"},
})

// Example: TopK sampling (Anthropic/Gemini/Ollama)
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model:    omnillm.ModelClaude3Sonnet,
    Messages: messages,
    TopK:     ptr(40), // Consider only top 40 tokens
})
```

### Logging Configuration

OmniLLM supports injectable logging via Go's standard `log/slog` package. If no logger is provided, a null logger is used (no output).

```go
import (
    "log/slog"
    "os"

    "github.com/plexusone/omnillm"
)

// Use a custom logger
logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-api-key"},
    },
    Logger: logger, // Optional: defaults to null logger if not provided
})

// Access the logger if needed
client.Logger().Info("client initialized", slog.String("provider", "openai"))
```

The logger is used internally for non-critical errors (e.g., memory save failures) that shouldn't interrupt the main request flow.

### Context-Aware Logging

OmniLLM supports request-scoped logging via context. This allows you to attach trace IDs, user IDs, or other request-specific attributes to all log output within a request:

```go
import (
    "log/slog"

    "github.com/plexusone/omnillm"
    "github.com/grokify/mogo/log/slogutil"
)

// Create a request-scoped logger with trace/user context
reqLogger := slog.Default().With(
    slog.String("trace_id", traceID),
    slog.String("user_id", userID),
    slog.String("request_id", requestID),
)

// Attach logger to context
ctx = slogutil.ContextWithLogger(ctx, reqLogger)

// All internal logging will now include trace_id, user_id, and request_id
response, err := client.CreateChatCompletionWithMemory(ctx, sessionID, req)
```

The context-aware logger is retrieved using `slogutil.LoggerFromContext(ctx, fallback)`, which returns the context logger if present, or falls back to the client's configured logger.

### Retry with Backoff

OmniLLM supports automatic retries for transient failures (rate limits, 5xx errors) via a custom HTTP client. This uses the `retryhttp` package from `github.com/grokify/mogo`.

```go
import (
    "net/http"
    "time"

    "github.com/plexusone/omnillm"
    "github.com/grokify/mogo/net/http/retryhttp"
)

// Create retry transport with exponential backoff
rt := retryhttp.NewWithOptions(
    retryhttp.WithMaxRetries(5),                           // Max 5 retries
    retryhttp.WithInitialBackoff(500 * time.Millisecond),  // Start with 500ms
    retryhttp.WithMaxBackoff(30 * time.Second),            // Cap at 30s
    retryhttp.WithOnRetry(func(attempt int, req *http.Request, resp *http.Response, err error, backoff time.Duration) {
        log.Printf("Retry attempt %d, waiting %v", attempt, backoff)
    }),
)

// Create client with retry-enabled HTTP client
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {
            Provider: omnillm.ProviderNameOpenAI,
            APIKey:   os.Getenv("OPENAI_API_KEY"),
            HTTPClient: &http.Client{
                Transport: rt,
                Timeout:   2 * time.Minute, // Allow time for retries
            },
        },
    },
})
```

**Retry Transport Features:**
| Feature | Default | Description |
|---------|---------|-------------|
| Max Retries | 3 | Maximum retry attempts |
| Initial Backoff | 1s | Starting backoff duration |
| Max Backoff | 30s | Cap on backoff duration |
| Backoff Multiplier | 2.0 | Exponential growth factor |
| Jitter | 10% | Randomness to prevent thundering herd |
| Retryable Status Codes | 429, 500, 502, 503, 504 | Rate limits + 5xx errors |

**Additional Options:**
- `WithRetryableStatusCodes(codes)` - Custom status codes to retry
- `WithShouldRetry(fn)` - Custom retry decision function
- `WithLogger(logger)` - Structured logging for retry events
- Respects `Retry-After` headers from API responses

**Provider Support:** Works with OpenAI, Anthropic, X.AI, and Ollama providers. Gemini and Bedrock use SDK clients with their own retry mechanisms.

## 🏗️ Adding New Providers

### 🎯 3rd Party Providers (Recommended)

External packages can create providers without modifying the core library. This is the recommended approach for most use cases:

#### Step 1: Create Your Provider Package

```go
// In your external package (e.g., github.com/yourname/omnillm-gemini)
package gemini

import (
    "context"
    "github.com/plexusone/omnillm/provider"
)

// Step 1: HTTP Client (like providers/openai/openai.go)
type Client struct {
    apiKey string
    // your HTTP client implementation
}

func New(apiKey string) *Client {
    return &Client{apiKey: apiKey}
}

// Step 2: Provider Adapter (like providers/openai/adapter.go)
type Provider struct {
    client *Client
}

func NewProvider(apiKey string) provider.Provider {
    return &Provider{client: New(apiKey)}
}

func (p *Provider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
    // Convert provider.ChatCompletionRequest to your API format
    // Make HTTP call via p.client
    // Convert response back to provider.ChatCompletionResponse
}

func (p *Provider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
    // Your streaming implementation
}

func (p *Provider) Close() error { return p.client.Close() }
func (p *Provider) Name() string { return "gemini" }
```

#### Step 2: Use Your Provider

```go
import (
    "github.com/plexusone/omnillm"
    "github.com/yourname/omnillm-gemini"
)

func main() {
    // Create your custom provider
    customProvider := gemini.NewProvider("your-api-key")

    // Inject it directly into omnillm - no core modifications needed!
    client, err := omnillm.NewClient(omnillm.ClientConfig{
        Providers: []omnillm.ProviderConfig{
            {CustomProvider: customProvider},
        },
    })

    // Use the same omnillm API
    response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
        Model: "gemini-pro",
        Messages: []omnillm.Message{{Role: omnillm.RoleUser, Content: "Hello!"}},
    })
}
```

### 🔧 Built-in Providers (For Core Contributors)

To add a built-in provider to the core library, follow the same structure as existing providers:

1. **Create Provider Package**: `providers/newprovider/`
   - `newprovider.go` - HTTP client implementation
   - `types.go` - Provider-specific request/response types
   - `adapter.go` - `provider.Provider` interface implementation

2. **Update Core Files**:
   - Add factory function in `providers.go`
   - Add provider constant in `constants.go`
   - Add model constants if needed

3. **Reference Implementation**: Look at any existing provider (e.g., `providers/openai/`) as they all follow the exact same pattern that external providers should use

### 🎯 Why This Architecture?

- **🔌 No Core Changes**: External providers don't require modifying the core library
- **🏗️ Reference Pattern**: Internal providers demonstrate the exact structure external providers should follow
- **🧪 Easy Testing**: Both internal and external providers use the same `provider.Provider` interface
- **📦 Self-Contained**: Each provider manages its own HTTP client, types, and adapter logic
- **🔧 Direct Injection**: Clean dependency injection via `ProviderConfig.CustomProvider`

## 📊 Model Support

| Provider | Models | Features |
|----------|--------|----------|
| OpenAI | GPT-5, GPT-4.1, GPT-4o, GPT-4o-mini, GPT-4-turbo, GPT-3.5-turbo | Chat, Streaming, Functions |
| Anthropic | Claude-Opus-4.1, Claude-Opus-4, Claude-Sonnet-4, Claude-3.7-Sonnet, Claude-3.5-Haiku | Chat, Streaming, System messages |
| Gemini | Gemini-2.5-Pro, Gemini-2.5-Flash, Gemini-1.5-Pro, Gemini-1.5-Flash | Chat, Streaming |
| X.AI | Grok-4.1-Fast, Grok-4, Grok-4-Fast, Grok-Code-Fast, Grok-3, Grok-3-Mini, Grok-2 | Chat, Streaming, 2M context, Tool calling |
| GLM | GLM-5, GLM-4.7, GLM-4.6, GLM-4.5 series | Chat, Streaming, Thinking modes |
| Kimi | Kimi K2.5, K2 series, Moonshot V1 | Chat, Streaming, 256K context, Vision |
| Qwen | Qwen3 Max, QwQ, Qwen3.5, Qwen3, Qwen2.5 series | Chat, Streaming, Thinking modes |
| Ollama | Llama 3, Mistral, CodeLlama, Gemma, Qwen2.5, DeepSeek-Coder | Chat, Streaming, Local inference |
| Bedrock* | Claude models, Titan models | Chat, Multiple model families |

*Available as [external module](https://github.com/plexusone/omnillm-bedrock)

## 🚨 Error Handling

OmniLLM provides comprehensive error handling with provider-specific context:

```go
response, err := client.CreateChatCompletion(ctx, request)
if err != nil {
    if apiErr, ok := err.(*omnillm.APIError); ok {
        fmt.Printf("Provider: %s, Status: %d, Message: %s\n", 
            apiErr.Provider, apiErr.StatusCode, apiErr.Message)
    }
}
```

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests to ensure everything works:
   ```bash
   go test ./... -short        # Run unit tests
   go build ./...              # Verify build
   go vet ./...                # Run static analysis
   ```
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Adding Tests

When contributing new features:
- Add unit tests for core logic
- Add integration tests for provider implementations (with API key checks)
- Ensure tests pass without API keys using `-short` flag
- Mock external dependencies when possible

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Projects

- [Anthropic Go SDK](https://github.com/anthropics/anthropic-sdk-go) - Official Anthropic SDK
- [OpenAI Go SDK](https://github.com/openai/openai-go) - Official OpenAI SDK
- [AWS SDK for Go](https://github.com/aws/aws-sdk-go-v2) - Official AWS SDK

---

**Made with ❤️ for the Go and AI community**
