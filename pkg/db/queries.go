package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetCollections retrieves the list of collections from the given database.
func GetCollections(client *mongo.Client, dbName string) ([]string, error) {
	collections, err := client.Database(dbName).ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	return collections, nil
}

// GetDocuments retrieves documents from the specified collection.
func GetDocuments(client *mongo.Client, dbName, collectionName string) ([]bson.M, error) {
	collection := client.Database(dbName).Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
