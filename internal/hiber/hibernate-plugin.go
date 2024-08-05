package hiber

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, mediator types.Mediator) types.Plugin {
	if o.Hibernate.Wake != nil {
		return &Plugin{
			BasePlugin: kernel.BasePlugin{
				Mediator:      mediator,
				ActivatedRole: enums.RoleClientHiberSleep, // TODO: or wake; to be resolved
			},
		}
	}

	return nil
}

type Plugin struct {
	kernel.BasePlugin
}

func (p *Plugin) Name() string {
	return "hibernation"
}

func (p *Plugin) Register(kc types.KernelController) error {
	p.Kontroller = kc

	return nil
}

func (p *Plugin) Next(node *core.Node, inspection core.Inspection) (bool, error) {
	_, _ = node, inspection

	return true, nil
}

func (p *Plugin) Init(_ *types.PluginInit) error {
	return p.Mediator.Decorate(p)
}
