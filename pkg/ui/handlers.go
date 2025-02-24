// Package ui handles the user interface and keybindings
package ui

import (
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/ksiezykm/FerretMate/pkg/db"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

// CursorUp moves the cursor up in the given view
func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	if cy > 0 {
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

	switch currentView {
	case "collections":
		max = len(model.State.Collections) - 1
	case "documents":
		max = len(model.State.Documents) - 1
	}

	cx, cy := v.Cursor()
	if cy < max {
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
		model.State.Documents, err = db.GetDocuments("testDB", selected)
		if err != nil {
			log.Fatalf("Failed to retrieve collections: %v", err)
		}
		updateDocuments(g)
	case "documents":
		model.State.DocumentDetails, err = db.GetDocumentByID("testDB", "testCollection", selected)
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
	return nil
}
