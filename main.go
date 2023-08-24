package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/gesture"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

const (
	cellWidth   unit.Dp = 130
	cellHeight  unit.Dp = 30
	cellPadding unit.Dp = 2
)

type Matrix struct {
	Pos,
	Size image.Point
	Color color.NRGBA
	Cells [][]int
}

func (m *Matrix) Layout(gtx layout.Context) layout.Dimensions {
	cellWidthPx := gtx.Dp(cellWidth)
	cellHeightPx := gtx.Dp(cellHeight)
	cellPaddingPx := gtx.Dp(cellPadding)
	totalSize := image.Point{}
	for x, row := range m.Cells {
		for y := range row {
			cell := image.Rect(m.Pos.X+(cellWidthPx*x)+cellPaddingPx, m.Pos.Y+(y*cellHeightPx)+cellPaddingPx, m.Pos.X+((cellWidthPx*x)+cellWidthPx), m.Pos.Y+((cellHeightPx*y)+cellHeightPx))
			cell.Min = cell.Min.Add(image.Pt(cellPaddingPx, cellPaddingPx))
			cell.Max = cell.Max.Add(image.Pt(cellPaddingPx, cellPaddingPx))
			rect := clip.Rect{Min: cell.Min, Max: cell.Max}
			cl := rect.Push(gtx.Ops)
			totalSize.X += rect.Max.X - rect.Min.X
			totalSize.Y += rect.Max.Y - rect.Min.Y
			paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			cl.Pop()
		}
	}
	m.Size = totalSize
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
	m := &Matrix{
		Pos:   image.Pt(20, 20),
		Color: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
		Cells: [][]int{{0}},
	}

	/*
		for x := 0; x < len(m.Cells); x++ {
			ext := make([]int, 3)
			m.Cells[x] = append(m.Cells[x], ext...)
		}
	*/
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	var drag gesture.Drag
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ops.Reset()
			gtx := layout.NewContext(&ops, e)

			paint.ColorOp{Color: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			m.Layout(gtx)

			ma := image.Rect(m.Pos.X, m.Pos.Y, m.Pos.X+m.Size.X, m.Pos.Y+m.Size.Y)
			stack := clip.Rect(ma).Push(gtx.Ops)
			drag.Add(gtx.Ops)
			stack.Pop()

			de := drag.Events(unit.Metric{}, gtx.Queue, gesture.Both)
			for i := 0; i < len(de); i++ {
				fmt.Printf("DE: %+v\n", de)
			}

			e.Frame(gtx.Ops)
		}
	}
}
