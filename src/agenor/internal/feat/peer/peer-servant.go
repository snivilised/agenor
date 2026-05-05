package peer

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
)

// peerServant wraps the original servant, adding resolved PeerInfo.
// The original servant is embedded so all existing methods pass through
// unchanged.
type peerServant struct {
	core.Servant
	info *core.PeerInfo
}

func (s *peerServant) Peer() *core.PeerInfo {
	return s.info
}
