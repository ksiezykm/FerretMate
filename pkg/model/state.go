package model

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type AppState struct {
	DBname             string
	DBclient           *mongo.Client
	Collections        []string
	Documents          []string
	DocumentContent    string
	SelectedCollection string
	SelectedDocument   string
	Messages           string
}

var State AppState

func InitAppState() {
	State = AppState{
		Collections: []string{},
		Documents:   []string{},
	}
}
