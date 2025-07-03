package database

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

func SetMongoClient(client *mongo.Client) {
	MongoClient = client
}

func GetCollection(collectionName string) *mongo.Collection {
	database := os.Getenv("MONGODB")
	if database == "" {
		database = "shop_dashboard"
	}
	return MongoClient.Database(database).Collection(collectionName)
}
