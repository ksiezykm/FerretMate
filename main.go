package main

import (
	"log"

	"github.com/ksiezykm/ferretmate/pkg/config"
	"github.com/ksiezykm/ferretmate/pkg/db"
	"github.com/ksiezykm/ferretmate/pkg/ui"

	"github.com/jroimartin/gocui"
)

func main() {
	// Wczytanie zmiennych środowiskowych z pliku .env
	err := config.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Połączenie z bazą danych
	client, err := db.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer client.Disconnect(nil)

	// Inicjalizacja interfejsu użytkownika
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("Failed to create GUI: %v", err)
	}
	defer g.Close()

	g.SetManagerFunc(ui.Layout)

	// Rejestracja klawiszy
	if err := ui.RegisterKeyBindings(g, client); err != nil {
		log.Fatalf("Failed to set key bindings: %v", err)
	}

	// Start interfejsu
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("Error in main loop: %v", err)
	}
}
