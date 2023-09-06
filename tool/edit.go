package tool

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/tauraamui/nebula/f32x"
	"github.com/tauraamui/nebula/gesturex"
)

type Edit struct {
	input                  *gesturex.InputEvents
	pendingSelectionBounds f32x.Rectangle
}

func (e *Edit) Update(gtx layout.Context) {
	if e.input == nil {
		e.input = &gesturex.InputEvents{
			Tag: e,
		}
	}

	e.input.Add(gtx.Ops)
	e.input.Events(gtx.Metric, gtx.Ops, gtx.Queue, e.pressEvents(gtx.Dp), e.releaseEvents(gtx.Dp), e.primaryButtonDragEvents(gtx.Dp), nil)

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

func (e *Edit) releaseEvents(dp func(v unit.Dp) int) func(pos f32.Point, buttons pointer.Buttons) {
	return func(pos f32.Point, buttons pointer.Buttons) {
		if buttons == pointer.ButtonPrimary {
			selectionArea := e.pendingSelectionBounds.SwappedBounds()
			if !selectionArea.Empty() {
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

func renderPendingSelectionSpan(gtx layout.Context, posx, posy int, span f32x.Rectangle, color color.NRGBA) {
	selectionArea := image.Rect(posx+gtx.Dp(unit.Dp(span.Min.X)), posy+gtx.Dp(unit.Dp(span.Min.Y)), posx+gtx.Dp(unit.Dp(span.Max.X)), posy+gtx.Dp(unit.Dp(span.Max.Y)))
	selectionClip := clip.Rect{Min: selectionArea.Min, Max: selectionArea.Max}.Push(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	selectionClip.Pop()
}
