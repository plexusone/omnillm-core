# Kimi (Moonshot AI)

## Overview

- **Models**: Kimi K2.5, Kimi K2 series, Moonshot V1 series
- **Features**: Chat completions, streaming, long context (up to 256K), multimodal support

## Configuration

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameKimi, APIKey: "your-kimi-api-key"},
    },
})
```

## Available Models

### Kimi K2.5 Series

| Model | Context Window | Description |
|-------|----------------|-------------|
| `kimi-k2.5` | 256K | Most intelligent multimodal model |

### Kimi K2 Series

| Model | Context Window | Description |
|-------|----------------|-------------|
| `kimi-k2-0905-preview` | 256K | Enhanced agentic coding model |
| `kimi-k2-0711-preview` | 128K | MoE base model with code/agent capabilities |
| `kimi-k2-turbo-preview` | 256K | High-speed version (60-100 tokens/sec) |
| `kimi-k2-thinking` | 256K | Long-term thinking with multi-step tool usage |
| `kimi-k2-thinking-turbo` | 256K | High-speed thinking model |

### Moonshot V1 Series

| Model | Context Window | Description |
|-------|----------------|-------------|
| `moonshot-v1-8k` | 8K | Short text generation |
| `moonshot-v1-32k` | 32K | Long text generation |
| `moonshot-v1-128k` | 128K | Very long text generation |
| `moonshot-v1-*-vision-preview` | 8K/32K/128K | Image understanding |

## OpenAI Compatibility

Kimi uses an OpenAI-compatible API:

```go
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model: "moonshot-v1-8k",
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
    Model: "kimi-k2.5",
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

## API Documentation

- [Kimi Models](https://platform.moonshot.cn/docs/api/chat)
- [Kimi API Reference](https://platform.moonshot.cn/docs)
