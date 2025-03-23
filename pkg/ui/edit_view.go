package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

var lineToEdit string
var lineToEditNumber int

func editView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("edit", maxX/2-30, maxY/2-5, maxX/2+30, maxY/2-3, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Editor"
		v.Editable = true
		v.SelBgColor = gocui.ColorGreen
		v.Editor = gocui.DefaultEditor
		updateEdit(g, insertChar(lineToEdit, '█', 0))
		v.SetCursor(0, 0)
		fmt.Fprintln(v, lineToEdit)
		if _, err := g.SetCurrentView("edit"); err != nil {
			return err
		}
		v.FrameColor = gocui.ColorGreen
	}
	if v, err := g.SetView("edit_blank_line", maxX/2-30, maxY/2-3, maxX/2+30, maxY/2-2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Frame = false
	}
	if v, err := g.SetView("edit_info", maxX/2-30, maxY/2-2, maxX/2+30, maxY/2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Keyboard shortcuts"
		fmt.Fprintln(v, "Ctrl+S: Save | Esc: Cancel")
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

func editCursorRight(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()

	var l string
	var err error

	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	l = strings.ReplaceAll(l, "█", "")
	lineToEdit = l

	if cx >= 0 {

		updateEdit(g, insertChar(lineToEdit, '█', cx+1))
		return v.SetCursor(cx+1, cy)
	}
	return nil
}
func editCursorLeft(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	var l string
	var err error

	cx, cy := v.Cursor()

	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	l = strings.ReplaceAll(l, "█", "")
	lineToEdit = l

	if cx >= 0 {
		updateEdit(g, insertChar(lineToEdit, '█', cx-1))

		if cx == 0 {
			return v.SetCursor(0, cy)
		} else {
			return v.SetCursor(cx-1, cy)
		}
	}
	return nil
}

func closeEditView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("edit"); err != nil {
		return err
	}
	if err := g.DeleteView("edit_blank_line"); err != nil {
		return err
	}
	if err := g.DeleteView("edit_info"); err != nil {
		return err
	}
	if err := g.DeleteView("messages"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("content"); err != nil {
		return err
	}
	return nil

}
