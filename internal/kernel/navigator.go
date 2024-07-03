package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type navigator struct {
	o         *pref.Options
	using     *pref.Using
	resources *types.Resources
	session   core.Session
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

func (n *navigator) descend(*navigationInfo) bool {
	return true
}

func (n *navigator) ascend(navi *navigationInfo, permit bool) {
	_, _ = navi, permit
}

func (n *navigator) Starting(session core.Session) {
	n.session = session
}

func (n *navigator) Top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	_, _ = ctx, ns

	// this is duff, will be removed after all the navigators
	// have been implemented
	return n.Result(ctx, nil), nil
}

func (n *navigator) Travel(context.Context,
	*navigationStatic,
	*core.Node,
) (bool, error) {
	return continueTraversal, nil
}

func (n *navigator) Result(ctx context.Context, err error) *types.KernelResult {
	complete := n.session.IsComplete()
	result := types.NewResult(n.session,
		n.resources.Supervisor,
		err,
		complete,
	)

	if complete {
		_ = services.Broker.Emit(ctx, services.TopicNavigationComplete, result)
	}

	return result
}
