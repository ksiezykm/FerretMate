package db

import (
	"context"

	"github.com/ksiezykm/FerretMate/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connectToDB connects to FerretDB using the URI from .env file
func ConnectToDB() (*mongo.Client, error) {
	uri, err := config.ReadDBURI(".env")
	if err != nil {
		return nil, err
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client, nil
}
