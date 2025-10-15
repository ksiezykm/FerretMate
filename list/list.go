package list

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

// List is a custom widget that displays selectable items
type List struct {
	Name     string
	Title    string
	Items    []string
	Selected int
	OnSelect func(item string) // callback when Enter is pressed
	OnBack   func()            // callback when Esc is pressed
}

// Update replaces list items and redraws the view.
func (l *List) Update(g *gocui.Gui) error {
	l.Selected = 0 // reset selection to first item (optional)

	v, err := g.View(l.Name)
	if err != nil {
		return err
	}
	v.Title = l.Title
	v.Clear()
	for _, item := range l.Items {
		v.Write([]byte(item + "\n"))
	}

	return nil
}

// Layout draws the list widget
func (l *List) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(l.Name, 1, 3, maxX/2-2, maxY-3, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = l.Title
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.FrameColor = gocui.ColorGreen // Set initial frame color to green
		v.Clear()
		for i, item := range l.Items {
			if i == l.Selected {
				// highlight selection manually
				v.SetCursor(0, i)
			}
			v.Write([]byte(item + "\n"))
		}
		if _, err := g.SetCurrentView(l.Name); err != nil {
			return err
		}
	}
	return nil
}

// SetActive sets the border color for active state
func (l *List) SetActive(g *gocui.Gui, active bool) {
	v, err := g.View(l.Name)
	if err != nil {
		return
	}
	if active {
		v.FrameColor = gocui.ColorGreen
	} else {
		v.FrameColor = gocui.ColorDefault
	}
}

// Move cursor up
func (l *List) CursorUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()

	if cy > 0 {
		l.Selected--
		v.SetCursor(cx, cy-1)
	} else if oy > 0 {
		l.Selected--
		v.SetOrigin(ox, oy-1)
	}
	return nil
}

// Move cursor down
func (l *List) CursorDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	_, h := v.Size()

	if cy+oy < len(l.Items)-1 {
		// jeśli jeszcze jest miejsce w widoku -> przesuwamy kursor
		if cy < h-1 {
			l.Selected++
			v.SetCursor(cx, cy+1)
		} else {
			// w przeciwnym wypadku przewijamy
			l.Selected++
			v.SetOrigin(ox, oy+1)
		}
	}
	return nil
}

// Select current item
func (l *List) Select(g *gocui.Gui, v *gocui.View) error {
	if l.OnSelect != nil && l.Selected < len(l.Items) {
		l.OnSelect(l.Items[l.Selected])
	}
	return nil
}

// GoBack navigates to the previous level
func (l *List) GoBack(g *gocui.Gui, v *gocui.View) error {
	if l.OnBack != nil {
		l.OnBack()
	}
	return nil
}

// BindKeys registers list-specific keybindings
func (l *List) BindKeys(g *gocui.Gui) {
	if err := g.SetKeybinding(l.Name, gocui.KeyArrowUp, gocui.ModNone, l.CursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(l.Name, gocui.KeyArrowDown, gocui.ModNone, l.CursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(l.Name, gocui.KeyEnter, gocui.ModNone, l.Select); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(l.Name, gocui.KeyEsc, gocui.ModNone, l.GoBack); err != nil {
		log.Panicln(err)
	}
}
