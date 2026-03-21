package kernel

import (
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/internal/enclave"
	"github.com/snivilised/jaywalk/src/agenor/pref"
)

// BasePlugin is a base struct for plugins.
type BasePlugin struct {
	// O is the options for the plugin.
	O *pref.Options
	// Mediator is the mediator for the plugin.
	Mediator enclave.Mediator
	// Kontroller is the kernel controller for the plugin.
	Kontroller enclave.KernelController
	// ActivatedRole is the role of the plugin.
	ActivatedRole enums.Role
}

// Role returns the role of the plugin.
func (p *BasePlugin) Role() enums.Role {
	return p.ActivatedRole
}

// Register registers the kernel controller with the plugin.
func (p *BasePlugin) Register(kc enclave.KernelController) error {
	p.Kontroller = kc

	return nil
}
