package tool

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/tauraamui/nebula/context"
	"github.com/tauraamui/nebula/f32x"
	"github.com/tauraamui/nebula/gesturex"
)

const (
	cellWidth  unit.Dp = 80
	cellHeight unit.Dp = 25
	cols       int     = 15
	rows       int     = 10
)

var cellSize f32.Point = f32.Pt(float32(cellWidth), float32(cellHeight))

type Edit struct {
	input                 *gesturex.InputEvents
	pendingCreationBounds f32x.Rectangle
}

func (e *Edit) Update(gtx *context.Context) {
	if e.input == nil {
		e.input = &gesturex.InputEvents{
			Tag: e,
		}
	}

	e.input.Add(gtx.Ops)
	e.input.Events(gtx.Metric, gtx.Ops, gtx.Queue, e.pressEvents(gtx.Dp), e.releaseEvents(gtx.Dp, gtx.PushEvent), e.primaryButtonDragEvents(gtx.Dp), nil)

	selectionBounds := e.pendingCreationBounds.SwappedBounds()
	if !selectionBounds.Empty() {
		renderPendingCreationSpan(gtx, selectionBounds, color.NRGBA{50, 110, 220, 80})
	}

}

func (e *Edit) pressEvents(dp func(v unit.Dp) int) func(pos f32.Point, buttons pointer.Buttons) {
	return func(pos f32.Point, buttons pointer.Buttons) {
		if buttons != pointer.ButtonPrimary {
			return
		}

		pos = pos.Div(float32(dp(1)))
		e.pendingCreationBounds = f32x.Rectangle{Min: f32.Pt(pos.X, pos.Y)}
		e.pendingCreationBounds.Max = e.pendingCreationBounds.Min
	}
}

func (e *Edit) releaseEvents(dp func(v unit.Dp) int, pushEvent func(e any)) func(pos f32.Point, buttons pointer.Buttons) {
	return func(pos f32.Point, buttons pointer.Buttons) {
		if buttons == pointer.ButtonPrimary {
			selectionArea := e.pendingCreationBounds.SwappedBounds()
			if !selectionArea.Empty() {
				pushEvent(context.CreateMatrix{
					Pos:  selectionArea.Min,
					Rows: 5,
					Cols: 5,
				})
				e.pendingCreationBounds = f32x.Rectangle{}
				return
			}
		}
	}
}

func (e *Edit) primaryButtonDragEvents(dp func(v unit.Dp) int) func(diff f32.Point) {
	return func(diff f32.Point) {
		scaledDiff := diff.Div(float32(dp(1)))
		e.pendingCreationBounds.Max = e.pendingCreationBounds.Max.Add(scaledDiff)

	}
}

func renderPendingCreationSpan(gtx *context.Context, span f32x.Rectangle, bgcolor color.NRGBA) {
	selectionArea := image.Rect(gtx.Dp(unit.Dp(span.Min.X)), gtx.Dp(unit.Dp(span.Min.Y)), gtx.Dp(unit.Dp(span.Max.X)), gtx.Dp(unit.Dp(span.Max.Y)))
	selectionClip := clip.Rect{Min: selectionArea.Min, Max: selectionArea.Max}.Push(gtx.Ops)
	paint.ColorOp{Color: bgcolor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	// render pending matrix cells

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			renderCell(gtx, x, y, gtx.Dp(unit.Dp(span.Min.X)), gtx.Dp(unit.Dp(span.Min.Y)), gtx.Dp(unit.Dp(cellSize.X)), gtx.Dp(unit.Dp(cellSize.Y)), color.NRGBA{255, 255, 255, 255})
		}
	}

	selectionClip.Pop()
}

func renderCell(gtx *context.Context, x, y int, posx, posy, cellwidth, cellheight int, bgcolor color.NRGBA) {
	// render background of cell
	cell := image.Rect(posx+(cellwidth*x), posy+(y*cellheight), posx+((cellwidth*x)+cellwidth), posy+((cellheight*y)+cellheight))
	cl1 := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
	paint.ColorOp{Color: bgcolor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl1.Pop()

	// render cell border
	borderWidth := float32(.25) / float32(gtx.Dp(1))
	borderColor := color.NRGBA{R: 55, G: 55, B: 55, A: 255}
	cl3 := clip.Stroke{Path: clip.RRect{Rect: cell}.Path(gtx.Ops), Width: borderWidth}.Op().Push(gtx.Ops)
	paint.ColorOp{Color: borderColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl3.Pop()
}

func calcRowsAndColsFromSpan(span f32x.Rectangle) (int, int) {
	return 0, 0
}
