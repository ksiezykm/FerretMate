package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/db"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

func RegisterKeyBindingsCollections(g *gocui.Gui) error {
	if err := g.SetKeybinding("collections", gocui.KeyEnter, gocui.ModNone, selectCollection); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlI, gocui.ModNone, setCurrentViewCollections); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyF3, gocui.ModNone, setCurrentViewCollections); err != nil {
		return err
	}
	return nil
}

func setCurrentViewCollections(g *gocui.Gui, v *gocui.View) error {

	v.FrameColor = gocui.ColorDefault
	v.SelFgColor = gocui.ColorDefault
	var nextView *gocui.View
	var err error
	if nextView, err = g.SetCurrentView("collections"); err != nil {
		return err
	}
	nextView.Highlight = true
	nextView.FrameColor = gocui.ColorGreen
	nextView.SelFgColor = gocui.ColorGreen
	// nextView.SetOrigin(0, 0)
	// nextView.SetCursor(0, 0)

	model.State.Messages = "Enter: view | Delete: delete | Ctrl+n: new"
	updateMessages(g)

	return nil
}

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

func selectCollection(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	lines := strings.Split(v.Buffer(), "\n")

	selected := ""

	if cy >= 0 && cy < len(lines)-1 {
		selected = lines[cy]
	}

	if selected == "" {
		return nil
	}

	var err error

	model.State.SelectedCollection = selected
	model.State.Documents, err = db.GetDocuments(model.State.SelectedDB, model.State.SelectedCollection, model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to retrieve collection: %v", err)
	}
	model.State.DocumentContent = ""
	updateDocumentContent(g)
	updateDocuments(g)

	return nil
}
