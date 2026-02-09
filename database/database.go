package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Database

// Connect initializes the MongoDB connection and stores the database reference.
func Connect() {
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDBName := os.Getenv("MONGODB_DATABASE")

	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable must be set")
	}
	if mongoDBName == "" {
		mongoDBName = "questweaverpro"
	}

	opts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal("Failed to create MongoDB client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	DB = client.Database(mongoDBName)
	log.Println("Successfully connected to MongoDB")
}

// GetCollection returns a handle to the named collection.
func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
