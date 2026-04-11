package omnillm

import (
	"testing"

	"github.com/plexusone/omnillm-core/provider"
)

func TestTokenEstimator_EstimateTokens(t *testing.T) {
	estimator := NewTokenEstimator(DefaultTokenEstimatorConfig())

	tests := []struct {
		name      string
		messages  []provider.Message
		minTokens int
		maxTokens int
	}{
		{
			name:      "empty messages",
			messages:  []provider.Message{},
			minTokens: 0,
			maxTokens: 0,
		},
		{
			name: "single short message",
			messages: []provider.Message{
				{Role: "user", Content: "Hello"},
			},
			minTokens: 1,
			maxTokens: 20,
		},
		{
			name: "conversation",
			messages: []provider.Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "What is the capital of France?"},
				{Role: "assistant", Content: "The capital of France is Paris."},
			},
			minTokens: 15,
			maxTokens: 50,
		},
		{
			name: "long message",
			messages: []provider.Message{
				{Role: "user", Content: "This is a much longer message that contains many words and should result in a higher token count than the shorter messages. We want to make sure the estimation scales appropriately with message length."},
			},
			minTokens: 30,
			maxTokens: 80,
		},
		{
			name: "message with tool calls",
			messages: []provider.Message{
				{
					Role:    "assistant",
					Content: "",
					ToolCalls: []provider.ToolCall{
						{
							ID:   "call_123",
							Type: "function",
							Function: provider.ToolFunction{
								Name:      "get_weather",
								Arguments: `{"location": "Paris", "unit": "celsius"}`,
							},
						},
					},
				},
			},
			minTokens: 10,
			maxTokens: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := estimator.EstimateTokens("gpt-4o", tt.messages)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tokens < tt.minTokens || tokens > tt.maxTokens {
				t.Errorf("token estimate %d not in expected range [%d, %d]",
					tokens, tt.minTokens, tt.maxTokens)
			}
		})
	}
}

func TestTokenEstimator_GetContextWindow(t *testing.T) {
	estimator := NewTokenEstimator(DefaultTokenEstimatorConfig())

	tests := []struct {
		model          string
		expectedWindow int
	}{
		{"gpt-4o", 128000},
		{"gpt-4o-mini", 128000},
		{"gpt-4", 8192},
		{"gpt-3.5-turbo", 16385},
		{"claude-3-opus", 200000},
		{"claude-3-sonnet", 200000},
		{"gemini-1.5-pro", 2000000},
		{"gemini-2.5-pro", 1000000},
		{"grok-4", 128000},
		{"llama3", 8192},
		{"mistral", 32768},
		{"unknown-model", 4096}, // Default fallback
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			window := estimator.GetContextWindow(tt.model)
			if window != tt.expectedWindow {
				t.Errorf("expected context window %d for %s, got %d",
					tt.expectedWindow, tt.model, window)
			}
		})
	}
}

func TestTokenEstimator_CustomContextWindows(t *testing.T) {
	config := TokenEstimatorConfig{
		CharactersPerToken: 4.0,
		CustomContextWindows: map[string]int{
			"my-custom-model": 1000000,
			"gpt-4o":          999999, // Override built-in
		},
	}
	estimator := NewTokenEstimator(config)

	// Custom model should use custom value
	if window := estimator.GetContextWindow("my-custom-model"); window != 1000000 {
		t.Errorf("expected 1000000 for custom model, got %d", window)
	}

	// Overridden model should use custom value
	if window := estimator.GetContextWindow("gpt-4o"); window != 999999 {
		t.Errorf("expected 999999 for overridden gpt-4o, got %d", window)
	}

	// Non-overridden model should use built-in
	if window := estimator.GetContextWindow("claude-3-opus"); window != 200000 {
		t.Errorf("expected 200000 for claude-3-opus, got %d", window)
	}
}

func TestTokenEstimator_CustomCharactersPerToken(t *testing.T) {
	// More conservative estimate (3 chars per token)
	conservative := NewTokenEstimator(TokenEstimatorConfig{
		CharactersPerToken: 3.0,
	})

	// Less conservative estimate (5 chars per token)
	liberal := NewTokenEstimator(TokenEstimatorConfig{
		CharactersPerToken: 5.0,
	})

	messages := []provider.Message{
		{Role: "user", Content: "This is a test message with some content."},
	}

	conservativeTokens, _ := conservative.EstimateTokens("gpt-4o", messages)
	liberalTokens, _ := liberal.EstimateTokens("gpt-4o", messages)

	if conservativeTokens <= liberalTokens {
		t.Errorf("conservative estimate (%d) should be higher than liberal (%d)",
			conservativeTokens, liberalTokens)
	}
}

func TestValidateTokens(t *testing.T) {
	estimator := NewTokenEstimator(DefaultTokenEstimatorConfig())

	tests := []struct {
		name                  string
		model                 string
		messages              []provider.Message
		maxCompletionTokens   int
		expectExceedsLimit    bool
		expectExceedsWithComp bool
	}{
		{
			name:  "small request fits",
			model: "gpt-4o",
			messages: []provider.Message{
				{Role: "user", Content: "Hello"},
			},
			maxCompletionTokens:   1000,
			expectExceedsLimit:    false,
			expectExceedsWithComp: false,
		},
		{
			name:  "request exceeds context",
			model: "phi", // Only 2048 context
			messages: func() []provider.Message {
				// Create a message with lots of content
				content := ""
				for i := 0; i < 10000; i++ {
					content += "word "
				}
				return []provider.Message{{Role: "user", Content: content}}
			}(),
			maxCompletionTokens: 100,
			expectExceedsLimit:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validation, err := ValidateTokens(estimator, tt.model, tt.messages, tt.maxCompletionTokens)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if validation.ExceedsLimit != tt.expectExceedsLimit {
				t.Errorf("ExceedsLimit: expected %v, got %v",
					tt.expectExceedsLimit, validation.ExceedsLimit)
			}
		})
	}
}

func TestTokenValidation_Fields(t *testing.T) {
	estimator := NewTokenEstimator(DefaultTokenEstimatorConfig())

	messages := []provider.Message{
		{Role: "user", Content: "Hello, how are you?"},
	}

	validation, err := ValidateTokens(estimator, "gpt-4o", messages, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if validation.EstimatedTokens <= 0 {
		t.Error("expected positive estimated tokens")
	}

	if validation.ContextWindow != 128000 {
		t.Errorf("expected context window 128000, got %d", validation.ContextWindow)
	}

	if validation.MaxCompletionTokens != 1000 {
		t.Errorf("expected max completion tokens 1000, got %d", validation.MaxCompletionTokens)
	}

	expectedAvailable := validation.ContextWindow - validation.EstimatedTokens
	if validation.AvailableTokens != expectedAvailable {
		t.Errorf("expected available tokens %d, got %d",
			expectedAvailable, validation.AvailableTokens)
	}
}

func TestTokenLimitError(t *testing.T) {
	err := &TokenLimitError{
		EstimatedTokens: 150000,
		ContextWindow:   128000,
		AvailableTokens: -22000,
		Model:           "gpt-4o",
	}

	expected := "estimated tokens (150000) exceed context window (128000) for model gpt-4o"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestDefaultTokenEstimatorConfig(t *testing.T) {
	config := DefaultTokenEstimatorConfig()

	if config.CharactersPerToken != 4.0 {
		t.Errorf("expected CharactersPerToken=4.0, got %f", config.CharactersPerToken)
	}

	if config.TokenOverheadPerMessage != 4 {
		t.Errorf("expected TokenOverheadPerMessage=4, got %d", config.TokenOverheadPerMessage)
	}
}

func TestEstimatePromptTokens(t *testing.T) {
	messages := []provider.Message{
		{Role: "user", Content: "Hello"},
	}

	tokens, err := EstimatePromptTokens("gpt-4o", messages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tokens <= 0 {
		t.Error("expected positive token count")
	}
}

func TestGetModelContextWindow(t *testing.T) {
	window := GetModelContextWindow("gpt-4o")
	if window != 128000 {
		t.Errorf("expected 128000, got %d", window)
	}
}
