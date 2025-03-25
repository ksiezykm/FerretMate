// Package ui handles the user interface and keybindings
package ui

import (
	"strings"

	"github.com/awesome-gocui/gocui"
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
	case "connections":
		max = len(model.State.Connections) - 1
	case "databases":
		max = len(model.State.Connections) - 1
	case "collections":
		max = len(model.State.Collections) - 1
	case "documents":
		max = len(model.State.Documents) - 1
	case "content":
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
