package resume

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

type (
	From struct {
		Active   *core.ActiveState
		Mediator enclave.Mediator
		Strategy enums.ResumeStrategy
		IfResult core.ResultCompletion
	}

	Plugin struct {
		kernel.BasePlugin
		IfResult   core.ResultCompletion
		Active     *core.ActiveState
		kontroller enclave.KernelController
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

func (p *Plugin) Init(pi *enclave.PluginInit) error {
	p.kontroller = pi.Kontroller

	return nil
}

func (p *Plugin) IsComplete() bool {
	return p.IfResult.IsComplete()
}

func Load(restoration *enclave.RestoreState,
	settings ...pref.Option,
) (*opts.LoadInfo, *opts.Binder, error) {
	result, err := persist.Unmarshal(&persist.UnmarshalRequest{
		Restore: restoration,
	})

	_ = err // TODO: don't forget to handle this

	return opts.Bind(result.O, result.Active, settings...)
}

func Artefacts(relic *pref.Relic, harvest enclave.OptionHarvest,
	resources *enclave.Resources,
) *kernel.Artefacts {
	sealer := lo.Ternary(relic.Strategy == enums.ResumeStrategyFastward,
		enclave.GuardianSealer(&fastwardGuardianSealer{}),
		enclave.GuardianSealer(&kernel.Benign{}),
	)

	// TODO: create a general type that carries all this info; pass
	// this into WithArtefacts
	//
	controller := kernel.New(relic, harvest.Options(), resources, sealer)
	strategy := newStrategy(relic, harvest, controller, sealer, resources)

	return &kernel.Artefacts{
		Kontroller: &Controller{
			kc:       controller,
			relic:    relic,
			load:     harvest.Loaded(),
			strategy: strategy,
		},
		Mediator:  controller.Mediator(),
		Resources: resources,
		IfResult:  strategy.ifResult,
	}
}
