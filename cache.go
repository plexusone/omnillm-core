package omnillm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grokify/sogo/database/kvs"

	"github.com/plexusone/omnillm-core/provider"
)

// CacheConfig configures response caching behavior
type CacheConfig struct {
	// TTL is the time-to-live for cached responses.
	// Default: 1 hour
	TTL time.Duration

	// KeyPrefix is the prefix for cache keys in the KVS.
	// Default: "omnillm:cache"
	KeyPrefix string

	// SkipStreaming skips caching for streaming requests.
	// Default: true (streaming responses are not cached)
	SkipStreaming bool

	// CacheableModels limits caching to specific models.
	// If nil or empty, all models are cached.
	CacheableModels []string

	// ExcludeParameters lists parameters to exclude from cache key calculation.
	// Common exclusions: "user" (user ID shouldn't affect cache)
	// Default: ["user"]
	ExcludeParameters []string

	// IncludeTemperature includes temperature in cache key.
	// Set to false if you want to cache regardless of temperature setting.
	// Default: true
	IncludeTemperature bool

	// IncludeSeed includes seed in cache key.
	// Default: true
	IncludeSeed bool
}

// DefaultCacheConfig returns a CacheConfig with sensible defaults
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		TTL:                1 * time.Hour,
		KeyPrefix:          "omnillm:cache",
		SkipStreaming:      true,
		ExcludeParameters:  []string{"user"},
		IncludeTemperature: true,
		IncludeSeed:        true,
	}
}

// CacheEntry represents a cached response with metadata
type CacheEntry struct {
	// Response is the cached chat completion response
	Response *provider.ChatCompletionResponse `json:"response"`

	// CachedAt is when the response was cached
	CachedAt time.Time `json:"cached_at"`

	// ExpiresAt is when the cache entry expires
	ExpiresAt time.Time `json:"expires_at"`

	// Model is the model used for the request
	Model string `json:"model"`

	// RequestHash is the hash of the request (for verification)
	RequestHash string `json:"request_hash"`
}

// IsExpired returns true if the cache entry has expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// CacheManager handles response caching using a KVS backend
type CacheManager struct {
	kvs    kvs.Client
	config CacheConfig
}

// NewCacheManager creates a new cache manager with the given KVS client and configuration.
// If config has zero values, defaults are used for those fields.
func NewCacheManager(kvsClient kvs.Client, config CacheConfig) *CacheManager {
	if config.TTL == 0 {
		config.TTL = 1 * time.Hour
	}
	if config.KeyPrefix == "" {
		config.KeyPrefix = "omnillm:cache"
	}
	if config.ExcludeParameters == nil {
		config.ExcludeParameters = []string{"user"}
	}

	return &CacheManager{
		kvs:    kvsClient,
		config: config,
	}
}

// Get retrieves a cached response for the given request.
// Returns nil if no valid cache entry exists.
func (m *CacheManager) Get(ctx context.Context, req *provider.ChatCompletionRequest) (*CacheEntry, error) {
	key := m.BuildCacheKey(req)

	var entry CacheEntry
	if err := m.kvs.GetAny(ctx, key, &entry); err != nil {
		// Cache miss or error
		return nil, nil
	}

	// Check expiration
	if entry.IsExpired() {
		return nil, nil
	}

	return &entry, nil
}

// Set stores a response in the cache for the given request.
func (m *CacheManager) Set(ctx context.Context, req *provider.ChatCompletionRequest, resp *provider.ChatCompletionResponse) error {
	key := m.BuildCacheKey(req)
	now := time.Now()

	entry := CacheEntry{
		Response:    resp,
		CachedAt:    now,
		ExpiresAt:   now.Add(m.config.TTL),
		Model:       req.Model,
		RequestHash: m.hashRequest(req),
	}

	return m.kvs.SetAny(ctx, key, entry)
}

// Delete removes a cache entry for the given request.
func (m *CacheManager) Delete(ctx context.Context, req *provider.ChatCompletionRequest) error {
	key := m.BuildCacheKey(req)
	return m.kvs.SetString(ctx, key, "") // KVS doesn't have Delete, use empty string
}

// ShouldCache determines if a request should be cached.
// Returns false for streaming requests (if configured), non-cacheable models, etc.
func (m *CacheManager) ShouldCache(req *provider.ChatCompletionRequest) bool {
	// Skip streaming requests if configured
	if m.config.SkipStreaming && req.Stream != nil && *req.Stream {
		return false
	}

	// Check model allowlist if configured
	if len(m.config.CacheableModels) > 0 {
		for _, model := range m.config.CacheableModels {
			if req.Model == model {
				return true
			}
		}
		return false
	}

	return true
}

// BuildCacheKey generates a deterministic cache key for a request.
// The key is a hash of the normalized request parameters.
func (m *CacheManager) BuildCacheKey(req *provider.ChatCompletionRequest) string {
	hash := m.hashRequest(req)
	return fmt.Sprintf("%s:%s", m.config.KeyPrefix, hash)
}

// normalizedRequest is used for cache key generation
type normalizedRequest struct {
	Model       string              `json:"model"`
	Messages    []normalizedMessage `json:"messages"`
	MaxTokens   *int                `json:"max_tokens,omitempty"`
	Temperature *float64            `json:"temperature,omitempty"`
	TopP        *float64            `json:"top_p,omitempty"`
	TopK        *int                `json:"top_k,omitempty"`
	Seed        *int                `json:"seed,omitempty"`
	Stop        []string            `json:"stop,omitempty"`
}

type normalizedMessage struct {
	Role       string  `json:"role"`
	Content    string  `json:"content"`
	Name       *string `json:"name,omitempty"`
	ToolCallID *string `json:"tool_call_id,omitempty"`
}

// hashRequest creates a deterministic hash of the request for caching
func (m *CacheManager) hashRequest(req *provider.ChatCompletionRequest) string {
	normalized := normalizedRequest{
		Model: req.Model,
	}

	// Normalize messages
	for _, msg := range req.Messages {
		normalized.Messages = append(normalized.Messages, normalizedMessage{
			Role:       string(msg.Role),
			Content:    msg.Content,
			Name:       msg.Name,
			ToolCallID: msg.ToolCallID,
		})
	}

	// Include parameters that affect output
	if req.MaxTokens != nil {
		normalized.MaxTokens = req.MaxTokens
	}

	if m.config.IncludeTemperature && req.Temperature != nil {
		normalized.Temperature = req.Temperature
	}

	if req.TopP != nil {
		normalized.TopP = req.TopP
	}

	if req.TopK != nil {
		normalized.TopK = req.TopK
	}

	if m.config.IncludeSeed && req.Seed != nil {
		normalized.Seed = req.Seed
	}

	if len(req.Stop) > 0 {
		normalized.Stop = req.Stop
	}

	// Hash the normalized request
	data, _ := json.Marshal(normalized)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:16]) // Use first 16 bytes for shorter keys
}

// Config returns the cache configuration
func (m *CacheManager) Config() CacheConfig {
	return m.config
}

// CacheStats contains statistics about cache usage
type CacheStats struct {
	Hits   int64
	Misses int64
}

// CacheHitError is a marker type to indicate a cache hit (not an actual error)
type CacheHitError struct {
	Entry *CacheEntry
}

func (e *CacheHitError) Error() string {
	return "cache hit"
}
