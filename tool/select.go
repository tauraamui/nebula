package tool

import "gioui.org/layout"

type Tool interface {
	Update(gtx layout.Context)
}

type Select struct {
}

func (s *Select) Update(gtx layout.Context) {}
