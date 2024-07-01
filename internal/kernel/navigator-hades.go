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

func (n *navigatorHades) Result(_ context.Context, err error) *types.KernelResult {
	return types.NewFailed(err)
}

func (n *navigatorHades) Starting(_ core.Session) {
}

func (n *navigatorHades) Navigate(ctx context.Context) (core.TraverseResult, error) {
	return n.Result(ctx, n.err), n.err
}
