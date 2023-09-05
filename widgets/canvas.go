package widgets

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/tauraamui/nebula/f32x"
	"github.com/tauraamui/nebula/gesturex"
	"gonum.org/v1/gonum/mat"
)

type Canvas struct {
	debug                  bool
	toolbar                Toolbar
	matrices               []*Matrix[float64]
	theme                  *material.Theme
	input                  *gesturex.InputEvents
	offset                 f32.Point
	pendingSelectionBounds f32x.Rectangle
}

func NewCanvas() *Canvas {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	return &Canvas{
		theme:   th,
		toolbar: Toolbar{Pos: f32.Pt(10, 10), Size: f32.Pt(400, 55)},
		matrices: []*Matrix[float64]{
			{
				Pos:           f32.Pt(200, 200),
				SelectedCells: []image.Point{image.Pt(0, 0)},
				Color:         color.NRGBA{R: 245, G: 245, B: 245, A: 255},
				Data: mat.NewDense(4, 3, []float64{
					12, 353, 11,
					87, 258, 93,
					29, 679, 224,
					229, 6945, 685,
				}),
			},
		},
	}
}

func (c *Canvas) Update(ops *op.Ops, e system.FrameEvent) {
	gtx := layout.NewContext(ops, e)

	key.InputOp{
		Tag: "root",
	}.Add(gtx.Ops)
	for _, e := range gtx.Queue.Events("root") {
		if ke, ok := e.(key.Event); ok {
			if ke.State == key.Press {
				if strings.EqualFold(ke.Name, "x") {
					c.debug = !c.debug
				}
			}
		}
	}

	dpScale := gtx.Dp(1)
	zoomLevelPx := float32(dpScale / dpScale)
	zoomLevelPx = zoomLevelPx - (zoomLevelPx * .1)
	scale := op.Affine(f32.Affine2D{}.Scale(f32.Point{}, f32.Point{X: float32(zoomLevelPx), Y: float32(zoomLevelPx)})).Push(gtx.Ops)

	paint.ColorOp{Color: color.NRGBA{R: 18, G: 18, B: 18, A: 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	if c.input == nil {
		c.input = &gesturex.InputEvents{
			Tag: c,
		}
	}
	ma := image.Rect(0, 0, e.Size.X, e.Size.Y)
	stack := clip.Rect(ma).Push(gtx.Ops)
	c.input.Add(gtx.Ops)
	c.input.Events(gtx.Metric, gtx.Ops, gtx.Queue, c.pressEvents(gtx.Dp), c.releaseEvents(gtx.Dp), c.primaryButtonDragEvents(gtx.Dp), c.secondaryButtonDragEvents(gtx.Dp))
	stack.Pop()

	th := c.theme
	for _, m := range c.matrices {
		m.Layout(gtx, th, c.offset, c.debug)
		m.Update(gtx, c.offset, c.debug)
	}

	selectionBounds := c.pendingSelectionBounds.SwappedBounds()
	if !selectionBounds.Empty() {
		renderPendingSelectionSpan(gtx, 0, 0, selectionBounds, color.NRGBA{50, 110, 220, 80})
	}

	scale.Pop()

	c.toolbar.Pos.X = float32(e.Size.X)/2 - c.toolbar.Size.X/2
	//c.toolbar.Size.X = float32(e.Size.X) * .4
	c.toolbar.Layout(gtx, th, c.debug)
}

func (c *Canvas) pressEvents(dp func(v unit.Dp) int) func(pos f32.Point, buttons pointer.Buttons) {
	return func(pos f32.Point, buttons pointer.Buttons) {
		if buttons != pointer.ButtonPrimary {
			return
		}

		pos = pos.Div(float32(dp(1)))
		c.pendingSelectionBounds = f32x.Rectangle{Min: f32.Pt(pos.X, pos.Y)}
		c.pendingSelectionBounds.Max = c.pendingSelectionBounds.Min
	}
}

func (c *Canvas) releaseEvents(dp func(v unit.Dp) int) func(pos f32.Point, buttons pointer.Buttons) {
	return func(pos f32.Point, buttons pointer.Buttons) {
		if buttons == pointer.ButtonPrimary {
			selectionArea := c.pendingSelectionBounds.SwappedBounds()
			if !selectionArea.Empty() {
				c.pendingSelectionBounds = f32x.Rectangle{}
				return
			}
		}
	}
}

func (c *Canvas) primaryButtonDragEvents(dp func(v unit.Dp) int) func(diff f32.Point) {
	return func(diff f32.Point) {
		scaledDiff := diff.Div(float32(dp(1)))
		c.pendingSelectionBounds.Max = c.pendingSelectionBounds.Max.Add(scaledDiff)

	}
}

func (c *Canvas) secondaryButtonDragEvents(dp func(v unit.Dp) int) func(diff f32.Point) {
	return func(diff f32.Point) {
		scaledDiff := diff.Div(float32(dp(1)))
		c.offset = c.offset.Add(scaledDiff)
	}
}
