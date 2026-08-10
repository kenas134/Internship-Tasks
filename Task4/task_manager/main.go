package main

import (
	"context"
	"log"
	"os"
	"time"

	"task_manager/data"
	"task_manager/database"
	"task_manager/router"

	"github.com/joho/godotenv"
)

func main() {

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	// Get MongoDB URI
	URI := os.Getenv("MONGODB_URI")

	if URI == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	// Connect to MongoDB
	client, db, err := database.Connect(URI)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	log.Println("Connected to MongoDB")

	// Disconnect MongoDB when application stops
	defer func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting MongoDB:", err)
		}

		log.Println("Disconnected from MongoDB")
	}()

	// Get tasks collection
	collection := db.Collection("tasks")

	// Give collection to data package
	data.InitCollection(collection)

	// Insert seed data if database is empty
	if err := data.SeedTasks(); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	log.Println("Database initialized")

	// Create Gin router
	r := router.SetupRouter()

	// Start server
	log.Println("Server running on :8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
