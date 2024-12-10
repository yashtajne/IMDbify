package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB Client and Context
var Client *mongo.Client
var Ctx context.Context

// ConnectToDatabase connects to MongoDB and initializes the global client and context variables
func ConnectToDatabase() error {

	// Attempt to load the .env file
	if err := godotenv.Load(); err != nil {
		// Check if the URI environment variable is not empty
		uri := os.Getenv("URI")
		if uri == "" {
			log.Fatalf("error loading .env file and URI environment variable is not set: %v", err)
		}
		// Log a warning and continue
		log.Println("Warning: .env file not found, but URI environment variable is set. Continuing...")
	}

	// Get the URI from the environment variable
	uri := os.Getenv("URI")
	if uri == "" {
		return fmt.Errorf("URI environment variable is not set")
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Create a context with timeout
	Ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	var connectionError error
	Client, connectionError = mongo.Connect(Ctx, clientOptions)
	if connectionError != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", connectionError)
	}

	// Ping the MongoDB server to verify connection
	if err := Client.Ping(Ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	return nil
}

func GetCollection() *mongo.Collection {
	return Client.Database("imdb").Collection("titles")
}
