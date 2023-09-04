package widgets

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"gonum.org/v1/gonum/mat"
)

type Canvas struct {
	matrices []*Matrix[float64]
}

func (c *Canvas) Run(w *app.Window) error {
	c.matrices = append(c.matrices, &Matrix[float64]{
		Pos:           f32.Pt(200, 200),
		SelectedCells: []image.Point{image.Pt(0, 0)},
		Color:         color.NRGBA{R: 245, G: 245, B: 245, A: 255},
		Data: mat.NewDense(4, 3, []float64{
			12, 353, 11,
			87, 258, 93,
			29, 679, 224,
			229, 6945, 685,
		}),
	})

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	debug := false
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ops.Reset()

			gtx := layout.NewContext(&ops, e)

			dpScale := gtx.Dp(1)
			zoomLevelPx := float32(dpScale / dpScale)
			zoomLevelPx = zoomLevelPx - (zoomLevelPx * .1)
			scale := op.Affine(f32.Affine2D{}.Scale(f32.Point{}, f32.Point{X: float32(zoomLevelPx), Y: float32(zoomLevelPx)})).Push(gtx.Ops)

			key.InputOp{
				Tag: "root",
			}.Add(gtx.Ops)
			for _, e := range gtx.Queue.Events("root") {
				if ke, ok := e.(key.Event); ok {
					if ke.State == key.Press {
						if strings.EqualFold(ke.Name, "x") {
							debug = !debug
						}
					}
				}
			}

			paint.ColorOp{Color: color.NRGBA{R: 18, G: 18, B: 18, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			for _, m := range c.matrices {
				m.Layout(gtx, th, debug)
				m.Update(gtx, debug)
			}

			scale.Pop()

			e.Frame(gtx.Ops)
		}
	}
}
