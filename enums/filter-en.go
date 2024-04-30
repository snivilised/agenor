package enums

import (
	"math"
	"strings"

	"github.com/samber/lo"
	"github.com/snivilised/extendio/collections"
)

// FilterScope allows client to define which node types should be filtered.
// Filters can be applied to multiple node types by bitwise or-ing the XXXNodes
// definitions. A node may have multiple scope designations, eg a node may be top
// level and leaf if the top level directory does not itself contain further
// sub-directories thereby making it also a leaf.
// It should be noted a file is only a leaf node all of its siblings are all files
// only
type FilterScope uint32

const (
	ScopeUndefined FilterScope = 0

	// ScopeRoot, the Root scope
	//
	ScopeRoot FilterScope = 1 << (iota - 1)

	// ScopeTop, any node that is a direct descendent of the root node
	//
	ScopeTop

	// ScopeLeaf, for directories, any node that has no sub folders. For files, any node
	// that appears under a leaf directory node
	//
	ScopeLeaf

	// ScopeIntermediate, apply filter to nodes which are neither leaf or top nodes
	//
	ScopeIntermediate

	// ScopeFile attributed to file nodes
	//
	ScopeFile

	// ScopeFolder attributed to directory nodes
	//
	ScopeFolder

	// ScopeCustom, client defined categorisation (yet to be confirmed)
	//
	ScopeCustom

	// ScopeAll represents any node type
	//
	ScopeAll = math.MaxUint32
)

// String converts enum value to a string
func (f FilterScope) String() string {
	result := make([]string, 0, len(filterScopeKeys))

	for _, en := range filterScopeKeys {
		if en == ScopeAll {
			continue
		}

		if (en & f) > 0 {
			result = append(result, filterScopes[en])
		}
	}

	return strings.Join(result, "|")
}

// Set sets the bit position indicated by mask
func (f *FilterScope) Set(mask FilterScope) {
	*f |= mask
}

// Clear clears the bit position indicated by mask
func (f *FilterScope) Clear(mask FilterScope) {
	*f &^= mask
}

type FilterType uint

const (
	FilterTypeUndefined FilterType = iota

	// FilterTypeExtendedGlob is the preferred filter type as it the most
	// user friendly. The base part of the name is filtered by a glob
	// and the suffix is filtered by a list of defined extensions. The pattern
	// for the extended filter type is composed of 2 parts; the first is a
	// glob, which is applied to the base part of the name. The second part
	// is a csv of required extensions to filter for. The pattern is specified
	// in the form: "<base-glob>|ext1,ext2...". Each extension may include a
	// a leading dot. An example pattern definition would be:
	// "cover.*|.jpg,jpeg"
	FilterTypeExtendedGlob

	// FilterTypeRegex regex filter
	FilterTypeRegex

	// FilterTypeGlob glob filter
	FilterTypeGlob

	// FilterTypeCustom client definable filter
	FilterTypeCustom

	// FilterTypePoly poly filter
	FilterTypePoly
)

type allOrderedFilterScopes collections.OrderedKeysMap[FilterScope, string]

var (
	filterScopes    allOrderedFilterScopes
	filterScopeKeys []FilterScope
)

func init() {
	filterScopes = allOrderedFilterScopes{
		ScopeUndefined:    "Undefined",
		ScopeRoot:         "Root",
		ScopeTop:          "Top",
		ScopeLeaf:         "Leaf",
		ScopeIntermediate: "Intermediate",
		ScopeFile:         "File",
		ScopeFolder:       "Folder",
		ScopeCustom:       "Custom",
		ScopeAll:          "All",
	}

	filterScopeKeys = lo.Keys(filterScopes)
}
