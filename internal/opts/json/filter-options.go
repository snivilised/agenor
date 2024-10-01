package json

import (
	"github.com/snivilised/traverse/enums"
)

type (
	PolyFilterDef struct {
		File   FilterDef
		Folder FilterDef
	}

	FilterDef struct {
		// Type specifies the type of filter (mandatory)
		Type enums.FilterType `json:"filter-type"`

		// Description describes filter (optional)
		Description string `json:"filter-description"`

		// Pattern filter definition (mandatory)
		Pattern string `json:"pattern"`

		// Scope which file system entries this filter applies to (defaults
		// to ScopeAllEn)
		Scope enums.FilterScope `json:"filter-scope"`

		// Negate, reverses the applicability of the filter (Defaults to false)
		Negate bool `json:"negate"`

		// IfNotApplicable, when the filter does not apply to a directory entry,
		// this value determines whether the callback is invoked for this entry
		// or not (defaults to true).
		IfNotApplicable enums.TriStateBool `json:"if-not-applicable"`

		// Poly allows for the definition of a PolyFilter which contains separate
		// filters that target files and folders separately. If present, then
		// all other fields are redundant, since the filter definitions inside
		// Poly should be referred to instead.
		Poly *PolyFilterDef
	}

	ChildFilterDef struct {
		// Type specifies the type of filter (mandatory)
		Type enums.FilterType `json:"child-filter-type"`

		// Description describes filter (optional)
		Description string `json:"child-filter-description"`

		// Pattern filter definition (mandatory)
		Pattern string `json:"child-pattern"`

		// Negate, reverses the applicability of the filter (Defaults to false)
		Negate bool `json:"negate"`
	}

	SampleFilterDef struct {
		// Type specifies the type of filter (mandatory)
		Type enums.FilterType `json:"sample-filter-type"`

		// Description describes filter (optional)
		Description string `json:"sample-description"`

		// Pattern filter definition (mandatory except if using Custom)
		Pattern string `json:"sample-filter"`

		// Scope which file system entries this filter applies to;
		// for sampling, only ScopeFile and ScopeFolder are valid.
		Scope enums.FilterScope `json:"sample-filter-scope"`

		// Negate, reverses the applicability of the filter (Defaults to false)
		Negate bool `json:"negate"`

		// Poly allows for the definition of a PolyFilter which contains separate
		// filters that target files and folders separately. If present, then
		// all other fields are redundant, since the filter definitions inside
		// Poly should be referred to instead.
		Poly *PolyFilterDef
	}

	FilterOptions struct {
		// Node filter definitions that applies to the current file system node
		//
		Node *FilterDef

		// Child denotes the Child filter that is applied to the files which
		// are direct descendants of the current directory node being visited.
		//
		Child *ChildFilterDef

		// Sample is the filter used for sampling
		//
		Sample *SampleFilterDef
	}
)
