package http

import (
	"encoding/json"
	"net/http"

	"vijaya-samskritam/backend/internal/corpus"
	"vijaya-samskritam/backend/internal/service"
)

// Handler handles HTTP requests
type Handler struct {
	searchService *service.SearchService
}

// NewHandler creates a new HTTP handler
func NewHandler(searchService *service.SearchService) *Handler {
	return &Handler{
		searchService: searchService,
	}
}

// SearchResponse represents the HTTP response
type SearchResponse struct {
	Results []corpus.SearchResult `json:"results"`
	Query   string                `json:"query"`
	Total   int                   `json:"total"`
}

// Search handles the /search endpoint
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	// Default limit per corpus
	limitPerCorpus := 10
	if limit := r.URL.Query().Get("limit"); limit != "" {
		// Parse limit if provided (simplified - no error handling for brevity)
		// In production, add proper parsing
	}

	// Perform search
	results, err := h.searchService.Search(r.Context(), query, limitPerCorpus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send response
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
