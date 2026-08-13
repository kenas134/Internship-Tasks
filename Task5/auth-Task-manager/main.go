package main

import (
	"context"
	"log"
	"os"
	"time"

	"auth-Task-Manager/controllers"
	"auth-Task-Manager/data"
	"auth-Task-Manager/router"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {

	// Load .env
	err := godotenv.Load()

	if err != nil {
		log.Println(
			"No .env file found. Using environment variables.",
		)
	}

	// Get MongoDB URI.
	mongoURI := os.Getenv("MONGODB_URI")

	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	// Get database name.
	databaseName := os.Getenv("DATABASE_NAME")

	if databaseName == "" {
		databaseName = "task_manager"
	}

	// MongoDB client options.
	serverAPI := options.ServerAPI(
		options.ServerAPIVersion1,
	)

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(serverAPI)

	// Create MongoDB client.
	client, err := mongo.Connect(
		clientOptions,
	)

	if err != nil {
		log.Fatal(
			"Failed to create MongoDB client:",
			err,
		)
	}

	// Disconnect when application exits.
	defer func() {

		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Println(
				"Error disconnecting MongoDB:",
				err,
			)
		}

	}()

	// Ping MongoDB.
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

	// Get database.
	db := client.Database(databaseName)

	// Create services.
	userService := data.NewUserService(db)
	taskService := data.NewTaskService(db)

	// Create controller.
	controller := controllers.NewController(
		userService,
		taskService,
	)

	// Create router.
	r := router.SetupRouter(controller)

	log.Println(
		"Server running on http://localhost:8080",
	)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}