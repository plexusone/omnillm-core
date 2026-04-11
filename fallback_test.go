package omnillm

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/plexusone/omnillm-core/provider"
)

// mockProvider is a test provider with configurable behavior
type mockProvider struct {
	name           string
	completionResp *provider.ChatCompletionResponse
	completionErr  error
	streamResp     provider.ChatCompletionStream
	streamErr      error
	callCount      int
	failUntil      int     // Fail first N calls
	errorSequence  []error // Specific errors for each call
}

func newMockProvider(name string) *mockProvider {
	finishReason := "stop"
	return &mockProvider{
		name: name,
		completionResp: &provider.ChatCompletionResponse{
			ID:    "mock-response-" + name,
			Model: "mock-model",
			Choices: []provider.ChatCompletionChoice{
				{
					Index: 0,
					Message: provider.Message{
						Role:    provider.RoleAssistant,
						Content: "Hello from " + name,
					},
					FinishReason: &finishReason,
				},
			},
			Usage: provider.Usage{
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
			},
		},
	}
}

func (m *mockProvider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
	m.callCount++

	// Check error sequence first
	if len(m.errorSequence) > 0 && m.callCount <= len(m.errorSequence) {
		if err := m.errorSequence[m.callCount-1]; err != nil {
			return nil, err
		}
	}

	// Check failUntil
	if m.callCount <= m.failUntil {
		if m.completionErr != nil {
			return nil, m.completionErr
		}
		return nil, errors.New("mock failure")
	}

	if m.completionErr != nil {
		return nil, m.completionErr
	}

	return m.completionResp, nil
}

func (m *mockProvider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
	m.callCount++

	if m.callCount <= m.failUntil {
		if m.streamErr != nil {
			return nil, m.streamErr
		}
		return nil, errors.New("mock stream failure")
	}

	if m.streamErr != nil {
		return nil, m.streamErr
	}

	if m.streamResp != nil {
		return m.streamResp, nil
	}

	return &mockStream{chunks: []string{"Hello ", "from ", m.name}}, nil
}

func (m *mockProvider) Close() error {
	return nil
}

func (m *mockProvider) Name() string {
	return m.name
}

// mockStream is a test stream
type mockStream struct {
	chunks []string
	index  int
	closed bool
}

func (s *mockStream) Recv() (*provider.ChatCompletionChunk, error) {
	if s.closed {
		return nil, errors.New("stream closed")
	}
	if s.index >= len(s.chunks) {
		return nil, io.EOF
	}

	chunk := &provider.ChatCompletionChunk{
		ID:    "chunk",
		Model: "mock-model",
		Choices: []provider.ChatCompletionChoice{
			{
				Index: 0,
				Delta: &provider.Message{
					Role:    provider.RoleAssistant,
					Content: s.chunks[s.index],
				},
			},
		},
	}
	s.index++
	return chunk, nil
}

func (s *mockStream) Close() error {
	s.closed = true
	return nil
}

func TestFallbackProvider_PrimarySuccess(t *testing.T) {
	primary := newMockProvider("primary")
	fallback := newMockProvider("fallback")

	fp := NewFallbackProvider(primary, []provider.Provider{fallback}, nil)

	req := &provider.ChatCompletionRequest{
		Model:    "test-model",
		Messages: []provider.Message{{Role: "user", Content: "Hello"}},
	}

	resp, err := fp.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "mock-response-primary" {
		t.Errorf("expected response from primary, got %s", resp.ID)
	}

	if primary.callCount != 1 {
		t.Errorf("expected primary to be called once, got %d", primary.callCount)
	}

	if fallback.callCount != 0 {
		t.Errorf("expected fallback not to be called, got %d", fallback.callCount)
	}
}

func TestFallbackProvider_FallbackOnError(t *testing.T) {
	primary := newMockProvider("primary")
	primary.completionErr = NewAPIErrorFull("primary", 500, "server error", "server_error", "500")

	fallback := newMockProvider("fallback")

	fp := NewFallbackProvider(primary, []provider.Provider{fallback}, nil)

	req := &provider.ChatCompletionRequest{
		Model:    "test-model",
		Messages: []provider.Message{{Role: "user", Content: "Hello"}},
	}

	resp, err := fp.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "mock-response-fallback" {
		t.Errorf("expected response from fallback, got %s", resp.ID)
	}

	if primary.callCount != 1 {
		t.Errorf("expected primary to be called once, got %d", primary.callCount)
	}

	if fallback.callCount != 1 {
		t.Errorf("expected fallback to be called once, got %d", fallback.callCount)
	}
}

func TestFallbackProvider_NoFallbackOnAuthError(t *testing.T) {
	primary := newMockProvider("primary")
	primary.completionErr = NewAPIErrorFull("primary", 401, "unauthorized", "auth_error", "401")

	fallback := newMockProvider("fallback")

	fp := NewFallbackProvider(primary, []provider.Provider{fallback}, nil)

	req := &provider.ChatCompletionRequest{
		Model:    "test-model",
		Messages: []provider.Message{{Role: "user", Content: "Hello"}},
	}

	_, err := fp.CreateChatCompletion(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if fallback.callCount != 0 {
		t.Errorf("expected fallback not to be called for auth error, got %d", fallback.callCount)
	}
}

func TestFallbackProvider_AllProvidersFail(t *testing.T) {
	primary := newMockProvider("primary")
	primary.completionErr = NewAPIErrorFull("primary", 500, "server error", "server_error", "500")

	fallback := newMockProvider("fallback")
	fallback.completionErr = NewAPIErrorFull("fallback", 503, "service unavailable", "unavailable", "503")

	fp := NewFallbackProvider(primary, []provider.Provider{fallback}, nil)

	req := &provider.ChatCompletionRequest{
		Model:    "test-model",
		Messages: []provider.Message{{Role: "user", Content: "Hello"}},
	}

	_, err := fp.CreateChatCompletion(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var fbErr *FallbackError
	if !errors.As(err, &fbErr) {
		t.Fatalf("expected FallbackError, got %T", err)
	}

	if len(fbErr.Attempts) != 2 {
		t.Errorf("expected 2 attempts, got %d", len(fbErr.Attempts))
	}
}

func TestFallbackProvider_WithCircuitBreaker(t *testing.T) {
	primary := newMockProvider("primary")
	primary.completionErr = NewAPIErrorFull("primary", 500, "server error", "server_error", "500")

	fallback := newMockProvider("fallback")

	cbConfig := &CircuitBreakerConfig{
		FailureThreshold: 2,
		SuccessThreshold: 1,
		Timeout:          50 * time.Millisecond,
		MinimumRequests:  10,
	}

	fp := NewFallbackProvider(primary, []provider.Provider{fallback}, &FallbackProviderConfig{
		CircuitBreakerConfig: cbConfig,
	})

	req := &provider.ChatCompletionRequest{
		Model:    "test-model",
		Messages: []provider.Message{{Role: "user", Content: "Hello"}},
	}

	// First two calls should hit primary (and fail), then fallback (and succeed)
	for i := 0; i < 2; i++ {
		resp, err := fp.CreateChatCompletion(context.Background(), req)
		if err != nil {
			t.Fatalf("call %d: unexpected error: %v", i+1, err)
		}
		if resp.ID != "mock-response-fallback" {
			t.Errorf("call %d: expected fallback response", i+1)
		}
	}

	// After 2 failures, primary circuit should be open
	cb := fp.CircuitBreaker("primary")
	if cb == nil {
		t.Fatal("expected circuit breaker for primary")
	}

	if cb.State() != CircuitOpen {
		t.Errorf("expected primary circuit to be open, got %v", cb.State())
	}

	// Third call should skip primary entirely
	primary.callCount = 0
	fallback.callCount = 0

	resp, err := fp.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "mock-response-fallback" {
		t.Error("expected fallback response")
	}

	if primary.callCount != 0 {
		t.Errorf("expected primary to be skipped due to open circuit, got %d calls", primary.callCount)
	}
}

func TestFallbackProvider_StreamingSuccess(t *testing.T) {
	primary := newMockProvider("primary")
	fallback := newMockProvider("fallback")

	fp := NewFallbackProvider(primary, []provider.Provider{fallback}, nil)

	req := &provider.ChatCompletionRequest{
		Model:    "test-model",
		Messages: []provider.Message{{Role: "user", Content: "Hello"}},
	}

	stream, err := fp.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer stream.Close()

	var content string
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("unexpected stream error: %v", err)
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			content += chunk.Choices[0].Delta.Content
		}
	}

	if content != "Hello from primary" {
		t.Errorf("expected 'Hello from primary', got %q", content)
	}
}

func TestFallbackProvider_StreamingFallback(t *testing.T) {
	primary := newMockProvider("primary")
	primary.streamErr = NewAPIErrorFull("primary", 500, "server error", "server_error", "500")

	fallback := newMockProvider("fallback")

	fp := NewFallbackProvider(primary, []provider.Provider{fallback}, nil)

	req := &provider.ChatCompletionRequest{
		Model:    "test-model",
		Messages: []provider.Message{{Role: "user", Content: "Hello"}},
	}

	stream, err := fp.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer stream.Close()

	var content string
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("unexpected stream error: %v", err)
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			content += chunk.Choices[0].Delta.Content
		}
	}

	if content != "Hello from fallback" {
		t.Errorf("expected 'Hello from fallback', got %q", content)
	}
}

func TestFallbackProvider_Name(t *testing.T) {
	primary := newMockProvider("openai")
	fallback := newMockProvider("anthropic")

	fp := NewFallbackProvider(primary, []provider.Provider{fallback}, nil)

	expected := "openai+fallback"
	if fp.Name() != expected {
		t.Errorf("expected name %q, got %q", expected, fp.Name())
	}
}

func TestFallbackProvider_Close(t *testing.T) {
	primary := newMockProvider("primary")
	fallback := newMockProvider("fallback")

	fp := NewFallbackProvider(primary, []provider.Provider{fallback}, nil)

	err := fp.Close()
	if err != nil {
		t.Errorf("unexpected error on close: %v", err)
	}
}

func TestFallbackError(t *testing.T) {
	err := &FallbackError{
		Attempts: []FallbackAttempt{
			{Provider: "primary", Error: errors.New("primary failed")},
			{Provider: "fallback", Error: errors.New("fallback failed")},
		},
		LastError: errors.New("fallback failed"),
	}

	expected := "all 2 providers failed, last error: fallback failed"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}

	unwrapped := err.Unwrap()
	if unwrapped.Error() != "fallback failed" {
		t.Errorf("expected unwrapped error 'fallback failed', got %q", unwrapped.Error())
	}
}

func TestFallbackProvider_ProviderMetadata(t *testing.T) {
	primary := newMockProvider("primary")
	fallback := newMockProvider("fallback")

	fp := NewFallbackProvider(primary, []provider.Provider{fallback}, nil)

	req := &provider.ChatCompletionRequest{
		Model:    "test-model",
		Messages: []provider.Message{{Role: "user", Content: "Hello"}},
	}

	resp, err := fp.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ProviderMetadata == nil {
		t.Fatal("expected provider metadata")
	}

	providerUsed, ok := resp.ProviderMetadata["fallback_provider_used"]
	if !ok || providerUsed != "primary" {
		t.Errorf("expected fallback_provider_used=primary, got %v", providerUsed)
	}

	attemptCount, ok := resp.ProviderMetadata["fallback_attempt_count"]
	if !ok || attemptCount != 1 {
		t.Errorf("expected fallback_attempt_count=1, got %v", attemptCount)
	}
}
