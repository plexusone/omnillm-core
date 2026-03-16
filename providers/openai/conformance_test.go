package openai

import (
	"os"
	"testing"

	"github.com/plexusone/omnillm/provider/providertest"
)

func TestConformance(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	p := NewProvider(apiKey, "", nil)

	providertest.RunAll(t, providertest.Config{
		Provider:        p,
		SkipIntegration: apiKey == "",
		TestModel:       "gpt-4o-mini",
	})
}
