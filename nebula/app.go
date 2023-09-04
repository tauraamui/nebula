package nebula

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"github.com/tauraamui/nebula/widgets"
)

type App struct {
	w *app.Window
	c *widgets.Canvas
}

func New() App {
	return App{
		w: app.NewWindow(
			app.Title("github.com/tauraamui/nebula"),
		),
		c: widgets.NewCanvas(),
	}
}

func (a *App) Run() error {
	var ops op.Ops
	var err error
updateProc:
	for {
		e := <-a.w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			err = e.Err
			break updateProc
		case system.FrameEvent:
			ops.Reset()
			a.c.Update(&ops, e)
			e.Frame(&ops)
		}
	}
	return err
}
