package nanny

// 📦 pkg: nanny - handles a node's children for directories with children subscription

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/pref"
)

func IfActive(o *pref.Options,
	sub enums.Subscription, mediator enclave.Mediator,
) enclave.Plugin {
	if sub == enums.SubscribeDirectoriesWithFiles &&
		!o.Filter.IsFilteringActive() {
		return &plugin{
			BasePlugin: kernel.BasePlugin{
				O:             o,
				Mediator:      mediator,
				ActivatedRole: enums.RoleNanny,
			},
		}
	}

	return nil
}

type plugin struct {
	kernel.BasePlugin
	crate enclave.Crate
}

func (p *plugin) Next(servant core.Servant,
	inspection enclave.Inspection,
) (bool, error) {
	node := servant.Node()
	files := inspection.Sort(enums.EntryTypeFile)
	node.Children = files
	p.crate.Metrics[enums.MetricNoChildFilesFound].Times(uint(len(files)))

	return true, nil
}

func (p *plugin) Init(_ *enclave.PluginInit) error {
	p.crate.Metrics = p.Mediator.Supervisor().Many(
		enums.MetricNoChildFilesFound,
	)

	return p.Mediator.Decorate(p)
}
