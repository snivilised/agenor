package refine

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/lo"
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
			sink: o.Filtering.FilterSink,
		}
	}

	return nil
}

// Plugin manages all filtering aspects of navigation
type Plugin struct {
	kernel.BasePlugin
	filters NavigationFilters
	sink    pref.FilteringSink
	owner   measure.Owned
}

func (p *Plugin) Name() string {
	return "filtering"
}

func (p *Plugin) Register(kc types.KernelController) error {
	p.Kontroller = kc

	if p.O.Core.Filter.IsNodeFilteringActive() {
		filter, err := newNodeFilter(p.O.Core.Filter.Node, &p.O.Filtering)
		if err != nil {
			return err
		}

		p.filters.Node = filter
	}

	if p.O.Filtering.IsCustomFilteringActive() {
		p.O.Filtering.Custom.Validate()
	}

	p.filters.Node = pref.ResolveFilter(p.filters.Node, p.O.Filtering)

	if p.O.Core.Filter.IsChildFilteringActive() {
		filter, err := newChildFilter(p.O.Core.Filter.Child)
		if err != nil {
			return err
		}
		p.filters.Children = filter
	}

	if p.sink != nil {
		p.sink(pref.FilterReply{
			Node:  p.filters.Node,
			Child: p.filters.Children,
		})
	}

	return nil
}

func (p *Plugin) Next(node *core.Node) (bool, error) {
	if p.filters.Node == nil {
		return true, nil
	}

	matched := p.filters.Node.IsMatch(node)

	if !matched {
		filteredOutMetric := lo.Ternary(node.IsFolder(),
			enums.MetricNoFoldersFilteredOut,
			enums.MetricNoFilesFilteredOut,
		)
		p.owner.Mums[filteredOutMetric].Tick()
	}

	return matched, nil
}

func (p *Plugin) Init(pi *types.PluginInit) error {
	// [KEEP-FILTER-IN-SYNC] keep this in sync with the default
	// behaviour in builders.override.Actions
	p.owner.Mums = p.Mediator.Supervisor().Many(
		enums.MetricNoFoldersFilteredOut,
		enums.MetricNoFilesFilteredOut,
		enums.MetricNoChildFilesFilteredOut,
	)

	pi.Actions.HandleChildren.Intercept(
		func(inspection core.Inspection, mums measure.MutableMetrics) {
			files := inspection.Sort(enums.EntryTypeFile)
			matching := lo.TernaryF(p.filters.Children != nil,
				func() []fs.DirEntry {
					return p.filters.Children.Matching(files)
				},
				func() []fs.DirEntry {
					return files
				},
			)

			inspection.AssignChildren(matching)
			mums[enums.MetricNoChildFilesFound].Times(uint(len(files)))

			filteredOut := len(files) - len(matching)
			p.owner.Mums[enums.MetricNoChildFilesFilteredOut].Times(uint(filteredOut))
		},
	)

	return p.Mediator.Decorate(p)
}
