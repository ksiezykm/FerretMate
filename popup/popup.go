package popup

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

// Popup is a modal dialog for editing text
type Popup struct {
	Name     string
	Title    string
	Content  string
	OnSave   func(newContent string) // callback when content is saved
	OnCancel func()                  // callback when cancelled
}

// Show displays the popup
func (p *Popup) Show(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	width := maxX * 2 / 3
	height := maxY / 3
	x0 := (maxX - width) / 2
	y0 := (maxY - height) / 2
	x1 := x0 + width
	y1 := y0 + height

	if v, err := g.SetView(p.Name, x0, y0, x1, y1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = p.Title
		v.Editable = true
		v.Wrap = true
		v.Clear()
		v.Write([]byte(p.Content))

		// Set cursor at the end of content
		lines := len(v.BufferLines())
		if lines > 0 {
			lastLine := v.BufferLines()[lines-1]
			v.SetCursor(len(lastLine), lines-1)
		}
	}

	// Enable cursor for editing
	g.Cursor = true

	// Set focus to popup
	if _, err := g.SetCurrentView(p.Name); err != nil {
		return err
	}

	return nil
}

// Hide removes the popup
func (p *Popup) Hide(g *gocui.Gui) error {
	// Disable cursor when closing popup
	g.Cursor = false

	if err := g.DeleteView(p.Name); err != nil {
		return err
	}
	// Unbind keys
	g.DeleteKeybindings(p.Name)
	return nil
}

// Save saves the edited content
func (p *Popup) Save(g *gocui.Gui, v *gocui.View) error {
	if p.OnSave != nil {
		// Get the edited content from the view
		content := v.Buffer()
		// Remove trailing newline that gocui adds
		if len(content) > 0 && content[len(content)-1] == '\n' {
			content = content[:len(content)-1]
		}
		p.OnSave(content)
	}
	return p.Hide(g)
}

// Cancel closes the popup without saving
func (p *Popup) Cancel(g *gocui.Gui, v *gocui.View) error {
	if p.OnCancel != nil {
		p.OnCancel()
	}
	return p.Hide(g)
}

// BindKeys registers keybindings for the popup
func (p *Popup) BindKeys(g *gocui.Gui) {
	// Ctrl+S to save
	if err := g.SetKeybinding(p.Name, gocui.KeyCtrlS, gocui.ModNone, p.Save); err != nil {
		log.Panicln(err)
	}
	// ESC to cancel
	if err := g.SetKeybinding(p.Name, gocui.KeyEsc, gocui.ModNone, p.Cancel); err != nil {
		log.Panicln(err)
	}
}
