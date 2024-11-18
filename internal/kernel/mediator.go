package kernel

import (
	"context"
	"io/fs"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/level"
	"github.com/snivilised/agenor/life"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/stock"
)

// mediator controls traversal, sends notifications and emits
// life-cycle events
type mediator struct {
	tree         string
	subscription enums.Subscription
	facade       pref.Facade
	impl         NavigatorImpl
	guardian     *guardian
	periscope    *level.Periscope
	o            *pref.Options
	resources    *enclave.Resources
	metrics      core.Metrics
	order        []enums.Role
}

func NewMediator(inception *Inception,
	sealer enclave.GuardianSealer,
) (enclave.Mediator, error) {
	o := inception.Harvest.Options()
	facade := inception.Facade
	resources := inception.Resources
	impl, err := newImpl(o, inception)

	metrics := resources.Supervisor.Many(
		enums.MetricNoFilesInvoked,
		enums.MetricNoDirectoriesInvoked,
		enums.MetricNoChildFilesFound,
	)

	return &mediator{
		tree:         inception.NavigationTree(),
		subscription: inception.Subscription,
		facade:       facade,
		impl:         impl,
		guardian: newGuardian(&guardianInfo{
			subscription: inception.Subscription,
			client:       facade.Client(),
			master:       sealer,
			metrics:      metrics,
		}),
		periscope: level.New(),
		o:         o,
		resources: resources,
		metrics:   metrics,
	}, err
}

func (m *mediator) Decorate(link enclave.Link) error {
	return m.guardian.Decorate(link)
}

func (m *mediator) Unwind(role enums.Role) error {
	return m.guardian.Unwind(role)
}

func (m *mediator) Invoke(servant core.Servant,
	inspection enclave.Inspection,
) error {
	return m.guardian.Invoke(servant, inspection)
}

func (m *mediator) Arrange(active, order []enums.Role) {
	m.order = order
	m.guardian.arrange(active, order)
}

func (m *mediator) Read(path string) ([]fs.DirEntry, error) {
	return m.o.Hooks.ReadDirectory.Invoke()(m.resources.Forest.T, path)
}

func (m *mediator) Spawn(ctx context.Context,
	tree string,
) (*enclave.KernelResult, error) {
	// NB:
	return m.impl.Top(ctx, &navigationStatic{
		mediator:     m,
		tree:         tree,
		calc:         m.resources.Forest.T.Calc(),
		subscription: m.subscription,
		magnitude:    m.facade.Magnitude(),
	})
}

// Bridge combines information gleaned from the previous traversal that was
// interrupted, into the resume traversal
func (m *mediator) Bridge(active *core.ActiveState) {
	m.tree = active.Tree
	m.periscope.Offset(active.Depth)
	m.Supervisor().Load(active.Metrics)
}

func (m *mediator) Supervisor() *enclave.Supervisor {
	return m.resources.Supervisor
}

func (m *mediator) Navigate(ctx context.Context) (result *enclave.KernelResult, err error) {
	result, err = m.impl.Top(ctx, &navigationStatic{
		mediator:     m,
		tree:         m.tree,
		calc:         m.resources.Forest.T.Calc(),
		subscription: m.subscription,
		magnitude:    m.facade.Magnitude(),
	})

	if !stock.IsBenignError(err) && m.o != nil {
		m.o.Monitor.Log.Error(err.Error())
	}

	return result, err
}

func (m *mediator) Ignite(ignition *enclave.Ignition) {
	m.impl.Ignite(ignition)
	m.resources.Binder.Controls.Begin.Dispatch()(&life.BeginState{
		Tree: m.tree,
	})
}

func (m *mediator) Result(ctx context.Context,
) *enclave.KernelResult {
	return m.impl.Result(ctx)
}

func (m *mediator) Snooze(ctx context.Context,
	active *core.ActiveState,
) (result *enclave.KernelResult, err error) {
	m.tree = active.Tree

	return m.impl.Top(ctx, &navigationStatic{
		mediator:     m,
		tree:         active.Tree,
		calc:         m.resources.Forest.T.Calc(),
		subscription: m.subscription,
		magnitude:    m.facade.Magnitude(),
	})
}

func (m *mediator) Conclude(result core.TraverseResult) {
	m.resources.Binder.Controls.End.Dispatch()(result)
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
