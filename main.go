package main

import (
	"log"

	"github.com/ksiezykm/FerretMate/pkg/config"
	"github.com/ksiezykm/FerretMate/pkg/db"
	"github.com/ksiezykm/FerretMate/pkg/ui"

	"github.com/jroimartin/gocui"
)

func main() {
	// Load temp linkt o test DB from file .env
	uri, err := config.ReadDBURI(".env")
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Connection with DB
	client, err := db.ConnectToDB(uri)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer client.Disconnect(nil)

	// CUI initialization
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("Failed to create GUI: %v", err)
	}
	defer g.Close()

	g.SetManagerFunc(ui.Layout)

	// Register key bindings
	if err := ui.RegisterKeyBindings(g, client); err != nil {
		log.Fatalf("Failed to set key bindings: %v", err)
	}

	// Interface initialization
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("Error in main loop: %v", err)
	}
}
