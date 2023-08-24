package widgets

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/tauraamui/nebula/gesturex"
)

type Widget interface {
	Layout(layout.Context) layout.Dimensions
}

const (
	cellWidth   unit.Dp = 80
	cellHeight          = 25
	cellPadding         = 1
)

type Matrix struct {
	Pos,
	Size image.Point
	Color color.NRGBA
	Cells [][]int
	cellWidth,
	cellHeight,
	cellPadding int
}

func (m *Matrix) Layout(gtx layout.Context) layout.Dimensions {
	dimensions := render(gtx, m)

	drag := gesturex.Drag{}
	ma := image.Rect(m.Pos.X, m.Pos.Y, m.Pos.X+m.Size.X, m.Pos.Y+m.Size.Y)
	stack := clip.Rect(ma).Push(gtx.Ops)
	drag.Add(gtx.Ops)
	stack.Pop()

	drag.Events(unit.Metric{PxPerDp: 1, PxPerSp: 1}, gtx.Queue, func(diff f32.Point) {
		m.Pos = m.Pos.Sub(image.Pt(diff.Round().X, diff.Round().Y))
	})

	return dimensions
}

func render(gtx layout.Context, m *Matrix) layout.Dimensions {
	m.cellWidth = gtx.Dp(cellWidth)
	m.cellHeight = gtx.Dp(cellHeight)
	m.cellPadding = gtx.Dp(cellPadding)

	totalSize := image.Point{}
	totalX := 0
	totalY := 0
	cellSize := image.Point{X: gtx.Dp(cellWidth), Y: gtx.Dp(cellHeight)}
	for x, column := range m.Cells {
		totalX += 1
		for y := range column {
			if totalX == 1 {
				totalY += 1
			}
			cell := image.Rect(m.Pos.X+(m.cellWidth*x)+m.cellPadding, m.Pos.Y+(y*m.cellHeight)+m.cellPadding, m.Pos.X+((m.cellWidth*x)+m.cellWidth), m.Pos.Y+((m.cellHeight*y)+m.cellHeight))
			cell.Min = cell.Min.Add(image.Pt(m.cellPadding, m.cellPadding))
			cell.Max = cell.Max.Add(image.Pt(m.cellPadding, m.cellPadding))
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
