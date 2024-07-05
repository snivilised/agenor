package kernel

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
)

type BasePlugin struct {
	Mediator      types.Mediator
	Kontroller    types.KernelController
	ActivatedRole enums.Role
}

func (p *BasePlugin) Role() enums.Role {
	return p.ActivatedRole
}
