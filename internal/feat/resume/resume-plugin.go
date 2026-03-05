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
	// From is a collection of information that defines where the
	// resume will continue from.
	From struct {
		// ActiveState represents state that needs to be persisted alongside
		// the options in order for resume to work.
		Active *core.ActiveState

		// Mediator controls interactions between different entities of
		// of the navigator
		Mediator enclave.Mediator

		// Strategy denotes which resume strategy to use
		Strategy enums.ResumeStrategy

		// IfResult is a ResultCompletion used to determine if the result really
		// represents final navigation completion. This is pertinent to spawn
		// resume where a completion event, may or may not mark the end of total
		// navigation.
		IfResult core.ResultCompletion
	}

	// Plugin the resume plugin
	Plugin struct {
		kernel.BasePlugin
		IfResult   core.ResultCompletion
		Active     *core.ActiveState
		kontroller enclave.KernelController
	}
)

// New creates a new plugin.
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

// Init initializes the plugin.
func (p *Plugin) Init(pi *enclave.PluginInit) error {
	p.kontroller = pi.Kontroller

	return nil
}

// IsComplete returns true if the plugin is complete.
func (p *Plugin) IsComplete() bool {
	return p.IfResult.IsComplete()
}

// Load loads the plugin.
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

// Artefacts creates the artefacts for the plugin.
func Artefacts(inception *kernel.Inception) *kernel.Artefacts {
	// the error from the following facade typecast is ignored, because
	// this is already checked by the inception of the scaffolding.
	//
	relic, _ := inception.Facade.(*pref.Relic)

	sealer := lo.Ternary(relic.Strategy == enums.ResumeStrategyFastward,
		enclave.GuardianSealer(&FastwardGuardianSealer{}),
		enclave.GuardianSealer(&kernel.Benign{}),
	)

	mediator, err := kernel.NewMediator(inception, sealer)
	strategy := newStrategy(inception, sealer, mediator)

	return &kernel.Artefacts{
		Kontroller: &Controller{
			med:      mediator,
			relic:    relic,
			load:     inception.Harvest.Loaded(),
			strategy: strategy,
		},
		Mediator:  mediator,
		Resources: inception.Resources,
		IfResult:  strategy.ifResult,
		Error:     err,
	}
}
