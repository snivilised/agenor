package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/level"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
	"github.com/snivilised/traverse/pref"
)

// mediator controls traversal events, sends notifications and emits
// life-cycle events
type mediator struct {
	root      string
	using     *pref.Using
	impl      NavigatorImpl
	guardian  *guardian
	frame     *navigationFrame
	pad       *scratchPad // gets created just before nav begins
	o         *pref.Options
	resources *types.Resources
}

func newMediator(using *pref.Using,
	o *pref.Options,
	impl NavigatorImpl,
	sealer types.GuardianSealer,
	resources *types.Resources,
) *mediator {
	return &mediator{
		root:  using.Root,
		using: using,
		impl:  impl,
		guardian: newGuardian(using.Handler, sealer, resources.Supervisor.Many(
			enums.MetricNoFilesInvoked,
			enums.MetricNoFoldersInvoked,
		)),
		frame: &navigationFrame{
			periscope: level.New(),
		},
		pad:       newScratch(o),
		o:         o,
		resources: resources,
	}
}

func (m *mediator) Decorate(link types.Link) error {
	return m.guardian.Decorate(link)
}

func (m *mediator) Unwind(role enums.Role) error {
	return m.guardian.Unwind(role)
}

func (m *mediator) Arrange(activeRoles []enums.Role) {
	m.guardian.arrange(activeRoles)
}

func (m *mediator) Ignite(ignition *types.Ignition) {
	m.impl.Ignite(ignition)
}

func (m *mediator) Navigate(ctx context.Context) (core.TraverseResult, error) {
	result, err := m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		root:     m.root,
	})

	return result, err
}

func (m *mediator) Spawn(ctx context.Context, root string) (core.TraverseResult, error) {
	// TODO: send a message indicating spawn
	//
	return m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		root:     root,
	})
}

func (m *mediator) Invoke(node *core.Node) error {
	return m.guardian.Invoke(node)
}

func (m *mediator) Supervisor() *measure.Supervisor {
	return m.resources.Supervisor
}
