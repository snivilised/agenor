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
		// TODO: setup iterators (extendio: sampling-adapters):
		// - slice -> forward/reverse
		// - pre-defined filter iterator -> forward/reverse
		// - custom iterator
		//
		return &Plugin{
			BasePlugin: kernel.BasePlugin{
				Mediator:      mediator,
				ActivatedRole: enums.RoleSampler,
			},
			ctrl: controller{
				o: &samplingOptions{
					sampling: &o.Core.Sampling,
					sampler:  &o.Sampler,
				},
				on: handlers{
					descend: func(_ *core.Node) {},
					ascend:  func(_ *core.Node) {},
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

func (p *Plugin) Init(pi *types.PluginInit) error {
	pi.O.Events.Descend.On(p.ctrl.on.descend)
	pi.O.Events.Ascend.On(p.ctrl.on.ascend)

	return p.Mediator.Decorate(&p.ctrl)
}
