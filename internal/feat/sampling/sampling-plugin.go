package sampling

import (
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/pref"
)

// IfActive returns a new plugin if sampling is active, otherwise nil
func IfActive(o *pref.Options, _ enums.Subscription, mediator enclave.Mediator) enclave.Plugin {
	if o.Sampling.IsSamplingActive() {
		return &plugin{
			BasePlugin: kernel.BasePlugin{
				O:             o,
				Mediator:      mediator,
				ActivatedRole: enums.RoleSampler,
			},
			ctrl: controller{
				o: &o.Sampling,
			},
		}
	}

	return nil
}

type plugin struct {
	kernel.BasePlugin
	ctrl controller
}

func (p *plugin) Init(_ *enclave.PluginInit) error {
	p.O.Hooks.ReadDirectory.Chain(
		p.ctrl.sample,
	)

	return p.Mediator.Decorate(&p.ctrl)
}
