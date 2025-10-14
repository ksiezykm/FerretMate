package model

type Model struct {
	SelectedListView string

	LoadedConnections  []Connection
	Connections        []string
	SelectedConnection string

	DBs        []string
	SelectedDB string

	Collections        []string
	SelectedCollection string

	Documents        []string
	SelectedDocument string

	DocumentContent map[string]string
}
