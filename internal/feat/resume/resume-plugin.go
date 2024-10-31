package resume

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type Plugin struct {
	kernel.BasePlugin
	IfResult core.ResultCompletion
	Active   *core.ActiveState
	guardian types.Guardian
}

func (p *Plugin) Init(pi *types.PluginInit) error {
	p.guardian = pi.Kontroller.Mediator()

	return nil
}

func (p *Plugin) IsComplete() bool {
	return p.IfResult.IsComplete()
}

func GetSealer(was *pref.Was) types.GuardianSealer {
	if was.Strategy == enums.ResumeStrategyFastward {
		return &fastwardGuardianSealer{}
	}

	return &kernel.Benign{}
}

func Load(restoration *types.RestoreState,
	settings ...pref.Option,
) (*opts.LoadInfo, *opts.Binder, error) {
	result, err := persist.Unmarshal(&persist.UnmarshalRequest{
		Restore: restoration,
	})

	_ = err // TODO: don't forget to handle this

	return opts.Bind(result.O, result.Active, settings...)
}

// this is not named correctly
func NewController(was *pref.Was, harvest types.OptionHarvest,
	artefacts *kernel.Artefacts,
) *kernel.Artefacts {
	// The Controller on the incoming artefacts is the core navigator. It is
	// decorated here for resume. The strategy only needs access to the core navigator.
	// The resume navigator delegates to the strategy.
	//
	var (
		strategy resumeStrategy
		err      error
	)

	if strategy, err = newStrategy(was, harvest, artefacts.Kontroller); err != nil {
		return artefacts
	}

	return &kernel.Artefacts{
		Kontroller: &Controller{
			kc:         artefacts.Kontroller,
			was:        was,
			load:       harvest.Loaded(),
			strategy:   strategy,
			facilities: artefacts.Facilities,
		},
		Mediator:  artefacts.Mediator,
		Resources: artefacts.Resources,
		IfResult:  strategy.ifResult,
	}
}
