package main

import (
	"log"

	"context"

	"github.com/jroimartin/gocui"
	"github.com/ksiezykm/FerretMate/pkg/db"
	"github.com/ksiezykm/FerretMate/pkg/ui"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

func main() {
	// Connect to DB and get collections
	client, err := db.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to FerretDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	model.State.Collections, err = db.GetCollections(client, "testDB")
	if err != nil {
		log.Fatalf("Failed to retrieve collections: %v", err)
	}

	// Create the GUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	g.SetManagerFunc(ui.Layout)

	// Register all keybindings
	if err := ui.RegisterKeyBindings(g); err != nil {
		log.Fatalf("Failed to register keybindings: %v", err)
	}

	// Start the main loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatal(err)
	}
}
