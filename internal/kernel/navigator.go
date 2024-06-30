package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type navigator struct {
	o     *pref.Options
	using *pref.Using
	res   *types.Resources
}

/*
func (n *navigator) descend(navi *NavigationInfo) bool {
	if !navi.frame.periscope.descend(n.o.Store.Behaviours.Cascade.Depth) {
		return false
	}

	navi.frame.notifiers.descend.invoke(navi.Item)

	return true
}

func (n *navigator) ascend(navi *NavigationInfo, permit bool) {
	if permit {
		navi.frame.periscope.ascend()
		navi.frame.notifiers.ascend.invoke(navi.Item)
	}
}
*/

func (n *navigator) descend(navi *navigationInfo) bool {
	_ = navi

	return true
}

func (n *navigator) ascend(navi *navigationInfo, permit bool) {
	_, _ = navi, permit
}

func (n *navigator) Top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	_, _ = ctx, ns

	return &types.KernelResult{}, nil
}

func (n *navigator) Travel(context.Context,
	*navigationStatic,
	*core.Node,
) (bool, error) {
	return continueTraversal, nil
}
