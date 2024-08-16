package config

import (
	"log"
	"os"
	domain "testing_task-manager_api/Domain"

	"github.com/joho/godotenv"
)

func GetConfig() *domain.Config {
	_, err := os.Stat(".env")

	if !os.IsNotExist(err) {
		err := godotenv.Load(".env")

		if err != nil {
			log.Println("Error while reading the env file", err)
			panic(err)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080" // default port if not set in environment
	}

	config := &domain.Config{
		MongoDBURI: "mongodb://localhost:27017/taskmanager",
		Port:       port,
		TimeZone:   "Asia/Jakarta",
		SecretKey:  os.Getenv("SECRET_KEY"),
		DatabaseName: "test_db",
	}

	return config
}
