package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/level"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/life"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/stock"
)

// mediator controls traversal events, sends notifications and emits
// life-cycle events
type mediator struct {
	tree         string
	subscription enums.Subscription
	using        *pref.Using
	impl         NavigatorImpl
	guardian     *guardian
	periscope    *level.Periscope
	o            *pref.Options
	resources    *types.Resources
	mums         measure.MutableMetrics
	order        []enums.Role
}

type mediatorInfo struct {
	using     *pref.Using
	o         *pref.Options
	impl      NavigatorImpl
	sealer    types.GuardianSealer
	resources *types.Resources
}

func newMediator(info *mediatorInfo) *mediator {
	mums := info.resources.Supervisor.Many(
		enums.MetricNoFilesInvoked,
		enums.MetricNoDirectoriesInvoked,
		enums.MetricNoChildFilesFound,
	)

	return &mediator{
		tree:         info.using.Tree,
		subscription: info.using.Subscription,
		using:        info.using,
		impl:         info.impl,
		guardian: newGuardian(&guardianInfo{
			subscription: info.using.Subscription,
			client:       info.using.Handler,
			master:       info.sealer,
			mums:         mums,
		}),
		periscope: level.New(),
		o:         info.o,
		resources: info.resources,
		mums:      mums,
	}
}

func (m *mediator) descend(node *core.Node) bool {
	if !m.periscope.Descend(m.o.Behaviours.Cascade.Depth) {
		return false
	}

	m.resources.Binder.Controls.Descend.Dispatch()(node)

	return true
}

func (m *mediator) ascend(node *core.Node, permit bool) {
	if permit {
		m.periscope.Ascend()
		m.resources.Binder.Controls.Ascend.Dispatch()(node)
	}
}

func (m *mediator) Decorate(link types.Link) error {
	return m.guardian.Decorate(link)
}

func (m *mediator) Unwind(role enums.Role) error {
	return m.guardian.Unwind(role)
}

func (m *mediator) Arrange(active, order []enums.Role) {
	m.order = order
	m.guardian.arrange(active, order)
}

func (m *mediator) Ignite(ignition *types.Ignition) {
	m.impl.Ignite(ignition)
	m.resources.Binder.Controls.Begin.Dispatch()(&life.BeginState{
		Tree: m.tree,
	})
}

func (m *mediator) Conclude(result core.TraverseResult) {
	m.resources.Binder.Controls.End.Dispatch()(result)
}

func (m *mediator) Navigate(ctx context.Context) (*types.KernelResult, error) {
	result, err := m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		tree:     m.tree,
	})

	if !stock.IsBenignError(err) && m.o != nil {
		m.o.Monitor.Log.Error(err.Error())
	}

	return result, err
}

func (m *mediator) Resume(ctx context.Context,
	active *core.ActiveState,
) (*types.KernelResult, error) {
	// TODO: there is something missing here...
	// we need to do more with the loaded active state
	//
	// - mute notifications
	// - combine metrics
	// - load the periscope with an adjusted depth from active state
	// - we might need to define a callback param for the strategy
	//
	return m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		tree:     active.Tree,
	})
}

func (m *mediator) Spawn(ctx context.Context,
	active *core.ActiveState,
) (*types.KernelResult, error) {
	// TODO: send a message indicating spawn
	// we need to reset the active state, eg synchronise
	// the depth
	return m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		tree:     active.Tree,
	})
}

func (m *mediator) Invoke(servant core.Servant,
	inspection types.Inspection,
) error {
	return m.guardian.Invoke(servant, inspection)
}

func (m *mediator) Supervisor() *measure.Supervisor {
	return m.resources.Supervisor
}
