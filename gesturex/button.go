package gesturex

import (
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/unit"
)

type ButtonEvents struct {
	Tag            event.Tag
	io             pointer.InputOp
	pid            pointer.ID
	ptr            pointer.Cursor
	pressedButtons pointer.Buttons
	pressed,
	leftWhilstPressed bool
}

// Add the handler to the operation list to receive drag events.
func (d *ButtonEvents) Add(ops *op.Ops) {
	d.io = pointer.InputOp{
		Tag:   d.Tag,
		Types: pointer.Enter | pointer.Leave | pointer.Press | pointer.Release,
	}
	d.io.Add(ops)
}

// Events returns the next drag events, if any.
func (d *ButtonEvents) Events(
	cfg unit.Metric, ops *op.Ops, q event.Queue, pressCallback, leaveWhilePressedCallback, releaseCallback func(),
) {
	for _, e := range q.Events(d.Tag) {
		if pe, ok := e.(pointer.Event); ok {
			d.ptr = d.handlePointerEvent(pe, pressCallback, leaveWhilePressedCallback, releaseCallback)
		}
	}

	pointer.Cursor.Add(d.ptr, ops)
}

func (d *ButtonEvents) handlePointerEvent(
	e pointer.Event,
	pressCallback func(),
	leaveWhilePressedCallback func(),
	releaseCallback func(),
) pointer.Cursor {
	ptr := pointer.CursorDefault

	switch e.Type {
	case pointer.Enter:
		ptr = pointer.CursorPointer
		d.leftWhilstPressed = false
		break
	case pointer.Leave:
		if d.pressed && leaveWhilePressedCallback != nil {
			leaveWhilePressedCallback()
			d.pressed = false
			d.leftWhilstPressed = true
		}
		break
	case pointer.Press:
		if e.Buttons == pointer.ButtonPrimary {
			d.pressed = true
			d.io.Grab = true
			if pressCallback != nil {
				pressCallback()
			}
			ptr = pointer.CursorPointer
			break
		}
	case pointer.Release, pointer.Cancel:
		if d.leftWhilstPressed {
			break
		}
		d.pressed = false
		d.io.Grab = false
		ptr = pointer.CursorDefault
		if releaseCallback != nil && e.Type != pointer.Cancel {
			releaseCallback()
		}
	}

	return ptr
}
