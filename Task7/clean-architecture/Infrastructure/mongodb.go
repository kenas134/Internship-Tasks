package infrastructure

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongoDB() *mongo.Client {
	mongoURI := os.Getenv("MONGODB_URI")

	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	databaseName := os.Getenv("DATABASE_NAME")

	if databaseName == "" {
		databaseName = "task_manager"
	}

	serverAPI := options.ServerAPI(
		options.ServerAPIVersion1,
	)

	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		log.Fatal(
			"Failed to create MongoDB client:",
			err,
		)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Println(
				"Error disconnecting MongoDB:",
				err,
			)
		}
	}()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(
			"Failed to connect to MongoDB:",
			err,
		)
	}

	log.Println("Connected to MongoDB")
	return client
}
