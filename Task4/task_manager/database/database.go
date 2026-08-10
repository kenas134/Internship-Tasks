package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(uri string) (*mongo.Client, *mongo.Database, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, nil, err
	}

	db := client.Database("task_management_api")

	return client, db, nil
}
