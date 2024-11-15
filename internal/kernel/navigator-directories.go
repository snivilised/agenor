package kernel

import (
	"context"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
)

type navigatorDirectories struct {
	navigatorAgent
}

func (n *navigatorDirectories) Top(ctx context.Context,
	ns *navigationStatic,
) (*enclave.KernelResult, error) {
	return n.top(ctx, ns)
}

func (n *navigatorDirectories) Traverse(ctx context.Context,
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

	vapour, err := n.inspect(ns, servant)

	if e := ns.mediator.Invoke(servant, vapour); e != nil {
		return continueTraversal, e
	}

	if skip, e := ns.mediator.o.Defects.Skip.Ask(
		current, vapour.Contents(), err,
	); skip == enums.SkipAllTraversal {
		return continueTraversal, e
	}

	return n.travel(ctx, ns, vapour)
}

func (n *navigatorDirectories) inspect(ns *navigationStatic,
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

	// for the directories navigator, we ignore the user defined setting in
	// (Options).Core.Behaviours.Sort.DirectoryEntryOrder, as we're only
	// interested in directories and therefore forced to use
	// NavigationBehaviours.SortBehaviour.SortFilesFirst=true instead.
	//
	vapour.cargo, err = read(ns.mediator.resources.Forest.T,
		n.ro,
		current.Path,
	)

	vapour.Sort(enums.EntryTypeDirectory)
	vapour.Pick(enums.EntryTypeDirectory)

	extend(ns, vapour)

	return vapour, err
}
