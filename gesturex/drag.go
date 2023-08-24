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
