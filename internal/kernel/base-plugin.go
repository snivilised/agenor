package kernel

import (
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/pref"
)

type BasePlugin struct {
	O             *pref.Options
	Mediator      enclave.Mediator
	Kontroller    enclave.KernelController
	ActivatedRole enums.Role
}

func (p *BasePlugin) Role() enums.Role {
	return p.ActivatedRole
}

func (p *BasePlugin) Register(kc enclave.KernelController) error {
	p.Kontroller = kc

	return nil
}
