# Qwen (Alibaba Cloud)

## Overview

- **Models**: Qwen3 Max, QwQ, Qwen3.5, Qwen3, Qwen2.5 series
- **Features**: Chat completions, streaming, thinking modes, wide global availability

## Configuration

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameQwen, APIKey: "your-qwen-api-key"},
    },
})
```

## Available Models

### Flagship Models

| Model | Description |
|-------|-------------|
| `qwen3-max` | Latest flagship with thinking capability |
| `qwen3-max-preview` | Preview with extended thinking |
| `qwen-max` | Flagship reasoning (International/China) |

### Balanced Models

| Model | Description |
|-------|-------------|
| `qwen3.5-plus` | Advanced balanced with thinking mode |
| `qwen-plus` | Balanced performance (wide availability) |
| `qwen3.5-flash` | Ultra-fast lightweight |
| `qwen-flash` | Fast with context caching |

### Deep Reasoning (QwQ)

| Model | Description |
|-------|-------------|
| `qwq-plus` | Deep reasoning with extended chain-of-thought |
| `qwq-32b` | Open-source 32B reasoning model |

### Open Source Models

| Model | Parameters | Description |
|-------|------------|-------------|
| `qwen3.5-397b-a17b` | 397B | Largest open-source model |
| `qwen3-235b-a22b` | 235B | Dual-mode thinking/non-thinking |
| `qwen3-32b` | 32B | Versatile dual-mode model |
| `qwen2.5-72b-instruct` | 72B | Large instruction-following |
| `qwen2.5-14b-instruct-1m` | 14B | Extended 1M token context |

## OpenAI Compatibility

Qwen uses an OpenAI-compatible API via DashScope:

```go
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model: "qwen-plus",
    Messages: []omnillm.Message{
        {Role: omnillm.RoleUser, Content: "Hello!"},
    },
    Temperature: &[]float64{0.7}[0],
    MaxTokens:   &[]int{1000}[0],
})
```

## Streaming

```go
stream, err := client.CreateChatCompletionStream(ctx, &omnillm.ChatCompletionRequest{
    Model: "qwen-plus",
    Messages: messages,
})
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for {
    chunk, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(chunk.Choices[0].Delta.Content)
}
```

## Regional Availability

Qwen models have wide availability including US regions:

| Model | Regions |
|-------|---------|
| `qwen-plus` | International, US, China |
| `qwen-flash` | International, US, China |
| `qwen-max` | International, China |

## API Documentation

- [Qwen Models](https://help.aliyun.com/zh/model-studio/getting-started/models)
- [Qwen API Reference](https://help.aliyun.com/zh/model-studio/developer-reference/use-qwen-by-calling-api)
