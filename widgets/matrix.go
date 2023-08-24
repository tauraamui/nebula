package widgets

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/tauraamui/nebula/gesturex"
)

type Widget interface {
	Layout(layout.Context) layout.Dimensions
}

const (
	cellWidth   = 130
	cellHeight  = 30
	cellPadding = 2
)

type Matrix struct {
	Pos,
	Size image.Point
	Color color.NRGBA
	Cells [][]int
	drag  *gesturex.Drag
}

func (m *Matrix) Layout(gtx layout.Context) layout.Dimensions {
	if m.drag == nil {
		m.drag = &gesturex.Drag{}
	}

	totalSize := image.Point{}
	totalX := 0
	totalY := 0
	cellSize := image.Point{X: cellWidth, Y: cellHeight}
	for x, column := range m.Cells {
		totalX += 1
		for y := range column {
			if totalX == 1 {
				totalY += 1
			}
			cell := image.Rect(m.Pos.X+(cellWidth*x)+cellPadding, m.Pos.Y+(y*cellHeight)+cellPadding, m.Pos.X+((cellWidth*x)+cellWidth), m.Pos.Y+((cellHeight*y)+cellHeight))
			cell.Min = cell.Min.Add(image.Pt(cellPadding, cellPadding))
			cell.Max = cell.Max.Add(image.Pt(cellPadding, cellPadding))
			cl := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
			paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			cl.Pop()
		}
	}
	totalSize.X = totalX * cellSize.X
	totalSize.Y = totalY * cellSize.Y
	m.Size = totalSize
	return layout.Dimensions{Size: m.Size}
}
