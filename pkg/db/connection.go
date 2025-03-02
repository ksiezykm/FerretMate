package db

import (
	"context"
	"fmt"

	"github.com/ksiezykm/FerretMate/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connectToDB connects to FerretDB using DatabaseConfig fields
func ConnectToDB(dbConfig config.DatabaseConfig) (*mongo.Client, error) {

	uri := "mongodb://" + dbConfig.Username + ":" + dbConfig.Password + "@" + dbConfig.Host + "/" + dbConfig.Database

	fmt.Println(uri)
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
