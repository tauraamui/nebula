package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
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

const (
	cellWidth   = 150
	cellHeight  = 50
	cellPadding = 4
	zoomLevel   = 0.6
)

func (m Matrix) Layout(gtx layout.Context) layout.Dimensions {
	for x, row := range m.Cells {
		for y := range row {
			cell := image.Rect((cellWidth*x)+cellPadding, (y*cellHeight)+cellPadding, ((cellWidth * x) + cellWidth), ((cellHeight * y) + cellHeight))
			cl := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
			paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			cl.Pop()
		}
	}
	return layout.Dimensions{Size: image.Pt(100, 200)}
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
		Color: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
		Cells: [][]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
	}
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			scale := op.Affine(f32.Affine2D{}.Scale(f32.Point{}, f32.Point{X: zoomLevel, Y: zoomLevel})).Push(gtx.Ops)

			paint.ColorOp{Color: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			m.Layout(gtx)

			e.Frame(gtx.Ops)

			scale.Pop()
		}
	}
}
