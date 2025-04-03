package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MyClient *mongo.Client

func init() {
	ConnectDB()
}

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(GetMongoURI()))
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("MongoDB ping failed: ", err)
	}

	fmt.Println("Connected to MongoDB!")
	MyClient = client
}

func GetCollection(collectionName string) *mongo.Collection {
	if MyClient == nil {
		log.Fatal("Mongo client is not initialized")
	}
	return MyClient.Database("cards").Collection(collectionName)
}
