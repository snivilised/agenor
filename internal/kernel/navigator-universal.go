package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
)

type navigatorUniversal struct {
	navigatorAgent
}

func (n *navigatorUniversal) Top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	return n.top(ctx, ns)
}

func (n *navigatorUniversal) Traverse(ctx context.Context,
	ns *navigationStatic,
	current *core.Node,
) (bool, error) {
	descended := ns.mediator.descend(current)

	defer func(permit bool) {
		ns.mediator.ascend(current, permit)
	}(descended)

	if !descended {
		return continueTraversal, nil
	}

	vapour, err := n.inspect(ns, current)

	if e := ns.mediator.Invoke(current); e != nil {
		return continueTraversal, e
	}

	if skip, e := ns.mediator.o.Defects.Skip.Ask(
		current, vapour.Contents(), err,
	); skip == enums.SkipAllTraversal {
		return continueTraversal, e
	} else if skip == enums.SkipDirTraversal {
		return skipTraversal, e
	}

	return n.travel(ctx, ns, vapour)
}

func (n *navigatorUniversal) inspect(ns *navigationStatic, current *core.Node) (inspection, error) {
	var (
		vapour = &navigationVapour{
			ns:      ns,
			present: current,
		}
		err error
	)

	if current.IsFolder() {
		vapour.cargo, err = read(ns.mediator.resources.FS.N,
			n.ro,
			current.Path,
		)

		vapour.Sort(enums.EntryTypeAll)
		vapour.Pick(enums.EntryTypeAll)
	} else {
		vapour.clear()
	}

	extend(ns, vapour)

	return vapour, err
}
