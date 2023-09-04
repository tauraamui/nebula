package widgets

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Toolbar struct {
	Pos, Size f32.Point
}

func (t *Toolbar) Layout(gtx layout.Context, th *material.Theme, debug bool) layout.Dimensions {
	posx, posy := t.Pos.Round().X, t.Pos.Round().Y
	size := t.Size
	background := image.Rect(posx, posy, posx+gtx.Dp(unit.Dp(size.X)), posy+gtx.Dp(unit.Dp(size.Y)))
	bgClip := clip.RRect{Rect: background, NE: 5, SE: 5, SW: 5, NW: 5}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{50, 50, 50, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	bgClip.Pop()
	return layout.Dimensions{}
}

func (t *Toolbar) Update(gtx layout.Context, offset f32.Point, debug bool) {

}
