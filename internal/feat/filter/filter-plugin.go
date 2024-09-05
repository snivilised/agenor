package filter

// 📦 pkg: filter - defines filters

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, _ *pref.Using, mediator types.Mediator) types.Plugin {
	if o.Filter.IsFilteringActive() {
		return &plugin{
			BasePlugin: kernel.BasePlugin{
				O:             o,
				Mediator:      mediator,
				ActivatedRole: enums.RoleClientFilter,
			},
			sink:   o.Filter.Sink,
			scheme: newScheme(o),
		}
	}

	return nil
}

// plugin manages all filtering aspects of navigation
type plugin struct {
	kernel.BasePlugin
	sink   pref.FilteringSink
	crate  measure.Crate
	scheme scheme
}

func (p *plugin) Name() string {
	return "filtering"
}

func (p *plugin) Register(kc types.KernelController) error {
	if err := p.BasePlugin.Register(kc); err != nil {
		return err
	}

	return p.scheme.create()
}

func (p *plugin) Next(node *core.Node, inspection types.Inspection) (bool, error) {
	return p.scheme.next(node, inspection)
}

func (p *plugin) Init(pi *types.PluginInit) error {
	p.crate.Mums = p.Mediator.Supervisor().Many(
		enums.MetricNoFoldersFilteredOut,
		enums.MetricNoFilesFilteredOut,
		enums.MetricNoChildFilesFound,
		enums.MetricNoChildFilesFilteredOut,
	)

	p.scheme.init(pi, &p.crate)

	return p.Mediator.Decorate(p)
}
