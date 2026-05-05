package peer

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/internal/enclave"
)

type peerBuffer struct {
	mediator  enclave.Mediator
	pending   *peerServant
	liveStack []bool
	first     bool
}

func (b *peerBuffer) init(mediator enclave.Mediator) {
	b.mediator = mediator
	b.first = true
}

// onAscend marks the buffered node as the last peer in its directory.
// Delivery of the node is not triggered here - it occurs on the next
// call to next, or via end for the final node.
func (b *peerBuffer) onAscend(depth int) {
	if depth+1 < len(b.liveStack) {
		b.liveStack = b.liveStack[:depth+1]
	}

	if b.pending != nil {
		b.pending.info.IsLast = true
	}
}

// next buffers the current servant. If a previous servant is pending,
// it is delivered directly to the client via Poke, bypassing the
// guardian chain. Always returns false so the current servant is never
// forwarded through the chain.
func (b *peerBuffer) next(servant core.Servant, _ enclave.Inspection) (bool, error) {
	if !b.first {
		if err := b.poke(); err != nil {
			return false, err
		}
	}

	b.first = false

	depth := servant.Node().Extension.Depth
	for len(b.liveStack) <= depth {
		b.liveStack = append(b.liveStack, false)
	}

	b.pending = &peerServant{
		Servant: servant,
		info: &core.PeerInfo{
			IsLast:      false,
			IndentStack: append([]bool{}, b.liveStack...),
		},
	}

	return false, nil
}

// end delivers the final buffered node directly to the client.
// The preceding Ascend will have already set IsLast=true on the
// buffered node's stack position.
func (b *peerBuffer) end() error {
	if b.pending == nil {
		return nil
	}

	return b.poke()
}

func (b *peerBuffer) poke() error {
	enriched := b.pending
	b.pending = nil

	return b.mediator.Poke(enriched)
}
