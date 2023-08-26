package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/tauraamui/nebula/widgets"
	"gonum.org/v1/gonum/mat"
)

func main() {
	go func() {
		w := app.NewWindow()
		err := loop(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	/*
		m := &widgets.Matrix{
			Pos:   f32.Pt(20, 20),
			Color: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
			Cells: [][][]byte{{[]byte("A"), []byte("A2"), []byte("A3")}, {[]byte("B"), []byte("B2"), []byte("B3")}},
		}
	*/

	m := &widgets.Matrix{
		Pos:   f32.Pt(100, 200),
		Color: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
		Data: mat.NewDense(4, 1, []float64{
			3,
			9,
			12,
			48,
		}),
	}

	m2 := &widgets.Matrix{
		Pos:   f32.Pt(200, 200),
		Color: color.NRGBA{R: 110, G: 0xff, B: 0xff, A: 255},
		Data: mat.NewDense(4, 3, []float64{
			12, 353, 11,
			87, 258, 93,
			29, 679, 224,
			229, 6945, 685,
		}),
	}

	c1 := mat.Col(nil, 0, m.Data)

	c2 := mat.Col(nil, 1, m2.Data)
	for i := 0; i < len(c1); i++ {
		c1[i] += c2[i]
	}

	c1r, _ := m.Data.Dims()
	m3 := &widgets.Matrix{
		Pos:   f32.Pt(460, 200),
		Color: color.NRGBA{R: 225, G: 110, B: 90, A: 255},
		Data:  mat.NewDense(c1r, 1, c1),
	}

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ops.Reset()
			gtx := layout.NewContext(&ops, e)

			paint.ColorOp{Color: color.NRGBA{R: 18, G: 18, B: 18, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			m.Layout(gtx, th)
			m.Update(gtx)
			m2.Layout(gtx, th)
			m2.Update(gtx)
			m3.Layout(gtx, th)
			m3.Update(gtx)

			e.Frame(gtx.Ops)
		}
	}
}
