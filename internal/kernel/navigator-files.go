package kernel

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
)

type navigatorFiles struct {
	navigatorBase
}

func (n *navigatorFiles) Top(ctx context.Context,
	static *navigationStatic,
) (*types.NavigateResult, error) {
	return n.navigatorBase.Top(ctx, static)
}
