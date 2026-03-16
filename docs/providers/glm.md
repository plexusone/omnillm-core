# GLM (Zhipu AI)

## Overview

- **Models**: GLM-5, GLM-4.7, GLM-4.6, GLM-4.5 series
- **Features**: Chat completions, streaming, OpenAI-compatible API, thinking modes

## Configuration

```go
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameGLM, APIKey: "your-glm-api-key"},
    },
})
```

## Available Models

| Model | Context Window | Description |
|-------|----------------|-------------|
| `glm-5` | 200K | Flagship MoE model (744B/40B active), forced thinking |
| `glm-5-code` | 200K | Code-specialized GLM-5 variant |
| `glm-4.7` | 200K | Premium model with Interleaved Thinking |
| `glm-4.7-flashx` | 200K | High-speed paid version with priority GPU |
| `glm-4.7-flash` | 200K | Free SOTA model with hybrid thinking |
| `glm-4.6` | 200K | Balanced model with auto-thinking |
| `glm-4.5` | 128K | Unified reasoning/coding/agent model |
| `glm-4.5-flash` | 128K | Free model with function calling |

## OpenAI Compatibility

GLM uses an OpenAI-compatible API endpoint, so standard parameters work seamlessly:

```go
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model: "glm-4.7-flash",
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
    Model: "glm-4.7-flash",
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

- [GLM Models](https://bigmodel.cn/dev/api/normal-model/glm-4)
- [GLM API Reference](https://open.bigmodel.cn/dev/howuse/introduction)
