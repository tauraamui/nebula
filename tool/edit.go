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
	input                  *gesturex.InputEvents
	pendingSelectionBounds f32x.Rectangle
}

func (e *Edit) Update(gtx *context.Context) {
	if e.input == nil {
		e.input = &gesturex.InputEvents{
			Tag: e,
		}
	}

	e.input.Add(gtx.Ops)
	e.input.Events(gtx.Metric, gtx.Ops, gtx.Queue, e.pressEvents(gtx.Dp), e.releaseEvents(gtx.Dp, gtx.PushEvent), e.primaryButtonDragEvents(gtx.Dp), nil)

	selectionBounds := e.pendingSelectionBounds.SwappedBounds()
	if !selectionBounds.Empty() {
		renderPendingSelectionSpan(gtx, 0, 0, selectionBounds, color.NRGBA{50, 110, 220, 80})
	}

}

func (e *Edit) pressEvents(dp func(v unit.Dp) int) func(pos f32.Point, buttons pointer.Buttons) {
	return func(pos f32.Point, buttons pointer.Buttons) {
		if buttons != pointer.ButtonPrimary {
			return
		}

		pos = pos.Div(float32(dp(1)))
		e.pendingSelectionBounds = f32x.Rectangle{Min: f32.Pt(pos.X, pos.Y)}
		e.pendingSelectionBounds.Max = e.pendingSelectionBounds.Min
	}
}

func (e *Edit) releaseEvents(dp func(v unit.Dp) int, pushEvent func(e any)) func(pos f32.Point, buttons pointer.Buttons) {
	return func(pos f32.Point, buttons pointer.Buttons) {
		if buttons == pointer.ButtonPrimary {
			selectionArea := e.pendingSelectionBounds.SwappedBounds()
			if !selectionArea.Empty() {
				// TODO:(tauraamui) -> implement pushing of "create matrix event" which will be read by canvas and actioned
				pushEvent(context.CreateMatrix{
					Pos:    selectionArea.Min,
					Bounds: selectionArea,
				})
				e.pendingSelectionBounds = f32x.Rectangle{}
				return
			}
		}
	}
}

func (e *Edit) primaryButtonDragEvents(dp func(v unit.Dp) int) func(diff f32.Point) {
	return func(diff f32.Point) {
		scaledDiff := diff.Div(float32(dp(1)))
		e.pendingSelectionBounds.Max = e.pendingSelectionBounds.Max.Add(scaledDiff)

	}
}

func renderPendingSelectionSpan(gtx *context.Context, posx, posy int, span f32x.Rectangle, bgcolor color.NRGBA) {
	selectionArea := image.Rect(posx+gtx.Dp(unit.Dp(span.Min.X)), posy+gtx.Dp(unit.Dp(span.Min.Y)), posx+gtx.Dp(unit.Dp(span.Max.X)), posy+gtx.Dp(unit.Dp(span.Max.Y)))
	selectionClip := clip.Rect{Min: selectionArea.Min, Max: selectionArea.Max}.Push(gtx.Ops)
	paint.ColorOp{Color: bgcolor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	// render pending matrix cells

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			renderCell(gtx, x, y, posx+gtx.Dp(unit.Dp(span.Min.X)), posy+gtx.Dp(unit.Dp(span.Min.Y)), gtx.Dp(unit.Dp(cellSize.X)), gtx.Dp(unit.Dp(cellSize.Y)), color.NRGBA{255, 255, 255, 255})
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
