package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Task struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	DueDate     time.Time     `bson:"dueDate" json:"dueDate"`
	Status      string        `bson:"status" json:"status"`
}
