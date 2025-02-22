package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectToDB establishes a connection to the FerretDB instance.
func ConnectToDB() (*mongo.Client, error) {
	uri := os.Getenv("DB_URI")
	if uri == "" {
		return nil, ErrMissingEnv
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// ErrMissingEnv is returned when the DB_URI environment variable is missing.
var ErrMissingEnv = &EnvError{"DB_URI not set"}

type EnvError struct{ msg string }

func (e *EnvError) Error() string { return e.msg }
