package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Matrix struct {
	Color color.NRGBA
	Cells [][]int
}

const (
	cellWidth   unit.Dp = 130
	cellHeight  unit.Dp = 30
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
		err := loop(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	m := Matrix{
		Color: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
		Cells: [][]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
	}

	for x := 0; x < len(m.Cells); x++ {
		ext := make([]int, 8)
		m.Cells[x] = append(m.Cells[x], ext...)
	}
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
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
