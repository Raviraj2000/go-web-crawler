package mongodriver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Raviraj2000/go-web-crawler/pkg/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDriver struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// Config represents the configuration required to initialize a MongoDriver instance.
type Config struct {
	URI        string
	Database   string
	Collection string
}

func NewMongoDriver(config Config) (models.StorageDriver, error) {
	// Set MongoDB Server API options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.URI).SetServerAPIOptions(serverAPI)

	// Create a context with timeout for connecting
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to confirm the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB successfully.")

	// Ensure the collection exists
	collection := client.Database(config.Database).Collection(config.Collection)
	log.Printf("Using collection: %s in database: %s", config.Collection, config.Database)

	return &MongoDriver{
		client:     client,
		collection: collection,
	}, nil
}

// Save inserts PageData into the MongoDB collection.
func (m *MongoDriver) Save(data models.PageData) error {
	doc := bson.M{
		"url":         data.URL,
		"title":       data.Title,
		"description": data.Description,
		"timestamp":   time.Now(),
	}

	_, err := m.collection.InsertOne(context.TODO(), doc)
	if err != nil {
		return fmt.Errorf("failed to insert document into MongoDB: %w", err)
	}

	log.Println("Document inserted successfully.")
	return nil
}

// Close closes the MongoDB connection.
func (m *MongoDriver) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("failed to close MongoDB connection: %w", err)
	}

	log.Println("MongoDB connection closed successfully.")
	return nil
}
