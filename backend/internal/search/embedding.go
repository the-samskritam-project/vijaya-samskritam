package search

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// EmbeddingGenerator generates vector embeddings for text queries
type EmbeddingGenerator interface {
	GenerateEmbedding(ctx context.Context, text string, dimension int) ([]float32, error)
}

// OpenAIEmbeddingGenerator generates embeddings using OpenAI API
type OpenAIEmbeddingGenerator struct {
	client *openai.Client
	model  openai.EmbeddingModel
}

// NewOpenAIEmbeddingGenerator creates a new OpenAI embedding generator
func NewOpenAIEmbeddingGenerator(apiKey string) *OpenAIEmbeddingGenerator {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	client := openai.NewClient(apiKey)
	return &OpenAIEmbeddingGenerator{
		client: client,
		model:  openai.SmallEmbedding3, // text-embedding-3-small
	}
}

// GenerateEmbedding generates an embedding for the given text with the specified dimension
// Note: OpenAI's text-embedding-3-small supports dimensions 512-1536.
// For 384 dimensions (Apte dictionary), we generate 512-dim embeddings and truncate to 384.
// This is a limitation - ideally, Apte embeddings should be regenerated with 512+ dimensions
// or a different embedding model should be used that natively supports 384 dimensions.
func (g *OpenAIEmbeddingGenerator) GenerateEmbedding(ctx context.Context, text string, dimension int) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: g.model,
	}

	// Handle different dimension requirements
	if dimension == 384 {
		// OpenAI's minimum is 512, so generate 512 and truncate to 384
		// Note: This may not be optimal for semantic search quality
		req.Dimensions = 512
	} else if dimension == 1536 {
		// Use default 1536 for text-embedding-3-small
		req.Dimensions = 1536
	} else if dimension >= 512 && dimension <= 1536 {
		// Use the specified dimension if within valid range
		req.Dimensions = dimension
	} else {
		// Default to 1536 if dimension is out of range
		req.Dimensions = 1536
	}

	resp, err := g.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	embedding := resp.Data[0].Embedding

	// Truncate to 384 if needed (for Apte dictionary compatibility)
	if dimension == 384 && len(embedding) == 512 {
		embedding = embedding[:384]
	}

	return embedding, nil
}
