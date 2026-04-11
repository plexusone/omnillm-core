# Conversation Memory

OmniLLM supports persistent conversation memory using any Key-Value Store that implements the [Sogo KVS interface](https://github.com/grokify/sogo/blob/master/database/kvs/definitions.go).

## Configuration

```go
memoryConfig := omnillm.MemoryConfig{
    MaxMessages: 50,                    // Keep last 50 messages per session
    TTL:         24 * time.Hour,        // Messages expire after 24 hours
    KeyPrefix:   "myapp:conversations", // Custom key prefix
}

client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-api-key"},
    },
    Memory:       kvsClient,   // Your KVS implementation
    MemoryConfig: &memoryConfig,
})
```

## Memory-Aware Completions

```go
// Create a session with system message
err = client.CreateConversationWithSystemMessage(ctx, "user-123",
    "You are a helpful assistant that remembers our conversation history.")

// Use memory-aware completion - automatically loads conversation history
response, err := client.CreateChatCompletionWithMemory(ctx, "user-123", &omnillm.ChatCompletionRequest{
    Model: omnillm.ModelGPT4o,
    Messages: []omnillm.Message{
        {Role: omnillm.RoleUser, Content: "What did we discuss last time?"},
    },
    MaxTokens: &[]int{200}[0],
})
```

## Memory Management

```go
// Load conversation history
conversation, err := client.LoadConversation(ctx, "user-123")

// Get just the messages
messages, err := client.GetConversationMessages(ctx, "user-123")

// Manually append messages
err = client.AppendMessage(ctx, "user-123", omnillm.Message{
    Role:    omnillm.RoleUser,
    Content: "Remember this important fact: I prefer JSON responses.",
})

// Delete conversation
err = client.DeleteConversation(ctx, "user-123")
```

## KVS Backend Support

Memory works with any KVS implementation:

- **Redis**: High-performance, distributed memory
- **DynamoDB**: AWS-native storage
- **In-Memory**: Testing and development
- **Custom**: Any implementation of the Sogo KVS interface

```go
// Example with Redis
redisKVS := redis.NewKVSClient("localhost:6379")
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "your-key"},
    },
    Memory: redisKVS,
})
```

## Mock KVS for Testing

```go
import omnillmtest "github.com/plexusone/omnillm-core/testing"

mockKVS := omnillmtest.NewMockKVS()

client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameOpenAI, APIKey: "test-key"},
    },
    Memory: mockKVS,
})
```
