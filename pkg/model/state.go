package model

type AppState struct {
	Collections []string
	Documents   []string
}

var State AppState

func InitAppState() {
	State = AppState{
		Collections: []string{},
		Documents:   []string{},
	}
}
