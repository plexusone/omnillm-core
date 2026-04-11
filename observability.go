package omnillm

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/plexusone/omnillm-core/provider"
)

// LLMCallInfo provides metadata about the LLM call for observability
type LLMCallInfo struct {
	CallID       string    // Unique identifier for correlating BeforeRequest/AfterResponse
	ProviderName string    // e.g., "openai", "anthropic"
	StartTime    time.Time // When the call started
}

// newCallID generates a unique call ID for correlation
func newCallID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}

// ObservabilityHook allows external packages to observe LLM calls.
// Implementations can use this to add tracing, logging, or metrics
// without modifying the core OmniLLM library.
type ObservabilityHook interface {
	// BeforeRequest is called before each LLM call.
	// Returns a new context for trace/span propagation.
	// The hook should not modify the request.
	BeforeRequest(ctx context.Context, info LLMCallInfo, req *provider.ChatCompletionRequest) context.Context

	// AfterResponse is called after each LLM call completes.
	// This is called for both successful and failed requests.
	AfterResponse(ctx context.Context, info LLMCallInfo, req *provider.ChatCompletionRequest, resp *provider.ChatCompletionResponse, err error)

	// WrapStream wraps a stream for observability.
	// This allows the hook to observe streaming responses.
	// The returned stream must implement the same interface as the input.
	//
	// Note: For streaming, AfterResponse is only called if stream creation fails.
	// To track streaming completion timing and content, the wrapper returned here
	// should handle Close() or detect EOF in Recv() to finalize metrics/traces.
	WrapStream(ctx context.Context, info LLMCallInfo, req *provider.ChatCompletionRequest, stream provider.ChatCompletionStream) provider.ChatCompletionStream
}
