package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
)

type NavigationController struct {
	mediator *mediator
}

func (nc *NavigationController) Navigate(ctx context.Context) (core.TraverseResult, error) {
	return nc.mediator.Navigate(ctx)
}

func (nc *NavigationController) Impl() NavigatorImpl {
	return nc.mediator.impl
}

func (nc *NavigationController) Register(types.Plugin) error {
	return nil
}
