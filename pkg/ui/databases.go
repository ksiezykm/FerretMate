package ui

import (
	"fmt"
	"math/rand"

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
	if err := g.SetKeybinding("databases", gocui.KeyDelete, gocui.ModNone, deleteDatabase); err != nil {
		return err
	}
	return nil
}

func setCurrentViewDatabases(g *gocui.Gui, v *gocui.View) error {
	return SetCurrentView(g, v, "databases", "Enter: view | Delete: delete | Ctrl+n: new")
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
	selected, err := GetSelectedLine(v)
	if err != nil {
		// Nothing selected or view is empty
		return nil
	}

	model.State.SelectedDB = selected
	model.State.Collections = nil

	// The client is already connected from the connections view.
	// No need to reconnect here.

	model.State.Collections, err = db.GetCollections(model.State.SelectedDB, model.State.DBclient)
	if err != nil {
		DisplayError(g, err)
		return nil
	}
	updateCollections(g)
	model.State.DocumentContent = ""
	updateDocumentContent(g)
	model.State.Documents = nil
	updateDocuments(g)

	return nil
}

func createNewDatabase(g *gocui.Gui, v *gocui.View) error {

	random := rand.Intn(100) + 1

	lineToEdit = "new_db" + fmt.Sprint(random)
	mode = "createDB"

	if err := editView(g); err != nil {
		return err
	}

	return nil
}

func deleteDatabase(g *gocui.Gui, v *gocui.View) error {
	selected, err := GetSelectedLine(v)
	if err != nil {
		// Nothing selected
		return nil
	}

	if err := db.DeleteDatabase(model.State.DBclient, selected); err != nil {
		DisplayError(g, err)
		return nil
	}

	model.State.DBnames, err = db.GetDBs(model.State.DBclient)
	if err != nil {
		DisplayError(g, err)
		return nil
	}

	updateDatabases(g)
	return nil
}
