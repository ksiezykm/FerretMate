package db

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ListDatabases(client *mongo.Client) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ListCollections(client *mongo.Client, dbName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(dbName)
	result, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ListDocuments(client *mongo.Client, dbName, collName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection(collName)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []string
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		jsonBytes, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			continue
		}
		docs = append(docs, string(jsonBytes))
	}
	return docs, nil
}
