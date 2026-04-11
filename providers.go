package omnillm

import (
	"net/http"

	"github.com/plexusone/omnillm/provider"
	"github.com/plexusone/omnillm/providers/anthropic"
	"github.com/plexusone/omnillm/providers/glm"
	"github.com/plexusone/omnillm/providers/kimi"
	"github.com/plexusone/omnillm/providers/ollama"
	"github.com/plexusone/omnillm/providers/openai"
	"github.com/plexusone/omnillm/providers/qwen"
	"github.com/plexusone/omnillm/providers/xai"
)

// getHTTPClientFromProviderConfig returns the HTTPClient from config, or creates one with the
// configured Timeout. Returns nil if neither is set (provider will use defaults).
func getHTTPClientFromProviderConfig(config ProviderConfig) *http.Client {
	if config.HTTPClient != nil {
		return config.HTTPClient
	}
	if config.Timeout > 0 {
		return &http.Client{Timeout: config.Timeout}
	}
	return nil
}

// newOpenAIProvider creates a new OpenAI provider adapter
func newOpenAIProvider(config ProviderConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return openai.NewProvider(config.APIKey, config.BaseURL, getHTTPClientFromProviderConfig(config)), nil
}

// newAnthropicProvider creates a new Anthropic provider adapter
func newAnthropicProvider(config ProviderConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return anthropic.NewProvider(config.APIKey, config.BaseURL, getHTTPClientFromProviderConfig(config)), nil
}

// newOllamaProvider creates a new Ollama provider adapter
func newOllamaProvider(config ProviderConfig) (provider.Provider, error) { //nolint:unparam // `error` added to fulfill interface requirements
	return ollama.NewProvider(config.BaseURL, getHTTPClientFromProviderConfig(config)), nil
}

// newXAIProvider creates a new X.AI provider adapter
func newXAIProvider(config ProviderConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return xai.NewProvider(config.APIKey, config.BaseURL, getHTTPClientFromProviderConfig(config)), nil
}

// newKimiProvider creates a new Kimi (Moonshot AI) provider adapter
func newKimiProvider(config ProviderConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return kimi.NewProvider(config.APIKey, config.BaseURL, getHTTPClientFromProviderConfig(config)), nil
}

// newGLMProvider creates a new GLM (Zhipu AI) provider adapter
func newGLMProvider(config ProviderConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return glm.NewProvider(config.APIKey, config.BaseURL, getHTTPClientFromProviderConfig(config)), nil
}

// newQwenProvider creates a new Qwen (Alibaba Cloud) provider adapter
func newQwenProvider(config ProviderConfig) (provider.Provider, error) {
	if config.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}
	return qwen.NewProvider(config.APIKey, config.BaseURL, getHTTPClientFromProviderConfig(config)), nil
}
