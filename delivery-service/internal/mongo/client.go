package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"delivery-service/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewClient creates a new MongoDB connection
func NewClient(uri, database, collection string) (*Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("✓ Connected to MongoDB")

	return &Client{
		client:     client,
		collection: client.Database(database).Collection(collection),
	}, nil
}

// Close closes the MongoDB connection
func (c *Client) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	log.Println("✓ MongoDB connection closed")
	return nil
}

// SaveDeliveryResult inserts a delivery result into MongoDB
func (c *Client) SaveDeliveryResult(ctx context.Context, result models.DeliveryResult) error {
	_, err := c.collection.InsertOne(ctx, result)
	if err != nil {
		return fmt.Errorf("failed to save delivery result: %w", err)
	}
	return nil
}
