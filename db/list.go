package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

		// Display _id as summary
		if idStr, ok := docID.(string); ok {
			summary = idStr
		} else {
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

// ExportDocument exports a single document to a JSON file
func ExportDocument(client *mongo.Client, dbName, collName string, docID interface{}, filePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection(collName)
	var doc bson.M
	err := coll.FindOne(ctx, bson.M{"_id": docID}).Decode(&doc)
	if err != nil {
		return fmt.Errorf("failed to find document: %w", err)
	}

	jsonBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filePath[:strings.LastIndex(filePath, "/")]
	if dir != "" && dir != filePath {
		if err := createDirIfNotExists(dir); err != nil {
			return err
		}
	}

	// Write to file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(jsonBytes)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// ExportCollection exports all documents from a collection to a directory
func ExportCollection(client *mongo.Client, dbName, collName, dirPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create directory
	if err := createDirIfNotExists(dirPath); err != nil {
		return err
	}

	coll := client.Database(dbName).Collection(collName)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		// Generate filename based on _id or index
		var filename string
		if id, ok := doc["_id"]; ok {
			filename = fmt.Sprintf("%v.json", id)
			// Clean filename from invalid characters
			filename = strings.ReplaceAll(filename, "/", "_")
			filename = strings.ReplaceAll(filename, "\\", "_")
			filename = strings.ReplaceAll(filename, ":", "_")
		} else {
			filename = fmt.Sprintf("doc_%d.json", count)
		}

		jsonBytes, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", dirPath, filename)
		file, err := os.Create(filePath)
		if err != nil {
			continue
		}

		file.Write(jsonBytes)
		file.Close()
		count++
	}

	if count == 0 {
		return fmt.Errorf("no documents exported")
	}

	return nil
}

// ExportDatabase exports all collections from a database to a directory structure
func ExportDatabase(client *mongo.Client, dbName, dirPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create base directory
	if err := createDirIfNotExists(dirPath); err != nil {
		return err
	}

	// Get all collections
	collections, err := client.Database(dbName).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	if len(collections) == 0 {
		return fmt.Errorf("no collections found in database")
	}

	// Export each collection
	for _, collName := range collections {
		collPath := fmt.Sprintf("%s/%s", dirPath, collName)
		if err := ExportCollection(client, dbName, collName, collPath); err != nil {
			// Log error but continue with other collections
			fmt.Printf("Warning: failed to export collection %s: %v\n", collName, err)
		}
	}

	return nil
}

// createDirIfNotExists creates a directory if it doesn't exist
func createDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

// UploadDocument uploads document(s) from a JSON file to MongoDB
// Supports both single document (object) and multiple documents (array)
// If a document doesn't have an _id, MongoDB will automatically generate one
func UploadDocument(client *mongo.Client, dbName, collName, filePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Read the JSON file
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	coll := client.Database(dbName).Collection(collName)

	// Try to parse as array first
	var docArray []bson.M
	if err := json.Unmarshal(jsonBytes, &docArray); err == nil {
		// It's an array of documents
		if len(docArray) == 0 {
			return fmt.Errorf("empty document array")
		}

		// Process each document
		var docs []interface{}
		for _, doc := range docArray {
			processDocumentID(&doc)
			docs = append(docs, doc)
		}

		// Insert all documents
		_, err := coll.InsertMany(ctx, docs)
		if err != nil {
			return fmt.Errorf("failed to insert documents: %w", err)
		}

		return nil
	}

	// Not an array, try to parse as single document
	var doc bson.M
	if err := json.Unmarshal(jsonBytes, &doc); err != nil {
		return fmt.Errorf("failed to parse JSON (not a valid object or array): %w", err)
	}

	// Process the document's _id
	processDocumentID(&doc)

	// Insert the single document
	_, err = coll.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	return nil
}

// processDocumentID converts _id field to proper ObjectID format if needed
func processDocumentID(doc *bson.M) {
	if rawID, ok := (*doc)["_id"]; ok {
		if idMap, ok := rawID.(map[string]interface{}); ok {
			if oidStr, ok := idMap["$oid"].(string); ok {
				// Convert {"$oid": "..."} format to ObjectID
				oid, err := primitive.ObjectIDFromHex(oidStr)
				if err == nil {
					(*doc)["_id"] = oid
				}
			}
		} else if idStr, ok := rawID.(string); ok {
			// Try to convert plain string to ObjectID
			oid, err := primitive.ObjectIDFromHex(idStr)
			if err == nil {
				(*doc)["_id"] = oid
			}
		}
	}
	// If no _id is present, MongoDB will automatically generate one during insert
}
