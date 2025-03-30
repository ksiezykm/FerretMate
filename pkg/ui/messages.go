package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/pkg/model"
)

// updateMessages
func updateMessages(g *gocui.Gui) error {
	v, err := g.View("messages")
	if err != nil {
		return err
	}
	v.Clear()

	fmt.Fprintln(v, model.State.Messages)

	return nil
}
