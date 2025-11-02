package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type Document struct {
	ID      interface{}
	JSON    string
	Summary string
}

func ListDocuments(client *mongo.Client, dbName, collName string) ([]Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection(collName)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []Document
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		jsonBytes, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			continue
		}

		docID := doc["_id"]
		summary := ""

		// Try to create a meaningful summary
		if name, ok := doc["name"].(string); ok && name != "" {
			summary = name
		} else if username, ok := doc["username"].(string); ok && username != "" {
			summary = username
		} else if title, ok := doc["title"].(string); ok && title != "" {
			summary = title
		} else if email, ok := doc["email"].(string); ok && email != "" {
			summary = email
		} else if idStr, ok := docID.(string); ok {
			summary = idStr
		} else {
			// Fallback: show _id type and first few fields
			summary = fmt.Sprintf("%v", docID)
			if len(summary) > 50 {
				summary = summary[:50] + "..."
			}
		}

		docs = append(docs, Document{
			ID:      docID,
			JSON:    string(jsonBytes),
			Summary: summary,
		})
	}
	return docs, nil
}

func GetDocument(client *mongo.Client, dbName, collName string, docID interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection(collName)
	var doc bson.M
	err := coll.FindOne(ctx, bson.M{"_id": docID}).Decode(&doc)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func UpdateDocument(client *mongo.Client, dbName, collName string, docJSON string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var doc bson.M
	if err := json.Unmarshal([]byte(docJSON), &doc); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	rawID, ok := doc["_id"]
	if !ok {
		return fmt.Errorf("document has no _id field")
	}

	// Convert _id to proper type if it's an ObjectID string
	var docID interface{}
	if idMap, ok := rawID.(map[string]interface{}); ok {
		if oidStr, ok := idMap["$oid"].(string); ok {
			// Handle {"$oid": "..."} format
			oid, err := primitive.ObjectIDFromHex(oidStr)
			if err != nil {
				docID = rawID
			} else {
				docID = oid
			}
		} else {
			docID = rawID
		}
	} else if idStr, ok := rawID.(string); ok {
		// Handle plain string format
		oid, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			docID = rawID
		} else {
			docID = oid
		}
	} else {
		docID = rawID
	}

	delete(doc, "_id")

	coll := client.Database(dbName).Collection(collName)
	result, err := coll.ReplaceOne(ctx, bson.M{"_id": docID}, doc)
	if err != nil {
		return fmt.Errorf("failed to replace document: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with _id: %v", docID)
	}

	return nil
}

// CreateDatabase creates a new database by creating an initial collection
// MongoDB requires at least one collection for a database to exist
func CreateDatabase(client *mongo.Client, dbName, collName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if client == nil {
		return fmt.Errorf("database client is nil")
	}

	if dbName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	if collName == "" {
		return fmt.Errorf("collection name cannot be empty")
	}

	// Create the collection (this will also create the database)
	err := client.Database(dbName).CreateCollection(ctx, collName)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

// CreateCollection creates a new collection in the specified database
func CreateCollection(client *mongo.Client, dbName, collName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Database(dbName).CreateCollection(ctx, collName)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	return nil
}

// CreateDocument creates a new document with a new ObjectID
func CreateDocument(client *mongo.Client, dbName, collName, docJSON string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var doc bson.M
	if err := json.Unmarshal([]byte(docJSON), &doc); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Generate new ObjectID if _id is not provided or is in template format
	if rawID, ok := doc["_id"]; ok {
		if idMap, ok := rawID.(map[string]interface{}); ok {
			if _, hasOid := idMap["$oid"]; hasOid {
				// Replace template ObjectID with new one
				doc["_id"] = primitive.NewObjectID()
			}
		}
	} else {
		// No _id provided, generate one
		doc["_id"] = primitive.NewObjectID()
	}

	coll := client.Database(dbName).Collection(collName)
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	return nil
}

// DeleteDocument deletes a document from a collection
func DeleteDocument(client *mongo.Client, dbName, collName string, docID interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection(collName)
	result, err := coll.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found with _id: %v", docID)
	}

	return nil
}

// DeleteCollection deletes a collection from a database
func DeleteCollection(client *mongo.Client, dbName, collName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection(collName)
	err := coll.Drop(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	return nil
}

// DeleteDatabase deletes a database
func DeleteDatabase(client *mongo.Client, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Database(dbName).Drop(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete database: %w", err)
	}

	return nil
}
