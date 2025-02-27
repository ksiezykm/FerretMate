package ui

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

func editView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("edit", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Edit"
		v.Editable = true
		v.Highlight = true
		v.Wrap = true
		v.SelBgColor = gocui.ColorGreen
		fmt.Fprintln(v, model.State.LineToEdit)
		if _, err := g.SetCurrentView("edit"); err != nil {
			return err
		}
	}
	return nil
}

func insertChar(s string, char rune, index int) string {
	if index < 0 || index > len(s) {
		return s // Index out of range
	}

	runes := []rune(s)
	runes = append(runes[:index], append([]rune{char}, runes[index:]...)...)
	return string(runes)
}

func EditCursorRight(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	if cx >= 0 {

		updateEdit(g, insertChar(model.State.LineToEdit, '█', cx+1))

		model.State.Messages = fmt.Sprintf("%d", cx)
		updateMessages(g)

		return v.SetCursor(cx+1, cy)
	}
	return nil
}
func EditCursorLeft(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	if cx >= 0 {

		updateEdit(g, insertChar(model.State.LineToEdit, '█', cx))

		model.State.Messages = fmt.Sprintf("%d", cx)
		updateMessages(g)

		return v.SetCursor(cx-1, cy)
	}
	return nil
}
