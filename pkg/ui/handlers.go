package ui

import (
	"encoding/json"
	"fmt"

	"github.com/ksiezykm/github.com/yourusername/ferretmate/pkg/db"

	"github.com/jroimartin/gocui"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterKeyBindings defines key bindings for navigation.
func RegisterKeyBindings(g *gocui.Gui, client *mongo.Client) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	// Enter to select collection
	if err := g.SetKeybinding("collections", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return showDocuments(g, v, client)
	}); err != nil {
		return err
	}

	// Enter to select document
	if err := g.SetKeybinding("documents", gocui.KeyEnter, gocui.ModNone, showDocumentDetails); err != nil {
		return err
	}

	return nil
}

// showDocuments retrieves documents from the selected collection.
func showDocuments(g *gocui.Gui, v *gocui.View, client *mongo.Client) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}

	collectionName := line
	docs, err := db.GetDocuments(client, "testDB", collectionName)
	if err != nil {
		return err
	}

	// Update document view
	docView, _ := g.View("documents")
	docView.Clear()
	for _, doc := range docs {
		docJSON, _ := json.MarshalIndent(doc, "", "  ")
		fmt.Fprintln(docView, string(docJSON))
	}

	return nil
}

// showDocumentDetails shows the selected document's details.
func showDocumentDetails(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}

	detailsView, _ := g.View("details")
	detailsView.Clear()
	fmt.Fprintln(detailsView, line)

	return nil
}

// quit exits the application.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
