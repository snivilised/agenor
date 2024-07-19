package refine

import (
	"github.com/snivilised/traverse/core"
)

// 📚 package: refine defines filters and should be used by rx to alter observables

type NavigationFilters struct {
	// Node denotes the filter object that represents the Node file system item
	// being visited.
	//
	Node core.TraverseFilter

	// Children denotes the Compound filter that is applied to the direct descendants
	// of the current file system item being visited.
	//
	Children core.ChildTraverseFilter
}
