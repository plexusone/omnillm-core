package anthropic

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/plexusone/omnillm-core/provider"
)

// TestAnthropicIntegration_ChatCompletion tests actual API calls
func TestAnthropicIntegration_ChatCompletion(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: ANTHROPIC_API_KEY not set")
	}

	p := NewProvider(apiKey, "", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: "claude-3-haiku-20240307",
		Messages: []provider.Message{
			{
				Role:    provider.RoleUser,
				Content: "Say 'test successful' if you can read this.",
			},
		},
		MaxTokens:   intPtr(50),
		Temperature: float64Ptr(0.5),
	}

	resp, err := p.CreateChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletion failed: %v", err)
	}

	// Verify response structure
	if resp.ID == "" {
		t.Error("Response ID is empty")
	}
	if resp.Model == "" {
		t.Error("Response model is empty")
	}
	if len(resp.Choices) == 0 {
		t.Fatal("No choices in response")
	}
	if resp.Choices[0].Message.Content == "" {
		t.Error("Response content is empty")
	}
	if resp.Usage.TotalTokens == 0 {
		t.Error("Usage tokens is zero")
	}

	t.Logf("Response: %s", resp.Choices[0].Message.Content)
	t.Logf("Tokens used: %d", resp.Usage.TotalTokens)
}

// TestAnthropicIntegration_ChatCompletionWithSystemMessage tests system message handling
func TestAnthropicIntegration_ChatCompletionWithSystemMessage(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: ANTHROPIC_API_KEY not set")
	}

	p := NewProvider(apiKey, "", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: "claude-3-haiku-20240307",
		Messages: []provider.Message{
			{
				Role:    provider.RoleSystem,
				Content: "You are a helpful assistant. Always respond with exactly 3 words.",
			},
			{
				Role:    provider.RoleUser,
				Content: "What is AI?",
			},
		},
		MaxTokens:   intPtr(20),
		Temperature: float64Ptr(0.7),
	}

	resp, err := p.CreateChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletion failed: %v", err)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("No choices in response")
	}

	t.Logf("Response with system message: %s", resp.Choices[0].Message.Content)
}

// TestAnthropicIntegration_Streaming tests actual streaming API calls
func TestAnthropicIntegration_Streaming(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: ANTHROPIC_API_KEY not set")
	}

	p := NewProvider(apiKey, "", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: "claude-3-haiku-20240307",
		Messages: []provider.Message{
			{
				Role:    provider.RoleUser,
				Content: "Count from 1 to 5, one number per line.",
			},
		},
		MaxTokens:   intPtr(50),
		Temperature: float64Ptr(0.5),
	}

	stream, err := p.CreateChatCompletionStream(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletionStream failed: %v", err)
	}
	defer stream.Close()

	var chunks []string
	var totalContent string
	chunkCount := 0

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Stream recv error: %v", err)
		}

		chunkCount++

		// Verify chunk structure
		if chunk.ID == "" && chunkCount > 1 {
			// First chunk might not have ID yet
			t.Logf("Warning: chunk %d has no ID", chunkCount)
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				chunks = append(chunks, content)
				totalContent += content
			}
		}
	}

	if chunkCount == 0 {
		t.Fatal("No chunks received from stream")
	}

	if len(chunks) == 0 {
		t.Fatal("No content chunks received")
	}

	t.Logf("Received %d total chunks, %d with content", chunkCount, len(chunks))
	t.Logf("Complete response: %s", totalContent)
}

// TestAnthropicIntegration_StreamingWithSystemMessage tests streaming with system messages
func TestAnthropicIntegration_StreamingWithSystemMessage(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: ANTHROPIC_API_KEY not set")
	}

	p := NewProvider(apiKey, "", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: "claude-3-haiku-20240307",
		Messages: []provider.Message{
			{
				Role:    provider.RoleSystem,
				Content: "You are a poet. Write concisely.",
			},
			{
				Role:    provider.RoleUser,
				Content: "Write a two-line poem about testing.",
			},
		},
		MaxTokens:   intPtr(100),
		Temperature: float64Ptr(0.8),
	}

	stream, err := p.CreateChatCompletionStream(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletionStream failed: %v", err)
	}
	defer stream.Close()

	var fullResponse string
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Stream recv error: %v", err)
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fullResponse += chunk.Choices[0].Delta.Content
		}
	}

	if fullResponse == "" {
		t.Fatal("No response content received")
	}

	t.Logf("Streamed poem: %s", fullResponse)
}

// TestAnthropicIntegration_ErrorHandling tests API error responses
func TestAnthropicIntegration_ErrorHandling(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: ANTHROPIC_API_KEY not set")
	}

	tests := []struct {
		name      string
		request   *provider.ChatCompletionRequest
		wantError bool
	}{
		{
			name: "empty model",
			request: &provider.ChatCompletionRequest{
				Model: "",
				Messages: []provider.Message{
					{Role: provider.RoleUser, Content: "Hello"},
				},
			},
			wantError: true,
		},
		{
			name: "empty messages",
			request: &provider.ChatCompletionRequest{
				Model:    "claude-3-haiku-20240307",
				Messages: []provider.Message{},
			},
			wantError: true,
		},
		{
			name: "invalid model name",
			request: &provider.ChatCompletionRequest{
				Model: "nonexistent-model-12345",
				Messages: []provider.Message{
					{Role: provider.RoleUser, Content: "Hello"},
				},
				MaxTokens: intPtr(10),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProvider(apiKey, "", nil)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err := p.CreateChatCompletion(ctx, tt.request)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if err != nil {
				t.Logf("Expected error received: %v", err)
			}
		})
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
