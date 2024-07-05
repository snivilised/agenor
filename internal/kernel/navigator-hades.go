package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
)

func HadesNav(err error) types.KernelController {
	return &navigatorHades{
		err: err,
	}
}

type navigatorHades struct {
	err error
}

func (n *navigatorHades) Rank() {
}

func (n *navigatorHades) Ignite(*types.Ignition) {
}

func (n *navigatorHades) Navigate(ctx context.Context) (core.TraverseResult, error) {
	return n.Result(ctx, n.err), n.err
}

func (n *navigatorHades) Result(_ context.Context, err error) *types.KernelResult {
	return types.NewFailed(err)
}
