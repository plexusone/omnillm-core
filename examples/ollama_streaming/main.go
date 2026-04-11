package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/plexusone/omnillm-core"
)

func main() {
	// Create a client for Ollama
	client, err := omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{
			{Provider: omnillm.ProviderNameOllama, BaseURL: "http://localhost:11434"}, // Optional - default is localhost
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("Testing Ollama streaming with OmniLLM...")
	fmt.Println("Make sure you have Ollama running locally with a model installed.")
	fmt.Println("Example: ollama run llama3:8b")
	fmt.Println()

	// Create a streaming chat completion request
	stream, err := client.CreateChatCompletionStream(context.Background(), &omnillm.ChatCompletionRequest{
		Model: omnillm.ModelOllamaLlama3_8B, // You can use any model you have installed
		Messages: []omnillm.Message{
			{
				Role:    omnillm.RoleUser,
				Content: "Tell me a short story about AI assistants helping developers.",
			},
		},
		MaxTokens:   &[]int{200}[0],
		Temperature: &[]float64{0.8}[0],
	})
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	fmt.Print("AI Response: ")
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}

		// Print usage information when stream is complete
		if chunk.Usage != nil {
			fmt.Printf("\n\nTokens used: %d (prompt: %d, completion: %d)\n",
				chunk.Usage.TotalTokens,
				chunk.Usage.PromptTokens,
				chunk.Usage.CompletionTokens)
		}
	}
	fmt.Println()
}
