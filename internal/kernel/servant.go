package kernel

import (
	"github.com/snivilised/traverse/core"
)

type (
	servant struct {
		node *core.Node
	}
)

func (s servant) Node() *core.Node {
	return s.node
}
