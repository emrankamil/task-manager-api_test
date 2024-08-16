package config

import (
	"context"
	"log"
	domain "testing_task-manager_api/Domain"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(config *domain.Config) (*mongo.Client, *mongo.Database){
	client := DBinstance(config.MongoDBURI)
	db := client.Database(config.DatabaseName)

	return client, db
}

func DBinstance(mongodb string) *mongo.Client{

	clientOptions := options.Client().ApplyURI(mongodb)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CloseMongoDBConnection(client *mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}