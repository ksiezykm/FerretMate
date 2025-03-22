package ui

import "github.com/awesome-gocui/gocui"

// RegisterKeyBindings sets up all keybindings
func RegisterKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	RegisterKeyBindingsConnections(g)
	RegisterKeyBindingsDatabases(g)
	RegisterKeyBindingsCollections(g)
	RegisterKeyBindingsDocuments(g)

	if err := g.SetKeybinding("connections", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("databases", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("collections", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("details", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	// if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, selectItem); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("collections", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("databases", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("documents", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("details", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding("details", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("edit", gocui.KeyArrowRight, gocui.ModNone, editCursorRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("edit", gocui.KeyArrowLeft, gocui.ModNone, editCursorLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("edit", gocui.KeyCtrlS, gocui.ModNone, saveChangesToEditedDocument); err != nil {
		return err
	}
	if err := g.SetKeybinding("edit", gocui.KeyEsc, gocui.ModNone, closeEditView); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyCtrlN, gocui.ModNone, createNewDocument); err != nil {
		return err
	}
	if err := g.SetKeybinding("documents", gocui.KeyDelete, gocui.ModNone, deleteDocument); err != nil {
		return err
	}
	if err := g.SetKeybinding("databases", gocui.KeyCtrlN, gocui.ModNone, createNewDatabase); err != nil {
		return err
	}
	return nil
}
