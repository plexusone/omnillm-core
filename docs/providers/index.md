# Providers

OmniLLM supports multiple LLM providers through a unified interface. Each provider is configured via `ProviderConfig` and implements the same `Provider` interface.

## Built-in Providers

| Provider | Package | Description |
|----------|---------|-------------|
| [OpenAI](openai.md) | Built-in | GPT-5, GPT-4o, GPT-4-turbo, GPT-3.5-turbo |
| [Anthropic](anthropic.md) | Built-in | Claude Opus 4, Sonnet 4, Claude 3.x series |
| [Google Gemini](gemini.md) | Built-in | Gemini 2.5/1.5 Pro and Flash |
| [X.AI](xai.md) | Built-in | Grok 4, Grok 3, 2M context window |
| [GLM](glm.md) | Built-in | Zhipu AI GLM-5, GLM-4.7, GLM-4.5 series |
| [Kimi](kimi.md) | Built-in | Moonshot AI Kimi K2.5, K2, Moonshot V1 |
| [Qwen](qwen.md) | Built-in | Alibaba Cloud Qwen3, QwQ, Qwen2.5 |
| [Ollama](ollama.md) | Built-in | Local models (Llama, Mistral, etc.) |

## Batteries Included

For a batteries-included experience with additional providers and official SDK implementations, use [omnillm](https://github.com/plexusone/omnillm):

```bash
go get github.com/plexusone/omnillm
```

This adds:

| Provider | Module | Description |
|----------|--------|-------------|
| Google Gemini | `omnillm-gemini` | Thick provider using official `google.golang.org/genai` SDK |
| AWS Bedrock | `omnillm-bedrock` | Thick provider using official AWS SDK v2 |
| OpenAI | `omnillm-openai` | Thick provider using official `openai-go` SDK |
| Anthropic | `omnillm-anthropic` | Thick provider using official `anthropic-sdk-go` SDK |

Thick providers use official SDKs for full API coverage, automatic retries, and SDK-managed authentication. They automatically override thin providers when imported.

## External Providers

Thick providers can also be imported individually:

| Provider | Module | Why External |
|----------|--------|--------------|
| [AWS Bedrock](bedrock.md) | `github.com/plexusone/omnillm-bedrock` | AWS SDK v2 adds 17+ transitive dependencies |
| Google Gemini | `github.com/plexusone/omnillm-gemini` | Official genai SDK |
| OpenAI | `github.com/plexusone/omnillm-openai` | Official openai-go SDK |
| Anthropic | `github.com/plexusone/omnillm-anthropic` | Official anthropic-sdk-go SDK |

## Multi-Provider Configuration

Configure multiple providers for fallback support:

```go
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

## Model Support Summary

| Provider | Models | Context Window | Features |
|----------|--------|----------------|----------|
| OpenAI | GPT-5, GPT-4.1, GPT-4o | 128K-200K | Chat, Streaming, Tools |
| Anthropic | Claude Opus 4, Sonnet 4 | 200K | Chat, Streaming, System |
| Gemini | Gemini 2.5 Pro/Flash | 1M-2M | Chat, Streaming |
| X.AI | Grok 4, Grok 3 | 128K-2M | Chat, Streaming, Tools |
| GLM | GLM-5, GLM-4.7, GLM-4.5 | 128K-200K | Chat, Streaming, Thinking |
| Kimi | Kimi K2.5, Moonshot V1 | 8K-256K | Chat, Streaming, Vision |
| Qwen | Qwen3, QwQ, Qwen2.5 | 128K-1M | Chat, Streaming, Thinking |
| Ollama | Llama 3, Mistral | Varies | Chat, Streaming, Local |
| Bedrock | Claude, Titan | Varies | Chat |
