# Retry with Backoff

OmniLLM supports automatic retries for transient failures via a custom HTTP client.

## Configuration

```go
import (
    "net/http"
    "time"

    "github.com/plexusone/omnillm-core"
    "github.com/grokify/mogo/net/http/retryhttp"
)

// Create retry transport with exponential backoff
rt := retryhttp.NewWithOptions(
    retryhttp.WithMaxRetries(5),
    retryhttp.WithInitialBackoff(500 * time.Millisecond),
    retryhttp.WithMaxBackoff(30 * time.Second),
    retryhttp.WithOnRetry(func(attempt int, req *http.Request, resp *http.Response, err error, backoff time.Duration) {
        log.Printf("Retry attempt %d, waiting %v", attempt, backoff)
    }),
)

// Create client with retry-enabled HTTP client
client, err := omnillm.NewClient(omnillm.ClientConfig{
    Providers: []omnillm.ProviderConfig{
        {
            Provider: omnillm.ProviderNameOpenAI,
            APIKey:   os.Getenv("OPENAI_API_KEY"),
            HTTPClient: &http.Client{
                Transport: rt,
                Timeout:   2 * time.Minute,
            },
        },
    },
})
```

## Retry Transport Options

| Option | Default | Description |
|--------|---------|-------------|
| Max Retries | 3 | Maximum retry attempts |
| Initial Backoff | 1s | Starting backoff duration |
| Max Backoff | 30s | Cap on backoff duration |
| Backoff Multiplier | 2.0 | Exponential growth factor |
| Jitter | 10% | Randomness to prevent thundering herd |
| Retryable Status Codes | 429, 500, 502, 503, 504 | Rate limits + 5xx errors |

## Additional Options

```go
retryhttp.WithRetryableStatusCodes([]int{429, 500, 502, 503, 504})
retryhttp.WithShouldRetry(func(resp *http.Response, err error) bool {
    // Custom retry logic
})
retryhttp.WithLogger(logger)
```

## Retry-After Header

The retry transport automatically respects `Retry-After` headers from API responses.

## Provider Support

| Provider | Custom HTTP Client |
|----------|-------------------|
| OpenAI | Yes |
| Anthropic | Yes |
| X.AI | Yes |
| Ollama | Yes |
| Gemini | No (SDK-managed) |
| Bedrock | No (SDK-managed) |
