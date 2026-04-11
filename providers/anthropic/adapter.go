// Package anthropic provides Anthropic provider adapter for the OmniLLM unified interface
package anthropic

import (
	"context"
	"net/http"
	"time"

	"github.com/plexusone/omnillm-core/provider"
)

// Provider represents the Anthropic provider adapter
type Provider struct {
	client *Client
}

// NewProvider creates a new Anthropic provider adapter
func NewProvider(apiKey, baseURL string, httpClient *http.Client) provider.Provider {
	client := New(apiKey, baseURL, httpClient)
	return &Provider{client: client}
}

// Name returns the provider name
func (p *Provider) Name() string {
	return p.client.Name()
}

// CreateChatCompletion creates a chat completion
func (p *Provider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
	// Convert from unified format to Anthropic format
	anthropicReq := &Request{
		Model:       req.Model,
		MaxTokens:   4096, // Default
		Temperature: req.Temperature,
		TopP:        req.TopP,
		TopK:        req.TopK,
	}

	if req.MaxTokens != nil {
		anthropicReq.MaxTokens = *req.MaxTokens
	}

	// Convert messages (Anthropic separates system messages)
	var systemMessage string
	for _, msg := range req.Messages {
		switch msg.Role {
		case provider.RoleSystem:
			systemMessage = msg.Content
		case provider.RoleUser, provider.RoleAssistant:
			anthropicReq.Messages = append(anthropicReq.Messages, Message{
				Role:    string(msg.Role),
				Content: msg.Content,
			})
		}
	}

	if systemMessage != "" {
		anthropicReq.System = systemMessage
	}

	resp, err := p.client.CreateCompletion(ctx, anthropicReq)
	if err != nil {
		return nil, err
	}

	// Convert back to unified format
	var content string
	if len(resp.Content) > 0 && resp.Content[0].Type == "text" {
		content = resp.Content[0].Text
	}

	// Preserve Anthropic-specific metadata
	metadata := map[string]any{
		"anthropic_type":        resp.Type,
		"anthropic_role":        resp.Role,
		"anthropic_content":     resp.Content, // Full content array
		"anthropic_stop_reason": resp.StopReason,
	}

	return &provider.ChatCompletionResponse{
		ID:      resp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   resp.Model,
		Choices: []provider.ChatCompletionChoice{
			{
				Index: 0,
				Message: provider.Message{
					Role:    provider.RoleAssistant,
					Content: content,
				},
				FinishReason: &resp.StopReason,
			},
		},
		Usage: provider.Usage{
			PromptTokens:     resp.Usage.InputTokens,
			CompletionTokens: resp.Usage.OutputTokens,
			TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
		},
		ProviderMetadata: metadata,
	}, nil
}

// CreateChatCompletionStream creates a streaming chat completion
func (p *Provider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
	// Convert from unified format to Anthropic format
	anthropicReq := &Request{
		Model:       req.Model,
		MaxTokens:   4096, // Default
		Temperature: req.Temperature,
		TopP:        req.TopP,
		TopK:        req.TopK,
	}

	if req.MaxTokens != nil {
		anthropicReq.MaxTokens = *req.MaxTokens
	}

	// Convert messages (Anthropic separates system messages)
	var systemMessage string
	for _, msg := range req.Messages {
		switch msg.Role {
		case provider.RoleSystem:
			systemMessage = msg.Content
		case provider.RoleUser, provider.RoleAssistant:
			anthropicReq.Messages = append(anthropicReq.Messages, Message{
				Role:    string(msg.Role),
				Content: msg.Content,
			})
		}
	}

	if systemMessage != "" {
		anthropicReq.System = systemMessage
	}

	stream, err := p.client.CreateCompletionStream(ctx, anthropicReq)
	if err != nil {
		return nil, err
	}

	return &StreamAdapter{stream: stream}, nil
}

// Close closes the provider
func (p *Provider) Close() error {
	return p.client.Close()
}

// StreamAdapter adapts Anthropic stream to unified interface
type StreamAdapter struct {
	stream    *Stream
	messageID string
	model     string
}

// Recv receives the next chunk from the stream
func (s *StreamAdapter) Recv() (*provider.ChatCompletionChunk, error) {
	event, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}

	// Handle different event types
	switch event.Type {
	case "message_start":
		// Store message metadata for future chunks
		if event.Message != nil {
			s.messageID = event.Message.ID
			s.model = event.Message.Model
		}
		// Return empty chunk for message_start
		metadata := map[string]any{
			"anthropic_event_type": event.Type,
			"anthropic_message":    event.Message,
		}
		return &provider.ChatCompletionChunk{
			ID:               s.messageID,
			Object:           "chat.completion.chunk",
			Created:          time.Now().Unix(),
			Model:            s.model,
			Choices:          []provider.ChatCompletionChoice{},
			ProviderMetadata: metadata,
		}, nil

	case "content_block_delta":
		// This contains the actual text content
		var content string
		if event.Delta != nil && event.Delta.Type == "text_delta" {
			content = event.Delta.Text
		}

		metadata := map[string]any{
			"anthropic_event_type": event.Type,
			"anthropic_delta":      event.Delta,
			"anthropic_index":      event.Index,
		}

		return &provider.ChatCompletionChunk{
			ID:      s.messageID,
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   s.model,
			Choices: []provider.ChatCompletionChoice{
				{
					Index: 0,
					Delta: &provider.Message{
						Role:    provider.RoleAssistant,
						Content: content,
					},
				},
			},
			ProviderMetadata: metadata,
		}, nil

	case "message_delta":
		// Contains stop reason and usage info
		var finishReason *string
		if event.Delta != nil && event.Delta.StopReason != "" {
			finishReason = &event.Delta.StopReason
		}

		metadata := map[string]any{
			"anthropic_event_type": event.Type,
			"anthropic_delta":      event.Delta,
			"anthropic_usage":      event.Usage,
		}

		chunk := &provider.ChatCompletionChunk{
			ID:      s.messageID,
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   s.model,
			Choices: []provider.ChatCompletionChoice{
				{
					Index:        0,
					FinishReason: finishReason,
				},
			},
			ProviderMetadata: metadata,
		}

		// Add usage if available
		if event.Usage != nil {
			chunk.Usage = &provider.Usage{
				CompletionTokens: event.Usage.OutputTokens,
			}
		}

		return chunk, nil

	case "message_stop":
		// End of stream - return empty chunk
		metadata := map[string]any{
			"anthropic_event_type": event.Type,
		}
		return &provider.ChatCompletionChunk{
			ID:               s.messageID,
			Object:           "chat.completion.chunk",
			Created:          time.Now().Unix(),
			Model:            s.model,
			Choices:          []provider.ChatCompletionChoice{},
			ProviderMetadata: metadata,
		}, nil

	default:
		// For other event types, continue to next event
		return s.Recv()
	}
}

// Close closes the stream
func (s *StreamAdapter) Close() error {
	return s.stream.Close()
}
