package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func HadesNav(o *pref.Options, err error) types.KernelController {
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

func (n *navigatorHades) Ignite(*types.Ignition) {
}

func (n *navigatorHades) Navigate(ctx context.Context) (*types.KernelResult, error) {
	return n.Result(ctx, n.err), n.err
}

func (n *navigatorHades) Result(_ context.Context, err error) *types.KernelResult {
	if !IsBenignError(err) && n.o != nil {
		n.o.Monitor.Log.Error(err.Error())
	}

	return types.NewFailed(err)
}

func (n *navigatorHades) Mediator() types.Mediator {
	return nil
}

func (n *navigatorHades) Resume(context.Context,
	*core.ActiveState,
) (*types.KernelResult, error) {
	return &types.KernelResult{}, nil
}

func (n *navigatorHades) Conclude(_ core.TraverseResult) {
}
