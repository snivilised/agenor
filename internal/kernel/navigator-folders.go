package kernel

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
)

type navigatorFolders struct {
	navigatorBase
}

func (n *navigatorFolders) Top(ctx context.Context,
	static *navigationStatic,
) (*types.NavigateResult, error) {
	return n.navigatorBase.Top(ctx, static)
}
