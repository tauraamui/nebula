package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Matrix struct {
	Color color.NRGBA
	Cells [][]int
}

func (m Matrix) Layout(gtx layout.Context) layout.Dimensions {
	defer clip.Rect{Max: image.Pt(150, 50)}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: image.Pt(150, 50)}
}

func main() {
	go func() {
		w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	m := Matrix{
		Color: color.NRGBA{R: 127, G: 0, B: 0, A: 255},
		Cells: [][]int{{0}, {0}, {0}, {0}},
	}
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			m.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}
