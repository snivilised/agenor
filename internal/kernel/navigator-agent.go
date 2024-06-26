package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
)

// navigatorAgent does work on behalf of the navigator. It is distinct
// from navigatorBase and should only be used when the limited polymorphism
// on base is inadequate.
type navigatorAgent struct {
}

func newAgent() *navigatorAgent {
	return &navigatorAgent{}
}

func top(ctx context.Context,
	static *navigationStatic,
) (*types.NavigateResult, error) {
	info, err := static.mediator.o.Hooks.QueryStatus.Invoke()(static.root)

	if err != nil {
		return &types.NavigateResult{
			Err: err,
		}, err
	}

	node := core.Root(static.root, info)

	_, _ = static.mediator.impl.Traverse(ctx, static, node)

	return &types.NavigateResult{}, nil
}

func traverse(_ navigationStatic) {

}
