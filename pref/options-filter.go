package pref

import (
	"github.com/snivilised/traverse/core"
)

type (
	FilterReply struct {
		Node  core.TraverseFilter
		Child core.ChildTraverseFilter
	}

	// FilteringSink is represents the callback function a client
	// can provide to enable them to receive the filter that has been
	// created from the definition specified.
	FilteringSink func(reply FilterReply)

	FilteringOptions struct {
		// FilterSink allows client access to the filter that is derived from the
		// filter definition
		//
		FilterSink FilteringSink

		// Custom client define-able filter. When restoring for resume feature,
		// its the client's responsibility to restore this themselves (see
		// PersistenceRestorer)
		Custom core.TraverseFilter
	}

	FilterOptions struct {
		// Node filter definitions that applies to the current file system node
		//
		Node *core.FilterDef

		// Child denotes the Child filter that is applied to the files which
		// are direct descendants of the current directory node being visited.
		//
		Child *core.ChildFilterDef
	}
)

func WithFilter(filter *FilterOptions) Option {
	return func(o *Options) error {
		o.Core.Filter = *filter

		return nil
	}
}

func WithFilterSink(sink FilteringSink) Option {
	return func(o *Options) error {
		o.Filtering.FilterSink = sink

		return nil
	}
}

func WithFilterCustom(filter core.TraverseFilter) Option {
	return func(o *Options) error {
		o.Filtering.Custom = filter

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

func (fo FilterOptions) IsFilteringActive() bool {
	return fo.IsNodeFilteringActive() || (fo.IsChildFilteringActive())
}

func (fo FilteringOptions) IsCustomFilteringActive() bool {
	return fo.Custom != nil
}

func IsFilteringActive(fo FilterOptions, fog FilteringOptions) bool {
	return fo.IsNodeFilteringActive() || fog.IsCustomFilteringActive()
}

func ResolveFilter(node core.TraverseFilter, fog FilteringOptions) core.TraverseFilter {
	if node != nil {
		return node
	}

	return fog.Custom
}
