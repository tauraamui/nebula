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
)

type Toolbar struct {
	Pos, Size  f32.Point
	TestButton widget.Clickable
}

func (t *Toolbar) Layout(gtx layout.Context, th *material.Theme, debug bool) layout.Dimensions {
	pos := t.Pos
	posX := gtx.Dp(unit.Dp(pos.X))
	posY := gtx.Dp(unit.Dp(pos.Y))
	size := t.Size
	background := image.Rect(posX, posY, posX+gtx.Dp(unit.Dp(size.X)), posY+gtx.Dp(unit.Dp(size.Y)))

	rounded := gtx.Dp(5)
	bgClip := clip.RRect{Rect: background, NE: rounded, SE: rounded, SW: rounded, NW: rounded}.Push(gtx.Ops)

	paint.ColorOp{Color: color.NRGBA{50, 50, 50, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	off := op.Offset(image.Pt(t.Pos.Round().X, t.Pos.Round().Y)).Push(gtx.Ops)
	btnClip := clip.RRect{Rect: image.Rect(gtx.Dp(6), gtx.Dp(6), gtx.Dp(75), gtx.Dp(unit.Dp(t.Size.Y))-gtx.Dp(6)), NE: rounded, SE: rounded, SW: rounded, NW: rounded}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{125, 10, 210, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	btnClip.Pop()
	off.Pop()

	bgClip.Pop()

	return layout.Dimensions{}
}

func (t *Toolbar) Update(gtx layout.Context, offset f32.Point, debug bool) {

}
