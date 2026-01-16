package corpus

import (
	"context"

	"vijaya-samskritam/backend/internal/db"
	"vijaya-samskritam/backend/internal/search"
)

// Pancatantra implements the Corpus interface for Pancatantra
type Pancatantra struct {
	client       *db.Client
	vectorSearch *search.VectorSearch
	mapper       *FieldMapper
}

// NewPancatantra creates a new Pancatantra corpus instance
func NewPancatantra(client *db.Client, vectorSearch *search.VectorSearch) *Pancatantra {
	return &Pancatantra{
		client:       client,
		vectorSearch: vectorSearch,
		mapper: &FieldMapper{
			DevanagariField:  "transliterated_devanagari",
			TranslationField: "full_translation",
		},
	}
}

// Name returns the corpus name
func (p *Pancatantra) Name() string {
	return "pancatantra"
}

// Database returns the database name
func (p *Pancatantra) Database() string {
	return "pancatantra"
}

// Collection returns the collection name
func (p *Pancatantra) Collection() string {
	return "pancatantra_vector_search"
}

// EmbeddingField returns the embedding field name
func (p *Pancatantra) EmbeddingField() string {
	return "embedding"
}

// IndexName returns the vector search index name
func (p *Pancatantra) IndexName() string {
	return "vector_index_pancatantra"
}

// EmbeddingDimension returns the embedding dimension for this corpus
func (p *Pancatantra) EmbeddingDimension() int {
	return 1536
}

// SearchLimit returns the maximum number of results for this corpus
func (p *Pancatantra) SearchLimit() int {
	return 4
}

// Search performs vector search on the corpus
func (p *Pancatantra) Search(ctx context.Context, queryEmbedding []float32, limit int) ([]SearchResult, error) {
	docs, err := p.vectorSearch.Search(ctx, p.Database(), p.Collection(), p.EmbeddingField(), p.IndexName(), queryEmbedding, limit)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(docs))
	for _, doc := range docs {
		score, _ := doc["score"].(float64)
		result, err := p.mapper.MapDocument(p.Name(), doc, score)
		if err != nil {
			continue
		}
		results = append(results, *result)
	}

	return results, nil
}
