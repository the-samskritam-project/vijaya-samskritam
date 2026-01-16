package corpus

import "context"

// Corpus defines the interface for corpus implementations
type Corpus interface {
	Name() string
	Database() string
	Collection() string
	EmbeddingField() string
	IndexName() string
	EmbeddingDimension() int
	SearchLimit() int // Returns the maximum number of results to return for this corpus
	Search(ctx context.Context, queryEmbedding []float32, limit int) ([]SearchResult, error)
}

// SearchResult represents a standardized search result
type SearchResult struct {
	// Standard mapped fields
	CorpusName      string
	DevanagariText  string
	FullTranslation string
	DocumentID      string
	ChapterNumber   *int
	VerseNumber     *string
	ProseNumber     *string

	// Search metadata
	Score float64
}
