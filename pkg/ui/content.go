package ui

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

func RegisterKeyBindingsContent(g *gocui.Gui) error {
	// if err := g.SetKeybinding("content", gocui.KeyEnter, gocui.ModNone, selectContent); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModNone, setCurrentViewContent); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyF5, gocui.ModNone, setCurrentViewContent); err != nil {
		return err
	}
	if err := g.SetKeybinding("content", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	return nil
}

func setCurrentViewContent(g *gocui.Gui, v *gocui.View) error {

	v.FrameColor = gocui.ColorDefault
	v.SelFgColor = gocui.ColorDefault
	var nextView *gocui.View
	var err error
	if nextView, err = g.SetCurrentView("content"); err != nil {
		return err
	}
	nextView.Highlight = true
	nextView.FrameColor = gocui.ColorGreen
	nextView.SelFgColor = gocui.ColorGreen
	// nextView.SetOrigin(0, 0)
	// nextView.SetCursor(0, 0)

	model.State.Messages = "Enter: edit"
	updateMessages(g)

	return nil
}

// updateDocument Content
func updateDocumentContent(g *gocui.Gui) error {
	v, err := g.View("content")
	if err != nil {
		return err
	}
	v.Clear()

	fmt.Fprintln(v, model.State.DocumentContent)
	v.SetOrigin(0, 0)
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	if l == "" {
		return nil
	}

	if !strings.Contains(l, "_id") {

		lineToEditNumber = cy
		lineToEdit = l

		if err := editView(g); err != nil {
			return err
		}
	}
	return nil
}

// func selectContent(g *gocui.Gui, v *gocui.View) error {
// 	var l string
// 	var err error

// 	_, cy := v.Cursor()
// 	if l, err = v.Line(cy); err != nil {
// 		l = ""
// 	}

// 	if !strings.Contains(l, "_id") {

// 		lineToEditNumber = cy
// 		lineToEdit = l

// 		if err := editView(g); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
