# Quick Start

## Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/plexusone/omnillm-core"
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

## Provider Switching

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

// Same API call for all providers
response1, _ := openaiClient.CreateChatCompletion(ctx, request)
response2, _ := anthropicClient.CreateChatCompletion(ctx, request)
```

## Running Examples

The repository includes comprehensive examples:

```bash
go run examples/basic/main.go
go run examples/streaming/main.go
go run examples/conversation/main.go
go run examples/memory_demo/main.go
```
