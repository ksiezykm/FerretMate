package ui

import "github.com/awesome-gocui/gocui"

func RegisterKeyBindingsDatabases(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlU, gocui.ModNone, setCurrentViewDatabases); err != nil {
		return err
	}
	return nil
}

func setCurrentViewDatabases(g *gocui.Gui, v *gocui.View) error {

	v.FrameColor = gocui.ColorDefault
	v.SelFgColor = gocui.ColorDefault
	var nextView *gocui.View
	var err error
	if nextView, err = g.SetCurrentView("databases"); err != nil {
		return err
	}
	nextView.Highlight = true
	nextView.FrameColor = gocui.ColorGreen
	nextView.SelFgColor = gocui.ColorGreen
	// nextView.SetOrigin(0, 0)
	// nextView.SetCursor(0, 0)
	return nil
}