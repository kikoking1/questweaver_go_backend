package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectTest() {

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDBName := os.Getenv("MONGODB_DATABASE")

	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable must be set")
	}
	if mongoDBName == "" {
		mongoDBName = "questweaverpro" // default database name
	}

	opts := options.Client().ApplyURI(mongoURI)

	// 2. Connect to MongoDB
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Verify connection with a ping (Context with timeout is best practice)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	fmt.Println("Successfully connected to MongoDB locally!")

	// Example: Access a collection
	// collection := client.Database(mongoDBName).Collection("users")
	// _ = collection // Ready for CRUD operations

	// // Remember to disconnect when the app closes
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
}
