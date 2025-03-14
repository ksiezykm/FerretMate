package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetCollections retrieves the list of collections from the given database.
func GetCollections(dbName string, client *mongo.Client) ([]string, error) {
	collections, err := client.Database(dbName).ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	return collections, nil

}

// GetDocuments retrieves documents from the specified collection.
func GetDocuments(dbName, collectionName string, client *mongo.Client) ([]string, error) {
	collection := client.Database(dbName).Collection(collectionName)

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

func GetDocumentByID(dbName, collectionName, documentID string, client *mongo.Client) (bson.M, error) {
	collection := client.Database(dbName).Collection(collectionName)

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

func UpdateDocumentByID(dbName, collectionName, documentID, documentContent string, client *mongo.Client) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Attempt to convert the identifier to bson.ObjectID
	objID, err := primitive.ObjectIDFromHex(documentID)
	filter := bson.M{"_id": documentID} // Default filter by string
	if err == nil {
		filter = bson.M{"_id": objID} // Filter by bson.ObjectID
	}

	var data map[string]interface{}

	// JSON-a to map.
	err = json.Unmarshal([]byte(documentContent), &data)
	if err != nil {
		return fmt.Errorf("błąd podczas dekodowania JSON: %w", err)
	}

	// delete "_id" from map, if exist.
	delete(data, "_id")

	// map to JSON-a.
	modifiedJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("błąd podczas kodowania JSON: %w", err)
	}

	// Parsing JSON string to bson.M
	var update bson.M
	err = json.Unmarshal([]byte(modifiedJSON), &update)
	if err != nil {
		return fmt.Errorf("error during JSON parsing: %w", err)
	}

	// Performing the update
	updateResult, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("error during document update: %w", err)
	}

	if updateResult.MatchedCount == 0 {
		return fmt.Errorf("document with ID %s not found", documentID)
	}

	return nil
}

// DeleteDocumentByID deletes a document from the specified collection by its ID.
func DeleteDocumentByID(dbName, collectionName, documentID string, client *mongo.Client) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Attempt to convert the identifier to bson.ObjectID
	objID, err := primitive.ObjectIDFromHex(documentID)
	filter := bson.M{"_id": documentID} // Default filter by string
	if err == nil {
		filter = bson.M{"_id": objID} // Filter by bson.ObjectID
	}

	// Execute the delete operation
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}

	// Check if any document was deleted
	if result.DeletedCount == 0 {
		return fmt.Errorf("document with ID %s not found", documentID)
	}

	return nil
}

// CreateDocument creates a new document with basic fields in the specified collection.
func CreateDocument(dbName, collectionName string, client *mongo.Client) (error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Create a basic document
	document := bson.M{
		"new":    "document",
		"status": "pending",
	}

	// Execute the insert operation
	_, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return fmt.Errorf("error creating document: %w", err)
	}

	return nil
}
