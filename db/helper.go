package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbSetup() (client *mongo.Client, err error) {
	connectionString := "mongodb+srv://" + os.Getenv("MONGO_USERNAME") + ":" + os.Getenv("MONGO_PASSWORD") + "@cluster1.5jqwhvz.mongodb.net/?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(connectionString)

	// Establish a connection to the MongoDB server.
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
