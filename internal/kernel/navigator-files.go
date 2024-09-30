package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
)

type navigatorFiles struct {
	navigatorAgent
}

func (n *navigatorFiles) Top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	return n.top(ctx, ns)
}

func (n *navigatorFiles) Traverse(ctx context.Context,
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

	// TODO: check this comment
	// For files, the registered callback will only be invoked for file entries. This means
	// that the client will have no way to skip the descending of a particular directory. In
	// this case, the client should use the OnDescend callback (yet to be implemented) and
	// return SkipDir from there.
	//
	vapour, err := n.inspect(ns, current)

	if !current.IsFolder() {
		// Effectively, this is the file only filter
		//
		return false, ns.mediator.Invoke(current, vapour)
	}

	if skip, e := ns.mediator.o.Defects.Skip.Ask(
		current, vapour.Contents(), err,
	); skip == enums.SkipAllTraversal || err != nil {
		return continueTraversal, e
	} else if skip == enums.SkipDirTraversal {
		return true, e
	}

	return n.travel(ctx, ns, vapour)
}

func (n *navigatorFiles) inspect(ns *navigationStatic,
	current *core.Node,
) (inspection, error) {
	var (
		vapour = &navigationVapour{
			ns:      ns,
			present: current,
		}
		err error
	)

	if vapour.present.IsFolder() {
		vapour.cargo, err = read(ns.mediator.resources.FS.T,
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
