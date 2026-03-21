package nanny

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/internal/enclave"
	"github.com/snivilised/jaywalk/src/agenor/internal/kernel"
	"github.com/snivilised/jaywalk/src/agenor/pref"
)

// IfActive returns a new nanny plugin if the subscription is directories with files
// and no filtering is active, otherwise it returns nil.
// The nanny plugin is responsible for handling the children of a node when the
// directories with files subscription is active and no filtering is active.
// It will sort the children of a node into files and directories, and then
// assign the files to the node's children. The nanny plugin will also track
// the number of child files found for each node. A new plugin if filtering
// is active, otherwise nil.
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

// Next determines whether the servant should be filtered out or not,
// and returns true if it should be filtered out.
func (p *plugin) Next(servant core.Servant,
	inspection enclave.Inspection,
) (bool, error) {
	node := servant.Node()
	files := inspection.Sort(enums.EntryTypeFile)
	node.Children = files
	p.crate.Metrics[enums.MetricNoChildFilesFound].Times(uint(len(files)))

	return true, nil
}

// Init initializes the plugin, setting up the metrics and decorating the plugin.
func (p *plugin) Init(_ *enclave.PluginInit) error {
	p.crate.Metrics = p.Mediator.Supervisor().Many(
		enums.MetricNoChildFilesFound,
	)

	return p.Mediator.Decorate(p)
}
