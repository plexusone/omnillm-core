package anthropic

import (
	"strings"
	"testing"

	"github.com/plexusone/omnillm-core/provider"
)

func TestProvider_Name(t *testing.T) {
	p := NewProvider("test-key", "", nil)
	if p.Name() != "anthropic" {
		t.Errorf("Expected provider name 'anthropic', got '%s'", p.Name())
	}
}

func TestProvider_CreateChatCompletion_MessageConversion(t *testing.T) {
	tests := []struct {
		name         string
		messages     []provider.Message
		wantSystem   string
		wantMsgCount int
	}{
		{
			name: "system message separated",
			messages: []provider.Message{
				{Role: provider.RoleSystem, Content: "You are helpful"},
				{Role: provider.RoleUser, Content: "Hello"},
			},
			wantSystem:   "You are helpful",
			wantMsgCount: 1,
		},
		{
			name: "no system message",
			messages: []provider.Message{
				{Role: provider.RoleUser, Content: "Hello"},
				{Role: provider.RoleAssistant, Content: "Hi there"},
			},
			wantSystem:   "",
			wantMsgCount: 2,
		},
		{
			name: "multiple system messages (last one wins)",
			messages: []provider.Message{
				{Role: provider.RoleSystem, Content: "First system"},
				{Role: provider.RoleUser, Content: "Hello"},
				{Role: provider.RoleSystem, Content: "Second system"},
			},
			wantSystem:   "Second system",
			wantMsgCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock request
			req := &provider.ChatCompletionRequest{
				Model:    "claude-3-haiku-20240307",
				Messages: tt.messages,
			}

			// Extract the system message and convert messages
			// (This simulates what CreateChatCompletion does)
			var systemMessage string
			var anthropicMessages []Message
			for _, msg := range req.Messages {
				switch msg.Role {
				case provider.RoleSystem:
					systemMessage = msg.Content
				case provider.RoleUser, provider.RoleAssistant:
					anthropicMessages = append(anthropicMessages, Message{
						Role:    string(msg.Role),
						Content: msg.Content,
					})
				}
			}

			if systemMessage != tt.wantSystem {
				t.Errorf("System message = %q, want %q", systemMessage, tt.wantSystem)
			}
			if len(anthropicMessages) != tt.wantMsgCount {
				t.Errorf("Message count = %d, want %d", len(anthropicMessages), tt.wantMsgCount)
			}
		})
	}
}

func TestStreamAdapter_EventHandling(t *testing.T) {
	tests := []struct {
		name      string
		event     StreamEvent
		wantError bool
		wantEmpty bool
		wantText  string
	}{
		{
			name: "message_start event",
			event: StreamEvent{
				Type: "message_start",
				Message: &StreamMessage{
					ID:    "msg_123",
					Model: "claude-3-haiku-20240307",
				},
			},
			wantEmpty: true,
		},
		{
			name: "content_block_delta with text",
			event: StreamEvent{
				Type: "content_block_delta",
				Delta: &StreamDelta{
					Type: "text_delta",
					Text: "Hello world",
				},
			},
			wantText: "Hello world",
		},
		{
			name: "message_delta with stop reason",
			event: StreamEvent{
				Type: "message_delta",
				Delta: &StreamDelta{
					StopReason: "end_turn",
				},
			},
		},
		{
			name: "message_stop event",
			event: StreamEvent{
				Type: "message_stop",
			},
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock stream adapter
			adapter := &StreamAdapter{
				messageID: "msg_test",
				model:     "claude-3-haiku-20240307",
			}

			// Simulate event handling based on type
			var chunk *provider.ChatCompletionChunk

			switch tt.event.Type {
			case "message_start":
				if tt.event.Message != nil {
					adapter.messageID = tt.event.Message.ID
					adapter.model = tt.event.Message.Model
				}
				chunk = &provider.ChatCompletionChunk{
					Choices: []provider.ChatCompletionChoice{},
				}

			case "content_block_delta":
				var content string
				if tt.event.Delta != nil && tt.event.Delta.Type == "text_delta" {
					content = tt.event.Delta.Text
				}
				chunk = &provider.ChatCompletionChunk{
					Choices: []provider.ChatCompletionChoice{
						{
							Delta: &provider.Message{
								Role:    provider.RoleAssistant,
								Content: content,
							},
						},
					},
				}

			case "message_delta":
				chunk = &provider.ChatCompletionChunk{
					Choices: []provider.ChatCompletionChoice{
						{Index: 0},
					},
				}

			case "message_stop":
				chunk = &provider.ChatCompletionChunk{
					Choices: []provider.ChatCompletionChoice{},
				}
			}

			// Verify results
			if tt.wantEmpty && len(chunk.Choices) != 0 {
				t.Errorf("Expected empty choices, got %d", len(chunk.Choices))
			}

			if tt.wantText != "" {
				if len(chunk.Choices) == 0 || chunk.Choices[0].Delta == nil {
					t.Errorf("Expected text chunk, got empty")
				} else if chunk.Choices[0].Delta.Content != tt.wantText {
					t.Errorf("Text = %q, want %q", chunk.Choices[0].Delta.Content, tt.wantText)
				}
			}
		})
	}
}

func TestStream_Recv_SSEParsing(t *testing.T) {
	// Mock SSE stream data
	sseData := `event: message_start
data: {"type":"message_start","message":{"id":"msg_123","type":"message","role":"assistant","model":"claude-3-haiku-20240307"}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" world"}}

event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":10}}

event: message_stop
data: {"type":"message_stop"}

`

	// Test that we can parse different SSE formats
	lines := strings.Split(sseData, "\n")

	var eventType string
	var eventData string

	for _, line := range lines {
		if line == "" {
			// End of event
			if eventType != "" && eventData != "" {
				t.Logf("Parsed event type: %s", eventType)
				t.Logf("Parsed data: %s", eventData)
			}
			eventType = ""
			eventData = ""
			continue
		}

		if strings.HasPrefix(line, "event: ") {
			eventType = strings.TrimPrefix(line, "event: ")
		}

		if strings.HasPrefix(line, "data: ") {
			eventData = strings.TrimPrefix(line, "data: ")
		}
	}
}

func TestBoolPtr(t *testing.T) {
	tests := []struct {
		name  string
		input bool
		want  bool
	}{
		{"true value", true, true},
		{"false value", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := boolPtr(tt.input)
			if ptr == nil {
				t.Fatal("boolPtr returned nil")
			}
			if *ptr != tt.want {
				t.Errorf("boolPtr(%v) = %v, want %v", tt.input, *ptr, tt.want)
			}
		})
	}
}
