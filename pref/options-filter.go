package pref

import (
	"github.com/snivilised/agenor/core"
)

type (

	// FilterReply represents the filter that is derived from the
	// filter definition specified by the client.
	FilterReply struct {
		// Node is the filter that is derived from the filter definition
		// specified by the client.
		Node core.TraverseFilter

		// Child is the filter that is applied to the files which are direct
		// descendants of the current directory node being visited.
		Child core.ChildTraverseFilter

		// Sample is the filter used for sampling
		Sampler core.SampleTraverseFilter
	}

	// FilteringSink is represents the callback function a client
	// can provide to enable them to receive the filter that has been
	// created from the definition specified.
	FilteringSink func(reply FilterReply)

	// FilterOptions represents the options for filtering that can be
	// specified by the client. These options are used to determine which file
	// system nodes (files or directories) the client defined handler is invoked
	// for. Note that the filter does not determine navigation, it only determines
	// wether the callback is invoked.
	FilterOptions struct {
		// Node filter definitions that applies to the current file system node
		//
		Node *core.FilterDef

		// Child denotes the Child filter that is applied to the files which
		// are direct descendants of the current directory node being visited.
		//
		Child *core.ChildFilterDef

		// Sample is the filter used for sampling
		//
		Sample *core.SampleFilterDef

		// Sink allows client access to the filter that is derived from the
		// filter definition
		//
		Sink FilteringSink

		// Custom client define-able filter. When restoring for resume feature,
		// its the client's responsibility to restore this themselves (see
		// PersistenceRestorer)
		Custom core.TraverseFilter
	}
)

// WithFilter used to determine which file system nodes (files or directories)
// the client defined handler is invoked for. Note that the filter does not
// determine navigation, it only determines wether the callback is invoked.
func WithFilter(filter *FilterOptions) Option {
	return func(o *Options) error {
		o.Filter = *filter

		return nil
	}
}

// IsNodeFilteringActive returns true if the filter options
// contain a node filter definition.
func (fo FilterOptions) IsNodeFilteringActive() bool {
	return (fo.Node != nil) &&
		((fo.Node.Pattern != "") || fo.Node.Poly != nil)
}

// IsChildFilteringActive returns true if the filter options
// contain a child filter definition.
func (fo FilterOptions) IsChildFilteringActive() bool {
	return (fo.Child != nil) && (fo.Child.Pattern != "")
}

// IsSampleFilteringActive returns true if the filter options
// contain a sample filter definition.
func (fo FilterOptions) IsSampleFilteringActive() bool {
	return fo.Sample != nil
}

// IsFilteringActive returns true if any of the filter options are active.
func (fo FilterOptions) IsFilteringActive() bool {
	return fo.IsNodeFilteringActive() ||
		fo.IsChildFilteringActive() ||
		fo.IsSampleFilteringActive() ||
		fo.IsCustomFilteringActive()
}

// IsCustomFilteringActive returns true if the filter options contain
// a custom filter.
func (fo FilterOptions) IsCustomFilteringActive() bool {
	return fo.Custom != nil
}
