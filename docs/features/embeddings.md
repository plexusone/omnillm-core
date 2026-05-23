# Embeddings

OmniLLM supports text embeddings for converting text into vector representations. Embeddings are essential for semantic search, similarity matching, and retrieval-augmented generation (RAG) workflows.

## Overview

The embeddings API provides a unified interface across providers, allowing you to:

- Convert text into dense vector representations
- Process multiple texts in a single request (batch embedding)
- Control output dimensions for supported models
- Choose encoding format (float or base64)

## Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/plexusone/omnillm-core"
    "github.com/plexusone/omnillm-core/provider"
)

func main() {
    // Get an embedding provider
    embeddingProvider, err := omnillm.GetEmbeddingProvider(
        omnillm.ProviderNameOpenAI,
        omnillm.ProviderConfig{APIKey: "your-api-key"},
    )
    if err != nil {
        log.Fatal(err)
    }
    defer embeddingProvider.Close()

    // Create embeddings
    resp, err := embeddingProvider.CreateEmbedding(context.Background(), &provider.EmbeddingRequest{
        Model: "text-embedding-3-small",
        Input: []string{
            "Hello world",
            "How are you?",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Access the vectors
    for _, data := range resp.Data {
        fmt.Printf("Index %d: %d dimensions\n", data.Index, len(data.Embedding))
    }

    fmt.Printf("Total tokens: %d\n", resp.Usage.TotalTokens)
}
```

## EmbeddingProvider Interface

```go
type EmbeddingProvider interface {
    // CreateEmbedding creates embeddings for the given input texts
    CreateEmbedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)

    // Close closes the provider and cleans up resources
    Close() error

    // Name returns the provider name
    Name() string
}
```

## Request Types

### EmbeddingRequest

```go
type EmbeddingRequest struct {
    // Model is the embedding model to use (required)
    Model string `json:"model"`

    // Input is the text(s) to embed (required)
    Input []string `json:"input"`

    // EncodingFormat specifies the output format (optional)
    // "float" (default) or "base64"
    EncodingFormat EmbeddingEncodingFormat `json:"encoding_format,omitempty"`

    // Dimensions specifies output vector dimensions (optional)
    // Only supported by some models (e.g., text-embedding-3-small/large)
    Dimensions *int `json:"dimensions,omitempty"`

    // User is an optional end-user identifier
    User *string `json:"user,omitempty"`
}
```

### EmbeddingResponse

```go
type EmbeddingResponse struct {
    Object           string            `json:"object"`            // Always "list"
    Data             []EmbeddingData   `json:"data"`              // Embedding vectors
    Model            string            `json:"model"`             // Model used
    Usage            EmbeddingUsage    `json:"usage"`             // Token usage
    ProviderMetadata map[string]any    `json:"provider_metadata"` // Provider-specific
}

type EmbeddingData struct {
    Object    string    `json:"object"`    // Always "embedding"
    Index     int       `json:"index"`     // Position in input array
    Embedding []float64 `json:"embedding"` // Vector values
}

type EmbeddingUsage struct {
    PromptTokens int `json:"prompt_tokens"`
    TotalTokens  int `json:"total_tokens"`
}
```

## Provider Registry

Embedding providers are managed through a separate registry from chat providers:

```go
// List available embedding providers
providers := omnillm.ListEmbeddingProviders()
for _, name := range providers {
    fmt.Println(name)
}

// Get a provider factory
factory := omnillm.GetEmbeddingProviderFactory(omnillm.ProviderNameOpenAI)

// Create a provider instance
provider, err := omnillm.GetEmbeddingProvider(
    omnillm.ProviderNameOpenAI,
    omnillm.ProviderConfig{APIKey: apiKey},
)
```

### Registering Custom Embedding Providers

```go
omnillm.RegisterEmbeddingProvider(
    "custom-provider",
    func(config omnillm.ProviderConfig) (provider.EmbeddingProvider, error) {
        return NewCustomEmbeddingProvider(config.APIKey), nil
    },
    omnillm.PriorityThin, // 0 for thin, 10 for thick (SDK-based)
)
```

## Supported Providers

| Provider | Models | Dimensions |
|----------|--------|------------|
| OpenAI | text-embedding-3-small | 512, 1536 (default) |
| OpenAI | text-embedding-3-large | 256, 1024, 3072 (default) |
| OpenAI | text-embedding-ada-002 | 1536 (fixed) |

## Advanced Usage

### Custom Dimensions

Some models support variable output dimensions for storage/performance tradeoffs:

```go
dims := 512
resp, err := embeddingProvider.CreateEmbedding(ctx, &provider.EmbeddingRequest{
    Model:      "text-embedding-3-small",
    Input:      []string{"Hello world"},
    Dimensions: &dims, // Reduce from 1536 to 512 dimensions
})
```

### Batch Processing

Process multiple texts efficiently in a single request:

```go
texts := []string{
    "First document content...",
    "Second document content...",
    "Third document content...",
    // Up to 2048 inputs per request for OpenAI
}

resp, err := embeddingProvider.CreateEmbedding(ctx, &provider.EmbeddingRequest{
    Model: "text-embedding-3-small",
    Input: texts,
})

// Each embedding corresponds to the input at the same index
for i, data := range resp.Data {
    fmt.Printf("Text %d -> %d-dimensional vector\n", data.Index, len(data.Embedding))
}
```

## Use Cases

### Semantic Search

```go
// 1. Embed your documents
docEmbeddings := make([][]float64, len(documents))
for i, doc := range documents {
    resp, _ := embeddingProvider.CreateEmbedding(ctx, &provider.EmbeddingRequest{
        Model: "text-embedding-3-small",
        Input: []string{doc},
    })
    docEmbeddings[i] = resp.Data[0].Embedding
}

// 2. Embed the query
queryResp, _ := embeddingProvider.CreateEmbedding(ctx, &provider.EmbeddingRequest{
    Model: "text-embedding-3-small",
    Input: []string{"user search query"},
})
queryVector := queryResp.Data[0].Embedding

// 3. Find similar documents using cosine similarity
// (use a vector database for production workloads)
```

### RAG (Retrieval-Augmented Generation)

```go
// 1. Embed query and retrieve relevant context
queryResp, _ := embeddingProvider.CreateEmbedding(ctx, &provider.EmbeddingRequest{
    Model: "text-embedding-3-small",
    Input: []string{userQuery},
})

// 2. Search vector database for similar content
relevantDocs := vectorDB.Search(queryResp.Data[0].Embedding, topK)

// 3. Augment LLM prompt with retrieved context
chatResp, _ := chatClient.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model: omnillm.ModelGPT4o,
    Messages: []omnillm.Message{
        {Role: omnillm.RoleSystem, Content: "Use the following context: " + relevantDocs},
        {Role: omnillm.RoleUser, Content: userQuery},
    },
})
```

## Error Handling

```go
resp, err := embeddingProvider.CreateEmbedding(ctx, req)
if err != nil {
    if apiErr, ok := err.(*omnillm.APIError); ok {
        fmt.Printf("API Error: %s (status %d)\n", apiErr.Message, apiErr.StatusCode)
    }
    return err
}
```

## Best Practices

1. **Batch requests**: Process multiple texts together to reduce API calls
2. **Choose appropriate dimensions**: Use smaller dimensions when storage/speed matters more than precision
3. **Cache embeddings**: Store computed embeddings to avoid recomputation
4. **Use the same model**: Always use the same embedding model for documents and queries
5. **Normalize vectors**: Some similarity metrics work better with normalized vectors
