package main

import (
	"context"
	"log"
	"os"
	"testing_task-manager_api/Delivery/routers"
	domain "testing_task-manager_api/Domain"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main(){

	err := godotenv.Load("../.env")
	if err != nil{
		log.Fatal("Error loading enviromental variables")
	}

	mongo_url := os.Getenv("MONGODB_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	client := DBinstance(mongo_url)
	db := client.Database(domain.DatabaseName)
	defer CloseMongoDBConnection(client)

	timeout := time.Duration(10) * time.Second

	gin := gin.Default()

	routers.Setup(timeout, db, gin)

	gin.Run(port)
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