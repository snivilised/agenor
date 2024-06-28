package kernel

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
)

type navigatorFiles struct {
	navigatorBase
}

func (n *navigatorFiles) Top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	return n.navigatorBase.Top(ctx, ns)
}
