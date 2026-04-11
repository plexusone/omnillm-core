package anthropic

import (
	"os"
	"testing"

	"github.com/plexusone/omnillm-core/provider/providertest"
)

func TestConformance(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	p := NewProvider(apiKey, "", nil)

	providertest.RunAll(t, providertest.Config{
		Provider:        p,
		SkipIntegration: apiKey == "",
		TestModel:       "claude-3-haiku-20240307",
	})
}
