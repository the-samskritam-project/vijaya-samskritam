package service

import (
	"context"
	"sync"

	"vijaya-samskritam/backend/internal/corpus"
	"vijaya-samskritam/backend/internal/search"
)

// SearchService orchestrates search across multiple corpora
type SearchService struct {
	registry           *corpus.Registry
	embeddingGenerator search.EmbeddingGenerator
	scoreThreshold     float64 // Minimum score threshold to filter results
}

// NewSearchService creates a new search service
func NewSearchService(registry *corpus.Registry, embeddingGenerator search.EmbeddingGenerator, scoreThreshold float64) *SearchService {
	return &SearchService{
		registry:           registry,
		embeddingGenerator: embeddingGenerator,
		scoreThreshold:     scoreThreshold,
	}
}

// SearchResult represents aggregated search results
type SearchResult struct {
	Results []corpus.SearchResult `json:"results"`
	Query   string                `json:"query"`
	Total   int                   `json:"total"`
}

// Search performs semantic search across all corpora in parallel
func (s *SearchService) Search(ctx context.Context, query string, limitPerCorpus int) (*SearchResult, error) {
	// Get all corpora
	corpora := s.registry.GetAll()

	// Use WaitGroup to coordinate parallel searches
	var wg sync.WaitGroup
	var mu sync.Mutex
	allResults := make([]corpus.SearchResult, 0)

	// Search each corpus in parallel
	for _, corp := range corpora {
		wg.Add(1)
		go func(c corpus.Corpus) {
			defer wg.Done()

			// Generate query embedding with the correct dimension for this corpus
			dimension := c.EmbeddingDimension()
			queryEmbedding, err := s.embeddingGenerator.GenerateEmbedding(ctx, query, dimension)
			if err != nil {
				// Log error but continue with other corpora
				return
			}

			// Use corpus-specific limit
			corpusLimit := c.SearchLimit()
			results, err := c.Search(ctx, queryEmbedding, corpusLimit)
			if err != nil {
				// Log error but continue with other corpora
				return
			}

			// Filter results by score threshold
			filteredResults := make([]corpus.SearchResult, 0)
			for _, result := range results {
				if result.Score >= s.scoreThreshold {
					filteredResults = append(filteredResults, result)
				}
			}

			mu.Lock()
			allResults = append(allResults, filteredResults...)
			mu.Unlock()
		}(corp)
	}

	// Wait for all searches to complete
	wg.Wait()

	// Sort results by score (highest first)
	// Simple insertion sort for now - can be optimized
	for i := 1; i < len(allResults); i++ {
		key := allResults[i]
		j := i - 1
		for j >= 0 && allResults[j].Score < key.Score {
			allResults[j+1] = allResults[j]
			j--
		}
		allResults[j+1] = key
	}

	return &SearchResult{
		Results: allResults,
		Query:   query,
		Total:   len(allResults),
	}, nil
}
