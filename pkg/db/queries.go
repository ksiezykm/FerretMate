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

	// Opcje dla zapytania (możesz dostosować)
	findOptions := options.Find()
	findOptions.SetProjection(bson.D{{Key: "_id", Value: 1}}) // Pobieramy tylko pole _id

	// Wykonanie zapytania
	cursor, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("błąd podczas wykonywania zapytania: %w", err)
	}
	defer cursor.Close(context.TODO())

	// Przetwarzanie wyników
	var documentIDs []string
	for cursor.Next(context.TODO()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Println("błąd podczas dekodowania dokumentu:", err)
			continue // Przechodzimy do następnego dokumentu w przypadku błędu
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
			log.Println("błąd: _id nie znalezione")
		}
	}
	
	// Sprawdzenie błędów kursora
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("błąd kursora: %w", err)
	}

	return documentIDs, nil
}

func GetDocumentByID(dbName, collectionName, documentID string) (bson.M, error) {
	collection := model.State.DBclient.Database(dbName).Collection(collectionName)

	// Próba konwersji identyfikatora na bson.ObjectID
	objID, err := primitive.ObjectIDFromHex(documentID)
	filter := bson.M{"_id": documentID} // Domyślnie filtr po stringu
	if err == nil {
			filter = bson.M{"_id": objID} // Filtr po bson.ObjectID
	}

	// Wykonanie zapytania
	var document bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&document)
	if err != nil {
			return nil, fmt.Errorf("błąd podczas pobierania dokumentu: %w", err)
	}

	return document, nil
}
