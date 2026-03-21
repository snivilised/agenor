package filtering

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/pref"
)

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
