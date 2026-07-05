# Anthropic (Claude)

## Overview

- **Models**: Claude-Opus-4.1, Claude-Opus-4, Claude-Sonnet-4, Claude-3.7-Sonnet, Claude-3.5-Haiku, Claude-3-Opus, Claude-3-Sonnet, Claude-3-Haiku
- **Features**: Chat completions, streaming, system message support

## Configuration

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameAnthropic, APIKey: "your-anthropic-api-key"},
    },
})
```

## Available Models

| Model | Context Window | Description |
|-------|----------------|-------------|
| `omnillm.ModelClaudeOpus41` | 200K | Claude Opus 4.1 (most capable) |
| `omnillm.ModelClaudeOpus4` | 200K | Claude Opus 4 |
| `omnillm.ModelClaudeSonnet4` | 200K | Claude Sonnet 4 |
| `omnillm.ModelClaude37Sonnet` | 200K | Claude 3.7 Sonnet |
| `omnillm.ModelClaude35Haiku` | 200K | Claude 3.5 Haiku (fast) |
| `omnillm.ModelClaude3Opus` | 200K | Claude 3 Opus |
| `omnillm.ModelClaude3Sonnet` | 200K | Claude 3 Sonnet |
| `omnillm.ModelClaude3Haiku` | 200K | Claude 3 Haiku |

## TopK Sampling

Anthropic supports TopK sampling:

```go
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model: omnillm.ModelClaude3Sonnet,
    Messages: messages,
    TopK: &[]int{40}[0], // Consider only top 40 tokens
})
```

## System Messages

System messages are fully supported:

```go
messages := []omnillm.Message{
    {Role: omnillm.RoleSystem, Content: "You are a helpful assistant."},
    {Role: omnillm.RoleUser, Content: "Hello!"},
}
```

## Extended Thinking

Anthropic Claude models support extended thinking for complex reasoning tasks. Use the `Thinking` field for native control:

```go
budget := int64(8192)
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model: omnillm.ModelClaudeSonnet4,
    Thinking: &omnillm.ThinkingConfig{
        Type:         omnillm.ThinkingTypeEnabled,
        BudgetTokens: &budget,
    },
    Messages: []omnillm.Message{
        {Role: omnillm.RoleUser, Content: "Analyze this complex scenario..."},
    },
})
```

### ThinkingType Values

| Constant | Value | Description |
|----------|-------|-------------|
| `omnillm.ThinkingTypeEnabled` | `"enabled"` | Enable thinking with explicit budget |
| `omnillm.ThinkingTypeDisabled` | `"disabled"` | Disable thinking |
| `omnillm.ThinkingTypeAdaptive` | `"adaptive"` | Let the model decide |

### Using ReasoningEffort

For cross-provider compatibility, you can also use `ReasoningEffort`, which maps to Anthropic thinking:

| ReasoningEffort | Anthropic Behavior |
|-----------------|-------------------|
| `"none"` | Thinking disabled |
| `"low"` | Adaptive thinking |
| `"medium"` | Adaptive thinking |
| `"high"` | Enabled with budget derived from MaxTokens |

```go
effort := omnillm.ReasoningEffortHigh
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model:           omnillm.ModelClaudeSonnet4,
    ReasoningEffort: &effort,
    Messages:        messages,
})
```

See [Reasoning Feature Guide](../features/reasoning.md) for cross-provider usage.
