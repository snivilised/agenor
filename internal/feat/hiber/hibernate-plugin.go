package hiber

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, _ *pref.Using, mediator enclave.Mediator) enclave.Plugin {
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

func (p *plugin) Next(servant core.Servant,
	inspection enclave.Inspection,
) (bool, error) {
	return p.profile.next(servant, servant.Node(), inspection)
}

func (p *plugin) Init(pi *enclave.PluginInit) error {
	if err := p.profile.init(pi.Controls); err != nil {
		return err
	}

	return p.Mediator.Decorate(p)
}
