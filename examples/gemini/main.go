//go:build ignore

// This example requires omnillm-gemini to be installed:
//
//	go get github.com/plexusone/omnillm-gemini
//
// Then remove the "//go:build ignore" line to enable building.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	omnillm "github.com/plexusone/omnillm-core"

	// Import omnillm-gemini to register the Gemini provider.
	// This is required because Gemini uses a heavy SDK dependency
	// and is not included in omnillm-core by default.
	_ "github.com/plexusone/omnillm-gemini"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	// Create a Gemini client
	client, err := omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{
			{Provider: omnillm.ProviderNameGemini, APIKey: apiKey},
		},
	})
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer client.Close()

	// Create a chat completion request
	response, err := client.CreateChatCompletion(context.Background(), &omnillm.ChatCompletionRequest{
		Model: omnillm.ModelGemini1_5Flash,
		Messages: []omnillm.Message{
			{
				Role:    omnillm.RoleUser,
				Content: "Hello! Can you tell me a short joke?",
			},
		},
		MaxTokens:   &[]int{150}[0],
		Temperature: &[]float64{0.7}[0],
	})
	if err != nil {
		log.Fatal("Failed to create completion:", err)
	}

	// Print the response
	fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)
}
