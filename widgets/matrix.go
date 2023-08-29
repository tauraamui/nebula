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

func (m *Matrix[T]) Layout(gtx layout.Context, th *material.Theme, debug bool) layout.Dimensions {
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

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			renderCell(gtx, strconv.FormatFloat(m.Data.At(y, x), 'f', -1, 64), x, y, posX, posY, m.cellWidth, m.cellHeight, m.Color, th)
		}
	}
	return layout.Dimensions{Size: m.Size.Round()}
}

func renderCell(gtx layout.Context, content string, x, y int, posx, posy, cellwidth, cellheight int, bgcolor color.NRGBA, th *material.Theme) {
	// render background of cell
	cell := image.Rect(posx+(cellwidth*x), posy+(y*cellheight), posx+((cellwidth*x)+cellwidth), posy+((cellheight*y)+cellheight))
	cl1 := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
	paint.ColorOp{Color: bgcolor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl1.Pop()

	// render cell content as text label
	cl2 := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
	l := material.Label(th, unit.Sp(23), content)
	l.Color = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
	off := op.Offset(cell.Min.Add(image.Pt(gtx.Sp(3), 0))).Push(gtx.Ops)
	l.Layout(gtx)
	off.Pop()
	cl2.Pop()

	cl3 := clip.Stroke{Path: clip.RRect{Rect: cell}.Path(gtx.Ops), Width: 1}.Op().Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{R: 12, G: 12, B: 12, A: 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl3.Pop()

}

func (m *Matrix[T]) Update(gtx layout.Context, debug bool) {
	if m.drag == nil {
		m.drag = &gesturex.Drag{Tag: m}
	}

	pos := f32.Pt(float32(gtx.Dp(unit.Dp(m.Pos.X))), float32(gtx.Dp(unit.Dp(m.Pos.Y))))
	size := f32.Pt(m.Size.X, m.Size.Y)

	posPt := pos.Round()
	sizePt := size.Round()
	ma := image.Rect(posPt.X, posPt.Y, posPt.X+sizePt.X, posPt.Y+sizePt.Y)
	if debug {
		ma.Min = ma.Min.Sub(image.Pt(10, 10))
		ma.Max = ma.Max.Add(image.Pt(10, 10))
		cl := clip.Stroke{Path: clip.RRect{Rect: ma}.Path(gtx.Ops), Width: 3}.Op().Push(gtx.Ops)
		paint.ColorOp{Color: color.NRGBA{R: 120, G: 12, B: 12, A: 255}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		cl.Pop()
	}
	stack := clip.Rect(ma).Push(gtx.Ops)
	m.drag.Add(gtx.Ops)

	m.drag.Events(gtx.Metric, gtx.Ops, gtx.Queue, func(diff f32.Point) {
		scaledDiff := diff.Div(float32(gtx.Dp(1)))
		m.Pos = m.Pos.Sub(scaledDiff)
	})
	stack.Pop()
}
