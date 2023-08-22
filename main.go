package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Matrix struct {
	Color color.NRGBA
	Cells [][]int
}

const (
	cellWidth   unit.Dp = 150
	cellHeight  unit.Dp = 50
	cellPadding unit.Dp = 2
)

var zoomLevel unit.Dp = 100

func (m Matrix) Layout(gtx layout.Context) layout.Dimensions {
	cellWidthPx := gtx.Dp(cellWidth)
	cellHeightPx := gtx.Dp(cellHeight)
	cellPaddingPx := gtx.Dp(cellPadding)
	totalSize := image.Point{}
	for x, row := range m.Cells {
		for y := range row {
			cell := image.Rect((gtx.Dp(cellWidth)*x)+cellPaddingPx, (y*cellHeightPx)+cellPaddingPx, ((cellWidthPx * x) + cellWidthPx), ((cellHeightPx * y) + cellHeightPx))
			cell.Min = cell.Min.Add(image.Pt(cellPaddingPx, cellPaddingPx))
			cell.Max = cell.Max.Add(image.Pt(cellPaddingPx, cellPaddingPx))
			cl := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
			totalSize.Add(cell.Bounds().Size())
			paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			cl.Pop()
		}
	}
	return layout.Dimensions{Size: totalSize}
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

	for x := 0; x < len(m.Cells); x++ {
		ext := make([]int, 8)
		m.Cells[x] = append(m.Cells[x], ext...)
	}

	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			po := pointer.InputOp{Tag: 0, Types: pointer.Scroll, ScrollBounds: image.Rect(-3, -3, 3, 3)}
			po.Add(gtx.Ops)

			for _, event := range gtx.Queue.Events(0) {
				if pe, ok := event.(pointer.Event); ok {
					if pe.Type == pointer.Scroll {
						zoomLevel += unit.Dp(pe.Scroll.Y)
						if zoomLevel < 5 {
							zoomLevel = 5
						} else if zoomLevel > 300 {
							zoomLevel = 300
						}
					}
				}
			}

			zoomLevelPx := gtx.Dp(zoomLevel)
			scale := op.Affine(f32.Affine2D{}.Scale(f32.Point{}, f32.Point{X: float32(zoomLevelPx) / 110, Y: float32(zoomLevelPx) / 110})).Push(gtx.Ops)

			paint.ColorOp{Color: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			m.Layout(gtx)

			e.Frame(gtx.Ops)

			scale.Pop()
		}
	}
}
