package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/tauraamui/nebula/widgets"
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
	m := &widgets.Matrix{
		Pos:   image.Pt(20, 20),
		Color: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
		Cells: [][][]byte{{[]byte("A"), []byte("A2"), []byte("A3")}, {[]byte("B"), []byte("B2"), []byte("B3")}},
	}

	m2 := &widgets.Matrix{
		Pos:   image.Pt(200, 200),
		Color: color.NRGBA{R: 110, G: 0xff, B: 0xff, A: 255},
		Cells: [][][]byte{{[]byte("random"), []byte("0.433"), []byte("shitzyxcfe3kqt")}, {[]byte("A1"), []byte("A3"), []byte("C")}},
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

			paint.ColorOp{Color: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			m.Layout(gtx, th)
			m.Update(gtx)
			m2.Layout(gtx, th)
			m2.Update(gtx)

			e.Frame(gtx.Ops)
		}
	}
}
