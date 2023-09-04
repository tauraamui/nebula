package widgets

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
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
	rounded := gtx.Dp(5)
	background := image.Rect(posX, posY, posX+gtx.Dp(unit.Dp(size.X)), posY+gtx.Dp(unit.Dp(size.Y)))
	bgClip := clip.RRect{Rect: background, NE: rounded, SE: rounded, SW: rounded, NW: rounded}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{50, 50, 50, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	bgClip.Pop()

	return layout.Dimensions{}
}

func (t *Toolbar) Update(gtx layout.Context, offset f32.Point, debug bool) {

}
