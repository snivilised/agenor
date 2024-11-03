package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/enclave"
)

type NavigationController struct {
	med *mediator
}

func newNavigationController(med *mediator) *NavigationController {
	return &NavigationController{
		med: med,
	}
}

func (nc *NavigationController) Register(enclave.Plugin) error {
	return nil
}

func (nc *NavigationController) Ignite(ignition *enclave.Ignition) {
	nc.med.Ignite(ignition)
}

func (nc *NavigationController) Resume(ctx context.Context,
	active *core.ActiveState,
) (*enclave.KernelResult, error) {
	return nc.med.Resume(ctx, active)
}

func (nc *NavigationController) Conclude(result core.TraverseResult) {
	nc.med.Conclude(result)
}

func (nc *NavigationController) Impl() NavigatorImpl {
	return nc.med.impl
}

func (nc *NavigationController) Navigate(ctx context.Context,
) (*enclave.KernelResult, error) {
	return nc.med.Navigate(ctx)
}

func (nc *NavigationController) Result(ctx context.Context,
	err error,
) *enclave.KernelResult {
	return nc.med.impl.Result(ctx, err)
}

func (nc *NavigationController) Mediator() enclave.Mediator {
	return nc.med
}
