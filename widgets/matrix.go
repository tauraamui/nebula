package widgets

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
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
	Cells [][][]byte
	cellWidth,
	cellHeight,
	cellPadding int
	drag *gesturex.Drag
}

func (m *Matrix) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	m.cellWidth = gtx.Dp(cellWidth)
	m.cellHeight = gtx.Dp(cellHeight)
	m.cellPadding = gtx.Dp(cellPadding)

	posX := gtx.Dp(unit.Dp(m.Pos.X))
	posY := gtx.Dp(unit.Dp(m.Pos.Y))

	totalSize := image.Point{}
	totalX := 0
	totalY := 0
	cellSize := image.Point{X: gtx.Dp(cellWidth), Y: gtx.Dp(cellHeight)}
	for x, column := range m.Cells {
		totalX += 1
		for y, content := range column {
			if totalX == 1 {
				totalY += 1
			}
			cell := image.Rect(posX+(m.cellWidth*x)+m.cellPadding, posY+(y*m.cellHeight)+m.cellPadding, posX+((m.cellWidth*x)+m.cellWidth), posY+((m.cellHeight*y)+m.cellHeight))
			cell.Min = cell.Min.Add(image.Pt(m.cellPadding, m.cellPadding))
			cell.Max = cell.Max.Add(image.Pt(m.cellPadding, m.cellPadding))
			cl := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
			paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			cl.Pop()

			cl = clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
			l := material.Label(th, unit.Sp(23), string(content))
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			l.Color = maroon
			off := op.Offset(cell.Min.Add(image.Pt(gtx.Sp(3), 0))).Push(gtx.Ops)
			l.Layout(gtx)
			off.Pop()
			cl.Pop()
		}
	}
	totalSize.X = totalX * cellSize.X
	totalSize.Y = totalY * cellSize.Y
	m.Size = totalSize
	return layout.Dimensions{Size: m.Size}
}

func (m *Matrix) Update(gtx layout.Context) {
	if m.drag == nil {
		m.drag = &gesturex.Drag{}
	}

	posX := gtx.Dp(unit.Dp(m.Pos.X))
	posY := gtx.Dp(unit.Dp(m.Pos.Y))
	sizeX := m.Size.X
	sizeY := m.Size.Y

	ma := image.Rect(posX, posY, posX+sizeX, posY+sizeY)
	stack := clip.Rect(ma).Push(gtx.Ops)
	m.drag.Add(gtx.Ops)
	stack.Pop()

	m.drag.Events(unit.Metric{}, gtx.Queue, func(diff f32.Point) {
		diff = diff.Div(float32(gtx.Dp(1)))
		m.Pos = m.Pos.Sub(image.Pt(diff.Round().X, diff.Round().Y))
	})
}
