package corpus

import (
	"context"

	"vijaya-samskritam/backend/internal/db"
	"vijaya-samskritam/backend/internal/search"
)

// Apte implements the Corpus interface for Apte dictionary
type Apte struct {
	client       *db.Client
	vectorSearch *search.VectorSearch
	mapper       *FieldMapper
}

// NewApte creates a new Apte dictionary corpus instance
func NewApte(client *db.Client, vectorSearch *search.VectorSearch) *Apte {
	return &Apte{
		client:       client,
		vectorSearch: vectorSearch,
		mapper: &FieldMapper{
			DevanagariField:  "sanskritString",
			TranslationField: "meaning",
		},
	}
}

// Name returns the corpus name
func (a *Apte) Name() string {
	return "apte_dictionary"
}

// Database returns the database name
func (a *Apte) Database() string {
	return "apte_dictionary"
}

// Collection returns the collection name
func (a *Apte) Collection() string {
	return "entries"
}

// EmbeddingField returns the embedding field name
func (a *Apte) EmbeddingField() string {
	return "embedding"
}

// IndexName returns the vector search index name
func (a *Apte) IndexName() string {
	return "vector_index_apte"
}

// EmbeddingDimension returns the embedding dimension for this corpus
func (a *Apte) EmbeddingDimension() int {
	return 384
}

// SearchLimit returns the maximum number of results for this corpus
func (a *Apte) SearchLimit() int {
	return 4
}

// Search performs vector search on the corpus
func (a *Apte) Search(ctx context.Context, queryEmbedding []float32, limit int) ([]SearchResult, error) {
	docs, err := a.vectorSearch.Search(ctx, a.Database(), a.Collection(), a.EmbeddingField(), a.IndexName(), queryEmbedding, limit)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(docs))
	for _, doc := range docs {
		score, _ := doc["score"].(float64)
		result, err := a.mapper.MapDocument(a.Name(), doc, score)
		if err != nil {
			continue
		}
		results = append(results, *result)
	}

	return results, nil
}
