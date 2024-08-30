package kernel

import (
	"github.com/snivilised/traverse/collections"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
)

type (
	invocationChain   = map[enums.Role]types.Link
	positionalRoleSet = collections.PositionalSet[enums.Role]
)

// anchor is a specialised link that should always be the
// last in the chain and contains the original client's handler.
type anchor struct {
	client core.Client
	owner  measure.Owned
}

func (t *anchor) Next(node *core.Node, _ core.Inspection) (bool, error) {
	if metric := lo.Ternary(node.IsFolder(),
		t.owner.Mums[enums.MetricNoFoldersInvoked],
		t.owner.Mums[enums.MetricNoFilesInvoked],
	); metric != nil {
		metric.Tick()
	}

	return false, t.client(node)
}

func (t *anchor) Role() enums.Role {
	return enums.RoleAnchor
}

type iterationContainer struct {
	invoker   Invokable
	positions *positionalRoleSet
	chain     invocationChain
}

// guardian controls access to the client callback
type guardian struct {
	container iterationContainer
	master    types.GuardianSealer
	anchor    *anchor
}

func newGuardian(client core.Client,
	master types.GuardianSealer,
	mums measure.MutableMetrics,
) *guardian {
	anchor := &anchor{
		client: client,
		owner: measure.Owned{
			Mums: mums,
		},
	}

	return &guardian{
		container: iterationContainer{
			chain: make(invocationChain),
		},
		master: master,
		anchor: anchor,
	}
}

func (g *guardian) arrange(active, order []enums.Role) {
	g.container.chain[enums.RoleAnchor] = g.anchor

	if len(active) == 0 {
		g.container.invoker = NodeInvoker(func(node *core.Node, inspection core.Inspection) error {
			_, err := g.anchor.Next(node, inspection)
			return err
		})

		return
	}

	g.container.positions = collections.NewPositionalSet(order, enums.RoleAnchor)
	g.container.invoker = NodeInvoker(func(node *core.Node, inspection core.Inspection) error {
		return g.iterate(node, inspection)
	})
}

// role indicates the guise under which the decorator is being applied.
// Not all roles can be decorated (sealed). The fastward-resume decorator is
// sealed. If an attempt is made to Decorate a sealed decorator,
// an error is returned.
func (g *guardian) Decorate(link types.Link) error {
	top := g.container.chain[g.container.positions.Top()]

	if g.master.IsSealed(top) {
		return core.ErrGuardianCantDecorateItemSealed
	}

	role := link.Role()
	g.container.chain[role] = link
	g.container.positions.Insert(role)

	return nil
}

func (g *guardian) Unwind(role enums.Role) error {
	if role == enums.RoleAnchor {
		return nil
	}

	delete(g.container.chain, role)
	g.container.positions.Delete(role)

	// TODO: required only for fastward resume or hibernation
	//
	return nil
}

// Invoke executes the chain which may or may not end up resulting in
// the invocation of the client's callback, depending on the contents
// of the chain.
func (g *guardian) Invoke(node *core.Node, inspection core.Inspection) error {
	return g.container.invoker.Invoke(node, inspection)
}

func (g *guardian) iterate(node *core.Node, inspection core.Inspection) error {
	for _, role := range g.container.positions.Items() {
		link := g.container.chain[role]

		if next, err := link.Next(node, inspection); !next || err != nil {
			return err
		}
	}

	return nil
}

// Benign is used when a master sealer has not been registered. It is
// permissive in nature.
type Benign struct {
}

func (m *Benign) Seal(types.Link) error {
	return nil
}

func (m *Benign) IsSealed(types.Link) bool {
	return false
}
