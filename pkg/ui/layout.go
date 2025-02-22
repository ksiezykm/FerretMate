package ui

import (
	"github.com/jroimartin/gocui"
)

// Layout defines the structure of the UI windows.
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Kolekcje (lewe okno)
	if v, err := g.SetView("collections", 0, 0, maxX/3-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Collections"
		v.Wrap = false
	}

	// Dokumenty (środkowe okno)
	if v, err := g.SetView("documents", maxX/3, 0, 2*maxX/3-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Documents"
		v.Wrap = false
	}

	// Szczegóły dokumentu (prawe okno)
	if v, err := g.SetView("details", 2*maxX/3, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Details"
		v.Wrap = true
	}

	return nil
}
