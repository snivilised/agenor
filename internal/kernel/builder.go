package kernel

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/enclave"
)

type (
	Artefacts struct {
		Kontroller enclave.KernelController
		Mediator   enclave.Mediator
		Resources  *enclave.Resources
		IfResult   core.Completion
	}

	NavigatorBuilder interface {
		Build(artefacts enclave.OptionHarvest,
			resources *enclave.Resources,
		) *Artefacts
	}

	Builder func(artefacts enclave.OptionHarvest,
		resources *enclave.Resources,
	) *Artefacts
)

func (fn Builder) Build(artefacts enclave.OptionHarvest,
	resources *enclave.Resources,
) *Artefacts {
	return fn(artefacts, resources)
}
