package pref

import (
	"github.com/snivilised/agenor/core"
)

type (
	FilterReply struct {
		Node    core.TraverseFilter
		Child   core.ChildTraverseFilter
		Sampler core.SampleTraverseFilter
	}

	// FilteringSink is represents the callback function a client
	// can provide to enable them to receive the filter that has been
	// created from the definition specified.
	FilteringSink func(reply FilterReply)

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

func (fo FilterOptions) IsNodeFilteringActive() bool {
	return (fo.Node != nil) &&
		((fo.Node.Pattern != "") || fo.Node.Poly != nil)
}

func (fo FilterOptions) IsChildFilteringActive() bool {
	return (fo.Child != nil) && (fo.Child.Pattern != "")
}

func (fo FilterOptions) IsSampleFilteringActive() bool {
	return fo.Sample != nil
}

func (fo FilterOptions) IsFilteringActive() bool {
	return fo.IsNodeFilteringActive() ||
		fo.IsChildFilteringActive() ||
		fo.IsSampleFilteringActive() ||
		fo.IsCustomFilteringActive()
}

func (fo FilterOptions) IsCustomFilteringActive() bool {
	return fo.Custom != nil
}
