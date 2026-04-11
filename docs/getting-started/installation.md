# Installation

## Requirements

- Go 1.21 or later

## Install

```bash
go get github.com/plexusone/omnillm-core
```

## Thick Providers (Official SDKs)

OmniLLM-core includes thin providers that use native HTTP (no external dependencies). For official SDK support, install the thick provider packages:

```bash
# OpenAI (uses openai-go SDK)
go get github.com/plexusone/omnillm-openai

# Anthropic (uses anthropic-sdk-go)
go get github.com/plexusone/omnillm-anthropic

# Google Gemini (uses google.golang.org/genai)
go get github.com/plexusone/omnillm-gemini
```

Thick providers automatically override thin providers when imported:

```go
import (
    "github.com/plexusone/omnillm-core"
    _ "github.com/plexusone/omnillm-openai"    // Use official OpenAI SDK
    _ "github.com/plexusone/omnillm-anthropic" // Use official Anthropic SDK
)
```

## External Providers

Some providers with heavy SDK dependencies are available as separate modules:

```bash
# AWS Bedrock (requires AWS SDK)
go get github.com/plexusone/omnillm-bedrock
```

## Batteries-Included

For the simplest setup with all thick providers included:

```bash
go get github.com/plexusone/omnillm
```

This imports omnillm-core plus all thick providers automatically.

## Environment Variables

Set API keys for the providers you want to use:

```bash
export OPENAI_API_KEY="your-openai-api-key"
export ANTHROPIC_API_KEY="your-anthropic-api-key"
export GEMINI_API_KEY="your-gemini-api-key"
export XAI_API_KEY="your-xai-api-key"
export GLM_API_KEY="your-glm-api-key"
export KIMI_API_KEY="your-kimi-api-key"
export QWEN_API_KEY="your-qwen-api-key"
```

Ollama runs locally and doesn't require an API key.
