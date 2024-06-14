package kernel

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type navigatorBase struct {
	o     *pref.Options
	using *pref.Using
}

func (n *navigatorBase) Top(ctx context.Context, root string) (*types.NavigateResult, error) {
	_, _ = ctx, root

	return &types.NavigateResult{}, nil
}
