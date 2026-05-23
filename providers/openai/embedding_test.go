package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plexusone/omnillm-core/provider"
)

func TestEmbeddingProvider_CreateEmbedding(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/embeddings" {
			t.Errorf("Expected path /embeddings, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got %s", r.Header.Get("Authorization"))
		}

		// Decode request body
		var req EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Model != "text-embedding-3-small" {
			t.Errorf("Expected model 'text-embedding-3-small', got %s", req.Model)
		}
		if len(req.Input) != 2 {
			t.Errorf("Expected 2 inputs, got %d", len(req.Input))
		}

		// Return mock response
		resp := EmbeddingResponse{
			Object: "list",
			Model:  "text-embedding-3-small",
			Data: []EmbeddingData{
				{
					Object:    "embedding",
					Index:     0,
					Embedding: []float64{0.1, 0.2, 0.3},
				},
				{
					Object:    "embedding",
					Index:     1,
					Embedding: []float64{0.4, 0.5, 0.6},
				},
			},
			Usage: EmbeddingUsage{
				PromptTokens: 10,
				TotalTokens:  10,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create provider
	p := NewEmbeddingProvider("test-key", server.URL, nil)

	// Test CreateEmbedding
	resp, err := p.CreateEmbedding(context.Background(), &provider.EmbeddingRequest{
		Model: "text-embedding-3-small",
		Input: []string{"Hello world", "How are you?"},
	})
	if err != nil {
		t.Fatalf("CreateEmbedding failed: %v", err)
	}

	// Verify response
	if resp.Object != "list" {
		t.Errorf("Expected object 'list', got %s", resp.Object)
	}
	if resp.Model != "text-embedding-3-small" {
		t.Errorf("Expected model 'text-embedding-3-small', got %s", resp.Model)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("Expected 2 embeddings, got %d", len(resp.Data))
	}
	if resp.Data[0].Index != 0 {
		t.Errorf("Expected first embedding index 0, got %d", resp.Data[0].Index)
	}
	if len(resp.Data[0].Embedding) != 3 {
		t.Errorf("Expected embedding length 3, got %d", len(resp.Data[0].Embedding))
	}
	if resp.Usage.PromptTokens != 10 {
		t.Errorf("Expected prompt tokens 10, got %d", resp.Usage.PromptTokens)
	}
}

func TestEmbeddingProvider_Name(t *testing.T) {
	p := NewEmbeddingProvider("test-key", "", nil)
	if p.Name() != "openai" {
		t.Errorf("Expected name 'openai', got %s", p.Name())
	}
}

func TestEmbeddingProvider_Close(t *testing.T) {
	p := NewEmbeddingProvider("test-key", "", nil)
	if err := p.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestCreateEmbedding_EmptyModel(t *testing.T) {
	p := NewEmbeddingProvider("test-key", "", nil)
	_, err := p.CreateEmbedding(context.Background(), &provider.EmbeddingRequest{
		Model: "",
		Input: []string{"test"},
	})
	if err == nil {
		t.Error("Expected error for empty model, got nil")
	}
}

func TestCreateEmbedding_EmptyInput(t *testing.T) {
	p := NewEmbeddingProvider("test-key", "", nil)
	_, err := p.CreateEmbedding(context.Background(), &provider.EmbeddingRequest{
		Model: "text-embedding-3-small",
		Input: []string{},
	})
	if err == nil {
		t.Error("Expected error for empty input, got nil")
	}
}
