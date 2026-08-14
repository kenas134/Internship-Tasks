package data

import (
	"context"
	"errors"
	"time"

	"auth-Task-Manager/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)




type TaskService struct {
	collection *mongo.Collection
}

func NewTaskService(db *mongo.Database) *TaskService {
	return &TaskService{
		collection: db.Collection("tasks"),
	}
}

// GetTasks gets all tasks.
func (service *TaskService) GetTasks() ([]models.Task, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	cursor, err := service.collection.Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var tasks []models.Task

	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTaskByID gets one task.
func (service *TaskService) GetTaskByID(
	id string,
) (models.Task, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	var task models.Task

	err := service.collection.FindOne(
		ctx,
		bson.M{
			"_id": id,
		},
	).Decode(&task)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Task{}, errors.New(
				"task not found",
			)
		}

		return models.Task{}, err
	}

	return task, nil
}

// CreateTask creates a task.
func (service *TaskService) CreateTask(
	task models.Task,
) (models.Task, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	task.ID = bson.NewObjectID().Hex()

	_, err := service.collection.InsertOne(
		ctx,
		task,
	)

	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

// UpdateTaskByID updates a task.
func (service *TaskService) UpdateTaskByID(
	id string,
	task models.Task,
) (models.Task, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	task.ID = id

	result, err := service.collection.ReplaceOne(
		ctx,
		bson.M{
			"_id": id,
		},
		task,
	)

	if err != nil {
		return models.Task{}, err
	}

	if result.MatchedCount == 0 {
		return models.Task{}, errors.New(
			"task not found",
		)
	}

	return task, nil
}

// DeleteTaskByID deletes a task.
func (service *TaskService) DeleteTaskByID(
	id string,
) error {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	result, err := service.collection.DeleteOne(
		ctx,
		bson.M{
			"_id": id,
		},
	)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New(
			"task not found",
		)
	}

	return nil
}