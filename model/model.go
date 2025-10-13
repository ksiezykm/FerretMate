package model

type Model struct {
	SelectedListView string

	Connections        []string
	SelectedConnection string

	DBs        []string
	SelectedDB string

	Collections        []string
	SelectedCollection string

	Documents        []string
	SelectedDocument string

	// Mockup documents - map of document name to JSON content
	DocumentContent map[string]string
}
