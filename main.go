// main.go
package main

import (
    "fmt"
    "log"

    "github.com/jroimartin/gocui"
)

func main() {
    // Initialize gocui GUI
    g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        log.Fatal(err)
    }
    defer g.Close()

    // Set the layout function
    g.SetManagerFunc(layout)

    // Set keybinding for quitting the app (Ctrl+C)
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        log.Fatal(err)
    }

    // Start the main event loop
    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Fatal(err)
    }
}

// layout defines the structure of the UI
func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()

    // Create the main view
    if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "FerretMate"
        v.Wrap = true
        fmt.Fprintln(v, "Welcome to FerretMate!")
    }
    return nil
}

// quit handles app termination
func quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}
