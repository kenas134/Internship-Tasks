package data

import (
	"context"
	"errors"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var collection *mongo.Collection

func InitCollection(col *mongo.Collection) {
	collection = col
}

// Get all tasks
func GetTasks() ([]models.Task, error) {
	var tasks []models.Task

	ctx := context.TODO()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Get one task
func GetTaskByID(id string) (models.Task, error) {
	var task models.Task

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, err
	}

	err = collection.FindOne(
		context.TODO(),
		bson.M{"_id": objectID},
	).Decode(&task)

	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

// Create task
func CreateTask(task models.Task) (models.Task, error) {
	result, err := collection.InsertOne(
		context.TODO(),
		task,
	)

	if err != nil {
		return models.Task{}, err
	}

	// MongoDB generated the ObjectID.
	task.ID = result.InsertedID.(bson.ObjectID)

	return task, nil
}

// Update task
func UpdateTaskByID(id string, updatedTask models.Task) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objectID,
	}

	update := bson.M{
		"$set": bson.M{
			"title":       updatedTask.Title,
			"description": updatedTask.Description,
			"dueDate":     updatedTask.DueDate,
			"status":      updatedTask.Status,
		},
	}

	result, err := collection.UpdateOne(
		context.TODO(),
		filter,
		update,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

// Delete task
func DeleteTaskByID(id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objectID,
	}

	result, err := collection.DeleteOne(
		context.TODO(),
		filter,
	)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

// Seed database with initial tasks
func SeedTasks() error {
	ctx := context.TODO()

	// Check whether the collection already contains tasks.
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	// Don't insert seed data if tasks already exist.
	if count > 0 {
		return nil
	}

	tasks := []models.Task{
		{
			Title:       "Learn Go",
			Description: "Learn Go fundamentals",
			Status:      "pending",
		},
		{
			Title:       "Learn MongoDB",
			Description: "Learn MongoDB CRUD operations",
			Status:      "pending",
		},
		{
			Title:       "Build REST API",
			Description: "Build Task Management REST API",
			Status:      "in-progress",
		},
		{
			Title:       "Test API",
			Description: "Test all endpoints using Postman",
			Status:      "pending",
		},
		{
			Title:       "Write Documentation",
			Description: "Document the Task Management API",
			Status:      "completed",
		},
	}

	_, err = collection.InsertMany(ctx, tasks)

	return err
}
