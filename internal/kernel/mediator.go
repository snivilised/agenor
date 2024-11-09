package kernel

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/level"
	"github.com/snivilised/traverse/life"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/stock"
)

// mediator controls traversal events, sends notifications and emits
// life-cycle events
type mediator struct {
	tree         string
	subscription enums.Subscription
	impl         NavigatorImpl
	guardian     *guardian
	periscope    *level.Periscope
	o            *pref.Options
	resources    *enclave.Resources
	metrics      core.Metrics
	order        []enums.Role
}

type mediatorInfo struct {
	facade       pref.Facade
	subscription enums.Subscription
	o            *pref.Options
	impl         NavigatorImpl
	sealer       enclave.GuardianSealer
	resources    *enclave.Resources
}

func newMediator(info *mediatorInfo) *mediator {
	metrics := info.resources.Supervisor.Many(
		enums.MetricNoFilesInvoked,
		enums.MetricNoDirectoriesInvoked,
		enums.MetricNoChildFilesFound,
	)

	return &mediator{
		tree:         info.facade.Path(), // TODO: ??? is this right for resume?
		subscription: info.subscription,
		impl:         info.impl,
		guardian: newGuardian(&guardianInfo{
			subscription: info.subscription,
			client:       info.facade.Client(),
			master:       info.sealer,
			metrics:      metrics,
		}),
		periscope: level.New(),
		o:         info.o,
		resources: info.resources,
		metrics:   metrics,
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

func (m *mediator) Decorate(link enclave.Link) error {
	return m.guardian.Decorate(link)
}

func (m *mediator) Unwind(role enums.Role) error {
	return m.guardian.Unwind(role)
}

func (m *mediator) Arrange(active, order []enums.Role) {
	m.order = order
	m.guardian.arrange(active, order)
}

func (m *mediator) Ignite(ignition *enclave.Ignition) {
	m.impl.Ignite(ignition)
	m.resources.Binder.Controls.Begin.Dispatch()(&life.BeginState{
		Tree: m.tree,
	})
}

func (m *mediator) Conclude(result core.TraverseResult) {
	m.resources.Binder.Controls.End.Dispatch()(result)
}

func (m *mediator) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	result, err := m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		tree:     m.tree,
		calc:     m.resources.Forest.T.Calc(),
	})

	if !stock.IsBenignError(err) && m.o != nil {
		m.o.Monitor.Log.Error(err.Error())
	}

	return result, err
}

func (m *mediator) Read(path string) ([]fs.DirEntry, error) {
	return m.o.Hooks.ReadDirectory.Invoke()(m.resources.Forest.T, path)
}

func (m *mediator) Resume(ctx context.Context,
	active *core.ActiveState,
) (*enclave.KernelResult, error) {
	m.tree = active.Tree
	// TODO: there is something missing here...
	// we need to do more with the loaded active state
	//
	// - mute notifications
	// - load the periscope with an adjusted depth from active state
	//
	return m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		tree:     active.Tree,
		calc:     m.resources.Forest.T.Calc(),
	})
}

// Bridge combines information gleaned from the previous traversal that was
// interrupted, into the resume traversal
//
// tree, current string
func (m *mediator) Bridge(tree, current string) {
	m.tree = tree
	fmt.Printf("---> mediator.Bridge - tree %q, current %q\n", tree, current)
}

func (m *mediator) Spawn(ctx context.Context,
	active *core.ActiveState, // TODO: this should not be ActiveState, ActiveState is being abused
) (*enclave.KernelResult, error) {
	m.tree = active.Tree
	offset := 0 // TODO: not sure what to set this to yet
	m.periscope = level.Restore(offset, active.Depth)

	return m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		tree:     active.Tree,
		calc:     m.resources.Forest.T.Calc(),
	})
}

func (m *mediator) Invoke(servant core.Servant,
	inspection enclave.Inspection,
) error {
	return m.guardian.Invoke(servant, inspection)
}

func (m *mediator) Supervisor() *core.Supervisor {
	return m.resources.Supervisor
}
