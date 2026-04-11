# AWS Bedrock

## Overview

AWS Bedrock is available as an external module to avoid pulling AWS SDK dependencies for users who don't need it.

- **Models**: Claude models, Titan models
- **Features**: Chat completions, multiple model families

## Installation

```bash
go get github.com/plexusone/omnillm-bedrock
```

## Configuration

```go
import (
    "github.com/plexusone/omnillm-core"
    "github.com/plexusone/omnillm-bedrock"
)

// Create the Bedrock provider
bedrockProvider, err := bedrock.NewProvider("us-east-1")
if err != nil {
    log.Fatal(err)
}

// Use it with omnillm via CustomProvider
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {CustomProvider: bedrockProvider},
    },
})
```

## AWS Credentials

The Bedrock provider uses the standard AWS credential chain:

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. Shared credentials file (`~/.aws/credentials`)
3. IAM role (when running on AWS)

## Why External?

AWS SDK v2 adds 17+ transitive dependencies. By keeping Bedrock as an external module, users who don't need AWS can keep their dependency tree lean.

## Source Code

See [github.com/plexusone/omnillm-bedrock](https://github.com/plexusone/omnillm-bedrock) for the full implementation.
