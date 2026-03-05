package filter

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/pref"
)

// IfActive returns a new plugin if filtering is active, otherwise nil
func IfActive(o *pref.Options, _ enums.Subscription, mediator enclave.Mediator) enclave.Plugin {
	if o.Filter.IsFilteringActive() {
		return &plugin{
			BasePlugin: kernel.BasePlugin{
				O:             o,
				Mediator:      mediator,
				ActivatedRole: enums.RoleClientFilter,
			},
			sink:   o.Filter.Sink,
			scheme: newScheme(o),
		}
	}

	return nil
}

// plugin manages all filtering aspects of navigation
type plugin struct {
	kernel.BasePlugin
	sink   pref.FilteringSink
	crate  enclave.Crate
	scheme scheme
}

// Register registers the plugin with the kernel controller.
func (p *plugin) Register(kc enclave.KernelController) error {
	if err := p.BasePlugin.Register(kc); err != nil {
		return err
	}

	return p.scheme.create()
}

// Next determines whether the servant should be filtered out or not,
// and returns true if it should be filtered out.
func (p *plugin) Next(servant core.Servant,
	inspection enclave.Inspection,
) (bool, error) {
	return p.scheme.next(servant, inspection)
}

// Init initializes the plugin, setting up necessary metrics and
// preparing the scheme for operation. It also decorates the plugin
// with any additional functionality provided by the mediator.
func (p *plugin) Init(pi *enclave.PluginInit) error {
	p.crate.Metrics = p.Mediator.Supervisor().Many(
		enums.MetricNoDirectoriesFilteredOut,
		enums.MetricNoFilesFilteredOut,
		enums.MetricNoChildFilesFound,
		enums.MetricNoChildFilesFilteredOut,
	)

	p.scheme.init(pi, &p.crate)

	return p.Mediator.Decorate(p)
}
