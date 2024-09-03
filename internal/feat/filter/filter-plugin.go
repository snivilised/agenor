package filter

// 📦 pkg: filter - defines filters

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/override"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, mediator types.Mediator) types.Plugin {
	if o.Filter.IsFilteringActive() {
		return &Plugin{
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

// Plugin manages all filtering aspects of navigation
type Plugin struct {
	kernel.BasePlugin
	sink   pref.FilteringSink
	crate  measure.Crate
	scheme scheme
}

func (p *Plugin) Name() string {
	return "filtering"
}

func (p *Plugin) Register(kc types.KernelController) error {
	p.Kontroller = kc
	return p.scheme.create()
}

func (p *Plugin) Next(node *core.Node, inspection override.Inspection) (bool, error) {
	return p.scheme.next(node, inspection)
}

func (p *Plugin) Init(pi *types.PluginInit) error {
	p.crate.Mums = p.Mediator.Supervisor().Many(
		enums.MetricNoFoldersFilteredOut,
		enums.MetricNoFilesFilteredOut,
		enums.MetricNoChildFilesFound,
		enums.MetricNoChildFilesFilteredOut,
	)

	p.scheme.init(pi, &p.crate)

	return p.Mediator.Decorate(p)
}
