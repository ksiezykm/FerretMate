package main

import (
	"log"

	"context"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/config"
	"github.com/ksiezykm/FerretMate/pkg/db"
	"github.com/ksiezykm/FerretMate/pkg/model"
	"github.com/ksiezykm/FerretMate/pkg/ui"
)

func main() {

	// Read config
	configMap, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Connect to DB and get collections
	model.State.DBclient, err = db.ConnectToDB(configMap["produkcja"])
	if err != nil {
		log.Fatalf("Failed to connect to FerretDB: %v", err)
	}
	defer model.State.DBclient.Disconnect(context.TODO())

	model.State.DBname = "testDB"

	model.State.Collections, err = db.GetCollections(model.State.DBname)
	if err != nil {
		log.Fatalf("Failed to retrieve collections: %v", err)
	}

	// Create the GUI
	g, err := gocui.NewGui(gocui.OutputNormal, true)
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
