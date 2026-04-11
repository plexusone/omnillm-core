# Google Gemini

!!! warning "Moved to Separate Package"
    As of v0.15.0, Gemini has been moved to a separate package to reduce dependencies in omnillm-core.

    Install the Gemini provider:
    ```bash
    go get github.com/plexusone/omnillm-gemini
    ```

## Overview

- **Models**: Gemini-2.5-Pro, Gemini-2.5-Flash, Gemini-1.5-Pro, Gemini-1.5-Flash
- **Features**: Chat completions, streaming, massive context windows
- **SDK**: Uses official `google.golang.org/genai` SDK

## Installation

```bash
go get github.com/plexusone/omnillm-gemini
```

## Configuration

```go
import (
    "github.com/plexusone/omnillm-core"
    _ "github.com/plexusone/omnillm-gemini" // Register Gemini provider
)

client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {Provider: omnillm.ProviderNameGemini, APIKey: "your-gemini-api-key"},
    },
})
```

## Available Models

| Model | Context Window | Description |
|-------|----------------|-------------|
| `omnillm.ModelGemini25Pro` | 1M | Gemini 2.5 Pro |
| `omnillm.ModelGemini25Flash` | 1M | Gemini 2.5 Flash (fast) |
| `omnillm.ModelGemini15Pro` | 2M | Gemini 1.5 Pro (largest context) |
| `omnillm.ModelGemini15Flash` | 1M | Gemini 1.5 Flash |

## JSON Mode

Gemini supports JSON mode for structured outputs:

```go
response, err := client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
    Model: omnillm.ModelGemini25Pro,
    Messages: messages,
    ResponseFormat: &omnillm.ResponseFormat{Type: "json_object"},
})
```

## Large Context

Gemini 1.5 Pro supports up to 2 million tokens of context, making it ideal for:

- Long document analysis
- Large codebase understanding
- Extended conversation history
