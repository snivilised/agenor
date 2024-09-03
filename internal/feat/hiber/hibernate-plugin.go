package hiber

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/override"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, mediator types.Mediator) types.Plugin {
	if o.Hibernate.IsHibernateActive() {
		return &Plugin{
			BasePlugin: kernel.BasePlugin{
				Mediator:      mediator,
				ActivatedRole: enums.RoleHibernate,
			},
			profile: &simple{
				common: common{
					ho: &o.Hibernate,
				},
			},
		}
	}

	return nil
}

type Plugin struct {
	kernel.BasePlugin
	profile profile
}

func (p *Plugin) Name() string {
	return "hibernation"
}

func (p *Plugin) Register(kc types.KernelController) error {
	p.Kontroller = kc

	return nil
}

func (p *Plugin) Next(node *core.Node, inspection override.Inspection) (bool, error) {
	return p.profile.next(node, inspection)
}

func (p *Plugin) Init(pi *types.PluginInit) error {
	if err := p.profile.init(pi.Controls); err != nil {
		return err
	}

	return p.Mediator.Decorate(p)
}
