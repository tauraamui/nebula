package context

import (
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/tauraamui/nebula/f32x"
)

/*
type Context interface {
	Dp(v unit.Dp) int
	Sp(v unit.Sp) int
	AppEvents() []struct{}
	AddAppEvent(e struct{})
	Events(k event.Tag) []event.Event
	Disabled() layout.Context
}
*/

type Context struct {
	layout.Context
	appEvents []any
}

func NewContext(ops *op.Ops, e system.FrameEvent) *Context {
	return &Context{
		layout.NewContext(ops, e),
		[]any{},
	}
}

func (c *Context) Events() []any {
	return c.appEvents
}

func (c *Context) PushEvent(e any) {
	c.appEvents = append(c.appEvents, e)
}

type CreateMatrix struct {
	Pos        f32.Point
	Rows, Cols int
	Bounds     f32x.Rectangle
}
