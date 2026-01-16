package search

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// VectorSearch performs MongoDB Atlas vector search
type VectorSearch struct {
	client *mongo.Client
}

// NewVectorSearch creates a new vector search instance
func NewVectorSearch(client *mongo.Client) *VectorSearch {
	return &VectorSearch{client: client}
}

// Search performs vector search on a collection
func (vs *VectorSearch) Search(ctx context.Context, database, collection, embeddingField, indexName string, queryEmbedding []float32, limit int) ([]map[string]interface{}, error) {
	coll := vs.client.Database(database).Collection(collection)

	// Convert []float32 to []float64 for MongoDB
	queryVector := make([]float64, len(queryEmbedding))
	for i, v := range queryEmbedding {
		queryVector[i] = float64(v)
	}

	// Build the $vectorSearch aggregation pipeline
	pipeline := []bson.M{
		{
			"$vectorSearch": bson.M{
				"index":         indexName,
				"path":          embeddingField,
				"queryVector":   queryVector,
				"numCandidates": limit * 10, // MongoDB recommends numCandidates >= limit * 10
				"limit":         limit,
			},
		},
		{
			"$addFields": bson.M{
				"score": bson.M{"$meta": "vectorSearchScore"},
			},
		},
		{
			"$project": bson.M{
				"embedding": 0, // Exclude embedding field to reduce response size
			},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
