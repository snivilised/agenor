package pref

import (
	"github.com/snivilised/traverse/core"
)

type (

	// EnoughAlready when using the universal navigator within a sampling operation, set
	// these accordingly from inside the while predicate to indicate wether the iteration
	// loop should continue to consider more entries to be included in the sample. So
	// set Files/Folders flags to true, when enough of those items have been included.
	EnoughAlready struct {
		Files   bool
		Folders bool
	}

	// FilteredInfo used within the sampling process during a traversal; more specifically,
	// they should be set inside the while predicate. Note, the Enough field is only
	// appropriate when using the universal navigator.
	FilteredInfo struct {
		Counts EntryQuantities
		Enough EnoughAlready
	}

	// EachDirectoryEntryPredicate callback to invoke for each child node event
	EachDirectoryEntryPredicate func(node *core.Node) bool

	// WhileDirectoryPredicate determines when to terminate the loop
	WhileDirectoryPredicate func(fi *FilteredInfo) bool

	SamplingIterationOptions struct {
		// Each enables customisation of the sampling functionality, instead of using
		// the defined filter. A directory's contents is sampled according to this
		// function. The function receives the TraverseItem being considered and should
		// return true to include in the sample, false otherwise.
		Each EachDirectoryEntryPredicate

		// While enables customisation of the sampling functionality, instead of using
		// the defined filter. The sampling loop will continue to run while this
		// condition is true. The predicate function should return false once condition
		// has been met to complete the sample. Of course, the loop will only run while
		// there are still remaining items to consider (ie there are no more entries
		// to consider for the current traverse item).
		While WhileDirectoryPredicate
	}

	SamplerOptions struct {
		// Iteration allows the client to customise how a directory's contents are sampled.
		// The default way to sample is either by slicing the directory's contents or
		// by using the filter to select either the first/last n entries (using the
		// SamplingOptions). If the client requires an alternative way of creating a
		// sample, eg to take all files greater than a certain size, then this
		// can be achieved by specifying Each and While inside Iteration.
		Iteration SamplingIterationOptions
	}
)

func WithSamplerOptions(so *SamplerOptions) Option {
	return func(o *Options) error {
		o.Sampler = *so

		return nil
	}
}
