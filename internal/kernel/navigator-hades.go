package kernel

import (
	"context"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/stock"
)

func HadesNav(o *pref.Options, err error) enclave.KernelController {
	return &navigatorHades{
		o:   o,
		err: err,
	}
}

type navigatorHades struct {
	o   *pref.Options
	err error
}

func (n *navigatorHades) Rank() {
}

func (n *navigatorHades) Ignite(*enclave.Ignition) {
}

func (n *navigatorHades) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	return n.Result(ctx, n.err), n.err
}

func (n *navigatorHades) Result(_ context.Context, err error) *enclave.KernelResult {
	if !stock.IsBenignError(err) && n.o != nil {
		n.o.Monitor.Log.Error(err.Error())
	}

	return enclave.NewFailed(err)
}

func (n *navigatorHades) Mediator() enclave.Mediator {
	return nil
}

func (n *navigatorHades) Resume(context.Context,
	*core.ActiveState,
) (*enclave.KernelResult, error) {
	return &enclave.KernelResult{}, nil
}

func (n *navigatorHades) Conclude(_ core.TraverseResult) {
}
