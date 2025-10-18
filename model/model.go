package model

type Model struct {
	SelectedListView string

	LoadedConnections       []Connection
	Connections             []string
	SelectedConnection      string
	SelectedConnectionIndex int

	DBs             []string
	SelectedDB      string
	SelectedDBIndex int

	Collections             []string
	SelectedCollection      string
	SelectedCollectionIndex int

	Documents             []string
	SelectedDocument      string
	SelectedDocumentIndex int

	DocumentContent map[string]string
}
