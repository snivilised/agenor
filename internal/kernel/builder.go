package kernel

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/pref"
)

type (
	Artefacts struct {
		Kontroller enclave.KernelController
		Mediator   enclave.Mediator
		Resources  *enclave.Resources
		IfResult   core.Completion
	}

	Inception struct {
		Facade       pref.Facade
		Subscription enums.Subscription
		Harvest      enclave.OptionHarvest
		Resources    *enclave.Resources
	}

	NavigatorBuilder interface {
		Build(inception *Inception) *Artefacts
	}

	Builder func(inception *Inception) *Artefacts
)

func (fn Builder) Build(inception *Inception) *Artefacts {
	return fn(inception)
}
