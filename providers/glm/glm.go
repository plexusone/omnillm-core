// Package glm provides GLM (Zhipu AI BigModel) API client implementation.
// GLM uses an OpenAI-compatible API endpoint.
package glm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DefaultBaseURL is the Zhipu AI (GLM) API base URL
const DefaultBaseURL = "https://open.bigmodel.cn/api/paas/v4"

// Client implements GLM API client
type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// New creates a new GLM client
func New(apiKey, baseURL string, httpClient *http.Client) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 60 * time.Second}
	}

	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  httpClient,
	}
}

// Name returns the provider name
func (c *Client) Name() string {
	return "glm"
}

// CreateCompletion creates a chat completion
func (c *Client) CreateCompletion(ctx context.Context, req *Request) (*Response, error) {
	if req.Model == "" {
		return nil, fmt.Errorf("model cannot be empty")
	}
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("messages cannot be empty")
	}

	req.Stream = boolPtr(false)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq) //nolint:gosec // G704: baseURL is configured at client init, not user-controlled per-request
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// CreateCompletionStream creates a streaming chat completion
func (c *Client) CreateCompletionStream(ctx context.Context, req *Request) (*Stream, error) {
	if req.Model == "" {
		return nil, fmt.Errorf("model cannot be empty")
	}
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("messages cannot be empty")
	}

	req.Stream = boolPtr(true)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.client.Do(httpReq) //nolint:gosec // G704: baseURL is configured at client init, not user-controlled per-request
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, c.handleErrorResponse(resp)
	}

	return &Stream{
		response: resp,
		scanner:  bufio.NewScanner(resp.Body),
	}, nil
}

// Close closes the client
func (c *Client) Close() error {
	return nil
}

// handleErrorResponse handles error responses from GLM API
func (c *Client) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response")
	}

	var errorResp struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &errorResp); err != nil {
		return fmt.Errorf("GLM API error (status %d): %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("GLM API error: %s", errorResp.Error.Message)
}

// Stream implements streaming for GLM
type Stream struct {
	response *http.Response
	scanner  *bufio.Scanner
	closed   bool
}

// Recv receives the next chunk from the stream
func (s *Stream) Recv() (*StreamChunk, error) {
	if s.closed {
		return nil, fmt.Errorf("stream is closed")
	}

	for s.scanner.Scan() {
		line := s.scanner.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return nil, io.EOF
			}

			var chunk StreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			return &chunk, nil
		}
	}

	if err := s.scanner.Err(); err != nil {
		return nil, fmt.Errorf("stream error: %w", err)
	}

	return nil, io.EOF
}

// Close closes the stream
func (s *Stream) Close() error {
	if !s.closed {
		s.closed = true
		return s.response.Body.Close()
	}
	return nil
}

// boolPtr is a helper to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}
