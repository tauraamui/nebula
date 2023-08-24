package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

const (
	cellWidth   = 130
	cellHeight  = 30
	cellPadding = 2
)

type Matrix struct {
	Pos,
	Size image.Point
	Color color.NRGBA
	Cells [][]int
}

func (m *Matrix) Layout(gtx layout.Context) layout.Dimensions {
	totalSize := image.Point{}
	totalX := 0
	totalY := 0
	cellSize := image.Point{X: cellWidth, Y: cellHeight}
	for x, column := range m.Cells {
		totalX += 1
		for y := range column {
			if totalX == 1 {
				totalY += 1
			}
			cell := image.Rect(m.Pos.X+(cellWidth*x)+cellPadding, m.Pos.Y+(y*cellHeight)+cellPadding, m.Pos.X+((cellWidth*x)+cellWidth), m.Pos.Y+((cellHeight*y)+cellHeight))
			cell.Min = cell.Min.Add(image.Pt(cellPadding, cellPadding))
			cell.Max = cell.Max.Add(image.Pt(cellPadding, cellPadding))
			cl := clip.Rect{Min: cell.Min, Max: cell.Max}.Push(gtx.Ops)
			paint.ColorOp{Color: m.Color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			cl.Pop()
		}
	}
	totalSize.X = totalX * cellSize.X
	totalSize.Y = totalY * cellSize.Y
	m.Size = totalSize
	return layout.Dimensions{Size: m.Size}
}

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
	m := &Matrix{
		Pos:   image.Pt(20, 20),
		Color: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 255},
		Cells: [][]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
	}

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	var drag Drag
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
			m.Layout(gtx)

			ma := image.Rect(m.Pos.X, m.Pos.Y, m.Pos.X+m.Size.X, m.Pos.Y+m.Size.Y)
			stack := clip.Rect(ma).Push(gtx.Ops)
			drag.Add(gtx.Ops)
			stack.Pop()

			drag.Events(unit.Metric{PxPerDp: 1, PxPerSp: 1}, gtx.Queue, func(diff f32.Point) {
				m.Pos = m.Pos.Sub(image.Pt(drag.diff.Round().X, drag.diff.Round().Y))
			})

			e.Frame(gtx.Ops)
		}
	}
}

// Drag detects drag gestures in the form of pointer.Drag events.
type Drag struct {
	diff        f32.Point
	dragging    bool
	lastDragPos f32.Point
	pressed     bool
	pid         pointer.ID
	start       f32.Point
}

// Add the handler to the operation list to receive drag events.
func (d *Drag) Add(ops *op.Ops) {
	pointer.InputOp{
		Tag:   d,
		Types: pointer.Press | pointer.Drag | pointer.Release,
	}.Add(ops)
}

// Events returns the next drag events, if any.
func (d *Drag) Events(cfg unit.Metric, q event.Queue, diffUpdated func(diff f32.Point)) []pointer.Event {
	var events []pointer.Event
	for _, e := range q.Events(d) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}

		switch e.Type {
		case pointer.Press:
			if !(e.Buttons == pointer.ButtonPrimary || e.Source == pointer.Touch) {
				continue
			}
			d.pressed = true
			if d.dragging {
				continue
			}
			d.dragging = true
			d.pid = e.PointerID
			d.start = e.Position
			d.lastDragPos = d.start
		case pointer.Drag:
			if !d.dragging || e.PointerID != d.pid {
				continue
			}
			diff := d.lastDragPos.Sub(e.Position)
			d.diff = diff
			d.lastDragPos = e.Position
			diffUpdated(d.diff)
		case pointer.Release, pointer.Cancel:
			d.pressed = false
			if !d.dragging || e.PointerID != d.pid {
				continue
			}
			d.dragging = false
		}

		events = append(events, e)
	}

	return events
}

// Dragging reports whether it is currently in use.
func (d *Drag) Dragging() bool { return d.dragging }

// Pressed returns whether a pointer is pressing.
func (d *Drag) Pressed() bool { return d.pressed }
