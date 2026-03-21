package jason

import (
	"github.com/snivilised/jaywalk/src/agenor/enums"
)

type (
	// PolyFilterDef allows for the definition of separate filters that target files and
	// directories separately.
	PolyFilterDef struct {
		// File filter definition (mandatory)
		File FilterDef

		// Directory filter definition (mandatory)
		Directory FilterDef
	}

	// FilterDef defines a filter that can be applied to a file system entry.
	// The filter definition includes the type of filter, a pattern that defines
	// the filter, and other optional fields that provide additional information
	// about the filter. The filter can be applied to either files or directories,
	// depending on the specified scope.
	// If the Poly field is present, it allows for the definition of a PolyFilter
	// which contains separate filters that target files and directories separately.
	// In this case, all other fields in FilterDef are redundant, since the filter
	// definitions inside Poly should be referred to instead.
	// The ChildFilterDef is a simplified version of FilterDef that is used for defining
	// filters that apply to the child entries of a directory. It does not include
	// the Scope field, since it is implicitly understood that the filter applies
	// to the child entries of a directory.
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
		// filters that target files and directories separately. If present, then
		// all other fields are redundant, since the filter definitions inside
		// Poly should be referred to instead.
		Poly *PolyFilterDef
	}

	// ChildFilterDef defines the filter that should be applied to the children
	// of the current node.
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

	// The SampleFilterDef is similar to FilterDef but is specifically used for
	// defining filters that are applied during sampling. It includes a Scope field
	// that specifies whether the filter applies to files or directories, as well
	// as an IfNotApplicable field that determines whether the callback is invoked
	// for entries that do not match the filter criteria.
	// The FilterOptions struct is a container for the different types of filters that can be
	// applied during navigation. It includes a Node filter that applies to the current
	// file system node, a Child filter that applies to the direct descendants of a
	// directory node, and a Sample filter that is used for sampling.
	// The design of these filter definitions allows for flexible and powerful filtering
	// capabilities during file system navigation, enabling users to specify complex criteria
	// for which entries should be processed or ignored.
	SampleFilterDef struct {
		// Type specifies the type of filter (mandatory)
		Type enums.FilterType `json:"sample-filter-type"`

		// Description describes filter (optional)
		Description string `json:"sample-description"`

		// Pattern filter definition (mandatory except if using Custom)
		Pattern string `json:"sample-filter"`

		// Scope which file system entries this filter applies to;
		// for sampling, only ScopeFile and ScopeDirectory are valid.
		Scope enums.FilterScope `json:"sample-filter-scope"`

		// Negate reverses the applicability of the filter (Defaults to false)
		Negate bool `json:"negate"`

		// Poly allows for the definition of a PolyFilter which contains separate
		// filters that target files and directories separately. If present, then
		// all other fields are redundant, since the filter definitions inside
		// Poly should be referred to instead.
		Poly *PolyFilterDef
	}

	// FilterOptions is a container for the different types of filters that can be
	// applied during navigation. It includes a Node filter that applies to the current
	// file system node, a Child filter that applies to the direct descendants of a
	// directory node, and a Sample filter that is used for sampling.
	// The design of these filter definitions allows for flexible and powerful filtering
	// capabilities during file system navigation, enabling users to specify complex criteria
	// for which entries should be processed or ignored.
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
