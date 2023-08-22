package main

import (
	"fmt"
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

const (
	cellWidth   = 150
	cellHeight  = 50
	cellPadding = 4
)

func (m Matrix) Layout(gtx layout.Context) layout.Dimensions {
	for i := 0; i < 10; i++ {
		cell := image.Rect((cellWidth*i)+cellPadding, 0, ((cellWidth * i) + cellWidth), cellHeight)
		fmt.Printf("CELL: %+v\n", cell)
		d := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
		paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		d.Pop()
	}
	/*
		for ri, row := range m.Cells {
			y := ri
			for ci := range row {
				x := 150 * ci
				cell := image.Rect(x, y, x+150, y+50)
				fmt.Printf("CELL min: %v, max: %v\n", cell.Min, cell.Max)
				defer clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops).Pop()
			}
		}
	*/
	/*
		defer clip.Rect{Max: image.Pt(150, 50)}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: image.Pt(150, 50)}
	*/
	return layout.Dimensions{Size: image.Pt(0, 200)}
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

			paint.ColorOp{Color: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			m.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}
