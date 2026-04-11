package omnillm

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/grokify/mogo/log/slogutil"

	"github.com/plexusone/omnillm/provider"
)

// ProviderConfig holds configuration for a single provider instance.
// Used in the Providers slice where index 0 is primary and 1+ are fallbacks.
type ProviderConfig struct {
	// Provider is the provider type (e.g., ProviderNameOpenAI).
	// Ignored if CustomProvider is set.
	Provider ProviderName

	// APIKey is the API key for the provider
	APIKey string //nolint:gosec // G117: config field for API key, not a hardcoded credential

	// BaseURL is an optional custom base URL
	BaseURL string

	// Region is for providers that require a region (e.g., AWS Bedrock)
	Region string

	// Timeout sets the HTTP client timeout for this provider
	Timeout time.Duration

	// HTTPClient is an optional custom HTTP client
	HTTPClient *http.Client

	// Extra holds provider-specific configuration
	Extra map[string]any

	// CustomProvider allows injecting a custom provider implementation.
	// When set, Provider, APIKey, BaseURL, etc. are ignored.
	CustomProvider provider.Provider
}

// FallbackProvider wraps multiple providers with fallback logic.
// It implements provider.Provider and tries providers in order until one succeeds.
type FallbackProvider struct {
	primary         provider.Provider
	fallbacks       []provider.Provider
	circuitBreakers map[string]*CircuitBreaker
	cbConfig        *CircuitBreakerConfig
	logger          *slog.Logger
}

// FallbackProviderConfig configures the fallback provider behavior
type FallbackProviderConfig struct {
	// CircuitBreakerConfig configures circuit breaker behavior.
	// If nil, circuit breaker is disabled.
	CircuitBreakerConfig *CircuitBreakerConfig

	// Logger for logging fallback events
	Logger *slog.Logger
}

// NewFallbackProvider creates a provider that tries fallbacks on failure.
// The primary provider is tried first, then fallbacks in order.
func NewFallbackProvider(
	primary provider.Provider,
	fallbacks []provider.Provider,
	config *FallbackProviderConfig,
) *FallbackProvider {
	if config == nil {
		config = &FallbackProviderConfig{}
	}

	fp := &FallbackProvider{
		primary:   primary,
		fallbacks: fallbacks,
		cbConfig:  config.CircuitBreakerConfig,
		logger:    config.Logger,
	}

	if fp.logger == nil {
		fp.logger = slogutil.Null()
	}

	// Initialize circuit breakers if configured
	if config.CircuitBreakerConfig != nil {
		fp.circuitBreakers = make(map[string]*CircuitBreaker)
		fp.circuitBreakers[primary.Name()] = NewCircuitBreaker(*config.CircuitBreakerConfig)
		for _, fb := range fallbacks {
			fp.circuitBreakers[fb.Name()] = NewCircuitBreaker(*config.CircuitBreakerConfig)
		}
	}

	return fp
}

// CreateChatCompletion tries the primary provider first, then fallbacks on retryable errors.
func (fp *FallbackProvider) CreateChatCompletion(
	ctx context.Context,
	req *provider.ChatCompletionRequest,
) (*provider.ChatCompletionResponse, error) {
	attempts := make([]FallbackAttempt, 0, 1+len(fp.fallbacks))

	// Try primary first
	resp, err := fp.tryProvider(ctx, fp.primary, req, &attempts)
	if err == nil {
		return resp, nil
	}

	// Don't fallback for non-retryable errors
	if IsNonRetryableError(err) {
		fp.logger.Debug("non-retryable error from primary, not attempting fallback",
			slog.String("provider", fp.primary.Name()),
			slog.String("error", err.Error()))
		return nil, err
	}

	// Try fallbacks in order
	for _, fb := range fp.fallbacks {
		resp, err = fp.tryProvider(ctx, fb, req, &attempts)
		if err == nil {
			return resp, nil
		}

		// Stop on non-retryable errors
		if IsNonRetryableError(err) {
			fp.logger.Debug("non-retryable error from fallback, stopping",
				slog.String("provider", fb.Name()),
				slog.String("error", err.Error()))
			break
		}
	}

	// All providers failed
	return nil, &FallbackError{
		Attempts:  attempts,
		LastError: err,
	}
}

// CreateChatCompletionStream tries the primary provider first, then fallbacks on retryable errors.
func (fp *FallbackProvider) CreateChatCompletionStream(
	ctx context.Context,
	req *provider.ChatCompletionRequest,
) (provider.ChatCompletionStream, error) {
	attempts := make([]FallbackAttempt, 0, 1+len(fp.fallbacks))

	// Try primary first
	stream, err := fp.tryProviderStream(ctx, fp.primary, req, &attempts)
	if err == nil {
		return stream, nil
	}

	// Don't fallback for non-retryable errors
	if IsNonRetryableError(err) {
		fp.logger.Debug("non-retryable error from primary, not attempting fallback",
			slog.String("provider", fp.primary.Name()),
			slog.String("error", err.Error()))
		return nil, err
	}

	// Try fallbacks in order
	for _, fb := range fp.fallbacks {
		stream, err = fp.tryProviderStream(ctx, fb, req, &attempts)
		if err == nil {
			return stream, nil
		}

		// Stop on non-retryable errors
		if IsNonRetryableError(err) {
			fp.logger.Debug("non-retryable error from fallback, stopping",
				slog.String("provider", fb.Name()),
				slog.String("error", err.Error()))
			break
		}
	}

	// All providers failed
	return nil, &FallbackError{
		Attempts:  attempts,
		LastError: err,
	}
}

// Close closes all providers
func (fp *FallbackProvider) Close() error {
	var lastErr error

	if err := fp.primary.Close(); err != nil {
		lastErr = err
	}

	for _, fb := range fp.fallbacks {
		if err := fb.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// Name returns a composite name indicating fallback configuration
func (fp *FallbackProvider) Name() string {
	return fp.primary.Name() + "+fallback"
}

// PrimaryProvider returns the primary provider
func (fp *FallbackProvider) PrimaryProvider() provider.Provider {
	return fp.primary
}

// FallbackProviders returns the fallback providers
func (fp *FallbackProvider) FallbackProviders() []provider.Provider {
	return fp.fallbacks
}

// CircuitBreaker returns the circuit breaker for a provider, or nil if not configured
func (fp *FallbackProvider) CircuitBreaker(providerName string) *CircuitBreaker {
	if fp.circuitBreakers == nil {
		return nil
	}
	return fp.circuitBreakers[providerName]
}

// shouldTryProvider checks if the provider should be tried based on circuit breaker state
func (fp *FallbackProvider) shouldTryProvider(providerName string) bool {
	if fp.circuitBreakers == nil {
		return true
	}

	cb, ok := fp.circuitBreakers[providerName]
	if !ok {
		return true
	}

	return cb.AllowRequest()
}

// recordSuccess records a successful request for the circuit breaker
func (fp *FallbackProvider) recordSuccess(providerName string) {
	if fp.circuitBreakers == nil {
		return
	}

	if cb, ok := fp.circuitBreakers[providerName]; ok {
		cb.RecordSuccess()
	}
}

// recordFailure records a failed request for the circuit breaker
func (fp *FallbackProvider) recordFailure(providerName string, err error) {
	if fp.circuitBreakers == nil {
		return
	}

	// Only record retryable errors as failures for circuit breaker
	if !IsRetryableError(err) {
		return
	}

	if cb, ok := fp.circuitBreakers[providerName]; ok {
		cb.RecordFailure()
	}
}

// tryProvider attempts a request to a single provider
func (fp *FallbackProvider) tryProvider(
	ctx context.Context,
	p provider.Provider,
	req *provider.ChatCompletionRequest,
	attempts *[]FallbackAttempt,
) (*provider.ChatCompletionResponse, error) {
	providerName := p.Name()
	start := time.Now()

	// Check circuit breaker
	if !fp.shouldTryProvider(providerName) {
		cb := fp.circuitBreakers[providerName]
		err := &CircuitOpenError{
			Provider:    providerName,
			State:       cb.State(),
			LastFailure: cb.Stats().LastFailure,
			RetryAfter:  fp.cbConfig.Timeout - time.Since(cb.Stats().LastFailure),
		}
		*attempts = append(*attempts, FallbackAttempt{
			Provider: providerName,
			Error:    err,
			Duration: time.Since(start),
			Skipped:  true,
		})
		fp.logger.Debug("skipping provider due to open circuit",
			slog.String("provider", providerName))
		return nil, err
	}

	// Try the provider
	resp, err := p.CreateChatCompletion(ctx, req)
	duration := time.Since(start)

	*attempts = append(*attempts, FallbackAttempt{
		Provider: providerName,
		Error:    err,
		Duration: duration,
	})

	if err != nil {
		fp.recordFailure(providerName, err)
		fp.logger.Debug("provider request failed",
			slog.String("provider", providerName),
			slog.Duration("duration", duration),
			slog.String("error", err.Error()))
		return nil, err
	}

	fp.recordSuccess(providerName)
	fp.logger.Debug("provider request succeeded",
		slog.String("provider", providerName),
		slog.Duration("duration", duration))

	// Add metadata about which provider was used
	if resp.ProviderMetadata == nil {
		resp.ProviderMetadata = make(map[string]any)
	}
	resp.ProviderMetadata["fallback_provider_used"] = providerName
	resp.ProviderMetadata["fallback_attempt_count"] = len(*attempts)

	return resp, nil
}

// tryProviderStream attempts a streaming request to a single provider
func (fp *FallbackProvider) tryProviderStream(
	ctx context.Context,
	p provider.Provider,
	req *provider.ChatCompletionRequest,
	attempts *[]FallbackAttempt,
) (provider.ChatCompletionStream, error) {
	providerName := p.Name()
	start := time.Now()

	// Check circuit breaker
	if !fp.shouldTryProvider(providerName) {
		cb := fp.circuitBreakers[providerName]
		err := &CircuitOpenError{
			Provider:    providerName,
			State:       cb.State(),
			LastFailure: cb.Stats().LastFailure,
			RetryAfter:  fp.cbConfig.Timeout - time.Since(cb.Stats().LastFailure),
		}
		*attempts = append(*attempts, FallbackAttempt{
			Provider: providerName,
			Error:    err,
			Duration: time.Since(start),
			Skipped:  true,
		})
		fp.logger.Debug("skipping provider due to open circuit",
			slog.String("provider", providerName))
		return nil, err
	}

	// Try the provider
	stream, err := p.CreateChatCompletionStream(ctx, req)
	duration := time.Since(start)

	*attempts = append(*attempts, FallbackAttempt{
		Provider: providerName,
		Error:    err,
		Duration: duration,
	})

	if err != nil {
		fp.recordFailure(providerName, err)
		fp.logger.Debug("provider stream request failed",
			slog.String("provider", providerName),
			slog.Duration("duration", duration),
			slog.String("error", err.Error()))
		return nil, err
	}

	fp.recordSuccess(providerName)
	fp.logger.Debug("provider stream request succeeded",
		slog.String("provider", providerName),
		slog.Duration("duration", duration))

	// Wrap stream to record circuit breaker success on completion
	return &fallbackAwareStream{
		stream:       stream,
		fp:           fp,
		providerName: providerName,
	}, nil
}

// fallbackAwareStream wraps a stream to track circuit breaker state
type fallbackAwareStream struct {
	stream       provider.ChatCompletionStream
	fp           *FallbackProvider
	providerName string
	closed       bool
}

func (s *fallbackAwareStream) Recv() (*provider.ChatCompletionChunk, error) {
	chunk, err := s.stream.Recv()
	if err != nil && err.Error() != "EOF" {
		// Record failure on non-EOF errors
		s.fp.recordFailure(s.providerName, err)
	}
	return chunk, err
}

func (s *fallbackAwareStream) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	return s.stream.Close()
}

// FallbackAttempt records information about a single fallback attempt
type FallbackAttempt struct {
	// Provider is the name of the provider that was tried
	Provider string

	// Error is the error returned, or nil on success
	Error error

	// Duration is how long the attempt took
	Duration time.Duration

	// Skipped indicates the provider was skipped (e.g., circuit open)
	Skipped bool
}

// FallbackError is returned when all providers fail
type FallbackError struct {
	// Attempts contains information about each provider attempt
	Attempts []FallbackAttempt

	// LastError is the last error encountered
	LastError error
}

func (e *FallbackError) Error() string {
	if len(e.Attempts) == 0 {
		return "all providers failed"
	}
	return fmt.Sprintf("all %d providers failed, last error: %v", len(e.Attempts), e.LastError)
}

func (e *FallbackError) Unwrap() error {
	return e.LastError
}

// buildProviderFromConfig creates a provider from a ProviderConfig
func buildProviderFromConfig(config ProviderConfig) (provider.Provider, error) {
	// Check for custom provider injection first
	if config.CustomProvider != nil {
		return config.CustomProvider, nil
	}

	// Special case: Bedrock requires external module
	if config.Provider == ProviderNameBedrock {
		return nil, ErrBedrockExternal
	}

	// Look up provider in registry
	factory := GetProviderFactory(config.Provider)
	if factory == nil {
		return nil, ErrUnsupportedProvider
	}

	return factory(config)
}
