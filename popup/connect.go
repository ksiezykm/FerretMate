package popup

import (
	"fmt"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/ksiezykm/FerretMate/db"
	"github.com/ksiezykm/FerretMate/model"
)

type ConnectPopup struct {
	g          *gocui.Gui
	cancel     chan struct{}
	cancelOnce sync.Once
}

func ShowConnect(g *gocui.Gui, conn model.Connection, onSuccess func() error) {
	cp := &ConnectPopup{
		g:      g,
		cancel: make(chan struct{}),
	}

	maxX, maxY := g.Size()
	width := 40
	height := 7
	x0 := (maxX - width) / 2
	y0 := (maxY - height) / 2
	x1 := x0 + width
	y1 := y0 + height

	g.Update(func(g *gocui.Gui) error {
		v, err := g.SetView("connect_popup", x0, y0, x1, y1, 0)
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Connecting "
		v.Clear()
		v.Write([]byte("\n  Connecting...\n\n  Press ESC to cancel"))
		g.SetCurrentView("connect_popup")
		return nil
	})

	closePopup := func() {
		g.Update(func(g *gocui.Gui) error {
			g.DeleteView("connect_popup")
			g.DeleteKeybindings("connect_popup")
			g.SetCurrentView("listView")
			return nil
		})
	}

	handleEsc := func(g *gocui.Gui, v *gocui.View) error {
		cp.cancelOnce.Do(func() {
			close(cp.cancel)
		})
		closePopup()
		return nil
	}

	g.Update(func(g *gocui.Gui) error {
		g.SetKeybinding("connect_popup", gocui.KeyEsc, gocui.ModNone, handleEsc)
		return nil
	})

	go func() {
		i := 0
		for {
			select {
			case <-cp.cancel:
				return
			default:
			}

			g.Update(func(g *gocui.Gui) error {
				v, _ := g.View("connect_popup")
				if v != nil {
					v.Clear()
					v.Write([]byte(fmt.Sprintf("\n  Connecting... %ds\n\n  Press ESC to cancel", i)))
				}
				return nil
			})

			i++
			time.Sleep(time.Second)
		}
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)

		select {
		case <-cp.cancel:
			return
		default:
		}

		connDone := make(chan error, 1)
		go func() {
			connDone <- db.Connect(conn)
		}()

		var err error
		select {
		case err = <-connDone:
		case <-time.After(10 * time.Second):
			err = fmt.Errorf("connection timeout")
		case <-cp.cancel:
			return
		}

		cp.cancelOnce.Do(func() {
			close(cp.cancel)
		})

		g.Update(func(g *gocui.Gui) error {
			v, _ := g.View("connect_popup")
			if v == nil {
				return nil
			}

			v.Clear()
			if err != nil {
				v.Title = " Error "
				v.Write([]byte("\n  Connection failed\n\n  Press ESC to close"))
			} else {
				v.Title = " Success "
				v.Write([]byte("\n  Connected!\n\n  Press ESC to close"))
				if onSuccess != nil {
					onSuccess()
				}
			}

			handleClose := func(g *gocui.Gui, v *gocui.View) error {
				closePopup()
				return nil
			}

			g.SetKeybinding("connect_popup", gocui.KeyEsc, gocui.ModNone, handleClose)

			if err == nil {
				time.AfterFunc(500*time.Millisecond, func() {
					handleClose(g, v)
				})
			}

			return nil
		})
	}()
}
