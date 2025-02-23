package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

// updateCollections
func updateCollections(g *gocui.Gui) error {
	v, err := g.View("collections")
	if err != nil {
		return err
	}
	v.Clear()
	for _, collection := range model.State.Collections {
		fmt.Fprintln(v, collection)
	}
	return nil
}

// updateDocuments
func updateDocuments(g *gocui.Gui) error {
	v, err := g.View("documents")
	if err != nil {
		return err
	}
	v.Clear()
	for _, doc := range model.State.Documents {
		fmt.Fprintln(v, doc)
	}
	return nil
}


// layout defines the UI layout with three panels
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Left panel for collections
	if v, err := g.SetView("collections", 0, 0, maxX/3, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Collections"
		v.Highlight = true
		v.SelFgColor = gocui.ColorGreen
		for _, name := range model.State.Collections {
			fmt.Fprintln(v, name)
		}
		if _, err := g.SetCurrentView("collections"); err != nil {
			return err
		}
	}

	// Middle panel for documents list
	if v, err := g.SetView("documents", maxX/3+1, 0, 2*maxX/3, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Documents"
		v.Highlight = true
		v.SelFgColor = gocui.ColorGreen
		// for _, doc := range documents {
		// 	// Serialize BSON to JSON for display
		// 	docJSON, _ := json.MarshalIndent(doc, "", "  ")
		// 	fmt.Fprintln(v, string(docJSON))
		// }
	}

	// Right panel for selected document details
	if v, err := g.SetView("details", 2*maxX/3+1, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Document Details"
		// if selectedDocument != nil {
		// 	// Serialize the selected document to JSON for display
		// 	docJSON, _ := json.MarshalIndent(selectedDocument, "", "  ")
		// 	fmt.Fprintln(v, string(docJSON))
		// } else {
		// 	fmt.Fprintln(v, "Select a document to view details.")
		// }
	}

	return nil
}
