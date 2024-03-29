package widgets

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/tauraamui/nebula/context"
	"github.com/tauraamui/nebula/f32x"
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
	Color                  color.NRGBA
	Data                   *mat.Dense
	Data2                  nmat.Matrix[T]
	cellSize               f32.Point
	inputEvents            *gesturex.InputEvents
	selectedCell           image.Point
	SelectedCells          []image.Point
	pendingSelectionBounds f32x.Rectangle
	wasMovingMinLast       bool
	cachedOps              *op.Ops
	call                   op.CallOp
}

func (m *Matrix[T]) Layout(gtx *context.Context, th *material.Theme, debug bool) layout.Dimensions {
	m.cellSize.X = float32(cellWidth)
	m.cellSize.Y = float32(cellHeight)

	off := op.Offset(image.Pt(gtx.Dp(unit.Dp(m.Pos.X)), gtx.Dp(unit.Dp(m.Pos.Y)))).Push(gtx.Ops)

	cellSize := f32.Point{X: float32(gtx.Dp(cellWidth)), Y: float32(gtx.Dp(cellHeight))}

	rows, cols := m.Data.Dims()
	totalSize := f32.Point{
		X: float32(cols) * cellSize.X,
		Y: float32(rows) * cellSize.Y,
	}
	m.Size = totalSize

	bgnd := clip.Rect{Min: image.Pt(0, 0), Max: image.Pt(m.Size.Round().X, m.Size.Round().Y)}.Push(gtx.Ops)
	paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	bgnd.Pop()

	if m.cachedOps == nil {
		m.cachedOps = &op.Ops{}
		macro := op.Record(m.cachedOps)
		cells := clip.Path{}
		cells.Begin(m.cachedOps)

		cellwidth := gtx.Dp(unit.Dp(m.cellSize.X))
		cellheight := gtx.Dp(unit.Dp(m.cellSize.Y))

		for x := 0; x < cols; x++ {
			for y := 0; y < rows; y++ {
				cells.MoveTo(f32.Pt(float32(cellwidth*x), float32(cellheight*y)))
				cells.LineTo(f32.Pt(float32((cellwidth*x)+cellwidth), float32(cellheight*y)))
				cells.LineTo(f32.Pt(float32((cellwidth*x)+cellwidth), float32(cellheight*y+cellheight)))
				cells.LineTo(f32.Pt(float32(cellwidth*x), float32(cellheight*y+cellheight)))
			}
		}
		cells.Close()

		borderWidth := float32(.35) / float32(gtx.Dp(1))
		borderColor := color.NRGBA{R: 55, G: 55, B: 55, A: 255}
		cStroke := clip.Stroke{Path: cells.End(), Width: borderWidth}.Op().Push(m.cachedOps)
		paint.ColorOp{Color: borderColor}.Add(m.cachedOps)
		paint.PaintOp{}.Add(m.cachedOps)
		cStroke.Pop()

		m.call = macro.Stop()
	}

	m.call.Add(gtx.Ops)

	/*
		for x := 0; x < cols; x++ {
			for y := 0; y < rows; y++ {
				renderCell(gtx, strconv.FormatFloat(m.Data.At(y, x), 'f', -1, 64), x, y, gtx.Dp(unit.Dp(m.cellSize.X)), gtx.Dp(unit.Dp(m.cellSize.Y)), m.Color, th)
			}
		}
	*/

	for _, selectedCell := range m.SelectedCells {
		renderCellSelection(gtx, selectedCell.X, selectedCell.Y, gtx.Dp(unit.Dp(m.cellSize.X)), gtx.Dp(unit.Dp(m.cellSize.Y)))
	}

	selectionBounds := m.pendingSelectionBounds.SwappedBounds()
	if !selectionBounds.Empty() {
		area := image.Rect(0, 0, m.Size.Round().X, m.Size.Round().Y)
		clip := clip.Rect{Min: area.Min, Max: area.Max}.Push(gtx.Ops)
		renderPendingSelectionSpan(gtx, selectionBounds, color.NRGBA{224, 63, 222, 110})
		clip.Pop()
	}

	off.Pop()

	return layout.Dimensions{Size: m.Size.Round()}
}

func renderPendingSelectionSpan(gtx *context.Context, span f32x.Rectangle, color color.NRGBA) {
	selectionArea := image.Rect(gtx.Dp(unit.Dp(span.Min.X)), gtx.Dp(unit.Dp(span.Min.Y)), gtx.Dp(unit.Dp(span.Max.X)), gtx.Dp(unit.Dp(span.Max.Y)))
	selectionClip := clip.Rect{Min: selectionArea.Min, Max: selectionArea.Max}.Push(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	selectionClip.Pop()
}

func renderCellSelection(gtx *context.Context, x, y, cellwidth, cellheight int) {
	cell := image.Rect((cellwidth * x), (y * cellheight), ((cellwidth * x) + cellwidth), ((cellheight * y) + cellheight))
	// render cell border
	borderWidth := 2 * float32(gtx.Dp(1))
	borderColor := color.NRGBA{R: 230, G: 90, B: 90, A: 255}
	cl3 := clip.Stroke{Path: clip.RRect{Rect: cell}.Path(gtx.Ops), Width: borderWidth}.Op().Push(gtx.Ops)
	paint.ColorOp{Color: borderColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl3.Pop()
}

func renderCell(gtx *context.Context, content string, x, y, cellwidth, cellheight int, bgcolor color.NRGBA, th *material.Theme) {
	// render background of cell
	cell := image.Rect(cellwidth*x, y*cellheight, ((cellwidth * x) + cellwidth), ((cellheight * y) + cellheight))
	cl1 := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
	paint.ColorOp{Color: bgcolor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl1.Pop()

	// render cell content as text label
	cl2 := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
	l := material.Label(th, unit.Sp(14), content)
	lineHeightPx := gtx.Sp(14)
	l.Color = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
	off := op.Offset(cell.Min.Add(image.Pt(gtx.Sp(3), (cellheight/2)-(lineHeightPx/2)))).Push(gtx.Ops)
	l.Layout(gtx.Context)
	off.Pop()
	cl2.Pop()

	// render cell border
	borderWidth := float32(.25) / float32(gtx.Dp(1))
	borderColor := color.NRGBA{R: 55, G: 55, B: 55, A: 255}
	cl3 := clip.Stroke{Path: clip.RRect{Rect: cell}.Path(gtx.Ops), Width: borderWidth}.Op().Push(gtx.Ops)
	paint.ColorOp{Color: borderColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl3.Pop()
}

func (m *Matrix[T]) Update(gtx layout.Context, debug bool) {
	if m.inputEvents == nil {
		m.inputEvents = &gesturex.InputEvents{Tag: m}
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
	m.inputEvents.Add(gtx.Ops)

	m.inputEvents.Events(gtx.Metric, gtx.Ops, gtx.Queue, m.pressEvents(gtx.Dp), m.releaseEvents(gtx.Dp), m.primaryButtonDragEvents(gtx.Dp), m.secondaryButtonDragEvents(gtx.Dp))
	stack.Pop()
}

func (m *Matrix[T]) pressEvents(dp func(v unit.Dp) int) func(pos f32.Point, buttons pointer.Buttons) {
	return func(pos f32.Point, buttons pointer.Buttons) {
		if buttons != pointer.ButtonPrimary {
			return
		}

		m.makeCellSelection(dp, pos)

		pos = pos.Div(float32(dp(1)))
		// wip pending selection implementation
		m.pendingSelectionBounds = f32x.Rectangle{Min: f32.Pt(pos.X, pos.Y)}
		m.pendingSelectionBounds.Min = m.pendingSelectionBounds.Min.Sub(m.Pos)
		m.pendingSelectionBounds.Max = m.pendingSelectionBounds.Min
	}
}

func (m *Matrix[T]) releaseEvents(dp func(v unit.Dp) int) func(pos f32.Point, buttons pointer.Buttons) {
	return func(pos f32.Point, buttons pointer.Buttons) {
		if buttons == pointer.ButtonPrimary {
			selectionArea := m.pendingSelectionBounds.SwappedBounds()
			if !selectionArea.Empty() {
				m.SelectedCells = resolveSelectedCells(m.Data.Dims())(dp, m.Pos, m.cellSize, selectionArea)
				m.pendingSelectionBounds = f32x.Rectangle{}
				return
			}
			m.SelectedCells = []image.Point{resolvePressedCell(m.Data.Dims())(dp, m.Pos, m.cellSize, pos)}
		}
	}
}

func in(p f32.Point, r f32x.Rectangle) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

func resolvePressedCell(rows, cols int) func(dp func(v unit.Dp) int, pos, cellSize f32.Point, pressPos f32.Point) image.Point {
	return func(dp func(v unit.Dp) int, pos, cellSize f32.Point, pressPos f32.Point) image.Point {
		pressPos = pressPos.Div(float32(dp(1)))

		var x, y float32
		for x = 0; x < float32(cols); x++ {
			for y = 0; y < float32(rows); y++ {
				cell := f32x.Rectangle{Min: f32.Pt(pos.X+(cellSize.X*x), pos.Y+(cellSize.Y*y)), Max: f32.Pt(pos.X+((cellSize.X*x)+cellSize.X), pos.Y+((cellSize.Y*y)+cellSize.Y))}
				if in(pressPos, cell) {
					return image.Pt(int(x), int(y))
				}
			}
		}
		return image.Pt(0, 0)
	}
}

func resolveSelectedCells(rows, cols int) func(dp func(v unit.Dp) int, pos, cellSize f32.Point, selection f32x.Rectangle) []image.Point {
	return func(dp func(v unit.Dp) int, pos, cellSize f32.Point, selection f32x.Rectangle) []image.Point {
		selection.Min = selection.Min.Add(pos)
		selection.Max = selection.Max.Add(pos)

		selectedCells := []image.Point{}
		var x, y float32
		for x = 0; x < float32(cols); x++ {
			for y = 0; y < float32(rows); y++ {
				cell := f32x.Rectangle{Min: f32.Pt(pos.X+(cellSize.X*x), pos.Y+(cellSize.Y*y)), Max: f32.Pt(pos.X+((cellSize.X*x)+cellSize.X), pos.Y+((cellSize.Y*y)+cellSize.Y))}
				if selection.Overlaps(cell) {
					selectedCells = append(selectedCells, image.Pt(int(x), int(y)))
				}
			}
		}
		return selectedCells
	}
}

func (m *Matrix[T]) makeCellSelection(dp func(v unit.Dp) int, pos f32.Point) {
	// make press postion relative to this matrix
	pos = pos.Sub(f32.Pt(m.Pos.X, m.Pos.Y))
	scaledDiff := pos.Div(float32(dp(1)))
	cellx := math.Floor(float64(scaledDiff.X) / float64(m.cellSize.X))
	celly := math.Floor(float64(scaledDiff.Y) / float64(m.cellSize.Y))
	m.selectedCell.X = int(cellx)
	m.selectedCell.Y = int(celly)
}

func (m *Matrix[T]) primaryButtonDragEvents(dp func(v unit.Dp) int) func(diff f32.Point) {
	return func(diff f32.Point) {
		scaledDiff := diff.Div(float32(dp(1)))
		m.pendingSelectionBounds.Max = m.pendingSelectionBounds.Max.Add(scaledDiff)

	}
}

func (m *Matrix[T]) secondaryButtonDragEvents(dp func(v unit.Dp) int) func(diff f32.Point) {
	return func(diff f32.Point) {
		scaledDiff := diff.Div(float32(dp(1)))
		m.Pos = m.Pos.Add(scaledDiff)
	}
}
