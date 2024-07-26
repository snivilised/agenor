package sampling

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options, mediator types.Mediator) types.Plugin {
	if o.Core.Sampling.IsSamplingActive() {
		return &Plugin{
			BasePlugin: kernel.BasePlugin{
				O:             o,
				Mediator:      mediator,
				ActivatedRole: enums.RoleSampler,
			},
			ctrl: controller{
				o: &samplingOptions{
					sampling: &o.Core.Sampling,
					sampler:  &o.Sampler,
				},
			},
		}
	}

	return nil
}

type samplingOptions struct {
	sampling *pref.SamplingOptions
	sampler  *pref.SamplerOptions
}

type Plugin struct {
	kernel.BasePlugin
	ctrl controller
}

func (p *Plugin) Name() string {
	return "sampling"
}

func (p *Plugin) Register(kc types.KernelController) error {
	p.Kontroller = kc

	return nil
}

func (p *Plugin) Init(_ *types.PluginInit) error {
	p.O.Hooks.ReadDirectory.Chain(
		p.ctrl.sample,
	)

	return p.Mediator.Decorate(&p.ctrl)
}
