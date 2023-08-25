package gesturex

import (
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/unit"
)

// Drag detects drag gestures in the form of pointer.Drag events.
type Drag struct {
	pid pointer.ID
	ptr pointer.Cursor
	pressed,
	dragging bool
	start f32.Point
}

// Add the handler to the operation list to receive drag events.
func (d *Drag) Add(ops *op.Ops) {
	pointer.InputOp{
		Tag:   d,
		Types: pointer.Press | pointer.Enter | pointer.Leave | pointer.Drag | pointer.Move | pointer.Release,
	}.Add(ops)
}

// Events returns the next drag events, if any.
func (d *Drag) Events(cfg unit.Metric, ops *op.Ops, q event.Queue, diffUpdated func(diff f32.Point)) {
	for _, e := range q.Events(d) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}

		d.ptr = pointer.CursorGrab

		switch e.Type {
		case pointer.Press:
			if !(e.Buttons == pointer.ButtonPrimary || e.Source == pointer.Touch) {
				continue
			}

			d.pressed = true
			d.ptr = pointer.CursorGrabbing
			d.start = e.Position
		case pointer.Move:
			d.start = e.Position
		case pointer.Drag:
			d.dragging = d.pressed
			if d.dragging {
				d.ptr = pointer.CursorGrabbing
				diff := d.start.Sub(e.Position)
				diffUpdated(diff)
			}
			d.start = e.Position
		case pointer.Release, pointer.Cancel:
			d.pressed = false
			d.ptr = pointer.CursorGrab
		}
	}
	pointer.Cursor.Add(d.ptr, ops)
}

// Dragging reports whether it is currently in use.
func (d *Drag) Dragging() bool { return d.dragging }

// Pressed returns whether a pointer is pressing.
func (d *Drag) Pressed() bool { return d.pressed }
