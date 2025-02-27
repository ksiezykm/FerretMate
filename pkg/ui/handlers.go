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
func CursorUp(g *gocui.Gui, v *gocui.View) error {
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
func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	currentView := v.Name()
	max := 0

	_, vSize := v.Size()

	switch currentView {
	case "collections":
		max = len(model.State.Collections) - 1
	case "documents":
		max = len(model.State.Documents) - 1
	case "details":
		jsonDoc, err := json.MarshalIndent(model.State.DocumentDetails, "", "  ")
		if err != nil {
			log.Println("Error converting document to JSON:", err)
			return nil
		}
		lines := strings.Split(string(jsonDoc), "\n")
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
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func selectItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	_, cy := v.Cursor()
	lines := strings.Split(v.Buffer(), "\n")

	selected := ""

	if cy >= 0 && cy < len(lines)-1 {
		selected = lines[cy]
	}
	currentView := v.Name()

	switch currentView {
	case "collections":
		model.State.SelectedCollection = selected
		model.State.Documents, err = db.GetDocuments(model.State.DBname, selected)
		if err != nil {
			log.Fatalf("Failed to retrieve collections: %v", err)
		}
		updateDocuments(g)
	case "documents":
		model.State.SelectedDocument = selected
		model.State.DocumentDetails, err = db.GetDocumentByID(model.State.DBname, model.State.SelectedCollection, selected)
		if err != nil {
			log.Fatalf("Failed to retrieve Document: %v", err)
		}
		updateDocumentDetails(g)
	}

	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentView := v.Name()
	nextView := ""

	switch currentView {
	case "collections":
		nextView = "documents"
	case "documents":
		nextView = "details"
	case "details":
		nextView = "collections"
	}

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
		nextV.SetCursor(0, 0)
	}

	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	model.State.LineToEdit = l

	if err := editView(g); err != nil {
		return err
	}

	return nil
}

// RegisterKeyBindings sets up all keybindings
func RegisterKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, CursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("collections", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("details", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, selectItem); err != nil {
		return err
	}
	if err := g.SetKeybinding("collections", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("details", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("details", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("edit", gocui.KeyArrowRight, gocui.ModNone, EditCursorRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("edit", gocui.KeyArrowLeft, gocui.ModNone, EditCursorLeft); err != nil {
		return err
	}
	return nil
}
