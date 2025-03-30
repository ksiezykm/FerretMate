package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/model"
)



// layout defines the UI layout with three panels
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Title panel
	if v, err := g.SetView("title", 0, 0, maxX/5, 2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = false
		v.Editable = false
		v.SelFgColor = gocui.ColorGreen

		fmt.Fprintln(v, "*****FerretMate*****")
	}

	// Left panel for connections
	if v, err := g.SetView("connections", 0, 3, maxX/5, maxY/2-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Connections (F1) "
		v.FrameColor = gocui.ColorGreen
		v.Highlight = true
		v.Autoscroll = false
		v.Editable = false
		v.SelFgColor = gocui.ColorGreen
		for k, _ := range model.State.Connections {
			fmt.Fprintln(v, k)
		}
		if _, err := g.SetCurrentView("connections"); err != nil {
			return err
		}
	}

	// Left panel for database
	if v, err := g.SetView("databases", 0, maxY/2, maxX/5, maxY-3, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Databases (F2) "
		//v.FrameColor = gocui.ColorGreen
		//v.Highlight = true
		v.Autoscroll = false
		v.Editable = false
		//v.SelFgColor = gocui.ColorGreen
		// for k, _ := range model.State.Config {
		// 	fmt.Fprintln(v, k)
		// }
	}

	// Left panel for collections
	if v, err := g.SetView("collections", maxX/5+1, 0, 2*maxX/5, maxY-3, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Collections (F3) "
		v.Highlight = false
		v.Autoscroll = false
		v.Editable = false
		v.SelFgColor = gocui.ColorGreen
		// for _, name := range model.State.Collections {
		// 	fmt.Fprintln(v, name)
		// }
		// if _, err := g.SetCurrentView("collections"); err != nil {
		// 	return err
		// }
	}

	// Middle panel for documents list
	if v, err := g.SetView("documents", 2*maxX/5+1, 0, 3*maxX/5-1, maxY-3, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Documents (F4) "
		v.Highlight = false
		v.SelFgColor = gocui.ColorGreen
		v.Autoscroll = false
		v.Editable = false
		v.Wrap = true
		// for _, doc := range documents {
		//      // Serialize BSON to JSON for display
		//      docJSON, _ := json.MarshalIndent(doc, "", "  ")
		//      fmt.Fprintln(v, string(docJSON))
		// }
	}

	// Right panel for selected document content
	if v, err := g.SetView("content", 3*maxX/5, 0, maxX-1, maxY-3, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Document Content (F5) "
		v.Highlight = false
		v.SelFgColor = gocui.ColorGreen
		v.Autoscroll = false
		v.Editable = false
		v.Wrap = true
		// if selectedDocument != nil {
		//      // Serialize the selected document to JSON for display
		//      docJSON, _ := json.MarshalIndent(selectedDocument, "", "  ")
		//      fmt.Fprintln(v, string(docJSON))
		// } else {
		//      fmt.Fprintln(v, "Select a document to view content.")
		// }
	}

	// Bottom panel for messages

	if v, err := g.SetView("messages", 0, maxY-3, maxX-1, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Messages"
		v.Autoscroll = false
		v.Editable = false
		v.Frame = false
	}

	return nil
}
