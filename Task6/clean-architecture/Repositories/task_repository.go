package repositories

import (
	"context"
	"errors"
	"time"

	domain "clean-architecture/Domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TaskRepository interface {
	Create(task domain.Task) error
	GetTaskByID(id string) (domain.Task, error)
	GetTasks() ([]domain.Task, error)
	Update(task domain.Task) error
	Delete(id string) error
}

type TaskDocument struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	DueDate     time.Time     `bson:"due_date"`
	Status      string        `bson:"status"`
}



type TaskRepositoryImpl struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database)*TaskRepositoryImpl{
	return &TaskRepositoryImpl{
		collection: db.Collection("task"),
	}
}


func (service *TaskRepositoryImpl) GetTasks() ([]domain.Task, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
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

	var tasks []TaskDocument

	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	domainTasks := make([]domain.Task, 0, len(tasks))

	for _, task := range tasks {
		domainTasks = append(domainTasks, domain.Task{
			ID:          task.ID.Hex(),
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Status:      task.Status,
		})
	}

	return domainTasks, nil
}


func (service *TaskRepositoryImpl) GetTaskByID(id string) (domain.Task,error) {
	

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	var task TaskDocument

	err = service.collection.FindOne(
		ctx,
		bson.M{"_id": objectID},
	).Decode(&task)

	if  err != nil {
		return domain.Task{}, err
	}


	res := domain.Task{
			ID:          task.ID.Hex(),
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Status:      task.Status,
		}

	return res, nil
}


func (service *TaskRepositoryImpl) Create(task domain.Task)  error {
	newTask := TaskDocument{
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Status:      task.Status,
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	_, err := service.collection.InsertOne(
		ctx,
		newTask,
	)

	if err != nil {
		return err
	}

	return nil
}


func (service *TaskRepositoryImpl) Update(task domain.Task) error {

    objectID, err := bson.ObjectIDFromHex(task.ID)
    if err != nil {
        return err
    }

    filter := bson.M{
        "_id": objectID,
    }

    update := bson.M{
        "$set": bson.M{
            "title":       task.Title,
            "description": task.Description,
            "due_date":    task.DueDate,
            "status":      task.Status,
        },
    }

    ctx, cancel := context.WithTimeout(
        context.Background(),
        5*time.Second,
    )
    defer cancel()

    result, err := service.collection.UpdateOne(
        ctx,
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


func (service *TaskRepositoryImpl) Delete(id string) error {
	objectID,err := bson.ObjectIDFromHex(id)

	if err != nil {
		return err
	}


	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	result, err := service.collection.DeleteOne(
		ctx,
		bson.M{
			"_id": objectID,
		},
	)

	if  err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New(
			"task not found",
		)
	}

	return nil
}