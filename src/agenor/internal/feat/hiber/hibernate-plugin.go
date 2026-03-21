package hiber

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/internal/enclave"
	"github.com/snivilised/jaywalk/src/agenor/internal/kernel"
	"github.com/snivilised/jaywalk/src/agenor/pref"
)

// IfActive returns a new plugin if the hibernate feature is active, otherwise nil.
func IfActive(o *pref.Options, _ enums.Subscription, mediator enclave.Mediator) enclave.Plugin {
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

// Next determines whether the servant should be filtered out or not, and
// returns true if it should be filtered out.
func (p *plugin) Next(servant core.Servant,
	inspection enclave.Inspection,
) (bool, error) {
	return p.profile.next(servant, servant.Node(), inspection)
}

// Init initializes the plugin, setting up the profile and decorating the plugin.
func (p *plugin) Init(pi *enclave.PluginInit) error {
	if err := p.profile.init(pi.Controls); err != nil {
		return err
	}

	return p.Mediator.Decorate(p)
}
