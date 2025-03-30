package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/db"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

func RegisterKeyBindingsDatabases(g *gocui.Gui) error {
	if err := g.SetKeybinding("databases", gocui.KeyEnter, gocui.ModNone, selectDatabase); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlU, gocui.ModNone, setCurrentViewDatabases); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyF2, gocui.ModNone, setCurrentViewDatabases); err != nil {
		return err
	}
	if err := g.SetKeybinding("databases", gocui.KeyCtrlN, gocui.ModNone, createNewDatabase); err != nil {
		return err
	}
	return nil
}

func setCurrentViewDatabases(g *gocui.Gui, v *gocui.View) error {

	v.FrameColor = gocui.ColorDefault
	v.SelFgColor = gocui.ColorDefault
	var nextView *gocui.View
	var err error
	if nextView, err = g.SetCurrentView("databases"); err != nil {
		return err
	}
	nextView.Highlight = true
	nextView.FrameColor = gocui.ColorGreen
	nextView.SelFgColor = gocui.ColorGreen
	// nextView.SetOrigin(0, 0)
	// nextView.SetCursor(0, 0)

	model.State.Messages = "Enter: view | Delete: delete | Ctrl+n: new"
	updateMessages(g)

	return nil
}

func updateDatabases(g *gocui.Gui) error {
	v, err := g.View("databases")
	if err != nil {
		return err
	}
	v.Clear()
	for _, dbname := range model.State.DBnames {
		fmt.Fprintln(v, dbname)
	}
	return nil
}

func selectDatabase(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	lines := strings.Split(v.Buffer(), "\n")

	selected := ""

	if cy >= 0 && cy < len(lines)-1 {
		selected = lines[cy]
	}
	var err error
	model.State.Collections = nil

	if selected == "" {
		return nil
	}

	model.State.SelectedDB = selected
	model.State.DBclient, err = db.Connect2(model.State.SelectedConnection, model.State.SelectedDB)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	//defer model.State.DBclient.Disconnect(context.TODO())

	model.State.Collections, err = db.GetCollections(model.State.SelectedDB, model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to retrieve collections: %v", err)
	}
	updateCollections(g)
	model.State.DocumentContent = ""
	updateDocumentContent(g)
	model.State.Documents = nil
	updateDocuments(g)

	return nil
}

func createNewDatabase(g *gocui.Gui, v *gocui.View) error {
	var err error

	err = db.CreateDatabase(model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// model.State.Documents, err = db.GetDocuments(model.State.DBname, model.State.SelectedCollection, model.State.DBclient)
	// if err != nil {
	// 	log.Fatalf("Failed to retrieve collection: %v", err)
	// }
	// model.State.DocumentContent = ""
	// updateDocumentDetails(g)
	// updateDocuments(g)

	return nil
}
