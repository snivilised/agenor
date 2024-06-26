package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type navigatorBase struct {
	o     *pref.Options
	using *pref.Using
}

func (n *navigatorBase) Top(ctx context.Context,
	static *navigationStatic,
) (*types.NavigateResult, error) {
	_, _ = ctx, static

	return &types.NavigateResult{}, nil
}

func (n *navigatorBase) Traverse(ctx context.Context,
	static *navigationStatic,
	current *core.Node,
) (*core.Node, error) {
	_, _ = ctx, static
	_ = current

	return nil, nil
}
