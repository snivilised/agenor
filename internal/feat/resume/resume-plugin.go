package resume

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type (
	From struct {
		Active   *core.ActiveState
		Mediator types.Mediator
		Strategy enums.ResumeStrategy
		IfResult core.ResultCompletion
	}

	Plugin struct {
		kernel.BasePlugin
		IfResult   core.ResultCompletion
		Active     *core.ActiveState
		kontroller types.KernelController
	}
)

func New(from *From,
) *Plugin {
	return &Plugin{
		Active: from.Active,
		BasePlugin: kernel.BasePlugin{
			Mediator: from.Mediator,
			ActivatedRole: lo.Ternary(from.Strategy == enums.ResumeStrategyFastward,
				enums.RoleFastward, enums.RoleUndefined,
			),
		},
		IfResult: from.IfResult,
	}
}

func (p *Plugin) Init(pi *types.PluginInit) error {
	p.kontroller = pi.Kontroller

	return nil
}

func (p *Plugin) IsComplete() bool {
	return p.IfResult.IsComplete()
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

func WithArtefacts(was *pref.Was, harvest types.OptionHarvest,
	resources *types.Resources,
) *kernel.Artefacts {
	sealer := lo.Ternary(was.Strategy == enums.ResumeStrategyFastward,
		types.GuardianSealer(&fastwardGuardianSealer{}),
		types.GuardianSealer(&kernel.Benign{}),
	)

	// TODO: create a general type that carries all this info; pass
	// this into WithArtefacts
	//
	controller := kernel.New(&was.Using, harvest.Options(), resources, sealer)
	strategy := newStrategy(was, harvest, controller, sealer, resources)

	return &kernel.Artefacts{
		Kontroller: &Controller{
			kc:       controller,
			was:      was,
			load:     harvest.Loaded(),
			strategy: strategy,
		},
		Mediator:  controller.Mediator(),
		Resources: resources,
		IfResult:  strategy.ifResult,
	}
}
