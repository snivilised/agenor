package kernel

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
)

type navigatorFiles struct {
	navigatorBase
}

func (n *navigatorFiles) Top(ctx context.Context,
	root string,
) (*types.NavigateResult, error) {
	return n.navigatorBase.Top(ctx, root)
}
