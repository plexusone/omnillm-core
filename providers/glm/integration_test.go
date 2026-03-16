package glm

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/plexusone/omnillm/provider"
)

// TestGLMIntegration_ChatCompletion tests actual API calls to GLM
func TestGLMIntegration_ChatCompletion(t *testing.T) {
	apiKey := os.Getenv("GLM_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: GLM_API_KEY not set")
	}

	p := NewProvider(apiKey, "", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: "glm-4.7-flash",
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

	t.Logf("Response: %s", resp.Choices[0].Message.Content)
	t.Logf("Tokens used: %d", resp.Usage.TotalTokens)
}

// TestGLMIntegration_Streaming tests actual streaming API calls to GLM
func TestGLMIntegration_Streaming(t *testing.T) {
	apiKey := os.Getenv("GLM_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: GLM_API_KEY not set")
	}

	p := NewProvider(apiKey, "", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: "glm-4.7-flash",
		Messages: []provider.Message{
			{
				Role:    provider.RoleUser,
				Content: "Count from 1 to 5, one number per line.",
			},
		},
		MaxTokens:   intPtr(100),
		Temperature: float64Ptr(0.5),
	}

	stream, err := p.CreateChatCompletionStream(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletionStream failed: %v", err)
	}
	defer stream.Close()

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
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			totalContent += chunk.Choices[0].Delta.Content
		}
	}

	if chunkCount == 0 {
		t.Fatal("No chunks received from stream")
	}

	t.Logf("Received %d chunks", chunkCount)
	t.Logf("Complete response: %s", totalContent)
}

// TestGLMIntegration_ErrorHandling tests API error responses
func TestGLMIntegration_ErrorHandling(t *testing.T) {
	apiKey := os.Getenv("GLM_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: GLM_API_KEY not set")
	}

	tests := []struct {
		name      string
		request   *provider.ChatCompletionRequest
		wantError bool
	}{
		{
			name: "empty model",
			request: &provider.ChatCompletionRequest{
				Model:    "",
				Messages: []provider.Message{{Role: provider.RoleUser, Content: "Hello"}},
			},
			wantError: true,
		},
		{
			name: "empty messages",
			request: &provider.ChatCompletionRequest{
				Model:    "glm-4.7-flash",
				Messages: []provider.Message{},
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
