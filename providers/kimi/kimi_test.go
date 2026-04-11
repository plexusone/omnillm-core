// Package kimi provides unit tests for the Kimi API client using a mock HTTP server.
package kimi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/plexusone/omnillm-core/provider"
)

// testModel is the model used in unit tests
const testModel = "moonshot-v1-8k"

// newTestServer creates a mock HTTP server that returns the given response JSON and status code.
func newTestServer(t *testing.T, statusCode int, responseBody string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))
}

// chatCompletionResponse builds a valid Kimi completion JSON response.
func chatCompletionResponse(content string) string {
	return `{
		"id": "test-id-123",
		"object": "chat.completion",
		"created": 1700000000,
		"model": "` + testModel + `",
		"choices": [{
			"index": 0,
			"message": {"role": "assistant", "content": "` + content + `"},
			"finish_reason": "stop"
		}],
		"usage": {"prompt_tokens": 10, "completion_tokens": 20, "total_tokens": 30}
	}`
}

// streamingResponse builds a mock SSE streaming body.
func streamingResponse() string {
	chunk1 := `{"id":"chunk-1","object":"chat.completion.chunk","created":1700000000,"model":"` + testModel + `","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}`
	chunk2 := `{"id":"chunk-2","object":"chat.completion.chunk","created":1700000000,"model":"` + testModel + `","choices":[{"index":0,"delta":{"content":" world"},"finish_reason":null}]}`
	chunk3 := `{"id":"chunk-3","object":"chat.completion.chunk","created":1700000000,"model":"` + testModel + `","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`
	return "data: " + chunk1 + "\n\ndata: " + chunk2 + "\n\ndata: " + chunk3 + "\n\ndata: [DONE]\n\n"
}

// ─── Client Unit Tests ───────────────────────────────────────────────────────

func TestKimiClient_New_DefaultBaseURL(t *testing.T) {
	c := New("test-key", "", nil)
	if c.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, DefaultBaseURL)
	}
	if c.apiKey != "test-key" {
		t.Errorf("apiKey = %q, want test-key", c.apiKey)
	}
	if c.client == nil {
		t.Error("http.Client should not be nil")
	}
}

func TestKimiClient_New_CustomBaseURL(t *testing.T) {
	c := New("key", "https://custom.example.com", nil)
	if c.baseURL != "https://custom.example.com" {
		t.Errorf("baseURL = %q, want custom URL", c.baseURL)
	}
}

func TestKimiClient_Name(t *testing.T) {
	c := New("key", "", nil)
	if c.Name() != "kimi" {
		t.Errorf("Name() = %q, want kimi", c.Name())
	}
}

func TestKimiClient_CreateCompletion_ValidationErrors(t *testing.T) {
	c := New("key", "", nil)
	ctx := context.Background()

	t.Run("empty model", func(t *testing.T) {
		_, err := c.CreateCompletion(ctx, &Request{
			Messages: []Message{{Role: "user", Content: "hello"}},
		})
		if err == nil || !strings.Contains(err.Error(), "model cannot be empty") {
			t.Errorf("expected 'model cannot be empty' error, got %v", err)
		}
	})

	t.Run("empty messages", func(t *testing.T) {
		_, err := c.CreateCompletion(ctx, &Request{Model: "moonshot-v1-8k"})
		if err == nil || !strings.Contains(err.Error(), "messages cannot be empty") {
			t.Errorf("expected 'messages cannot be empty' error, got %v", err)
		}
	})
}

func TestKimiClient_CreateCompletion_Success(t *testing.T) {
	srv := newTestServer(t, http.StatusOK, chatCompletionResponse("test successful"))
	defer srv.Close()

	c := New("test-api-key", srv.URL, nil)
	resp, err := c.CreateCompletion(context.Background(), &Request{
		Model:    "moonshot-v1-8k",
		Messages: []Message{{Role: "user", Content: "say test successful"}},
	})
	if err != nil {
		t.Fatalf("CreateCompletion failed: %v", err)
	}
	if resp.ID != "test-id-123" {
		t.Errorf("ID = %q, want test-id-123", resp.ID)
	}
	if len(resp.Choices) == 0 {
		t.Fatal("no choices in response")
	}
	if resp.Choices[0].Message.Content != "test successful" {
		t.Errorf("content = %q, want 'test successful'", resp.Choices[0].Message.Content)
	}
	if resp.Usage.TotalTokens != 30 {
		t.Errorf("TotalTokens = %d, want 30", resp.Usage.TotalTokens)
	}
}

func TestKimiClient_CreateCompletion_APIError(t *testing.T) {
	errBody := `{"error":{"message":"invalid api key","type":"auth_error","code":"401"}}`
	srv := newTestServer(t, http.StatusUnauthorized, errBody)
	defer srv.Close()

	c := New("bad-key", srv.URL, nil)
	_, err := c.CreateCompletion(context.Background(), &Request{
		Model:    "moonshot-v1-8k",
		Messages: []Message{{Role: "user", Content: "hello"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid api key") {
		t.Errorf("error = %q, want it to contain 'invalid api key'", err.Error())
	}
}

func TestKimiClient_CreateCompletion_AuthorizationHeader(t *testing.T) {
	var capturedHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(chatCompletionResponse("ok")))
	}))
	defer srv.Close()

	c := New("my-secret-key", srv.URL, nil)
	_, err := c.CreateCompletion(context.Background(), &Request{
		Model:    "moonshot-v1-8k",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedHeader != "Bearer my-secret-key" {
		t.Errorf("Authorization header = %q, want 'Bearer my-secret-key'", capturedHeader)
	}
}

func TestKimiClient_CreateCompletion_StreamSetToFalse(t *testing.T) {
	var capturedStream *bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body Request
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &body)
		capturedStream = body.Stream
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(chatCompletionResponse("ok")))
	}))
	defer srv.Close()

	c := New("key", srv.URL, nil)
	_, _ = c.CreateCompletion(context.Background(), &Request{
		Model:    "moonshot-v1-8k",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if capturedStream == nil || *capturedStream {
		t.Error("stream should be set to false for non-streaming requests")
	}
}

// ─── Streaming Unit Tests ────────────────────────────────────────────────────

func TestKimiClient_CreateCompletionStream_ValidationErrors(t *testing.T) {
	c := New("key", "", nil)
	ctx := context.Background()

	t.Run("empty model", func(t *testing.T) {
		_, err := c.CreateCompletionStream(ctx, &Request{
			Messages: []Message{{Role: "user", Content: "hello"}},
		})
		if err == nil || !strings.Contains(err.Error(), "model cannot be empty") {
			t.Errorf("expected 'model cannot be empty' error, got %v", err)
		}
	})

	t.Run("empty messages", func(t *testing.T) {
		_, err := c.CreateCompletionStream(ctx, &Request{Model: "moonshot-v1-8k"})
		if err == nil || !strings.Contains(err.Error(), "messages cannot be empty") {
			t.Errorf("expected 'messages cannot be empty' error, got %v", err)
		}
	})
}

func newStreamServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(streamingResponse()))
	}))
}

func TestKimiClient_CreateCompletionStream_Success(t *testing.T) {
	srv := newStreamServer(t)
	defer srv.Close()

	c := New("key", srv.URL, nil)
	stream, err := c.CreateCompletionStream(context.Background(), &Request{
		Model:    "moonshot-v1-8k",
		Messages: []Message{{Role: "user", Content: "hello"}},
	})
	if err != nil {
		t.Fatalf("CreateCompletionStream failed: %v", err)
	}
	defer stream.Close()

	var content strings.Builder
	chunkCount := 0
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv error: %v", err)
		}
		chunkCount++
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			content.WriteString(chunk.Choices[0].Delta.Content)
		}
	}

	if chunkCount == 0 {
		t.Error("expected at least one chunk")
	}
	if content.String() != "Hello world" {
		t.Errorf("content = %q, want 'Hello world'", content.String())
	}
}

func TestKimiStream_Close(t *testing.T) {
	srv := newStreamServer(t)
	defer srv.Close()

	c := New("key", srv.URL, nil)
	stream, err := c.CreateCompletionStream(context.Background(), &Request{
		Model:    "moonshot-v1-8k",
		Messages: []Message{{Role: "user", Content: "hello"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := stream.Close(); err != nil {
		t.Errorf("Close() error: %v", err)
	}
	// Double-close should be a no-op
	if err := stream.Close(); err != nil {
		t.Errorf("second Close() error: %v", err)
	}
}

func TestKimiStream_RecvAfterClose(t *testing.T) {
	srv := newStreamServer(t)
	defer srv.Close()

	c := New("key", srv.URL, nil)
	stream, err := c.CreateCompletionStream(context.Background(), &Request{
		Model:    "moonshot-v1-8k",
		Messages: []Message{{Role: "user", Content: "hello"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stream.Close()
	_, err = stream.Recv()
	if err == nil || !strings.Contains(err.Error(), "stream is closed") {
		t.Errorf("expected 'stream is closed' error, got %v", err)
	}
}

// ─── Adapter (Provider) Unit Tests ──────────────────────────────────────────

func TestKimiProvider_Name(t *testing.T) {
	p := NewProvider("key", "", nil)
	if p.Name() != "kimi" {
		t.Errorf("Name() = %q, want kimi", p.Name())
	}
}

func TestKimiProvider_CreateChatCompletion_Success(t *testing.T) {
	srv := newTestServer(t, http.StatusOK, chatCompletionResponse("hello there"))
	defer srv.Close()

	p := NewProvider("key", srv.URL, nil)
	resp, err := p.CreateChatCompletion(context.Background(), &provider.ChatCompletionRequest{
		Model: "moonshot-v1-8k",
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: "hello"},
		},
	})
	if err != nil {
		t.Fatalf("CreateChatCompletion failed: %v", err)
	}
	if resp.ID != "test-id-123" {
		t.Errorf("ID = %q, want test-id-123", resp.ID)
	}
	if len(resp.Choices) == 0 {
		t.Fatal("no choices")
	}
	if resp.Choices[0].Message.Content != "hello there" {
		t.Errorf("content = %q, want 'hello there'", resp.Choices[0].Message.Content)
	}
	if resp.Usage.TotalTokens != 30 {
		t.Errorf("TotalTokens = %d, want 30", resp.Usage.TotalTokens)
	}
}

func TestKimiProvider_CreateChatCompletionStream_Success(t *testing.T) {
	srv := newStreamServer(t)
	defer srv.Close()

	p := NewProvider("key", srv.URL, nil)
	stream, err := p.CreateChatCompletionStream(context.Background(), &provider.ChatCompletionRequest{
		Model: "moonshot-v1-8k",
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: "hello"},
		},
	})
	if err != nil {
		t.Fatalf("CreateChatCompletionStream failed: %v", err)
	}
	defer stream.Close()

	var content strings.Builder
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv error: %v", err)
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			content.WriteString(chunk.Choices[0].Delta.Content)
		}
	}
	if content.String() != "Hello world" {
		t.Errorf("content = %q, want 'Hello world'", content.String())
	}
}

func TestKimiProvider_Close(t *testing.T) {
	p := NewProvider("key", "", nil)
	if err := p.Close(); err != nil {
		t.Errorf("Close() error: %v", err)
	}
}

func TestKimiProvider_CreateChatCompletion_EmptyChoices(t *testing.T) {
	body := `{"id":"x","object":"chat.completion","created":1,"model":"moonshot-v1-8k","choices":[],"usage":{}}`
	srv := newTestServer(t, http.StatusOK, body)
	defer srv.Close()

	p := NewProvider("key", srv.URL, nil)
	_, err := p.CreateChatCompletion(context.Background(), &provider.ChatCompletionRequest{
		Model:    testModel,
		Messages: []provider.Message{{Role: provider.RoleUser, Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error for empty choices, got nil")
	}
	if !strings.Contains(err.Error(), "empty choices") {
		t.Errorf("expected 'empty choices' error, got: %v", err)
	}
}

// ─── CustomHTTPClient Test ───────────────────────────────────────────────────

func TestKimiClient_CustomHTTPClient(t *testing.T) {
	srv := newTestServer(t, http.StatusOK, chatCompletionResponse("ok"))
	defer srv.Close()

	custom := &http.Client{Timeout: 30 * time.Second}
	c := New("key", srv.URL, custom)
	if c.client != custom {
		t.Error("expected custom http.Client to be used")
	}
	_, err := c.CreateCompletion(context.Background(), &Request{
		Model:    "moonshot-v1-8k",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
