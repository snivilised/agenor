package refine

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, mediator types.Mediator) types.Plugin {
	if o.Core.Filter.Node != nil {
		return &Plugin{
			BasePlugin: kernel.BasePlugin{
				Mediator: mediator,
			},
		}
	}

	return nil
}

type Plugin struct {
	kernel.BasePlugin
}

func (p *Plugin) Name() string {
	return "filtering"
}

func (p *Plugin) Register() error {
	return nil
}

func (p *Plugin) Next(node *core.Node) (bool, error) {
	_ = node
	// if filtered in send filtered in message
	// if filtered out send filtered out message
	// these filtered in/out messages could be handled
	// by the metrics plugin, so that it can record
	// counts

	return true, nil
}

func (p *Plugin) Role() enums.Role {
	return enums.RoleClientFilter
}

func (p *Plugin) Init() error {
	p.Mediator.Supervisor().Many(
		enums.MetricNoFoldersFilteredOut,
		enums.MetricNoFilesFilteredOut,
	)

	return p.Mediator.Decorate(p)
}
