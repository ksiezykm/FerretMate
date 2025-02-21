package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// readDBURI reads the MongoDB URI from the .env file
func readDBURI(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", fmt.Errorf("empty or invalid .env file")
}

// connectToDB connects to FerretDB using the URI from .env file
func connectToDB() (*mongo.Client, error) {
	uri, err := readDBURI(".env")
	if err != nil {
		return nil, err
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	// Connect to the database
	client, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to FerretDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	log.Println("Connected to FerretDB!")

	// Initialize gocui
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	// Set layout and keybindings
	g.SetManagerFunc(layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Fatal(err)
	}

	// Start the main loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatal(err)
	}
}

// layout defines the UI layout
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Create the main view
	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "FerretMate"
		v.Wrap = true
		fmt.Fprintln(v, "Welcome to FerretMate!\nConnected to FerretDB successfully!")
	}

	return nil
}

// quit exits the application
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
