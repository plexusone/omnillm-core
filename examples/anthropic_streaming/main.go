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
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is required")
	}

	// Create Anthropic client
	client, err := omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{
			{Provider: omnillm.ProviderNameAnthropic, APIKey: apiKey},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Println("=== Anthropic Claude Streaming Demo ===")
	fmt.Println()

	// Example 1: Simple streaming
	fmt.Println("Example 1: Simple question with streaming")
	fmt.Println("Question: Explain what streaming responses are in one sentence.")
	fmt.Print("\nClaude: ")

	stream, err := client.CreateChatCompletionStream(context.Background(), &omnillm.ChatCompletionRequest{
		Model: omnillm.ModelClaude3Haiku,
		Messages: []omnillm.Message{
			{
				Role:    omnillm.RoleUser,
				Content: "Explain what streaming responses are in one sentence.",
			},
		},
		MaxTokens:   intPtr(100),
		Temperature: float64Ptr(0.7),
	})
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	defer stream.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Stream error: %v", err)
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
	fmt.Println()

	// Example 2: Streaming with system message
	fmt.Println("Example 2: Creative writing with system message")
	fmt.Println("Task: Write a haiku about AI")
	fmt.Print("\nClaude: ")

	stream2, err := client.CreateChatCompletionStream(context.Background(), &omnillm.ChatCompletionRequest{
		Model: omnillm.ModelClaude3Sonnet,
		Messages: []omnillm.Message{
			{
				Role:    omnillm.RoleSystem,
				Content: "You are a creative poet who writes thoughtful and elegant haikus.",
			},
			{
				Role:    omnillm.RoleUser,
				Content: "Write a haiku about artificial intelligence.",
			},
		},
		MaxTokens:   intPtr(100),
		Temperature: float64Ptr(0.9),
	})
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	defer stream2.Close()

	for {
		chunk, err := stream2.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Stream error: %v", err)
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
	fmt.Println()

	// Example 3: Longer streaming response
	fmt.Println("Example 3: Longer response with streaming")
	fmt.Println("Question: List 5 benefits of using Go for backend development")
	fmt.Print("\nClaude: ")

	stream3, err := client.CreateChatCompletionStream(context.Background(), &omnillm.ChatCompletionRequest{
		Model: omnillm.ModelClaude3Haiku,
		Messages: []omnillm.Message{
			{
				Role:    omnillm.RoleUser,
				Content: "List 5 key benefits of using Go for backend development. Be concise.",
			},
		},
		MaxTokens:   intPtr(300),
		Temperature: float64Ptr(0.5),
	})
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	defer stream3.Close()

	totalChunks := 0
	for {
		chunk, err := stream3.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Stream error: %v", err)
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fmt.Print(chunk.Choices[0].Delta.Content)
			totalChunks++
		}
	}
	fmt.Printf("\n\n(Received %d chunks)\n", totalChunks)

	fmt.Println("\n=== Streaming Demo Complete ===")
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
