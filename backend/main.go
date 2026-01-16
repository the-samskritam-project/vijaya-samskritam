package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vijaya-samskritam/backend/internal/config"
	"vijaya-samskritam/backend/internal/corpus"
	"vijaya-samskritam/backend/internal/db"
	httphandler "vijaya-samskritam/backend/internal/http"
	"vijaya-samskritam/backend/internal/search"
	"vijaya-samskritam/backend/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize MongoDB client
	mongoClient, err := db.NewClient(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Close(ctx)

	// Get underlying MongoDB client for vector search
	var mongoClientRaw *mongo.Client
	// We need to access the underlying client - we'll need to modify db.Client
	// For now, let's create a wrapper that exposes it
	mongoClientRaw = mongoClient.Client()

	// Initialize vector search
	vectorSearch := search.NewVectorSearch(mongoClientRaw)

	// Initialize OpenAI embedding generator
	// OPENAI_API_KEY should be set as environment variable
	embeddingGenerator := search.NewOpenAIEmbeddingGenerator("")

	// Initialize corpus registry
	corpusRegistry := corpus.NewRegistry(mongoClient, vectorSearch)

	// Initialize search service with score threshold from config
	searchService := service.NewSearchService(corpusRegistry, embeddingGenerator, cfg.ScoreThreshold)

	// Initialize HTTP handler with API key
	handler := httphandler.NewHandler(searchService, cfg.APIKey)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/search", handler.Search)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
