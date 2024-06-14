package kernel

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
)

type navigatorUniversal struct {
	navigatorBase
}

func (n *navigatorUniversal) Top(ctx context.Context,
	root string,
) (*types.NavigateResult, error) {
	return n.navigatorBase.Top(ctx, root)
}
