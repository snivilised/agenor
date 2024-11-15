package resume

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/internal/persist"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
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

func New(from *From) *Plugin {
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

	if err != nil {
		return &opts.LoadInfo{}, nil, err
	}

	return opts.Bind(result.O, result.Active, settings...)
}

func Artefacts(creation *kernel.Creation) *kernel.Artefacts {
	// the error from the following facade typecast is ignored, because
	// this is already checked by the creation of the scaffolding.
	//
	relic, _ := creation.Facade.(*pref.Relic)

	sealer := lo.Ternary(relic.Strategy == enums.ResumeStrategyFastward,
		enclave.GuardianSealer(&FastwardGuardianSealer{}),
		enclave.GuardianSealer(&kernel.Benign{}),
	)

	controller := kernel.New(creation, sealer)
	strategy := newStrategy(creation, sealer, controller)

	return &kernel.Artefacts{
		Kontroller: &Controller{
			kc:       controller,
			relic:    relic,
			load:     creation.Harvest.Loaded(),
			strategy: strategy,
		},
		Mediator:  controller.Mediator(),
		Resources: creation.Resources,
		IfResult:  strategy.ifResult,
	}
}
