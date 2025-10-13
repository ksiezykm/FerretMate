package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

// SetCurrentView switches the active view, updating highlights and messages.
func SetCurrentView(g *gocui.Gui, currentView *gocui.View, nextViewName string, helpMessage string) error {
	// Deselect the current view if it's not the global view
	if currentView != nil && currentView.Name() != "" {
		currentView.FrameColor = gocui.ColorDefault
		currentView.SelFgColor = gocui.ColorDefault
		currentView.Highlight = false
	}

	nextView, err := g.SetCurrentView(nextViewName)
	if err != nil {
		return err
	}

	nextView.Highlight = true
	nextView.FrameColor = gocui.ColorGreen
	nextView.SelFgColor = gocui.ColorGreen

	model.State.Messages = helpMessage
	updateMessages(g)

	return nil
}

// GetSelectedLine returns the currently selected line from a view.
func GetSelectedLine(v *gocui.View) (string, error) {
	if v == nil {
		return "", errors.New("view is nil")
	}
	_, cy := v.Cursor()
	lines := strings.Split(v.Buffer(), "\n")

	if cy >= 0 && cy < len(lines) {
		line := strings.TrimSpace(lines[cy])
		if line != "" {
			return line, nil
		}
	}
	return "", errors.New("no line selected or line is empty")
}

// DisplayError shows an error message in the messages view.
func DisplayError(g *gocui.Gui, err error) {
	model.State.Messages = fmt.Sprintf("Error: %v", err)
	updateMessages(g)
}