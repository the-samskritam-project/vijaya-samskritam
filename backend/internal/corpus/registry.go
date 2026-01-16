package corpus

import (
	"vijaya-samskritam/backend/internal/db"
	"vijaya-samskritam/backend/internal/search"
)

// Registry manages corpus implementations
type Registry struct {
	corpora map[string]Corpus
}

// NewRegistry creates a new corpus registry with all registered corpora
func NewRegistry(client *db.Client, vectorSearch *search.VectorSearch) *Registry {
	registry := &Registry{
		corpora: make(map[string]Corpus),
	}

	// Register all corpora
	registry.Register(NewBhagavadGita(client, vectorSearch))
	registry.Register(NewPancatantra(client, vectorSearch))
	registry.Register(NewApte(client, vectorSearch))

	return registry
}

// Register registers a corpus in the registry
func (r *Registry) Register(c Corpus) {
	r.corpora[c.Name()] = c
}

// Get retrieves a corpus by name
func (r *Registry) Get(name string) (Corpus, bool) {
	corpus, ok := r.corpora[name]
	return corpus, ok
}

// GetAll returns all registered corpora
func (r *Registry) GetAll() []Corpus {
	corpora := make([]Corpus, 0, len(r.corpora))
	for _, corpus := range r.corpora {
		corpora = append(corpora, corpus)
	}
	return corpora
}
