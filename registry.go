package omnillm

import (
	"sync"

	"github.com/plexusone/omnillm-core/provider"
)

// ProviderFactory is a function that creates a provider from config.
type ProviderFactory func(config ProviderConfig) (provider.Provider, error)

// registeredProvider holds a factory with its priority.
type registeredProvider struct {
	factory  ProviderFactory
	priority int
}

var (
	providerRegistry = make(map[ProviderName]registeredProvider)
	registryMu       sync.RWMutex
)

// RegisterProvider registers a provider factory with the given name and priority.
// Higher priority values override lower priority registrations.
// Thin (stdlib) providers should use priority 0.
// Thick (SDK) providers should use priority 10.
//
// Example:
//
//	// In omnillm-core/providers/openai/init.go (thin, priority 0)
//	func init() {
//	    omnillm.RegisterProvider(omnillm.ProviderNameOpenAI, NewProvider, 0)
//	}
//
//	// In omnillm-openai/init.go (thick, priority 10)
//	func init() {
//	    omnillm.RegisterProvider(omnillm.ProviderNameOpenAI, NewProvider, 10)
//	}
func RegisterProvider(name ProviderName, factory ProviderFactory, priority int) {
	registryMu.Lock()
	defer registryMu.Unlock()

	existing, ok := providerRegistry[name]
	if !ok || priority >= existing.priority {
		providerRegistry[name] = registeredProvider{
			factory:  factory,
			priority: priority,
		}
	}
}

// GetProviderFactory returns the registered factory for the given provider name.
// Returns nil if no provider is registered with that name.
func GetProviderFactory(name ProviderName) ProviderFactory {
	registryMu.RLock()
	defer registryMu.RUnlock()

	if rp, ok := providerRegistry[name]; ok {
		return rp.factory
	}
	return nil
}

// ListRegisteredProviders returns a list of all registered provider names.
func ListRegisteredProviders() []ProviderName {
	registryMu.RLock()
	defer registryMu.RUnlock()

	names := make([]ProviderName, 0, len(providerRegistry))
	for name := range providerRegistry {
		names = append(names, name)
	}
	return names
}

// GetProviderPriority returns the priority of the registered provider.
// Returns -1 if the provider is not registered.
func GetProviderPriority(name ProviderName) int {
	registryMu.RLock()
	defer registryMu.RUnlock()

	if rp, ok := providerRegistry[name]; ok {
		return rp.priority
	}
	return -1
}

// Priority constants for provider registration.
const (
	// PriorityThin is the priority for thin (stdlib-only) provider implementations.
	PriorityThin = 0

	// PriorityThick is the priority for thick (official SDK) provider implementations.
	PriorityThick = 10
)

// init registers the built-in thin providers.
func init() {
	RegisterProvider(ProviderNameOpenAI, newOpenAIProvider, PriorityThin)
	RegisterProvider(ProviderNameAnthropic, newAnthropicProvider, PriorityThin)
	RegisterProvider(ProviderNameOllama, newOllamaProvider, PriorityThin)
	RegisterProvider(ProviderNameXAI, newXAIProvider, PriorityThin)
	RegisterProvider(ProviderNameKimi, newKimiProvider, PriorityThin)
	RegisterProvider(ProviderNameGLM, newGLMProvider, PriorityThin)
	RegisterProvider(ProviderNameQwen, newQwenProvider, PriorityThin)

	// Register embedding providers
	RegisterEmbeddingProvider(ProviderNameOpenAI, newOpenAIEmbeddingProvider, PriorityThin)
}

// EmbeddingProviderFactory is a function that creates an embedding provider from config.
type EmbeddingProviderFactory func(config ProviderConfig) (provider.EmbeddingProvider, error)

// registeredEmbeddingProvider holds a factory with its priority.
type registeredEmbeddingProvider struct {
	factory  EmbeddingProviderFactory
	priority int
}

var (
	embeddingProviderRegistry = make(map[ProviderName]registeredEmbeddingProvider)
	embeddingRegistryMu       sync.RWMutex
)

// RegisterEmbeddingProvider registers an embedding provider factory with the given name and priority.
// Higher priority values override lower priority registrations.
// Thin (stdlib) providers should use priority 0.
// Thick (SDK) providers should use priority 10.
func RegisterEmbeddingProvider(name ProviderName, factory EmbeddingProviderFactory, priority int) {
	embeddingRegistryMu.Lock()
	defer embeddingRegistryMu.Unlock()

	existing, ok := embeddingProviderRegistry[name]
	if !ok || priority >= existing.priority {
		embeddingProviderRegistry[name] = registeredEmbeddingProvider{
			factory:  factory,
			priority: priority,
		}
	}
}

// GetEmbeddingProviderFactory returns the registered factory for the given provider name.
// Returns nil if no embedding provider is registered with that name.
func GetEmbeddingProviderFactory(name ProviderName) EmbeddingProviderFactory {
	embeddingRegistryMu.RLock()
	defer embeddingRegistryMu.RUnlock()

	if rp, ok := embeddingProviderRegistry[name]; ok {
		return rp.factory
	}
	return nil
}

// ListEmbeddingProviders returns a list of all registered embedding provider names.
func ListEmbeddingProviders() []ProviderName {
	embeddingRegistryMu.RLock()
	defer embeddingRegistryMu.RUnlock()

	names := make([]ProviderName, 0, len(embeddingProviderRegistry))
	for name := range embeddingProviderRegistry {
		names = append(names, name)
	}
	return names
}

// GetEmbeddingProvider creates an embedding provider instance from the registry.
// Returns an error if the provider is not registered or if creation fails.
func GetEmbeddingProvider(name ProviderName, config ProviderConfig) (provider.EmbeddingProvider, error) {
	factory := GetEmbeddingProviderFactory(name)
	if factory == nil {
		return nil, &APIError{
			Provider: name,
			Type:     "provider_not_found",
			Message:  "embedding provider not registered: " + string(name),
		}
	}
	return factory(config)
}
