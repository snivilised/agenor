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

func (n *navigatorFolders) Travel(ctx context.Context,
	ns *navigationStatic,
	current *core.Node,
) (bool, error) {
	vapour, err := n.inspect(ns, current)

	if e := ns.mediator.Invoke(current); e != nil {
		return continueTraversal, e
	}

	if skip, e := ns.mediator.o.Defects.Skip.Ask(
		current, vapour.contents(), err,
	); skip == enums.SkipAllTraversal {
		return continueTraversal, e
	}

	return n.travel(ctx, ns, vapour)
}

func (n *navigatorFolders) inspect(ns *navigationStatic,
	current *core.Node,
) (inspection, error) {
	var (
		vapour = &navigationVapour{
			ns:      ns,
			present: current,
		}
		err error
	)

	// for the folders navigator, we ignore the user defined setting in
	// n.o.Store.Behaviours.Sort.DirectoryEntryOrder, as we're only interested in
	// folders and therefore force to use DirectoryEntryOrderFoldersFirstEn instead
	//
	vapour.cargo, err = read(ns.mediator.resources.FS.N,
		n.ro,
		current.Path,
	)

	vapour.cargo.Sort(enums.EntryTypeFolder)
	vapour.pick(enums.EntryTypeFolder)

	// TODO: implement directory with files

	extend(ns, vapour)

	return vapour, err
}
