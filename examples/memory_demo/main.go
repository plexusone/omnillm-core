package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/grokify/sogo/database/kvs"

	"github.com/plexusone/omnillm-core"
)

// mockKVS is a simple in-memory implementation of the KVS interface for demonstration
type mockKVS struct {
	data map[string]string
}

func newMockKVS() *mockKVS {
	return &mockKVS{
		data: make(map[string]string),
	}
}

func (m *mockKVS) SetString(ctx context.Context, key, val string) error {
	m.data[key] = val
	return nil
}

func (m *mockKVS) GetString(ctx context.Context, key string) (string, error) {
	val, exists := m.data[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return val, nil
}

func (m *mockKVS) GetOrDefaultString(ctx context.Context, key, def string) string {
	val, err := m.GetString(ctx, key)
	if err != nil {
		return def
	}
	return val
}

func (m *mockKVS) SetAny(ctx context.Context, key string, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	m.data[key] = string(data)
	return nil
}

func (m *mockKVS) GetAny(ctx context.Context, key string, val any) error {
	data, exists := m.data[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}
	return json.Unmarshal([]byte(data), val)
}

func main() {
	fmt.Println("=== omnillm Memory Demo ===")
	fmt.Println("This example demonstrates conversation memory using a KVS backend.")
	fmt.Println("Type 'quit' to exit, 'new' to start a new session, 'sessions' to list sessions")
	fmt.Println("Note: This uses a mock in-memory KVS for demonstration.")
	fmt.Println()

	if err := runMemoryDemo(); err != nil {
		log.Fatal(err)
	}
}

func runMemoryDemo() error {
	// Create mock KVS (in production, you'd use Redis, DynamoDB, etc.)
	mockKVS := newMockKVS()

	// Configure memory settings
	memoryConfig := omnillm.MemoryConfig{
		MaxMessages: 20, // Keep last 20 messages per session
		TTL:         2 * time.Hour,
		KeyPrefix:   "demo:session",
	}

	// Create client with memory
	client, err := createClientWithMemory(mockKVS, &memoryConfig)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	scanner := bufio.NewScanner(os.Stdin)
	currentSessionID := "session-1"

	// Initialize the first session with a system message
	err = client.CreateConversationWithSystemMessage(
		context.Background(),
		currentSessionID,
		"You are a helpful assistant with persistent memory. You remember our previous conversations in this session.",
	)
	if err != nil {
		log.Printf("Warning: Failed to create initial conversation: %v", err)
	}

	fmt.Printf("Current session: %s\n", currentSessionID)
	fmt.Print("You: ")

	sessionCounter := 1

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			fmt.Print("You: ")
			continue
		}

		if input == "quit" {
			break
		}

		if input == "new" {
			sessionCounter++
			currentSessionID = fmt.Sprintf("session-%d", sessionCounter)

			// Initialize new session
			err = client.CreateConversationWithSystemMessage(
				context.Background(),
				currentSessionID,
				"You are a helpful assistant with persistent memory. This is a new conversation session.",
			)
			if err != nil {
				log.Printf("Warning: Failed to create new conversation: %v", err)
			}

			fmt.Printf("\nStarted new session: %s\n", currentSessionID)
			fmt.Print("You: ")
			continue
		}

		if input == "sessions" {
			fmt.Printf("\nAvailable sessions: ")
			for i := 1; i <= sessionCounter; i++ {
				sessionID := fmt.Sprintf("session-%d", i)
				messages, err := client.GetConversationMessages(context.Background(), sessionID)
				if err == nil {
					marker := ""
					if sessionID == currentSessionID {
						marker = " (current)"
					}
					fmt.Printf("%s (%d messages)%s, ", sessionID, len(messages), marker)
				}
			}
			fmt.Println()
			fmt.Print("You: ")
			continue
		}

		// Check if switching to an existing session
		if strings.HasPrefix(input, "switch ") {
			newSessionID := strings.TrimPrefix(input, "switch ")
			// Check if session exists by trying to load it
			_, err := client.LoadConversation(context.Background(), newSessionID)
			if err == nil {
				currentSessionID = newSessionID
				fmt.Printf("Switched to session: %s\n", currentSessionID)
			} else {
				fmt.Printf("Session '%s' not found. Type 'sessions' to see available sessions.\n", newSessionID)
			}
			fmt.Print("You: ")
			continue
		}

		// Create user message
		userMessage := omnillm.Message{
			Role:    omnillm.RoleUser,
			Content: input,
		}

		// Use memory-aware completion
		response, err := client.CreateChatCompletionWithMemory(
			context.Background(),
			currentSessionID,
			&omnillm.ChatCompletionRequest{
				Model:       getAvailableModel(),
				Messages:    []omnillm.Message{userMessage},
				MaxTokens:   intPtr(200),
				Temperature: float64Ptr(0.7),
			},
		)
		if err != nil {
			log.Printf("Error: %v", err)
			fmt.Print("You: ")
			continue
		}

		if len(response.Choices) > 0 {
			fmt.Printf("Assistant: %s\n", response.Choices[0].Message.Content)
		}

		// Show conversation stats
		messages, err := client.GetConversationMessages(context.Background(), currentSessionID)
		if err == nil {
			fmt.Printf("(Session: %s, Messages: %d)\n", currentSessionID, len(messages))
		}

		fmt.Print("You: ")
	}

	return nil
}

func createClientWithMemory(kvsClient kvs.Client, memoryConfig *omnillm.MemoryConfig) (*omnillm.ChatClient, error) {
	// Try to get API keys from environment
	var provider omnillm.ProviderName
	var apiKey string

	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey != "" {
		provider = omnillm.ProviderNameOpenAI
		apiKey = openaiKey
	} else if anthropicKey := os.Getenv("ANTHROPIC_API_KEY"); anthropicKey != "" {
		provider = omnillm.ProviderNameAnthropic
		apiKey = anthropicKey
	} else {
		// Fall back to Bedrock (which doesn't require API key in config)
		provider = omnillm.ProviderNameBedrock
	}

	providerConfig := omnillm.ProviderConfig{
		Provider: provider,
		APIKey:   apiKey,
	}
	if provider == omnillm.ProviderNameBedrock {
		providerConfig.Region = "us-east-1"
	}

	return omnillm.NewClient(omnillm.ClientConfig{
		Providers:    []omnillm.ProviderConfig{providerConfig},
		Memory:       kvsClient,
		MemoryConfig: memoryConfig,
	})
}

func getAvailableModel() string {
	if os.Getenv("OPENAI_API_KEY") != "" {
		return omnillm.ModelGPT4oMini
	} else if os.Getenv("ANTHROPIC_API_KEY") != "" {
		return omnillm.ModelClaude3Haiku
	} else {
		return omnillm.ModelBedrockClaude3Sonnet
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
