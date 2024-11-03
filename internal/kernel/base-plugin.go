package kernel

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/pref"
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
