package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
)

type NavigationController struct {
	Med *mediator
}

func (nc *NavigationController) Register(types.Plugin) error {
	return nil
}

func (nc *NavigationController) Ignite(ignition *types.Ignition) {
	nc.Med.Ignite(ignition)
}

func (nc *NavigationController) Conclude(result core.TraverseResult) {
	nc.Med.Conclude(result)
}

func (nc *NavigationController) Impl() NavigatorImpl {
	return nc.Med.impl
}

func (nc *NavigationController) Navigate(ctx context.Context,
) (core.TraverseResult, error) {
	return nc.Med.Navigate(ctx)
}

func (nc *NavigationController) Result(ctx context.Context,
	err error,
) *types.KernelResult {
	return nc.Med.impl.Result(ctx, err)
}

func (nc *NavigationController) Mediator() types.Mediator {
	return nc.Med
}
