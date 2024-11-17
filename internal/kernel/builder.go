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

func (i *Inception) NavigationTree() string {
	if using, ok := i.Facade.(*pref.Using); ok {
		return using.Tree
	}

	return i.Harvest.Loaded().State.Tree
}
