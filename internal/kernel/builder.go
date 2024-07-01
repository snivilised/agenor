package kernel

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
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
		Build(o *pref.Options,
			resources *types.Resources,
		) (*Artefacts, error)
	}

	Builder func(o *pref.Options,
		resources *types.Resources,
	) (*Artefacts, error)
)

func (fn Builder) Build(o *pref.Options,
	resources *types.Resources,
) (*Artefacts, error) {
	return fn(o, resources)
}
