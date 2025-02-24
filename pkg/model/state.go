package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var DBclient *mongo.Client

type AppState struct {
	Collections     []string
	Documents       []string
	DocumentDetails bson.M
}

var State AppState

func InitAppState() {
	State = AppState{
		Collections: []string{},
		Documents:   []string{},
	}
}
