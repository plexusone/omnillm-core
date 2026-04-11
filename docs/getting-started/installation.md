# Installation

## Requirements

- Go 1.21 or later

## Install

```bash
go get github.com/plexusone/omnillm-core
```

## External Providers

Some providers with heavy SDK dependencies are available as separate modules:

```bash
# AWS Bedrock (requires AWS SDK)
go get github.com/plexusone/omnillm-bedrock
```

## Environment Variables

Set API keys for the providers you want to use:

```bash
export OPENAI_API_KEY="your-openai-api-key"
export ANTHROPIC_API_KEY="your-anthropic-api-key"
export GEMINI_API_KEY="your-gemini-api-key"
export XAI_API_KEY="your-xai-api-key"
```

Ollama runs locally and doesn't require an API key.
