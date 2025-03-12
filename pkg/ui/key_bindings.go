package ui

import "github.com/awesome-gocui/gocui"

// RegisterKeyBindings sets up all keybindings
func RegisterKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, CursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("databases", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("collections", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("details", gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, selectItem); err != nil {
		return err
	}
	if err := g.SetKeybinding("collections", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("databases", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("details", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("details", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("edit", gocui.KeyArrowRight, gocui.ModNone, EditCursorRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("edit", gocui.KeyArrowLeft, gocui.ModNone, EditCursorLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("edit", gocui.KeyCtrlS, gocui.ModNone, SaveChangesToEditedDocument); err != nil {
		return err
	}
	if err := g.SetKeybinding("edit", gocui.KeyEsc, gocui.ModNone, closeEditView); err != nil {
		return err
	}
	return nil
}
