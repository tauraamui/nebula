package widgets

import (
	"image"
	"image/color"
	"strconv"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/tauraamui/nebula/gesturex"
	nmat "github.com/tauraamui/nebula/mat"
	"gonum.org/v1/gonum/mat"
)

type Widget interface {
	Layout(layout.Context) layout.Dimensions
}

const (
	cellWidth   unit.Dp = 80
	cellHeight          = 25
	cellPadding         = 1
)

type Matrix[T any] struct {
	Pos,
	Size f32.Point
	Color color.NRGBA
	Data  *mat.Dense
	Data2 nmat.Matrix[T]
	cellWidth,
	cellHeight,
	cellPadding int
	drag *gesturex.Drag
}

func (m *Matrix[T]) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	m.cellWidth = gtx.Dp(cellWidth)
	m.cellHeight = gtx.Dp(cellHeight)
	m.cellPadding = gtx.Dp(cellPadding)

	posX := gtx.Dp(unit.Dp(m.Pos.X))
	posY := gtx.Dp(unit.Dp(m.Pos.Y))

	cellSize := f32.Point{X: float32(gtx.Dp(cellWidth)), Y: float32(gtx.Dp(cellHeight))}

	rows, cols := m.Data.Dims()
	totalSize := f32.Point{
		X: float32(cols) * cellSize.X,
		Y: float32(rows) * cellSize.Y,
	}
	m.Size = totalSize

	background := image.Rect(posX, posY, posX+gtx.Dp(unit.Dp(m.Size.X))+(m.cellPadding), posY+gtx.Dp(unit.Dp(m.Size.Y))+m.cellPadding)
	background.Max = background.Max.Add(image.Pt(m.cellPadding*2, m.cellPadding*2))
	cl := clip.Rect{Min: background.Min, Max: background.Max}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{200, 200, 200, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl.Pop()

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			cell := image.Rect(posX+(m.cellWidth*x)+m.cellPadding, posY+(y*m.cellHeight)+m.cellPadding, posX+((m.cellWidth*x)+m.cellWidth), posY+((m.cellHeight*y)+m.cellHeight))
			cell.Min = cell.Min.Add(image.Pt(m.cellPadding, m.cellPadding))
			cell.Max = cell.Max.Add(image.Pt(m.cellPadding, m.cellPadding))
			cl := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
			paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			cl.Pop()

			cl = clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
			l := material.Label(th, unit.Sp(23), strconv.FormatFloat(m.Data.At(y, x), 'f', -1, 64))
			l.Color = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
			off := op.Offset(cell.Min.Add(image.Pt(gtx.Sp(3), 0))).Push(gtx.Ops)
			l.Layout(gtx)
			off.Pop()
			cl.Pop()
		}
	}
	return layout.Dimensions{Size: m.Size.Round()}
}

func (m *Matrix[T]) Update(gtx layout.Context) {
	if m.drag == nil {
		m.drag = &gesturex.Drag{}
	}

	pos := f32.Pt(float32(gtx.Dp(unit.Dp(m.Pos.X))), float32(gtx.Dp(unit.Dp(m.Pos.Y))))
	size := f32.Pt(m.Size.X, m.Size.Y)

	posPt := pos.Round()
	sizePt := size.Round()
	ma := image.Rect(posPt.X, posPt.Y, posPt.X+sizePt.X, posPt.Y+sizePt.Y)
	stack := clip.Rect(ma).Push(gtx.Ops)
	m.drag.Add(gtx.Ops)

	m.drag.Events(gtx.Metric, gtx.Ops, gtx.Queue, func(diff f32.Point) {
		scaledDiff := diff.Div(float32(gtx.Dp(1)))
		m.Pos = m.Pos.Sub(scaledDiff)
	})
	stack.Pop()
}
