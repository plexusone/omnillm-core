package omnillm

import (
	"context"
	"testing"
	"time"

	"github.com/plexusone/omnillm-core/provider"
	testutil "github.com/plexusone/omnillm-core/testing"
)

func TestCacheManager_SetAndGet(t *testing.T) {
	kvs := testutil.NewMockKVS()
	cache := NewCacheManager(kvs, DefaultCacheConfig())
	ctx := context.Background()

	req := &provider.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []provider.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	resp := &provider.ChatCompletionResponse{
		ID:    "resp-123",
		Model: "gpt-4o",
		Choices: []provider.ChatCompletionChoice{
			{
				Index: 0,
				Message: provider.Message{
					Role:    "assistant",
					Content: "Hi there!",
				},
			},
		},
		Usage: provider.Usage{
			PromptTokens:     5,
			CompletionTokens: 3,
			TotalTokens:      8,
		},
	}

	// Set cache
	err := cache.Set(ctx, req, resp)
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	// Get cache
	entry, err := cache.Get(ctx, req)
	if err != nil {
		t.Fatalf("failed to get cache: %v", err)
	}

	if entry == nil {
		t.Fatal("expected cache entry, got nil")
	}

	if entry.Response.ID != "resp-123" {
		t.Errorf("expected response ID 'resp-123', got %q", entry.Response.ID)
	}

	if entry.Model != "gpt-4o" {
		t.Errorf("expected model 'gpt-4o', got %q", entry.Model)
	}
}

func TestCacheManager_CacheMiss(t *testing.T) {
	kvs := testutil.NewMockKVS()
	cache := NewCacheManager(kvs, DefaultCacheConfig())
	ctx := context.Background()

	req := &provider.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []provider.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	// Get non-existent entry
	entry, err := cache.Get(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entry != nil {
		t.Error("expected nil entry for cache miss")
	}
}

func TestCacheManager_Expiration(t *testing.T) {
	kvs := testutil.NewMockKVS()
	config := CacheConfig{
		TTL:       50 * time.Millisecond,
		KeyPrefix: "test:cache",
	}
	cache := NewCacheManager(kvs, config)
	ctx := context.Background()

	req := &provider.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []provider.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	resp := &provider.ChatCompletionResponse{
		ID:    "resp-123",
		Model: "gpt-4o",
	}

	// Set cache
	err := cache.Set(ctx, req, resp)
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	// Verify it's there
	entry, _ := cache.Get(ctx, req)
	if entry == nil {
		t.Fatal("expected cache entry immediately after set")
	}

	// Wait for expiration
	time.Sleep(60 * time.Millisecond)

	// Should be expired
	entry, _ = cache.Get(ctx, req)
	if entry != nil {
		t.Error("expected nil entry after expiration")
	}
}

func TestCacheManager_ShouldCache(t *testing.T) {
	cache := NewCacheManager(testutil.NewMockKVS(), DefaultCacheConfig())

	tests := []struct {
		name     string
		req      *provider.ChatCompletionRequest
		expected bool
	}{
		{
			name: "regular request",
			req: &provider.ChatCompletionRequest{
				Model:    "gpt-4o",
				Messages: []provider.Message{{Role: "user", Content: "Hello"}},
			},
			expected: true,
		},
		{
			name: "streaming request",
			req: &provider.ChatCompletionRequest{
				Model:    "gpt-4o",
				Messages: []provider.Message{{Role: "user", Content: "Hello"}},
				Stream:   boolPtr(true),
			},
			expected: false, // Streaming is not cached by default
		},
		{
			name: "non-streaming explicit",
			req: &provider.ChatCompletionRequest{
				Model:    "gpt-4o",
				Messages: []provider.Message{{Role: "user", Content: "Hello"}},
				Stream:   boolPtr(false),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cache.ShouldCache(tt.req)
			if result != tt.expected {
				t.Errorf("ShouldCache() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCacheManager_ShouldCache_ModelAllowlist(t *testing.T) {
	config := CacheConfig{
		CacheableModels: []string{"gpt-4o", "claude-3-opus"},
	}
	cache := NewCacheManager(testutil.NewMockKVS(), config)

	tests := []struct {
		model    string
		expected bool
	}{
		{"gpt-4o", true},
		{"claude-3-opus", true},
		{"gpt-3.5-turbo", false},
		{"other-model", false},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			req := &provider.ChatCompletionRequest{
				Model:    tt.model,
				Messages: []provider.Message{{Role: "user", Content: "Hello"}},
			}
			result := cache.ShouldCache(req)
			if result != tt.expected {
				t.Errorf("ShouldCache(%s) = %v, expected %v", tt.model, result, tt.expected)
			}
		})
	}
}

func TestCacheManager_BuildCacheKey(t *testing.T) {
	cache := NewCacheManager(testutil.NewMockKVS(), DefaultCacheConfig())

	req1 := &provider.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []provider.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	req2 := &provider.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []provider.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	req3 := &provider.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []provider.Message{
			{Role: "user", Content: "Different message"},
		},
	}

	key1 := cache.BuildCacheKey(req1)
	key2 := cache.BuildCacheKey(req2)
	key3 := cache.BuildCacheKey(req3)

	// Same requests should have same key
	if key1 != key2 {
		t.Errorf("identical requests should have same key: %s != %s", key1, key2)
	}

	// Different requests should have different keys
	if key1 == key3 {
		t.Errorf("different requests should have different keys: %s == %s", key1, key3)
	}

	// Key should have prefix
	if len(key1) < len("omnillm:cache:") {
		t.Errorf("key should have prefix, got %s", key1)
	}
}

func TestCacheManager_KeyIncludesParameters(t *testing.T) {
	cache := NewCacheManager(testutil.NewMockKVS(), DefaultCacheConfig())

	baseReq := &provider.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []provider.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	temp := 0.7
	reqWithTemp := &provider.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []provider.Message{
			{Role: "user", Content: "Hello"},
		},
		Temperature: &temp,
	}

	key1 := cache.BuildCacheKey(baseReq)
	key2 := cache.BuildCacheKey(reqWithTemp)

	// Different temperature should produce different key
	if key1 == key2 {
		t.Error("requests with different temperature should have different keys")
	}
}

func TestCacheManager_KeyExcludesTemperatureWhenConfigured(t *testing.T) {
	config := CacheConfig{
		IncludeTemperature: false,
	}
	cache := NewCacheManager(testutil.NewMockKVS(), config)

	temp1 := 0.5
	temp2 := 1.0

	req1 := &provider.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []provider.Message{
			{Role: "user", Content: "Hello"},
		},
		Temperature: &temp1,
	}

	req2 := &provider.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []provider.Message{
			{Role: "user", Content: "Hello"},
		},
		Temperature: &temp2,
	}

	key1 := cache.BuildCacheKey(req1)
	key2 := cache.BuildCacheKey(req2)

	// With IncludeTemperature=false, different temperatures should have same key
	if key1 != key2 {
		t.Error("requests should have same key when temperature is excluded")
	}
}

func TestCacheEntry_IsExpired(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "not expired",
			expiresAt: now.Add(1 * time.Hour),
			expected:  false,
		},
		{
			name:      "expired",
			expiresAt: now.Add(-1 * time.Hour),
			expected:  true,
		},
		{
			name:      "just expired",
			expiresAt: now.Add(-1 * time.Millisecond),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := &CacheEntry{ExpiresAt: tt.expiresAt}
			if entry.IsExpired() != tt.expected {
				t.Errorf("IsExpired() = %v, expected %v", entry.IsExpired(), tt.expected)
			}
		})
	}
}

func TestDefaultCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()

	if config.TTL != 1*time.Hour {
		t.Errorf("expected TTL=1h, got %v", config.TTL)
	}

	if config.KeyPrefix != "omnillm:cache" {
		t.Errorf("expected KeyPrefix='omnillm:cache', got %q", config.KeyPrefix)
	}

	if !config.SkipStreaming {
		t.Error("expected SkipStreaming=true")
	}

	if !config.IncludeTemperature {
		t.Error("expected IncludeTemperature=true")
	}

	if !config.IncludeSeed {
		t.Error("expected IncludeSeed=true")
	}
}

func boolPtr(b bool) *bool {
	return &b
}
