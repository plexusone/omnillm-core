# Reasoning and Extended Thinking

OmniLLM provides a unified interface for configuring reasoning and extended thinking capabilities across different LLM providers. This feature enables models to "think" before responding, improving accuracy for complex tasks.

## Overview

Different providers implement reasoning/thinking in different ways:

| Provider | Implementation | Field |
|----------|---------------|-------|
| OpenAI, X.AI, Qwen, GLM, Kimi | `reasoning_effort` parameter | `ReasoningEffort` |
| Anthropic | Extended thinking with budget | `Thinking` |
| Google Gemini | Thinking levels with budget | Maps from both |
| OpenRouter | Extended reasoning effort | `ReasoningEffort` |

OmniLLM abstracts these differences with two unified fields in `ChatCompletionRequest`:

- **`ReasoningEffort`** - OpenAI-compatible string: `"none"`, `"low"`, `"medium"`, `"high"`
- **`Thinking`** - Anthropic-style config with type and optional budget

## Basic Usage

### Using ReasoningEffort (OpenAI-style)

```go
import omnillm "github.com/plexusone/omnillm-core"

effort := omnillm.ReasoningEffortHigh
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model: "o1-preview", // or other reasoning-capable model
    Messages: []omnillm.Message{
        {Role: omnillm.RoleUser, Content: "Solve this complex math problem..."},
    },
    ReasoningEffort: &effort,
})
```

### Using Thinking (Anthropic-style)

```go
import omnillm "github.com/plexusone/omnillm-core"

budget := int64(8192)
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model: "claude-sonnet-4-20250514",
    Messages: []omnillm.Message{
        {Role: omnillm.RoleUser, Content: "Analyze this complex scenario..."},
    },
    Thinking: &omnillm.ThinkingConfig{
        Type:         omnillm.ThinkingTypeEnabled,
        BudgetTokens: &budget,
    },
})
```

## ReasoningEffort Constants

Use these constants for the `ReasoningEffort` field:

| Constant | Value | Description |
|----------|-------|-------------|
| `omnillm.ReasoningEffortNone` | `"none"` | Disable reasoning entirely |
| `omnillm.ReasoningEffortLow` | `"low"` | Light reasoning (default for reasoning models) |
| `omnillm.ReasoningEffortMedium` | `"medium"` | Moderate reasoning for complex tasks |
| `omnillm.ReasoningEffortHigh` | `"high"` | Deep reasoning for demanding problems |

## ThinkingConfig

The `Thinking` field accepts a `ThinkingConfig` struct:

```go
type ThinkingConfig struct {
    Type         string // Use ThinkingType* constants
    BudgetTokens *int64 // Optional token budget for thinking
}
```

### ThinkingType Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `omnillm.ThinkingTypeEnabled` | `"enabled"` | Enable thinking with explicit budget |
| `omnillm.ThinkingTypeDisabled` | `"disabled"` | Disable thinking entirely |
| `omnillm.ThinkingTypeAdaptive` | `"adaptive"` | Let the model decide thinking depth |

## Provider-Specific Behavior

### OpenAI / X.AI / Qwen / GLM / Kimi

These providers use the OpenAI-compatible `reasoning_effort` parameter directly. The `ReasoningEffort` field is passed through as-is.

```go
effort := omnillm.ReasoningEffortMedium
req := &omnillm.ChatCompletionRequest{
    Model:           "o1-mini",
    ReasoningEffort: &effort,
    // ...
}
```

### Anthropic

Anthropic uses extended thinking with a budget-based approach. OmniLLM maps:

| Input | Anthropic Behavior |
|-------|-------------------|
| `Thinking.Type = "enabled"` | Enabled with specified `BudgetTokens` |
| `Thinking.Type = "disabled"` | Thinking disabled |
| `Thinking.Type = "adaptive"` | Adaptive thinking |
| `ReasoningEffort = "none"` | Maps to disabled |
| `ReasoningEffort = "low"` or `"medium"` | Maps to adaptive |
| `ReasoningEffort = "high"` | Maps to enabled with budget derived from MaxTokens |

### Google Gemini

Gemini uses `ThinkingLevel` (`MINIMAL`, `LOW`, `MEDIUM`, `HIGH`) and optional `ThinkingBudget`. OmniLLM maps:

| ReasoningEffort | Gemini ThinkingLevel |
|-----------------|---------------------|
| `"none"` | `MINIMAL` |
| `"low"` | `LOW` |
| `"medium"` | `MEDIUM` |
| `"high"` | `HIGH` |

| ThinkingType | Gemini ThinkingLevel |
|--------------|---------------------|
| `"enabled"` | `HIGH` |
| `"disabled"` | `MINIMAL` |
| `"adaptive"` | `MEDIUM` |

If `Thinking.BudgetTokens` is set, it maps to Gemini's `ThinkingBudget`.

### OpenRouter

OpenRouter supports extended `reasoning_effort` values beyond the standard four:

- `max`, `xhigh`, `high`, `medium`, `low`, `minimal`, `none`

The standard `ReasoningEffort` values work directly. For extended values, set the string directly.

## Streaming with Reasoning

Reasoning works with streaming responses:

```go
effort := omnillm.ReasoningEffortHigh
stream, err := client.CreateChatCompletionStream(ctx, &omnillm.ChatCompletionRequest{
    Model:           "o1-preview",
    ReasoningEffort: &effort,
    Messages: []omnillm.Message{
        {Role: omnillm.RoleUser, Content: "Explain quantum computing..."},
    },
})
```

## Best Practices

1. **Use appropriate effort levels** - Higher reasoning effort increases latency and cost. Use `"high"` only for genuinely complex tasks.

2. **Set token budgets wisely** - For Anthropic's extended thinking, set `BudgetTokens` based on task complexity. Start with 4096-8192 for moderate tasks.

3. **Model selection matters** - Not all models support reasoning. Use reasoning-capable models like:
   - OpenAI: `o1-preview`, `o1-mini`, `o3-mini`
   - Anthropic: Claude 3.5+ models
   - X.AI: Grok models
   - Google: Gemini 2.0+ with thinking

4. **Prefer ReasoningEffort for portability** - If you might switch providers, use `ReasoningEffort` as it maps to all providers. Use `Thinking` only when you need Anthropic-specific budget control.

## Provider Support Matrix

| Provider | ReasoningEffort | Thinking | Notes |
|----------|----------------|----------|-------|
| OpenAI | Yes | Mapped | Direct support |
| Anthropic | Mapped | Yes | Native extended thinking |
| X.AI | Yes | Mapped | Direct support |
| Google Gemini | Mapped | Mapped | Uses ThinkingLevel/Budget |
| Qwen | Yes | Mapped | OpenAI-compatible |
| GLM | Yes | Mapped | OpenAI-compatible |
| Kimi | Yes | Mapped | OpenAI-compatible |
| OpenRouter | Yes | Mapped | Extended values available |
| Ollama | No | No | Not supported |
| AWS Bedrock | Model-dependent | Model-dependent | Via AdditionalModelRequestFields |
