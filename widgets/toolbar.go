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
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/inkeliz/giosvg"
)

type Toolbar struct {
	Pos, Size        f32.Point
	TestButton       widget.Clickable
	MousePointerIcon *giosvg.Icon
}

func (t *Toolbar) Layout(gtx layout.Context, th *material.Theme, debug bool) layout.Dimensions {
	pos := t.Pos
	posX := gtx.Dp(unit.Dp(pos.X))
	posY := gtx.Dp(unit.Dp(pos.Y))
	size := t.Size.Round()
	background := image.Rect(posX, posY, posX+gtx.Dp(unit.Dp(size.X)), posY+gtx.Dp(unit.Dp(size.Y)))

	rounded := gtx.Dp(8)
	bgClip := clip.RRect{Rect: background, NE: rounded, SE: rounded, SW: rounded, NW: rounded}.Push(gtx.Ops)

	paint.ColorOp{Color: color.NRGBA{50, 50, 50, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	off := op.Offset(image.Pt(posX+gtx.Dp(5), posY+gtx.Dp(5))).Push(gtx.Ops)
	btn := image.Rect(0, 0, gtx.Dp(30), (gtx.Dp(unit.Dp(size.Y)) - gtx.Dp(10)))
	btnClip := clip.RRect{Rect: btn, NE: rounded, SE: rounded, SW: rounded, NW: rounded}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{110, 50, 180, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	btnClip.Pop()

	iconOff := op.Offset(image.Pt(gtx.Dp(5), gtx.Dp(5))).Push(gtx.Ops)
	gtx.Constraints.Min = image.Point{}
	gtx.Constraints.Max = image.Pt(gtx.Dp(20), gtx.Dp(20))
	paint.ColorOp{Color: color.NRGBA{200, 200, 200, 255}}.Add(gtx.Ops)
	t.MousePointerIcon.Layout(gtx)
	iconOff.Pop()

	off.Pop()

	bgClip.Pop()

	return layout.Dimensions{}
}

func (t *Toolbar) Update(gtx layout.Context, offset f32.Point, debug bool) {

}
