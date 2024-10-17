package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
)

type navigatorFolders struct {
	navigatorAgent
}

func (n *navigatorFolders) Top(ctx context.Context,
	ns *navigationStatic,
) (*types.KernelResult, error) {
	return n.top(ctx, ns)
}

func (n *navigatorFolders) Traverse(ctx context.Context,
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

func (n *navigatorFolders) inspect(ns *navigationStatic,
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

	// for the folders navigator, we ignore the user defined setting in
	// (Options).Core.Behaviours.Sort.DirectoryEntryOrder, as we're only
	// interested in folders and therefore forced to use
	// enums.DirectoryEntryOrderFoldersFirst instead.
	//
	vapour.cargo, err = read(ns.mediator.resources.FS.T,
		n.ro,
		current.Path,
	)

	vapour.Sort(enums.EntryTypeFolder)
	vapour.Pick(enums.EntryTypeFolder)

	extend(ns, vapour)

	return vapour, err
}
