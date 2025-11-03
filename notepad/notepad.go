package notepad

import (
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
)

// Notepad is our custom editor widget
type Notepad struct {
	Name       string
	Title      string
	Editable   bool
	Content    string
	Lines      []string
	OnEditLine func(lineNum int, oldLine string) // callback when Enter is pressed on a line
	OnBack     func()                            // callback when Esc is pressed
}

// Layout draws the notepad
func (n *Notepad) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(n.Name, maxX/2+1, 4, maxX-2, maxY-3, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = n.Title
		v.Editable = false // Make it non-editable directly
		v.Wrap = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		v.FrameColor = gocui.ColorDefault // Set initial frame color to default (inactive)
		v.Clear()
		v.Write([]byte(n.Content))
	}
	return nil
}

// Update replaces content dynamically
func (n *Notepad) Update(g *gocui.Gui, text string) error {
	n.Content = text

	// Split content into lines for line-by-line editing
	n.Lines = strings.Split(text, "\n")

	v, err := g.View(n.Name)
	if err != nil {
		return err
	}
	v.Clear()
	v.Write([]byte(n.Content))
	return nil
}

// SetActive sets the border color for active state
func (n *Notepad) SetActive(g *gocui.Gui, active bool) {
	v, err := g.View(n.Name)
	if err != nil {
		return
	}
	if active {
		v.FrameColor = gocui.ColorCyan
	} else {
		v.FrameColor = gocui.ColorDefault
	}
}

// CursorDown moves cursor down in the notepad
func (n *Notepad) CursorDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	_, h := v.Size()

	if cy+oy < len(n.Lines)-1 {
		if cy < h-1 {
			v.SetCursor(cx, cy+1)
		} else {
			v.SetOrigin(ox, oy+1)
		}
	}
	return nil
}

// CursorUp moves cursor up in the notepad
func (n *Notepad) CursorUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()

	if cy > 0 {
		v.SetCursor(cx, cy-1)
	} else if oy > 0 {
		v.SetOrigin(ox, oy-1)
	}
	return nil
}

// EditLine triggers the edit line callback
func (n *Notepad) EditLine(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, oy := v.Origin()
	lineNum := cy + oy

	if lineNum < len(n.Lines) && n.OnEditLine != nil {
		n.OnEditLine(lineNum, n.Lines[lineNum])
	}
	return nil
}

// GoBack navigates back (e.g., to document list)
func (n *Notepad) GoBack(g *gocui.Gui, v *gocui.View) error {
	if n.OnBack != nil {
		n.OnBack()
	}
	return nil
}

// BindKeys registers keybindings for notepad
func (n *Notepad) BindKeys(g *gocui.Gui) {
	if err := g.SetKeybinding(n.Name, gocui.KeyArrowDown, gocui.ModNone, n.CursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(n.Name, gocui.KeyArrowUp, gocui.ModNone, n.CursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(n.Name, gocui.KeyEnter, gocui.ModNone, n.EditLine); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(n.Name, gocui.KeyEsc, gocui.ModNone, n.GoBack); err != nil {
		log.Panicln(err)
	}
}
