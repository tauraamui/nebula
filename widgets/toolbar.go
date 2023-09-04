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
	selectionArea := image.Rect(posx, posy, posx+gtx.Dp(unit.Dp(size.X)), posy+gtx.Dp(unit.Dp(size.Y)))
	selectionClip := clip.Rect{Min: selectionArea.Min, Max: selectionArea.Max}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{10, 10, 10, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	selectionClip.Pop()
	return layout.Dimensions{}
}

func (t *Toolbar) Update(gtx layout.Context, offset f32.Point, debug bool) {

}
