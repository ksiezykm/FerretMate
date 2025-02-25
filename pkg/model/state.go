package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppState struct {
	DBname             string
	DBclient           *mongo.Client
	Collections        []string
	Documents          []string
	DocumentDetails    bson.M
	SelectedCollection string
	SelectedDocument   string
	Messages		   string
}

var State AppState

func InitAppState() {
	State = AppState{
		Collections: []string{},
		Documents:   []string{},
	}
}
