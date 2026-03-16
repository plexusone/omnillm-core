package xai

import (
	"os"
	"testing"

	"github.com/plexusone/omnillm/provider/providertest"
)

func TestConformance(t *testing.T) {
	apiKey := os.Getenv("XAI_API_KEY")

	p := NewProvider(apiKey, "", nil)

	providertest.RunAll(t, providertest.Config{
		Provider:        p,
		SkipIntegration: apiKey == "",
		TestModel:       "grok-3-mini-fast",
	})
}
