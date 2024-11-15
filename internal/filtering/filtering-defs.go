package filtering

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/pref"
)

// 📦 pkg: filtering - this package is required because filters are required
// not just but the filter plugin, but others too like hibernation. The filter
// required by hibernation could have been implemented by the filter plugin,
// but doing so in this fashion would have mean introducing coupling of
// hibernation on filter; ie how to allow hibernation to access the filter(s)
// created by filter?
//	Instead, we factor out the filter creation code to this package, so that
// hibernation can create and apply filters as it needs, without depending on
// filter. So filter, now doesn't own the filter implementations, rather it's
// simply responsible for the plugin aspects of filtering, not implementation
// or creation.
//

// OrFuncE returns the first func that returns a value that is not equal to the
// zero value and does not return an error. If no argument is non-zero, it returns
// the zero value. All functions, must return an error value
func OrFuncE[T comparable](funcs ...func() (T, error)) (T, error) {
	var zero T
	for _, fn := range funcs {
		result, err := fn()
		if err != nil {
			return zero, err
		}

		if result != zero {
			return result, err
		}
	}

	return zero, nil
}

type (
	// filterNativeFunc implies that the filter has to be constructed from the
	// filter definition only.
	filterNativeFunc func(definition *core.FilterDef) (core.TraverseFilter, error)

	// filterUsingOptionsFunc implies that the filter options object is required
	// to obtain the filter. The filter may be created or just retrieved
	// if custom.
	filterUsingOptionsFunc func(definition *core.FilterDef,
		fo *pref.FilterOptions,
	) (core.TraverseFilter, error)
)
