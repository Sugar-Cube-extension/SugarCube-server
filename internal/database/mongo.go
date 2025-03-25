package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"os"
	"time"
)

var Client *mongo.Client

func ConnectDB() {
	uri := os.Getenv("MONGO_URI") // edit
	if uri == "" {
		uri = "mongodb://localhost:27017" // edit
	}

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB not responding:", err)
	}

	fmt.Println("Connected to MongoDB")
	Client = client
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("sugarcube").Collection(collectionName)
}
