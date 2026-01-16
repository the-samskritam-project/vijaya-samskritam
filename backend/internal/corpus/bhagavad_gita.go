package corpus

import (
	"context"

	"vijaya-samskritam/backend/internal/db"
	"vijaya-samskritam/backend/internal/search"
)

// BhagavadGita implements the Corpus interface for Bhagavad Gita
type BhagavadGita struct {
	client       *db.Client
	vectorSearch *search.VectorSearch
	mapper       *FieldMapper
}

// NewBhagavadGita creates a new Bhagavad Gita corpus instance
func NewBhagavadGita(client *db.Client, vectorSearch *search.VectorSearch) *BhagavadGita {
	return &BhagavadGita{
		client:       client,
		vectorSearch: vectorSearch,
		mapper: &FieldMapper{
			DevanagariField:  "transliterated_devanagari",
			TranslationField: "full_translation",
		},
	}
}

// Name returns the corpus name
func (bg *BhagavadGita) Name() string {
	return "bhagavad_gita"
}

// Database returns the database name
func (bg *BhagavadGita) Database() string {
	return "bhagavad_gita_shankara_bhasya"
}

// Collection returns the collection name
func (bg *BhagavadGita) Collection() string {
	return "bhagavad_gita_vector_search"
}

// EmbeddingField returns the embedding field name
func (bg *BhagavadGita) EmbeddingField() string {
	return "embedding"
}

// IndexName returns the vector search index name
func (bg *BhagavadGita) IndexName() string {
	return "vector_index_bhagavad_gita"
}

// EmbeddingDimension returns the embedding dimension for this corpus
func (bg *BhagavadGita) EmbeddingDimension() int {
	return 1536
}

// SearchLimit returns the maximum number of results for this corpus
func (bg *BhagavadGita) SearchLimit() int {
	return 4
}

// Search performs vector search on the corpus
func (bg *BhagavadGita) Search(ctx context.Context, queryEmbedding []float32, limit int) ([]SearchResult, error) {
	docs, err := bg.vectorSearch.Search(ctx, bg.Database(), bg.Collection(), bg.EmbeddingField(), bg.IndexName(), queryEmbedding, limit)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(docs))
	for _, doc := range docs {
		score, _ := doc["score"].(float64)
		result, err := bg.mapper.MapDocument(bg.Name(), doc, score)
		if err != nil {
			continue
		}
		results = append(results, *result)
	}

	return results, nil
}
