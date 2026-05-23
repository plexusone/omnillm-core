package omnillm

import (
	"testing"

	"github.com/plexusone/omnillm-core/provider"
)

func TestListEmbeddingProviders(t *testing.T) {
	providers := ListEmbeddingProviders()

	// OpenAI should be registered by default
	found := false
	for _, p := range providers {
		if p == ProviderNameOpenAI {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected OpenAI to be registered as an embedding provider")
	}
}

func TestGetEmbeddingProviderFactory(t *testing.T) {
	// Test existing provider
	factory := GetEmbeddingProviderFactory(ProviderNameOpenAI)
	if factory == nil {
		t.Error("Expected factory for OpenAI, got nil")
	}

	// Test non-existent provider
	factory = GetEmbeddingProviderFactory("nonexistent")
	if factory != nil {
		t.Error("Expected nil for non-existent provider")
	}
}

func TestGetEmbeddingProvider(t *testing.T) {
	// Test successful creation
	p, err := GetEmbeddingProvider(ProviderNameOpenAI, ProviderConfig{
		APIKey: "test-key",
	})
	if err != nil {
		t.Fatalf("GetEmbeddingProvider failed: %v", err)
	}
	if p == nil {
		t.Error("Expected provider, got nil")
	}
	if p.Name() != "openai" {
		t.Errorf("Expected name 'openai', got %s", p.Name())
	}
	p.Close()

	// Test missing API key
	_, err = GetEmbeddingProvider(ProviderNameOpenAI, ProviderConfig{})
	if err == nil {
		t.Error("Expected error for missing API key, got nil")
	}

	// Test non-existent provider
	_, err = GetEmbeddingProvider("nonexistent", ProviderConfig{APIKey: "test"})
	if err == nil {
		t.Error("Expected error for non-existent provider, got nil")
	}
}

func TestRegisterEmbeddingProvider_Priority(t *testing.T) {
	// Custom provider name for testing
	testProvider := ProviderName("test-embedding-provider")

	// Register with low priority
	RegisterEmbeddingProvider(testProvider, func(config ProviderConfig) (provider.EmbeddingProvider, error) {
		return nil, nil
	}, 0)

	// Register with higher priority
	var highPriorityUsed bool
	RegisterEmbeddingProvider(testProvider, func(config ProviderConfig) (provider.EmbeddingProvider, error) {
		highPriorityUsed = true
		return nil, nil
	}, 10)

	// Get factory and verify high priority was used
	factory := GetEmbeddingProviderFactory(testProvider)
	if factory == nil {
		t.Fatal("Expected factory, got nil")
	}

	_, _ = factory(ProviderConfig{})
	if !highPriorityUsed {
		t.Error("Expected high priority factory to be used")
	}
}
