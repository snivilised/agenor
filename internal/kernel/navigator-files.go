package kernel

import (
	"context"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
)

type navigatorFiles struct {
	navigatorAgent
}

func (n *navigatorFiles) Top(ctx context.Context,
	ns *navigationStatic,
) (*enclave.KernelResult, error) {
	return n.top(ctx, ns)
}

func (n *navigatorFiles) Traverse(ctx context.Context,
	ns *navigationStatic,
	servant core.Servant,
) (bool, error) {
	current := servant.Node()
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
	vapour, err := n.inspect(ns, servant)

	if !current.IsDirectory() {
		// Effectively, this is the file only filter
		//
		return false, ns.mediator.Invoke(servant, vapour)
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
	servant core.Servant,
) (inspection, error) {
	var (
		current = servant.Node()
		vapour  = &navigationVapour{
			ns:      ns,
			present: current,
		}
		err error
	)

	if vapour.present.IsDirectory() {
		vapour.cargo, err = read(ns.mediator.resources.Forest.T,
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
