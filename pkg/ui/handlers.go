// Package ui handles the user interface and keybindings
package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// CursorUp moves the cursor up in the given view
func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	if cy > 0 {
		return v.SetCursor(cx, cy-1)
	}
	return nil
}

// CursorDown moves the cursor down in the given view
func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	if cy < 1 {
		return v.SetCursor(cx, cy+1)
	}
	fmt.Fprintln(v, "name")
	return nil
}

// Quit exits the application
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// RegisterKeyBindings sets up all keybindings
func RegisterKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("collections", gocui.KeyArrowUp, gocui.ModNone, CursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("collections", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyArrowUp, gocui.ModNone, CursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}
	return nil
}
