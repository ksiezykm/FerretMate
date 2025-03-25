package ui

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/db"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

func RegisterKeyBindingsDocuments(g *gocui.Gui) error {
	if err := g.SetKeybinding("documents", gocui.KeyEnter, gocui.ModNone, selectDocument); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, setCurrentViewDocuments); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyF4, gocui.ModNone, setCurrentViewDocuments); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyCtrlN, gocui.ModNone, createNewDocument); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyDelete, gocui.ModNone, deleteDocument); err != nil {
		return err
	}
	return nil
}

func setCurrentViewDocuments(g *gocui.Gui, v *gocui.View) error {

	v.FrameColor = gocui.ColorDefault
	v.SelFgColor = gocui.ColorDefault
	var nextView *gocui.View
	var err error
	if nextView, err = g.SetCurrentView("documents"); err != nil {
		return err
	}
	nextView.Highlight = true
	nextView.FrameColor = gocui.ColorGreen
	nextView.SelFgColor = gocui.ColorGreen
	// nextView.SetOrigin(0, 0)
	// nextView.SetCursor(0, 0)
	return nil
}

func selectDocument(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	lines := strings.Split(v.Buffer(), "\n")

	selected := ""

	if cy >= 0 && cy < len(lines)-1 {
		selected = lines[cy]
	}
	var err error

	model.State.SelectedDocument = selected
	documentDromDB, err := db.GetDocumentByID(model.State.SelectedDB, model.State.SelectedCollection, selected, model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to retrieve Document: %v", err)
	}
	jsonDoc, err := json.MarshalIndent(documentDromDB, "", "  ")
	if err != nil {
		log.Println("Error converting document to JSON:", err)
		return nil
	}
	model.State.DocumentContent = string(jsonDoc)
	updateDocumentContent(g)
	return nil
}

func createNewDocument(g *gocui.Gui, v *gocui.View) error {
	var err error

	err = db.CreateDocument(model.State.SelectedDB, model.State.SelectedCollection, model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	model.State.Documents, err = db.GetDocuments(model.State.SelectedDB, model.State.SelectedCollection, model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to retrieve collection: %v", err)
	}
	model.State.DocumentContent = ""
	updateDocumentContent(g)
	updateDocuments(g)

	return nil
}

func deleteDocument(g *gocui.Gui, v *gocui.View) error {
	var err error
	_, cy := v.Cursor()
	lines := strings.Split(v.Buffer(), "\n")

	selected := ""

	if cy >= 0 && cy < len(lines)-1 {
		selected = lines[cy]
	}

	err = db.DeleteDocumentByID(model.State.SelectedDB, model.State.SelectedCollection, selected, model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	model.State.Documents, err = db.GetDocuments(model.State.SelectedDB, model.State.SelectedCollection, model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to retrieve collection: %v", err)
	}
	model.State.DocumentContent = ""
	updateDocumentContent(g)
	updateDocuments(g)

	return nil
}
