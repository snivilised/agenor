package filtering

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/pref"
)

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
