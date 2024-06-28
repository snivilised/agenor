package sampling

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, mediator types.Mediator) types.Plugin {
	if (o.Core.Sampling.NoOf.Files > 0) || (o.Core.Sampling.NoOf.Folders > 0) {
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
	return "sampling"
}

func (p *Plugin) Register() error {
	return nil
}

func (p *Plugin) Next(node *core.Node) (bool, error) {
	_ = node

	// apply the filter to the node
	return true, nil
}

func (p *Plugin) Role() enums.Role {
	return enums.RoleSampler
}

func (p *Plugin) Init() error {
	return p.Mediator.Decorate(p)
}
