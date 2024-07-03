package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
)

type navigatorFiles struct {
	navigator
}

func (n *navigatorFiles) Top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	return top(ctx, ns)
}

func (n *navigatorFiles) Travel(ctx context.Context,
	ns *navigationStatic,
	current *core.Node,
) (bool, error) {
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
		return false, ns.mediator.Invoke(current)
	}

	if skip, e := ns.mediator.o.Defects.Skip.Ask(
		current, vapour.contents(), err,
	); skip == enums.SkipAllTraversal || err != nil {
		return continueTraversal, e
	} else if skip == enums.SkipDirTraversal {
		return true, e
	}

	return travel(ctx, ns, vapour)
}

func (n *navigatorFiles) inspect(ns *navigationStatic, current *core.Node) (inspection, error) {
	var (
		vapour = &navigationVapour{
			ns:      ns,
			present: current,
		}
		err error
	)

	if vapour.present.IsFolder() {
		vapour.cargo, err = read(ns.mediator.resources.FS.N,
			ns.mediator.o,
			current.Path,
		)

		vapour.cargo.Sort(enums.EntryTypeAll)
		vapour.pick(enums.EntryTypeAll)
	} else {
		vapour.clear()
	}

	extend(ns, vapour)

	return vapour, err
}
