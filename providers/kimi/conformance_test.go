package kimi

import (
	"os"
	"testing"

	"github.com/plexusone/omnillm-core/provider/providertest"
)

func TestConformance(t *testing.T) {
	apiKey := os.Getenv("KIMI_API_KEY")

	p := NewProvider(apiKey, "", nil)

	providertest.RunAll(t, providertest.Config{
		Provider:        p,
		SkipIntegration: apiKey == "",
		TestModel:       "moonshot-v1-8k",
	})
}
