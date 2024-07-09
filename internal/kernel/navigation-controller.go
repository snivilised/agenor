package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
)

type NavigationController struct {
	Mediator *mediator
}

func (nc *NavigationController) Register(types.Plugin) error {
	return nil
}

func (nc *NavigationController) Ignite(ignition *types.Ignition) {
	nc.Mediator.Ignite(ignition)
}

func (nc *NavigationController) Impl() NavigatorImpl {
	return nc.Mediator.impl
}

func (nc *NavigationController) Navigate(ctx context.Context) (core.TraverseResult, error) {
	return nc.Mediator.Navigate(ctx)
}

func (nc *NavigationController) Result(ctx context.Context, err error) *types.KernelResult {
	return nc.Mediator.impl.Result(ctx, err)
}
