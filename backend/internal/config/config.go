package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	MongoURI       string
	Port           string
	ScoreThreshold float64
}

// Load loads configuration from environment variables
// It first tries to load from .env file (if it exists), then reads from environment
// The .env file should be in the backend directory (same level as go.mod)
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		// .env file is optional - in production, env vars are typically set directly
		log.Println("No .env file found, using environment variables")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable is required")
	}

	scoreThreshold := 0.5 // Default threshold
	if thresholdStr := os.Getenv("SCORE_THRESHOLD"); thresholdStr != "" {
		if threshold, err := strconv.ParseFloat(thresholdStr, 64); err == nil {
			scoreThreshold = threshold
		}
	}

	return &Config{
		MongoURI:       mongoURI,
		Port:           getEnv("PORT", "8080"),
		ScoreThreshold: scoreThreshold,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
