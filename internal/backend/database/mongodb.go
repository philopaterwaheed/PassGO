package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	client   *mongo.Client
	database *mongo.Database
)

// Connect establishes a connection to MongoDB
func Connect(uri, dbName string) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	var err error
	client, err = mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	database = client.Database(dbName)
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return nil
}

// Disconnect closes the MongoDB connection
func Disconnect() error {
	if client == nil {
		return nil
	}

	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	log.Println("Successfully disconnected from MongoDB")
	return nil
}

// GetDatabase returns the database instance
func GetDatabase() *mongo.Database {
	return database
}

// GetClient returns the MongoDB client instance
func GetClient() *mongo.Client {
	return client
}

// GetCollection returns a collection from the database
func GetCollection(collectionName string) *mongo.Collection {
	if database == nil {
		log.Fatal("Database not initialized. Call Connect() first.")
	}
	return database.Collection(collectionName)
}

// HealthCheck verifies the database connection is alive
func HealthCheck() error {
	if client == nil {
		return fmt.Errorf("database client not initialized")
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}
