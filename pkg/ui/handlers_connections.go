package ui

import (
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/db"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

func selectConnextion(g *gocui.Gui, v *gocui.View) error {
	var err error
	_, cy := v.Cursor()
	lines := strings.Split(v.Buffer(), "\n")

	selected := ""

	if cy >= 0 && cy < len(lines)-1 {
		selected = lines[cy]
	}
	model.State.SelectedDB = model.State.Config[selected].Database

	model.State.DBclient, err = db.Connect(model.State.Config[selected])
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	model.State.DBnames, err = db.GetDBs(model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to get DBs: %v", err)
	}

	updateDatabases(g)

	return nil
}
