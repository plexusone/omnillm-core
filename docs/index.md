# OmniLLM

**Unified Go SDK for Large Language Models**

OmniLLM is a unified Go SDK that provides a consistent interface for interacting with multiple Large Language Model (LLM) providers including OpenAI, Anthropic (Claude), Google Gemini, X.AI (Grok), and Ollama. It implements the Chat Completions API pattern and offers both synchronous and streaming capabilities.

## Features

- **Multi-Provider Support**: OpenAI, Anthropic (Claude), Google Gemini, X.AI (Grok), Ollama, plus external providers (AWS Bedrock, etc.)
- **Unified API**: Same interface across all providers
- **Streaming Support**: Real-time response streaming for all providers
- **Conversation Memory**: Persistent conversation history using Key-Value Stores
- **Fallback Providers**: Automatic failover to backup providers when primary fails
- **Circuit Breaker**: Prevent cascading failures by temporarily skipping unhealthy providers
- **Token Estimation**: Pre-flight token counting to validate requests before sending
- **Response Caching**: Cache identical requests with configurable TTL to reduce costs
- **Observability Hooks**: Extensible hooks for tracing, logging, and metrics
- **Retry with Backoff**: Automatic retries for transient failures (rate limits, 5xx errors)
- **Tool Calling**: Function/tool calling support for agentic workflows
- **Type Safe**: Full Go type safety with comprehensive error handling

## Quick Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/plexusone/omnillm"
)

func main() {
    client, err := omnillm.NewClient(omnillm.ClientConfig{
        Providers: []omnillm.ProviderConfig{
            {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-openai-api-key"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    response, err := client.CreateChatCompletion(context.Background(), &omnillm.ChatCompletionRequest{
        Model: omnillm.ModelGPT4o,
        Messages: []omnillm.Message{
            {Role: omnillm.RoleUser, Content: "Hello! How can you help me today?"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
}
```

## Supported Providers

| Provider | Models | Features |
|----------|--------|----------|
| [OpenAI](providers/openai.md) | GPT-5, GPT-4.1, GPT-4o, GPT-4o-mini | Chat, Streaming, Tool Calling |
| [Anthropic](providers/anthropic.md) | Claude-Opus-4.1, Claude-Sonnet-4, Claude-3.7-Sonnet | Chat, Streaming, System messages |
| [Google Gemini](providers/gemini.md) | Gemini-2.5-Pro, Gemini-2.5-Flash, Gemini-1.5-Pro | Chat, Streaming |
| [X.AI](providers/xai.md) | Grok-4.1-Fast, Grok-4, Grok-3 | Chat, Streaming, 2M context |
| [Ollama](providers/ollama.md) | Llama 3, Mistral, CodeLlama | Chat, Streaming, Local inference |
| [AWS Bedrock](providers/bedrock.md) | Claude models, Titan models | Chat (external module) |

## Next Steps

- [Installation](getting-started/installation.md) - Get OmniLLM set up
- [Quick Start](getting-started/quickstart.md) - Your first LLM call
- [Providers](providers/index.md) - Configure specific providers
- [Features](features/streaming.md) - Explore advanced features
