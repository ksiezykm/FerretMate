package ui

import (
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/db"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

var conView string = "connections"

func RegisterKeyBindingsConnections(g *gocui.Gui) error {
	if err := g.SetKeybinding(conView, gocui.KeyEnter, gocui.ModNone, selectConnection); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlY, gocui.ModNone, activateConnections); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlU, gocui.ModNone, activateDatabases); err != nil {
		return err
	}
	return nil
}
func activateConnections(g *gocui.Gui, v *gocui.View) error {

	v.FrameColor = gocui.ColorDefault
	var nextView *gocui.View
	var err error
	if nextView, err = g.SetCurrentView("connections"); err != nil {
		return err
	}
	nextView.Highlight = true
	nextView.FrameColor = gocui.ColorGreen
	// nextView.SetCursor(0, 0)
	// // nextView.SetOrigin(0, 0)
	return nil

}
func activateDatabases(g *gocui.Gui, v *gocui.View) error {

	v.FrameColor = gocui.ColorDefault
	v.SelFgColor = gocui.ColorDefault
	var nextView *gocui.View
	var err error
	if nextView, err = g.SetCurrentView("databases"); err != nil {
		return err
	}
	nextView.Highlight = true
	nextView.FrameColor = gocui.ColorGreen
	nextView.SelFgColor = gocui.ColorGreen
	// nextView.SetOrigin(0, 0)
	// nextView.SetCursor(0, 0)
	return nil
}

func selectConnection(g *gocui.Gui, v *gocui.View) error {
	var err error
	_, cy := v.Cursor()
	lines := strings.Split(v.Buffer(), "\n")

	selected := ""

	if cy >= 0 && cy < len(lines)-1 {
		selected = lines[cy]
	}

	model.State.DBclient, err = db.Connect(model.State.Connections[selected])
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	model.State.DBnames, err = db.GetDBs(model.State.DBclient)
	if err != nil {
		log.Fatalf("Failed to get DBs: %v", err)
	}

	updateDatabases(g)

	return nil
}
