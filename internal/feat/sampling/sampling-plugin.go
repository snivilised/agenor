package sampling

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, _ *pref.Using, mediator types.Mediator) types.Plugin {
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

func (p *plugin) Name() string {
	return "sampling"
}

func (p *plugin) Init(_ *types.PluginInit) error {
	p.O.Hooks.ReadDirectory.Chain(
		p.ctrl.sample,
	)

	return p.Mediator.Decorate(&p.ctrl)
}
