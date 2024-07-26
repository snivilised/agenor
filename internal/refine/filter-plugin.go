package refine

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, mediator types.Mediator) types.Plugin {
	if o.Core.Filter.IsFilteringActive() || o.Filtering.IsCustomFilteringActive() {
		return &Plugin{
			BasePlugin: kernel.BasePlugin{
				O:             o,
				Mediator:      mediator,
				ActivatedRole: enums.RoleClientFilter,
			},
			sink:   o.Filtering.FilterSink,
			scheme: newScheme(o),
		}
	}

	return nil
}

// Plugin manages all filtering aspects of navigation
type Plugin struct {
	kernel.BasePlugin
	sink   pref.FilteringSink
	owner  measure.Owned
	scheme scheme
}

func (p *Plugin) Name() string {
	return "filtering"
}

func (p *Plugin) Register(kc types.KernelController) error {
	p.Kontroller = kc
	return p.scheme.create()
}

func (p *Plugin) Next(node *core.Node, inspection core.Inspection) (bool, error) {
	return p.scheme.next(node, inspection)
}

func (p *Plugin) Init(pi *types.PluginInit) error {
	p.owner.Mums = p.Mediator.Supervisor().Many(
		enums.MetricNoFoldersFilteredOut,
		enums.MetricNoFilesFilteredOut,
		enums.MetricNoChildFilesFound,
		enums.MetricNoChildFilesFilteredOut,
	)

	p.scheme.init(pi, &p.owner)

	return p.Mediator.Decorate(p)
}
