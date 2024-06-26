package kernel

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
)

type navigatorUniversal struct {
	navigatorBase
}

func (n *navigatorUniversal) Top(ctx context.Context,
	static *navigationStatic,
) (*types.NavigateResult, error) {
	return top(ctx, static)
}
