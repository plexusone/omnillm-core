// Package openai provides OpenAI provider adapter for the OmniLLM unified interface
package openai

import (
	"context"
	"net/http"

	"github.com/plexusone/omnillm-core/provider"
)

// Provider represents the OpenAI provider adapter
type Provider struct {
	client *Client
}

// NewProvider creates a new OpenAI provider adapter
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
	// Convert from unified format to OpenAI format
	openaiReq := &Request{
		Model:            req.Model,
		MaxTokens:        req.MaxTokens,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		Stop:             req.Stop,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		LogitBias:        req.LogitBias,
		User:             req.User,
		Seed:             req.Seed,
		N:                req.N,
		Logprobs:         req.Logprobs,
		TopLogprobs:      req.TopLogprobs,
	}

	// Convert response format if provided
	if req.ResponseFormat != nil {
		openaiReq.ResponseFormat = &ResponseFormat{
			Type: req.ResponseFormat.Type,
		}
	}

	// Convert tools
	for _, tool := range req.Tools {
		openaiReq.Tools = append(openaiReq.Tools, Tool{
			Type: tool.Type,
			Function: ToolSpec{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				Parameters:  tool.Function.Parameters,
			},
		})
	}
	openaiReq.ToolChoice = req.ToolChoice

	// Convert messages
	for _, msg := range req.Messages {
		openaiMsg := Message{
			Role:       string(msg.Role),
			Content:    msg.Content,
			Name:       msg.Name,
			ToolCallID: msg.ToolCallID,
		}
		// Convert tool calls if present
		for _, tc := range msg.ToolCalls {
			openaiMsg.ToolCalls = append(openaiMsg.ToolCalls, ToolCall{
				ID:   tc.ID,
				Type: tc.Type,
				Function: ToolFunction{
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			})
		}
		openaiReq.Messages = append(openaiReq.Messages, openaiMsg)
	}

	resp, err := p.client.CreateCompletion(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	// Convert tool calls from response
	var toolCalls []provider.ToolCall
	for _, tc := range resp.Choices[0].Message.ToolCalls {
		toolCalls = append(toolCalls, provider.ToolCall{
			ID:   tc.ID,
			Type: tc.Type,
			Function: provider.ToolFunction{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			},
		})
	}

	// Convert back to unified format
	return &provider.ChatCompletionResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: []provider.ChatCompletionChoice{
			{
				Index: 0,
				Message: provider.Message{
					Role:      provider.Role(resp.Choices[0].Message.Role),
					Content:   resp.Choices[0].Message.Content,
					ToolCalls: toolCalls,
				},
				FinishReason: resp.Choices[0].FinishReason,
			},
		},
		Usage: provider.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}

// CreateChatCompletionStream creates a streaming chat completion
func (p *Provider) CreateChatCompletionStream(ctx context.Context, req *provider.ChatCompletionRequest) (provider.ChatCompletionStream, error) {
	// Convert from unified format to OpenAI format
	openaiReq := &Request{
		Model:            req.Model,
		MaxTokens:        req.MaxTokens,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		Stop:             req.Stop,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		LogitBias:        req.LogitBias,
		User:             req.User,
		Seed:             req.Seed,
		N:                req.N,
		Logprobs:         req.Logprobs,
		TopLogprobs:      req.TopLogprobs,
	}

	// Convert response format if provided
	if req.ResponseFormat != nil {
		openaiReq.ResponseFormat = &ResponseFormat{
			Type: req.ResponseFormat.Type,
		}
	}

	// Convert messages
	for _, msg := range req.Messages {
		openaiReq.Messages = append(openaiReq.Messages, Message{
			Role:    string(msg.Role),
			Content: msg.Content,
			Name:    msg.Name,
		})
	}

	stream, err := p.client.CreateCompletionStream(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	return &StreamAdapter{stream: stream}, nil
}

// Close closes the provider
func (p *Provider) Close() error {
	return p.client.Close()
}

// StreamAdapter adapts OpenAI stream to unified interface
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
		ID:      chunk.ID,
		Object:  chunk.Object,
		Created: chunk.Created,
		Model:   chunk.Model,
	}

	if chunk.Usage != nil {
		result.Usage = &provider.Usage{
			PromptTokens:     chunk.Usage.PromptTokens,
			CompletionTokens: chunk.Usage.CompletionTokens,
			TotalTokens:      chunk.Usage.TotalTokens,
		}
	}

	for _, choice := range chunk.Choices {
		result.Choices = append(result.Choices, provider.ChatCompletionChoice{
			Index:        choice.Index,
			FinishReason: choice.FinishReason,
		})
		if choice.Delta != nil {
			result.Choices[len(result.Choices)-1].Delta = &provider.Message{
				Role:    provider.Role(choice.Delta.Role),
				Content: choice.Delta.Content,
			}
		}
	}

	return result, nil
}

// Close closes the stream
func (s *StreamAdapter) Close() error {
	return s.stream.Close()
}

// EmbeddingProvider represents the OpenAI embedding provider adapter
type EmbeddingProvider struct {
	client *Client
}

// NewEmbeddingProvider creates a new OpenAI embedding provider adapter
func NewEmbeddingProvider(apiKey, baseURL string, httpClient *http.Client) provider.EmbeddingProvider {
	client := New(apiKey, baseURL, httpClient)
	return &EmbeddingProvider{client: client}
}

// Name returns the provider name
func (p *EmbeddingProvider) Name() string {
	return p.client.Name()
}

// CreateEmbedding creates embeddings for the given input texts
func (p *EmbeddingProvider) CreateEmbedding(ctx context.Context, req *provider.EmbeddingRequest) (*provider.EmbeddingResponse, error) {
	// Convert from unified format to OpenAI format
	openaiReq := &EmbeddingRequest{
		Model:      req.Model,
		Input:      req.Input,
		Dimensions: req.Dimensions,
		User:       req.User,
	}

	if req.EncodingFormat != "" {
		openaiReq.EncodingFormat = string(req.EncodingFormat)
	}

	resp, err := p.client.CreateEmbedding(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	// Convert back to unified format
	result := &provider.EmbeddingResponse{
		Object: resp.Object,
		Model:  resp.Model,
		Usage: provider.EmbeddingUsage{
			PromptTokens: resp.Usage.PromptTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
	}

	for _, data := range resp.Data {
		result.Data = append(result.Data, provider.EmbeddingData{
			Object:    data.Object,
			Index:     data.Index,
			Embedding: data.Embedding,
		})
	}

	return result, nil
}

// Close closes the provider
func (p *EmbeddingProvider) Close() error {
	return p.client.Close()
}
