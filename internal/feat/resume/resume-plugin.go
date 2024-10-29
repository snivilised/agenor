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
}

func (p *Plugin) Next(servant core.Servant,
	inspection types.Inspection,
) (bool, error) {
	_, _ = servant, inspection
	// apply the wake filter

	return true, nil
}

func (p *Plugin) Role() enums.Role {
	return enums.RoleHibernate
}

func (p *Plugin) Init(_ *types.PluginInit) error {
	return p.Mediator.Decorate(p)
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
