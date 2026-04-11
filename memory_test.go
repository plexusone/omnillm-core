package omnillm

import (
	"context"
	"testing"
	"time"

	mocktest "github.com/plexusone/omnillm-core/testing"
)

func TestMemoryManager_LoadConversation(t *testing.T) {
	mockKVS := mocktest.NewMockKVS()
	config := MemoryConfig{
		MaxMessages: 10,
		TTL:         time.Hour,
		KeyPrefix:   "test",
	}
	mm := NewMemoryManager(mockKVS, config)

	ctx := context.Background()

	// Test loading non-existent conversation
	conv, err := mm.LoadConversation(ctx, "session1")
	if err != nil {
		t.Fatalf("LoadConversation failed: %v", err)
	}
	if conv.SessionID != "session1" {
		t.Errorf("SessionID = %s, want session1", conv.SessionID)
	}
	if len(conv.Messages) != 0 {
		t.Errorf("New conversation has %d messages, want 0", len(conv.Messages))
	}
}

func TestMemoryManager_SaveAndLoadConversation(t *testing.T) {
	mockKVS := mocktest.NewMockKVS()
	config := DefaultMemoryConfig()
	mm := NewMemoryManager(mockKVS, config)

	ctx := context.Background()

	// Create and save a conversation
	conv := &ConversationMemory{
		SessionID: "session1",
		Messages: []Message{
			{Role: RoleUser, Content: "Hello"},
			{Role: RoleAssistant, Content: "Hi there!"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]any),
	}

	err := mm.SaveConversation(ctx, conv)
	if err != nil {
		t.Fatalf("SaveConversation failed: %v", err)
	}

	// Load it back
	loaded, err := mm.LoadConversation(ctx, "session1")
	if err != nil {
		t.Fatalf("LoadConversation failed: %v", err)
	}

	if loaded.SessionID != "session1" {
		t.Errorf("SessionID = %s, want session1", loaded.SessionID)
	}
	if len(loaded.Messages) != 2 {
		t.Fatalf("Messages count = %d, want 2", len(loaded.Messages))
	}
	if loaded.Messages[0].Content != "Hello" {
		t.Errorf("First message = %s, want Hello", loaded.Messages[0].Content)
	}
	if loaded.Messages[1].Content != "Hi there!" {
		t.Errorf("Second message = %s, want Hi there!", loaded.Messages[1].Content)
	}
}

func TestMemoryManager_AppendMessage(t *testing.T) {
	mockKVS := mocktest.NewMockKVS()
	config := DefaultMemoryConfig()
	mm := NewMemoryManager(mockKVS, config)

	ctx := context.Background()
	sessionID := "session1"

	// Append first message
	err := mm.AppendMessage(ctx, sessionID, Message{
		Role:    RoleUser,
		Content: "First message",
	})
	if err != nil {
		t.Fatalf("AppendMessage failed: %v", err)
	}

	// Append second message
	err = mm.AppendMessage(ctx, sessionID, Message{
		Role:    RoleAssistant,
		Content: "Second message",
	})
	if err != nil {
		t.Fatalf("AppendMessage failed: %v", err)
	}

	// Load and verify
	conv, err := mm.LoadConversation(ctx, sessionID)
	if err != nil {
		t.Fatalf("LoadConversation failed: %v", err)
	}

	if len(conv.Messages) != 2 {
		t.Fatalf("Messages count = %d, want 2", len(conv.Messages))
	}
}

func TestMemoryManager_AppendMessages(t *testing.T) {
	mockKVS := mocktest.NewMockKVS()
	config := DefaultMemoryConfig()
	mm := NewMemoryManager(mockKVS, config)

	ctx := context.Background()
	sessionID := "session1"

	messages := []Message{
		{Role: RoleUser, Content: "Message 1"},
		{Role: RoleAssistant, Content: "Message 2"},
		{Role: RoleUser, Content: "Message 3"},
	}

	err := mm.AppendMessages(ctx, sessionID, messages)
	if err != nil {
		t.Fatalf("AppendMessages failed: %v", err)
	}

	// Load and verify
	conv, err := mm.LoadConversation(ctx, sessionID)
	if err != nil {
		t.Fatalf("LoadConversation failed: %v", err)
	}

	if len(conv.Messages) != 3 {
		t.Fatalf("Messages count = %d, want 3", len(conv.Messages))
	}
}

func TestMemoryManager_MaxMessages(t *testing.T) {
	mockKVS := mocktest.NewMockKVS()
	config := MemoryConfig{
		MaxMessages: 5,
		TTL:         time.Hour,
		KeyPrefix:   "test",
	}
	mm := NewMemoryManager(mockKVS, config)

	ctx := context.Background()
	sessionID := "session1"

	// Add system message
	err := mm.AppendMessage(ctx, sessionID, Message{
		Role:    RoleSystem,
		Content: "You are helpful",
	})
	if err != nil {
		t.Fatalf("AppendMessage failed: %v", err)
	}

	// Add more messages than the limit
	for i := 0; i < 10; i++ {
		err := mm.AppendMessage(ctx, sessionID, Message{
			Role:    RoleUser,
			Content: "Message " + string(rune('A'+i)),
		})
		if err != nil {
			t.Fatalf("AppendMessage failed: %v", err)
		}
	}

	// Load and verify
	conv, err := mm.LoadConversation(ctx, sessionID)
	if err != nil {
		t.Fatalf("LoadConversation failed: %v", err)
	}

	if len(conv.Messages) > config.MaxMessages {
		t.Errorf("Messages count = %d, want <= %d", len(conv.Messages), config.MaxMessages)
	}

	// Verify system message is preserved
	hasSystem := false
	for _, msg := range conv.Messages {
		if msg.Role == RoleSystem {
			hasSystem = true
			break
		}
	}
	if !hasSystem {
		t.Error("System message was not preserved")
	}
}

func TestMemoryManager_GetMessages(t *testing.T) {
	mockKVS := mocktest.NewMockKVS()
	config := DefaultMemoryConfig()
	mm := NewMemoryManager(mockKVS, config)

	ctx := context.Background()
	sessionID := "session1"

	// Add messages
	messages := []Message{
		{Role: RoleUser, Content: "Hello"},
		{Role: RoleAssistant, Content: "Hi"},
	}
	err := mm.AppendMessages(ctx, sessionID, messages)
	if err != nil {
		t.Fatalf("AppendMessages failed: %v", err)
	}

	// Get messages
	retrieved, err := mm.GetMessages(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetMessages failed: %v", err)
	}

	if len(retrieved) != 2 {
		t.Fatalf("Messages count = %d, want 2", len(retrieved))
	}
}

func TestMemoryManager_CreateConversationWithSystemMessage(t *testing.T) {
	mockKVS := mocktest.NewMockKVS()
	config := DefaultMemoryConfig()
	mm := NewMemoryManager(mockKVS, config)

	ctx := context.Background()
	sessionID := "session1"

	err := mm.CreateConversationWithSystemMessage(ctx, sessionID, "You are a helpful assistant")
	if err != nil {
		t.Fatalf("CreateConversationWithSystemMessage failed: %v", err)
	}

	// Load and verify
	conv, err := mm.LoadConversation(ctx, sessionID)
	if err != nil {
		t.Fatalf("LoadConversation failed: %v", err)
	}

	if len(conv.Messages) != 1 {
		t.Fatalf("Messages count = %d, want 1", len(conv.Messages))
	}
	if conv.Messages[0].Role != RoleSystem {
		t.Errorf("Message role = %s, want %s", conv.Messages[0].Role, RoleSystem)
	}
	if conv.Messages[0].Content != "You are a helpful assistant" {
		t.Errorf("Message content = %s, want 'You are a helpful assistant'", conv.Messages[0].Content)
	}
}

func TestMemoryManager_SetMetadata(t *testing.T) {
	mockKVS := mocktest.NewMockKVS()
	config := DefaultMemoryConfig()
	mm := NewMemoryManager(mockKVS, config)

	ctx := context.Background()
	sessionID := "session1"

	// Set metadata
	metadata := map[string]any{
		"user_id": "user123",
		"tags":    []string{"test", "demo"},
	}
	err := mm.SetMetadata(ctx, sessionID, metadata)
	if err != nil {
		t.Fatalf("SetMetadata failed: %v", err)
	}

	// Load and verify
	conv, err := mm.LoadConversation(ctx, sessionID)
	if err != nil {
		t.Fatalf("LoadConversation failed: %v", err)
	}

	if conv.Metadata["user_id"] != "user123" {
		t.Errorf("Metadata user_id = %v, want user123", conv.Metadata["user_id"])
	}
}

func TestMemoryManager_BuildKey(t *testing.T) {
	config := MemoryConfig{
		KeyPrefix: "myapp:chat",
	}
	mm := NewMemoryManager(nil, config)

	key := mm.buildKey("session123")
	expected := "myapp:chat:session123"
	if key != expected {
		t.Errorf("buildKey = %s, want %s", key, expected)
	}
}

func TestDefaultMemoryConfig(t *testing.T) {
	config := DefaultMemoryConfig()

	if config.MaxMessages != 50 {
		t.Errorf("MaxMessages = %d, want 50", config.MaxMessages)
	}
	if config.TTL != 24*time.Hour {
		t.Errorf("TTL = %v, want 24h", config.TTL)
	}
	if config.KeyPrefix != "omnillm:session" {
		t.Errorf("KeyPrefix = %s, want omnillm:session", config.KeyPrefix)
	}
}
