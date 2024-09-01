package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


//var MONGO_URI string = "mongodb+srv://{user_name}:{password}@cluster0.yqtmf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"




// OpenCollection opens a collection from the specified database
func OpenCollection(databaseName, collectionName string) *mongo.Collection {
	client := DBInstance()

	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}


// DBInstance initializes a new MongoDB client instance
func DBInstance() *mongo.Client {

	// Get MongoDB URI from environment variables
	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		// can set some Default value like below or 
		// MONGO_URI = "mongodb://localhost:27017"
		log.Fatal("MONGO_URI not found in environment variables")
		return nil
	}

	// Create a new context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	// verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("error pinging MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

