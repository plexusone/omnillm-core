package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/plexusone/omnillm-core"
)

// ProviderDemo holds configuration for demonstrating a specific provider
type ProviderDemo struct {
	Name     string
	Config   omnillm.ClientConfig
	Model    string
	Messages []omnillm.Message
}

func main() {
	// Define all provider demonstrations
	demos := []ProviderDemo{
		{
			Name: "OpenAI",
			Config: omnillm.ClientConfig{
				Providers: []omnillm.ProviderConfig{
					{Provider: omnillm.ProviderNameOpenAI, APIKey: os.Getenv("OPENAI_API_KEY")},
				},
			},
			Model: omnillm.ModelGPT4o,
			Messages: []omnillm.Message{
				{
					Role:    omnillm.RoleUser,
					Content: "Hello! Can you explain what a unified LLM SDK is?",
				},
			},
		},
		{
			Name: "Anthropic",
			Config: omnillm.ClientConfig{
				Providers: []omnillm.ProviderConfig{
					{Provider: omnillm.ProviderNameAnthropic, APIKey: os.Getenv("ANTHROPIC_API_KEY")},
				},
			},
			Model: omnillm.ModelClaude3Sonnet,
			Messages: []omnillm.Message{
				{
					Role:    omnillm.RoleSystem,
					Content: "You are a helpful assistant that explains technical concepts clearly.",
				},
				{
					Role:    omnillm.RoleUser,
					Content: "What are the benefits of using a unified SDK for multiple LLM providers?",
				},
			},
		},
		{
			Name: "AWS Bedrock",
			Config: omnillm.ClientConfig{
				Providers: []omnillm.ProviderConfig{
					{Provider: omnillm.ProviderNameBedrock, Region: "us-east-1"},
				},
			},
			Model: omnillm.ModelBedrockClaude3Sonnet,
			Messages: []omnillm.Message{
				{
					Role:    omnillm.RoleUser,
					Content: "Explain the advantages of using AWS Bedrock for LLM deployments.",
				},
			},
		},
		{
			Name: "Ollama (Local)",
			Config: omnillm.ClientConfig{
				Providers: []omnillm.ProviderConfig{
					{Provider: omnillm.ProviderNameOllama, BaseURL: "http://localhost:11434"},
				},
			},
			Model: "llama3", // Use the model name as it appears in "ollama list"
			Messages: []omnillm.Message{
				{
					Role:    omnillm.RoleUser,
					Content: "What are the benefits of running LLMs locally with Ollama?",
				},
			},
		},
	}

	// Run all demonstrations
	for _, demo := range demos {
		fmt.Printf("=== %s Example ===\n", demo.Name)
		if err := demonstrateProvider(demo); err != nil {
			log.Printf("%s error: %v", demo.Name, err)
		}
		fmt.Println()
	}
}

// demonstrateProvider is a generic function that works with any provider
func demonstrateProvider(demo ProviderDemo) error {
	client, err := omnillm.NewClient(demo.Config)
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.CreateChatCompletion(context.Background(), &omnillm.ChatCompletionRequest{
		Model:       demo.Model,
		Messages:    demo.Messages,
		MaxTokens:   intPtr(150),
		Temperature: float64Ptr(0.7),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)

	return nil
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
