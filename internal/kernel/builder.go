package kernel

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
)

type (
	Artefacts struct {
		Kontroller types.KernelController
		Mediator   types.Mediator
		Resources  *types.Resources
		IfResult   core.Completion
	}

	NavigatorBuilder interface {
		Build(artefacts types.OptionHarvest,
			resources *types.Resources,
		) *Artefacts
	}

	Builder func(artefacts types.OptionHarvest,
		resources *types.Resources,
	) *Artefacts
)

func (fn Builder) Build(artefacts types.OptionHarvest,
	resources *types.Resources,
) *Artefacts {
	return fn(artefacts, resources)
}
