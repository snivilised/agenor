package sampling

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/pref"
)

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

type samplingOptions struct {
	sampling *pref.SamplingOptions
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
