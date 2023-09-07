package context

import (
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
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
	appEvents []struct{}
}

func NewContext(ops *op.Ops, e system.FrameEvent) Context {
	return Context{
		layout.NewContext(ops, e),
		[]struct{}{},
	}
}

func (c Context) AppEvents() []struct{} {
	return c.appEvents
}

func (c Context) AddAppEvent(e struct{}) {
	c.appEvents = append(c.appEvents, e)
}
