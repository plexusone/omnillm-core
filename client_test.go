package omnillm

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/plexusone/omnillm-core/provider"
	mocktest "github.com/plexusone/omnillm-core/testing"
)

// MockProvider implements provider.Provider for testing
type MockProvider struct {
	name                   string
	completionError        error
	streamError            error
	completionResp         *provider.ChatCompletionResponse
	streamChunks           []*provider.ChatCompletionChunk
	createCompletionCalled bool
	createStreamCalled     bool
}

func NewMockProvider(name string) *MockProvider {
	return &MockProvider{
		name: name,
		completionResp: &provider.ChatCompletionResponse{
			ID:      "test-id",
			Model:   "test-model",
			Created: time.Now().Unix(),
			Choices: []provider.ChatCompletionChoice{
				{
					Index: 0,
					Message: provider.Message{
						Role:    provider.RoleAssistant,
						Content: "Mock response",
					},
					FinishReason: stringPtr("stop"),
				},
			},
			Usage: provider.Usage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
		},
	}
}

func (m *MockProvider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
	m.createCompletionCalled = true
	if m.completionError != nil {
		return nil, m.completionError
	}
	return m.completionResp, nil
}

func (m *MockProvider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
	m.createStreamCalled = true
	if m.streamError != nil {
		return nil, m.streamError
	}
	return &MockStream{chunks: m.streamChunks}, nil
}

func (m *MockProvider) Close() error {
	return nil
}

func (m *MockProvider) Name() string {
	return m.name
}

// MockStream implements provider.ChatCompletionStream for testing
type MockStream struct {
	chunks []*provider.ChatCompletionChunk
	index  int
	closed bool
}

func (m *MockStream) Recv() (*provider.ChatCompletionChunk, error) {
	if m.closed {
		return nil, io.EOF
	}
	if m.index >= len(m.chunks) {
		return nil, io.EOF
	}
	chunk := m.chunks[m.index]
	m.index++
	return chunk, nil
}

func (m *MockStream) Close() error {
	m.closed = true
	return nil
}

func TestNewClient_CustomProvider(t *testing.T) {
	mockProv := NewMockProvider("test-provider")

	client, err := NewClient(ClientConfig{
		Providers: []ProviderConfig{
			{CustomProvider: mockProv},
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	if client.Provider().Name() != "test-provider" {
		t.Errorf("Provider name = %s, want test-provider", client.Provider().Name())
	}
}

func TestNewClient_UnsupportedProvider(t *testing.T) {
	_, err := NewClient(ClientConfig{
		Providers: []ProviderConfig{
			{Provider: "unsupported-provider"},
		},
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestNewClient_NoProviders(t *testing.T) {
	_, err := NewClient(ClientConfig{})
	if err != ErrNoProviders {
		t.Errorf("Expected ErrNoProviders, got %v", err)
	}
}

func TestChatClient_CreateChatCompletion(t *testing.T) {
	mockProv := NewMockProvider("test")
	client := &ChatClient{provider: mockProv}

	ctx := context.Background()
	req := &provider.ChatCompletionRequest{
		Model: "test-model",
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: "Hello"},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletion failed: %v", err)
	}

	if !mockProv.createCompletionCalled {
		t.Error("CreateChatCompletion was not called on provider")
	}
	if resp.ID != "test-id" {
		t.Errorf("Response ID = %s, want test-id", resp.ID)
	}
	if len(resp.Choices) == 0 {
		t.Fatal("No choices in response")
	}
	if resp.Choices[0].Message.Content != "Mock response" {
		t.Errorf("Response content = %s, want 'Mock response'", resp.Choices[0].Message.Content)
	}
}

func TestChatClient_CreateChatCompletionStream(t *testing.T) {
	mockProv := NewMockProvider("test")
	mockProv.streamChunks = []*provider.ChatCompletionChunk{
		{
			ID:    "chunk1",
			Model: "test-model",
			Choices: []provider.ChatCompletionChoice{
				{
					Delta: &provider.Message{Content: "Hello"},
				},
			},
		},
		{
			ID:    "chunk2",
			Model: "test-model",
			Choices: []provider.ChatCompletionChoice{
				{
					Delta: &provider.Message{Content: " world"},
				},
			},
		},
	}

	client := &ChatClient{provider: mockProv}

	ctx := context.Background()
	req := &provider.ChatCompletionRequest{
		Model: "test-model",
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: "Hello"},
		},
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletionStream failed: %v", err)
	}
	defer stream.Close()

	if !mockProv.createStreamCalled {
		t.Error("CreateChatCompletionStream was not called on provider")
	}

	// Read chunks
	var fullContent string
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Stream recv error: %v", err)
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fullContent += chunk.Choices[0].Delta.Content
		}
	}

	if fullContent != "Hello world" {
		t.Errorf("Full content = %s, want 'Hello world'", fullContent)
	}
}

func TestChatClient_WithMemory(t *testing.T) {
	mockProv := NewMockProvider("test")
	mockKVS := mocktest.NewMockKVS()

	client, err := NewClient(ClientConfig{
		Providers: []ProviderConfig{
			{CustomProvider: mockProv},
		},
		Memory:       mockKVS,
		MemoryConfig: &MemoryConfig{MaxMessages: 10, KeyPrefix: "test"},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	if !client.HasMemory() {
		t.Error("Client should have memory configured")
	}

	if client.Memory() == nil {
		t.Error("Memory() returned nil")
	}
}

func TestChatClient_CreateChatCompletionWithMemory(t *testing.T) {
	mockProv := NewMockProvider("test")
	mockKVS := mocktest.NewMockKVS()

	client, err := NewClient(ClientConfig{
		Providers: []ProviderConfig{
			{CustomProvider: mockProv},
		},
		Memory: mockKVS,
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	sessionID := "session1"

	// Create conversation with system message
	err = client.CreateConversationWithSystemMessage(ctx, sessionID, "You are helpful")
	if err != nil {
		t.Fatalf("CreateConversationWithSystemMessage failed: %v", err)
	}

	// Make a completion request
	req := &provider.ChatCompletionRequest{
		Model: "test-model",
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: "Hello"},
		},
	}

	resp, err := client.CreateChatCompletionWithMemory(ctx, sessionID, req)
	if err != nil {
		t.Fatalf("CreateChatCompletionWithMemory failed: %v", err)
	}

	if resp.Choices[0].Message.Content != "Mock response" {
		t.Errorf("Response content = %s, want 'Mock response'", resp.Choices[0].Message.Content)
	}

	// Verify conversation was saved
	messages, err := client.GetConversationMessages(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetConversationMessages failed: %v", err)
	}

	// Should have: system message + user message + assistant response
	if len(messages) != 3 {
		t.Errorf("Messages count = %d, want 3", len(messages))
	}
}

func TestChatClient_CreateChatCompletionStreamWithMemory(t *testing.T) {
	mockProv := NewMockProvider("test")
	mockProv.streamChunks = []*provider.ChatCompletionChunk{
		{
			Choices: []provider.ChatCompletionChoice{
				{Delta: &provider.Message{Content: "Streaming"}},
			},
		},
		{
			Choices: []provider.ChatCompletionChoice{
				{Delta: &provider.Message{Content: " response"}},
			},
		},
	}

	mockKVS := mocktest.NewMockKVS()

	client, err := NewClient(ClientConfig{
		Providers: []ProviderConfig{
			{CustomProvider: mockProv},
		},
		Memory: mockKVS,
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	sessionID := "session1"

	req := &provider.ChatCompletionRequest{
		Model: "test-model",
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: "Hello"},
		},
	}

	stream, err := client.CreateChatCompletionStreamWithMemory(ctx, sessionID, req)
	if err != nil {
		t.Fatalf("CreateChatCompletionStreamWithMemory failed: %v", err)
	}

	var fullContent string
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Stream recv error: %v", err)
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fullContent += chunk.Choices[0].Delta.Content
		}
	}
	stream.Close()

	if fullContent != "Streaming response" {
		t.Errorf("Full content = %s, want 'Streaming response'", fullContent)
	}

	// Verify conversation was saved
	messages, err := client.GetConversationMessages(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetConversationMessages failed: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("Messages count = %d, want 2", len(messages))
	}
}

func TestChatClient_ConversationManagement(t *testing.T) {
	mockProv := NewMockProvider("test")
	mockKVS := mocktest.NewMockKVS()

	client, err := NewClient(ClientConfig{
		Providers: []ProviderConfig{
			{CustomProvider: mockProv},
		},
		Memory: mockKVS,
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	sessionID := "session1"

	// Append a message
	err = client.AppendMessage(ctx, sessionID, provider.Message{
		Role:    provider.RoleUser,
		Content: "Test message",
	})
	if err != nil {
		t.Fatalf("AppendMessage failed: %v", err)
	}

	// Load conversation
	conv, err := client.LoadConversation(ctx, sessionID)
	if err != nil {
		t.Fatalf("LoadConversation failed: %v", err)
	}

	if len(conv.Messages) != 1 {
		t.Errorf("Messages count = %d, want 1", len(conv.Messages))
	}

	// Delete conversation
	err = client.DeleteConversation(ctx, sessionID)
	if err != nil {
		t.Fatalf("DeleteConversation failed: %v", err)
	}
}

func TestChatClient_NoMemory(t *testing.T) {
	mockProv := NewMockProvider("test")

	client, err := NewClient(ClientConfig{
		Providers: []ProviderConfig{
			{CustomProvider: mockProv},
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	if client.HasMemory() {
		t.Error("Client should not have memory configured")
	}

	ctx := context.Background()

	// These should return errors when memory is not configured
	_, err = client.LoadConversation(ctx, "session1")
	if err == nil {
		t.Error("LoadConversation should fail without memory")
	}

	err = client.AppendMessage(ctx, "session1", provider.Message{})
	if err == nil {
		t.Error("AppendMessage should fail without memory")
	}

	err = client.DeleteConversation(ctx, "session1")
	if err == nil {
		t.Error("DeleteConversation should fail without memory")
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
