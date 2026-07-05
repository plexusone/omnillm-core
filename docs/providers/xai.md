# X.AI (Grok)

## Overview

- **Models**: Grok-4.1-Fast, Grok-4, Grok-4-Fast, Grok-Code-Fast, Grok-3, Grok-3-Mini, Grok-2, Grok-2-Vision
- **Features**: Chat completions, streaming, OpenAI-compatible API, 2M context window

## Configuration

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameXAI, APIKey: "your-xai-api-key"},
    },
})
```

## Available Models

| Model | Context Window | Description |
|-------|----------------|-------------|
| `omnillm.ModelGrok41Fast` | 2M | Grok 4.1 Fast (reasoning/non-reasoning) |
| `omnillm.ModelGrok4` | 128K | Grok 4 (0709) |
| `omnillm.ModelGrok4Fast` | 2M | Grok 4 Fast (reasoning/non-reasoning) |
| `omnillm.ModelGrokCodeFast` | 128K | Grok Code Fast |
| `omnillm.ModelGrok3` | 128K | Grok 3 |
| `omnillm.ModelGrok3Mini` | 128K | Grok 3 Mini |
| `omnillm.ModelGrok2` | 128K | Grok 2 |
| `omnillm.ModelGrok2Vision` | 128K | Grok 2 Vision |

## OpenAI Compatibility

X.AI uses an OpenAI-compatible API, so parameters like `Seed`, `PresencePenalty`, and `FrequencyPenalty` are supported:

```go
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model: omnillm.ModelGrok4Fast,
    Messages: messages,
    Seed: &[]int{42}[0],
    PresencePenalty: &[]float64{0.5}[0],
})
```

## Reasoning

X.AI Grok models support the `reasoning_effort` parameter (OpenAI-compatible):

```go
effort := omnillm.ReasoningEffortHigh
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model:           omnillm.ModelGrok4Fast,
    ReasoningEffort: &effort,
    Messages: []omnillm.Message{
        {Role: omnillm.RoleUser, Content: "Solve this complex problem..."},
    },
})
```

### ReasoningEffort Values

| Constant | Value | Description |
|----------|-------|-------------|
| `omnillm.ReasoningEffortNone` | `"none"` | Disable reasoning |
| `omnillm.ReasoningEffortLow` | `"low"` | Light reasoning |
| `omnillm.ReasoningEffortMedium` | `"medium"` | Moderate reasoning |
| `omnillm.ReasoningEffortHigh` | `"high"` | Deep reasoning |

### Reasoning-Capable Models

- `grok-4.1-fast` - Supports reasoning mode
- `grok-4-fast` - Supports reasoning mode

See [Reasoning Feature Guide](../features/reasoning.md) for cross-provider usage.
