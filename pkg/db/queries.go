package db

import (
	"context"
	"fmt"
	"log"

	"github.com/ksiezykm/FerretMate/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetCollections retrieves the list of collections from the given database.
func GetCollections(dbName string) ([]string, error) {
	collections, err := model.State.DBclient.Database(dbName).ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	return collections, nil
}

// GetDocuments retrieves documents from the specified collection.
func GetDocuments(dbName, collectionName string) ([]string, error) {
	collection := model.State.DBclient.Database(dbName).Collection(collectionName)

	// Query options (you can customize)
	findOptions := options.Find()
	findOptions.SetProjection(bson.D{{Key: "_id", Value: 1}}) // Retrieve only the _id field

	// Execute the query
	cursor, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer cursor.Close(context.TODO())

	// Process results
	var documentIDs []string
	for cursor.Next(context.TODO()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Println("error decoding document:", err)
			continue // Proceed to the next document in case of error
		}
		if id, ok := doc["_id"]; ok {
			switch id.(type) {
			case string:
			case primitive.ObjectID:
				documentIDs = append(documentIDs, id.(primitive.ObjectID).Hex())
			default:
				documentIDs = append(documentIDs, fmt.Sprintf("%v", id))
			}
		} else {
			log.Println("error: _id not found")
		}
	}

	// Check cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return documentIDs, nil
}

func GetDocumentByID(dbName, collectionName, documentID string) (bson.M, error) {
	collection := model.State.DBclient.Database(dbName).Collection(collectionName)

	// Attempt to convert the identifier to bson.ObjectID
	objID, err := primitive.ObjectIDFromHex(documentID)
	filter := bson.M{"_id": documentID} // Default filter by string
	if err == nil {
		filter = bson.M{"_id": objID} // Filter by bson.ObjectID
	}

	// Execute the query
	var document bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&document)
	if err != nil {
		return nil, fmt.Errorf("error retrieving document: %w", err)
	}

	return document, nil
}
