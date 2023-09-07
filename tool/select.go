package tool

import (
	"github.com/tauraamui/nebula/context"
)

type Tool interface {
	Update(gtx context.Context)
}

type Select struct {
}

func (s *Select) Update(gtx context.Context) {}
