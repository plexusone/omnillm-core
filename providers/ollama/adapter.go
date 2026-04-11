// Package ollama provides Ollama provider adapter for the OmniLLM unified interface
package ollama

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/plexusone/omnillm-core/provider"
)

// Provider represents the Ollama provider adapter
type Provider struct {
	client *Client
}

// NewProvider creates a new Ollama provider adapter
func NewProvider(baseURL string, httpClient *http.Client) provider.Provider {
	client := New(baseURL, httpClient)
	return &Provider{client: client}
}

// Name returns the provider name
func (p *Provider) Name() string {
	return p.client.Name()
}

// CreateChatCompletion creates a chat completion
func (p *Provider) CreateChatCompletion(ctx context.Context, req *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
	// Convert from unified format to Ollama format
	ollamaReq := &Request{
		Model: req.Model,
	}

	// Set options if provided
	if req.MaxTokens != nil || req.Temperature != nil || req.TopP != nil || req.TopK != nil || req.Seed != nil || len(req.Stop) > 0 {
		ollamaReq.Options = &Options{
			Temperature: req.Temperature,
			TopP:        req.TopP,
			TopK:        req.TopK,
			Stop:        req.Stop,
			Seed:        req.Seed,
		}
		if req.MaxTokens != nil {
			ollamaReq.Options.NumPredict = req.MaxTokens
		}
	}

	// Convert messages
	for _, msg := range req.Messages {
		ollamaReq.Messages = append(ollamaReq.Messages, Message{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	resp, err := p.client.CreateCompletion(ctx, ollamaReq)
	if err != nil {
		return nil, err
	}

	// Convert back to unified format
	return &provider.ChatCompletionResponse{
		ID:      fmt.Sprintf("ollama-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   resp.Model,
		Choices: []provider.ChatCompletionChoice{
			{
				Index: 0,
				Message: provider.Message{
					Role:    provider.Role(resp.Message.Role),
					Content: resp.Message.Content,
				},
				FinishReason: func() *string {
					if resp.Done {
						reason := "stop"
						return &reason
					}
					return nil
				}(),
			},
		},
		Usage: provider.Usage{
			PromptTokens:     resp.PromptEvalCount,
			CompletionTokens: resp.EvalCount,
			TotalTokens:      resp.PromptEvalCount + resp.EvalCount,
		},
	}, nil
}

// CreateChatCompletionStream creates a streaming chat completion
func (p *Provider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
	// Convert from unified format to Ollama format
	ollamaReq := &Request{
		Model: req.Model,
	}

	// Set options if provided
	if req.MaxTokens != nil || req.Temperature != nil || req.TopP != nil || req.TopK != nil || req.Seed != nil || len(req.Stop) > 0 {
		ollamaReq.Options = &Options{
			Temperature: req.Temperature,
			TopP:        req.TopP,
			TopK:        req.TopK,
			Stop:        req.Stop,
			Seed:        req.Seed,
		}
		if req.MaxTokens != nil {
			ollamaReq.Options.NumPredict = req.MaxTokens
		}
	}

	// Convert messages
	for _, msg := range req.Messages {
		ollamaReq.Messages = append(ollamaReq.Messages, Message{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	stream, err := p.client.CreateCompletionStream(ctx, ollamaReq)
	if err != nil {
		return nil, err
	}

	return &StreamAdapter{stream: stream}, nil
}

// Close closes the provider
func (p *Provider) Close() error {
	return p.client.Close()
}

// StreamAdapter adapts Ollama stream to unified interface
type StreamAdapter struct {
	stream *Stream
}

// Recv receives the next chunk from the stream
func (s *StreamAdapter) Recv() (*provider.ChatCompletionChunk, error) {
	chunk, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}

	// Convert to unified format
	result := &provider.ChatCompletionChunk{
		ID:      fmt.Sprintf("ollama-stream-%d", time.Now().Unix()),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   chunk.Model,
		Choices: []provider.ChatCompletionChoice{
			{
				Index: 0,
				Delta: &provider.Message{
					Role:    provider.Role(chunk.Message.Role),
					Content: chunk.Message.Content,
				},
				FinishReason: func() *string {
					if chunk.Done {
						reason := "stop"
						return &reason
					}
					return nil
				}(),
			},
		},
	}

	if chunk.Done && chunk.EvalCount > 0 {
		result.Usage = &provider.Usage{
			PromptTokens:     chunk.PromptEvalCount,
			CompletionTokens: chunk.EvalCount,
			TotalTokens:      chunk.PromptEvalCount + chunk.EvalCount,
		}
	}

	return result, nil
}

// Close closes the stream
func (s *StreamAdapter) Close() error {
	return s.stream.Close()
}
