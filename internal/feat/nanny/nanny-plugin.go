package nanny

// 📦 pkg: nanny - handles a node's children for directories with children subscription

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options,
	facade pref.Facade, mediator enclave.Mediator,
) enclave.Plugin {
	if facade.Sub() == enums.SubscribeDirectoriesWithFiles &&
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
	crate core.Crate
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
