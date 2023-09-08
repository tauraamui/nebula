package widgets

import (
	"image"
	"image/color"
	"log"
	"strings"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/tauraamui/nebula/context"
	"github.com/tauraamui/nebula/f32x"
	"github.com/tauraamui/nebula/gesturex"
	"gonum.org/v1/gonum/mat"
)

type Canvas struct {
	debug                  bool
	toolbar                *Toolbar
	matrices               []*Matrix[float64]
	theme                  *material.Theme
	input                  *gesturex.InputEvents
	offset                 f32.Point
	pendingSelectionBounds f32x.Rectangle
}

func NewCanvas() *Canvas {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	tlbar, err := NewToolbar(f32.Pt(300, 40))
	if err != nil {
		log.Fatalf("unable to load toolbar: %v\n", err)
	}

	return &Canvas{
		theme:   th,
		toolbar: tlbar,
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
	gtx := context.NewContext(ops, e)

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
	c.input.Events(gtx.Metric, gtx.Ops, gtx.Queue, nil, nil, nil, c.secondaryButtonDragEvents(gtx.Dp))
	activeTool := c.toolbar.GetActiveTool()
	activeTool.Update(gtx)
	stack.Pop()

	scale := op.Affine(f32.Affine2D{}.Scale(f32.Point{}, f32.Point{X: float32(zoomLevelPx), Y: float32(zoomLevelPx)})).Push(gtx.Ops)

	th := c.theme
	canvasOff := op.Offset(image.Pt(gtx.Dp(unit.Dp(c.offset.Round().X)), gtx.Dp(unit.Dp(c.offset.Round().Y)))).Push(gtx.Ops)
	for _, m := range c.matrices {
		m.Layout(gtx, th, c.debug)
		m.Update(gtx.Context, c.debug)
	}
	canvasOff.Pop()

	selectionBounds := c.pendingSelectionBounds.SwappedBounds()
	if !selectionBounds.Empty() {
		renderPendingSelectionSpan(gtx, selectionBounds, color.NRGBA{50, 110, 220, 80})
	}

	scale.Pop()

	off := op.Offset(image.Pt((e.Size.X/2)-gtx.Dp(unit.Dp(c.toolbar.Size.Round().X))/2, gtx.Dp(10))).Push(gtx.Ops)
	c.toolbar.Layout(gtx.Context, th, c.debug)
	off.Pop()

	for _, e := range gtx.Events() {
		switch evt := e.(type) {
		case context.CreateMatrix:
			c.matrices = append(c.matrices, &Matrix[float64]{
				Pos:   evt.Pos.Div(float32(zoomLevelPx)),
				Color: color.NRGBA{R: 245, G: 245, B: 245, A: 255},
				Data:  mat.NewDense(evt.Rows, evt.Cols, make([]float64, evt.Rows*evt.Cols)),
			})
		}
	}
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
