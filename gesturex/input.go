package gesturex

import (
	"fmt"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/unit"
)

// InputEvents detects drag gestures in the form of pointer.InputEvents events.
type InputEvents struct {
	Tag event.Tag
	io  pointer.InputOp
	pid pointer.ID
	ptr pointer.Cursor
	pressed,
	dragging bool
	start f32.Point
}

// Add the handler to the operation list to receive drag events.
func (d *InputEvents) Add(ops *op.Ops) {
	d.io = pointer.InputOp{
		Tag:   d.Tag,
		Types: pointer.Press | pointer.Drag | pointer.Move | pointer.Release,
	}
	d.io.Add(ops)
}

// Events returns the next drag events, if any.
func (d *InputEvents) Events(
	cfg unit.Metric, ops *op.Ops, q event.Queue, pressCallback func(pos f32.Point, buttons pointer.Buttons), releaseCallback func(pos f32.Point, buttons pointer.Buttons), primaryDragCallback, secondaryDragCallback func(diff f32.Point),
) {
	for _, e := range q.Events(d.Tag) {
		if pe, ok := e.(pointer.Event); ok {
			d.ptr = d.handlePointerEvent(pe, pressCallback, releaseCallback, primaryDragCallback, secondaryDragCallback)
		}
	}

	pointer.Cursor.Add(d.ptr, ops)
}

func (d *InputEvents) handlePointerEvent(
	e pointer.Event,
	pressCallback func(pos f32.Point, buttons pointer.Buttons),
	releaseCallback func(pos f32.Point, buttons pointer.Buttons),
	primaryDragCallback, secondaryDragCallback func(diff f32.Point),
) pointer.Cursor {
	ptr := pointer.CursorDefault

	switch e.Type {
	case pointer.Press:
		if e.Buttons == pointer.ButtonPrimary {
			if pressCallback != nil {
				pressCallback(e.Position, e.Buttons)
			}
			ptr = pointer.CursorDefault
		} else {
			ptr = pointer.CursorGrabbing
		}

		d.pressed = true
		d.start = e.Position
	case pointer.Move:
		d.start = e.Position
	case pointer.Drag:
		d.dragging = d.pressed
		if d.dragging {
			d.io.Grab = true
			ptr = pointer.CursorGrabbing
			diff := d.start.Sub(e.Position)
			if e.Buttons == pointer.ButtonPrimary {
				if primaryDragCallback != nil {
					primaryDragCallback(diff)
					ptr = pointer.CursorDefault
				}
			} else if e.Buttons == pointer.ButtonSecondary {
				if secondaryDragCallback != nil {
					secondaryDragCallback(diff)
					ptr = pointer.CursorGrabbing
				}
			}
		}
		d.start = e.Position
	case pointer.Release, pointer.Cancel:
		fmt.Printf("%+v\n", e)
		d.pressed = false
		d.io.Grab = false
		ptr = pointer.CursorDefault
		if releaseCallback != nil {
			releaseCallback(e.Position, e.Buttons)
		}
	}

	return ptr
}

// InputEventsging reports whether it is currently in use.
func (d *InputEvents) Dragging() bool { return d.dragging }

// Pressed returns whether a pointer is pressing.
func (d *InputEvents) Pressed() bool { return d.pressed }
