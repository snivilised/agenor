package kernel

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/pref"
)

type (
	Artefacts struct {
		Kontroller enclave.KernelController
		Mediator   enclave.Mediator
		Resources  *enclave.Resources
		IfResult   core.Completion
	}

	Creation struct {
		Facade       pref.Facade
		Subscription enums.Subscription
		Harvest      enclave.OptionHarvest
		Resources    *enclave.Resources
	}

	NavigatorBuilder interface {
		Build(creation *Creation) *Artefacts
	}

	Builder func(creation *Creation) *Artefacts
)

func (fn Builder) Build(creation *Creation) *Artefacts {
	return fn(creation)
}
