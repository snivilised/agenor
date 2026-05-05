package kernel

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
)

type (
	servant struct {
		node *core.Node
		peer *core.PeerInfo
	}
)

func (s servant) Node() *core.Node {
	return s.node
}

func (s servant) Peer() *core.PeerInfo {
	return s.peer
}
