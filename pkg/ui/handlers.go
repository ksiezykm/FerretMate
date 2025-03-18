// Package ui handles the user interface and keybindings
package ui

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/db"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

// CursorUp moves the cursor up in the given view
func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	_, vOrigin := v.Origin()
	cx, cy := v.Cursor()
	if cy >= 0 {
		if cy == 0 {
			v.SetOrigin(0, vOrigin-1)
			return v.SetCursor(cx, 0)
		}

		return v.SetCursor(cx, cy-1)
	}
	return nil
}

// CursorDown moves the cursor down in the given view
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	currentView := v.Name()
	max := 0

	_, vSize := v.Size()

	switch currentView {
	case "databases":
		max = len(model.State.Config) - 1
	case "collections":
		max = len(model.State.Collections) - 1
	case "documents":
		max = len(model.State.Documents) - 1
	case "details":
		lines := strings.Split(model.State.DocumentContent, "\n")
		max = len(lines) - 1
	}

	cx, cy := v.Cursor()
	_, vOrigin := v.Origin()

	if cy < max {
		if cy >= vSize-1 {
			v.SetOrigin(0, vOrigin+1)
			return v.SetCursor(cx, vSize-1)
		}
		return v.SetCursor(cx, cy+1)
	}

	g.Update(func(g *gocui.Gui) error { return nil })
	return nil
}

// Quit exits the application
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func selectItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	_, cy := v.Cursor()
	lines := strings.Split(v.Buffer(), "\n")

	// model.State.Messages = fmt.Sprint(cy)
	// updateMessages(g)

	selected := ""

	if cy >= 0 && cy < len(lines)-1 {
		selected = lines[cy]
	}
	currentView := v.Name()

	switch currentView {
	case "databases":
		model.State.Collections = nil

		model.State.SelectedDB = model.State.Config[selected].Database
		model.State.DBclient, err = db.Connect(model.State.Config[selected])
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
		updateDocumentDetails(g)
		model.State.Documents = nil
		updateDocuments(g)
	case "collections":
		model.State.SelectedCollection = selected
		model.State.Documents, err = db.GetDocuments(model.State.SelectedDB, model.State.SelectedCollection, model.State.DBclient)
		if err != nil {
			log.Fatalf("Failed to retrieve collection: %v", err)
		}
		model.State.DocumentContent = ""
		updateDocumentDetails(g)
		updateDocuments(g)
	case "documents":
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
		updateDocumentDetails(g)
	}

	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentView := v.Name()
	nextView := ""

	switch currentView {
	case "databases":
		nextView = "collections"
	case "collections":
		nextView = "documents"
	case "documents":
		nextView = "details"
	case "details":
		nextView = "databases"
	}

	v.FrameColor = gocui.ColorDefault

	if _, err := g.SetCurrentView(nextView); err != nil {
		return err
	}

	if nextView != "" {
		nextV, err := g.View(nextView)
		if err != nil {
			return err
		}

		v.Highlight = false

		if _, err := g.SetCurrentView(nextView); err != nil {
			return err
		}
		nextV.Highlight = true
		nextV.FrameColor = gocui.ColorGreen
		nextV.SetCursor(0, 0)
	}

	documentsInfo := "Tab: next view | Enter: view | Delete: delete | Ctrl+n: new"
	detailsInfo := "Tab: next view | Enter: edit line"

	switch nextView {

	case "documents":
		model.State.Messages = documentsInfo
	case "details":
		model.State.Messages = detailsInfo
	default:
		model.State.Messages = ""
	}
	updateMessages(g)
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	if !strings.Contains(l, "_id") {

		lineToEditNumber = cy
		lineToEdit = l

		if err := editView(g); err != nil {
			return err
		}
	}
	return nil
}

func saveChangesToEditedDocument(g *gocui.Gui, v *gocui.View) error {
	var line string
	var err error

	_, cy := v.Cursor()
	if line, err = v.Line(cy); err != nil {
		line = ""
	}

	line = strings.ReplaceAll(line, "█", "")
	lines := strings.Split(model.State.DocumentContent, "\n")

	lines[lineToEditNumber] = line

	model.State.DocumentContent = strings.Join(lines, "\n")

	updateDocumentDetails(g)

	db.UpdateDocumentByID(model.State.SelectedDB, model.State.SelectedCollection, model.State.SelectedDocument, model.State.DocumentContent, model.State.DBclient)

	return nil
}

func createNewDocument(g *gocui.Gui, v *gocui.View) error {
	var err error
	// _, cy := v.Cursor()
	// lines := strings.Split(v.Buffer(), "\n")

	// selected := ""

	// if cy >= 0 && cy < len(lines)-1 {
	// 	selected = lines[cy]
	// }

	// model.State.SelectedCollection = selected
	err = db.CreateDocument(model.State.SelectedDB, model.State.SelectedCollection, model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	// model.State.Messages = fmt.Sprint(insertedId)
	// updateMessages(g)

	model.State.Documents, err = db.GetDocuments(model.State.SelectedDB, model.State.SelectedCollection, model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to retrieve collection: %v", err)
	}
	model.State.DocumentContent = ""
	updateDocumentDetails(g)
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
	updateDocumentDetails(g)
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
