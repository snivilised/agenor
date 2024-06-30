package kernel

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
)

type navigatorFiles struct {
	navigator
}

func (n *navigatorFiles) Top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	return n.navigator.Top(ctx, ns)
}
