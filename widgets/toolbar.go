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
	"gioui.org/widget/material"
	"github.com/inkeliz/giosvg"
	"github.com/tauraamui/nebula/gesturex"
	"github.com/tauraamui/nebula/icons"
	"github.com/tauraamui/nebula/tool"
)

type Toolbar struct {
	Size   f32.Point
	tools  []tool.Tool
	btns   []*toolButton
	active int
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

	paint.ColorOp{Color: color.NRGBA{7, 7, 7, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	off := op.Offset(image.Pt(gtx.Dp(5), gtx.Dp(5))).Push(gtx.Ops)
	for i, btn := range t.btns {
		var btnoff op.TransformStack
		if i > 0 {
			btnoff = op.Offset(image.Pt(gtx.Dp(unit.Dp((5+btn.size.X)))*i, 0)).Push(gtx.Ops)
		}
		btn.Layout(gtx, t.Size.Y, i == t.active)
		btn.Update(i, gtx, debug, i == t.active, t.buttonClicked(i))
		if i > 0 {
			btnoff.Pop()
		}
	}
	off.Pop()

	bgClip.Pop()

	return layout.Dimensions{}
}

func (t *Toolbar) GetActiveTool() tool.Tool {
	return t.btns[t.active].tool
}

func (t *Toolbar) buttonClicked(index int) func() {
	return func() {
		t.active = index
	}
}

type toolButton struct {
	tool                     tool.Tool
	inactiveIcon, activeIcon *giosvg.Icon
	size                     f32.Point
	rounded                  int
	inputEvents              *gesturex.ButtonEvents
	beingPressed             bool
}

func (b *toolButton) Layout(gtx layout.Context, barHeight float32, active bool) layout.Dimensions {
	btn := image.Rect(0, 0, gtx.Dp(unit.Dp(b.size.X)), (gtx.Dp(unit.Dp(barHeight)) - gtx.Dp(10)))
	rounded := gtx.Dp(unit.Dp(b.rounded))
	btnClip := clip.RRect{Rect: btn, NE: rounded, SE: rounded, SW: rounded, NW: rounded}

	if b.beingPressed {
		cl := clip.Stroke{Path: btnClip.Path(gtx.Ops), Width: 3}.Op().Push(gtx.Ops)
		paint.ColorOp{Color: color.NRGBA{172, 155, 238, 255}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		cl.Pop()
	}

	if active {
		cl := btnClip.Push(gtx.Ops)
		paint.ColorOp{Color: color.NRGBA{172, 155, 238, 255}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		cl.Pop()

		if b.activeIcon != nil {
			drawIcon(gtx, b.activeIcon, active)
		}
	} else {
		if b.inactiveIcon != nil {
			drawIcon(gtx, b.inactiveIcon, active)
		}
	}

	return layout.Dimensions{}
}

func drawIcon(gtx layout.Context, ic *giosvg.Icon, active bool) {
	iconOff := op.Offset(image.Pt(gtx.Dp(7), gtx.Dp(7))).Push(gtx.Ops)
	gtx.Constraints.Min = image.Pt(gtx.Dp(10), gtx.Dp(10))
	gtx.Constraints.Max = image.Pt(gtx.Dp(100), gtx.Dp(16))
	iconcolor := color.NRGBA{255, 255, 255, 255}
	if active {
		iconcolor = color.NRGBA{82, 29, 228, 255}
	}
	paint.ColorOp{Color: iconcolor}.Add(gtx.Ops)
	ic.Layout(gtx)
	iconOff.Pop()
}

func (b *toolButton) Update(index int, gtx layout.Context, debug, active bool, clicked func()) {
	if b.inputEvents == nil {
		b.inputEvents = &gesturex.ButtonEvents{Tag: b}
	}

	if active {
		return
	}

	btn := image.Rect(0, 0, gtx.Dp(unit.Dp(b.size.X)), (gtx.Dp(unit.Dp(30))))
	stack := clip.Rect(btn).Push(gtx.Ops)
	b.inputEvents.Add(gtx.Ops)
	b.inputEvents.Events(gtx.Metric, gtx.Ops, gtx.Queue, func() {
		b.beingPressed = true
	}, func() {
		b.beingPressed = false
	}, func() {
		b.beingPressed = false
		if !active {
			clicked()
		}
	})
	stack.Pop()
}

func makeButton(inactiveIcon, activeIcon icons.IconResolver) (*toolButton, error) {
	inactive, err := inactiveIcon()
	if err != nil {
		return nil, err
	}

	active, err := activeIcon()
	if err != nil {
		return nil, err
	}
	return &toolButton{size: f32.Pt(30, 0), inactiveIcon: inactive, activeIcon: active, rounded: 10}, nil
}

func makeAllButtons() ([]*toolButton, error) {
	btns := []*toolButton{}

	pointAndSelect, err := makeButton(icons.MousePointerOutline, icons.MousePointer)
	if err != nil {
		return nil, err
	}
	pointAndSelect.tool = &tool.Select{}

	drawNewMatrix, err := makeButton(icons.SquareBorder, icons.Square)
	if err != nil {
		return nil, err
	}
	drawNewMatrix.tool = &tool.Edit{}

	btns = append(btns, pointAndSelect)
	btns = append(btns, drawNewMatrix)

	return btns, nil
}
