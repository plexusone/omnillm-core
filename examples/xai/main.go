package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/plexusone/omnillm-core"
)

func main() {
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	fmt.Println("=== X.AI Grok Demo ===")
	fmt.Println()

	// Example 1: Basic chat completion
	fmt.Println("Example 1: Basic chat completion with Grok")
	if err := demonstrateBasicCompletion(apiKey); err != nil {
		log.Printf("Basic completion error: %v", err)
	}
	fmt.Println()

	// Example 2: Streaming response
	fmt.Println("Example 2: Streaming response with Grok")
	if err := demonstrateStreaming(apiKey); err != nil {
		log.Printf("Streaming error: %v", err)
	}
	fmt.Println()

	// Example 3: System message
	fmt.Println("Example 3: Using system message for role-playing")
	if err := demonstrateSystemMessage(apiKey); err != nil {
		log.Printf("System message error: %v", err)
	}
}

func demonstrateBasicCompletion(apiKey string) error {
	client, err := omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{
			{Provider: omnillm.ProviderNameXAI, APIKey: apiKey},
		},
	})
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.CreateChatCompletion(context.Background(), &omnillm.ChatCompletionRequest{
		Model: omnillm.ModelGrok4_1FastReasoning,
		Messages: []omnillm.Message{
			{
				Role:    omnillm.RoleUser,
				Content: "What makes you different from other AI assistants?",
			},
		},
		MaxTokens:   intPtr(150),
		Temperature: float64Ptr(0.7),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Grok: %s\n", response.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)

	return nil
}

func demonstrateStreaming(apiKey string) error {
	client, err := omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{
			{Provider: omnillm.ProviderNameXAI, APIKey: apiKey},
		},
	})
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Print("Grok: ")

	stream, err := client.CreateChatCompletionStream(context.Background(), &omnillm.ChatCompletionRequest{
		Model: omnillm.ModelGrok4_1FastReasoning,
		Messages: []omnillm.Message{
			{
				Role:    omnillm.RoleUser,
				Content: "Write a haiku about artificial intelligence.",
			},
		},
		MaxTokens:   intPtr(100),
		Temperature: float64Ptr(0.8),
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println()

	return nil
}

func demonstrateSystemMessage(apiKey string) error {
	client, err := omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{
			{Provider: omnillm.ProviderNameXAI, APIKey: apiKey},
		},
	})
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.CreateChatCompletion(context.Background(), &omnillm.ChatCompletionRequest{
		Model: omnillm.ModelGrok4_1FastReasoning,
		Messages: []omnillm.Message{
			{
				Role:    omnillm.RoleSystem,
				Content: "You are a witty AI assistant with a good sense of humor. Keep responses concise.",
			},
			{
				Role:    omnillm.RoleUser,
				Content: "Tell me a programmer joke.",
			},
		},
		MaxTokens:   intPtr(150),
		Temperature: float64Ptr(0.9),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Grok: %s\n", response.Choices[0].Message.Content)

	return nil
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
