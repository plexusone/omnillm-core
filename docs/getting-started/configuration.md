# Configuration

## Client Configuration

```go
config := omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {
            Provider: omnillm.ProviderNameOpenAI,
            APIKey:   "your-api-key",
            BaseURL:  "https://custom-endpoint.com/v1", // Optional
            Extra: map[string]any{
                "timeout": 60, // Custom provider-specific settings
            },
        },
    },
}
```

## Request Parameters

`ChatCompletionRequest` supports the following parameters:

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
| `ResponseFormat` | `*ResponseFormat` | OpenAI, Gemini | JSON mode |
| `Logprobs` | `*bool` | OpenAI | Return log probabilities |
| `TopLogprobs` | `*int` | OpenAI | Top logprobs count (0-20) |
| `User` | `*string` | OpenAI | End-user identifier |
| `LogitBias` | `map[string]int` | OpenAI | Token bias adjustments |

## Example: Advanced Request

```go
// Helper for pointer values
func ptr[T any](v T) *T { return &v }

// Reproducible outputs with seed
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model:    omnillm.ModelGPT4o,
    Messages: messages,
    Seed:     ptr(42), // Same seed = same output
})

// JSON mode response
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model:    omnillm.ModelGPT4o,
    Messages: messages,
    ResponseFormat: &omnillm.ResponseFormat{Type: "json_object"},
})

// TopK sampling (Anthropic/Gemini/Ollama)
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model:    omnillm.ModelClaude3Sonnet,
    Messages: messages,
    TopK:     ptr(40), // Consider only top 40 tokens
})
```

## Logging Configuration

OmniLLM supports injectable logging via Go's standard `log/slog` package:

```go
import (
    "log/slog"
    "os"

    "github.com/plexusone/omnillm-core"
)

logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-api-key"},
    },
    Logger: logger,
})
```

## Context-Aware Logging

Attach request-scoped attributes to all log output:

```go
import "github.com/grokify/mogo/log/slogutil"

reqLogger := slog.Default().With(
    slog.String("trace_id", traceID),
    slog.String("user_id", userID),
)

ctx = slogutil.ContextWithLogger(ctx, reqLogger)

// All internal logging will now include trace_id and user_id
response, err := client.CreateChatCompletionWithMemory(ctx, sessionID, req)
```
