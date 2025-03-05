package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
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

// updateDocument Details
func updateDocumentDetails(g *gocui.Gui) error {
	v, err := g.View("details")
	if err != nil {
		return err
	}
	v.Clear()

	fmt.Fprintln(v, model.State.DocumentContent)
	v.SetOrigin(0, 0)
	return nil
}

// updateMessages
func updateMessages(g *gocui.Gui) error {
	v, err := g.View("messages")
	if err != nil {
		return err
	}
	v.Clear()

	fmt.Fprintln(v, model.State.Messages)

	return nil
}

func updateEdit(g *gocui.Gui, line string) error {
	v, err := g.View("edit")
	if err != nil {
		return err
	}
	v.Clear()

	fmt.Fprintln(v, line)

	return nil
}

// layout defines the UI layout with three panels
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Left panel for database
	if v, err := g.SetView("database", 0, 0, maxX/4, (maxY-3)/2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Databases"
		v.FrameColor = gocui.ColorGreen
		v.Highlight = true
		v.Autoscroll = false
		v.Editable = false
		v.SelFgColor = gocui.ColorGreen
		for _, name := range model.State.Collections {
			fmt.Fprintln(v, name)
		}
		// if _, err := g.SetCurrentView("database"); err != nil {
		// 	return err
		// }
	}

	// Left panel for collections
	if v, err := g.SetView("collections", 0, 1+(maxY-3)/2, maxX/4, maxY-3, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Collections"
		v.FrameColor = gocui.ColorGreen
		v.Highlight = true
		v.Autoscroll = false
		v.Editable = false
		v.SelFgColor = gocui.ColorGreen
		for _, name := range model.State.Collections {
			fmt.Fprintln(v, name)
		}
		if _, err := g.SetCurrentView("collections"); err != nil {
			return err
		}
	}

	// Middle panel for documents list
	if v, err := g.SetView("documents", maxX/4+1, 0, 2*maxX/4, maxY-3, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Documents"
		v.Highlight = false
		v.SelFgColor = gocui.ColorGreen
		v.Autoscroll = false
		v.Editable = false
		v.Wrap = true
		// for _, doc := range documents {
		//      // Serialize BSON to JSON for display
		//      docJSON, _ := json.MarshalIndent(doc, "", "  ")
		//      fmt.Fprintln(v, string(docJSON))
		// }
	}

	// Right panel for selected document details
	if v, err := g.SetView("details", 2*maxX/4+1, 0, maxX-1, maxY-3, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Document Details"
		v.Highlight = false
		v.SelFgColor = gocui.ColorGreen
		v.Autoscroll = false
		v.Editable = false
		v.Wrap = true
		// if selectedDocument != nil {
		//      // Serialize the selected document to JSON for display
		//      docJSON, _ := json.MarshalIndent(selectedDocument, "", "  ")
		//      fmt.Fprintln(v, string(docJSON))
		// } else {
		//      fmt.Fprintln(v, "Select a document to view details.")
		// }
	}

	// Bottom panel for messages

	if v, err := g.SetView("messages", 0, maxY-3, maxX-1, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Messages"
		v.Autoscroll = false
		v.Editable = false
		v.Wrap = true
	}

	return nil
}
