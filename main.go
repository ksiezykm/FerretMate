// main.go
package main

import (
    "fmt"
    "log"
    "github.com/jroimartin/gocui"
)

func main() {
    g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        log.Fatal(err)
    }
    defer g.Close()

    g.SetManagerFunc(layout)

    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Fatal(err)
    }
}

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        fmt.Fprintln(v, "Welcome to FerretMate!")
    }
    return nil
}
