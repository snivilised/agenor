package kernel

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
)

type (
	Artefacts struct {
		Kontroller types.KernelController
		Mediator   types.Mediator
		Facilities types.Facilities
		Resources  *types.Resources
		IfResult   core.Completion
	}

	NavigatorBuilder interface {
		Build(artefacts types.OptionHarvest,
			resources *types.Resources,
		) (*Artefacts, error)
	}

	Builder func(artefacts types.OptionHarvest,
		resources *types.Resources,
	) (*Artefacts, error)
)

func (fn Builder) Build(artefacts types.OptionHarvest,
	resources *types.Resources,
) (*Artefacts, error) {
	return fn(artefacts, resources)
}
