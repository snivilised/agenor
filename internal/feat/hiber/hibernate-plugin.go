package hiber

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, _ *pref.Using, mediator types.Mediator) types.Plugin {
	if o.Hibernate.IsHibernateActive() {
		return &plugin{
			BasePlugin: kernel.BasePlugin{
				Mediator:      mediator,
				ActivatedRole: enums.RoleHibernate,
			},
			profile: &simple{
				common: common{
					ho: &o.Hibernate,
					fo: &o.Filter,
				},
			},
		}
	}

	return nil
}

type plugin struct {
	kernel.BasePlugin
	profile profile
}

func (p *plugin) Name() string {
	return "hibernation"
}

func (p *plugin) Next(node *core.Node, inspection types.Inspection) (bool, error) {
	return p.profile.next(node, inspection)
}

func (p *plugin) Init(pi *types.PluginInit) error {
	if err := p.profile.init(pi.Controls); err != nil {
		return err
	}

	return p.Mediator.Decorate(p)
}
