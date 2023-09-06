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
	"github.com/tauraamui/nebula/icons"
)

type Toolbar struct {
	Size             f32.Point
	btns             []*toolButton
	moveAndSelect    toolButton
	TestButton       widget.Clickable
	MousePointerIcon *giosvg.Icon
}

func NewToolbar(size f32.Point) (*Toolbar, error) {
	tlbar := Toolbar{Size: size}

	btns, err := makeAllButtons()
	if err != nil {
		return nil, err
	}

	tlbar.btns = btns

	return &tlbar, nil
}

func (t *Toolbar) Layout(gtx layout.Context, th *material.Theme, debug bool) layout.Dimensions {
	size := t.Size.Round()
	background := image.Rect(0, 0, gtx.Dp(unit.Dp(size.X)), gtx.Dp(unit.Dp(size.Y)))

	rounded := gtx.Dp(8)
	bgClip := clip.RRect{Rect: background, NE: rounded, SE: rounded, SW: rounded, NW: rounded}.Push(gtx.Ops)

	paint.ColorOp{Color: color.NRGBA{50, 50, 50, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	off := op.Offset(image.Pt(gtx.Dp(5), gtx.Dp(5))).Push(gtx.Ops)
	off.Pop()

	bgClip.Pop()

	return layout.Dimensions{}
}

type toolButton struct {
	size    f32.Point
	rounded int
	icon    *giosvg.Icon
}

func (b *toolButton) Layout(gtx layout.Context) layout.Dimensions {
	btn := image.Rect(0, 0, gtx.Dp(30), (gtx.Dp(unit.Dp(b.size.Y)) - gtx.Dp(10)))
	btnClip := clip.RRect{Rect: btn, NE: b.rounded, SE: b.rounded, SW: b.rounded, NW: b.rounded}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{110, 50, 180, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	btnClip.Pop()

	if b.icon != nil {
		iconOff := op.Offset(image.Pt(gtx.Dp(7), gtx.Dp(7))).Push(gtx.Ops)
		gtx.Constraints.Min = image.Pt(gtx.Dp(10), gtx.Dp(10))
		gtx.Constraints.Max = image.Pt(gtx.Dp(100), gtx.Dp(16))
		paint.ColorOp{Color: color.NRGBA{200, 200, 200, 255}}.Add(gtx.Ops)
		b.icon.Layout(gtx)
		iconOff.Pop()
	}

	return layout.Dimensions{}
}

func makeButton(icon icons.IconResolver) (*toolButton, error) {
	ic, err := icon()
	if err != nil {
		return nil, err
	}
	return &toolButton{icon: ic}, nil
}

func makeAllButtons() ([]*toolButton, error) {
	btns := []*toolButton{}

	pointAndSelect, err := makeButton(icons.MousePointer)
	if err != nil {
		return nil, err
	}

	btns = append(btns, pointAndSelect)

	return btns, nil
}
