package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client wraps MongoDB client
type Client struct {
	client *mongo.Client
}

// NewClient creates a new MongoDB client
func NewClient(ctx context.Context, uri string) (*Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

// Database returns a database instance
func (c *Client) Database(name string) *mongo.Database {
	return c.client.Database(name)
}

// Collection returns a collection instance from a database
func (c *Client) Collection(databaseName, collectionName string) *mongo.Collection {
	return c.client.Database(databaseName).Collection(collectionName)
}

// Client returns the underlying MongoDB client
func (c *Client) Client() *mongo.Client {
	return c.client
}

// Close closes the MongoDB connection
func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}
