package peer

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/internal/enclave"
	"github.com/snivilised/jaywalk/src/agenor/internal/kernel"
	"github.com/snivilised/jaywalk/src/agenor/pref"
)

// IfActive returns a new plugin if peer info is active, otherwise nil.
func IfActive(o *pref.Options,
	_ enums.Subscription,
	mediator enclave.Mediator,
) enclave.Plugin {
	if o.View.Peer.IsActive {
		return &plugin{
			BasePlugin: kernel.BasePlugin{
				O:             o,
				Mediator:      mediator,
				ActivatedRole: enums.RolePeer,
			},
		}
	}

	return nil
}

// plugin manages peer info buffering during navigation.
type plugin struct {
	kernel.BasePlugin
	buffer peerBuffer
}

// Register registers the plugin with the kernel controller.
func (p *plugin) Register(kc enclave.KernelController) error {
	return p.BasePlugin.Register(kc)
}

// Next buffers the current servant and delivers the previously buffered
// one directly to the client via Poke, bypassing the guardian chain.
// Always returns false so the current servant is never forwarded through
// the chain.
func (p *plugin) Next(servant core.Servant,
	inspection enclave.Inspection,
) (bool, error) {
	return p.buffer.next(servant, inspection)
}

// Init sets up the buffer and registers life-cycle handlers for
// descend, ascend, and end events, then decorates the plugin into
// the guardian chain.
func (p *plugin) Init(pi *enclave.PluginInit) error {
	p.buffer.init(p.Mediator)

	pi.Controls.Ascend.On(func(node *core.Node) {
		p.buffer.onAscend(node.Extension.Depth)
	})

	pi.Controls.End.On(func(_ core.TraverseResult) {
		_ = p.buffer.end()
	})

	return p.Mediator.Decorate(p)
}
