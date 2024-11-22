package kernel

import (
	"github.com/snivilised/agenor/collections"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/third/lo"
)

type (
	invocationChain    = map[enums.Role]enclave.Link
	positionalRoleSet  = collections.PositionalSet[enums.Role]
	iterationContainer struct {
		invoker   Invokable
		positions *positionalRoleSet
		chain     invocationChain
	}
)

type (
	NodeInvoker func(servant core.Servant, inspection enclave.Inspection) error
)

func (fn NodeInvoker) Invoke(servant core.Servant,
	inspection enclave.Inspection,
) error {
	return fn(servant, inspection)
}

// anchor is a specialised link that should always be the
// last in the chain and contains the original client's handler.
type anchor struct {
	subscription enums.Subscription
	client       core.Client
	crate        enclave.Crate
}

func (a *anchor) Next(servant core.Servant, _ enclave.Inspection) (bool, error) {
	node := servant.Node()
	if metric := lo.Ternary(node.IsDirectory(),
		a.crate.Metrics[enums.MetricNoDirectoriesInvoked],
		a.crate.Metrics[enums.MetricNoFilesInvoked],
	); metric != nil {
		metric.Tick()
	}

	return false, a.client(servant)
}

func (a *anchor) Role() enums.Role {
	return enums.RoleAnchor
}

func (a *anchor) swap(decorator core.Client) core.Client {
	swap := a.client
	a.client = decorator

	return swap
}

// guardian controls access to the client callback
type guardian struct {
	container iterationContainer
	master    enclave.GuardianSealer
	anchor    *anchor
}

type guardianInfo struct {
	subscription enums.Subscription
	client       core.Client
	master       enclave.GuardianSealer
	metrics      core.Metrics
}

func newGuardian(info *guardianInfo) *guardian {
	return &guardian{
		container: iterationContainer{
			chain: make(invocationChain),
		},
		master: info.master,
		anchor: &anchor{
			subscription: info.subscription,
			client:       info.client,
			crate: enclave.Crate{
				Metrics: info.metrics,
			},
		},
	}
}

func (g *guardian) arrange(active, order []enums.Role) {
	g.container.chain[enums.RoleAnchor] = g.anchor

	if len(active) == 0 {
		g.container.invoker = NodeInvoker(func(servant core.Servant,
			inspection enclave.Inspection,
		) error {
			_, err := g.anchor.Next(servant, inspection)
			return err
		})

		return
	}

	g.container.positions = collections.NewPositionalSet(order, enums.RoleAnchor)
	g.container.invoker = NodeInvoker(func(servant core.Servant,
		inspection enclave.Inspection,
	) error {
		return g.iterate(servant, inspection)
	})
}

// role indicates the guise under which the decorator is being applied.
// Not all roles can be decorated (sealed). The fastward-resume decorator is
// sealed. If an attempt is made to Decorate a sealed decorator,
// an error is returned.
func (g *guardian) Decorate(link enclave.Link) error {
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
func (g *guardian) Invoke(servant core.Servant, inspection enclave.Inspection) error {
	return g.container.invoker.Invoke(servant, inspection)
}

func (g *guardian) iterate(servant core.Servant, inspection enclave.Inspection) error {
	for _, role := range g.container.positions.Items() {
		link := g.container.chain[role]

		if next, err := link.Next(servant, inspection); !next || err != nil {
			return err
		}
	}

	return nil
}

func (g *guardian) Swap(decorator core.Client) {
	_ = g.anchor.swap(decorator)
}

// Benign is used when a master sealer has not been registered. It is
// permissive in nature.
type Benign struct {
}

func (m *Benign) Seal(enclave.Link) error {
	return nil
}

func (m *Benign) IsSealed(enclave.Link) bool {
	return false
}
