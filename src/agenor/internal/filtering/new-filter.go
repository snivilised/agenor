package filtering

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/pref"
)

// New creates a new filter based on the given definition and options.
func New(definition *core.FilterDef,
	fo *pref.FilterOptions,
) (core.TraverseFilter, error) {
	return OrFuncE(
		func() (core.TraverseFilter, error) {
			return buildPolyNodeFilter(definition, fo, buildNativeNodeFilter, getCustomFilter)
		},
		func() (core.TraverseFilter, error) {
			return getCustomFilter(definition, fo)
		},
		func() (core.TraverseFilter, error) {
			return buildNativeNodeFilter(definition)
		},
	)
}
