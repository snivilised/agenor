package kernel

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
)

type navigatorFolders struct {
	navigatorBase
}

func (n *navigatorFolders) Top(ctx context.Context,
	root string,
) (*types.NavigateResult, error) {
	return n.navigatorBase.Top(ctx, root)
}
