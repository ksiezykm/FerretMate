package ui

import "github.com/awesome-gocui/gocui"

// RegisterKeyBindings sets up all keybindings
func RegisterKeyBindings(g *gocui.Gui) error {

	RegisterKeyBindingsConnections(g)
	RegisterKeyBindingsDatabases(g)
	RegisterKeyBindingsCollections(g)
	RegisterKeyBindingsDocuments(g)
	RegisterKeyBindingsContent(g)
	RegisterKeyBindingsEdit(g)

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}
