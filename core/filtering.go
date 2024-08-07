package core

import (
	"io/fs"

	"github.com/snivilised/traverse/enums"
)

// TraverseFilter filter that can be applied to file system entries. When specified,
// the callback will only be invoked for file system nodes that pass the filter.
type (
	TraverseFilter interface {
		// Description describes filter
		Description() string

		// Validate ensures the filter definition is valid, panics when invalid
		Validate()

		// Source, filter definition (comes from filter definition Pattern)
		Source() string

		// IsMatch does this node match the filter
		IsMatch(node *Node) bool

		// IsApplicable is this filter applicable to this node's scope
		IsApplicable(node *Node) bool

		// Scope, what items this filter applies to
		Scope() enums.FilterScope
	}

	FilterDef struct {
		// Type specifies the type of filter (mandatory)
		Type enums.FilterType

		// Description describes filter (optional)
		Description string

		// Pattern filter definition (mandatory)
		Pattern string

		// Scope which file system entries this filter applies to (defaults
		// to ScopeAllEn)
		Scope enums.FilterScope

		// Negate, reverses the applicability of the filter (Defaults to false)
		Negate bool

		// IfNotApplicable, when the filter does not apply to a directory entry,
		// this value determines whether the callback is invoked for this entry
		// or not (defaults to true).
		IfNotApplicable enums.TriStateBool

		// Poly allows for the definition of a PolyFilter which contains separate
		// filters that target files and folders separately. If present, then
		// all other fields are redundant, since the filter definitions inside
		// Poly should be referred to instead.
		Poly *PolyFilterDef
	}

	PolyFilterDef struct {
		File   FilterDef
		Folder FilterDef
	}

	// ChildTraverseFilter filter that can be applied to a folder's collection of entries
	// when subscription is

	ChildTraverseFilter interface {
		// Description describes filter
		Description() string

		// Validate ensures the filter definition is valid, panics when invalid
		Validate()

		// Source, filter definition (comes from filter definition Pattern)
		Source() string

		// Matching returns the collection of files contained within this
		// item's folder that matches this filter.
		Matching(children []fs.DirEntry) []fs.DirEntry
	}

	ChildFilterDef struct {
		// Type specifies the type of filter (mandatory)
		Type enums.FilterType

		// Description describes filter (optional)
		Description string

		// Pattern filter definition (mandatory)
		Pattern string

		// Negate, reverses the applicability of the filter (Defaults to false)
		Negate bool
	}

	SampleFilterDef struct {
		// Type specifies the type of filter (mandatory)
		Type enums.FilterType

		// Description describes filter (optional)
		Description string

		// Pattern filter definition (mandatory except if using Custom)
		Pattern string

		// Scope which file system entries this filter applies to;
		// for sampling, only ScopeFile and ScopeFolder are valid.
		Scope enums.FilterScope

		// Negate, reverses the applicability of the filter (Defaults to false)
		Negate bool

		// Poly allows for the definition of a PolyFilter which contains separate
		// filters that target files and folders separately. If present, then
		// all other fields are redundant, since the filter definitions inside
		// Poly should be referred to instead.
		Poly *PolyFilterDef

		// Custom client defined sampling filter
		//
		Custom SampleTraverseFilter
	}
	SampleTraverseFilter interface {
		// Description describes filter
		Description() string

		// Validate ensures the filter definition is valid, panics when invalid
		Validate()

		// Matching returns the collection of files contained within this
		// item's folder that matches this filter.
		Matching(children []fs.DirEntry) []fs.DirEntry
	}

	compoundCounters struct {
		filteredIn  uint
		filteredOut uint
	}
)

var BenignNodeFilterDef = FilterDef{
	Type:        enums.FilterTypeRegex,
	Description: "benign allow all",
	Pattern:     ".",
	Scope:       enums.ScopeRoot,
}
