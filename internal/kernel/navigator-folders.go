package kernel

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
)

type navigatorFolders struct {
	navigatorBase
}

func (n *navigatorFolders) Top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	return n.navigatorBase.Top(ctx, ns)
}
