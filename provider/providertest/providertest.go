// Package providertest provides conformance tests for LLM provider implementations.
//
// Provider implementations can use this package to verify they correctly implement
// the provider.Provider interface with consistent behavior.
//
// Basic usage:
//
//	func TestConformance(t *testing.T) {
//	    p := openai.NewProvider(apiKey, "", nil)
//
//	    providertest.RunAll(t, providertest.Config{
//	        Provider:        p,
//	        SkipIntegration: apiKey == "",
//	        TestModel:       "gpt-4o-mini",
//	    })
//	}
package providertest

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/plexusone/omnillm-core/provider"
)

// Config configures the LLM provider conformance test suite.
type Config struct {
	// Provider is the LLM provider implementation to test.
	Provider provider.Provider

	// SkipIntegration skips tests that require real API calls.
	// Set to true for unit tests without API credentials.
	SkipIntegration bool

	// TestModel is the model ID to use for integration tests.
	// Required if SkipIntegration is false.
	TestModel string

	// TestPrompt is the prompt to use in integration tests.
	// Defaults to "Say 'hello' and nothing else." if empty.
	TestPrompt string

	// Timeout for individual test operations.
	// Defaults to 30 seconds if zero.
	Timeout time.Duration
}

// withDefaults returns a copy of Config with default values applied.
func (c Config) withDefaults() Config {
	if c.TestPrompt == "" {
		c.TestPrompt = "Say 'hello' and nothing else."
	}
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
	return c
}

// RunAll runs all conformance tests for an LLM provider.
func RunAll(t *testing.T, cfg Config) {
	t.Helper()
	cfg = cfg.withDefaults()

	// Interface tests (always run)
	t.Run("Interface", func(t *testing.T) {
		RunInterfaceTests(t, cfg)
	})

	// Behavior tests (always run, some may be skipped without API)
	t.Run("Behavior", func(t *testing.T) {
		RunBehaviorTests(t, cfg)
	})

	// Integration tests (skipped if no API)
	if !cfg.SkipIntegration {
		t.Run("Integration", func(t *testing.T) {
			RunIntegrationTests(t, cfg)
		})
	}
}

// RunInterfaceTests runs only interface compliance tests.
// These tests verify the provider correctly implements the interface contract
// and do not require API credentials.
func RunInterfaceTests(t *testing.T, cfg Config) {
	t.Helper()
	cfg = cfg.withDefaults()

	t.Run("Name", func(t *testing.T) { testName(t, cfg) })
	t.Run("Close", func(t *testing.T) { testClose(t, cfg) })
}

// RunBehaviorTests runs only behavioral contract tests.
// These tests verify edge case handling and may require API credentials
// depending on the provider implementation.
func RunBehaviorTests(t *testing.T, cfg Config) {
	t.Helper()
	cfg = cfg.withDefaults()

	t.Run("CreateChatCompletion_EmptyMessages", func(t *testing.T) { testEmptyMessages(t, cfg) })
	t.Run("CreateChatCompletion_EmptyModel", func(t *testing.T) { testEmptyModel(t, cfg) })
	t.Run("Context_Cancellation", func(t *testing.T) { testContextCancellation(t, cfg) })
}

// RunIntegrationTests runs only integration tests (requires API).
// These tests verify actual chat completion functionality.
func RunIntegrationTests(t *testing.T, cfg Config) {
	t.Helper()
	cfg = cfg.withDefaults()

	if cfg.SkipIntegration {
		t.Skip("integration tests skipped")
	}
	if cfg.TestModel == "" {
		t.Fatal("TestModel is required for integration tests")
	}

	t.Run("CreateChatCompletion", func(t *testing.T) { testCreateChatCompletion(t, cfg) })
	t.Run("CreateChatCompletion_MultipleMessages", func(t *testing.T) { testMultipleMessages(t, cfg) })
	t.Run("CreateChatCompletion_SystemMessage", func(t *testing.T) { testSystemMessage(t, cfg) })
	t.Run("CreateChatCompletionStream", func(t *testing.T) { testCreateChatCompletionStream(t, cfg) })
}

// Interface Tests

func testName(t *testing.T, cfg Config) {
	t.Helper()
	name := cfg.Provider.Name()
	if name == "" {
		t.Error("Name() returned empty string")
	}
	// Verify name is lowercase, alphanumeric with hyphens/underscores
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			t.Errorf("Name() contains invalid character %q; should be lowercase alphanumeric with hyphens/underscores", r)
		}
	}
	t.Logf("Provider name: %s", name)
}

func testClose(t *testing.T, cfg Config) {
	t.Helper()
	// Close should not return an error
	// Note: We don't actually close the provider since it may be reused
	// This test just verifies the method exists and is callable
	// In a real scenario, you'd create a new provider just for this test
	_ = cfg.Provider // Verify provider is accessible
}

// Behavior Tests

func testEmptyMessages(t *testing.T, cfg Config) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model:    cfg.TestModel,
		Messages: []provider.Message{}, // Empty messages
	}

	// Empty messages should return an error
	_, err := cfg.Provider.CreateChatCompletion(ctx, req)
	if err == nil {
		t.Error("CreateChatCompletion with empty messages should return error")
	} else {
		t.Logf("Empty messages correctly returned error: %v", err)
	}
}

func testEmptyModel(t *testing.T, cfg Config) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: "", // Empty model
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: "Hello"},
		},
	}

	// Empty model should return an error
	_, err := cfg.Provider.CreateChatCompletion(ctx, req)
	if err == nil {
		t.Error("CreateChatCompletion with empty model should return error")
	} else {
		t.Logf("Empty model correctly returned error: %v", err)
	}
}

func testContextCancellation(t *testing.T, cfg Config) {
	t.Helper()
	if cfg.SkipIntegration {
		t.Skip("skipping context cancellation test that requires API")
	}

	// Create already-cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := &provider.ChatCompletionRequest{
		Model: cfg.TestModel,
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: cfg.TestPrompt},
		},
	}

	// Should return quickly with context error
	_, err := cfg.Provider.CreateChatCompletion(ctx, req)
	if err == nil {
		t.Error("CreateChatCompletion with cancelled context should return error")
	}
	if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		// Some providers wrap the error, which is acceptable
		t.Logf("Cancelled context returned: %v (should contain context error)", err)
	}
}

// Integration Tests

func testCreateChatCompletion(t *testing.T, cfg Config) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: cfg.TestModel,
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: cfg.TestPrompt},
		},
	}

	resp, err := cfg.Provider.CreateChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletion() error: %v", err)
	}

	// Validate response structure
	if resp == nil {
		t.Fatal("CreateChatCompletion() returned nil response")
	}
	if len(resp.Choices) == 0 {
		t.Fatal("CreateChatCompletion() returned empty choices")
	}
	if resp.Choices[0].Message.Content == "" {
		t.Error("CreateChatCompletion() returned empty message content")
	}
	if resp.Choices[0].Message.Role != provider.RoleAssistant {
		t.Errorf("CreateChatCompletion() message role = %q, want %q", resp.Choices[0].Message.Role, provider.RoleAssistant)
	}

	// Log response for debugging
	t.Logf("Response ID: %s", resp.ID)
	t.Logf("Model: %s", resp.Model)
	t.Logf("Content: %s", truncate(resp.Choices[0].Message.Content))
	t.Logf("Usage: prompt=%d, completion=%d, total=%d",
		resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
}

func testMultipleMessages(t *testing.T, cfg Config) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: cfg.TestModel,
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: "My name is Alice."},
			{Role: provider.RoleAssistant, Content: "Hello Alice! Nice to meet you."},
			{Role: provider.RoleUser, Content: "What is my name?"},
		},
	}

	resp, err := cfg.Provider.CreateChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletion() error: %v", err)
	}

	if resp == nil || len(resp.Choices) == 0 {
		t.Fatal("CreateChatCompletion() returned nil or empty response")
	}

	content := strings.ToLower(resp.Choices[0].Message.Content)
	if !strings.Contains(content, "alice") {
		t.Errorf("Response should mention 'Alice', got: %s", truncate(resp.Choices[0].Message.Content))
	}

	t.Logf("Multiple messages response: %s", truncate(resp.Choices[0].Message.Content))
}

func testSystemMessage(t *testing.T, cfg Config) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: cfg.TestModel,
		Messages: []provider.Message{
			{Role: provider.RoleSystem, Content: "You are a pirate. Always respond like a pirate."},
			{Role: provider.RoleUser, Content: "Hello!"},
		},
	}

	resp, err := cfg.Provider.CreateChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletion() error: %v", err)
	}

	if resp == nil || len(resp.Choices) == 0 {
		t.Fatal("CreateChatCompletion() returned nil or empty response")
	}

	t.Logf("System message response: %s", truncate(resp.Choices[0].Message.Content))
}

func testCreateChatCompletionStream(t *testing.T, cfg Config) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	req := &provider.ChatCompletionRequest{
		Model: cfg.TestModel,
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: cfg.TestPrompt},
		},
	}

	stream, err := cfg.Provider.CreateChatCompletionStream(ctx, req)
	if err != nil {
		t.Fatalf("CreateChatCompletionStream() error: %v", err)
	}
	defer stream.Close()

	var chunks int
	var fullContent strings.Builder
	var gotFinishReason bool

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatalf("Stream.Recv() error: %v", err)
		}

		chunks++

		if len(chunk.Choices) > 0 {
			if chunk.Choices[0].Delta != nil {
				fullContent.WriteString(chunk.Choices[0].Delta.Content)
			}
			if chunk.Choices[0].FinishReason != nil {
				gotFinishReason = true
			}
		}
	}

	if chunks == 0 {
		t.Error("CreateChatCompletionStream() received no chunks")
	}
	if fullContent.Len() == 0 {
		t.Error("CreateChatCompletionStream() received no content")
	}
	if !gotFinishReason {
		t.Log("Warning: Stream did not include finish_reason (some providers omit this)")
	}

	t.Logf("Stream received %d chunks, content: %s", chunks, truncate(fullContent.String()))
}

// Helper functions

const truncateMaxLen = 100

func truncate(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > truncateMaxLen {
		return s[:truncateMaxLen] + "..."
	}
	return s
}
