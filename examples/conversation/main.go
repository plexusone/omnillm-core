package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/plexusone/omnillm-core"
)

func main() {
	// Interactive conversation example
	fmt.Println("=== Interactive Conversation Example ===")
	fmt.Println("This example demonstrates maintaining conversation context across multiple providers.")
	fmt.Println("Type 'quit' to exit, 'switch' to change provider")
	fmt.Println()

	if err := runConversation(); err != nil {
		log.Fatal(err)
	}
}

func runConversation() error {
	scanner := bufio.NewScanner(os.Stdin)
	messages := []omnillm.Message{
		{
			Role:    omnillm.RoleSystem,
			Content: "You are a helpful assistant. Keep your responses concise and friendly.",
		},
	}

	currentProvider := omnillm.ProviderNameOpenAI
	client, err := createClient(currentProvider)
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Printf("Current provider: %s\n", currentProvider)
	fmt.Print("You: ")

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			fmt.Print("You: ")
			continue
		}

		if input == "quit" {
			break
		}

		if input == "switch" {
			// Switch to next provider
			client.Close()
			currentProvider = getNextProvider(currentProvider)
			client, err = createClient(currentProvider)
			if err != nil {
				log.Printf("Failed to switch provider: %v", err)
				fmt.Print("You: ")
				continue
			}
			fmt.Printf("\nSwitched to provider: %s\n", currentProvider)
			fmt.Print("You: ")
			continue
		}

		// Add user message
		messages = append(messages, omnillm.Message{
			Role:    omnillm.RoleUser,
			Content: input,
		})

		// Get response
		response, err := client.CreateChatCompletion(context.Background(), &omnillm.ChatCompletionRequest{
			Model:       getModelForProvider(currentProvider),
			Messages:    messages,
			MaxTokens:   intPtr(150),
			Temperature: float64Ptr(0.7),
		})
		if err != nil {
			log.Printf("Error: %v", err)
			fmt.Print("You: ")
			continue
		}

		assistantMessage := response.Choices[0].Message.Content
		fmt.Printf("Assistant (%s): %s\n", currentProvider, assistantMessage)

		// Add assistant response to conversation
		messages = append(messages, omnillm.Message{
			Role:    omnillm.RoleAssistant,
			Content: assistantMessage,
		})

		// Keep conversation history manageable (last 10 messages + system message)
		if len(messages) > 11 {
			messages = append(messages[:1], messages[len(messages)-10:]...)
		}

		fmt.Print("You: ")
	}

	return nil
}

func createClient(provider omnillm.ProviderName) (*omnillm.ChatClient, error) {
	var providerConfig omnillm.ProviderConfig

	switch provider {
	case omnillm.ProviderNameOpenAI:
		providerConfig = omnillm.ProviderConfig{Provider: provider, APIKey: os.Getenv("OPENAI_API_KEY")}
	case omnillm.ProviderNameAnthropic:
		providerConfig = omnillm.ProviderConfig{Provider: provider, APIKey: os.Getenv("ANTHROPIC_API_KEY")}
	case omnillm.ProviderNameBedrock:
		providerConfig = omnillm.ProviderConfig{Provider: provider, Region: "us-east-1"}
	default:
		providerConfig = omnillm.ProviderConfig{Provider: provider}
	}

	return omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{providerConfig},
	})
}

func getNextProvider(current omnillm.ProviderName) omnillm.ProviderName {
	switch current {
	case omnillm.ProviderNameOpenAI:
		return omnillm.ProviderNameAnthropic
	case omnillm.ProviderNameAnthropic:
		return omnillm.ProviderNameBedrock
	case omnillm.ProviderNameBedrock:
		return omnillm.ProviderNameOpenAI
	default:
		return omnillm.ProviderNameOpenAI
	}
}

func getModelForProvider(provider omnillm.ProviderName) string {
	switch provider {
	case omnillm.ProviderNameOpenAI:
		return omnillm.ModelGPT4oMini
	case omnillm.ProviderNameAnthropic:
		return omnillm.ModelClaude3Haiku
	case omnillm.ProviderNameBedrock:
		return omnillm.ModelBedrockClaude3Sonnet
	default:
		return omnillm.ModelGPT4oMini
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
